package main

import "net/http"

func handleError(w http.ResponseWriter, r *http.Request) {
	respondWithERROR(w, 400, "Something went wrong")
}
