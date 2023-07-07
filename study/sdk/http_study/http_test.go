package http_study

import (
	"io"
	"net/http"
	"testing"
)

func TestHttp(t *testing.T) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, "Hello, World\n")
	})

	http.ListenAndServe(":8080", nil)
}
