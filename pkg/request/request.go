package request

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
)

func Request(httpMethod string, url url.URL, headers map[string]string, body []byte) ([]byte, error) {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(httpMethod, url.String(), bytes.NewBuffer(body))

	if err != nil {
		log.Println("Failure : ", err)
		return []byte{}, err
	}

	// Headers
	for key, element := range headers {
		req.Header.Add(key, element)
	}

	// Fetch Request
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Failure : ", err)
		return []byte{}, err
	}

	// Read Response Body
	respBody, _ := io.ReadAll(resp.Body)

	return respBody, nil
}
