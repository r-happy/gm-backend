package model

type BorrowUser struct {
	GoodID string `json:"good_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

func CreateBorrowUser(borrowUser *BorrowUser) {
	db.Create(borrowUser)
}

func FindBorrowUser(borrowUser *BorrowUser) BorrowUser {
	var b BorrowUser
	db.Where(borrowUser).First(&b)
	return b
}

func FindBorrowUsers(borrowUser *BorrowUser) []BorrowUser {
	var b []BorrowUser
	db.Where(borrowUser).Find(&b)
	return b
}

func RemoveBorrowUser(good_id string) {
	db.Where("good_id = ?", good_id).Delete(&BorrowUser{})
}
