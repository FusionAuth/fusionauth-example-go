package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/thanhpk/randstr"
	"golang.org/x/oauth2"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
)

var (
	FusionAuthConfig *oauth2.Config
	oauthStateString = randstr.Hex(16)
	// initialize the code verifier
	CodeVerifier, _ = cv.CreateCodeVerifier()

	// Create code_challenge with S256 method
	codeChallenge = CodeVerifier.CodeChallengeS256()
)

func init() {
	FusionAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "c671e83a-8cc4-4444-b44d-d5a4c581633b",
		ClientSecret: "2HYT86lWSAntc-mvtHLX5XXEpk9ThcqZb4YEh65CLjA",
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "http://localhost:9011/oauth2/authorize",
			TokenURL:  "http://localhost:9011/oauth2/token",
			AuthStyle: oauth2.AuthStyleInHeader,
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
	url := FusionAuthConfig.AuthCodeURL(oauthStateString, oauth2.SetAuthURLParam("code_challenge", codeChallenge), oauth2.SetAuthURLParam("code_challenge_method", "S256"))
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

	token, err := FusionAuthConfig.Exchange(oauth2.NoContext, code, oauth2.SetAuthURLParam("code_verifier", CodeVerifier.String()))
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	url := "http://localhost:9011/oauth2/userinfo"
	var bearer = "Bearer " + token.AccessToken
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
