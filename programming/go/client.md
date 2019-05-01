---
title: Client
---

## Client

## 実装
### クライアントについて
* ほぼ必須で持つであろう要素は、以下あたり
    * http.Client
    * User Agent
    * URL
    * (logger)
* http.Client は外から渡せるようにしておくことに注意
    * timeout の設定など弄りたい部分が多い
    * 方式は、New~ で引数にする・functional option パターンで渡せるようにしておいて基本はデフォルト設定で初期化など

### よくある共通実装
* from go-github
``` go
// メソッド例
func (s *ActivityService) ListEvents(ctx context.Context, opt *ListOptions) ([]*Event, *Response, error) {
	u, err := addOptions("events", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*Event
	resp, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, resp, err
	}

	return events, resp, nil
}

// リクエスト
// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
    // https://github.com/google/go-github/pull/690#issuecomment-321879628
    // とか(自分なら / 付けて対応しそうだけど確かに書いてある通りな気がする)
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// cf. https://developer.github.com/v3/media/
	req.Header.Set("Accept", mediaTypeV3)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Do
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
    req = req.WithContext(ctx)
    resp, err := c.httpClient.Do(req)
    if err != nil {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }
        return nil, err
    }
    defer resp.Body.Close()
    err = json.NewDecoder(resp.Body).Decode(v)
    return resp, err
}

// Error
// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range or equal to 202 Accepted.
// API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other
// response body will be silently ignored.
//
// The error type will be *RateLimitError for rate limit exceeded errors,
// *AcceptedError for 202 Accepted status codes,
// and *TwoFactorAuthError for two-factor authentication errors.
func CheckResponse(r *http.Response) error {
	if r.StatusCode == http.StatusAccepted {
		return &AcceptedError{}
	}
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	switch {
	case r.StatusCode == http.StatusUnauthorized && strings.HasPrefix(r.Header.Get(headerOTP), "required"):
		return (*TwoFactorAuthError)(errorResponse)
	case r.StatusCode == http.StatusForbidden && r.Header.Get(headerRateRemaining) == "0" && strings.HasPrefix(errorResponse.Message, "API rate limit exceeded for "):
		return &RateLimitError{
			Rate:     parseRate(r),
			Response: errorResponse.Response,
			Message:  errorResponse.Message,
		}
	case r.StatusCode == http.StatusForbidden && strings.HasSuffix(errorResponse.DocumentationURL, "/v3/#abuse-rate-limits"):
		abuseRateLimitError := &AbuseRateLimitError{
			Response: errorResponse.Response,
			Message:  errorResponse.Message,
		}
		if v := r.Header["Retry-After"]; len(v) > 0 {
			// According to GitHub support, the "Retry-After" header value will be
			// an integer which represents the number of seconds that one should
			// wait before resuming making requests.
			retryAfterSeconds, _ := strconv.ParseInt(v[0], 10, 64) // Error handling is noop.
			retryAfter := time.Duration(retryAfterSeconds) * time.Second
			abuseRateLimitError.RetryAfter = &retryAfter
		}
		return abuseRateLimitError
	default:
		return errorResponse
	}
}

// オプション
// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
```

### メソッドが多い場合の対処
主に2つの対処が考えられる
1. ファイルを分けてることで、把握を容易にする
    * moby/moby 参照
2. Service に分ける
    * Service はクライアントのフィールド要素になり、一方で Service 側もクライアントをフィールドに持つ(じゃないと当然リクエストできない..)

``` go
// service
type service struct {
	client *Client
}

// Client のフィールド
{...
    common service
...
}

// クライアント生成関数内で
{...
    // ここの処理については、cf. https://github.com/google/go-github/issues/389
	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent, UploadURL: uploadURL}
	c.common.client = c
	c.Activity = (*ActivityService)(&c.common)
...
}
```

### テスト
* httptest が便利

## Reference
### クライアント一般
* [Writing a Go client for your RESTful API](https://medium.com/@marcus.olsson/writing-a-go-client-for-your-restful-api-c193a2f4998c)
* [GolangでAPI Clientを実装する](https://deeeet.com/writing/2016/11/01/go-api-client/)
* [Go, REST APIs, and Pointers](https://willnorris.com/2014/05/go-rest-apis-and-pointers)
    * go-github や、aws-sdk-go でよくみる pointer の使い方解説
* [google/go-github](https://github.com/google/go-github)
    * 良い例
* [moby/moby](https://github.com/moby/moby/tree/master/client)
    * 良い例

### context
* [How to correctly use context.Context in Go 1.7](https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39)
    * request scoped
    * If you are writing a library and your functions may block, it is a perfect use case for Context.
    * Try not to use context.Value
        * unknown-unknowns を増やすことになる

### test
* [Symmetric API Testing](https://blog.gopheracademy.com/advent-2015/symmetric-api-testing-in-go/)
