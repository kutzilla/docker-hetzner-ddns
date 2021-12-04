package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Request(httpMethod string, url url.URL, headers map[string]string, body []byte) []byte {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(httpMethod, url.String(), bytes.NewBuffer(body))

	if err != nil {
		fmt.Println("Failure : ", err)
		return []byte{}
	}

	// Headers
	for key, element := range headers {
		req.Header.Add(key, element)
	}

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failure : ", err)
		return []byte{}
	}

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)

	return respBody
}
