package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	// The URL of the file to download
	url := "https://gist.githubusercontent.com/mtavano/d8e5f3d98019f8a1415748bd8326fba6/raw/4f2e508f6943b5186752fc96b13fae6b53419df6/kubeconf.mnode2.yml"

	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Create the output file
	out, err := os.Create("kubeconf.yml")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Copy the contents of the response body to the output file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	// Print a message indicating success
	println("File downloaded as kubeconf.yml")
}
