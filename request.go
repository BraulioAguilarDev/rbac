package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
)

// UIDRequest struct
type UIDRequest struct {
	UID      string
	Method   string
	Endpoint string
	Body     []byte
}

// Role struct
type Role struct {
	ID       int    `json:"-"`
	Name     string `json:"name"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}

// GetRoles .
func GetRoles(uid string) ([]string, error) {
	req := &UIDRequest{
		UID:      uid,
		Method:   "GET",
		Endpoint: fmt.Sprintf("http://localhost:8080/api/roles/%v", uid),
	}

	res, err := req.MakePetition()
	if err != nil {
		return nil, err
	}

	rolesData := new([]Role)

	if err := mapstructure.Decode(res, &rolesData); err != nil {
		fmt.Println(err.Error())
	}

	var roles []string
	for _, r := range *rolesData {
		roles = append(roles, r.Name)
	}

	return roles, nil
}

// MakePetition .
func (ur *UIDRequest) MakePetition() (interface{}, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(ur.Method, ur.Endpoint, bytes.NewBuffer(ur.Body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

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
