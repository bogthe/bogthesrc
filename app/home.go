package app

import (
	"fmt"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprint(w, "You are now home ET")
	return nil
}
