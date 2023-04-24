# mysqltls
encapsulation of single table operation flow concept based on MySQL database

kernel uses sqlx package.

### USE
Your structure needs to implement this interface
``` go
// TableInterface
type TableInterface interface {
	TableName() string    // TableName returnes the table name
	PrimaryKey() []string // PrimaryKey returnes the table id columns
}
```

DEMO:
``` go
type UserInfo struct {
	Id   string `db:"id"`   // primary key id
	Name string `db:"name"` // name
	Age  int    `db:"age"`  // age
}

func (UserInfo) TableName() string {
    return "user_info"
}

func (UserInfo) PrimaryKey() []string {
    return []string{"id"}
}

func main() {
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "123456", "127.0.0.1:3306", "miajiodb")
    db := tls.MustOpen(dsn)
    en := db.Engine(UserInfo{})
    var userInfo UserInfo{}
    if err := en.FindById(&userInfo, "xxxx"); err != nil {
        panic(err)
    }
}
```