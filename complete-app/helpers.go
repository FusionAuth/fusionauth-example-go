package main

import (
  "html/template"
  "net/http"
  "path"
)

// a function for writing a server-rendered web page
func WriteWebPage(w http.ResponseWriter, tmpl string, vars interface{}) {
  fn := path.Join("templates", tmpl)
  parsed_tmpl, error := template.ParseFiles(fn)

  if error != nil {
    http.Error(w, error.Error(), http.StatusInternalServerError)
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
