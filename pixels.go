package pixels

import (
	"fmt"
	"net/http"
)

func PixelsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi")
}
