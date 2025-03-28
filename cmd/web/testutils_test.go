package main

import (
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/Ng1n3/go-further/pkg/models/mock"
	"github.com/golangcollege/sessions"
)

func newTestApplication(t *testing.T) *application {
	//? Create an instance of the template cache.
	templateCache, err := newTemplateCache(("./../../ui/html/"))

	if err != nil {
		t.Fatal(err)
	}

	//? Create a session manager instance, with the same settings as production.
	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	//? Initialize the dependecies, using the mocks for the loggers and database models.
	return &application{
		errorLog:      log.New(io.Discard, "", 0),
		infoLog:       log.New(io.Discard, "", 0),
		session:       session,
		snippets:      &mock.SnippetModel{},
		templateCache: templateCache,
		users:         &mock.UserModel{},
	}
}

// Define a custom testServer type which anonymously embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// Creat a newTestServer helper which initializes and returns a new instance of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	//? Add cookie jar to the client, so that response cookies are stored and then sent with subsequent requests.
	ts.Client().Jar = jar

	//? Disable redirect-following for the client.Essentially this function is called after a 3xx response is received by the client, and returning the http.ErrUseLastResponse error forces it to immediately return the received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
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

//? Define a regular expression which captures the CSRF token value from the HTML for our user signup page.
var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='([^']+)'/>`)

func extractCSRFToken(t *testing.T, body []byte) string {
  //? use findsubmatch method to extract the token from the html body.
  matches := csrfTokenRX.FindSubmatch(body)
  if len(matches) < 2 {
    t.Fatal("no csrf token found in the body")
  }

  return html.UnescapeString(string(matches[1]))
}
