package midwr

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
	tmpl, err := template.ParseFiles("../login.html")
	helper.Errors(err, "templateparsefiles(renderlogin)")
	tmpl.Execute(w, nil)
}

func BuildAuth(w http.ResponseWriter, r *http.Request) []byte {
	r.ParseForm()
	auth := models.Auth{
		Password: strings.TrimSpace(r.FormValue("password")),
		Email:    strings.TrimSpace(r.FormValue("email")),
	}
	return models.Encode(auth)
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
	response := Requester(r, http.MethodPost, "http://localhost:8383/authorize/?", BuildAuth(w, r))
	res := models.ResponseAuth{}
	json.NewDecoder(response.Body).Decode(&res)
	if res.Error == "invld" || res.Message == "" {
		Render_login(w, r)
		w.Write([]byte("\nInvalid password or email"))
		return
	}
	w.Write([]byte(res.Message))
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
