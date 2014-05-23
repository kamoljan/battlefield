package main

import (
	"fmt"
	fb "github.com/huandu/facebook"
)

func main() {
	res, _ := fb.Get("/538744468", fb.Params{
		"fields": "username",
	})
	fmt.Println("here is my facebook username:", res["username"])

	// create a global App var to hold your app id and secret.
	var globalApp = facebook.New("your-app-id", "your-app-secret")

	// facebook asks for a valid redirect uri when parsing signed request.
	// it's a new enforced policy starting in late 2013.
	// it can be omitted in a mobile app server.
	globalApp.RedirectUri = "http://your-site-canvas-url/"

	// here comes a client with a facebook signed request string in query string.
	// creates a new session with signed request.
	session, _ := globalApp.SessionFromSignedRequest(signedRequest)

	// or, you just get a valid access token in other way.
	// creates a session directly.
	seesion := globalApp.Session(token)

	// use session to send api request with your access token.
	res, _ := session.Get("/me/feed", nil)

	// validate access token. err is nil if token is valid.
	err := session.Validate()
}
