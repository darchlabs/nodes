package instance

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Record struct {
	ID          string     `db:"id"`
	UserID      string     `db:"user_id"`
	Network     string     `db:"network"`
	Environment string     `db:"environment"`
	Name        string     `db:"name"`
	ServiceURL  string     `db:"service_url"`
	Artifacts   *Artifacts `json:"artifacts"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

type Artifacts struct {
	Deployments []string `json:"deployments"`
	Pods        []string `json:"pods"`
	Services    []string `json:"services"`
}

// Implement the database/sql.Scanner interface for Artifacts
func (a *Artifacts) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("Artifacts.Scan: argument is not []byte")
	}
	return json.Unmarshal(b, a)
}

// Implement the database/sql/driver.Valuer interface for Artifacts
func (a Artifacts) Value() (driver.Value, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return b, nil
}
