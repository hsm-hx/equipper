package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nlopes/slack"
	"net/http"
	"strings"
	"time"
)

type Equip struct {
	Id      int
	Title   string
	Type    int
	Owner   string
	DueDate time.Time
	State   int
	Remark  string
}

func createDatabase(db *sql.DB) {
  _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS "EQUIPS" ("ID" INTEGER PRIMARY KEY, "TITLE" TEXT, "TYPE" INTEGER, "OWNER" TEXT, "DUE_DATE" TEXT DEFAULT CURRENT_DATE, "BORROWER" TEXT DEFAULT "", "STATE" INTEGER DEFAULT 0, "REMARK" TEXT DEFAULT "")`,
	)
  if err != nil {
    panic(err)
  }
}

func converseEquipType(s string) (n int) {
	switch s {
	case "BOOK":
		n = 1
	case "COMPUTER":
		n = 2
	case "SUPPLY":
		n = 3
	case "CABLE":
		n = 4
	default:
		n = 0
	}

	return
}

func parseAddText(s string) (e Equip) {
	a := strings.Split(s, " ")

	if len(a) <= 2 {
		a = append(a, "computer_club")
	}

	e = Equip{
		Title: a[0],
		Type:  converseEquipType(a[1]),
		Owner: a[2],
	}

	return
}

func addEquip(e Equip, db *sql.DB) (id int64) {
	res, err := db.Exec(
		`INSERT INTO EQUIPS (TITLE, TYPE, OWNER) VALUES (?, ?, ?)`,
		e.Title,
		e.Type,
		e.Owner,
	)
	if err != nil {
		panic(err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		panic(err)
	}

	return
}

func commandResponse(s slack.SlashCommand, db *sql.DB) (c int, params slack.Msg) {
	switch s.Command {
	case "/hello":
		params := slack.Msg{Text: "Hello"}

		return http.StatusOK, params

	case "/addEquip":
		e := parseAddText(s.Text)
		addEquip(e, db)

		params := slack.Msg{Text: "Added new Equipment: " + e.Title}

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

	// サーバー準備
	r := gin.Default()
	// DB準備
	db, err := sql.Open("sqlite3", "./equip.db")
	if err != nil {
		panic(err)
	}

  _, err = db.Exec(`SELECT * FROM EQUIPS`, )
  if err != nil {
    createDatabase(db)
    return
  }

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
		c.JSON(commandResponse(s, db))
	})

	fmt.Println("[INFO] Server Listening")
	r.Run(":3000")
}
