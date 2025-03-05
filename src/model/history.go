package model

import "time"

type History struct {
	SpaceID    string    `json:"space_id"`
	GoodID     string    `json:"good_id"`
	GoodName   string    `json:"good_name"`
	BorrowUser string    `json:"borrow_user"`
	WhenBorrow time.Time `json:"when_borrow"`
	WhenReturn time.Time `json:"when_return"`
}

func CreateHistory(history *History) {
	db.Create(history)
}

func FindHistories(history *History) []History {
	var h []History
	db.Where(history).Find(&h)
	return h
}