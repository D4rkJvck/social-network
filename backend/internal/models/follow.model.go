package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Follow struct {
	Id       int
	Followed *User
	Follower *User
	Archived bool
}

var ErrSelfFollow = errors.New("un utilisateur ne peut pas se suivre lui-même")

func (m *ConnDB) SetFollower(followerID, followedID int) error {
	if followerID == followedID {
		return ErrSelfFollow
	}
	// Vérifier d'abord si une relation archivée existe
	query := `
			SELECT id FROM follows
			WHERE id_follower = ? AND id_followed = ?
	`
	var followID int
	err := m.DB.QueryRow(query, followerID, followedID).Scan(&followID)

	if err == nil {
		updateQuery := `
					UPDATE follows
					SET archived = FALSE
					WHERE id = ?
			`
		_, err = m.DB.Exec(updateQuery, followID)
		if err != nil {
			return fmt.Errorf("error reactivating follow: %w", err)
		}
		return nil
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("error checking existing follow: %w", err)
	}

	insertQuery := `
			INSERT INTO follows (id_follower, id_followed, created_at, archived)
			VALUES (?, ?, CURRENT_TIMESTAMP, FALSE)
	`
	_, err = m.DB.Exec(insertQuery, followerID, followedID)
	if err != nil {
		return fmt.Errorf("error creating new follow: %w", err)
	}
	return nil
}

func (m *ConnDB) GetFollowers(userID int) ([]*Follow, error) {
	query := `
	SELECT f.id, f.archived, u.*
	FROM users u
	JOIN follows f ON u.id = f.id_follower
	WHERE f.id_followed = ? AND f.archived = FALSE
`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []*Follow

	for rows.Next() {
		f := &Follow{
			Follower: &User{},
		}
		var dateOfBirthFollowerStr string
		err := rows.Scan(
			&f.Id, &f.Archived, &f.Follower.Id, &f.Follower.Email, &f.Follower.Password, &f.Follower.FirstName, &f.Follower.LastName, &dateOfBirthFollowerStr, &f.Follower.ProfilePicture, &f.Follower.Nickname, &f.Follower.AboutMe, &f.Follower.Private, &f.Follower.CreatedAt)
		if err != nil {
			return nil, err
		}
		f.Follower.DateOfBirth, err = time.Parse("2006-01-02", dateOfBirthFollowerStr)
		if err != nil {
			return nil, err
		}
		followers = append(followers, f)
	}
	return followers, nil
}

func (m *ConnDB) GetFollowing(userID int) ([]*Follow, error) {
	query := `
	SELECT f.id, f.archived, u.*
	FROM users u
	JOIN follows f ON u.id = f.id_followed
	WHERE f.id_follower = ? AND f.archived = FALSE
`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Followed []*Follow

	for rows.Next() {
		f := &Follow{
			Followed: &User{},
		}
		var dateOfBirthFollowedStr string
		err := rows.Scan(
			&f.Id, &f.Archived, &f.Followed.Id, &f.Followed.Email, &f.Followed.Password, &f.Followed.FirstName, &f.Followed.LastName, &dateOfBirthFollowedStr, &f.Followed.ProfilePicture, &f.Followed.Nickname, &f.Followed.AboutMe, &f.Followed.Private, &f.Followed.CreatedAt)
		if err != nil {
			return nil, err
		}
		f.Followed.DateOfBirth, err = time.Parse("2006-01-02", dateOfBirthFollowedStr)
		if err != nil {
			return nil, err
		}
		Followed = append(Followed, f)
	}
	return Followed, nil
}

func (m *ConnDB) ArchiveFollower(followerID, followedID int) error {
	if followerID == followedID {
		return ErrSelfFollow
	}
	query := `
		UPDATE follows
		SET archived = TRUE
		WHERE id_follower = ? AND id_followed = ?
	`
	_, err := m.DB.Exec(query, followerID, followedID)
	if err != nil {
		// Affichage de l'erreur pour le débogage
		return fmt.Errorf("error executing query: %w", err)

	}
	fmt.Println("Successfully archive")
	return nil
}
