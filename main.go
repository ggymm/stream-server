package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for i := 0; i < 6; i++ {
		go func(i int) {
			stream := NewStream()
			router := http.NewServeMux()

			router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "text/html;charset=utf-8")
				_, _ = w.Write([]byte("<img width='100%' height='100%' src='/download'>"))
			})
			router.HandleFunc("/upload", func(_ http.ResponseWriter, r *http.Request) {
				origin, _ := io.ReadAll(r.Body)

				// fmt.Println(len(origin))
				// stream.UpdateJPEG(origin)

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
			})
			router.HandleFunc("/download", stream.ServeHTTP)

			port := 8000 + i
			log.Println("server listening on", port)

			server := &http.Server{
				Addr:    fmt.Sprintf("0.0.0.0:%d", port),
				Handler: router,
			}
			log.Fatalln(server.ListenAndServe())
		}(i)
	}

EXIT:
	for {
		sig := <-sc
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	time.Sleep(time.Second)
	os.Exit(state)
}
