package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/lxzan/gws"
)

func main() {
	socket, _, err := gws.NewClient(new(gws.BuiltinEventHandler), &gws.ClientOption{
		Addr:      "ws://127.0.0.1:6666/connect",
		TlsConfig: &tls.Config{InsecureSkipVerify: true},
	})
	if err != nil {
		log.Println(err.Error())
		return
	}
	for i := 0; i < 100; i++ {
		socket.WriteString("hello world" + string(i))
	}
	var channel = make(chan []byte, 8)
	var closer = make(chan struct{})
	socket.SessionStorage.Store("channel", channel)
	socket.SessionStorage.Store("closer", closer)
	go socket.ReadLoop()
	go func() {
		for {
			select {
			case p := <-channel:
				_ = socket.WriteMessage(gws.OpcodeText, p)
				fmt.Println("send:", string(p))
			case <-closer:
				return
			}
		}
	}()
	time.Sleep(10 * time.Second)
}
