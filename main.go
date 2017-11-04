package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/garyburd/go-oauth/oauth"
)

const callbackURL = "http://localhost:8080/access_token"

var oauthClient = oauth.Client{}
var tempCred *oauth.Credentials

func main() {
	oauthClient.TemporaryCredentialRequestURI = "https://api.twitter.com/oauth/request_token"
	oauthClient.ResourceOwnerAuthorizationURI = "https://api.twitter.com/oauth/authenticate"
	oauthClient.TokenRequestURI = "https://api.twitter.com/oauth/access_token"

	oauthClient.Credentials.Token = os.Getenv("CONSUMER_KEY")
	oauthClient.Credentials.Secret = os.Getenv("CONSUMER_SECRET")

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/access_token", AccessTokenHandler)
	http.ListenAndServe(":8080", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	c, err := oauthClient.RequestTemporaryCredentials(http.DefaultClient, callbackURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	tempCred = c
	url := oauthClient.AuthorizationURL(tempCred, nil)
	http.Redirect(w, r, url, http.StatusFound)
}

func AccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	c, _, err := oauthClient.RequestToken(http.DefaultClient, tempCred, r.URL.Query().Get("oauth_verifier"))
	if err != nil {
		log.Fatal(err)
	}
	const templ = `<h1>Your Access Token</h1><p>Access Token: {{.Token}}</p><p>Access Token Secret: {{.Secret}}</p>`
	t := template.Must(template.New("request").Parse(templ))
	t.Execute(w, c)
}
