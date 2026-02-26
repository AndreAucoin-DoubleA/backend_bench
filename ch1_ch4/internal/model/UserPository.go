package model

import (
	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

type UserRepository struct {
	Session *gocql.Session
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user User

	err := r.Session.Query(`
        SELECT id, email, password_hash, created_at
        FROM user_by_email
        WHERE email = ?`,
		email,
	).Scan(&user.UserID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
