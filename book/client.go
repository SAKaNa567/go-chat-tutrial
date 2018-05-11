package main

import (
    "github.com/gorilla/websocket"
)

//client はチャットをおこなっている１人のユーザを表します。
type client struct {
    //socketはこのクライアントのためのwebsocketです。
    socket *websocket.Conn
    //send はメッセージが送られるチャネルです。
    send chan []byte
    //room はこのクライアントが参加しているチャットルーム
    room *room
}

func ( c *client ) from_client() {
    for {
        if _, msg , err := c.socket.ReadMessage(); err == nil {
            c.room.forward <- msg //フロントから送られて来たmsgをforwardへ送信する。
           // fmt.Printf("%v -> %v",c.room,msg)
        }else{
            break
        }
    }
    c.socket.Close()
}

func (c *client ) write() {//メッセージを書き込んでいます。 Goで実行されているために、c.sendが送信された瞬間にこのメソッドは実行されます。
    for msg := range c.send {
        if err := c.socket.WriteMessage(websocket.TextMessage, msg);//フロントへ書き込んでいます。
        err != nil {
            break
        }
        //fmt.Println(msg)
    }
    c.socket.Close()
}

