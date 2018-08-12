// Package recaptcha provides support for reCaptcha 2.0 (https://www.google.com/recaptcha) user response verification.
// It allows the use of a custom http.Client, and will fall back to http.DefaultClient if none is supplied.
//
// This package does not send the optional "remoteip" parameter.
//
// Example Usage
//
// Using the recaptcha package with the default http.Client:
//  func handlePage(w http.ResponseWriter, r *http.Request) {
//     re := recapcha.Recaptcha{"SECRET_KEY"}
//     response, err := re.CheckRequest(r, nil)
//     if err != nil {
//         panic(err.Error())
//     }
//     if response.Success {
//         print("Success")
//     }
//  }
//
// To run this package on Google's Appengine, simply invoke CheckRequest with the appropriate client:
//   client := urlfetch.Client(appengine.NewContext(r))
//   response, err := re.CheckRequest(r, client)
package recaptcha

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

// The URL of the Google reCaptcha API.
const apiURL string = "https://www.google.com/recaptcha/api/siteverify"

// Response represents the response of the Google reCaptcha API.
type Response struct {
	Success     bool      `json:"success"`
	ChallengeTs time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

// The Recaptcha struct is used to store the Secret token, and to invoke the CheckRequest method on.
type Recaptcha struct {
	Secret string
}

// CheckRequest will extract the captcha response from the http.Request and attempt to verify the response with the
// reCaptcha servers.
//
// If the http.Client argument is nil, the default http.DefaultClient will be used.
func (re Recaptcha) CheckRequest(r *http.Request, c *http.Client) (Response, error) {
	if r == nil {
		return Response{}, errors.New("HTTP Request is nil")
	}

	captcha := r.FormValue("g-recaptcha-response")

	if captcha == "" {
		return Response{}, errors.New("No captcha response found in request body")
	}

	if c == nil {
		c = http.DefaultClient
	}

	remoteip, _, _ := net.SplitHostPort(r.RemoteAddr)
	serverResponse, err := c.PostForm(apiURL, url.Values{
		"secret":   {re.Secret},
		"response": {captcha},
		"remoteip": {remoteip},
	})

	if err != nil {
		return Response{}, errors.New("Could not connect to Recaptcha service: " + err.Error())
	}

	defer serverResponse.Body.Close()
	body, err := ioutil.ReadAll(serverResponse.Body)

	if err != nil {
		return Response{}, errors.New("Could not read response from Recaptcha service: " + err.Error())
	}

	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return Response{}, errors.New("Invalid response from Recaptcha service: " + err.Error())
	}

	return resp, nil
}
