package main

import (
    "net/http"
    "log"
    "sync"
    "text/template"
    "path/filepath"
    "flag"
    "github.com/stretchr/gomniauth"
    "github.com/stretchr/gomniauth/providers/facebook"
    "github.com/stretchr/gomniauth/providers/github"
    "github.com/stretchr/gomniauth/providers/google"
    "github.com/stretchr/signature"
    "github.com/vaughan0/go-ini"
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

    file, _ := ini.LoadFile(".ini")
    facebook_id, _ := file.Get("FaceBook", "id")
    facebook_callback, _ := file.Get("FaceBook", "callback")
    github_id, _ := file.Get("GitHub", "id")
    github_secret, _ := file.Get("GitHub", "secret")
    github_callback, _ := file.Get("GitHub", "callback")
    google_id, _ := file.Get("Google", "id")
    google_secret, _ := file.Get("Google", "secret")
    google_callback, _ := file.Get("Google", "callback")
    gomniauth.SetSecurityKey(signature.RandomKey(64))
    gomniauth.WithProviders(
        facebook.New(facebook_id, "", facebook_callback),
        github.New(github_id, github_secret, github_callback),
        google.New(google_id, google_secret, google_callback),
    )
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
