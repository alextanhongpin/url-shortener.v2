package health

import "time"

type Health struct {
	DeployedAt time.Time `json:"deployed_at"`
	Uptime     string    `json:"uptime"`
	Version    string    `json:"version"`
}
