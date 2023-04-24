package example_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/miacio/mysqltls/example"
	"github.com/miacio/mysqltls/tls"
	"github.com/miacio/mysqltls/types"
)

func PointString(val string) *string {
	return &val
}

func TestParserColumns(t *testing.T) {
	c := example.UserInfo{}
	columns := tls.ParserColumns(c, "db", true)
	fmt.Println(columns)
}

func TestFindById(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "miajio", "123456", "127.0.0.1:3306", "miajiodb")
	db, err := tls.Open(dsn, "db")
	if err != nil {
		t.Fatal(err)
	}
	var userInfo example.UserInfo
	err = db.Engine(example.UserInfo{}).FindById(&userInfo, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(*userInfo.Name)
}

func TestUpdateById(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "miajio", "123456", "127.0.0.1:3306", "miajiodb")
	db, err := tls.Open(dsn, "db")
	if err != nil {
		t.Fatal(err)
	}
	var userInfo example.UserInfo
	en := db.Engine(example.UserInfo{})
	err = en.FindById(&userInfo, 1)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	name := "这是我啊"
	fmt.Println(userInfo.Name)
	userInfo.UpdateTime = &now
	userInfo.Name = &name
	_, err = en.UpdateById(userInfo)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBatchInsert(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "miajio", "123456", "127.0.0.1:3306", "miajiodb")
	db, err := tls.Open(dsn, "db")
	if err != nil {
		t.Fatal(err)
	}
	params := make([]tls.TableInterface, 0)
	now := time.Now()
	bl := types.IBool(true)
	params = append(params, example.UserInfo{
		Name:       PointString("这是七名字"),
		Password:   PointString("MD5('123456')"),
		CreateTime: &now,
		Sex:        &bl,
	})

	en := db.Engine(example.UserInfo{})
	res, err := en.BatchInsert(params...)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.RowsAffected())
}

func TestDeleteById(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "miajio", "123456", "127.0.0.1:3306", "miajiodb")
	db, err := tls.Open(dsn, "db")
	if err != nil {
		t.Fatal(err)
	}
	res, err := db.Engine(example.UserInfo{}).DeleteById(1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.RowsAffected())
}

func TestSelect(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "miajio", "123456", "127.0.0.1:3306", "miajiodb")
	db, err := tls.Open(dsn, "db")
	if err != nil {
		t.Fatal(err)
	}
	var result []example.UserInfo
	err = db.Engine(example.UserInfo{}).Where().And("name like concat('%', ?, '%')", "这是").AndIn("id in (?)", []int{2, 4, 5}).And("sex = ?", types.IBool(false)).Find(&result)
	if err != nil {
		t.Fatal(err)
	}
	by, _ := json.Marshal(result)
	fmt.Println(string(by))

	res, err := db.Engine(example.UserInfo{}).Where().AndIn("id in (?)", []int{2, 4}).Delete()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

type ParserTagInStruct struct {
	Name   string  `db:"name"`
	Age    int     `db:"age"`
	Mobile *string `db:"mobile"`
}

type ParserTagStruct struct {
	Id       string                 `db:"id"`
	Children *ParserTagInStruct     `db:"children"`
	Params   map[string]interface{} `db:"params"`
	NowTime  time.Time              `db:"now_time"`
	LastTime *time.Time             `db:"last_time"`
}

func TestParserTagToMap(t *testing.T) {
	pc := ParserTagInStruct{
		Name:   "儿子",
		Age:    12,
		Mobile: nil,
	}
	now := time.Now()
	p := ParserTagStruct{
		Children: &pc,
		Params: map[string]interface{}{
			"age": 10,
			"sex": "男",
		},
		NowTime: now,
	}

	mp, err := tls.ParserTagToMap(p, "db")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mp)
}
