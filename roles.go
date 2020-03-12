package rbac

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

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
func (rbac *RBAC) GetRolesByUID() ([]string, error) {
	uid := rbac.Wrapper.UID

	req := &Request{
		Method:   "GET",
		Endpoint: fmt.Sprintf("%s/api/roles/%v", rbac.Config.RoleAPI, uid),
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
