package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	algorithm = "argon2id"
	version   = 19
)

type PasswordHasher struct {
	memory      uint32
	time        uint32
	parallelism uint8
	keyLen      uint32
	saltLen     uint32
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		memory:      64 * 1024, // Every iteration takes 64mb
		time:        3,
		parallelism: 2,
		keyLen:      32,
		saltLen:     16,
	}
}

func (p *PasswordHasher) Hash(password string) (string, error) {
	salt := make([]byte, p.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		p.time,
		p.memory,
		p.parallelism,
		p.keyLen,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"%s$v=%d$m=%d,p=%d$%s$%s",
		algorithm,
		version,
		p.memory,
		p.parallelism,
		b64Salt,
		b64Hash,
	)
	return encoded, nil
}

func (p *PasswordHasher) Verify(password, encoded string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 5 {
		return false, errors.New("Invalid hash format")
	}

	if parts[0] != algorithm {
		return false, errors.New("Unsupported algorithm")
	}

	var parsedVersion int
	_, err := fmt.Sscan(parts[1], "v=%d", &parsedVersion)
	if err != nil || parsedVersion != version {
		return false, errors.New("Incompatible version")
	}

	var memory uint32
	var time uint32
	var parallelism uint8

	_, err = fmt.Sscanf(parts[2], "m=%d,t=%d,p=%d", &memory, &time, &parallelism)
	if err != nil {
		return false, errors.New("Invalid password parameters")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[3])

	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[4])

	if err != nil {
		return false, err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		parallelism,
		uint32(len(expectedHash)), // This needs to be dynamic value
	)

	if subtle.ConstantTimeCompare(hash, expectedHash) == 1 {
		return true, nil
	}
	return false, nil
}
