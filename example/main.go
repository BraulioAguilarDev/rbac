package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	pka "github.com/braulioinf/pkgauth"
)

// GetNice .
func GetNice(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Nice XD"))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/example", GetNice)

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
