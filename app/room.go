package main

import (
    "github.com/gorilla/websocket"
    "net/http"
    "log"
    "Web_Application_By_Go/voyage/learning-in-goldenweek/trace"
    "github.com/stretchr/objx"
)

type room struct {
    //forwardは他のクライアントに転送するためのメッセージを保持するチャネルです
    forward chan *message
    //joinはチャットルームに参加しようとしているクライアントのためのチャネルです
    join chan *client 
    //leave はチャットルームから退室しようとしているクライアントのためのチャネルです。
    leave chan *client 
    //clientsには在室している全てのクライアントが保持されます。
    clients map[*client]bool
    //tracerはチャットルーム場で行われた操作ログを受け取ります。
    //in the room
    tracer trace.Tracer
}

//newRoomはすぐに利用できるチャットルームを生成して返します。
func newRoom() *room{
    return &room{
        forward: make(chan *message),
        join: make(chan *client),
        leave: make(chan *client),
        clients: make(map[*client]bool),
        tracer:	trace.Off(),//trace.TracerへOffメソッドを渡すことによりTraceメソッドの呼び出しを無視するTraceを実行できる。
    }
}


func (r *room) run() {
    for {
        select {
        case client := <- r.join://受信
            //参加
            r.clients[client] = true
            r.tracer.Trace("新しいクライアントが参加しました。")
        case client := <- r.leave://受信
            //退室
            delete(r.clients, client)
            close(client.send)
            r.tracer.Trace("クライアントが退室しました。")
        case msg := <- r.forward://受信
            //全てのクライアントにメッセージを転送
            r.tracer.Trace("メッセージを受信しました。",msg.Message)
            for client := range r.clients {
                select {
                case client.send <- msg://送信
                    r.tracer.Trace(" --- クライアントに送信されました。")
                    //メッセージを他の人たちへ送信。全員が共有できるようになる。
                default:
                    //送信に失敗
                    delete(r.clients,client)
                    close(client.send)
                    r.tracer.Trace(" --- 送信に失敗しました。クライアントをクリーンアップします。")
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
    authCookie, err := req.Cookie("auth")
    if err != nil {
        log.Fatal("クッキーの取得に失敗しました。:",err)
        return
    }


    client := &client{
        socket :socket,
        send: make(chan *message, messageBufferSize),
        room: r,
        userData: objx.MustFromBase64(authCookie.Value),
    }

    r.join <- client//r.joinへ送信
    defer func() {r.leave <- client} ()//送信
    go client.write()
    client.from_client()
}

