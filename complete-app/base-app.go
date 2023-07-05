//tag::baseApplication[]
package main

import (
  "fmt"
  "net/http"
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
//end::baseApplication[]

//tag::loginRoute[]
func handleFusionAuthLogin(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/", http.StatusFound)
}
//end::loginRoute[]

//tag::callbackRoute[]
func handleFusionAuthCallback(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/", http.StatusFound)
}
//end::callbackRoute[]

//tag::accountRoute[]
func handleAccount(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/", http.StatusFound)
}
//end::accountRoute[]

//tag::logoutRoute[]
func handleLogout(w http.ResponseWriter, r *http.Request) {
}
//end::logoutRoute[]
