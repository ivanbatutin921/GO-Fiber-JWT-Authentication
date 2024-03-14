package models

type Token struct {
	ID           int    `json:"id" gorm:"type:uuid;primary_key;autoIncrement"`
	RefreshToken string `json:"refresh_token"`
	UserID       int    `json:"user_id"`
	Expiry       int64  `json:"expiry"`
}

// Migration implements initializers.Table.
func (*Token) Migration(data interface{}) error {
	panic("unimplemented")
}
