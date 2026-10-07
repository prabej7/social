package main

import "net/http"

func (app *application) health(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ok"))
}
