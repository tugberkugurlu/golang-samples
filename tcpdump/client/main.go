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
	dialer := net.Dialer{
		// KeepAlive: 1*time.Second,
	}
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
	c := http.Client{
		Timeout: 2*time.Second,
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
	func (){
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
		val, _ := ioutil.ReadAll(r.Body)
		fmt.Println(val)
	}()
	select {}
}
