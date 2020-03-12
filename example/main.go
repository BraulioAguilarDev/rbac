package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	rbac "github.com/braulioinf/rbac"
)

var (
	FIREBASE_CONFIG string
	VAULT_API       string
	VAULT_USERNAME  string
	VAULT_PASSWORD  string
	ROLE_API        string
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}

	FIREBASE_CONFIG = fmt.Sprintf("firebase-admin.%v.json", os.Getenv("ENVIRONMENT"))
	VAULT_API = os.Getenv("VAULT_API")
	VAULT_USERNAME = os.Getenv("VAULT_USERNAME")
	VAULT_PASSWORD = os.Getenv("VAULT_PASSWORD")
	ROLE_API = os.Getenv("ROLE_API")
}

// GetProduct -> GET: https://localhost:5000/product
func GetProduct(w http.ResponseWriter, r *http.Request) {
	data := []string{"Go", "Js", "PyThon"}
	products, _ := json.Marshal(data)
	w.Write([]byte(products))
}

// PostProduct -> POST: http://localhost:5000/product
func PostProduct(w http.ResponseWriter, r *http.Request) {
	data := []string{"React"}
	new, _ := json.Marshal(data)
	w.Write([]byte(new))
}

// GetUser -> GET: http://localhost:5000/user
func GetUser(w http.ResponseWriter, r *http.Request) {
	data := []string{"Juan", "Pedro", "David"}
	users, _ := json.Marshal(data)
	w.Write([]byte(users))
}

// Run server $ go run main.go
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/product", PostProduct).Methods("POST")
	r.HandleFunc("/v1/product", GetProduct)

	r.HandleFunc("/v1/user", GetUser)

	rbac := rbac.RBAC{
		Config: &rbac.Config{
			VaultAPI: VAULT_API,
			Username: VAULT_USERNAME,
			Password: VAULT_PASSWORD,
			Firebase: FIREBASE_CONFIG,
			RoleAPI:  ROLE_API,
		},
	}
	if err := rbac.Initialize(); err != nil {
		fmt.Println(err.Error())
	}

	r.Use(rbac.Authorizer())

	if err := http.ListenAndServe(":5000", r); err != nil {
		fmt.Println(err)
	}
}
