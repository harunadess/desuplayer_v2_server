# desuplayer v2 (server)
Server for [Desuplayer V2](https://github.com/jordanjohnston/desuplayer_v2_frontend)


## General Info

Basic http backend for getting meta data from music files only some formats supported - namely .wav is not supported yet
and then serving files found via blobs.

### Deployment

- Have go (golang) installed
- Run `go build` in the root directory
- Run the executable (probably called main)
  - Alternatively, you can run `mv main desuplayer_v2` after to have a nicer named executable
- Server will listen on localhost:4444


## Todo:

- Make some hardcoded things configurable (from files, probably)
- Figure out a way to do diffing (i.e. what changes there are from what you have stored)
  - This could potentially just be a set where if there is a path that is different, you just update based on that.
- Potential upgrade to just storing everything in JSON files. Nothing major, but you could probably move to something like redis with JSON extension or something? If Go has a package for it.
  - Might be nicer and more scalable than big json files, but then at the end of the day, you're only reading it once and then holding it in memory so /__/
- General refactoring and clean up