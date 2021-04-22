package mysql_test

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/uccu/go-mysql"
)

func getPool() (*DB, error) {
	dbpool, err := Open("mysql", "")
	if err != nil {
		return nil, err
	}
	dbpool.SetMaxOpenConns(1)
	dbpool.SetMaxIdleConns(1)
	dbpool.SetConnMaxLifetime(2 * time.Second)
	err = dbpool.Ping()
	if err != nil {
		return nil, err
	}

	return dbpool, nil
}

func TestCount(t *testing.T) {

	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	count := dbpool.Table("cmf_users").Count()
	log.Println("TestCount", count)

}

func TestGetFields(t *testing.T) {

	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	e := []*string{}
	dbpool.Table("cmf_users").Query("WHERE id>?", 45175).Dest(&e).GetFields("create_time")
	b, _ := json.Marshal(e)
	log.Println("TestGetFields", string(b))

}

func TestGetFieldString(t *testing.T) {
	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	data := dbpool.Table("cmf_users").Query("WHERE id=?", 45175).GetFieldString("user_nicename")
	log.Println("TestGetFieldString", data)
}

func TestGetFieldInt(t *testing.T) {
	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	data := dbpool.Table("cmf_users").Query("WHERE id=?", 45175).GetFieldString("coin")
	log.Println("TestGetFieldInt", data)
}

func TestUpdate(t *testing.T) {
	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	data, err := dbpool.Table("cmf_users").Query("SET coin=coin+1 WHERE id=?", 45175).Update()
	if err != nil {
		t.Error(err)
	}
	log.Println("TestUpdate", data)
}

func TestInsert(t *testing.T) {
	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	data, err := dbpool.Table("cmf_faq").Query("SET name=?,content=?", 123, 345).Insert()
	if err != nil {
		t.Error(err)
	}
	log.Println("TestInsert", data)
}

func TestFetchOne(t *testing.T) {

	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	type User struct {
		Id int64  `db:"id" json:"id"`
		W  string `db:"user_nicename"`
	}

	user := &User{}

	err = dbpool.Table("cmf_users").Field("id", "user_nicename").Dest(&user).Query("WHERE id>?", 45175).FetchOne()
	if err != nil {
		t.Error(err)
	}
	b, _ := json.Marshal(user)
	log.Println("TestFetchOne", string(b))

}

func TestSelect(t *testing.T) {

	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	type User struct {
		Id int64  `db:"id" json:"id"`
		W  string `db:"user_nicename"`
	}

	user := []*User{}

	err = dbpool.Table("cmf_users").Field("id", "user_nicename").Dest(&user).Query("WHERE id>?", 45175).Select()
	if err != nil {
		t.Error(err)
	}
	b, _ := json.Marshal(user)
	log.Println("TestSelect", string(b))

}
