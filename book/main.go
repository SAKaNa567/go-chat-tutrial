package main

import (
    "log"
    "net/http"
    "sync"
    "text/template"
   "path/filepath"
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
            t.templ.Execute(w,nil)
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
    http.Handle("/",&templateHandler{filename: "chat.html"})
    //web サーバを開始します
    if err := http.ListenAndServe(":8080",nil); err != nil {
        log.Fatal("ListenAndServe:", err )
    }
}

