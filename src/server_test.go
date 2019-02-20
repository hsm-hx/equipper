package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestConverseEquipType(t *testing.T) {
	n, err := converseEquipType("BOOK")
	if n != 1 || err != nil {
		t.Fatal("Failed test: BOOK is type 1")
	}

	n, err = converseEquipType("COMPUTER")
	if n != 2 || err != nil {
		t.Fatal("Failed test: COMPUTER is type 2")
	}

	n, err = converseEquipType("SUPPLY")
	if n != 3 || err != nil {
		t.Fatal("Failed test: SUPPLY is type 3")
	}

	n, err = converseEquipType("CABLE")
	if n != 4 || err != nil {
		t.Fatal("Failed test: CABLE is type 4")
	}

	n, err = converseEquipType("OTHER")
	if n != 0 || err != nil {
		t.Fatal("Failed test: OTHER is type 0")
	}

	n, err = converseEquipType("UNDEFINED")
	if err != BorrowEquipError {
		t.Fatal("Failed test: UNDEFINED is type -: ", err)
	}
}

func TestParseAddText(t *testing.T) {
	s := "EQUIP_NAME BOOK OWNER_NAME"
	e, err := parseAddText(s)

	expectEquip := Equip{
		Title: "EQUIP_NAME",
		Type:  1,
		Owner: "OWNER_NAME",
	}
	if err != nil {
		t.Fatal("Failed test: ", err)
	}
	if e != expectEquip {
		t.Fatal("Failed test: parseAddText")
	}

	s = "EQUIP_NAME OTHER"
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
		t.Fatal("Failed test: parseAddText")
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

	e, err := borrowEquip(1, "user", db)

	if err != nil {
		t.Fatal("Failed test: ", err)
	}
	if e.State != 1 || e.Borrower != "user" {
		t.Fatal("Failed test: Cannot borrow equipment")
	}

	// 存在しない備品は借りられない
	e, err = borrowEquip(99, "user", db)
	if err != BorrowEquipError {
		t.Fatal("Failed test: ", err)
	}

	// 借りた上から借りられない
	e, err = borrowEquip(1, "user", db)
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
