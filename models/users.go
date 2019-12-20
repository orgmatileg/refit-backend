package models

import "time"

type User struct {
	ID          int
	FullName    string
	Email       string
	Password    string
	Image       string
	Height      int
	Gender      string
	CountryCode int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

// CreateUser ...
func CreateUser(m *User) (lastInsertedID int64, err error) {
	db := infrastructures.GetMySQLDBconn()
	q := fmt.Sprintf(`
		INSERT INTO %s
		(email, password, firstname, lastname, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, User{}.TableName())

	res, err := db.Exec(q,
		m.Email,
		m.Password,
		m.FirstName,
		m.LastName,
		m.CreatedAt)
	if err != nil {
		return lastInsertedID, err
	}

	lastInsertedID, err = res.LastInsertId()
	if err != nil {
		return lastInsertedID, err
	}

	return lastInsertedID, nil
}

// FindUserByEmail ...
func FindUserByEmail(email string) (user *User, err error) {
	q := `
	SELECT id, email, password, firstname, lastname, created_at
	FROM user 
	WHERE email = ?`
	db := infrastructures.GetMySQLDBconn()

	var mu User
	err = db.QueryRow(q, email).Scan(
		&mu.ID,
		&mu.Email,
		&mu.Password,
		&mu.FirstName,
		&mu.LastName,
		&mu.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &mu, nil
}
