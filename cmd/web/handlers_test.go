package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
  rr := httptest.NewRecorder()

  r, err := http.NewRequest("GET", "/", nil)
  if err != nil {
    t.Fatal(err)
  }

  // Call the ping handler function, passing in the httptest.ResponseRecorder and http.Request.
  ping(rr, r)

  rs := rr.Result()

  if rs.StatusCode != http.StatusOK {
    t.Errorf("want  %d; got %d", http.StatusOK, rs.StatusCode)
  }

  // ANd we can check that the response body written by tht ping handler is equals "OK".
  defer rs.Body.Close()
  body, err :=  ioutil.ReadAll(rs.Body)
  if err != nil {
    t.Fatal(err)
  }

  if string(body) != "OK"  {
    t.Errorf("want body to equal %q", "OK")
  }
}