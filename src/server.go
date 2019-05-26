package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	BorrowEquipError = errors.New("BorrowEquipError")
	numType          = map[string]int{"BOOK": 1, "COMPUTER": 2, "SUPPLY": 3, "CABLE": 4, "OTHER": 0}
)

func createDatabase(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS "EQUIPS" ("ID" INTEGER PRIMARY KEY, "TITLE" TEXT, "TYPE" INTEGER, "OWNER" TEXT, "DUE_DATE" TEXT DEFAULT CURRENT_DATE, "BORROWER" TEXT DEFAULT "", "STATE" INTEGER DEFAULT 0, "REMARK" TEXT DEFAULT "")`,
	)
	if err != nil {
		panic(err)
	}
}

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

func (e *Equip) ConverseEquipType(s string) (err error) {
	var ok bool
	e.Type, ok = numType[s]
	if ok != true {
		err = BorrowEquipError
	}

	return
}

func (e Equip) UnconverseEquipType() (s string, err error) {
	switch e.Type {
	case 1:
		s = "BOOK"
	case 2:
		s = "COMPUTER"
	case 3:
		s = "SUPPLY"
	case 4:
		s = "CABLE"
	case 0:
		s = "OTHER"
	default:
		err = BorrowEquipError
		s = ""
	}

	return
}

func (e Equip) UnconverseEquipState() (s string, err error) {
	switch e.State {
	case 0:
		s = "○"
	case 1:
		s = "×"
	default:
		err = BorrowEquipError
		s = ""
	}

	return
}

func parseAddText(s string) (e Equip, err error) {
	a := strings.Split(s, " ")

	// 所有者の指定がなければ
	if len(a) <= 2 {
		a = append(a, "computer_club")
	}

	e = Equip{
		Title: a[0],
		Type:  99,
		Owner: a[2],
	}

	err = e.ConverseEquipType(a[1])
	if err == BorrowEquipError {
		return
	}

	return
}

func selectAllEquips(db *sql.DB) (equips []Equip) {
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

func selectEquipFromId(id int, db *sql.DB) (e Equip, err error) {
	res := db.QueryRow(`SELECT * FROM EQUIPS WHERE ID = ?`, id)

	err = res.Scan(&e.Id, &e.Title, &e.Type, &e.Owner, &e.DueDate, &e.Borrower, &e.State, &e.Remark)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
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

func borrowEquip(id int, name string, due int, db *sql.DB) (e Equip, err error) {
	e, err = selectEquipFromId(id, db)
	if e.State == 1 || err == sql.ErrNoRows {
		err = BorrowEquipError
		return
	} else if err != nil {
		panic(err)
	}

	const layout = "2006-01-02"
	date := time.Now().AddDate(0, 0, due).Format(layout)

	_, err = db.Exec(`UPDATE EQUIPS SET STATE = 1, BORROWER = ?, DUE_DATE = ? WHERE ID = ?`,
		name,
		date,
		id,
	)
	if err != nil {
		panic(err)
	}

	e, err = selectEquipFromId(id, db)
	return
}

func returnEquip(id int, name string, db *sql.DB) (err error) {
	e, err := selectEquipFromId(id, db)
	if e.State != 1 || e.Borrower != name {
		err = BorrowEquipError
		return
	}
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		panic(err)
	}

	_, err = db.Exec(`UPDATE EQUIPS SET STATE = 0, BORROWER = ?, DUE_DATE = '' WHERE ID = ?`,
		"",
		id,
	)
	if err != nil {
		panic(err)
	}

	return
}

func commandResponse(s slack.SlashCommand, due int, db *sql.DB) (c int, params slack.Msg) {
	switch s.Command {
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
		e, err := selectEquipFromId(id, db)
		if err == BorrowEquipError {
			params := slack.Msg{
				Text: "id" + s.Text + "の備品は存在しません",
			}

			return http.StatusOK, params
		} else if err != nil {
			panic(err)
		}

		deleteEquip(id, db)

		params := slack.Msg{Text: "備品を削除しました: " + e.Title}

		return http.StatusOK, params

	case "/equipborrow":
		id, _ := strconv.Atoi(s.Text)
		e, err := borrowEquip(id, s.UserName, due, db)

		if err == sql.ErrNoRows {
			params := slack.Msg{
				Text: "id" + s.Text + "の備品は存在しません",
			}
			return http.StatusOK, params
		}

		params := slack.Msg{
			Text:         s.UserName + "が" + e.Title + "を貸出しました",
			ResponseType: "in_channel",
		}

		return http.StatusOK, params

	case "/equipreturn":
		id, _ := strconv.Atoi(s.Text)
		e, err := selectEquipFromId(id, db)
		if err == BorrowEquipError {
			params := slack.Msg{
				Text: "id" + s.Text + "の備品は存在しません",
			}
			return http.StatusOK, params
		}

		err = returnEquip(id, s.UserName, db)

		var params slack.Msg
		if err != nil {
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
	var borrowDue int

	// フラグ解析
	flag.StringVar(&verificationToken, "token", "YOUR_VERIFICATION_TOKEN_HERE", "Your Slash Verification Token")
	flag.IntVar(&borrowDue, "due", 14, "Your team's lending period")
	flag.Parse()
	fmt.Println("Your slash verification token ->", verificationToken)
	fmt.Println("Your team's lending period ->", borrowDue)

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
		e := selectAllEquips(db)

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
		c.JSON(commandResponse(s, borrowDue, db))
	})

	fmt.Println("[INFO] Server Listening")
	r.Run(":3000")
}
