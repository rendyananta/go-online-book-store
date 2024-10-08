package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	mathrand "math/rand"
	"time"
)

type cacheDriver interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

type Config struct {
	TokenLifetime time.Duration
	CipherKeys    []string
}

type UserSession struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	ExpiredAt time.Time `json:"expired_at"`
}

type Manager struct {
	config      Config
	cacheDriver cacheDriver
	ciphers     []cipher.Block
}

func NewAuthManager(conf Config, cacheDriver cacheDriver) (*Manager, error) {
	if len(conf.CipherKeys) == 0 {
		return nil, ErrCipherKeysIsEmpty
	}

	ciphers := make([]cipher.Block, 0, len(conf.CipherKeys))

	// loop in reverse
	for i := len(conf.CipherKeys) - 1; i >= 0; i-- {
		c, err := aes.NewCipher([]byte(conf.CipherKeys[i]))
		if err != nil {
			return nil, err
		}

		ciphers = append(ciphers, c)
	}

	if conf.TokenLifetime == 0 {
		conf.TokenLifetime = defaultTTL
	}

	return &Manager{
		config:      conf,
		cacheDriver: cacheDriver,
		ciphers:     ciphers,
	}, nil
}

func (a *Manager) Token(ctx context.Context, userID string) (string, error) {
	ttl := time.Now().Add(a.config.TokenLifetime)

	session := UserSession{
		ID:        userID,
		ExpiredAt: ttl,
	}

	contents, err := json.Marshal(session)
	if err != nil {
		return "", nil
	}

	gcm, err := cipher.NewGCM(a.ciphers[0])
	if err != nil {
		return "", nil
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	randomizer := mathrand.New(mathrand.NewSource(time.Now().UTC().Local().UnixMicro()))
	randVal := randomizer.Uint64()

	sessionKey := fmt.Sprintf("auth:%s_%s_%d", defaultUserType, userID, randVal)

	var encryptedSessionKey = gcm.Seal(nonce, nonce, []byte(sessionKey), nil)

	err = a.cacheDriver.Set(ctx, sessionKey, contents, a.config.TokenLifetime)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedSessionKey), nil
}

func (a *Manager) User(ctx context.Context, token string) (UserSession, error) {
	key, err := a.sessionKeyFor(ctx, token)
	if err != nil {
		return UserSession{}, ErrUnauthenticated
	}

	result, err := a.cacheDriver.Get(ctx, key)
	if err != nil || result == nil {
		slog.Info("err get cache", slog.String("error", err.Error()))
		return UserSession{}, ErrUnauthenticated
	}

	var authenticatedUser UserSession

	if err := json.Unmarshal(result, &authenticatedUser); err != nil {
		return UserSession{}, err
	}

	if time.Now().UnixMilli() > authenticatedUser.ExpiredAt.UnixMilli() {
		if err := a.cacheDriver.Del(ctx, key); err != nil {
			log.Printf("auth: unable to delete session for user key [%s], err: %s", key, err)
		}

		return UserSession{}, ErrTokenExpired
	}

	return authenticatedUser, nil
}

func (a *Manager) sessionKeyFor(_ context.Context, token string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(a.ciphers[0])
	if err != nil {
		return "", nil
	}

	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return "", ErrInvalidTokenSize
	}

	cipherText := []byte(decoded)

	nonce, ciphertext := cipherText[:nonceSize], cipherText[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (a *Manager) Revoke(ctx context.Context, token string) error {
	key, err := a.sessionKeyFor(ctx, token)

	// if we cannot find the key, then the token is already revoked.
	if err != nil {
		return nil
	}

	return a.cacheDriver.Del(ctx, key)
}
