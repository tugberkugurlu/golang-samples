package main

import (
	"log"
	"net"
	"net/http"
	"time"
)

const (
	fakeRequestProcessingTime = 10 * time.Second
	readTimeout               = fakeRequestProcessingTime / 5
	writeTimeout              = fakeRequestProcessingTime / 2
)

type handler struct {
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.NewTimer(fakeRequestProcessingTime) // simulate a long running process
	defer t.Stop()

	select {
	case <-r.Context().Done():
		log.Println("stopped processing")
	case <-t.C:
		log.Println("writing response", r.Proto)
		_, err := w.Write([]byte("hello world"))
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	srv := http.Server{
		Addr:         ":4400",
		Handler:      http.TimeoutHandler(handler{}, writeTimeout/2, ""),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		ConnState: func(conn net.Conn, state http.ConnState) {
			log.Printf("[%s,%s]: %s", conn.RemoteAddr().String(),
				conn.LocalAddr().String(),
				state.String())
		},
	}

	err := srv.ListenAndServe()
	log.Println(err)
}
