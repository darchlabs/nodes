package command

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

// Status represents the current command status
type Status string

func (s Status) String() string {
	return string(s)
}

// Command statuses
const (
	StatusIdle         Status = "idle"
	StatusBootingError Status = "error"
	StatusRunning      Status = "running"
	StatusCrashed      Status = "crashed"
	StatusUpdating     Status = "updating"
	StatusDone         Status = "succeed"
	StatusStopped      Status = "stopped"
)

// Command represents a os level command, which can also receive a logger file
// in order to dump the output to it.
type Command struct {
	Cmd    *exec.Cmd
	Finish chan<- error // sends command exit error provided by `exec.Cmd.Wait`

	Version string

	status   Status
	execName string
	execArgs []string

	exitError chan error

	mu sync.RWMutex // guards command status
}

// Clone clones the command by instantiate a new one with same attributes
// and returns it. This is handy if you need to restart the process, first
// you stop it, then clone it, then you start the new cloned process.
func (c *Command) Clone() *Command {
	cmd := New(c.execName, c.execArgs...)
	cmd.Version = c.Version
	cmd.Finish = c.Finish
	return cmd
}

// Updater knows how to update the codebase of a specific command codebase.
type Updater interface {
	Update() (newVersion string, err error)
}

// MarshalJSON implements the json interface
func (c *Command) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Slug    string `json:"slug"`
		Version string `json:"version"`
		Status  Status `json:"status"`
	}{
		Slug:    c.Slug(),
		Status:  c.Status(),
		Version: c.Version,
	})
}

// Slug combines the command execName and execArgs in order to return a verbose
// identifier.
func (c *Command) Slug() string {
	return fmt.Sprintf("%s %s", c.execName, strings.Join(c.execArgs, " "))
}

// SetStatus sets the command current status
func (c *Command) SetStatus(status Status) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.status = status
}

// Update uses the updater in order to update the code base and the command
// version. If no updater is found, it returns an error. Update function
// returns a boolean that indicate if the code was either updated or not.
// Knowing if the command was updated is important in order to decide if we
// need to restart it or not.
func (c *Command) Update(updater Updater) (updated bool, err error) {
	oldVersion := c.Version
	newVersion, err := updater.Update()
	if err != nil {
		return false, err
	}

	if newVersion != oldVersion {
		c.Version = newVersion
		updated = true
	}

	return updated, nil
}

// Status check the command's process state and returns a verbose status.
func (c *Command) Status() Status {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.status == StatusStopped {
		return c.status
	}

	if ps := c.Cmd.ProcessState; ps != nil {
		if ps.Success() {
			return StatusDone
		}

		return StatusCrashed
	}

	return c.status
}

// Stop stops the command and closes the log file if exists.
func (c *Command) Stop() error {
	if c.status == StatusStopped {
		return fmt.Errorf("commands: command %q is already stopped", c.Slug())
	}

	if c.status == StatusCrashed || c.status == StatusBootingError {
		return nil
	}

	if c.Cmd.Process == nil {
		log.Println("[DEBUG] Stopped command when c.Cmd.Process was nil")
		return nil
	}

	// the ProcessState only exists if either the process exited, or we called
	// Run or Wait functions.
	ps := c.Cmd.ProcessState
	if ps != nil && ps.Exited() {
		log.Println("[DEBUG] Stopped command when c.Cmd.ProcessState is set and process.Exited() is true")
		return nil
	}

	if err := c.Cmd.Process.Signal(os.Interrupt); err != nil {
		return err
	}

	c.SetStatus(StatusStopped)

	return <-c.exitError
}

// Wait only proxies the function call to the  os.Command.Wait function.
func (c *Command) Wait() error {
	return c.Cmd.Wait()
}

func (c *Command) streamOutputWithPrefix(prefix string, rc io.ReadCloser) {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", prefix, scanner.Text())
	}
}

func (c *Command) StreamOutput(id string) error {
	pipe, err := c.Cmd.StdoutPipe()
	if err != nil {
		return err
	}

	pipeErr, err := c.Cmd.StderrPipe()
	if err != nil {
		return err
	}

	go c.streamOutputWithPrefix(fmt.Sprintf("STDOUT/%s", id), pipe)
	go c.streamOutputWithPrefix(fmt.Sprintf("STDERR/%s", id), pipeErr)

	return nil
}

// Start starts the process and pipes the command's output to the log file.
// If at any point there is an error it also closes the file if exists.
func (c *Command) Start() error {
	if err := c.Cmd.Start(); err != nil {
		c.SetStatus(StatusBootingError)
		return err
	}

	go func() {
		err := c.Wait()
		if err != nil {
			fmt.Println("[ERROR] status changed to crashed due to", err.Error())
			c.SetStatus(StatusCrashed)
		}

		c.exitError <- err
		c.Finish <- err
	}()

	c.SetStatus(StatusRunning)

	return nil
}

// Success just proxies the function call to the command.ProcessState struct.
func (c *Command) Success() bool {
	return c.Cmd.ProcessState.Success()
}

// NewCommand returns an initalized command pointer.
func New(name string, args ...string) *Command {
	cmd := &Command{
		Cmd: exec.Command(name, args...),

		status:    StatusIdle,
		execName:  name,
		execArgs:  args,
		exitError: make(chan error, 1),
	}

	cmd.Cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	return cmd
}
