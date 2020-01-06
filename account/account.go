package account

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) CheckCredentials(cr *Credentials) (*Account, bool, error) {
	var accWithPassword FullAccount
	if err := s.db.Get(&accWithPassword,
		`SELECT id, email, password, username, avatar FROM account WHERE email = $1`,
		cr.Email,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}

		return nil, false, err
	}

	matched, err := comparePasswords(accWithPassword.Password, cr.Password)
	if err != nil {
		return nil, false, err
	}

	if matched {
		return &accWithPassword.Account, true, nil
	}

	return nil, false, nil
}

func (s *Service) GetByID(id uint64) (*Account, error) {
	acc := new(Account)
	if err := s.db.Get(acc, `SELECT id, email, username, avatar FROM account WHERE id = $1`, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return acc, nil
}

func (s *Service) New(info *FullAccount) (*Account, error) {
	acc := new(Account)

	hashedPassword, err := hashAndSalt(info.Password)
	if err != nil {
		return nil, err
	}

	if err = s.db.Get(
		acc,
		`INSERT INTO account (email, password, username) VALUES ($1, $2, $3) RETURNING id, email, username`,
		info.Email,
		hashedPassword,
		info.Username,
	); err != nil {
		pqErr := err.(*pq.Error)
		switch pqErr.Code {
		case "23505":
			return nil, UniqueConstraintViolationError{Field: s.db.PqGetKeyFromDetail(pqErr.Detail)}
		case "23502":
			return nil, NotNullConstraintViolationError{Field: s.db.PqGetKeyFromDetail(pqErr.Detail)}
		}

		return nil, err
	}

	return acc, nil
}

func (s *Service) Update(upd *FullAccount) (*Account, error) {
	var (
		queryBuilder strings.Builder
		n            = 1
		args         = make([]interface{}, 0, 5)
	)

	if err := buildUpdateAccountFieldsQuery(upd, &queryBuilder, &n, &args); err != nil {
		return nil, err
	}

	if n == 1 {
		return s.GetByID(upd.ID)
	}

	if _, err := queryBuilder.WriteString(
		fmt.Sprintf(` WHERE id = $%d RETURNING id, email, username, avatar`, n),
	); err != nil {
		return nil, err
	}

	args = append(args, upd.ID)

	acc := new(Account)
	if err := s.db.Get(acc, queryBuilder.String(), args...); err != nil {
		pqErr := err.(*pq.Error)
		switch pqErr.Code {
		case "23505":
			return nil, UniqueConstraintViolationError{Field: s.db.PqGetKeyFromDetail(pqErr.Detail)}
		case "23502":
			return nil, NotNullConstraintViolationError{Field: s.db.PqGetKeyFromDetail(pqErr.Detail)}
		}

		return nil, err
	}

	return acc, nil
}

func buildUpdateAccountFieldsQuery(
	upd *FullAccount,
	queryBuilder io.StringWriter,
	n *int,
	args *[]interface{},
) error {
	if _, err := queryBuilder.WriteString(`UPDATE account SET `); err != nil {
		return err
	}

	if upd.Email != "" {
		if err := updateSetOneField(queryBuilder, "email", upd.Email, n, args); err != nil {
			return err
		}
	}

	if upd.Username != "" {
		if err := updateSetOneField(queryBuilder, "username", upd.Username, n, args); err != nil {
			return err
		}
	}

	if upd.Password != "" {
		var err error
		if upd.Password, err = hashAndSalt(upd.Password); err != nil {
			return err
		}

		if err := updateSetOneField(queryBuilder, "password", upd.Password, n, args); err != nil {
			return err
		}
	}

	if upd.Avatar != nil && *upd.Avatar != "" {
		if err := updateSetOneField(queryBuilder, "avatar", upd.Avatar, n, args); err != nil {
			return err
		}
	}

	return nil
}

func hashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash), err
}

func comparePasswords(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	switch err {
	case nil:
		return true, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return false, nil
	default:
		return false, err
	}
}
