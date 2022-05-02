package controllers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/udonetsm/help/helper"
	"github.com/udonetsm/help/models"
)

type Caller func()

func TestRequestToApi(w http.ResponseWriter, r *http.Request) {
	response := Requester(r, http.MethodPost, "http://localhost:8484/test?", nil)
	res := ParseResponseBody(response)
	ResponseError(res, func() { w.Write([]byte(res.Message)) }, func() { ShowLoginAndError(w, r, res.Error) })
}

func CheckSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("Auth")
		if err != nil {
			Render_login(w, r)
			return
		}
		r.Header.Add("Auth", token.Value)
		next.ServeHTTP(w, r)
	})
}

func Render_login(w http.ResponseWriter, r *http.Request) {
	MakeTempl(w, r, "../login.html", nil)
}

func Render_signup(w http.ResponseWriter, r *http.Request) {
	MakeTempl(w, r, "../signup.html", nil)
}

func MakeTempl(w http.ResponseWriter, r *http.Request, filename string, fields interface{}) {
	tmpl, err := template.ParseFiles(filename)
	helper.Errors(err, "templateparsefiles")
	tmpl.Execute(w, fields)
}

func Requester(r *http.Request, method, url string, body []byte) *http.Response {
	url = strings.TrimSpace(fmt.Sprintf("\n%s\n", url))
	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	helper.Errors(err, "httpnewrequest(requester)")
	request.Header.Add("Auth", r.Header.Get("Auth"))
	response, err := Client().Do(request)
	helper.Errors(err, "clientdo(requester)")
	return response
}

func Client() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func GetToken(w http.ResponseWriter, r *http.Request) {
	auth := models.Auth{}
	response := Requester(r, http.MethodPost, "http://localhost:8383/authorize?", auth.BuildAuth(w, r))
	res := ParseResponseBody(response)
	ResponseError(res, func() { SetCookieAuthAndRedirect(w, r, res.Message) }, func() { ShowLoginAndError(w, r, res.Error) })
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	user := models.AUser{}
	encode_user := user.BuildUser(w, r)
	response := Requester(r, http.MethodPost, "http://localhost:8383/newuser?", encode_user)
	res := ParseResponseBody(response)
	ResponseError(res, func() { SetcookieAndRedirect(w, r, res.Message) }, func() { ShowSignupAndError(w, r, res.Error) })
}

func SetcookieAndRedirect(w http.ResponseWriter, r *http.Request, message string) {
	SetCookie(w, "Auth", message)
	http.Redirect(w, r, "/test?", http.StatusMovedPermanently) //redirecting to /test? must be changed(it for tests)
}

func ShowSignupAndError(w http.ResponseWriter, r *http.Request, error string) {
	Render_signup(w, r)
	w.Write([]byte(error))
}

func SetCookieAuthAndRedirect(w http.ResponseWriter, r *http.Request, message string) {
	SetCookie(w, "Auth", message)
	http.Redirect(w, r, "/test?", http.StatusMovedPermanently) //redirecting to /test for tests (must be changed)
}

func ShowLoginAndError(w http.ResponseWriter, r *http.Request, error string) {
	Render_login(w, r)
	w.Write([]byte(error))
}

func SetCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Expires:  time.Now().AddDate(1, 0, 0),
	})
}

func ResponseError(res models.ResponseAuth, doIfOk Caller, DoIfNotOk Caller) {
	if res.Error != "" || res.Message == "" {
		DoIfNotOk()
		return
	}
	doIfOk()
}

func ParseResponseBody(respone *http.Response) models.ResponseAuth {
	res := models.ResponseAuth{}
	json.NewDecoder(respone.Body).Decode(&res)
	return res
}
