package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type player struct {
	Id          string
	NickName    string
	ShoutColor  string
	ShoutCount  int
	TotalScore  int
	SeasonScore int
	TotalRank   int
	SeasonRank  int
	Country     int
	Gender      string
	Authority   string
	IsVip       string
	Exist       bool
}

func newPlayer(_nickname string) *player {
	return &player{
		Id:          "",
		NickName:    _nickname,
		ShoutColor:  "",
		ShoutCount:  0,
		TotalScore:  0,
		SeasonScore: 0,
		TotalRank:   0,
		SeasonRank:  0,
		Country:     0,
		Gender:      "0",
		Authority:   "0",
		Exist:       false,
	}
}

func (user *player) dbconn() (*sql.DB, error) {
	db, connerr := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/gbwc")
	if connerr != nil {
		return nil, fmt.Errorf("cant connect to mysql server")
	}
	return db, nil
}

func (user *player) read() {

	db, connerr := user.dbconn()
	if connerr != nil {
		log.Fatal(connerr)
	} else {
		defer db.Close()

		row := db.QueryRow("SELECT Id, NickName, ShoutColor, ShoutCount, Gender, Authority FROM user WHERE NickName = ?", user.NickName[1:])
		qryerr := row.Scan(&user.Id, &user.NickName, &user.ShoutColor, &user.ShoutCount, &user.Gender, &user.Authority)

		if qryerr != nil && qryerr != sql.ErrNoRows {
			log.Fatal(qryerr)
		} else {
			user.Exist = true
		}
	}
}

func (user *player) readscore() {

	db, connerr := user.dbconn()
	if connerr != nil {
		log.Fatal(connerr)
	} else {
		defer db.Close()

		row := db.QueryRow("SELECT TotalScore, SeasonScore, TotalRank, SeasonRank, Country FROM game WHERE NickName = ?", user.NickName[1:])
		qryerr := row.Scan(&user.TotalScore, &user.SeasonScore, &user.TotalRank, &user.SeasonRank, &user.Country)

		if qryerr != nil && qryerr != sql.ErrNoRows {
			log.Fatal(qryerr)
		} else {
			user.Exist = true
		}
	}
}

func (user *player) decshout() {

	db, connerr := user.dbconn()
	if connerr != nil {
		log.Fatal(connerr)
	} else {
		defer db.Close()

		cmd := "UPDATE user SET ShoutCount = ShoutCount - ? WHERE NickName = ?"
		res, err := db.Exec(cmd, 1, user.NickName)

		if err != nil {
			log.Fatal(err)
		}

		affect, err := res.RowsAffected()

		if err != nil || affect == 0 {
			log.Fatal(err)
		}
	}
}
