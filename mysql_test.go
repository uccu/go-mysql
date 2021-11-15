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

func getPool() *DB {
	dbpool, err := Open("mysql", "")
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

	return dbpool.WithPrefix("b_")
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
	assert.Equal(t, orm.Sql, "SELECT COUNT(`id`) FROM `b_user` LIMIT ?")

	orm = dbpool.Table("user")
	orm.Exec(false).Suffix(10).Count()
	assert.Equal(t, orm.Sql, "SELECT COUNT(1) FROM `b_user_1` LIMIT ?")
}

func TestGetFields(t *testing.T) {
	dbpool := getPool()
	orm := dbpool.Table("user")
	orm.Exec(false).GetFields("a")
	assert.Equal(t, orm.Sql, "SELECT `a` FROM `b_user`")
}

func TestGetField(t *testing.T) {
	dbpool := getPool()
	orm := dbpool.Table("user")
	orm.Exec(false).GetField("a")
	assert.Equal(t, orm.Sql, "SELECT `a` FROM `b_user` LIMIT ?")
}

func TestField(t *testing.T) {
	dbpool := getPool()
	orm := dbpool.Table("user")
	orm.Exec(false).Field("a", "w").FetchOne()
	assert.Equal(t, orm.Sql, "SELECT `a`, `w` FROM `b_user` LIMIT ?")

	orm = dbpool.Table("user")
	orm.Exec(false).Fields([]interface{}{"a", "b"}).FetchOne()
	assert.Equal(t, orm.Sql, "SELECT `a`, `b` FROM `b_user` LIMIT ?")

	orm = dbpool.Table("user")
	orm.Exec(false).Dest(user{Id: 1, Name: "name"}).FetchOne()
	assert.Equal(t, orm.Sql, "SELECT `id`, `name` FROM `b_user` LIMIT ?")
}

func TestQuery(t *testing.T) {
	dbpool := getPool()
	orm := dbpool.Table()
	orm.Exec(false).Query("UPDATE bb SET 1=1").Update()
	assert.Equal(t, orm.Sql, "UPDATE bb SET 1=1")

	orm = dbpool.Table("user")
	orm.Exec(false).Query("WHERE 1=1").Delete()
	assert.Equal(t, orm.Sql, "DELETE FROM `b_user` WHERE 1=1")

	orm = dbpool.Table()
	orm.Exec(false).Query("UPDATE bb").Query("SET a=1").Query("WHERE a=?", 1).FetchOne()
	assert.Equal(t, orm.Sql, "UPDATE bb SET a=1 WHERE a=? LIMIT ?")

}

func TestUpdate(t *testing.T) {
	dbpool := getPool()
	var orm *Orm

	orm = dbpool.Table("user")
	orm.Exec(false).Set("a", 1, "b", 3).Where("a", 1).Where("b", 3).Update()
	assert.Equal(t, orm.Sql, "UPDATE `b_user` SET `a`=?, `b`=? WHERE `a`=? AND `b`=?")

	orm = dbpool.Table("user")
	orm.Exec(false).Set(map[string]interface{}{"a": 2}).Where(map[string]interface{}{"a": 2}).Update()
	assert.Equal(t, orm.Sql, "UPDATE `b_user` SET `a`=? WHERE `a`=?")

	orm = dbpool.Table("user")
	orm.Exec(false).Set(user{Id: 1, Name: "name"}).Where(user{Id: 1, Name: "name"}).Update()
	assert.Equal(t, orm.Sql, "UPDATE `b_user` SET `name`=? WHERE `id`=?")
}

func TestInsert(t *testing.T) {
	dbpool := getPool()
	var orm *Orm

	orm = dbpool.Table("user")
	orm.Exec(false).Set("a", 1).Set("b", 3).Insert()
	assert.Equal(t, orm.Sql, "INSERT INTO `b_user` SET `a`=?, `b`=?")

	orm = dbpool.Table("user")
	orm.Exec(false).Set(map[string]interface{}{"a": 2}).Insert()
	assert.Equal(t, orm.Sql, "INSERT INTO `b_user` SET `a`=?")

	orm = dbpool.Table("user")
	orm.Exec(false).Set(user{Id: 1, Name: "name"}).Insert()
	assert.Equal(t, orm.Sql, "INSERT INTO `b_user` SET `name`=?")
}

