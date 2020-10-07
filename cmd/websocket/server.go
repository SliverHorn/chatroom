package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		n, _ := fmt.Fprintln(writer, "Hello SliverHorn")
		fmt.Println(n)
	})
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := websocket.Accept(writer, request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		defer conn.Close(websocket.StatusInternalError, "内部出错了")

		ctx, cancel := context.WithTimeout(request.Context(), time.Second*10)
		defer cancel()
		var v interface{}
		err = wsjson.Read(ctx, conn, &v)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("接收到客户端:%v\n", v)

		err = wsjson.Write(ctx, conn, "Hello Websocket Client!")
		if err != nil {
			log.Println(err)
			return
		}
		err = conn.Close(websocket.StatusNormalClosure, "")
		if err != nil {
			log.Println(err)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":2021", nil))
}