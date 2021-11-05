package mysql_test

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
	"time"

	. "github.com/uccu/go-mysql"
	"github.com/uccu/go-stringify"
)

func getPool() (*DB, error) {
	dbpool, err := Open("mysql", "billiards:HkhD3wyYxKM2JLh5@tcp(60.205.184.251:3306)/billiards?charset=utf8mb4&parseTime=true&loc=Asia%2FChongqing")
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

	return dbpool.WithPrefix("b_"), nil
}

func TestCount(t *testing.T) {

	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	count := dbpool.GetOrm("users").Count()
	log.Println("TestCount", count)

}

func TestGetFields(t *testing.T) {

	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	e := []*string{}
	dbpool.GetOrm("users").Query("WHERE id>?", 45175).Dest(&e).GetFields("create_time")
	b, _ := json.Marshal(e)
	log.Println("TestGetFields", string(b))

}

func TestFields(t *testing.T) {

	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	type User struct {
		Id int64  `db:"id" json:"id"`
		W  string `db:"user_nicename"`
	}

	e := User{}
	err = dbpool.GetOrm("users").Field("id").Fields([]string{"id"}).Query("WHERE id>?", 45175).Dest(&e).FetchOne()
	if err != nil {
		t.Error(err)
	}
	b, _ := json.Marshal(e)
	log.Println("TestGetFields", string(b))

}

func TestGetFieldString(t *testing.T) {
	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	data := dbpool.GetOrm("users").Query("WHERE id=?", 45175).GetFieldString("user_nicename")
	log.Println("TestGetFieldString", data)
}

func TestGetFieldInt(t *testing.T) {
	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	data := dbpool.GetOrm("users").Query("WHERE id=?", 45175).GetFieldString("coin")
	log.Println("TestGetFieldInt", data)
}

func TestUpdate(t *testing.T) {
	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	data, err := dbpool.GetOrm("users").Query("SET coin=coin+1 WHERE id=?", 45175).Update()
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

	data, err := dbpool.GetOrm("faq").Query("SET name=?,content=?", 123, 345).Insert()
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

	err = dbpool.GetOrm("users").Field("id", "user_nicename").Dest(&user).Query("WHERE id>?", 45175).FetchOne()
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

	err = dbpool.GetOrm("users").Field("id", "user_nicename").Dest(&user).Query("WHERE id>?", 45175).Select()
	if err != nil {
		t.Error(err)
	}
	b, _ := json.Marshal(user)
	log.Println("TestSelect", string(b))

}

type A struct {
	Name string `db:"name"`
}

type B struct {
	*A
	Id int `db:"id"`
}

func Test5(t *testing.T) {

	dest := []*B{}
	val := stringify.GetReflectValue(dest).Type()
	if val.Kind() == reflect.Slice {
		val = val.Elem()
		for val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
	}

	if val.Kind() == reflect.Struct {
		fields := []string{}

		loopStruct(val, func(s reflect.StructField) {
			name := s.Tag.Get("db")
			if name != "" {
				fields = append(fields, name)
			}
		})

		log.Println(fields)
	}
}

func loopStruct(val reflect.Type, f func(s reflect.StructField)) {
	if val.Kind() != reflect.Struct {
		return
	}
	for k := 0; k < val.NumField(); k++ {
		ft := val.Field(k).Type
		for ft.Kind() == reflect.Ptr || ft.Kind() == reflect.Interface {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.Struct {
			loopStruct(ft, f)
		} else {
			f(val.Field(k))
		}
	}
}

func Test6(t *testing.T) {

	var i []interface{}

	i = append(i, 1)

	log.Println(i)

}

func TestWhereMap(t *testing.T) {

	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	type User struct {
		Id int64 `db:"id" json:"id"`
	}

	user := &User{}

	err = dbpool.GetOrm("user").Dest(user).WhereMap("id", 2).FetchOne()
	if err != nil {
		t.Error(err)
	}
	b, _ := json.Marshal(user)
	log.Println("TestWhereMap", string(b))

}

type UserSimple struct {
	Id   int    `dbwhere:"id"`
	Name string `dbset:"name"`
}

type User struct {
	*UserSimple
}

func TestPointer(t *testing.T) {

	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	user := new(User)
	user.UserSimple = new(UserSimple)
	user.Id = 57
	user.Name = "dde"

	dbpool.GetOrm("user").Dest(user).WhereStru(user).SetStru(user).Update()

	if user.Id != 57 {
		t.Error("id not right")
	}

}

func TestNilPointer(t *testing.T) {
	dbpool, err := getPool()
	if err != nil {
		t.Error(err)
	}

	orm := dbpool.GetOrm("order").RawFields(true)

	total := orm.GetFieldString("SUM(price_total)")

	log.Println(orm.Err())
	log.Println(total)
}
