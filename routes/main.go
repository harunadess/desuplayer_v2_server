package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jordanjohnston/desuplayer_v2/directoryscaper"
	"github.com/jordanjohnston/desuplayer_v2/library"
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
	query := r.URL.Query()
	baseDir := query.Get("musicDir")

	musicLibrary, err := directoryscaper.GetAllInDirectory(baseDir)
	if err != nil {
		log.Println("error getting music library ", err)
		w.WriteHeader(http.StatusInternalServerError)
		// todo: better handling
		w.Write([]byte(""))
		return
	}
	// marshal to json before send
	jsonifiedLib, err := json.Marshal(musicLibrary)
	if err != nil {
		log.Println("error converting music library to json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(""))
		return
	}
	w.Write(jsonifiedLib)
}

func musicGetSingle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	key := query.Get("key")
	fmt.Printf("key:%v\n", key)

	// filePath := "/mnt/d/Users/Jorta/Music/Alestorm/Back Through Time (Limited Edition)/01 - Back Through Time.mp3"
	// filePath := "D:\\Users\\Jorta\\Music/Tanooki Suit\\Kosmonaut (FLAC)\\Tanuki - Kosmonaut . Diver - 01 Kosmonaut.flac"
	// bs, err := fileutil.ReadSingleFile(filePath)
	bs, err := library.GetFromLibrary(key)
	if err != nil {
		log.Printf("error getting track %v: %v", key, err)
		fmt.Fprintf(w, "failed to get file %v\n", key)
	}
	w.Write(bs)
}
