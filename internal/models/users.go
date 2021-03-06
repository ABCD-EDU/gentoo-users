package models

import (
	"fmt"
	"time"
)

type UserSchema struct {
	UserId    string           `form:"user_id" json:"user_id" xml:"user_id"  binding:"required"`
	UserInfo  UserRegistration `form:"user_info" json:"user_info" xml:"user_info"  binding:"required"`
	CreatedOn time.Time        `form:"created_on" json:"created_on" xml:"created_on"  binding:"required"`
	CanPost   bool             `form:"can_post" json:"can_post" xml:"can_post"  binding:"required"`
	IsAdmin   bool             `form:"is_admin" json:"is_admin" xml:"is_admin"`
}

type UserRegistration struct {
	Email       string `form:"email" json:"email" xml:"email"  binding:"required"`
	Username    string `form:"username" json:"username" xml:"username"  binding:"required"`
	GooglePhoto string `form:"google_photo" json:"google_photo" xml:"google_photo"  binding:"required"`
	Description string `form:"description" json:"description" xml:"description"  binding:"required"`
}

// HIGH PRIORITY
func RegisterUser(userInfo UserRegistration) (UserSchema, error) {
	var userRegistered UserSchema
	sqlQuery := `
		INSERT INTO users (email, username, google_photo, description, created_on, can_post, is_admin)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING user_id
	`

	createdOn := time.Now()
	var userId string
	if err := db.QueryRow(
		sqlQuery,
		userInfo.Email,
		userInfo.Username,
		userInfo.GooglePhoto,
		userInfo.Description,
		createdOn,
		true,
		false,
	).Scan(&userId); err != nil {
		fmt.Println("SOMETHING WENT WRONG WITH WRITING TO DB AND GETTING THE ID")
		fmt.Println(err)
		return userRegistered, err
	}

	userRegistered = UserSchema{UserId: userId, UserInfo: userInfo, CreatedOn: createdOn, CanPost: true, IsAdmin: false}

	return userRegistered, nil
}

// HIGH PRIORITY
func GetUserInfo(column string, param string) (*UserSchema, error) {
	userInfo := new(UserSchema)
	sqlQuery := fmt.Sprintf(`
	SELECT * FROM users
	WHERE %s=$1
	LIMIT 1
	`, column)

	rows, err := db.Query(sqlQuery, param)
	if err != nil {
		fmt.Println(err)
		return userInfo, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, email, username, googlePhoto, description string
		var createdOn time.Time
		var canPost, isAdmin bool
		err := rows.Scan(&id, &username, &email, &description, &googlePhoto, &createdOn, &canPost, &isAdmin)
		if err != nil {
			fmt.Println(err)
			return userInfo, err
		}

		userReg := &UserRegistration{Email: email, Username: username, GooglePhoto: googlePhoto, Description: description}
		user := &UserSchema{id, *userReg, createdOn, canPost, isAdmin}

		userInfo = user
	}

	return userInfo, nil
}

func MuteUser(userId string) error {
	sqlQuery := `
		UPDATE users
		SET can_post=false
		WHERE user_id=$1;
		`

	_, err := db.Exec(sqlQuery, userId)
	if err != nil {
		return err
	}

	return nil
}

func BanUser(userId string) error {
	sqlQuery := `
		DELETE FROM users
		WHERE user_id=$1;
		`

	_, err := db.Exec(sqlQuery, userId)
	if err != nil {
		return err
	}

	return nil
}
