package model

import "time"

type Good struct {
	GoodID        string    `json:"good_id"`
	SpaceID		  string	`json:"space_id"`
	AddEmail      string    `json:"add_email"`
	GoodName      string    `json:"good_name"`
	CanBorrow     bool      `json:"can_borrow"`
	Status        bool      `json:"status"`
	Description   string    `json:"description"`
	WhoBorrowUid  string    `json:"who_borrow_uid"`
	WhoBorrowName string    `json:"who_borrow_name"`
	WhoReturnUid  string    `json:"who_return_uid"`
	WhoReturnName string    `json:"who_return_name"`
	WhenBorrow    time.Time `json:"when_borrow"`
}

func CreateGood(good *Good) {
	db.Create(good)
}

func FindGood(good *Good) Good {
	var g Good
	db.Where(good).First(&g)
	return g
}

func FindGoods(good *Good) []Good {
	var g []Good
	db.Where(good).Find(&g)
	return g
}
