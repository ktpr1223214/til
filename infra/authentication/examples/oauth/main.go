package main

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
	"reflect"
)

// OAuth Client 側の情報
const (
	port     = 9001
	clientID = "oauth-client-1"
	// もちろん実際にこういう風に書いては駄目
	clientSecret = "oauth-client-secret-1"
	redirectURI  = "http://localhost:9000/callback"
)

// 認可サーバ側の構成情報
// 認可のやりとりを行うために、クライアントも知る必要がある
const (
	authorizationEndpoint = "http://localhost:9001/authorize"
	tokenEndpoint         = "http://localhost:9001/token"
)

func addOptions(s string, opt interface{}) (string, error) {
	// reflect.ValueOf(~).Kind() で変数に格納されている値の種類を得る
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}
	u, err := url.Parse(s)
	if err != nil {
		return s, errors.Wrapf(err, "failed to parse URL %s", s)
	}
	vs, err := query.Values(opt)
	if err != nil {
		return s, errors.Wrap(err, "failed to set query parameters")
	}
	u.RawQuery = vs.Encode()
	return u.String(), nil
}

// buildURL
// options: クエリパラメータの組
func buildURL(base string, options interface{}) (string, error) {
	// Parse parses rawurl into a URL structure
	/*u, err := url.Parse(base)
	if err != nil {
		return "", err
	}*/

	u, err := addOptions(base, options)
	if err != nil {
		return "", err
	}

	return u, nil
}

func buildRedirect(r *http.Request) {

}

type AuthClient struct {
	router *http.ServeMux
}

// Run server 実行
func (a *AuthClient) Run(port int) error {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), a.router); err != nil {
		return err
	}
	return nil
}

func (a *AuthClient) auth() http.HandlerFunc {
	type AuthOption struct {
		ResponseType string `url:"response_type"`
		ClientID     string `url:"client_id"`
		RedirectURI  string `url:"redirect_uri"`
	}
	opt := AuthOption{
		ResponseType: "code",
		ClientID:     clientID,
		RedirectURI:  redirectURI,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		authorizationURL, err := buildURL(authorizationEndpoint, opt)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, authorizationURL, 302)
	}
}

func (a *AuthClient) Routes() {
	a.router.HandleFunc("/authorize", a.auth())
}

func main() {
	a := &AuthClient{
		router: http.NewServeMux(),
	}
	a.Routes()

	if err := a.Run(port); err != nil {
		log.Fatal(err)
	}
}
