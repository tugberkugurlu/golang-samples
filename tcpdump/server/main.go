package main

import (
	"fmt"
	"net/http"
	"time"
)

type handler struct {
}

func (h handler) ServeHTTP(w http.ResponseWriter,  r *http.Request) {
	t := time.NewTimer(1*time.Second)
	defer t.Stop()

	select {
	case <-r.Context().Done():
		fmt.Println("stopped processing")
	case <-t.C:
		fmt.Println("writing response")
		_, err := w.Write([]byte("hello world"))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	srv := http.Server{
		Addr: ":4400",
		Handler: handler{},
	}

	err := srv.ListenAndServe()
	fmt.Println(err)
}
