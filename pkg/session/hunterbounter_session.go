package hunterbounter_session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
	"log"
	"time"
)

var (
	store = sqlite3.New()
)

var sessConfig = session.Config{
	Expiration:     60 * time.Minute,        // Expire sessions after 30 minutes of inactivity
	KeyLookup:      "cookie:hunter_session", // Recommended to use the __Host- prefix when serving the app over TLS
	CookieSecure:   true,
	CookieHTTPOnly: true,
	CookieSameSite: "Lax",
	Storage:        store,
}
var HunterSession = *session.New(sessConfig)

var GlobalSessionTempList = make(map[string]interface{})

func SetSessionValue(c *fiber.Ctx, key string, value interface{}) error {
	sess, err := HunterSession.Get(c)
	if err != nil {
		return err
	}

	sess.Set(key, value)
	log.Println("Session set: ", key, value)
	log.Println("Session: ", sess)
	return sess.Save()
}

func GetSessionValue(c *fiber.Ctx, key string) (interface{}, error) {
	sess, err := HunterSession.Get(c)
	if err != nil {
		return nil, err
	}

	return sess.Get(key), nil
}

func DestroySession(c *fiber.Ctx) error {
	sess, err := HunterSession.Get(c)
	if err != nil {
		return err
	}

	return sess.Destroy()
}

func init() {
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
	log.Println("hunterbounter package initialized.")
}
