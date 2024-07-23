package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Session struct {
	SessionID string
	Iduser    int
	Data      map[string]interface{}
	Expires   time.Time
	IsActive  bool
	Cookie    http.Cookie
}

func generateSessionID() string {
	u2, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}

	return u2.String()
}

type SessionManager struct {
	Sessions map[string]*Session
	Db       *sql.DB
	mu       sync.RWMutex
}

func NewSessionManager(db *sql.DB) *SessionManager {
	return &SessionManager{
		Sessions: make(map[string]*Session),
		Db:       db,
	}
}
func (sess *SessionManager) ClearSessions() {
	sess.Sessions = make(map[string]*Session)
}

func (sess *SessionManager) NewSession(w http.ResponseWriter, iduser int) (*Session, error) {
	sess.mu.RLock()         // Lock the map for reading
	defer sess.mu.RUnlock() // Defer the unlocking of the map
	sessionId := generateSessionID()

	session := &Session{
		SessionID: sessionId,
		Iduser:    iduser,
		IsActive:  true,
		Data:      make(map[string]interface{}),
		Expires:   time.Now().Add(120 * time.Minute),
		Cookie: http.Cookie{
			Name:     "session",
			Value:    sessionId,
			Expires:  time.Now().Add(120 * time.Minute),
			MaxAge:   3600,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		},
	}

	encodedData, err := encodeSessionData(session.Data)
	if err != nil {
		return nil, err
	}
	_, err = sess.Db.Exec("INSERT INTO session (sessionID, iduser, data, expires, isactive) VALUES (?, ?, ?, ?,true)", sessionId, iduser, encodedData, session.Expires)
	if err != nil {
		return nil, err
	}

	sess.Sessions[sessionId] = session

	http.SetCookie(w, &session.Cookie)

	return session, nil
}

func (sess *SessionManager) GetSession(r *http.Request) (*Session, error) {
	sess.mu.Lock()         // Lock the map for writing
	defer sess.mu.Unlock() // Defer the unlocking of the map
	cookie, err := r.Cookie("session")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, nil
		}
		return nil, err
	}

	sessionID := cookie.Value
	// session, ok := sess.Sessions[sessionID]
	// if !ok {
	row := sess.Db.QueryRow("SELECT iduser, data, expires FROM session WHERE sessionID = ?", sessionID)
	dataBytes := []byte{}
	expires := time.Time{}
	var iduser int
	err = row.Scan(&iduser, &dataBytes, &expires)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	decodedData, err := decodeSessionData(dataBytes)
	if err != nil {
		return nil, err
	}

	if time.Now().After(expires) {
		sess.DeleteSession(sessionID)
		return nil, nil
	}

	session := &Session{
		SessionID: sessionID,
		Iduser:    iduser,
		Data:      decodedData,
		Expires:   expires,
		Cookie:    *cookie,
	}

	sess.Sessions[sessionID] = session
	// }

	return session, nil
}

func decodeSessionData(dataBytes []byte) (map[string]interface{}, error) {
	var decodedData map[string]interface{}
	err := json.Unmarshal(dataBytes, &decodedData)
	return decodedData, err
}

func encodeSessionData(data map[string]interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (m *SessionManager) DeleteSession(sessionID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	stmt := `DELETE FROM session WHERE sessionID = ?`
	_, err := m.Db.Exec(stmt, sessionID)
	return err
}

func (sess *SessionManager) GetActiveSession(userID int) (*Session, error) {
	sess.mu.Lock()         // Lock the map for writing
	defer sess.mu.Unlock() // Defer the unlocking of the map
	stmt := `SELECT sessionID FROM session WHERE iduser = ? AND isactive = true`
	row := sess.Db.QueryRow(stmt, userID)

	s := &Session{}
	err := row.Scan(&s.SessionID)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("models: no matching record found")
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (u *SessionManager) GetSessions(IdUser int) (bool, error) {
	stmt := `SELECT sessionID FROM session WHERE iduser = ?`
	s := &Session{}
	err := u.Db.QueryRow(stmt, IdUser).Scan(&s.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false,errors.New("models: no matching record found")
		} else {
			return false, err
		}
	}

	return true, nil
}

func DeleteAllTableContent(db *sql.DB, tableName string) error {
	query := "DELETE FROM " + tableName

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error deleting all content from table %s: %v", tableName, err)
		return err
	}

	return nil
}
