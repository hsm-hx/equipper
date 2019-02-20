package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestConverseEquipType(t *testing.T) {
	n := converseEquipType("BOOK")
	if n != 1 {
		t.Fatal("Failed test: BOOK is type 1")
	}

	n = converseEquipType("COMPUTER")
	if n != 2 {
		t.Fatal("Failed test: COMPUTER is type 2")
	}

	n = converseEquipType("SUPPLY")
	if n != 3 {
		t.Fatal("Failed test: SUPPLY is type 3")
	}

	n = converseEquipType("CABLE")
	if n != 4 {
		t.Fatal("Failed test: CABLE is type 4")
	}

	n = converseEquipType("OTHER")
	if n != 0 {
		t.Fatal("Failed test: OTHER is type 0")
	}
}

func TestParseAddText(t *testing.T) {
	s := "EQUIP_NAME BOOK OWNER_NAME"
	e := parseAddText(s)

	expectEquip := Equip{
		Title: "EQUIP_NAME",
		Type:  1,
		Owner: "OWNER_NAME",
	}
	if e != expectEquip {
		t.Fatal("Failed test: parseAddText")
	}

	s = "EQUIP_NAME OTHER"
	e = parseAddText(s)

	expectEquip = Equip{
		Title: "EQUIP_NAME",
		Type:  0,
		Owner: "computer_club",
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
	e := selectEquipFromId(id, db)
	if e.Id != 1 {
		t.Fatal("Failed test: Cannot select equipment")
	}
}

func TestBorrowEquip(t *testing.T) {
	db, err := sql.Open("sqlite3", "./equip.db")
	if err != nil {
		panic(err)
	}

	e := borrowEquip(1, "user", db)

	if e.State != 1 || e.Borrower != "user" {
		t.Fatal("Failed test: Cannot borrow equipment")
	}
}

func TestReturnEquip(t *testing.T) {
	db, err := sql.Open("sqlite3", "./equip.db")
	if err != nil {
		panic(err)
	}

	str := returnEquip(1, "user", db)
	if str != "" {
		t.Fatal("Failed test: Cannot return equipment")
	}

	e := selectEquipFromId(1, db)
	if e.State != 0 || e.Borrower != "" {
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
