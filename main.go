package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	

	"golang.org/x/oauth2"
)

var (
	FusionAuthConfig *oauth2.Config
	// TODO: randomize it
	oauthStateString = "pseudo-random1"
)

func init() {
	FusionAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "7d2b4cb4-ccd5-42ac-8469-f802393c8f98",
		Scopes:       []string{"email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:9011/oauth2/authorize",
			TokenURL: "http://localhost:9011/oauth2/token",
		},
	}
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleFusionAuthLogin)
	http.HandleFunc("/callback", handleFusionAuthCallback)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
<body>
	<a href="/login">FusionAuth Log In</a>
</body>
</html>`

	fmt.Fprintf(w, htmlIndex)
}

func handleFusionAuthLogin(w http.ResponseWriter, r *http.Request) {
	url := FusionAuthConfig.AuthCodeURL(oauthStateString)
	
	// fmt.Println(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleFusionAuthCallback(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Content: %s\n", content)
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := FusionAuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	
	url := "http://localhost:9011/oauth2/userinfo"
	var bearer = "Bearer "+token.AccessToken
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close() 
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}