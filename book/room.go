package main

import (
    "github.com/gorilla/websocket"
    "net/http"
    "log"
    //"github.com/trace"
)

type room struct {
    //forwardは他のクライアントに転送するためのメッセージを保持するチャネルです
    forward chan []byte
    //joinはチャットルームに参加しようとしているクライアントのためのチャネルです
    join chan *client 
    //leave はチャットルームから退室しようとしているクライアントのためのチャネルです。
    leave chan *client 
    //clientsには在室している全てのクライアントが保持されます。
    clients map[*client]bool
    //tracer will recieve tracer information of activity
    //in the room
    //tracer trace.Tracer
}

//newRoomはすぐに利用できるチャットルームを生成して返します。
func newRoom() *room{
    return &room{
        forward: make(chan []byte),
        join: make(chan *client),
        leave: make(chan *client),
        clients: make(map[*client]bool),
        //tracer:	trace.Off(),
    }
}


func (r *room) run() {
    for {
        select {
        case client := <- r.join://受信
            //参加
            r.clients[client] = true
        case client := <- r.leave://受信
            //退室
            delete(r.clients, client)
            close(client.send)
        case msg := <- r.forward://受信
            //全てのクライアントにメッセージを転送
            for client := range r.clients {
                select {
                case client.send <- msg://送信
                    //メッセージを他の人たちへ送信。全員が共有できるようになる。
                default:
                    //送信に失敗
                    delete(r.clients,client)
                    close(client.send)
                }
            }
        }
    }
}

const (
    socketBufferSize = 1024
    messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize:
    socketBufferSize, WriteBufferSize: socketBufferSize}

func(r *room) ServeHTTP(w http.ResponseWriter, req *http.Request){
    socket ,err := upgrader.Upgrade(w,req,nil)
    if err != nil {
        log.Fatal("ServeHTTP:",err )
        return 
    }
    client := &client{
        socket :socket,
        send: make(chan []byte, messageBufferSize),
        room: r,
    }
    r.join <- client//r.joinへ送信
    defer func() {r.leave <- client} ()//送信
    go client.write()
    client.from_client()
}

