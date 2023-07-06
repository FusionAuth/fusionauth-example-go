//tag::baseApplication[]
package main

import (
  "fmt"
  "html/template"
  "net/http"
  "path"
)

func main() {
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
  http.HandleFunc("/", handleMain)
  http.HandleFunc("/login", handleFusionAuthLogin)
  http.HandleFunc("/callback", handleFusionAuthCallback)
  http.HandleFunc("/account", handleAccount)
  http.HandleFunc("/logout", handleLogout)

  port := "8080"

  fmt.Println("Starting HTTP server at http://localhost:" + port)
  fmt.Println(http.ListenAndServe(":" + port, nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
  WriteWebPage(w, "home.html", nil)
  return
}

func handleFusionAuthLogin(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/", http.StatusFound)
}

func handleFusionAuthCallback(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/", http.StatusFound)
}

func handleAccount(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/", http.StatusFound)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
}

// a function for writing a server-rendered web page
func WriteWebPage(w http.ResponseWriter, tmpl string, vars interface{}) {
  fn := path.Join("templates", tmpl)
  parsed_tmpl, error := template.ParseFiles(fn)

  if error != nil {
    http.Error(w, "Error reading template file " + tmpl + ": " + error.Error(), http.StatusInternalServerError)
    return
  }

  if error := parsed_tmpl.Execute(w, vars); error != nil {
    http.Error(w, error.Error(), http.StatusInternalServerError)
  }
}

func WriteCookie(w http.ResponseWriter, name string, value string, maxAge int) {
  cookie := http.Cookie{ Name: name, Domain: "localhost", Value: value, Path: "/", MaxAge: maxAge, HttpOnly: true, SameSite: http.SameSiteLaxMode, }
  http.SetCookie(w, &cookie)
}

//end::baseApplication[]
