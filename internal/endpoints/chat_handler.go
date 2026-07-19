package endpoints

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg *Config) HandleListConversations(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		RespondWithError(w, http.StatusUnauthorized, "unauthorized", "invalid user context", "", nil)
		return
	}

	conversations, err := cfg.Chat.ListConversations(
		r.Context(),
		userID,
	)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "database_error", "could not load conversations", "", err)
		return
	}

	response := make([]ConversationResponse, 0, len(conversations))

	for _, conversation := range conversations {
		response = append(
			response,
			toConversationResponse(conversation),
		)
	}

	RespondWithJSON(w, http.StatusOK, response)
}

func (cfg *Config) HandleConversationMessages(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		RespondWithError(
			w,
			http.StatusUnauthorized,
			"unauthorized",
			"invalid user context",
			"",
			nil,
		)
		return
	}

	conversationID, err := uuid.Parse(
		chi.URLParam(r, "conversationID"),
	)

	if err != nil {
		RespondWithError(
			w,
			http.StatusBadRequest,
			"invalid_conversation",
			"conversation id is invalid",
			"",
			err,
		)
		return
	}

	isMember, err := cfg.Chat.IsConversationMember(
		r.Context(),
		conversationID,
		userID,
	)

	if err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"database_error",
			"could not verify conversation",
			"",
			err,
		)
		return
	}

	if !isMember {
		RespondWithError(
			w,
			http.StatusForbidden,
			"forbidden",
			"you are not a member of this conversation",
			"",
			nil,
		)
		return
	}

	messages, err := cfg.Chat.GetMessages(
		r.Context(),
		conversationID,
		50,
		0,
	)

	if err != nil {
		RespondWithError(
			w,
			http.StatusInternalServerError,
			"database_error",
			"could not load messages",
			"",
			err,
		)
		return
	}

	response := make([]ChatMessageResponse, 0, len(messages))

	for _, message := range messages {
		response = append(response, toChatMessageResponse(message))
	}
	RespondWithJSON(w, http.StatusOK, response)

}
