// 业务概念的层次，不一定要与数据库中的数据完全对应
package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string
	Nickname string
	Birthday time.Time
	AboutMe  string

	Ctime time.Time
}
