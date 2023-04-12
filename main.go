package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

var (
	stream = NewStream()
)

func upload(_ http.ResponseWriter, r *http.Request) {
	bytes, _ := io.ReadAll(r.Body)
	stream.UpdateJPEG(bytes)
}

func main() {
	port := "8080"
	if len(os.Args) == 2 {
		port = os.Args[1]
	}

	http.HandleFunc("/upload", upload)
	http.HandleFunc("/download", stream.ServeHTTP)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html;charset=utf-8")
		_, _ = w.Write([]byte("<img width='100%' height='100%' src='/download'>"))
	})

	log.Println("server listening on", port)
	log.Fatalln(http.ListenAndServe("0.0.0.0:"+port, nil))
}
