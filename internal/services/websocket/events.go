package websocket

import (
	"context"

	"github.com/AliKefall/Somnambulist/internal/database"
)

type EventHandler interface {
	// These are friends endpoint events
	NotifyFriendAccepted(
		ctx context.Context,
		accepter database.User,
		requester database.User,
	) error

	NotifyFriendRejected(
		ctx context.Context,
		rejecter database.User,
		requester database.User,
	) error
	// These are general events we need nearly all around the project
	OnUserConnected(
		ctx context.Context,
		user database.User,
	) error

	OnUserDisconnected(
		ctx context.Context,
		user database.User,
	) error



}
