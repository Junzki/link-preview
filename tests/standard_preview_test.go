package tests

import (
	"fmt"
	"github.com/Junzki/link-preview/handlers"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPreviewStandard(t *testing.T) {
	server, serverErr := serveLocalFile("test_cases/case_standard.html")
	if nil != serverErr {
		t.Error(serverErr.Error())
	}

	link := "http://custom-domain.local/case.html"
	client, _ := http.NewRequest("GET", server.URL, nil)

	cxt := &handlers.LinkPreviewContext{
		TargetType: handlers.StandardMetaTags,
		Link:       link,
		Client:     client,
	}

	handler := handlers.StandardLinkPreview{
		cxt,
	}

	err := handler.Preview()
	if nil != err {
		t.Error(err)
	}

	assert.Equal(t, "test-title", handler.Title)
	assert.Equal(t, "test-desc", handler.Description)
	assert.Equal(t, "http://fake-image-url.image.png", handler.ImageURL)
	assert.Equal(t, link, handler.Link)
}

func TestPreviewFallback(t *testing.T) {
	server, serverErr := serveLocalFile("test_cases/case_standard_fallback.html")
	if nil != serverErr {
		t.Error(serverErr.Error())
	}

	link := "http://custom-domain.local/case.html"
	client, _ := http.NewRequest("GET", server.URL, nil)

	cxt := &handlers.LinkPreviewContext{
		TargetType: handlers.StandardMetaTags,
		Link:       link,
		Client:     client,
	}

	handler := handlers.StandardLinkPreview{
		cxt,
	}

	err := handler.Preview()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "Title", handler.Title)
	assert.Equal(t, "meta-desc", handler.Description)
	assert.Equal(t, "http://custom-domain.local/favicon.ico", handler.ImageURL)
	assert.Equal(t, link, handler.Link)
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

		buf, err := ioutil.ReadFile("test_cases/case_standard_fallback.html")
		if nil != err {
			panic(err)
		}

		content := string(buf)
		_, _ = fmt.Fprint(w, content)
	}))

	link := "http://custom-domain.local/case.html"
	client, _ := http.NewRequest("GET", server.URL, nil)

	cxt := &handlers.LinkPreviewContext{
		TargetType: handlers.StandardMetaTags,
		Link:       link,
		Client:     client,
	}

	handler := handlers.StandardLinkPreview{
		cxt,
	}

	err := handler.Preview()
	if nil != err {
		t.Error(err)
	}

	assert.Equal(t, "Title", handler.Title)
}
