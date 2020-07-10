package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// PgProductTagStore is the postgres implementation
type PgProductTagStore struct {
	PgStore
}

// NewPgProductTagStore creates the new user store
func NewPgProductTagStore(pgst *PgStore) store.ProductTagStore {
	return &PgProductTagStore{*pgst}
}

var (
	msgBUlkInsertTags = &i18n.Message{ID: "store.postgres.product_tag.bulk_insert.app_error", Other: "could not bulk insert product tags"}
)

// BulkInsert multiple tags in the db
func (s PgProductTagStore) BulkInsert(tags []*model.ProductTag) ([]int64, *model.AppErr) {
	q := `INSERT INTO public.product_tag (product_id, name, created_at, updated_at) VALUES(:tag_product_id, :tag_name, :tag_created_at, :tag_updated_at) RETURNING id`

	var ids []int64
	rows, err := s.db.NamedQuery(q, tags)
	defer rows.Close()
	if err != nil {
		return nil, model.NewAppErr("PgProductStore.BulkInsertTags", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, model.NewAppErr("PgProductStore.BulkInsertTags", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveProduct, http.StatusInternalServerError, nil)
	}

	return ids, nil
}
