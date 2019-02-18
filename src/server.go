package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"net/http"
  "flag"
)

func main() {
  var verificationToken string

  flag.StringVar(&verificationToken, "token", "YOUR_VERIFICATION_TOKEN_HERE", "Your Slash Verification Token")
  flag.Parse()
  fmt.Println("Your slash verification token ->", verificationToken)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    // コマンドをパースする
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

    // トークン認証
		if !s.ValidateToken(verificationToken) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

    // コマンドの種類によってレスポンスを変える
		switch s.Command {
		case "/hello":
      // 返すメッセージ
			params := &slack.Msg{Text: "Hello"}
			b, err := json.Marshal(params)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	fmt.Println("[INFO] Server Listening")
	http.ListenAndServe(":3000", nil)
}
