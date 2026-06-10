package endpoints

import (
	"database/sql"

	"github.com/AliKefall/Somnambulist/internal/auth"
	"github.com/AliKefall/Somnambulist/internal/database"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	DB      *sql.DB
	Queries *database.Queries
	JWT     *auth.JWTManager
	Hasher  *auth.PasswordHasher
	Redis   *redis.Client
}
