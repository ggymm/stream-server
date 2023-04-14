package main

import (
	"bytes"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	stream = NewStream()
)

func upload(_ http.ResponseWriter, r *http.Request) {
	origin, _ := io.ReadAll(r.Body)

	inputImage, err := jpeg.Decode(bytes.NewBuffer(origin))
	if err != nil {
		log.Fatalln(err)
	}

	var compress bytes.Buffer
	err = jpeg.Encode(&compress, inputImage, &jpeg.Options{Quality: 50})
	if err != nil {
		log.Fatalln(err)
	}

	stream.UpdateJPEG(compress.Bytes())
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
