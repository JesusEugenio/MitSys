package models

type User struct {
	ID         int
	Email      string
	FullName   string
	Password   string
	IsAdmin    bool
	IsCashier  bool
	IsDelivery bool
}
