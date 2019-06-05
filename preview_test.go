package link_preview

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestPreviewStandard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadFile("test_cases/case_standard.html")
		if nil != err {
			panic(err)
		}

		content := string(buf)
		_, _ = fmt.Fprint(w, content)
	}))

	link := "http://custom-domain.local/case.html"
	client, _ := http.NewRequest("GET", server.URL, nil)

 	result, err := PreviewLink(link, client)
 	if err != nil {
 		t.Error(err)
	}

 	assert.Equal(t, "test-title", result.Title)
 	assert.Equal(t, "test-desc", result.Description)
 	assert.Equal(t, "http://fake-image-url.image.png", result.ImageURL)
 	assert.Equal(t, link, result.Link)
}


func TestPreviewFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadFile("test_cases/case_fallback.html")
		if nil != err {
			panic(err)
		}

		content := string(buf)
		_, _ = fmt.Fprint(w, content)
	}))

	link := "http://custom-domain.local/case.html"
	client, _ := http.NewRequest("GET", server.URL, nil)

	result, err := PreviewLink(link, client)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "Title", result.Title)
	assert.Equal(t, "meta-desc", result.Description)
	assert.Equal(t, "http://custom-domain.local/favicon.ico", result.ImageURL)
	assert.Equal(t, link, result.Link)
}


func TestRedirect(t *testing.T) {
	// This test shows that Go's built in http client can follow redirects automatically.
	// Check out this for solution to disable this feature:
	// - https://stackoverflow.com/questions/23297520/how-can-i-make-the-go-http-client-not-follow-redirects-automatically
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := r.URL
		query := target.Query()
		redir := query.Get("redirect")

		query.Set("redirect", "true")
		target.RawQuery = query.Encode()


		if "" == redir {
			http.Redirect(w, r, target.String(), http.StatusFound)
		}

		buf, err := ioutil.ReadFile("test_cases/case_fallback.html")
		if nil != err {
			panic(err)
		}

		content := string(buf)
		_, _ = fmt.Fprint(w, content)
	}))

	result, err := PreviewLink(server.URL, nil)
	if nil != err {
		t.Error(err)
	}

	assert.Equal(t, "Title", result.Title)
}