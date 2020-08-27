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
	msgBUlkInsertTags   = &i18n.Message{ID: "store.postgres.product_tag.bulk_insert.app_error", Other: "could not bulk insert product tags"}
	msgGetProductTag    = &i18n.Message{ID: "store.postgres.product_tag.get.app_error", Other: "could not get product tag"}
	msgGetProductTags   = &i18n.Message{ID: "store.postgres.product_tag.get_all.app_error", Other: "could not get product tags"}
	msgUpdateProductTag = &i18n.Message{ID: "store.postgres.product_tag.update.app_error", Other: "could not update product tag"}
	msgDeleteProductTag = &i18n.Message{ID: "store.postgres.product_tag.delete.app_error", Other: "could not delete product tag"}
)

// BulkInsert multiple tags in the db
func (s PgProductTagStore) BulkInsert(tags []*model.ProductTag) ([]int64, *model.AppErr) {
	q := `INSERT INTO public.product_tag(tag_id, product_id) VALUES(:tag_id, :product_id) RETURNING id`

	var ids []int64
	rows, err := s.db.NamedQuery(q, tags)
	if err != nil {
		return nil, model.NewAppErr("PgProductTagStore.BulkInsertTags", model.ErrInternal, locale.GetUserLocalizer("en"), msgBUlkInsertTags, http.StatusInternalServerError, nil)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, model.NewAppErr("PgProductTagStore.BulkInsertTags", model.ErrInternal, locale.GetUserLocalizer("en"), msgBUlkInsertTags, http.StatusInternalServerError, nil)
	}

	return ids, nil
}

