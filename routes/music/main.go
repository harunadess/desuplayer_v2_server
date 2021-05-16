package music

import (
	"log"
	"net/http"

	"github.com/jordanjohnston/desuplayer_v2/library"
	"github.com/jordanjohnston/desuplayer_v2/routes/util"
)

const subRoute = "music"

func Routes() map[string]util.RequestHandler {
	return map[string]util.RequestHandler{
		util.FormatRoute(util.BaseRoute, subRoute, "get"): getSingle,
	}
}

func getSingle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")

	bs, err := library.GetFromLibrary(key)
	if err != nil {
		log.Printf("error getting track %v: %v", key, err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("error getting track " + err.Error()))
		return
	}
	w.Write(bs)
}
