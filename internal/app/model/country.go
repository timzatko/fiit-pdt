package model

type Country struct {
	Id   int64  `gorm:"primarykey"`
	Code string `gorm:"index:,unique"`
	Name string `gorm:"index:,unique"`
}
