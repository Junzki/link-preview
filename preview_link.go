package LinkPreview

import (
	"github.com/Junzki/link-preview/handlers"
	"net/http"
)

func Preview(link string, extraContent *http.Request) (*handlers.LinkPreviewContext, error) {
	return PreviewLink(link, extraContent)
}

func PreviewLink(link string, extraClient *http.Request) (*handlers.LinkPreviewContext, error) {
	cxt := &handlers.LinkPreviewContext{
		Link: link,
		TargetType: handlers.StandardMetaTags,
	}

	if nil != extraClient {
		cxt.Client = extraClient
	}

	handler, handlerErr := handlers.GetPreviewHandler(cxt)
	if nil != handlerErr {
		return nil, handlerErr
	}

	cxt, previewErr := handler.Preview()
	if nil != previewErr {
		return nil, previewErr
	}

	return cxt, nil
}
