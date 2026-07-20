package services

import (
	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/AliKefall/Somnambulist/internal/services/chat"
	"github.com/AliKefall/Somnambulist/internal/services/observability"
)

// Container groups service-layer dependencies that are shared by runtime systems
// such as the websocket hub. Keeping these dependencies here prevents the hub
// from growing every time a new service product is added.
type Container struct {
	Queries *database.Queries
	Metrics *observability.Metrics
	Chat    *chat.Service
}
