package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApplication(t *testing.T) *application {
  return &application{
    errorLog: log.New(io.Discard, "", 0),
    infoLog: log.New(io.Discard, "", 0),
  }
}

// Define a custom testServer type which anonymously embeds a httptest.Server instance.
type testServer struct {
  *httptest.Server
}

// Creat a newTestServer helper which initializes and returns a new instance of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
  ts := httptest.NewTLSServer(h)
  return &testServer{ts}
}

// Implement a get method on our custom testServer type. This makes a GET request to ga given url path on the test server, and returns the response status code, headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
  rs, err := ts.Client().Get(ts.URL + urlPath)
  if err != nil {
    t.Fatal(err)
  }

  defer rs.Body.Close()
  body, err := io.ReadAll(rs.Body)
  if err != nil {
    t.Fatal(err)
  }

  return rs.StatusCode, rs.Header, body
}