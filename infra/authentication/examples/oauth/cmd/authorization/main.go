package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	port = 9001
)

type AuthServer struct {
	router     *http.ServeMux
	httpClient *http.Client
}

func (a *AuthServer) auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("[Error] Failed to read request body %s", err)
			return
		}
		fmt.Println(string(body))
	}
}

func (a *AuthServer) Routes() {
	a.router.HandleFunc("/authorize", a.auth())
}

func (a *AuthServer) Run(port int) error {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), a.router); err != nil {
		return err
	}
	return nil
}

func main() {
	a := &AuthServer{
		router:     http.NewServeMux(),
		httpClient: http.DefaultClient,
	}
	a.Routes()

	if err := a.Run(port); err != nil {
		log.Fatal(err)
	}
}
