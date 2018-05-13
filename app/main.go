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
    "github.com/stretchr/gomniauth"
    "github.com/stretchr/gomniauth/providers/github"
    "github.com/stretchr/gomniauth/providers/google"
    "github.com/stretchr/gomniauth/providers/facebook"
    "github.com/stretchr/objx"
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
        data := map[string]interface{}{
            "Host": r.Host,
        }
        if authCookie, err := r.Cookie("auth"); err == nil {
            data["UserData"] = objx.MustFromBase64(authCookie.Value)
        }
            t.templ.Execute(w,data)
    }


func main() {

 var addr = flag.String("addr",":8080","アプリケーションのアドレス")
 flag.Parse() //parse the flags



 //Gomniauthのセットアップ
 gomniauth.SetSecurityKey("セキュリティキー")
 gomniauth.WithProviders(
     facebook.New("クライアントID","秘密の値","http://localhost:8080/auth/callback/facebook"),
     github.New("クライアントID","秘密の値","http://localhost:8080/auth/callback/github"),
     google.New("652120078645-lklhekp7f29qo895aqpikitksao2p05k.apps.googleusercontent.com","KFOo-MOctyM2nCsJRXa7_RqY","http://localhost:8080/auth/callback/google"),
 )

 r := newRoom()
 //r.tracer = trace.New(os.Stdout) //tracer.offをしない状態でこの一行がないと異常終了となる。

 http.Handle("/chat",MustAuth(&templateHandler{filename:"chat.html"}))//こうすることで、まずauthhandlerのServeHTTP-> templateHandelerのServeHTTPとなる。
 http.Handle("/login",&templateHandler{filename:"login.html"})
 http.HandleFunc("/auth/",loginHandler)
 http.Handle("/room",r)

 //get the room going
 go r.run()

 //start the web server
 log.Println("Starting web server on ",*addr)
 if err := http.ListenAndServe(*addr,nil); err != nil {
     log.Fatal("ListenAndServe:",err)
 }
}