package database

import (
	"time"

	"github.com/dhax/go-base/auth"
	"github.com/go-pg/pg"
)

// AuthStore implements database operations for account authentication.
type AuthStore struct {
	db *pg.DB
}

// NewAuthStore return an AuthStore.
func NewAuthStore(db *pg.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

// GetByID returns an account by ID.
func (s *AuthStore) GetByID(id int) (*auth.Account, error) {
	a := auth.Account{ID: id}
	err := s.db.Model(&a).
		Column("account.*").
		Where("id = ?id").
		First()
	return &a, err
}

// GetByEmail returns an account by email.
func (s *AuthStore) GetByEmail(e string) (*auth.Account, error) {
	a := auth.Account{Email: e}
	err := s.db.Model(&a).
		Column("id", "active", "email", "name").
		Where("email = ?email").
		First()
	return &a, err
}

// GetByRefreshToken returns an account and refresh token by token identifier.
func (s *AuthStore) GetByRefreshToken(t string) (*auth.Account, *auth.Token, error) {
	token := auth.Token{Token: t}
	err := s.db.Model(&token).
		Where("token = ?token").
		First()
	if err != nil {
		return nil, nil, err
	}

	a := auth.Account{ID: token.AccountID}
	err = s.db.Model(&a).
		Column("account.*").
		Where("id = ?id").
		First()

	return &a, &token, err
}

// UpdateAccount upates account data related to authentication.
func (s *AuthStore) UpdateAccount(a *auth.Account) error {
	_, err := s.db.Model(a).
		Column("last_login").
		Update()
	return err
}

// SaveRefreshToken creates or updates a refresh token.
func (s *AuthStore) SaveRefreshToken(t *auth.Token) error {
	var err error
	if t.ID == 0 {
		err = s.db.Insert(t)
	} else {
		err = s.db.Update(t)
	}
	return err
}

// DeleteRefreshToken deletes a refresh token.
func (s *AuthStore) DeleteRefreshToken(t *auth.Token) error {
	err := s.db.Delete(t)
	return err
}

// PurgeExpiredToken deletes expired refresh token.
func (s *AuthStore) PurgeExpiredToken() error {
	_, err := s.db.Model(&auth.Token{}).
		Where("expiry < ?", time.Now()).
		Delete()

	return err
}
