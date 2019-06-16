package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func ServeLocalFile(filename string) (*httptest.Server, error) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := ioutil.ReadFile(filename)
		if nil != err {
			panic(err)
		}

		content := string(buf)
		_, _ = fmt.Fprint(w, content)
	}))

	return server, nil
}
