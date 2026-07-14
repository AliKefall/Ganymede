package endpoints

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/google/uuid"
)

type AddFriendRequest struct {
	Username string `json:"username"`
}

type FriendResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Online   bool      `json:"online"`
}

type FriendListResponse struct {
	Friends []FriendResponse `json:"friends"`
}

type FriendRequestResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type FriendRequestListResponse struct {
	Incoming []FriendRequestResponse `json:"incoming"`
	Outgoing []FriendRequestResponse `json:"outgoing"`
}

func (cfg *Config) mustUserID(r *http.Request) (uuid.UUID, error) {
	v := r.Context().Value(UserIDKey)
	id, ok := v.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("missing user id")
	}

	return id, nil
}

func normalizeFriendPair(a, b uuid.UUID) (uuid.UUID, uuid.UUID) {
	if a.String() < b.String() {
		return a, b
	}
	return b, a
}
func mapIncomingRequests(
	users []database.ListIncomingFriendRequestsByUserIDRow,
) []FriendRequestResponse {

	result := make([]FriendRequestResponse, 0, len(users))

	for _, u := range users {
		result = append(result, FriendRequestResponse{
			ID:       u.ID,
			Username: u.Username,
		})
	}

	return result
}
func mapOutgoingRequests(
	users []database.ListOutgoingFriendRequestsByUserIDRow,
) []FriendRequestResponse {

	result := make([]FriendRequestResponse, 0, len(users))

	for _, u := range users {
		result = append(result, FriendRequestResponse{
			ID:       u.ID,
			Username: u.Username,
		})
	}

	return result
}
func (cfg *Config) HandleListFriends(w http.ResponseWriter, r *http.Request) {
	uid, err := cfg.mustUserID(r)
	if err != nil {
		RespondWithError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
			"unauthorized",
			"",
			err,
		)
		return
	}

	friends, err := cfg.Queries.ListFriendsByUserID(
		r.Context(),
		database.ListFriendsByUserIDParams{
			UserID:   uid,
			UserID_2: uid,
			FriendID: uid,
		},
	)

	if err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"database_error",
			"Could not load friends",
			"",
			err,
		)
		return
	}

	resp := make([]FriendResponse, 0, len(friends))

	for _, f := range friends {
		resp = append(resp, FriendResponse{
			ID:       f.ID,
			Username: f.Username,
			Online:   cfg.WS.IsOnline(f.ID.String()),
		})
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Username < resp[j].Username
	})

	RespondWithJSON(w, http.StatusOK, FriendListResponse{Friends: resp})
}

func (cfg *Config) HandleSendFriendRequest(w http.ResponseWriter, r *http.Request) {
	uid, err := cfg.mustUserID(r)
	if err != nil {
		RespondWithError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
			"unauthorized",
			"",
			err,
		)
		return
	}

	var req AddFriendRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(
			w,
			http.StatusBadRequest,
			"invalid_body",
			"Invalid request body",
			"",
			err,
		)
		return

	}

	req.Username = strings.TrimSpace(req.Username)

	if req.Username == "" {
		RespondWithError(
			w,
			http.StatusBadRequest,
			"empty_username",
			"Username is required",
			"",
			nil,
		)
		return
	}

	target, err := cfg.Queries.GetUserByUsername(
		r.Context(),
		req.Username,
	)

	if target.ID == uid {
		RespondWithError(
			w,
			http.StatusBadRequest,
			"self_request",
			"You cannot add yourself",
			"",
			nil,
		)
		return
	}

	a, b := normalizeFriendPair(uid, target.ID)

	friendshipExists, err := cfg.Queries.FriendshipExists(
		r.Context(),
		database.FriendshipExistsParams{
			UserID:   a,
			FriendID: b,
		},
	)

	if err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"database_error",
			"Database error",
			"",
			err,
		)
		return
	}

	if friendshipExists {
		RespondWithError(
			w,
			http.StatusConflict,
			"already_friends",
			"You are already friends",
			"",
			nil,
		)
		return
	}

	requestExists, err := cfg.Queries.FriendRequestExists(
		r.Context(),
		database.FriendRequestExistsParams{
			RequesterID: uid,
			TargetID:    target.ID,
		},
	)

	if requestExists {
		RespondWithError(
			w,
			http.StatusConflict,
			"request_exists",
			"friend request already exists",
			"",
			nil,
		)
		return
	}

	rows, err := cfg.Queries.CreateFriendRequest(
		r.Context(),
		database.CreateFriendRequestParams{
			RequesterID: uid,
			TargetID:    target.ID,
			CreatedAt:   time.Now().UTC(),
		},
	)

	if err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"database_error",
			"Could not create friend request",
			"",
			err,
		)
		return
	}

	if rows == 0 {
		RespondWithError(
			w,
			http.StatusConflict,
			"request_not_created",
			"Friend request was not created",
			"",
			nil,
		)
		return
	}

	RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "friend request created",
	})
}

