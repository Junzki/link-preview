package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Change working dir (cwd) to project root.
	_ = os.Chdir("../")

	os.Exit(m.Run())
}

func serveLocalFile(filename string) (*httptest.Server, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}

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
