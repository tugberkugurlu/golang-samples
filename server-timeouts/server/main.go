package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

const largeFileEnvVar = "SERVER_TIMEOUT_SAMPLE_LARGE_FILE_PATH"

const (
	readTimeout  = 2 * time.Second
	writeTimeout = 3 * time.Second
)

type handler struct {
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.NewTimer(writeTimeout * 2) // simulate a long running process
	defer t.Stop()

	select {
	case <-r.Context().Done():
		log.Println("stopped processing")
	case <-t.C:
		log.Println("writing response", r.Proto)

		openFile, err := os.Open(os.Getenv(largeFileEnvVar))
		defer openFile.Close()
		if err != nil {
			http.NotFound(w, r)
			return
		}

		fileStat, _ := openFile.Stat()
		fileSize := strconv.FormatInt(fileStat.Size(), 10)
		w.Header().Set("Content-Length", fileSize)

		writeStart := time.Now()
		_, err = io.Copy(w, openFile)
		log.Printf("took %s\n", time.Now().Sub(writeStart))

		if err != nil {
			log.Printf("response write err: %v\n", err)
		}
	}
}

func main() {
	if os.Getenv(largeFileEnvVar) == "" {
		log.Fatalf("set lage file path through '%s' env variable", largeFileEnvVar)
	}

	srv := http.Server{
		Addr:         ":4400",
		Handler:      handler{},
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
