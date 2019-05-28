package main

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestConverseEquipType(t *testing.T) {
	e := Equip{}

	err := e.ConverseEquipType("BOOK")
	if e.Type != 1 || err != nil {
		t.Fatal("Failed test: BOOK is type 1, but output e.Type is ", e.Type)
	}

	err = e.ConverseEquipType("UNDEFINED")
	if err != EquipConverseError {
		t.Fatal("Failed test: UNDEFINED is type -: ", err)
	}
}

func TestUnconverseEquipState(t *testing.T) {
	e := Equip{}

	e.State = 0
	s, err := e.UnconverseEquipState()
	if s != "○" || err != nil {
		t.Fatal("Failed test: State 0 is not borrowing, but output is ", s)
	}

	e.State = 1
	s, err = e.UnconverseEquipState()
	if s != "×" || err != nil {
		t.Fatal("Failed test: State 1 is borrowing, but output is ", s)
	}
}

func TestParseAddText(t *testing.T) {
	s := "\"EQUIP NAME\" BOOK OWNER_NAME"
	e, err := parseAddText(s)

	expectEquip := Equip{
		Title: "EQUIP NAME",
		Type:  1,
		Owner: "OWNER_NAME",
	}
	if err != nil {
    t.Fatal("Failed test: ", err, "in case s:", s)
	}
	if e != expectEquip {
		t.Fatal("Failed test: unexpected e:", e)
	}

	s = "\"EQUIP_NAME\" OTHER"
	e, err = parseAddText(s)

	expectEquip = Equip{
		Title: "EQUIP_NAME",
		Type:  0,
		Owner: "computer_club",
	}
	if err != nil {
		t.Fatal("Failed test: ", err)
	}
	if e != expectEquip {
		t.Fatal("Failed test: unexpected e:", e)
	}
}

func TestAddEquip(t *testing.T) {
	db, err := sql.Open("sqlite3", "./equip.db")
	if err != nil {
		panic(err)
	}

	e := Equip{
		Title: "EQUIP_NAME",
		Type:  1,
		Owner: "computer_club",
	}

	id := addEquip(e, db)

	res := db.QueryRow("SELECT * FROM EQUIPS WHERE ID = ?", id)

	var (
		title    string
		eType    int
		owner    string
		due      string
		borrower string
		state    int
		remark   string
	)

	err = res.Scan(&id, &title, &eType, &owner, &due, &borrower, &state, &remark)

	if err == sql.ErrNoRows {
		t.Fatal("Failed test: Cannot insert equipment")
	}
	if err != nil {
		panic(err)
	}

	if title != "EQUIP_NAME" || eType != 1 || owner != "computer_club" {
		t.Fatal("Failed test: Cannot insert equipment")
	}
}

func TestSelectEquipFromId(t *testing.T) {
	db, err := sql.Open("sqlite3", "./equip.db")
	if err != nil {
		panic(err)
	}

	id := 1
	e, err := selectEquipFromId(id, db)
	if err != nil {
		t.Fatal("Failed test: ", err)
	}
	if e.Id != 1 {
		t.Fatal("Failed test: Cannot select equipment")
	}

	id = 99
	e, err = selectEquipFromId(id, db)
	if err != sql.ErrNoRows {
		t.Fatal("Failed test: ", err)
	}
}

func TestBorrowEquip(t *testing.T) {
	db, err := sql.Open("sqlite3", "./equip.db")
	if err != nil {
		panic(err)
	}

	// 備品が借りられるかテスト
	e, err := borrowEquip(1, "user", 14, db)

	const layout = "2006-01-02"
	today := time.Now().AddDate(0, 0, 14).Format(layout)

	if err != nil {
		t.Fatal("Failed test: ", err)
	}
	if e.State != 1 || e.Borrower != "user" {
		t.Fatal("Failed test: Cannot borrow equipment")
	}
	if e.DueDate != today {
		t.Fatal("Failed test: Expected", today, "but result is", e.DueDate)
	}

	// 存在しない備品は借りられない
	e, err = borrowEquip(99, "user", 14, db)
	if err != BorrowEquipError {
		t.Fatal("Failed test: ", err)
	}

	// 借りた上から借りられない
	e, err = borrowEquip(1, "user", 14, db)
	if err != BorrowEquipError {
		t.Fatal("Failed test: ", err)
	}
}

func TestReturnEquip(t *testing.T) {
	db, err := sql.Open("sqlite3", "./equip.db")
	if err != nil {
		panic(err)
	}

	err = returnEquip(1, "user", db)
	if err != nil {
		t.Fatal("Failed test: ", err)
	}

	e, err := selectEquipFromId(1, db)
	if err != nil {
		t.Fatal("Failed test: ", err)
	}
	if e.State != 0 {
		t.Fatal("Failed test: Cannot return equipment")
	}
	if e.DueDate != "" {
		t.Fatal("Failed test: Expected DueDate is", "", "but it is", e.DueDate)
	}
}

func TestDeleteEquip(t *testing.T) {
	id := 0

	db, err := sql.Open("sqlite3", "./equip.db")
	if err != nil {
		panic(err)
	}

	deleteEquip(id, db)

	var (
		title    string
		eType    int
		owner    string
		due      string
		borrower string
		state    int
		remark   string
	)

	res := db.QueryRow(`SELECT * FROM EQUIPS WHERE ID = ?`, id)

	err = res.Scan(&id, &title, &eType, &owner, &due, &borrower, &state, &remark)

	if err != sql.ErrNoRows {
		t.Fatal("Failed test: Cannot delete equipment")
	}
}
