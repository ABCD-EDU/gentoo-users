package models

import "log"

type UserSchema struct {
	user_id  string `json:"email"`
	userInfo UserRegistration
}

type UserRegistration struct {
	email        string `json:"email"`
	username     string `json:"username"`
	google_photo string `json:"google_photo"`
	description  string `json:"description"`
}

// HIGH PRIORITY
func RegisterUser(userInfo *UserRegistration) error {
	sqlQuery := `
		INSERT INTO users (email, username, google_photo, description)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.Exec(
		sqlQuery,
		userInfo.email,
		userInfo.username,
		userInfo.google_photo,
		userInfo.description,
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// HIGH PRIORITY
func GetUserInfo(email string) (*UserSchema, error) {
	userInfo := new(UserSchema)
	sqlQuery := `
		SELECT * FROM users
		WHERE email=$1
		LIMIT 1
	`

	rows, err := db.Query(sqlQuery, email)
	if err != nil {
		log.Fatal(err)
		return userInfo, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, email, username, googlePhoto, description string
		err := rows.Scan(&id, &email, &username, &googlePhoto, &description)
		if err != nil {
			log.Fatal(err)
			return userInfo, nil
		}

		userReg := &UserRegistration{email, username, googlePhoto, description}
		user := &UserSchema{id, *userReg}

		userInfo = user
	}

	return userInfo, nil
}

// LOW PRIORITY
func UpdateUserInfo(userInfo UserRegistration) error {
	return nil
}

// LOW PRIORITY
func DeleteUserInfo(email string) error {
	return nil
}
