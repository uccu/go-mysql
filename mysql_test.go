package mysql_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	. "github.com/uccu/go-mysql"
	"github.com/uccu/go-stringify"
)

type user struct {
	Id   int64  `db:"id" dbset:"-"`
	Name string `db:"name" dbwhere:"-"`
}

var db *DB

func getPool() *DB {

	if db != nil {
		return db
	}

	dbpool, err := Open("mysql", "test2:53CNrrnY6N6Tk2yr@tcp(60.205.184.251:3306)/test2?charset=utf8mb4&parseTime=true&loc=Asia%2FChongqing")
	if err != nil {
		panic(err)
	}
	dbpool.SetMaxOpenConns(1)
	dbpool.SetMaxIdleConns(1)
	dbpool.SetConnMaxLifetime(2 * time.Second)
	err = dbpool.Ping()
	if err != nil {
		panic(err)
	}

	dbpool.WithErrHandler(func(e error, o *Orm) {
		// log.Println("sql: ", o.Sql, o.GetArgs())
		// log.Println(fmt.Sprintln(e))
	})
	dbpool.WithAfterQueryHandler(func(o *Orm) {
		// log.Println("sql: ", o.Sql, o.GetArgs())
	})

	dbpool.WithSuffix(func(i interface{}) string {
		return "_" + stringify.ToString(stringify.ToInt(i)%3)
	})

	db = dbpool.WithPrefix("t_")
	return db
}

func TestPool(t *testing.T) {
	_, err := Open("", "")
	assert.NotNil(t, err)
}

func TestCount(t *testing.T) {
	dbpool := getPool()
	var orm *Orm

	orm = dbpool.Table("user")
	orm.Exec(false).Count("id")
	assert.Equal(t, orm.Sql, "SELECT COUNT(`id`) FROM `t_user` LIMIT ?")

	orm = dbpool.Table("user")
	orm.Exec(false).Suffix(10).Count()
	assert.Equal(t, orm.Sql, "SELECT COUNT(1) FROM `t_user_1` LIMIT ?")

	orm = dbpool.Table("user")
	orm.Count()
	assert.Nil(t, orm.Err())
}

func TestField(t *testing.T) {
	dbpool := getPool()
	orm := dbpool.Table("user")
	orm.Exec(false).Field("a", "w").FetchOne()
	assert.Equal(t, orm.Sql, "SELECT `a`, `w` FROM `t_user` LIMIT ?")

	orm = dbpool.Table("user")
	orm.Exec(false).Fields([]interface{}{"a", "b"}).FetchOne()
	assert.Equal(t, orm.Sql, "SELECT `a`, `b` FROM `t_user` LIMIT ?")

	orm = dbpool.Table("user")
	orm.Exec(false).Dest(&user{Id: 1, Name: "name"}).FetchOne()
	assert.Equal(t, orm.Sql, "SELECT `id`, `name` FROM `t_user` LIMIT ?")
}

func TestQuery(t *testing.T) {
	dbpool := getPool()
	orm := dbpool.Table()
	orm.Exec(false).Query("UPDATE bb SET 1=1").Update()
	assert.Equal(t, orm.Sql, "UPDATE bb SET 1=1")

	orm = dbpool.Table("user")
	orm.Exec(false).Query("WHERE 1=1").Delete()
	assert.Equal(t, orm.Sql, "DELETE FROM `t_user` WHERE 1=1")

	orm = dbpool.Table()
	orm.Exec(false).Query("UPDATE bb").Query("SET a=1").Query("WHERE a=?", 1).FetchOne()
	assert.Equal(t, orm.Sql, "UPDATE bb SET a=1 WHERE a=? LIMIT ?")

}

func TestUpdate(t *testing.T) {
	dbpool := getPool()
	var orm *Orm

	orm = dbpool.Table("user")
	orm.Exec(false).Set("a", 1, "b", 3).Where("a", 1).Where("b", 3).Update()
	assert.Equal(t, orm.Sql, "UPDATE `t_user` SET `a`=?, `b`=? WHERE `a`=? AND `b`=?")

	orm = dbpool.Table("user")
	orm.Exec(false).Set(map[string]interface{}{"a": 2}).Where(map[string]interface{}{"a": 2}).Update()
	assert.Equal(t, orm.Sql, "UPDATE `t_user` SET `a`=? WHERE `a`=?")

	orm = dbpool.Table("user")
	orm.Exec(false).Set(user{Id: 1, Name: "name"}).Where(user{Id: 1, Name: "name"}).Update()
	assert.Equal(t, orm.Sql, "UPDATE `t_user` SET `name`=? WHERE `id`=?")
}

