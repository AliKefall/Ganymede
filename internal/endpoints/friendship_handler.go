package endpoints

import (
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
	Username string `json:"username"`
	Online   bool   `json:"online"`
}

type FriendListResponse struct {
	Friends []FriendResponse `json:"friends"`
}

type FriendRequestListResponse struct {
	Incoming []string `json:"incoming"`
	Outgoing []string `json:"outgoing"`
}

func (cfg *Config) mustUserID(r *http.Request) (uuid.UUID, error) {
	v := r.Context().Value(UserIDKey)
	id, ok := v.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("missing user id")
	}
	return id, nil
}

func normalizeFriendPair(userID, friendID string) (string, string) {
	if userID < friendID {
		return userID, friendID
	}
	return friendID, userID
}

func (cfg *Config) HandleListFriends(w http.ResponseWriter, r *http.Request) {
	uid, err := cfg.mustUserID(r)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "friends_error", "Unauthorized", "", err)
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
		RespondWithError(w, http.StatusInternalServerError, "database_error", "Database error at friends endpoint", "", err)
		return
	}

	resp := make([]FriendResponse, 0, len(friends))
	for _, friend := range friends {
		resp = append(resp, FriendResponse{
			Username: friend.Username,
			Online:   false,
		})
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Username < resp[j].Username
	})

	RespondWithJSON(w, 200, FriendListResponse{Friends: resp})

}

func (cfg *Config) HandleSendFriendRequest(w http.ResponseWriter, r *http.Request) {
	uid, err := cfg.mustUserID(r)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "friends_error", "Unauthorized", "", err)
		return
	}

	var req AddFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "json_error", "Invalid request body", "", err)
		return
	}
	username := strings.TrimSpace(req.Username)

	target, err := cfg.Queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "not_found", "User not found", "", err)
		return
	}

	if target.ID == uid {
		RespondWithError(w, http.StatusBadRequest, "self_request", "You can not add yourself as a friend", "", nil)
		return
	}

	a, b := normalizeFriendPair(uid.String(), target.ID.String())
	aid, _ := uuid.Parse(a)
	bid, _ := uuid.Parse(b)

	exists, _ := cfg.Queries.FriendshipExists(
		r.Context(),
		database.FriendshipExistsParams{
			UserID:   aid,
			FriendID: bid,
		},
	)

	if exists > 0 {
		RespondWithError(w, http.StatusBadRequest, "already_friends", "You are already friends", "", nil)
		return
	}

	existsReq, _ := cfg.Queries.FriendRequestExists(
		r.Context(),
		database.FriendRequestExistsParams{
			RequesterID: uid,
			TargetID:    target.ID,
		},
	)

	if existsReq > 0 {
		RespondWithError(w, http.StatusBadRequest, "request_exists", "You already have a request", "", nil)
		return
	}

	err = cfg.Queries.CreateFriendRequest(
		r.Context(),
		database.CreateFriendRequestParams{
			RequesterID: uid,
			TargetID:    target.ID,
			CreatedAt:   time.Now().UTC(),
		},
	)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "database_error", "Database error at friends endpoint", "", err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "created"})
}

func (cfg *Config) HandleAcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	uid, err := cfg.mustUserID(r)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "friends_error", "unauthorized", "", err)
		return
	}

	var req AddFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "json_error", "Invalid json body", "", err)
		return
	}

	username := strings.TrimSpace(req.Username)

	u, err := cfg.Queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "not_found", "User not found", "", err)
		return
	}

	tx, _ := cfg.DB.BeginTx(r.Context(), nil)
	defer tx.Rollback()

	qtx := cfg.Queries.WithTx(tx)

	_ = qtx.DeleteFriendRequest(r.Context(), database.DeleteFriendRequestParams{
		RequesterID: u.ID,
		TargetID:    uid,
	})

	a, b := normalizeFriendPair(uid.String(), u.ID.String())
	aid, _ := uuid.Parse(a)
	bid, _ := uuid.Parse(b)

	_ = qtx.CreateFriendship(r.Context(), database.CreateFriendshipParams{
		UserID:    aid,
		FriendID:  bid,
		CreatedAt: time.Now().UTC(),
	})

	tx.Commit()

	RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "accepted",
	})
}

func (cfg *Config) HandleRejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	uid, err := cfg.mustUserID(r)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "friends_error", "Unauthorized", "", err)
		return
	}

	var req AddFriendRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	u, err := cfg.Queries.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "not_found", "User not found", "", err)
		return
	}

	_ = cfg.Queries.DeleteFriendRequest(r.Context(),
		database.DeleteFriendRequestParams{
			RequesterID: u.ID,
			TargetID:    uid,
		})
	RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "rejected",
	})
}

func (cfg *Config) HandleFriendRequests(w http.ResponseWriter, r *http.Request) {
	uid, err := cfg.mustUserID(r)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "friends_error", "Unauthorized", "", err)
		return
	}

	incoming, _ := cfg.Queries.ListIncomingFriendRequestsByUserID(r.Context(), uid)
	outgoing, _ := cfg.Queries.ListOutgoingFriendRequestsByUserID(r.Context(), uid)

	//NOTE: I still haven't tasted if the user schema is actually json safe or not. So if it is change the
	// incoming and outgoing types.
	RespondWithJSON(w, http.StatusOK, FriendRequestListResponse{
		Incoming: getUsernames(incoming),
		Outgoing: getUsernames(outgoing),
	})
}
func getUsernames(users []database.User) []string {
    result := make([]string, 0, len(users))

    for _, u := range users {
        result = append(result, u.Username)
    }

    return result
}
