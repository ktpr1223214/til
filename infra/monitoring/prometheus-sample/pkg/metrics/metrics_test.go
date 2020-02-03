package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPrometheus_Metrics(t *testing.T) {
	resp := requestGet(promhttp.Handler(), "/test", nil)
	bBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body %s", err)
	}

	t.Log("Prometheus で独自メトリクスをまだ利用していない")
	{
		if strings.Contains(string(bBody), "http_requests_total") {
			t.Errorf("want %s not contained in response body, but got", "http_requests_total")
		}
		if strings.Contains(string(bBody), "http_requests_errors_total") {
			t.Errorf("want %s not contained in response body, but got", "http_requests_errors_total")
		}
		if strings.Contains(string(bBody), "http_requests_duration_seconds_bucket") {
			t.Errorf("want %s not contained in response body, but got", "http_requests_duration_seconds_bucket")
		}
	}

	p := New()
	p.ObserveHTTPRequest("200", "GET", "/test", time.Now())
	p.ObserveHTTPError("200", "GET", "/test")

	resp = requestGet(promhttp.Handler(), "/test", nil)
	bBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body %s", err)
	}

	t.Log("Prometheus で独自メトリクスを収集した")
	{
		if !strings.Contains(string(bBody), `http_requests_total{code="200",method="GET",path="/test"} 1`) {
			t.Errorf("want %s contained in response body, but not got", `http_requests_total{code="200",method="GET",path="/test"} 1`)
		}
		if !strings.Contains(string(bBody), `http_requests_errors_total{code="200",method="GET",path="/test"} 1`) {
			t.Errorf("want %s contained in response body, but not got", `http_requests_errors_total{code="200",method="GET",path="/test"} 1`)
		}
		if !strings.Contains(string(bBody), `http_requests_duration_seconds_bucket{code="200",method="GET",path="/test",le="0.005"} 1`) {
			t.Errorf("want %s contained in response body, but not got", `http_requests_duration_seconds_bucket{code="200",method="GET",path="/test",le="0.005"} 1`)
		}
	}
}

func requestGet(srv http.Handler, path string, cookie *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if cookie != nil {
		req.Header.Add("Cookie", fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}
