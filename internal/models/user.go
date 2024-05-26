package models

import (
	"database/sql"
	"fmt"

	"github.com/cjnghn/db-shard-example/internal/db"
	"github.com/google/uuid"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Shard string `json:"shard"`
}

func CreateUser(name string) (string, error) {
	userID := uuid.New().String()
	shard := db.GetShardByUUID(userID)

	query := "INSERT INTO users (id, name) VALUES (?, ?)"
	_, err := shard.DB.Exec(query, userID, name)
	if err != nil {
		return "", fmt.Errorf("failed to insert user into %s: %v", shard.Name, err)
	}

	return userID, nil
}

func GetUser(userID string) (string, error) {
	shard := db.GetShardByUUID(userID)

	query := "SELECT id, name FROM users WHERE id = ?"
	row := shard.DB.QueryRow(query, userID)

	var id string
	var name string
	err := row.Scan(&id, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user with ID %d not found in %s", userID, shard.Name)
		}
		return "", fmt.Errorf("failed to query user from %s: %v", shard.Name, err)
	}

	return name, nil
}

func GetAllUsers() ([]User, error) {
	var users []User
	query := "SELECT id, name FROM users"

	for _, shard := range db.Shards {
		rows, err := shard.DB.Query(query)
		if err != nil {
			return nil, fmt.Errorf("failed to query users from %s: %v", shard.Name, err)
		}
		defer rows.Close()

		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Name)
			if err != nil {
				return nil, fmt.Errorf("failed to read user data from %s: %v", shard.Name, err)
			}
			user.Shard = shard.Name
			users = append(users, user)
		}
	}

	return users, nil
}
