package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type connWrapper struct {
	net.Conn
}

func (ct *connWrapper) Close() error {
	err := ct.Conn.Close()
	if err != nil {
		fmt.Println("conn close failed", ct.Conn.LocalAddr(), ct.Conn.RemoteAddr())
	} else {
		fmt.Println("conn closed", ct.Conn.LocalAddr(), ct.Conn.RemoteAddr())
	}
	return err
}

func main() {
	time.Sleep(10*time.Second)
	t := http.DefaultTransport.(*http.Transport).Clone()
	dialer := net.Dialer{}
	t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := dialer.DialContext(ctx, network, addr)
		if err != nil {
			return conn, err
		}
		fmt.Println("conn established", conn.LocalAddr(), conn.RemoteAddr())
		return &connWrapper{
			conn,
		}, nil
	}
	t.MaxIdleConnsPerHost = 10
	c := http.Client{
		Timeout: 3*time.Second,
		Transport: t,
	}

	cc := make(chan os.Signal)
	signal.Notify(cc, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cc
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		c.CloseIdleConnections()
		os.Exit(0)
	}()

	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(1*time.Second)
			for i := 0; i < 5; i++ {
				func (){
					// fmt.Println("sending request in the background")
					r, err := c.Get("http://localhost:4400")
					defer func() {
						if r != nil && r.Body != nil {
							r.Body.Close()
						}
					}()
					if err != nil {
						fmt.Println(err)
						return
					}
					ioutil.ReadAll(r.Body)
				}()
			}
		}()
	}

	select {}
}