func TestInsert(t *testing.T) {
	dbpool := getPool()
	var orm *Orm

	orm = dbpool.Table("user")
	orm.Exec(false).Set("a", 1).Set("b", 3).Insert()
	assert.Equal(t, orm.Sql, "INSERT INTO `t_user` SET `a`=?, `b`=?")

	orm = dbpool.Table("user")
	orm.Exec(false).Set(map[string]interface{}{"a": 2}).Insert()
	assert.Equal(t, orm.Sql, "INSERT INTO `t_user` SET `a`=?")

	orm = dbpool.Table("user")
	orm.Exec(false).Set(user{Id: 1, Name: "name"}).Insert()
	assert.Equal(t, orm.Sql, "INSERT INTO `t_user` SET `name`=?")

	orm = dbpool.Table("user")
	orm.Exec(false).Set(user{Id: 1, Name: "name"}).Replace()
	assert.Equal(t, orm.Sql, "REPLACE INTO `t_user` SET `name`=?")
}

func TestUnion(t *testing.T) {
	dbpool := getPool()

	orm := dbpool.Table("order")
	orm2 := dbpool.Table("order2")
	orm.Union(orm2).Exec(false).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `t_order` UNION (SELECT * FROM `t_order2`)")

	orm = dbpool.Table("order")
	orm.Union().Exec(false).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `t_order`")

	orm = dbpool.Table("order")
	orm2 = dbpool.Table("order2")
	orm3 := dbpool.Table("order3")
	orm.Union(orm2).UnionAll(orm3).Exec(false).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `t_order` UNION (SELECT * FROM `t_order2`) UNION ALL (SELECT * FROM `t_order3`)")

	orm = dbpool.Table("order")
	orm2 = dbpool.Table("order2")
	orm3 = dbpool.Table("order3")
	orm.Union(orm2, orm3).Exec(false).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `t_order` UNION (SELECT * FROM `t_order2`) UNION (SELECT * FROM `t_order3`)")

}

func TestGroup(t *testing.T) {
	dbpool := getPool()
	var orm *Orm

	orm = dbpool.Table("user")
	orm.Exec(false).Group("type").Having("a", 1).Having("b", 3).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `t_user` GROUP BY `type` HAVING `a`=? AND `b`=?")

	orm = dbpool.Table("user")
	orm.Exec(false).Group("type").Having(map[string]interface{}{"a": 2}).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `t_user` GROUP BY `type` HAVING `a`=?")

	orm = dbpool.Table("user")
	orm.Exec(false).Group("type").Having(user{Id: 1, Name: "name"}).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `t_user` GROUP BY `type` HAVING `id`=?")
}

func TestChain(t *testing.T) {
	dbpool := getPool()
	var orm *Orm

	orm = dbpool.Table()
	orm.Exec(false).Field("a").Table("user").Alias("u").Order("id desc").Limit(1, 2).Select()
	assert.Equal(t, orm.Sql, "SELECT `a` FROM `t_user` ORDER BY `id` DESC LIMIT ?,?")

	orm = dbpool.Table()
	orm.Exec(false).Field("e.a vv").Table("user u").Alias("u1").Table("level e").Order("u.id").Page(1, 2).Select()
	assert.Equal(t, orm.Sql, "SELECT `e`.`a` `vv` FROM `t_user` `u1`, `t_level` `e` ORDER BY `u`.`id` LIMIT ?,?")

	orm = dbpool.Table("user")
	orm.Exec(false).Page(0, 2).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `t_user` LIMIT ?,?")

}

type userj struct {
	Id int64 `db:"a.id"`
}

func TestJoin(t *testing.T) {
	dbpool := getPool()
	var orm *Orm

	orm = dbpool.Table()
	orm.Join("user")
	assert.Equal(t, orm.Err(), ErrNoTable)

	orm.Exec(false).Field("a.a").Table("user").Alias("a").LeftJoin("user b", Raw("ON a.id=b.id")).Select()
	assert.Equal(t, orm.Sql, "SELECT `a`.`a` FROM `t_user` `a` LEFT JOIN `t_user` `b` ON a.id=b.id")

	orm = dbpool.Table()
	orm.Exec(false).Dest(&userj{}).Table("user").Alias("a").LeftJoin("user b", Raw("ON a.id=b.id")).Select()
	assert.Equal(t, orm.Sql, "SELECT `a`.`id` FROM `t_user` `a` LEFT JOIN `t_user` `b` ON a.id=b.id")

	orm = dbpool.Table()
	orm.Exec(false).Field("b.c").Table("user a").RightJoin("user b", Mix("ON %t=%t", Field("a.id"), Field("b.id"))).Select()
	assert.Equal(t, orm.Sql, "SELECT `b`.`c` FROM `t_user` `a` RIGHT JOIN `t_user` `b` ON `a`.`id`=`b`.`id`")

	orm = dbpool.Table("user")
	orm.Exec(false).Field("user.c").RightJoin("goods", Mix("ON %t=%t", "user.id", "goods.user_id")).Where("user.id", 1).Select()
	assert.Equal(t, orm.Sql, "SELECT `user`.`c` FROM `t_user` RIGHT JOIN `t_goods` ON `user`.`id`=`goods`.`user_id` WHERE `user`.`id`=?")

	orm = dbpool.Table("user u")
	orm.Exec(false).Join(Table("user_goods g"), Mix("ON %t=%t", Field("u.id"), Field("g.user_id"))).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `t_user` `u` JOIN `user_goods` `g` ON `u`.`id`=`g`.`user_id`")

	orm = dbpool.Table("user").Join(map[string]int{})
	assert.Equal(t, orm.Err(), ErrNoContainer)

	orm = dbpool.Table("user").Join("goods", "user.id")
	assert.Equal(t, orm.Err(), ErrOddNumberOfParams)

}

