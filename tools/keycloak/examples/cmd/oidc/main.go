package main

import (
	"fmt"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"net/http"

	"context"
	"encoding/json"
	"log"
)

const (
	clientID     = "oidc-sample"
	clientSecret = "75b8d827-24d1-49ba-8c03-4769d1c893c9"
	redirectURL  = "http://localhost:8000/private/callback"
)

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
		Scopes: []string{oidc.ScopeOpenID},
	}

	// TODO: 適切に設定
	state := "foobar" // Don't do this in production.

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
	})

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

		oauth2Token.AccessToken = "*REDACTED*"

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
		w.Write(data)
	})

	// http.Handle("/hello", http.HandlerFunc(hello))
	http.ListenAndServe(":8000", nil)
	fmt.Println(provider)
}
