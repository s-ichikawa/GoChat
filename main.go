package main

import (
    "net/http"
    "log"
    "sync"
    "text/template"
    "path/filepath"
    "flag"
)

type templateHandler struct {
    once     sync.Once
    filename string
    templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    t.once.Do(func() {
        t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
    })
    t.templ.Execute(w, r)
}

func main() {
    var addr = flag.String("host", ":8080", "アプリケーションのアドレス")
    flag.Parse()
    r := newRoom()
    //r.tracer = trace.New(os.Stdout)
    http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
    http.Handle("/login", &templateHandler{filename: "login.html"})
    http.HandleFunc("/auth/", loginHandler)
    http.Handle("/room", r)
    // チャットルームを開始
    go r.run()
    // Webサーバを起動
    log.Println("Webサーバを開始します。ポート:", *addr)
    if err := http.ListenAndServe(*addr, nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}
