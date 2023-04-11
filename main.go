package main

import (
	"io"
	"log"
	"net/http"
)

var (
	stream = NewStream()
)

func upload(_ http.ResponseWriter, r *http.Request) {
	bytes, _ := io.ReadAll(r.Body)
	stream.UpdateJPEG(bytes)
}

func main() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/download", stream.ServeHTTP)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html;charset=utf-8")
		_, _ = w.Write([]byte("<img width='100%' height='100%' src='/download'>"))
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
