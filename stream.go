package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Stream represents a single video feed.
type Stream struct {
	m    map[chan []byte]bool
	lock sync.Mutex

	frame         []byte
	FrameInterval time.Duration
}

const boundary = "MJPEGBOUNDARY"
const headerFormat = "\r\n" +
	"--" + boundary + "\r\n" +
	"Content-Type: image/jpeg\r\n" +
	"Content-Length: %d\r\n" +
	"X-Timestamp: 0.000000\r\n" +
	"\r\n"

// ServeHTTP responds to HTTP requests with the MJPEG stream, implementing the http.Handler interface.
func (s *Stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Stream:", r.RemoteAddr, "connected")
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+boundary)

	c := make(chan []byte)
	s.lock.Lock()
	s.m[c] = true
	s.lock.Unlock()

	for {
		time.Sleep(s.FrameInterval)
		b := <-c
		_, err := w.Write(b)
		if err != nil {
			break
		}
	}

	s.lock.Lock()
	delete(s.m, c)
	s.lock.Unlock()
	log.Println("Stream:", r.RemoteAddr, "disconnected")
}

// UpdateJPEG pushes a new JPEG frame onto the clients.
func (s *Stream) UpdateJPEG(jpeg []byte) {
	header := fmt.Sprintf(headerFormat, len(jpeg))
	if len(s.frame) < len(jpeg)+len(header) {
		s.frame = make([]byte, (len(jpeg)+len(header))*2)
	}

	copy(s.frame, header)
	copy(s.frame[len(header):], jpeg)

	s.lock.Lock()
	for c := range s.m {
		select {
		case c <- s.frame:
		default:
		}
	}
	s.lock.Unlock()
}

// NewStream initializes and returns a new Stream.
func NewStream() *Stream {
	return &Stream{
		m:             make(map[chan []byte]bool),
		frame:         make([]byte, len(headerFormat)),
		FrameInterval: 50 * time.Millisecond,
	}
}
