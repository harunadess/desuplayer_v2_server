# desuplayer v2 (server)
Server for [Desuplayer V2](https://github.com/jordanjohnston/desuplayer_v2_frontend) but could be repurposed for any frontend you want.


Basic http backend for getting meta data from music files only some formats supported - namely .wav is not supported yet
and then serving those files found via api requests and blobs.


## Releases
- Check releases tab and see if there is one (at time of writing there isn't).

## Manual Build/Installation
- Have [go (golang)](https://golang.org/dl/) installed
- Run `go build` in the root directory
- Run the executable (probably called main)
  - Alternatively, you can run `mv main desuplayer_v2` after to have a nicer named executable
  - Super alternatively, you can run `go install` instead and that will add it to your path
    - At some point, I might get around to figuring out how you name the executable properly.
- Server will listen on **localhost:4444**, and accept incoming requests from **localhost:8080**

### Endpoints
- *library/build*: builds library "database"
  - params:
    - musicDir: base directory to get files from
- *music/getAllArtists (to be renamed)*: gets all albums in library
  - params:
    - (none)
- *music/getSong*: gets song data (i.e. the raw bytes to then be used/played somewhere)
  - params:
    - path: path to song to get
- *music/getSongMeta*: gets song meta data (i.e. the album art, genre etc.)
  - params:
    - path: path to song to get
    - artist: artist of song
    - album: album of song

## Todo:

- Make some hardcoded things configurable (from files, probably)
- Figure out a way to do diffing (i.e. what changes there are from what you have stored)
  - This could potentially just be a set where if there is a path that is different, you just update based on that.
- Potential upgrade to just storing everything in JSON files. Nothing major, but you could probably move to something like redis with JSON extension or something? If Go has a package for it.
  - Might be nicer and more scalable than big json files, but then at the end of the day, you're only reading it once and then holding it in memory so /__/
- General refactoring and clean up