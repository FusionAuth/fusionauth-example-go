package main

import (
  "errors"
  "fmt"
  "os"
  "net/http"
  "net/url"
  "github.com/thanhpk/randstr"
  "golang.org/x/oauth2"
  "github.com/coreos/go-oidc/v3/oidc"
)

// This supplies data to the /account page
type AccountVars struct {
  LogoutUrl string
  Email string
}

//tag::oidcConstants[]
const (
  FusionAuthHost string = "http://localhost:9011"
  FusionAuthTenantID string = "d7d09513-a3f5-401c-9685-34ab6c552453"
  FusionAuthClientID string = "e9fdb985-9173-4e01-9d73-ac2d60d1dc8e"
  FusionAuthClientSecret string = "2HYT86lWSAntc-mvtHLX5XXEpk9ThcqZb4YEh65CLjA-not-for-prod"
  AccessTokenCookieName string = "cb_access_token"
  RefreshTokenCookieName string = "cb_refresh_token"
  IDTokenCookieName string = "cb_id_token"
)
//end::oidcConstants[]

//tag::oidcClient[]
var (
  oidcProvider *oidc.Provider 
  fusionAuthConfig *oauth2.Config

  // In a production application, we would persist a unique state string for each login request
  oauthStateString string = randstr.Hex(16)
)

func init() {
  // configure an OIDC provider
  provider, err := oidc.NewProvider(oauth2.NoContext, FusionAuthHost)

  if err != nil {
    fmt.Println("Error creating OIDC provider: " + err.Error())
  } else {

    oidcProvider = provider

    // initialize OAuth
    fusionAuthConfig = &oauth2.Config{
      ClientID:     FusionAuthClientID,
      ClientSecret: FusionAuthClientSecret,
      RedirectURL:  "http://localhost:8080/callback?this=that",
      Endpoint:     oidcProvider.Endpoint(),
      Scopes:       []string{oidc.ScopeOpenID, "offline_access"},
    }
  }
}
//end::oidcClient[]

func main() {
//  if fusionAuthConfig == nil {
//    fmt.Println("Error configuring OAuth, exiting");
//    os.Exit(1);
//  }

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
  http.HandleFunc("/", handleMain)
  http.HandleFunc("/login", handleLoginRequest)
  http.HandleFunc("/callback", handleFusionAuthCallback)
  http.HandleFunc("/account", handleAccount)
  http.HandleFunc("/logout", handleLogout)

  port := "8080"

  fmt.Println("Starting HTTP server at http://localhost:" + port)
  fmt.Println(http.ListenAndServe(":" + port, nil))
}

//tag::main[]
func handleMain(w http.ResponseWriter, r *http.Request) {

  // see if the user is authenticated. In a real application, we would validate the token signature and expiration
  _, err := r.Cookie(AccessTokenCookieName)

  if err != nil {
    WriteWebPage(w, "home.html", nil)
    return
  }

  // The user is authenticated so redirect to /account. We redirect so that the location in the browser shows /account.
  http.Redirect(w, r, "/account", http.StatusFound)
  return
}
//end::main[]

//tag::loginRoute[]
func handleLoginRequest(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, fusionAuthConfig.AuthCodeURL(oauthStateString), http.StatusFound)
}
//end::loginRoute[]

//tag::callbackRoute[]
func handleFusionAuthCallback(w http.ResponseWriter, r *http.Request) {

  // validate the state value, to make sure this came from us
  if r.FormValue("state") != oauthStateString {
    http.Error(w, "Bad request - incorrect state value", http.StatusBadRequest)
    return
  }

  // Exchange the authorization code for access, refresh, and id tokens
  token, err := fusionAuthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))

  if err != nil {
    http.Error(w, "Error getting access token: " + err.Error(), http.StatusInternalServerError)
    return
  }

  rawIDToken, ok := token.Extra("id_token").(string)

  if !ok {
    http.Error(w, "No ID token found in request to /callback", http.StatusBadRequest)
    return
  }

  // Write access, refresh, and id tokens to http-only cookies
  WriteCookie(w, AccessTokenCookieName, token.AccessToken, 3600)
  WriteCookie(w, RefreshTokenCookieName, token.RefreshToken, 3600)
  WriteCookie(w, IDTokenCookieName, rawIDToken, 3600)

  http.Redirect(w, r, "/account", http.StatusFound)
}
//end::callbackRoute[]

//tag::accountRoute[]
func getLogoutUrl() string {
  url := fmt.Sprintf("%s/oauth2/logout?client_id=%s&tenantId=%s", FusionAuthHost, url.QueryEscape(FusionAuthClientID), url.QueryEscape(FusionAuthTenantID))
  return url
}

func handleAccount(w http.ResponseWriter, r *http.Request) {

  // Make sure the user is authenticated. Note that in a production application, we would validate the token signature, 
  // make sure it wasn't expired, and attempt to refresh it if it were
  cookie, err := r.Cookie(AccessTokenCookieName)

  if err != nil || cookie == nil {
    http.Redirect(w, r, "/", http.StatusFound)
    return
  }

  // Now get the ID token so we can show the user's email address
  cookie, err = r.Cookie(IDTokenCookieName)

  if err != nil || cookie == nil{
    http.Error(w, "No ID token found", http.StatusBadRequest)
    return
  }

  var verifier = oidcProvider.Verifier(&oidc.Config{ClientID: FusionAuthClientID})

  idToken, err := verifier.Verify(oauth2.NoContext, cookie.Value)

  if err != nil {
    http.Error(w, "Error verifying ID token: " + err.Error(), http.StatusBadRequest)
    return
  }

  // Extract the email claim
  var claims struct {
    Email    string `json:"email"`
  }

  if err := idToken.Claims(&claims); err != nil {
    http.Error(w, "Error reading claims from ID token: " + err.Error(), http.StatusInternalServerError)
    return
  }

  templateVars := AccountVars{LogoutUrl: getLogoutUrl(), Email: claims.Email}

  WriteWebPage(w, "account.html", templateVars)
}
//end::accountRoute[]

//tag::logoutRoute[]
func handleLogout(w http.ResponseWriter, r *http.Request) {
  // Delete the cookies we set
  WriteCookie(w, AccessTokenCookieName, "", -1)
  WriteCookie(w, RefreshTokenCookieName, "", -1)
  WriteCookie(w, IDTokenCookieName, "", -1)

  http.Redirect(w, r, "/", http.StatusFound)
}
//end::logoutRoute[]
