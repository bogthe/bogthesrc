package app

import (
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) error {
	return renderTemplate(w, r, HomeTemplate, http.StatusOK, nil)
}
