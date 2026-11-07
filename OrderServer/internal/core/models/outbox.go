package models

import "time"

type OutboxMessage struct {
	ID          string
	OrderID     string
	MessageType string
	Payload     []byte
	Attempts    int
	CreatedAt   time.Time
}
