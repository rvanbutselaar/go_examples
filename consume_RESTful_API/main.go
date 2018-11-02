package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

func main() {
	// fmt.Println("Starting the application...")
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	response, err := client.Get("https://release-monitoring.org/api/v2/projects/?name=docker")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Printf("The HTTP response is not 200: %#v\n", response.StatusCode)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// fmt.Println(string(data))

		value := gjson.GetBytes(data, "items.#.version")
		fmt.Printf("Version: %s\n", value.String())
	}
	// fmt.Println("Terminating the application...")
}
