package rbac

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// Response struct
type Response struct {
	Data []Role `json:"data"`
}

// Role struct
type Role struct {
	ID       int    `json:"-"`
	Name     string `json:"name"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}

// Request struct
type Request struct {
	Method   string
	Endpoint string
	Body     []byte
}

// GetRolesByUID func
func (rbac *RBAC) GetRolesByUID(uid string) ([]string, error) {
	req := &Request{
		Method:   "GET",
		Endpoint: fmt.Sprintf("%s/api/roles/%v", rbac.RoleAPI, uid),
	}

	data, err := req.MakePetition()
	if err != nil {
		return nil, err
	}

	response := new(Response)

	if err := mapstructure.Decode(data, &response); err != nil {
		fmt.Println(err.Error())
	}

	var roles []string
	for _, r := range response.Data {
		roles = append(roles, r.Name)
	}

	return roles, nil
}
