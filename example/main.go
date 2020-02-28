package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	pka "github.com/braulioinf/pkgauth"
)

var FIREBASE_CREDENTIALS_JSON string

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}

	FIREBASE_CREDENTIALS_JSON = fmt.Sprintf("firebase-admin.%v.json", os.Getenv("ENVIRONMENT"))
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
	r.HandleFunc("/product", PostProduct).Methods("POST")
	r.HandleFunc("/product", GetProduct)

	r.HandleFunc("/user", GetUser)

	rbac := pka.NewRBAC("./model.conf", "./policy.csv", FIREBASE_CREDENTIALS_JSON)
	if err := rbac.Initialize(); err != nil {
		fmt.Println(err.Error())
	}

	r.Use(rbac.Authorizer())

	if err := http.ListenAndServe(":5000", r); err != nil {
		fmt.Println(err)
	}
}
