package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (cfg *Config) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid_json", "Invalid JSON body", "", err)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	if isValidEmail(req.Email) || isValidUsername(req.Username) || isValidPassword(req.Password) {
		RespondWithError(w, http.StatusBadRequest, "unfit_register_parameters", "Invalid parameters", "", nil)
		return
	}

	hashed, err := cfg.Hasher.Hash(req.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "hasher_error", "Password hashing failed", "", err)
		return
	}

	userID := uuid.New()
	_, err = cfg.Queries.CreateUser(ctx, database.CreateUserParams{
		ID:        userID,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashed,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // Unique constraint violation
				switch pgErr.ConstraintName {
				case "users_email_key":
					RespondWithError(w, http.StatusConflict, "conflict_error", "Email already in use", "", nil)
					return
				case "users_usename_key":
					RespondWithError(w, http.StatusConflict, "conflict_error", "Username already in use", "", nil)
					return
				default:
					RespondWithError(w, http.StatusConflict, "conflict_error", "User already exists", "", nil)
					return

				}
			}
		}
		RespondWithError(w, http.StatusInternalServerError, "internal_server_error", "Server error at register", "", err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, struct{}{})

}
