package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jordanjohnston/desuplayer_v2/directoryscaper"
	"github.com/jordanjohnston/desuplayer_v2/fileutil"
	"github.com/jordanjohnston/desuplayer_v2/library"
)

type requestHandler func(w http.ResponseWriter, r *http.Request)

const baseRoute string = "api"

var routes = map[string]requestHandler{
	formatRoute(baseRoute, "library", "build"): libraryBuildLibrary,
	formatRoute(baseRoute, "library", "get"):   libraryGetLibrary,
	formatRoute(baseRoute, "music", "get"):     musicGetSingle,
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

func libraryBuildLibrary(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	baseDir := query.Get("musicDir")

	musicLibrary, err := directoryscaper.GetAllInDirectory(baseDir)
	if err != nil {
		log.Println("error getting music library ", err)
		w.WriteHeader(http.StatusInternalServerError)
		// todo: better handling
		w.Write([]byte("error getting music library"))
		return
	}
	library.LoadLibrary()

	// marshal to json before send
	jsonifiedLib, err := json.Marshal(musicLibrary)
	if err != nil {
		log.Println("error converting music library to json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error converting music library to json"))
		return
	}
	w.Write(jsonifiedLib)
}

func musicGetSingle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")
	log.Printf("musicGetSingle -> key:%v\n", key)

	bs, err := library.GetFromLibrary(key)
	if err != nil {
		log.Printf("error getting track %v: %v", key, err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("error getting track " + err.Error()))
		return
	}
	w.Write(bs)
}

func libraryGetLibrary(w http.ResponseWriter, r *http.Request) {

	file, err := fileutil.ReadSingleFile("./library.json")
	if err != nil {
		log.Printf("error getting library: %v", err)
		if strings.Contains(err.Error(), "cannot find") {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("could not find library file - run 'build library' again." + err.Error()))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error getting track " + err.Error()))
		}
		return
	}
	w.Write(file)
}
