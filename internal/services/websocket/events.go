package websocket

import (
	"context"

	"github.com/AliKefall/Somnambulist/internal/database"
)

type EventHandler interface {
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

	OnUserConnected(
		ctx context.Context,
		user database.User,
	)	error

	OnUserDisconnected(
		ctx context.Context,
		user database.User,
	) error
}

