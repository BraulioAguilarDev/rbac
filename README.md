# RBAC Module

A middleware for request validation created with Go and Vault

## What do you need for this to work?

* Run [Vault Service](https://github.com/braulioinf/vault-poc)
* Run [Profile & Role Service](https://github.com/ExponentialEducation/roles-profiles-microservice)
* Run [Account Microservice](https://github.com/ExponentialEducation/account-microservice)

## Installation

`go get -u github.com/ExponentialEducation/go-rbac`

## Basic Usage

Middleware can be added to a router using `Router.Use()`:

```go
// package main

import "github.com/ExponentialEducation/go-rbac"

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/v1/example", handlerFunc).Methods("POST")

  // Config
  rbac := rbac.RBAC{
    Config: &rbac.Config{
      VaultAPI: "https://api.vault:8200",
      Username: "user",
      Password: "pass",
      Firebase: "firebase-admin.json",
      RoleAPI:  "https://api.roles.com",
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
```

## Full example

For more information [Example directory](https://github.com/ExponentialEducation/go-rbac/tree/develop/example)
