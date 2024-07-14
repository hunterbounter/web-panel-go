package hunterbounter_session

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
	"hunterbounter.com/web-panel/pkg/date_util"
	"log"
	"time"
)

var (
	store = sqlite3.New()
)

type SessionData struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	IpAddress    string `json:"ip_address"`
	LoginTime    string `json:"login_time"`
	LastActivity string `json:"last_activity"`
}

var sessConfig = session.Config{
	Expiration:     60 * time.Minute,        // Expire sessions after 30 minutes of inactivity
	KeyLookup:      "cookie:hunter_session", // Recommended to use the __Host- prefix when serving the app over TLS
	CookieSecure:   true,
	CookieHTTPOnly: true,
	CookieSameSite: "Lax",
	Storage:        store,
}
var HunterSession = *session.New(sessConfig)

var HunterBounterSession *session.Store

const SESSION_KEY = "hunterbounter_session"

// SetSessionValue adds a value to the session
func SetSessionValue(c *fiber.Ctx, key string, value *SessionData) error {
	sess, err := HunterBounterSession.Get(c)
	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Println("Error marshaling value:", err)
		return err
	}

	sess.Set(key, jsonValue)
	log.Println("Session set:", key, string(jsonValue))
	err = sess.Save()
	if err != nil {
		log.Println("Error saving session:", err)
	}
	log.Println("Session after save:", sess)
	return err
}

// GetSessionValue retrieves a value from the session
func GetSessionValue(c *fiber.Ctx, key string) (*SessionData, error) {

	sess, err := HunterBounterSession.Get(c)
	if err != nil {
		log.Println("Error getting session:", err)
		return nil, err
	}

	log.Println("Session get:", key)

	value := sess.Get(key)
	if value == nil {
		log.Println("No value found for key:", key)

		return nil, errors.New("Session value not found")
	}

	var data SessionData
	err = json.Unmarshal(value.([]byte), &data)
	if err != nil {
		log.Println("Error unmarshaling value:", err)
		return nil, err
	}
	// add activity
	data.LastActivity = date_util.DateYYYYMMDDHH24MISSWithTRTimezone()
	err = SetSessionValue(c, key, &data)
	return &data, nil
}

func DestroySession(c *fiber.Ctx) error {
	sess, err := HunterSession.Get(c)
	if err != nil {
		return err
	}

	return sess.Destroy()
}

func InitSession() {
	// Initialize custom config
	store = sqlite3.New(sqlite3.Config{
		Database:        "/fiber.sqlite3",
		Table:           "hunter_sessions",
		Reset:           false,
		GCInterval:      10 * time.Second,
		MaxOpenConns:    100,
		MaxIdleConns:    100,
		ConnMaxLifetime: 1 * time.Second,
	})

	HunterBounterSession = session.New(sessConfig)

	log.Println("hunterbounter package initialized.")
}