// Get gets single tag by id
func (s PgProductTagStore) Get(id int64) (*model.ProductTag, *model.AppErr) {
	q := `SELECT product_tag.*,
	tag.name AS name,
	tag.slug AS slug,
	tag.description AS description
	FROM public.product_tag LEFT JOIN public.tag on product_tag.tag_id = tag.id WHERE product_tag.product_id = $1`
	var tag model.ProductTag
	if err := s.db.Get(&tag, q, id); err != nil {
		return nil, model.NewAppErr("PgProductTagStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProductTag, http.StatusInternalServerError, nil)
	}
	return &tag, nil
}

// GetAll gets all product's tags
func (s PgProductTagStore) GetAll(pid int64) ([]*model.ProductTag, *model.AppErr) {
	q := `SELECT product_tag.*,
	tag.name AS name,
	tag.slug AS slug,
	tag.description AS description
	FROM public.product_tag LEFT JOIN public.tag on product_tag.tag_id = tag.id WHERE product_tag.product_id = $1`
	var tags []*model.ProductTag
	if err := s.db.Select(&tags, q, pid); err != nil {
		return nil, model.NewAppErr("PgProductTagStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProductTags, http.StatusInternalServerError, nil)
	}
	return tags, nil
}

// Update updates the tag
func (s PgProductTagStore) Update(id int64, pt *model.ProductTag) (*model.ProductTag, *model.AppErr) {
	q := `UPDATE public.product_tag SET tag_id=:tag_id WHERE id=:id`
	if _, err := s.db.NamedExec(q, pt); err != nil {
		return nil, model.NewAppErr("PgProductTagStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateProductTag, http.StatusInternalServerError, nil)
	}
	return pt, nil
}

// Delete deletes the tag
func (s PgProductTagStore) Delete(id int64) *model.AppErr {
	if _, err := s.db.NamedExec(`DELETE FROM public.product_tag WHERE id=:id`, map[string]interface{}{"id": id}); err != nil {
		return model.NewAppErr("PgProductTagStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteProductTag, http.StatusInternalServerError, nil)
	}
	return nil
}

// package postgres

// import (
// 	"net/http"

// 	"github.com/dankobgd/ecommerce-shop/model"
// 	"github.com/dankobgd/ecommerce-shop/store"
// 	"github.com/dankobgd/ecommerce-shop/utils/locale"
// 	"github.com/nicksnyder/go-i18n/v2/i18n"
// )

// // PgProductTagStore is the postgres implementation
// type PgProductTagStore struct {
// 	PgStore
// }

// // NewPgProductTagStore creates the new user store
// func NewPgProductTagStore(pgst *PgStore) store.ProductTagStore {
// 	return &PgProductTagStore{*pgst}
// }

// var (
// 	msgBUlkInsertTags   = &i18n.Message{ID: "store.postgres.product_tag.bulk_insert.app_error", Other: "could not bulk insert product tags"}
// 	msgGetProductTag    = &i18n.Message{ID: "store.postgres.product_tag.get.app_error", Other: "could not get product tag"}
// 	msgGetProductTags   = &i18n.Message{ID: "store.postgres.product_tag.get_all.app_error", Other: "could not get product tags"}
// 	msgUpdateProductTag = &i18n.Message{ID: "store.postgres.product_tag.update.app_error", Other: "could not update product tag"}
// 	msgDeleteProductTag = &i18n.Message{ID: "store.postgres.product_tag.delete.app_error", Other: "could not delete product tag"}
// )

// // BulkInsert multiple tags in the db
// func (s PgProductTagStore) BulkInsert(tags []*model.ProductTag) ([]int64, *model.AppErr) {
// 	q := `INSERT INTO public.product_tag(tag_id, product_id) VALUES(:tag_id, :product_id) RETURNING id`

// 	var ids []int64
// 	rows, err := s.db.NamedQuery(q, tags)
// 	if err != nil {
// 		return nil, model.NewAppErr("PgProductTagStore.BulkInsertTags", model.ErrInternal, locale.GetUserLocalizer("en"), msgBUlkInsertTags, http.StatusInternalServerError, nil)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var id int64
// 		rows.Scan(&id)
// 		ids = append(ids, id)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, model.NewAppErr("PgProductTagStore.BulkInsertTags", model.ErrInternal, locale.GetUserLocalizer("en"), msgBUlkInsertTags, http.StatusInternalServerError, nil)
// 	}

// 	return ids, nil
// }

// // Get gets single tag by id
// func (s PgProductTagStore) Get(id int64) (*model.ProductTag, *model.AppErr) {
// 	q := `SELECT * FROM public.product_tag tag WHERE tag.id = $1`
// 	var tag model.ProductTag
// 	if err := s.db.Get(&tag, q, id); err != nil {
// 		return nil, model.NewAppErr("PgProductTagStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProductTag, http.StatusInternalServerError, nil)
// 	}
// 	return &tag, nil
// }

// // GetAll gets all product's tags
// func (s PgProductTagStore) GetAll(pid int64) ([]*model.ProductTag, *model.AppErr) {
// 	q := `SELECT * FROM public.product_tag tag WHERE tag.product_id = $1`
// 	var tags []*model.ProductTag
// 	if err := s.db.Select(&tags, q, pid); err != nil {
// 		return nil, model.NewAppErr("PgProductTagStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProductTags, http.StatusInternalServerError, nil)
// 	}
// 	return tags, nil
// }

// // Update updates the tag
// func (s PgProductTagStore) Update(id int64, pt *model.ProductTag) (*model.ProductTag, *model.AppErr) {
// 	q := `UPDATE public.product_tag SET tag_id=:tag_id WHERE id=:id`
// 	if _, err := s.db.NamedExec(q, pt); err != nil {
// 		return nil, model.NewAppErr("PgProductTagStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateProductTag, http.StatusInternalServerError, nil)
// 	}
// 	return pt, nil
// }

// // Delete deletes the tag
// func (s PgProductTagStore) Delete(id int64) *model.AppErr {
// 	if _, err := s.db.NamedExec(`DELETE FROM public.product_tag WHERE id=:id`, map[string]interface{}{"id": id}); err != nil {
// 		return model.NewAppErr("PgProductTagStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteProductTag, http.StatusInternalServerError, nil)
// 	}
// 	return nil
// }
