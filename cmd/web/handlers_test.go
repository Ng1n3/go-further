package main

import (
	"bytes"
	"io"
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
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestPing2(t *testing.T) {
	//? Create a new instance of our application struct.
	// app := &application{
	// 	errorLog: log.New(io.Discard, "", 0),
	// 	infoLog:  log.New(io.Discard, "", 0),
	// }

	//? Use http.NewTLSServer() to create a new test server, passing in the value returned by our app.routes() method as the handler for the server. This starts up a HTTPS server which listens on a randomly-chosen port of your local machine for the duration of the test.We use ts.Close() to shutdown the server when the test finishes.
	// ts := httptest.NewTLSServer(app.routes())
	// defer ts.Close()

	//? newwork address that test sever is listening on is contained in the ts.URL. We can use this along with the ts.Client().Get() method to make a GET /ping request against the test server. This returns a http.Response struct containing the response.

	// rs, err := ts.get(ts.URL + "/ping")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	//? We can check the value of the response status code and body using the same code as before
	// if rs.StatusCode != http.StatusOK {
	// 	t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	// }

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}

}

func TestShowSnippet(t *testing.T) {
	//? create an instance of our app that uses the mocked dependencies
	app := newTestApplication(t)

	//? Establish a new test server for runing end-to-end tests
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	//? Setup a table driven test to check the reponses sent by our application for different URLs.
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func TestSignupUser(t *testing.T) {
	app  := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	t.Log(csrfToken)
}