func (cfg *Config) HandleAcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	uid, err := cfg.mustUserID(r)
	if err != nil {
		RespondWithError(
			w, http.StatusUnauthorized,
			"unauthorized",
			"unauthorized",
			"",
			err,
		)
		return
	}

	var req AddFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(
			w,
			http.StatusBadRequest,
			"invalid_body",
			"invalid request body",
			"",
			err,
		)
		return
	}

	req.Username = strings.TrimSpace(req.Username)

	user, err := cfg.Queries.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(
				w,
				http.StatusNotFound,
				"user_not_found",
				"User not found",
				"",
				nil,
			)
			return
		}
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"database_error",
			"Database error",
			"",
			err,
		)
		return
	}

	tx, err := cfg.DB.BeginTx(r.Context(), nil)
	if err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"transaction_error",
			"Could not start transaction",
			"",
			err,
		)

		return
	}

	defer func() {
		_ = tx.Rollback()
	}()

	qtx := cfg.Queries.WithTx(tx)
	if err := qtx.DeleteFriendRequest(
		r.Context(),
		database.DeleteFriendRequestParams{
			RequesterID: user.ID,
			TargetID:    uid,
		},
	); err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"database_error",
			"Could not delete friend request",
			"",
			err,
		)
		return
	}

	a, b := normalizeFriendPair(uid, user.ID)

	rows, err := qtx.CreateFriendship(
		r.Context(),
		database.CreateFriendshipParams{
			UserID:   a,
			FriendID: b,
			CreatedAt: time.Now().UTC(),
		},
	)
	if err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"database_error",
			"Could not create friendship",
			"",
			err,
		)

		return
	}

	if rows == 0 {
		RespondWithError(
			w,
			http.StatusConflict,
			"friendship_not_created",
			"Friendship already exists",
			"",
			nil,
		)
		return
	}

	if err := tx.Commit(); err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"transaction_error",
			"Transaction commit failed",
			"",
			err,
		)
		return
	}
	currentUser, err := cfg.Queries.GetUserByID(r.Context(), uid)
	if err != nil {
		RespondWithError(w,
			http.StatusInternalServerError,
			"database_error",
			"Could not load current user",
			"",
			err,
		)
		return
	}
	RespondWithJSON(
		w,
		http.StatusOK,
		map[string]string{
			"message": "friend request accepted",
		},
	)

	go func(accepter, requester database.User) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = cfg.WS.Events.NotifyFriendAccepted(ctx, accepter, requester)
	}(currentUser, user)

}

func (cfg *Config) HandleRejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	uid, err := cfg.mustUserID(r)
	if err != nil {
		RespondWithError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
			"unauthorized",
			"",
			err,
		)
		return
	}

	var req AddFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(
			w,
			http.StatusBadRequest,
			"invalid_body",
			"Invalid request body",
			"",
			err,
		)
		return
	}

	user, err := cfg.Queries.GetUserByUsername(r.Context(), strings.TrimSpace(req.Username))
	if err != nil {
		RespondWithError(
			w,
			http.StatusNotFound,
			"user_not_found",
			"User not found",
			"",
			err,
		)
		return
	}

	if err := cfg.Queries.DeleteFriendRequest(
		r.Context(),
		database.DeleteFriendRequestParams{
			RequesterID: user.ID,
			TargetID:    uid,
		},
	); err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"database_error",
			"Could not delete friend request",
			"",
			err,
		)
		return
	}

	rejecter, err := cfg.Queries.GetUserByID(r.Context(), uid)
	if err != nil {
		RespondWithError(
			w, http.StatusInternalServerError,
			"database_error",
			"Could not load current user",
			"",
			err,
		)
		return
	}

	RespondWithJSON(
		w,
		http.StatusOK,
		map[string]string{
			"message": "friend request rejected",
		},
	)

	go func(requester, rejecter database.User) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = cfg.WS.Events.NotifyFriendRejected(ctx, rejecter, requester)
	}(user, rejecter)

}

func (cfg *Config) HandleFriendRequests(w http.ResponseWriter, r *http.Request) {
	uid, err := cfg.mustUserID(r)
	if err != nil {
		RespondWithError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
			"unauthorized",
			"",
			err,
		)
		return
	}

	incoming, err := cfg.Queries.ListIncomingFriendRequestsByUserID(
		r.Context(),
		uid,
	)
	if err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"database_error",
			"Could not load incoming requests",
			"",
			err,
		)
		return
	}

	outgoing, err := cfg.Queries.ListOutgoingFriendRequestsByUserID(
		r.Context(),
		uid,
	)

	if err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"database_error",
			"Could not load outgoing requests",
			"",
			err,
		)
		return
	}

	RespondWithJSON(w, http.StatusOK, FriendRequestListResponse{
		Incoming: mapIncomingRequests(incoming),
		Outgoing: mapOutgoingRequests(outgoing),
	})

}

func (cfg *Config) HandleDeleteFriend(w http.ResponseWriter, r *http.Request) {
	uid, err := cfg.mustUserID(r)
	if err != nil {
		RespondWithError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
			"unauthorized",
			"",
			err,
		)
		return
	}

	username := strings.TrimSpace(r.PathValue("username"))
	if username == "" {
		RespondWithError(
			w,
			http.StatusBadRequest,
			"invalid_username",
			"username is required",
			"",
			nil,
		)
		return
	}

	friend, err := cfg.Queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		RespondWithError(
			w,
			http.StatusNotFound,
			"user_not_found",
			"user not found",
			"",
			nil,
		)
		return
	}

	a, b := normalizeFriendPair(uid, friend.ID)

	if err := cfg.Queries.DeleteFriendship(
		r.Context(),
		database.DeleteFriendshipParams{
			UserID:   a,
			FriendID: b,
		},
	); err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"database_error",
			"Could not delete friendship",
			"",
			err,
		)
		return
	}

	RespondWithJSON(
		w,
		http.StatusOK,
		map[string]string{
			"message": "friend deleted",
		},
	)
}
