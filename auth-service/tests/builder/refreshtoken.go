package builder

import (
	"database/sql"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/token"
	"testing"
	"time"
)

type RefreshTokenBuilder struct {
	t     *testing.T
	model *token.RefreshToken
}

func NewRefreshTokenBuilder(t *testing.T) *RefreshTokenBuilder {
	return &RefreshTokenBuilder{
		model: &token.RefreshToken{},
	}
}

func (b *RefreshTokenBuilder) WithUserID(userID string) *RefreshTokenBuilder {
	b.model.UserID = userID
	return b
}

func (b *RefreshTokenBuilder) WithHashedToken(hashedToken string) *RefreshTokenBuilder {
	b.model.HashedToken = hashedToken
	return b
}

func (b *RefreshTokenBuilder) WithExpiresAt(t time.Time) *RefreshTokenBuilder {
	b.model.ExpiresAt = t
	return b
}

func (b *RefreshTokenBuilder) Build(db *sql.DB) *token.RefreshToken {
	if b.model.UserID == "" {
		b.t.Fatal("builder: missing user_id")
	}
	if b.model.HashedToken == "" {
		b.t.Fatal("builder: missing hashed_token")
	}
	if b.model.ExpiresAt.IsZero() {
		b.t.Fatal("builder: missing expires_at")
	}

	q := "insert into refresh_tokens (user_id, hashed_token, expires_at) values (?, ?, ?) returning *;"
	row := db.QueryRow(q, b.model.UserID, b.model.HashedToken, b.model.ExpiresAt)
	err := row.Scan(&b.model.ID, &b.model.UserID, &b.model.HashedToken, &b.model.ExpiresAt, &b.model.Revoked, &b.model.CreatedAt)
	if err != nil {
		b.t.Fatalf("builder: failed to insert refresh_token: %v", err)
	}

	return b.model
}
