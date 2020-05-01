package rbac

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// MakePetition handler
func (ur *Request) MakePetition() (interface{}, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(ur.Method, ur.Endpoint, bytes.NewBuffer(ur.Body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "pwa.bedu.org")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		return nil, errors.New("Bad Request")
	}

	var response interface{}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
