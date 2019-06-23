package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// OAuth Client 側の情報
const (
	port     = 9000
	clientID = "oauth-client-1"
	// もちろん実際にこういう風に書いては駄目
	clientSecret = "oauth-client-secret-1"
	redirectURI  = "http://localhost:9000/callback"
	userAgent    = "oauth-client-1"
)

// 認可サーバ側の構成情報
// 認可のやりとりを行うために、クライアントも知る必要がある
const (
	authorizationURL      = "http://localhost:9001"
	authorizationEndpoint = "authorize"
	tokenEndpoint         = "token"
)

func addOptions(s string, opt interface{}) (string, error) {
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

type AuthClient struct {
	router     *http.ServeMux
	httpClient *http.Client

	userAgent        string
	authorizationURL *url.URL
}

// Run server 実行
func (a *AuthClient) Run(port int) error {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), a.router); err != nil {
		return err
	}
	return nil
}

// newJSONRequest 認可サーバに対するリクエストを行う
func (a *AuthClient) newJSONRequest(method string, urlPath string, body interface{}) (*http.Request, error) {
	u, err := a.authorizationURL.Parse(urlPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url %s", urlPath)
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, errors.Wrap(err, "failed to encode json body")
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new request")
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", a.userAgent)
	return req, nil
}

// newPostRequest 認可サーバに対するリクエストを行う
// x-www-form-urlencoded 形式の POST
func (a *AuthClient) newPostRequest(urlPath string, data url.Values) (*http.Request, error) {
	u, err := a.authorizationURL.Parse(urlPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", a.userAgent)
	return req, nil
}

// Do
func (a *AuthClient) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		// ctx.Done() は channel(<- chan)を返す
		select {
		case <-ctx.Done():
			return nil, errors.Wrap(ctx.Err(), "failed by context")
		default:
		}
		return nil, errors.Wrap(err, "failed to request")
	}
	defer resp.Body.Close()
	// TODO: ここは恐らくもっと拡張すべき
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("request failed(status code: %d status: %s)", resp.StatusCode, resp.Status)
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, errors.Wrap(err, "failed to decode response Body")
}

// 最初の認可リクエスト
// 認可サーバにリダイレクトを行う
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
		urlPath, err := addOptions(authorizationEndpoint, opt)
		if err != nil {
			log.Printf("[Error] Failed to build authorization URL %s", err)
			return
		}
		u, err := a.authorizationURL.Parse(urlPath)
		if err != nil {
			log.Printf("[Error] Failed to build authorization URL %s", err)
			return
		}

		http.Redirect(w, r, u.String(), 302)
	}
}

func (a *AuthClient) callback() http.HandlerFunc {
	type CallbackOption struct {
		GrantType   string `url:"grant_type"`
		Code        string `url:"code"`
		RedirectURI string `url:"redirect_uri"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		queryParam := r.URL.Query()

		urlEncodedBody, err := query.Values(CallbackOption{
			GrantType:   "authorization_code",
			Code:        queryParam.Get("code"),
			RedirectURI: redirectURI,
		})
		if err != nil {
			log.Printf("[Error] Failed to get url encodeed body %s", err)
			return
		}

		u, err := a.authorizationURL.Parse(tokenEndpoint)
		if err != nil {
			log.Printf("[Error] Failed to parse url %s", err)
			return
		}
		req, err := a.newPostRequest(u.String(), urlEncodedBody)
		if err != nil {
			log.Printf("[Error] Failed to build request %s", err)
			return
		}
		// TODO: basic auth
		// TODO: parse response body

		if _, err := a.do(context.Background(), req, nil); err != nil {
			log.Printf("[Error] Failed to send request %s", err)
			return
		}
	}
}

func (a *AuthClient) Routes() {
	a.router.HandleFunc("/authorize", a.auth())
	a.router.HandleFunc("/callback", a.callback())
}

func main() {
	au, _ := url.ParseRequestURI(authorizationURL)

	a := &AuthClient{
		router:           http.NewServeMux(),
		httpClient:       http.DefaultClient,
		userAgent:        userAgent,
		authorizationURL: au,
	}
	a.Routes()

	if err := a.Run(port); err != nil {
		log.Fatal(err)
	}
}
