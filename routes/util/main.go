package util

import (
	"fmt"
	"net/http"
)

type RequestHandler func(w http.ResponseWriter, r *http.Request)
type Route interface {
	Route() map[string]RequestHandler
}

const BaseRoute string = "api"

func FormatRoute(urlParts ...interface{}) string {
	fmtString := ""
	for range urlParts {
		fmtString = fmtString + "/%v"
	}
	return fmt.Sprintf(fmtString, urlParts...)
}
