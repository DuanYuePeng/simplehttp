package graphql

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type SomeResponse struct {
	Message string
}

func TestDoGraphQL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		reqBytes, _ := io.ReadAll(r.Body)
		reqStr := string(reqBytes)

		if reqStr != `{"query":"query($arg:String!){hero{name}}","variables":{"arg":"someArgument"}}` {
			t.Fatalf("bad request body: %s", reqStr)
		}

		reqAuthHeader := r.Header.Get("Authorization")
		if reqAuthHeader != "Bearer some-token" {
			t.Fatalf("bad request header: %s", reqAuthHeader)
		}

		fmt.Fprintln(w, `{"message":"Hello"}`)
	}))
	defer server.Close()

	respBody := new(SomeResponse)
	err := Query(server.URL, "query($arg:String!){hero{name}}", map[string]interface{}{
		"arg": "someArgument",
	}, map[string]string{
		"Authorization": "Bearer some-token",
	}, respBody)

	if err != nil {
		t.Fatalf("got error: %s", err.Error())
	}

	if respBody.Message != "Hello" {
		t.Fatalf("bad response body: %s", respBody.Message)
	}
}
