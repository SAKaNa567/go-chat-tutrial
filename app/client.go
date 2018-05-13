package main

import (
    "github.com/gorilla/websocket"
    "time"
)

//client はチャットをおこなっている１人のユーザを表します。
type client struct {
    //socketはこのクライアントのためのwebsocketです。
    socket *websocket.Conn
    //send はメッセージが送られるチャネルです。
    send chan *message
    //room はこのクライアントが参加しているチャットルーム
    room *room
    //usreDataはユーザーに関する情報を保持します。
    userData map[string]interface{}
}

func ( c *client ) from_client() {
    for {
        var msg *message
        if err := c.socket.ReadJSON(&msg); err == nil {
            msg.When = time.Now()
            msg.Name = c.userData["name"].(string)
            c.room.forward <- msg
        } else {
            break
        }
    }

    c.socket.Close()
}


func (c *client) write() {//メッセージを書き込んでいます。 Goで実行されているために、c.sendが送信された瞬間にこのメソッドは実行されます。
    for msg := range c.send{///フロントへ書き込んでいます。
        if err := c.socket.WriteJSON(msg); err != nil{
            break
        }
    }
    c.socket.Close()
}