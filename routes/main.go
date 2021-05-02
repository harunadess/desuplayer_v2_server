package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jordanjohnston/desuplayer_v2/fileUtil"
)

type requestHandler func(w http.ResponseWriter, r *http.Request)

const baseRoute string = "api"

var routes = map[string]requestHandler{
	formatRoute(baseRoute, "music", "getAll"): musicGetAll,
	formatRoute(baseRoute, "music", "get"):    musicGetSingle,
}

func formatRoute(urlParts ...interface{}) string {
	fmtString := ""
	for range urlParts {
		fmtString = fmtString + "/%v"
	}
	return fmt.Sprintf(fmtString, urlParts...)
}

func SetUpRequestHandlers(m *http.ServeMux) {
	for route, handler := range routes {
		m.HandleFunc(route, handler)
	}
}

func musicGetAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v: %v", r.Method, r.RequestURI)
	fmt.Printf("hit /api/music/getAll")
}

func musicGetSingle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	track := query.Get("track")
	artist := query.Get("artist")
	fmt.Printf("track:%v artist:%v\n", track, artist)

	// filePath := "/mnt/d/Users/Jorta/Music/Alestorm/Back Through Time (Limited Edition)/01 - Back Through Time.mp3"
	filePath := "/mnt/d/Users/Jorta/Music/Tanooki Suit/Kosmonaut (FLAC)/Tanuki - Kosmonaut . Diver - 01 Kosmonaut.flac"
	bs, err := fileUtil.ReadSingleFile(filePath)
	if err != nil {
		log.Printf("error getting track %v: %v", track, err)
		fmt.Fprintf(w, "failed to get file %v\n", track)
	}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(bs)
}
