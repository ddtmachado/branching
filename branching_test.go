package branching_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ddtmachado/branching"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
)

func TestBuilder(t *testing.T) {
	testCases := []struct {
		desc      string
		condition string
		wantErr   bool
	}{
		{
			desc:      "invalid condition",
			condition: "HumHum",
			wantErr:   true,
		},
	}
	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			cfg := branching.CreateConfig()
			cfg.Condition = test.condition

			_, err := branching.New(
				context.Background(),
				http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}),
				cfg,
				"test-plugin")

			if err != nil && !test.wantErr {
				t.Fatal("unexpected error", err)
			}
		})
	}
}

func TestBranching(t *testing.T) {
	testCases := []struct {
		desc                  string
		condition             string
		requestModifier       func(req *http.Request)
		branchResponseHeaders map[string]string
		wantResponseHeaders   map[string]string
	}{
		{
			desc:                  "match alt chain",
			condition:             "Header[`Foo`].0 == `bar`",
			branchResponseHeaders: map[string]string{"chain": "foo"},
			wantResponseHeaders:   map[string]string{"chain": "foo"},
			requestModifier: func(req *http.Request) {
				req.Header.Add("foo", "bar")
			},
		},
		{
			desc:                  "no match alt chain",
			condition:             "Header[`Foo`].0 == `bar`",
			branchResponseHeaders: map[string]string{"chain": "foo"},
			wantResponseHeaders:   map[string]string{"chain": "default"},
			requestModifier: func(req *http.Request) {
				req.Header.Add("foo", "notbar")
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			cfg := branching.CreateConfig()
			cfg.Condition = test.condition
			cfg.Chain = map[string]*dynamic.Middleware{
				"test-chain": {
					Headers: &dynamic.Headers{CustomResponseHeaders: test.branchResponseHeaders},
				}}

			backend := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Header().Add("chain", "default")
				rw.WriteHeader(http.StatusOK)
				return
			})

			ctx := context.Background()
			plugin, err := branching.New(ctx, backend, cfg, "test-plugin")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			if err != nil {
				t.Fatal(err)
			}
			test.requestModifier(req)

			plugin.ServeHTTP(recorder, req)

			response := recorder.Result()
			if response.StatusCode != http.StatusOK {
				t.Fatal("failed response")
			}

			for hk, hv := range test.wantResponseHeaders {
				if response.Header.Get(hk) != hv {
					t.Fatal("unexpected chain header")
				}
			}
		})
	}
}
