[![CircleCI](https://circleci.com/gh/uccu/go-mysql/tree/master.svg?style=svg)](https://circleci.com/gh/uccu/go-mysql/tree/master)
[![Maintainability](https://api.codeclimate.com/v1/badges/d36b2d73f54d02d99076/maintainability)](https://codeclimate.com/github/uccu/go-mysql/maintainability)
[![codecov](https://codecov.io/gh/uccu/go-mysql/branch/master/graph/badge.svg?token=NFBBXRMOEO)](https://codecov.io/gh/uccu/go-mysql)
[![GitHub issues](https://img.shields.io/github/issues/uccu/go-mysql)](https://github.com/uccu/go-mysql/issues)
![GitHub](https://img.shields.io/github/license/uccu/go-mysql)

### 前言

本项目是基于[github.com/go-sql-driver/mysql ](http://github.com/go-sql-driver/mysql)开发，扩展了链式查询的工具。


### 1 连接数据库

```go
db, err := mysql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
if err != nil {
    panic(err)
}
```
#### 快速启动
```go

type User struct{
    Id      int `db:"id" dbset:"-"`
    Name    string `db:"name" dbwhere:"-"`
}

db.WithPrefix("t_")

var user User
db.Table("user").Where("id", 2).Dest(&user).FetchOne()
// sql: SELECT `id`, `name` FROM `t_user` WHERE `id`=? LIMIT ?
// args: [2, 1]

var users []*User
db.Table("user").Where("id", 2).Dest(&users).Page(2, 10).Select()
// sql: SELECT `id`, `name` FROM `t_user` WHERE `id`=?
// args: [2, 10, 10]

user2 := &User{
    Id: 2,
    Name: "kitty",
}
db.Table("user").Set(user2).Where(user2).Update()
// Set和Where务必按照顺序排列，不可颠倒
// sql: UPDATE `t_user` SET `name`=? WHERE `id`=?
// args: ["kitty", 2]
```




#### 1.1 配置
所有go-sql-driver/mysql的配置基本都能兼容，如
```go
db.SetMaxOpenConns(1)
db.SetMaxIdleConns(1)
db.SetConnMaxLifetime(2 * time.Second)
```

#### 1.2 扩展配置

添加表名前缀:
```go
db.WithPrefix("pre_")
```
根据传入的值添加表名后缀:
```go
db.WithPrefix(func(i interface{}) string {
    return "_" + strconv.FormatInt(int64(i.(int)%3), 10)
})
```
添加错误处理函数:
```go
db.WithErrHandler(func(e error, o *Orm) {
    fmt.Println("exec sql: ", o.Sql)
    fmt.Println(error.Error())
})
```
添加AfterQuery事件回调:
```go
db.WithAfterQueryHandler(func(*Orm)) {
    fmt.Println("exec sql: ", o.Sql)
})
```
获取生成构造器:
```go
omr := db.Table("table_user")
```
获取生成默认构造器:
```go
omr := db.Default()
```

### 2 构造器

#### 2.1 查询

##### 2.1.1 FetchOne
查询单条数据, eg:
```go
user := &struct{
    Id int64 `db:"id"`
    Name string `db:"name"`
}{}
db.Table("user").Where("id", 2).Dest(user).FetchOne()
```

##### 2.1.2 Select
查询多条数据, eg:
```go
list := []struct{
    Id int64 `db:"id"`
    Name string `db:"name"`
}{}
db.Table("user").Where("id", 2).Dest(&list).Select()
```

##### 2.1.3 GetField, GetFieldInt, GetFieldString
获取单个字段的值, eg:
```go
var name string := db.Table("user").Where("id", 2).GetFieldString("name")
```

##### 2.1.4 GetFields, GetFieldsInt, GetFieldsString
获取单个字段的值的列表, eg:
```go
var ids []int64 := db.Table("user").Order("id").GetFieldsInt("id")
```

##### 2.1.5 Count, Sum, SumFloat
其他聚合, 基于GetField封装的快捷方法，`Count`查询数量，`Sum`查询总数(int64)，`SumFloat`查询总数浮点数版(float64)
```go
var userCount int64 := db.Table("user").Count()
var costTotal int64 := db.Table("user_cost").SumFloat("cost")
```

#### 2.2 增删改

##### 2.2.1 Insert
新增
```go
insertId, err := db.Table("user").Set("name", "fffff").Insert()
insertId, err := db.Table("user").Set(map[string]interface{}{
    "name": "ffffff",
}).Insert()
insertId, err := db.Table("user").Query("SET name=?", "ffffff").Insert()
```

##### 2.2.2 Update
修改 **注意set和where需要严格顺序**
```go
aff, err := db.Table("user").Set("name", "fffff").Where("id", 2).Update()
```

##### 2.2.3 Delete
删除
```go
aff, err := db.Table("user").Where("id", 2).Delete()
```

#### 2.3 复杂组合
##### 2.3.1 Mix
拼接
```go
db.Table("user").Where(mysql.Mix("UNIX_TIMESTAMP(%t)>?", "create_at", time.Now().Unix())).Delete()
```
##### 2.3.2 Raw
原生
```go
db.Table("user").Where(mysql.Raw("id=2")).Delete()
```

#### 2.4 链式操作

##### 2.4.1 Where
条件筛选
```go
Where(Mix)
Where(map[string]value)
Where(key, value, key, value...)
Where(struct) // 优先dbwhere, 后db标签
```
##### 2.4.2 Set
设置，插入和更新都使用此方法
```go
Set(Mix)
Set(map[string]value)
Set(key, value, key, value...)
Set(struct) // 优先dbset, 后db标签
```
##### 2.4.3 Query
原生查询
```go
Query(sql, args...)
```

##### 2.4.4 Field, Fields
查询字段
```go
Field(Field...)
Field(string...)  // [表名.]名字[ 别名]
Fields([]Field)
Fields([]string)  // [表名.]名字[ 别名]
```

##### 2.4.5 Table
添加表
```go
Table(Table...)
Table(string...)  // [库名.]名字[ 别名]
```

##### 2.4.6 Group
分组
```go
Group(string...)  // [库名.]名字[ 别名][,...]
```

##### 2.4.7 Having
配合分组的筛选
```go
Having(Mix)
Having(map[string]value)
Having(key, value, key, value...)
Having(struct) // 优先dbwhere, 后db标签
```

##### 2.4.8 Order
排序
```go
Order(string...)  // [库名.]名字[ 别名][,...]
```

##### 2.4.9 Limit
控制数量
```go
Limit(length)
Limit(offset, length)  
```

##### 2.4.10 Page
分页
```go
Page(page, length)
```
##### 2.4.10 Alias
别名，当有至少一个表存在时才有效，且只会修改第一个表
```go
Alias(string)
```

##### 2.4.11 Join, LeftJoin, RightJoin
链表查询
```go
Join(Table/string, Mix)
eg:
db.Table("user").LeftJoin("goods", mysql.Mix("ON %t=%t", "user.id", "goods.user_id")).Where("user.id", 1).Select()
db.Table("bag b").Join("goods g", mysql.Mix("USING(%t)", "user_id")).Where("b.id", 1).Select()
```

##### 2.4.12 Union, UnionAll
合并
```go
Union(*Orm...)
eg:
o1 := db.Table("user").Where("id", 2)
o2 := db.Table("user").Where("id", 3).Union(o1).Select()
```

##### 2.4.13 Exec
设置是否执行sql，默认true
```go
Exec(bool)
eg:
o1 := db.Table("user").Where("id", 2).Exec(false).Select()
fmt.Println(o1.Sql)
```

##### 2.4.14 Err
返回错误
```go
Err() error
```

##### 2.4.15 GetArgs
返回所有变量
```go
GetArgs() []interface{}
```

##### 2.4.16 Dest
映射结构体，如果没有指定字段，则会去搜索结构体内db标签,如果非结构体如map类型, 则取*, 指定fields优先
```go
user := &struct{
    Id int64 `db:"id"`
    Name string `db:"name"`
}{}
db.Table("user").Where("id", 2).Dest(user).FetchOne()
```

#### 2.4 事务处理

##### 2.4.1 创建事务
```go
tx := dbpool.Start()
```

##### 2.4.2 回滚
```go
tx := dbpool.Start()
id, err = tx.Table("user").Set("name", "123").Insert()
if false {
    tx.Rollback()
}
```

##### 2.4.3 提交
```go
tx := dbpool.Start()
id, err = tx.Table("user").Set("name", "123").Insert()
if true {
    tx.Commit()
}
```

