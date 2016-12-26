package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"../../kinli"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func init() {
	kinli.SessionStore = sessions.NewFilesystemStore("./sessions", []byte("some-secret-string"))

	kinli.CacheMode = false // set to true if you don't want live editing on views
	kinli.ClientConfig = make(map[string]string)
	kinli.ViewFuncs = template.FuncMap{
		"hello": hello,
	}
	kinli.InitTmpl()

}

func hello(name string) string {
	return fmt.Sprintf("hello %s", name)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		hc := &kinli.HttpContext{w, r}
		page := kinli.NewPage(hc, "hello page", "", "", nil)
		kinli.DisplayPage(w, "home", page)
	})
	srv := &http.Server{Handler: r, Addr: "localhost:3000"}
	fmt.Println("Open the page http://localhost:3000/home")
	log.Fatal(srv.ListenAndServe())
}
