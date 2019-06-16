package tests

import (
	LinkPreview "github.com/Junzki/link-preview"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// TestPreviewLink tests the generic entrance function.
func TestPreviewLink(t *testing.T) {
	server, serverErr := ServeLocalFile("test_cases/case_standard.html")
	if nil != serverErr {
		t.Error(serverErr.Error())
	}

	link := "http://custom-domain.local/case.html"
	client, _ := http.NewRequest("GET", server.URL, nil)

	handler, err := LinkPreview.Preview(link, client)
	if nil != err {
		t.Error(err)
	}

	assert.Equal(t, "test-title", handler.Title)
	assert.Equal(t, "test-desc", handler.Description)
	assert.Equal(t, "http://fake-image-url.image.png", handler.ImageURL)
	assert.Equal(t, link, handler.Link)
}
