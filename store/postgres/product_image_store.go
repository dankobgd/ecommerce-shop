package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// PgProductImageStore is the postgres implementation
type PgProductImageStore struct {
	PgStore
}

// NewPgProductImageStore creates the new user store
func NewPgProductImageStore(pgst *PgStore) store.ProductImageStore {
	return &PgProductImageStore{*pgst}
}

var (
	msgBulkInsertImages = &i18n.Message{ID: "store.postgres.product_image.bulk_insert.app_error", Other: "could not bulk insert product images"}
)

// BulkInsert inserts multiple images in the db
func (s PgProductImageStore) BulkInsert(images []*model.ProductImage) ([]int64, *model.AppErr) {
	q := `INSERT INTO public.product_image (product_id, url, created_at, updated_at) VALUES(:img_product_id, :img_url, :img_created_at, :img_updated_at) RETURNING id`

	var ids []int64
	rows, err := s.db.NamedQuery(q, images)
	defer rows.Close()
	if err != nil {
		return nil, model.NewAppErr("PgProductStore.BulkInsertImages", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertImages, http.StatusInternalServerError, nil)
	}
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, model.NewAppErr("PgProductStore.BulkInsertImages", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertImages, http.StatusInternalServerError, nil)
	}
	return ids, nil
}
