package link_preview

import (
	"fmt"
	"testing"
)


func TestPreviewLink(t *testing.T) {
	var url = "https://github.com/aakash4525/py_link_preview/blob/master/link_preview/link_preview.py"

	preview, err := PreviewLink(url)
	if nil != err {
		t.Error(err)
	}

	title := preview.Title
	fmt.Println(title)


}