func TestGetField(t *testing.T) {

	var err error
	dbpool := getPool()
	var data int

	orm := dbpool.Table("user")
	orm.Exec(false).GetFields("a")
	assert.Equal(t, orm.Sql, "SELECT `a` FROM `t_user`")

	dbpool.Table("user").Exec(false).GetFieldInt("id")
	dbpool.Table("user").GetFieldString("id")

	err = dbpool.Table("user").Dest(data).GetField("id")
	assert.NotNil(t, err)

	err = dbpool.Table("user").Dest(&data).GetField("id-")
	assert.NotNil(t, err)

}

func TestGetFields(t *testing.T) {

	var err error
	dbpool := getPool()
	data := []int{}
	data2 := []*int{}

	orm := dbpool.Table("user")
	orm.Exec(false).GetField("a")
	assert.Equal(t, orm.Sql, "SELECT `a` FROM `t_user` LIMIT ?")

	dbpool.Table("user").Exec(false).GetFieldsInt("id")
	dbpool.Table("user").GetFieldsString("id")

	err = dbpool.Table("user").Dest(data).GetFields("id")
	assert.NotNil(t, err)

	err = dbpool.Table("user").Dest(&data).GetFields("id-")
	assert.NotNil(t, err)

	err = dbpool.Table("user").Dest(&data2).GetFields("id")
	assert.Nil(t, err)

}

func TestResult(t *testing.T) {

	var err error
	dbpool := getPool()

	dbpool.Table("user").Exec(false).Sum("id")
	dbpool.Table("user").Exec(false).SumFloat("id")
	dbpool.Table("user").Exec(false).Exist()

	user := &user{}

	_, err = dbpool.Table("user").Set(Raw("name=''")).Where(Raw("id=")).Update()
	assert.NotNil(t, err)

	_, err = dbpool.Table("user").Set(Raw("name=''")).Where(Raw("id=0")).Update()
	assert.Nil(t, err)

	err = dbpool.Table("user").Where(Raw("id=0")).Dest(user).FetchOne()
	assert.Equal(t, err, ErrNoRows)
}

type user2 struct {
	Id   int64
	Name string
}

type user3 struct {
	Id int64 `json:"id"`
	user4
}

type user4 struct {
	Name string `json:"name"`
}

func TestSelect(t *testing.T) {

	var err error
	var id int64
	dbpool := getPool()

	id, err = dbpool.Table("user").Set("name", "123").Insert()
	assert.Nil(t, err)

	users := []*user{}
	err = dbpool.Table("user").Dest(&users).Where(Raw("id>0")).Limit(1).Select()
	assert.Nil(t, err)

	err = dbpool.Table("user").Dest(&users).Where(Raw("id=")).Select()
	assert.NotNil(t, err)

	err = dbpool.Table("user").Dest(users).Where(Raw("id=1")).Select()
	assert.NotNil(t, err)

	u1 := &user{}
	err = dbpool.Table("user").Dest(u1).FetchOne()
	assert.Nil(t, err)

	err = dbpool.Table("user").Dest(u1).Where(Raw("id=")).FetchOne()
	assert.NotNil(t, err)

	u2 := &user2{}
	err = dbpool.Table("user").Dest(u2).Where("id", id).FetchOne()
	assert.Nil(t, err)

	u3 := &user3{}
	err = dbpool.Table("user").Dest(u3).Where("id", id).FetchOne()
	assert.Nil(t, err)

	var mp1 map[string]string
	err = dbpool.Table("user").Dest(&mp1).Where("id", id).FetchOne()
	assert.Nil(t, err)

	mp2 := map[string]string{}
	err = dbpool.Table("user").Dest(&mp2).Where("id", id).FetchOne()
	assert.Nil(t, err)

	mp := []map[string]string{}
	err = dbpool.Table("user").Dest(&mp).Where("id", id).Select()
	assert.Nil(t, err)

	var mp3 []map[string]string
	err = dbpool.Table("user").Dest(&mp3).Where("id", id).Select()
	assert.Nil(t, err)

	_, err = dbpool.Table("user").Where("id", id).Delete()
	assert.Nil(t, err)

}

func TestUtil(t *testing.T) {
	dbpool := getPool()
	orm := dbpool.Table("user").Dest("123")
	assert.Equal(t, orm.Err(), ErrNotPointer)
}
