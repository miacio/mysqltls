package example

import (
	"time"

	"github.com/miacio/mysqltls/types"
)

type UserInfo struct {
	ID         int          `db:"id" json:"id" xml:"id"`                            // ID id 自增
	Name       *string      `db:"name" json:"name" xml:"name"`                      // Name name 用户名
	Password   *string      `db:"password" json:"password" xml:"password"`          // Password password 密码
	CreateTime *time.Time   `db:"create_time" json:"create_time" xml:"create_time"` // CreateTime create_time 创建时间
	UpdateTime *time.Time   `db:"update_time" json:"update_time" xml:"update_time"` // UpdateTime update_time 修改时间
	Sex        *types.IBool `db:"sex" json:"sex" xml:"sex"`                         // Sex sex 性别 true 男 false 女
}

// TableName UserInfo user_info
func (UserInfo) TableName() string {
	return "user_info"
}

func (UserInfo) PrimaryKey() []string {
	return []string{"id"}
}
