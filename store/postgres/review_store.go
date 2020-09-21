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
func NewPgReviewStore(pgst *PgStore) store.ReviewStore {
	return &PgReviewStore{*pgst}
}

var (
	msgUniqueConstraintReview = &i18n.Message{ID: "store.postgres.review.save.unique_constraint.app_error", Other: "invalid review foreign key"}
	msgSaveReview             = &i18n.Message{ID: "store.postgres.review.save.app_error", Other: "could not save review"}
	msgUpdateReview           = &i18n.Message{ID: "store.postgres.review.update.app_error", Other: "could not update review"}
	msgBulkInsertReviews      = &i18n.Message{ID: "store.postgres.review.bulk.insert.app_error", Other: "could not bulk insert reviews"}
	msgGetReview              = &i18n.Message{ID: "store.postgres.review.get.app_error", Other: "could not get the review"}
	msgGetReviews             = &i18n.Message{ID: "store.postgres.review.get.app_error", Other: "could not get the review"}
	msgDeleteReview           = &i18n.Message{ID: "store.postgres.review.delete.app_error", Other: "could not delete review"}
)

// BulkInsert inserts multiple reviews in the db
func (s PgReviewStore) BulkInsert(reviews []*model.Review) *model.AppErr {
	q := `INSERT INTO public.product_review(user_id, product_id, rating, comment, created_at, updated_at) VALUES(:user_id, :product_id, :rating, :comment, :created_at, :updated_at) RETURNING id`

	if _, err := s.db.NamedExec(q, reviews); err != nil {
		return model.NewAppErr("PgReviewStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertReviews, http.StatusInternalServerError, nil)
	}
	return nil
}

// Save inserts the new review in the db
func (s PgReviewStore) Save(review *model.Review) (*model.Review, *model.AppErr) {
	q := `INSERT INTO public.product_review(user_id, product_id, rating, comment, created_at, updated_at) VALUES(:user_id, :product_id, :rating, :comment, :created_at, :updated_at) RETURNING id`

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

// Update updates the review
func (s PgReviewStore) Update(id int64, rev *model.Review) (*model.Review, *model.AppErr) {
	q := `UPDATE public.product_review SET product_id=:product_id, rating=:rating, comment=:comment, updated_at=:updated_at WHERE id=:id`
	if _, err := s.db.NamedExec(q, rev); err != nil {
		return nil, model.NewAppErr("PgReviewStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateReview, http.StatusInternalServerError, nil)
	}
	return rev, nil
}

// Get gets one review by id
func (s PgReviewStore) Get(id int64) (*model.Review, *model.AppErr) {
	var review model.Review
	if err := s.db.Get(&review, "SELECT * FROM public.product_review WHERE id = $1", id); err != nil {
		return nil, model.NewAppErr("PgReviewStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetReview, http.StatusInternalServerError, nil)
	}
	return &review, nil
}

// GetAll returns all reviews
func (s PgReviewStore) GetAll(limit, offset int) ([]*model.Review, *model.AppErr) {
	var reviews = make([]*model.Review, 0)
	if err := s.db.Select(&reviews, `SELECT COUNT(*) OVER() AS total_count, * FROM public.product_review LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		return nil, model.NewAppErr("PgReviewStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetReviews, http.StatusInternalServerError, nil)
	}

	return reviews, nil
}

// Delete hard deletes the review
func (s PgReviewStore) Delete(id int64) *model.AppErr {
	if _, err := s.db.NamedExec("DELETE from public.product_review WHERE id = :id", map[string]interface{}{"id": id}); err != nil {
		return model.NewAppErr("PgReviewStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteReview, http.StatusInternalServerError, nil)
	}
	return nil
}
