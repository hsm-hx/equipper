package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
	"net/http"
)

func main() {
	var verificationToken string

	flag.StringVar(&verificationToken, "token", "YOUR_VERIFICATION_TOKEN_HERE", "Your Slash Verification Token")
	flag.Parse()
	fmt.Println("Your slash verification token ->", verificationToken)

	r := gin.Default()
	r.POST("/cmd", func(c *gin.Context) {
		// コマンドをパースする
		s, err := slack.SlashCommandParse(c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		// トークン認証
		if !s.ValidateToken(verificationToken) {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		// コマンドの種類によってレスポンスを変える
		switch s.Command {
		case "/hello":
			// 返すメッセージ
      params := &slack.Msg{Text: "Hello"}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}

			c.JSON(200, params)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
	})

	fmt.Println("[INFO] Server Listening")
	r.Run(":3000")
}
