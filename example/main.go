package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	rbac "github.com/ExponentialEducation/go-rbac"
)

var (
	FIREBASE_CREDENTIALS string
	VAULT_API            string
	VAULT_USERNAME       string
	VAULT_PASSWORD       string
	ROLE_API             string
	PORT                 string
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}

	FIREBASE_CREDENTIALS = os.Getenv("FIREBASE_CREDENTIALS")
	VAULT_API = os.Getenv("VAULT_API")
	VAULT_USERNAME = os.Getenv("VAULT_USERNAME")
	VAULT_PASSWORD = os.Getenv("VAULT_PASSWORD")
	ROLE_API = os.Getenv("ROLE_API")
	PORT = os.Getenv("APP_PORT")
}

func sessions(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"path":     "/v1/data/lms/capture/sessions",
		"granteds": "admin,insdes",
		"mehods":   "POST",
	}

	data, _ := json.Marshal(p)
	w.Write([]byte(data))
}

func publish(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"path":     "/v1/data/lms/capture/programs/publish",
		"granteds": "admin,insdessr",
		"methods":  "GET",
	}
	data, _ := json.Marshal(p)
	w.Write([]byte(data))
}

func courses(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"path":     "/v1/data/lms/capture/courses",
		"granteds": "admin,insdes,qa,insdessr",
		"methods":  "GET,UPDATE",
	}

	data, _ := json.Marshal(p)
	w.Write([]byte(data))
}

// Run server $ go run main.go
func main() {
	fmt.Printf("Run example in: :%v\n", PORT)

	r := mux.NewRouter()
	r.HandleFunc("/v1/data/lms/capture/sessions", sessions).Methods("GET", "POST", "PUT")
	r.HandleFunc("/v1/data/lms/capture/programs/publish", publish).Methods("GET", "POST", "PUT")
	r.HandleFunc("/v1/data/lms/capture/courses", courses).Methods("GET")

	rbac := rbac.RBAC{
		Config: &rbac.Config{
			VaultAPI: VAULT_API,
			Username: VAULT_USERNAME,
			Password: VAULT_PASSWORD,
			Firebase: FIREBASE_CREDENTIALS,
			RoleAPI:  ROLE_API,
		},
	}
	if err := rbac.Initialize(); err != nil {
		fmt.Println(err.Error())
	}

	r.Use(rbac.Authorizer())

	if err := http.ListenAndServe(fmt.Sprintf(":%v", PORT), r); err != nil {
		fmt.Println(err)
	}
}