func TestUnion(t *testing.T) {
	dbpool := getPool()

	orm := dbpool.Table("order")
	orm2 := dbpool.Table("order2")
	orm.Union(orm2).Exec(false).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `b_order` UNION (SELECT * FROM `b_order2`)")

	orm = dbpool.Table("order")
	orm.Union().Exec(false).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `b_order`")

	orm = dbpool.Table("order")
	orm2 = dbpool.Table("order2")
	orm3 := dbpool.Table("order3")
	orm.Union(orm2).UnionAll(orm3).Exec(false).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `b_order` UNION (SELECT * FROM `b_order2`) UNION ALL (SELECT * FROM `b_order3`)")

	orm = dbpool.Table("order")
	orm2 = dbpool.Table("order2")
	orm3 = dbpool.Table("order3")
	orm.Union(orm2, orm3).Exec(false).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `b_order` UNION (SELECT * FROM `b_order2`) UNION (SELECT * FROM `b_order3`)")

}

func TestGroup(t *testing.T) {
	dbpool := getPool()
	var orm *Orm

	orm = dbpool.Table("user")
	orm.Exec(false).Group("type").Having("a", 1).Having("b", 3).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `b_user` GROUP BY `type` HAVING `a`=? AND `b`=?")

	orm = dbpool.Table("user")
	orm.Exec(false).Group("type").Having(map[string]interface{}{"a": 2}).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `b_user` GROUP BY `type` HAVING `a`=?")

	orm = dbpool.Table("user")
	orm.Exec(false).Group("type").Having(user{Id: 1, Name: "name"}).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `b_user` GROUP BY `type` HAVING `id`=?")
}

func TestChain(t *testing.T) {
	dbpool := getPool()
	var orm *Orm

	orm = dbpool.Table()
	orm.Exec(false).Field("a").Table("user").Alias("u").Order("id desc").Limit(1, 2).Select()
	assert.Equal(t, orm.Sql, "SELECT `a` FROM `b_user` ORDER BY `id` DESC LIMIT ?,?")

	orm = dbpool.Table()
	orm.Exec(false).Field("e.a vv").Table("user u").Alias("u1").Table("level e").Order("u.id").Page(1, 2).Select()
	assert.Equal(t, orm.Sql, "SELECT `e`.`a` `vv` FROM `b_user` `u1`, `b_level` `e` ORDER BY `u`.`id` LIMIT ?,?")

	orm = dbpool.Table("user")
	orm.Exec(false).Page(0, 2).Select()
	assert.Equal(t, orm.Sql, "SELECT * FROM `b_user` LIMIT ?,?")

}

func TestJoin(t *testing.T) {
	dbpool := getPool()
	var orm *Orm

	orm = dbpool.Table()
	orm.Exec(false).Field("a.a").Table("user").Alias("a").LeftJoin("user b", Raw("ON a.id=b.id")).Select()
	assert.Equal(t, orm.Sql, "SELECT `a`.`a` FROM `b_user` `a` LEFT JOIN `b_user` `b` ON a.id=b.id")

	orm = dbpool.Table()
	orm.Exec(false).Field("b.c").Table("user a").RightJoin("user b", Mix("ON %t=%t", Field("a.id"), Field("b.id"))).Select()
	assert.Equal(t, orm.Sql, "SELECT `b`.`c` FROM `b_user` `a` RIGHT JOIN `b_user` `b` ON `a`.`id`=`b`.`id`")

	orm = dbpool.Table("user")
	orm.Exec(false).Field("user.c").RightJoin("goods", Mix("ON %t=%t", "user.id", "goods.user_id")).Where("user.id", 1).Select()
	assert.Equal(t, orm.Sql, "SELECT `user`.`c` FROM `b_user` RIGHT JOIN `b_goods` ON `user`.`id`=`goods`.`user_id` WHERE `user`.`id`=?")
}

func TestResult(t *testing.T) {
	dbpool := getPool()

	dbpool.Table("user").GetFieldInt("id")
	dbpool.Table("user").GetFieldString("id")
	dbpool.Table("user").GetFieldsInt("id")
	dbpool.Table("user").GetFieldsString("id")

	user := &struct {
		Id int64 `db:"id"`
	}{}

	dbpool.Table("user").Dest(user).FetchOne()
	dbpool.Table("user").Where(Raw("id=0")).Delete()
	dbpool.Table("user").Set(Raw("name=''")).Where(Raw("id=0")).Update()
	err := dbpool.Table("user").Where(Raw("id=0")).Dest(user).FetchOne()

	assert.Equal(t, err, ErrNoRows)
}
