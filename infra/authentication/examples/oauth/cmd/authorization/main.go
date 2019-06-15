package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	port                         = 9001
	userAgent                    = "oauth-server-1"
	authorizationURL             = "http://localhost:9001"
	authorizationApproveEndpoint = "http://localhost:9001/approve"
)

// RandomString 検証したいだけなので、実装は適当
func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// ClientInfo クライアントに関する情報
// この内容は本来、事前に何かしらの方法で共有されるべき情報だが、今回は固定としてコード中にハードコーディングしてある
type ClientInfo struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

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

type AuthServer struct {
	router     *http.ServeMux
	httpClient *http.Client

	userAgent        string
	authorizationURL *url.URL

	clients map[string]ClientInfo
	// reqID- の保存場所
	requests map[string]url.Values
}

// newPostRequest 認可サーバに対するリクエストを行う
// x-www-form-urlencoded 形式の POST
func (a *AuthServer) newPostRequest(urlPath string, data url.Values) (*http.Request, error) {
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

// do
func (a *AuthServer) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
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
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("request failed(status code: %d status: %s)", resp.StatusCode, resp.Status)
	}
	// v == nil で post するだけもある
	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, errors.Wrap(err, "failed to decode response Body")
}

func (a *AuthServer) auth() http.HandlerFunc {
	type AuthOption struct {
		ReqID   string `url:"req_id"`
		Approve string `url:"approve"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		queryParam := r.URL.Query()
		clientID := queryParam.Get("client_id")
		// 対象の client id が存在するか
		c, ok := a.clients[clientID]
		if !ok {
			log.Println("[Error] Failed to find client")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// 対象の redirect uri が登録済みのものと一致するか
		if c.RedirectURI != queryParam.Get("redirect_uri") {
			log.Println("[Error] The request send invalid redirect uri")
			return
		}

		// リクエスト ID-クエリパラメータを保存しておく
		reqID := RandomString(8)
		a.requests[reqID] = r.URL.Query()

		// TODO: 本来ここのステップは、ブラウザ経由でリソース権限保有者の承認ステップを挟むべき
		opt := AuthOption{
			ReqID:   reqID,
			Approve: "Approve",
		}
		urlEncodedBody, err := query.Values(opt)
		if err != nil {
			log.Printf("[Error] Failed to set query parameters %s", err)
			return
		}
		req, err := a.newPostRequest(authorizationApproveEndpoint, urlEncodedBody)
		if err != nil {
			log.Printf("[Error] Failed to build request %s", err)
			return
		}

		if _, err := a.do(context.Background(), req, nil); err != nil {
			log.Printf("[Error] Failed to send request %s", err)
			return
		}
	}
}

func (a *AuthServer) approve() http.HandlerFunc {
	type ApproveOption struct {
		Code string `url:"code"`
		// TODO: impl State
		// State string `url:"state"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("[Error] Failed to parse form %s", err)
			return
		}
		reqID := r.FormValue("req_id")
		queryParam, ok := a.requests[reqID]
		if !ok {
			log.Println("[Error] Failed to find past request")
			return
		}
		delete(a.requests, reqID)

		if r.FormValue("approve") != "Approve" {
			// TODO: ここは適切にクライアントに返すような処理を実装すること
			log.Println("[Error] Failed to get approval")
			return
		}

		if queryParam.Get("response_type") == "code" {
			code := RandomString(8)
			opt := ApproveOption{
				Code: code,
				// TODO: impl state at client
				// State: queryParam.Get("state"),
			}

			urlPath, err := addOptions(queryParam.Get("redirect_uri"), opt)
			if err != nil {
				log.Printf("[Error] Failed to build authorization URL %s", err)
				return
			}
			fmt.Println(urlPath)

			http.Redirect(w, r, urlPath, 302)

		} else {
			// TODO: implemente
		}
	}
}

func (a *AuthServer) Routes() {
	a.router.HandleFunc("/authorize", a.auth())
	a.router.HandleFunc("/approve", a.approve())
}

func (a *AuthServer) Run(port int) error {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), a.router); err != nil {
		return err
	}
	return nil
}

func main() {
	au, _ := url.ParseRequestURI(authorizationURL)

	a := &AuthServer{
		router:     http.NewServeMux(),
		httpClient: http.DefaultClient,
		clients: map[string]ClientInfo{
			"oauth-client-1": {
				ClientID:     "oauth-client-1",
				ClientSecret: "oauth-client-secret-1",
				RedirectURI:  "http://localhost:9000/callback",
			},
		},
		requests:         map[string]url.Values{},
		userAgent:        userAgent,
		authorizationURL: au,
	}
	a.Routes()

	if err := a.Run(port); err != nil {
		log.Fatal(err)
	}
}
