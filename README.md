# desuplayer v2 (server)
Server for Desuplayer V2


## Todo:
- Figure out a way to do diffing (i.e. what changes there are from what you have stored)
  - This could potentially just be a set where if there is a path that is different, you just update based on that.
- Potential upgrade to just storing everything in JSON files. Nothing major, but you could probably move to something like redis with JSON extension or something? If Go has a package for it.
  - Might be nicer and more scalable than big json files, but then at the end of the day, you're only reading it once and then holding it in memory so /__/
- General refactoring and clean up