package request

import (
	"fmt"
	"io"
	"net/http"
)

const keyServerAddr = "serverAddr"

func GetRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	has_first := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	has_second := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	fmt.Printf("%s: got / request. first(%t)=%s, second(%t)=%s\n",
		ctx.Value(keyServerAddr),
		has_first, first,
		has_second, second)

	io.WriteString(w, "Hello scraper\n")
}
