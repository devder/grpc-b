package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// creates a new token for a specific username and duration
	CreateToken(username, role string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
