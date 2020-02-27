package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	pka "github.com/braulioinf/pkgauth"
)

// GetProduct .
func GetProduct(w http.ResponseWriter, r *http.Request) {
	data := []string{"Go", "Js", "PyThon"}
	products, _ := json.Marshal(data)
	w.Write([]byte(products))
}

// PostProduct .
func PostProduct(w http.ResponseWriter, r *http.Request) {
	data := []string{"React"}
	new, _ := json.Marshal(data)
	w.Write([]byte(new))
}

// GetUser .
func GetUser(w http.ResponseWriter, r *http.Request) {
	data := []string{"Juan", "Pedro", "David"}
	users, _ := json.Marshal(data)
	w.Write([]byte(users))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/product", PostProduct).Methods("POST")
	r.HandleFunc("/product", GetProduct)

	r.HandleFunc("/user", GetUser)

	rbac := &pka.Auth{
		Model:  "./model.conf",
		Policy: "./policy.csv",
	}

	if err := rbac.NewEnforcer(); err != nil {
		fmt.Println(err.Error())
	}

	r.Use(rbac.Authorizer())

	if err := http.ListenAndServe(":5000", r); err != nil {
		fmt.Println(err)
	}
}
