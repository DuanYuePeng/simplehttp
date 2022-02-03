package graphql

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Send(url string, query string, variables map[string]interface{}, headers map[string]string, respBody interface{}) (err error) {
	reqBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")

	for key, element := range headers {
		req.Header.Add(key, element)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode >= 400 {
		err = errors.New(fmt.Sprintf("failed to search repos, got status code %d\n%s", resp.StatusCode, respBodyBytes))
		return
	}

	if err = json.Unmarshal(respBodyBytes, respBody); err != nil {
		return
	}

	return
}
