package main

import (
	"fmt"

	rbac "github.com/ExponentialEducation/go-rbac"
)

const (
	VAULT_API      = "https://vault.pitakill.net:8200"
	VAULT_USERNAME = "authorizer"
	VAULT_PASSWORD = "helloworld"

	APPROLE_LOGIN       = "auth/approle/login"
	LMS_PROGRAM         = "v1/data/lms/capture/programs"
	LMS_PROGRAM_PUBLISH = "v1/data/lms/capture/programs/publish"
	ACCOUNT_SERVICE     = "http://localhost:8080"
)

func main() {
	rbac := &rbac.RBAC{
		Config: &rbac.Config{
			VaultAPI: VAULT_API,
			Username: VAULT_USERNAME,
			Password: VAULT_PASSWORD,
			RoleAPI:  ACCOUNT_SERVICE,
		},
	}
	if err := rbac.Initialize(); err != nil {
		fmt.Printf("Error Firebase confirg, Detail: %s \n", err)
	}

	roles := []string{
		"qa",
		"insdes",
		"insdessr",
		"admin",
	}

	fmt.Printf("Use PATH: LMS PROGRAM => %s\n", LMS_PROGRAM)
	fmt.Printf("Use PATH: LMS PROGRAM PUBLISH => %s\n", LMS_PROGRAM_PUBLISH)
	fmt.Println(".........................................................")

	for _, role := range roles {
		if err := rbac.Wrapper.LoginAs(role); err != nil {
			fmt.Printf("LoginAs: %v\n", err.Error())
		}

		cantDelete := fmt.Sprintf("The role %q DOES NOT HAVE permissions to DELETE programs on the LMS\n", role)
		canDelete := fmt.Sprintf("The role %q HAS permissions to DELETE programs on the LMS\n", role)
		cantRead := fmt.Sprintf("The role %q DOES NOT HAVE permissions to READ programs on the LMS\n", role)
		canRead := fmt.Sprintf("The role %q HAS permissions to READ programs on the LMS\n", role)
		cantPublish := fmt.Sprintf("The role %q DOES NOT HAVE permissions to PUBLISH programs on the LMS\n", role)
		canPublish := fmt.Sprintf("The role %q HAS permissions to PUBLISH programs on the LMS\n", role)

		// DELETE
		if rbac.Wrapper.CanDelete(LMS_PROGRAM) {
			fmt.Printf(canDelete)
		} else {
			fmt.Printf(cantDelete)
		}

		// GET
		if rbac.Wrapper.CanRead(LMS_PROGRAM) {
			fmt.Printf(canRead)
		} else {
			fmt.Printf(cantRead)
		}

		// POST/PUT
		if rbac.Wrapper.CanWrite(LMS_PROGRAM_PUBLISH) {
			fmt.Printf(canPublish)
		} else {
			fmt.Printf(cantPublish)
		}
	}
}
