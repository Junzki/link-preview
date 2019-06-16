package tests

import (
	"github.com/Junzki/link-preview/handlers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetHandler(t *testing.T) {
	cxt := handlers.LinkPreviewContext{
		TargetType: handlers.StandardMetaTags,
		Link: "https://golang.org",
	}

	handler, err := handlers.GetPreviewHandler(&cxt)

	assert.NotNil(t, handler)
	assert.Nil(t, err)
}

func TestGetHandlerBadTargetType(t *testing.T) {
	cxt := handlers.LinkPreviewContext{
		TargetType: -1,
		Link: "https://golang.org",
	}

	handler, err := handlers.GetPreviewHandler(&cxt)

	assert.Nil(t, handler)
	assert.NotNil(t, err)
}
