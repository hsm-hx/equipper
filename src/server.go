package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Equip struct {
	Id       int
	Title    string
	Type     int
	Owner    string
	DueDate  string
	Borrower string
	State    int
	Remark   string
}

var (
	BorrowEquipError = errors.New("BorrowEquipError")
)

func createDatabase(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS "EQUIPS" ("ID" INTEGER PRIMARY KEY, "TITLE" TEXT, "TYPE" INTEGER, "OWNER" TEXT, "DUE_DATE" TEXT DEFAULT CURRENT_DATE, "BORROWER" TEXT DEFAULT "", "STATE" INTEGER DEFAULT 0, "REMARK" TEXT DEFAULT "")`,
	)
	if err != nil {
		panic(err)
	}
}

func converseEquipType(s string) (n int, err error) {
	switch s {
	case "BOOK":
		n = 1
	case "COMPUTER":
		n = 2
	case "SUPPLY":
		n = 3
	case "CABLE":
		n = 4
	case "OTHER":
		n = 0
	default:
		err = BorrowEquipError
		n = 99
	}

	return
}

func parseAddText(s string) (e Equip, err error) {
	a := strings.Split(s, " ")

	if len(a) <= 2 {
		a = append(a, "computer_club")
	}

	eType, err := converseEquipType(a[1])
	if err == BorrowEquipError {
		return
	}

	e = Equip{
		Title: a[0],
		Type:  eType,
		Owner: a[2],
	}

	return
}

func selectEquips(db *sql.DB) (equips []Equip) {
	rows, err := db.Query(
		`SELECT * FROM EQUIPS`,
	)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		var e Equip

		if err := rows.Scan(&e.Id, &e.Title, &e.Type, &e.Owner, &e.DueDate, &e.Borrower, &e.State, &e.Remark); err != nil {
			log.Fatal("rows.Scan()", err)
			return
		}

		equips = append(equips, e)
	}

	return
}

func selectEquipFromId(id int, db *sql.DB) (e Equip) {
	res := db.QueryRow(`SELECT * FROM EQUIPS WHERE ID = ?`, id)

	err := res.Scan(&e.Id, &e.Title, &e.Type, &e.Owner, &e.DueDate, &e.Borrower, &e.State, &e.Remark)
	if err != nil {
		panic(err)
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

func deleteEquip(id int, db *sql.DB) {
	_, err := db.Exec(`DELETE FROM EQUIPS WHERE ID = ?`, id)
	if err != nil {
		panic(err)
	}
}

func borrowEquip(id int, name string, db *sql.DB) (e Equip) {
	_, err := db.Exec(`UPDATE EQUIPS SET STATE = 1, BORROWER = ? WHERE ID = ?`,
		name,
		id,
	)
	if err != nil {
		panic(err)
	}

	e = selectEquipFromId(id, db)
	return
}

func returnEquip(id int, name string, db *sql.DB) string {
	e := selectEquipFromId(id, db)
	if e.State != 1 || e.Borrower != name {
		return "err"
	}

	_, err := db.Exec(`UPDATE EQUIPS SET STATE = 0, BORROWER = ? WHERE ID = ?`,
		"",
		id,
	)
	if err != nil {
		panic(err)
	}

	return ""
}

func commandResponse(s slack.SlashCommand, db *sql.DB) (c int, params slack.Msg) {
	switch s.Command {
	case "/hello":
		params := slack.Msg{Text: "Hello"}

		return http.StatusOK, params

	case "/equipadd":
		e, err := parseAddText(s.Text)

		if err == BorrowEquipError {
			params := slack.Msg{
				Text: "TYPEはBOOK, COMPUTER, SUPPLY, CABLE, OTHERより選択してください",
			}
			return http.StatusOK, params
		} else if err != nil {
			panic(err)
		}

		addEquip(e, db)

		params := slack.Msg{Text: "新しい備品を追加しました: " + e.Title}
		return http.StatusOK, params

	case "/equipdelete":
		id, _ := strconv.Atoi(s.Text)
		e := selectEquipFromId(id, db)

		deleteEquip(id, db)

		params := slack.Msg{Text: "備品を削除しました: " + e.Title}

		return http.StatusOK, params

	case "/equipborrow":
		id, _ := strconv.Atoi(s.Text)
		e := borrowEquip(id, s.UserName, db)

		params := slack.Msg{
			Text:         s.UserName + "が" + e.Title + "を貸出しました",
			ResponseType: "in_channel",
		}

		return http.StatusOK, params

	case "/equipreturn":
		id, _ := strconv.Atoi(s.Text)
		e := selectEquipFromId(id, db)
		str := returnEquip(id, s.UserName, db)

		var params slack.Msg
		if str != "" {
			params = slack.Msg{
				Text: e.Title + "は現在貸出中でない，またはあなた以外が貸出中です",
			}
		} else {
			params = slack.Msg{
				Text:         s.UserName + "が" + e.Title + "を返却しました",
				ResponseType: "in_channel",
			}
		}

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
	_, err = db.Exec(`SELECT * FROM EQUIPS`)
	if err != nil {
		createDatabase(db)
		return
	}
	// HTML準備
	r.LoadHTMLGlob("template/*.tmpl")

	r.GET("/equip", func(c *gin.Context) {
		e := selectEquips(db)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"equips": e,
		})
	})
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
