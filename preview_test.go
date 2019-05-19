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
		fmt.Fprint(w, content)
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
		fmt.Fprint(w, content)
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
