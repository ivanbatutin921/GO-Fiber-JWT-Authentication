package models

type Role string

const (
	ADMIN Role = "admin"
	USER  Role = "user"
)

type User struct {
	ID       int    `json:"id" gorm:"type:uuid;primary_key;autoIncrement"`
	Email    string `json:"email" gorm:"type:varchar(100);unique"`
	Password string `json:"password" gorm:"type:varchar(100)"`
	Role     string `json:"role" gorm:"type:varchar(50)"`
}

// Migration implements initializers.Table.
func (*User) Migration(data interface{}) error {
	panic("unimplemented")
}
