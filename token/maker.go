package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)

	verifyToken(token string) (*Payload, error)
}
