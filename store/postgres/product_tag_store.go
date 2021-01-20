package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/jmoiron/sqlx"
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
	msgBUlkInsertTags     = &i18n.Message{ID: "store.postgres.product_tag.bulk_insert.app_error", Other: "could not bulk insert product tags"}
	msgGetProductTag      = &i18n.Message{ID: "store.postgres.product_tag.get.app_error", Other: "could not get product tag"}
	msgGetProductTags     = &i18n.Message{ID: "store.postgres.product_tag.get_all.app_error", Other: "could not get product tags"}
	msgUpdateProductTag   = &i18n.Message{ID: "store.postgres.product_tag.update.app_error", Other: "could not update product tag"}
	msgReplaceProductTags = &i18n.Message{ID: "store.postgres.product_tag.replace.app_error", Other: "could not replace product tags"}
	msgDeleteProductTag   = &i18n.Message{ID: "store.postgres.product_tag.delete.app_error", Other: "could not delete product tag"}
)

// BulkInsert multiple tags in the db
func (s PgProductTagStore) BulkInsert(tags []*model.ProductTag) *model.AppErr {
	if _, err := s.db.NamedExec(`INSERT INTO public.product_tag(tag_id, product_id) VALUES(:tag_id, :product_id)`, tags); err != nil {
		return model.NewAppErr("PgProductTagStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBUlkInsertTags, http.StatusInternalServerError, nil)
	}
	return nil
}

// Save multiple tags in the db
func (s PgProductTagStore) Save(pid int64, pt *model.ProductTag) (*model.ProductTag, *model.AppErr) {
	if _, err := s.db.NamedExec(`INSERT INTO public.product_tag(tag_id, product_id) VALUES(:tag_id, :product_id)`, map[string]interface{}{"tag_id": pt.TagID, "product_id": pid}); err != nil {
		return nil, model.NewAppErr("PgProductTagStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgBUlkInsertTags, http.StatusInternalServerError, nil)
	}
	pt.ProductID = &pid
	return pt, nil
}

// Get gets single tag by id
func (s PgProductTagStore) Get(pid, tid int64) (*model.ProductTag, *model.AppErr) {
	q := `SELECT product_tag.*,
	tag.name AS name,
	tag.slug AS slug,
	tag.description AS description
	FROM public.product_tag LEFT JOIN public.tag on product_tag.tag_id = tag.id WHERE product_tag.product_id = $1 AND product_tag.tag_id = $2`
	var tag model.ProductTag
	if err := s.db.Get(&tag, q, pid, tid); err != nil {
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
	tags := make([]*model.ProductTag, 0)
	if err := s.db.Select(&tags, q, pid); err != nil {
		return nil, model.NewAppErr("PgProductTagStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProductTags, http.StatusInternalServerError, nil)
	}
	return tags, nil
}

// Update updates the tag
func (s PgProductTagStore) Update(pid, tid int64, pt *model.ProductTag) (*model.ProductTag, *model.AppErr) {
	q := `UPDATE public.product_tag SET tag_id=:tag_id WHERE product_id=:product_id AND tag_id=:tid`
	if _, err := s.db.NamedExec(q, map[string]interface{}{"product_id": pid, "tag_id": pt.TagID, "tid": tid}); err != nil {
		return nil, model.NewAppErr("PgProductTagStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateProductTag, http.StatusInternalServerError, nil)
	}
	return pt, nil
}

// Replace replaces the tags with provided ones (deletes and inserts)
func (s PgProductTagStore) Replace(pid int64, tagIDs []int) ([]*model.ProductTag, *model.AppErr) {
	ptags := make([]*model.ProductTag, 0)

	if len(tagIDs) == 0 {
		if _, err := s.db.NamedExec(`DELETE FROM product_tag WHERE product_id=:product_id `, map[string]interface{}{"product_id": pid}); err != nil {
			return nil, model.NewAppErr("PgProductStore.Replace", model.ErrInternal, locale.GetUserLocalizer("en"), msgReplaceProductTags, http.StatusInternalServerError, nil)
		}
		return ptags, nil
	}

	for _, id := range tagIDs {
		ptags = append(ptags, &model.ProductTag{TagID: model.NewInt64(int64(id)), ProductID: model.NewInt64(pid)})
	}

	tx, txErr := s.db.Beginx()

	if txErr != nil {
		tx.Rollback()
		return nil, model.NewAppErr("PgProductStore.Replace", model.ErrInternal, locale.GetUserLocalizer("en"), msgReplaceProductTags, http.StatusInternalServerError, nil)
	}

	if _, err := tx.NamedExec(`INSERT INTO product_tag(tag_id, product_id) VALUES(:tag_id, :product_id) ON CONFLICT (product_id, tag_id) DO NOTHING`, ptags); err != nil {
		tx.Rollback()
		return nil, model.NewAppErr("PgProductStore.Replace", model.ErrInternal, locale.GetUserLocalizer("en"), msgReplaceProductTags, http.StatusInternalServerError, nil)
	}

	q, args, e := sqlx.In(`DELETE FROM product_tag WHERE product_id = ? AND tag_id NOT IN (?)`, pid, tagIDs)
	if e != nil {
		tx.Rollback()
		return nil, model.NewAppErr("PgProductStore.Replace", model.ErrInternal, locale.GetUserLocalizer("en"), msgReplaceProductTags, http.StatusInternalServerError, nil)
	}
	if _, err := tx.Exec(s.db.Rebind(q), args...); err != nil {
		tx.Rollback()
		return nil, model.NewAppErr("PgProductStore.Replace", model.ErrInternal, locale.GetUserLocalizer("en"), msgReplaceProductTags, http.StatusInternalServerError, nil)
	}

	tx.Commit()

	return ptags, nil
}

// Delete deletes the tag
func (s PgProductTagStore) Delete(pid, tid int64) *model.AppErr {
	if _, err := s.db.NamedExec(`DELETE FROM public.product_tag WHERE product_id=:product_id AND tag_id=:tag_id`, map[string]interface{}{"product_id": pid, "tag_id": tid}); err != nil {
		return model.NewAppErr("PgProductTagStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteProductTag, http.StatusInternalServerError, nil)
	}
	return nil
}

// BulkDelete deletes tags with given ids
func (s PgProductTagStore) BulkDelete(pid int64, ids []int) *model.AppErr {
	q, args, err := sqlx.In(`DELETE FROM public.product_tag WHERE product_id = ? AND tag_id IN (?)`, pid, ids)
	if err != nil {
		return model.NewAppErr("PgProductTagStore.BulkDelete", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkDeleteTags, http.StatusInternalServerError, nil)
	}

	if _, err := s.db.Exec(s.db.Rebind(q), args...); err != nil {
		return model.NewAppErr("PgProductTagStore.BulkDelete", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkDeleteTags, http.StatusInternalServerError, nil)
	}

	return nil
}
