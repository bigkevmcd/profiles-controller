package httpapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	profilesv1alpha1 "github.com/bigkevmcd/profiles-controller/api/v1alpha1"
	"github.com/bigkevmcd/profiles-controller/controllers"
	"github.com/google/go-cmp/cmp"
)

func TestSearchProfiles(t *testing.T) {
	router := NewRouter(&controllers.InMemoryProfiles{})
	router.repository.Add(makeTestProfile("testing", "Just a test repository"))
	router.repository.Add(makeTestProfile("other", "Another repository"))
	ts := makeServer(t, router)
	req := makeClientRequest(t, ts, "/profiles", func(u *url.URL) {
		q := u.Query()
		q.Set("q", "test")
		u.RawQuery = q.Encode()
	})
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assertJSONResponse(t, res, map[string]interface{}{
		"profiles": []interface{}{
			map[string]interface{}{
				"description": "Just a test repository",
				"maturity":    "alpha",
				"name":        "testing",
				"publisher":   "WeaveWorks",
				"url":         "https://example.com/testing.git",
				"version":     "v0.0.1",
			},
		},
	})
}

func makeClientRequest(t *testing.T, ts *httptest.Server, path string, opts ...func(*url.URL)) *http.Request {
	r, err := http.NewRequest("GET", makeURL(t, ts, path, opts...), nil)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func makeURL(t *testing.T, ts *httptest.Server, path string, opts ...func(*url.URL)) string {
	t.Helper()
	parsed, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	parsed.Path = path
	for _, o := range opts {
		o(parsed)
	}
	return parsed.String()
}

func makeServer(t *testing.T, a *APIRouter) *httptest.Server {
	ts := httptest.NewTLSServer(a)
	t.Cleanup(ts.Close)
	return ts
}

func assertJSONResponse(t *testing.T, res *http.Response, want map[string]interface{}) {
	t.Helper()
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		errMsg, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		t.Fatalf("didn't get a successful response: %v (%s)", res.StatusCode, strings.TrimSpace(string(errMsg)))
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if h := res.Header.Get("Content-Type"); h != "application/json" {
		t.Fatalf("wanted 'application/json' got %s", h)
	}
	got := map[string]interface{}{}

	err = json.Unmarshal(b, &got)
	if err != nil {
		t.Fatalf("failed to parse %s: %s", b, err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("JSON response failed:\n%s", diff)
	}
}

func assertErrorResponse(t *testing.T, res *http.Response, status int, want string) {
	t.Helper()
	if res.StatusCode != status {
		defer res.Body.Close()
		errMsg, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		t.Fatalf("status code didn't match: %v (%s)", res.StatusCode, strings.TrimSpace(string(errMsg)))
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(b)); got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func makeTestProfile(name, description string) profilesv1alpha1.Profile {
	return profilesv1alpha1.Profile{
		Name:        name,
		URL:         "https://example.com/" + name + ".git",
		Version:     "v0.0.1",
		Description: description,
		Maturity:    "alpha",
		Publisher:   "WeaveWorks",
	}
}
