package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

// User is
type User struct {
	Username string `json:"username"`
}
type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未実装
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// 何らかの別のエラーが発生
		panic(err.Error())
	} else {
		// 成功。ラップされたハンドラを呼び出します
		h.next.ServeHTTP(w, r)
	}
}

// MustAuth は認証が必要とするページにauthHandlerを返します
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandlerはサードパーティへのログインの処理を受け持ちます。
// パスの形式: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	toml := getConfig()
	conf := &oauth2.Config{
		ClientID:     toml.Discord.ID,
		ClientSecret: toml.Discord.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discordapp.com/api/oauth2/authorize",
			TokenURL: "https://discordapp.com/api/oauth2/token",
		},
		RedirectURL: "http://localhost:8080/auth/callback",
		Scopes:      []string{"identify"},
	}
	switch action {
	case "login":
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		ctx := context.Background()
		query := r.URL.Query()
		code := query["code"][0]
		tok, err := conf.Exchange(ctx, code)
		if err != nil {
			log.Fatal(err)
		}
		client := conf.Client(ctx, tok)
		resp, err := client.Get("https://discordapp.com/api/users/@me")
		if err != nil {
			log.Fatalf("client get error")
		}
		defer resp.Body.Close()
		byteArray, _ := ioutil.ReadAll(resp.Body)
		var user User
		if err := json.Unmarshal(byteArray, &user); err != nil {
			log.Fatal(err)
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: user.Username,
			Path:  "/"})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
