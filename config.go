package main

import "github.com/satori/go.uuid"

type Config struct {
	WebPort          string    `json:"web_port"`
	NotificationPort string    `json:"notification_port"`
	DBPath           string    `json:"db_path"`
	NodeId           uuid.UUID `json:"node_id"`
	LogPath          string    `json:"log_path"`
	Name             string    `json:"name"`
}
