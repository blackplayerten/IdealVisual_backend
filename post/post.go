package post

import (
	"database/sql" // TODO: move to database (obertka)
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/blackplayerten/IdealVisual_backend/database"
)

func (s *Service) Get(userID uint64, postIDs []string) ([]Post, error) {
	var posts []Post

	var (
		q    = "SELECT id, photo, photo_index, date, place, text, last_updated FROM post WHERE "
		args []interface{}
	)

	if len(postIDs) == 0 {
		q += "acc = $1"
	} else {
		q += "id IN (?) AND acc = $"

		var err error
		if q, args, err = sqlx.In(q, postIDs); err != nil {
			return nil, err
		}

		q = s.db.Rebind(q)
		q += strconv.Itoa(len(args) + 1)
	}

	args = append(args, userID)

	if err := s.db.Select(&posts, q, args...); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *Service) getByID(userID uint64, postID string) (*Post, error) {
	post := new(Post)

	if err := s.db.Get(&post,
		"SELECT id, photo, photo_index, date, place, text, last_updated FROM post WHERE id = $1 AND acc = $2",
		userID,
		postID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, database.ErrNotFound
		}

		return nil, err
	}

	return post, nil
}

func (s *Service) New(post *Post) (*Post, error) {
	newP := new(Post)
	if err := s.db.Get(newP, `INSERT INTO post (acc, photo, photo_index, date, place, text, last_updated)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, photo, photo_index, date, place, text, last_updated`,
		post.Acc, post.Photo, post.PhotoIndex, post.Date, post.Place, post.Text,
		time.Unix(time.Now().Unix(), 0)); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23503" {
				return nil, database.ForeignKeyViolation{Field: s.db.PqGetKeyFromDetail(pqErr.Detail)}
			}
		}

		return nil, err
	}

	return newP, nil
}

func (s *Service) Update(post *Post) (*Post, error) {
	var (
		queryBuilder strings.Builder
		n            = 1
		args         = make([]interface{}, 0, 5)
	)

	if err := buildUpdatePostFieldsQuery(post, &queryBuilder, &n, &args); err != nil {
		return nil, err
	}

	if n == 1 {
		return s.getByID(post.Acc, post.ID)
	}

	if err := database.UpdateSetOneField(&queryBuilder, "last_updated",
		time.Unix(time.Now().Unix(), 0), &n, &args); err != nil {
		return nil, err
	}

	if _, err := queryBuilder.WriteString(
		fmt.Sprintf(` WHERE acc = $%d AND id = $%d RETURNING id, photo, photo_index, date, place, text`, n, n+1),
	); err != nil {
		return nil, err
	}

	args = append(args, post.Acc, post.ID)

	newP := new(Post)
	if err := s.db.Get(newP, queryBuilder.String(), args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, database.ErrNotFound
		}

		return nil, err
	}

	return newP, nil
}

func buildUpdatePostFieldsQuery(
	upd *Post,
	queryBuilder io.StringWriter,
	n *int,
	args *[]interface{},
) error {
	if _, err := queryBuilder.WriteString(`UPDATE post SET `); err != nil {
		return err
	}

	if upd.Photo != "" {
		if err := database.UpdateSetOneField(queryBuilder, "photo", upd.Photo, n, args); err != nil {
			return err
		}
	}

	if upd.PhotoIndex != nil {
		if err := database.UpdateSetOneField(queryBuilder, "photo_index", upd.PhotoIndex, n, args); err != nil {
			return err
		}
	}

	if upd.Date != nil {
		if err := database.UpdateSetOneField(queryBuilder, "date", upd.Date, n, args); err != nil {
			return err
		}
	}

	if upd.Place != nil {
		if err := database.UpdateSetOneField(queryBuilder, "place", upd.Place, n, args); err != nil {
			return err
		}
	}

	if upd.Text != nil {
		if err := database.UpdateSetOneField(queryBuilder, "text", upd.Text, n, args); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Delete(userID uint64, postIDs []string) error {
	q, args, err := sqlx.In(`DELETE FROM post WHERE id IN (?) AND acc = $`, postIDs)
	if err != nil {
		return err
	}

	q = s.db.Rebind(q)
	q += strconv.Itoa(len(args) + 1)
	args = append(args, userID)

	_, err = s.db.Exec(q, args...)

	return err
}
