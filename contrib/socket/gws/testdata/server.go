package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lxzan/gws"
)

func main() {
	upgrader := gws.NewUpgrader(new(gws.BuiltinEventHandler), &gws.ServerOption{
		// Authorize: func(r *http.Request, session gws.SessionStorage) bool {
		// 	session.Store("username", r.URL.Query().Get("username"))
		// 	return true
		// },
	})

	http.HandleFunc("/connect", func(writer http.ResponseWriter, request *http.Request) {
		socket, err := upgrader.Upgrade(writer, request)
		if err != nil {
			log.Printf(err.Error())
			return
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
		a, b := socket.SessionStorage.Load("channel")
		if b {
			for {
				fmt.Printf("%v", <-a.(chan []byte))
			}
		}
		fmt.Println(socket.SessionStorage.Load("channel"))
	})

	if err := http.ListenAndServe(":6666", nil); err != nil {
		log.Fatalf("%v", err)
	}
}
