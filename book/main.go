package main

import (
    "net/http"
    "sync"
    "text/template"
   "path/filepath"

    //"flag"
   // "os"
    "log"
    //"github.com/trace"
    "flag"
)

//tmp1は１つのテンプレートを表します
type templateHandler struct {
    once sync.Once 
    filename string 
    templ *template.Template
}
//serveHTTPはHTTPリクエストを処理します。
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    t.once.Do(func() {
        t.templ = 
            template.Must(template.ParseFiles(filepath.Join("templates",
                t.filename)))
            })
            t.templ.Execute(w,r)
    }


func main() {
//    http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {
//        w.Write([]byte(`
//            <html>
//                <head>
//                    <title>チャット</title>
//                    </head>
//                    <body>
//                    チャットしましょう
//                    </body>
//                    </html>
//                    `))
//                })
//

// /       //webサーバを開始します。
//        if err := http.ListenAndServe(":8080",nil); err != nil {
//            log.Fatal("ListenAndServe:",err)
 //       }
 //
 //   http.Handle("/",&templateHandler{filename: "chat.html"})
 //   //web サーバを開始します
 //   if err := http.ListenAndServe(":8080",nil); err != nil {
 //       log.Fatal("ListenAndServe:", err )
 //   }

 var addr = flag.String("addr",":8080","アプリケーションのアドレス")
 flag.Parse() //parse the flags

 r := newRoom()
 //r.tracer = trace.New(os.Stdout) //tracer.offをしない状態でこの一行がないと異常終了となる。

 http.Handle("/",&templateHandler{filename:"chat.html"})
 http.Handle("/room",r)

 //get the room going
 go r.run()

 //start the web server
 log.Println("Starting web server on ",*addr)
 if err := http.ListenAndServe(*addr,nil); err != nil {
     log.Fatal("ListenAndServe:",err)
 }
}