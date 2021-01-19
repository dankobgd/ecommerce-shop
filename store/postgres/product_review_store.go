package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// PgReviewStore is the postgres implementation
type PgReviewStore struct {
	PgStore
}

// NewPgReviewStore creates the new review store
func NewPgReviewStore(pgst *PgStore) store.ProductReviewStore {
	return &PgReviewStore{*pgst}
}

var (
	msgUniqueConstraintReview = &i18n.Message{ID: "store.postgres.review.save.unique_constraint.app_error", Other: "review already exists"}
	msgSaveReview             = &i18n.Message{ID: "store.postgres.review.save.app_error", Other: "could not save review"}
	msgUpdateReview           = &i18n.Message{ID: "store.postgres.review.update.app_error", Other: "could not update review"}
	msgBulkInsertReviews      = &i18n.Message{ID: "store.postgres.review.bulk.insert.app_error", Other: "could not bulk insert reviews"}
	msgGetReview              = &i18n.Message{ID: "store.postgres.review.get.app_error", Other: "could not get the review"}
	msgGetReviews             = &i18n.Message{ID: "store.postgres.review.get.app_error", Other: "could not get the reviews"}
	msgDeleteReview           = &i18n.Message{ID: "store.postgres.review.delete.app_error", Other: "could not delete review"}
)

// BulkInsert inserts multiple reviews in the db
func (s PgReviewStore) BulkInsert(reviews []*model.ProductReview) *model.AppErr {
	q := `INSERT INTO product_review(user_id, product_id, rating, title, comment, created_at, updated_at) VALUES(:user_id, :product_id, :rating, :title, :comment, :created_at, :updated_at) RETURNING id`

	if _, err := s.db.NamedExec(q, reviews); err != nil {
		return model.NewAppErr("PgReviewStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertReviews, http.StatusInternalServerError, nil)
	}
	return nil
}

// Save inserts the new review in the db
func (s PgReviewStore) Save(pid int64, review *model.ProductReview) (*model.ProductReview, *model.AppErr) {
	q := `INSERT INTO product_review(user_id, product_id, rating, title, comment, created_at, updated_at) VALUES(:user_id, :product_id, :rating, :title, :comment, :created_at, :updated_at) RETURNING id`

	var id int64
	rows, err := s.db.NamedQuery(q, review)
	if err != nil {
		return nil, model.NewAppErr("PgReviewStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveReview, http.StatusInternalServerError, nil)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
	}
	if err := rows.Err(); err != nil {
		if IsUniqueConstraintViolationError(err) {
			return nil, model.NewAppErr("PgReviewStore.Save", model.ErrConflict, locale.GetUserLocalizer("en"), msgUniqueConstraintReview, http.StatusInternalServerError, nil)
		}
		return nil, model.NewAppErr("PgReviewStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveReview, http.StatusInternalServerError, nil)
	}

	review.ID = id
	return review, nil
}

// Get gets one review by id
func (s PgReviewStore) Get(pid, rid int64) (*model.ProductReview, *model.AppErr) {
	q := `SELECT 
	r.*,
	u.id AS user_id,
	u.first_name AS user_first_name,
	u.last_name AS user_last_name,
	u.username AS user_username,
	u.avatar_url AS user_avatar_url,
	u.avatar_public_id AS user_avatar_public_id
	FROM product_review r 
	LEFT JOIN public.user u ON r.user_id = u.id 
	WHERE r.product_id = $1 AND r.id = $2`
	var rj reviewJoin
	if err := s.db.Get(&rj, q, pid, rid); err != nil {
		return nil, model.NewAppErr("PgReviewStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetReview, http.StatusInternalServerError, nil)
	}
	return rj.ToReview(), nil
}

// GetAll returns all reviews
func (s PgReviewStore) GetAll(pid int64) ([]*model.ProductReview, *model.AppErr) {
	q := `SELECT 
	r.*,
	u.id AS user_id,
	u.first_name AS user_first_name,
	u.last_name AS user_last_name,
	u.username AS user_username,
	u.avatar_url AS user_avatar_url,
	u.avatar_public_id AS user_avatar_public_id
	FROM product_review r 
	LEFT JOIN public.user u ON r.user_id = u.id
	WHERE r.product_id = $1`

	var rj []reviewJoin
	if err := s.db.Select(&rj, q, pid); err != nil {
		return nil, model.NewAppErr("PgReviewStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetReviews, http.StatusInternalServerError, nil)
	}

	var reviews = make([]*model.ProductReview, 0)
	for _, x := range rj {
		reviews = append(reviews, x.ToReview())
	}
	return reviews, nil
}

// Update updates the review
func (s PgReviewStore) Update(pid, rid int64, rev *model.ProductReview) (*model.ProductReview, *model.AppErr) {
	q := `UPDATE product_review SET rating=:rating, title=:title, comment=:comment, updated_at=:updated_at WHERE product_id=:product_id AND id=:review_id`
	m := map[string]interface{}{"product_id": pid, "review_id": rid, "rating": rev.Rating, "title": rev.Title, "comment": rev.Comment, "updated_at": rev.UpdatedAt}
	if _, err := s.db.NamedExec(q, m); err != nil {
		return nil, model.NewAppErr("PgReviewStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateReview, http.StatusInternalServerError, nil)
	}
	return rev, nil
}

// Delete hard deletes the review
func (s PgReviewStore) Delete(pid, rid int64) *model.AppErr {
	if _, err := s.db.NamedExec("DELETE from product_review WHERE product_id=:product_id AND id=:review_id", map[string]interface{}{"product_id": pid, "review_id": rid}); err != nil {
		return model.NewAppErr("PgReviewStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteReview, http.StatusInternalServerError, nil)
	}
	return nil
}
