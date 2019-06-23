package main

import (
	"fmt"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"net/http"

	"context"
	"encoding/json"
	"log"
	"strings"
)

const (
	clientID     = "oidc-sample"
	clientSecret = "75b8d827-24d1-49ba-8c03-4769d1c893c9"
	redirectURL  = "http://localhost:8000/private/callback"
)

func OidcAuthMiddleware() {

}

func main() {
	// https://github.com/coreos/go-oidc/blob/v2/oidc.go#L97
	// wellKnown := strings.TrimSuffix(issuer, "/") + "/.well-known/openid-configuration" なので、
	provider, err := oidc.NewProvider(context.Background(), "http://localhost/auth/realms/Sample")
	if err != nil {
		log.Fatal(err)
	}
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)

	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		// keycloak 上で、Assigned Optional Client Scopes に入れる場合は、こうして明示的に scopes に入れる
		// default に入れている場合は、気にする必要はない
		Scopes: []string{oidc.ScopeOpenID, "good-service"},
	}

	// TODO: 適切に設定
	state := "foobar" // Don't do this in production.

	// TODO: OAuth による保護は、Middleware なりで書くべきかと
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Authorization を要求するのか何を要求するのかはよく考えるべきでは？
		rawAccessToken := r.Header.Get("Authorization")
		if rawAccessToken == "" {
			http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
			return
		}

		parts := strings.Split(rawAccessToken, " ")
		if len(parts) != 2 {
			w.WriteHeader(400)
			return
		}
		_, err := verifier.Verify(context.Background(), parts[1])
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
			return
		}

		w.Write([]byte("hello world"))
	})

	// ここでリロードするとエラーとなるが、それは認可コードが一度きりしか使えないため(exchange で失敗)
	http.HandleFunc("/private/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		oauth2Token, err := oauth2Config.Exchange(context.Background(), r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}
		idToken, err := verifier.Verify(context.Background(), rawIDToken)
		if err != nil {
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 編集済み or 裸
		// oauth2Token.AccessToken = "*REDACTED*"

		resp := struct {
			OAuth2Token   *oauth2.Token
			IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
		}{oauth2Token, new(json.RawMessage)}

		if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := json.MarshalIndent(resp, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// TODO: ここどうするのか考える
		w.Write(data)
		/*
		fmt.Println(data)
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oauth2Token.AccessToken))
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusFound)
		// http.Redirect(w, r, "/", http.StatusFound)
		*/
	})
	http.ListenAndServe(":8000", nil)
	fmt.Println(provider)
}
