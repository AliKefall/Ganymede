package handlers

import (
	"context"
	"log"

	"github.com/AliKefall/Somnambulist/internal/database"
)

const (
	TypeFriendRequestReceived = "friend_request_received"
	TypeFriendRequestAccepted = "friend_request_accepted"
	TypeFriendRequestRejected = "friend_request_rejected"
	TypeFriendOnline          = "friend_online"
	TypeFriendOffline         = "friend_offline"
)

type FriendUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type FriendRequestPayload struct {
	User FriendUser `json:"user"`
}

type FriendPresencePayload struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Online   bool   `json:"online"`
}


func (cfg *Config) NotifyFriendAccepted(ctx context.Context, accepter database.User, requester database.User) error {
	return cfg.Hub.SendToUser(
		requester.ID.String(),
		TypeFriendRequestAccepted,
		FriendRequestPayload{
			User: FriendUser{
				ID:       accepter.ID.String(),
				Username: accepter.Username,
			},
		},
	)
}

func (cfg *Config) NotifyFriendRejected(ctx context.Context, rejecter database.User, requester database.User) error {
	return cfg.Hub.SendToUser(
		requester.ID.String(),
		TypeFriendRequestRejected,
		FriendRequestPayload{
			User: FriendUser{
				ID:       rejecter.ID.String(),
				Username: rejecter.Username,
			},
		},
	)
}

func (cfg *Config) OnUserConnected(ctx context.Context, user database.User) error {
	friends, err := cfg.Queries.ListFriendsByUserID(ctx, database.ListFriendsByUserIDParams{
		UserID:   user.ID,
		UserID_2: user.ID,
		FriendID: user.ID,
	})
	if err != nil {
		return err
	}

	payload := FriendPresencePayload{
		UserID:   user.ID.String(),
		Username: user.Username,
		Online:   true,
	}

	for _, friend := range friends {
		err = cfg.Hub.SendToUser(
			friend.ID.String(),
			TypeFriendOnline,
			payload,
		)
		if err := cfg.Hub.SendToUser(
			friend.ID.String(),
			TypeFriendOffline,
			payload,
		); err != nil {
			log.Printf(
				"Could not send offline event to %s: %v",
				friend.Username,
				err,
			)
		}
	}

	return nil
}

func (cfg *Config) OnUserDisconnected(ctx context.Context, user database.User) error {
	friends, err := cfg.Queries.ListFriendsByUserID(ctx, database.ListFriendsByUserIDParams{
		UserID:   user.ID,
		UserID_2: user.ID,
		FriendID: user.ID,
	})
	if err != nil {
		return err
	}

	payload := FriendPresencePayload{
		UserID:   user.ID.String(),
		Username: user.Username,
		Online:   false,
	}

	for _, friend := range friends {
		_ = cfg.Hub.SendToUser(
			friend.ID.String(),
			TypeFriendOffline,
			payload,
		)
	}
	return nil
}
