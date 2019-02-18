package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
	"net/http"
)

func commandResponse(s slack.SlashCommand) (c int, params slack.Msg) {
	switch s.Command {
	case "/hello":
		params := slack.Msg{Text: "Hello"}

		return http.StatusOK, params

	default:
		return http.StatusInternalServerError, slack.Msg{}
	}
}

func main() {
	var verificationToken string

	// フラグ解析
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

		// コマンドに応じてレスポンス
		c.JSON(commandResponse(s))
	})

	fmt.Println("[INFO] Server Listening")
	r.Run(":3000")
}
