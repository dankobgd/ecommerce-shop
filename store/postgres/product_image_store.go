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
	msgBulkInsertImages   = &i18n.Message{ID: "store.postgres.product_image.bulk_insert.app_error", Other: "could not bulk insert product images"}
	msgGetProductImage    = &i18n.Message{ID: "store.postgres.product_image.get.app_error", Other: "could not get product image"}
	msgGetProductImages   = &i18n.Message{ID: "store.postgres.product_image.get_all.app_error", Other: "could not get product images"}
	msgUpdateProductImage = &i18n.Message{ID: "store.postgres.product_image.update.app_error", Other: "could not update product image"}
	msgDeleteProductImage = &i18n.Message{ID: "store.postgres.product_image.delete.app_error", Other: "could not delete product image"}
)

// BulkInsert inserts multiple images in the db
func (s PgProductImageStore) BulkInsert(images []*model.ProductImage) ([]int64, *model.AppErr) {
	q := `INSERT INTO public.product_image (product_id, url, created_at, updated_at) VALUES(:img_product_id, :img_url, :img_created_at, :img_updated_at) RETURNING id`

	var ids []int64
	rows, err := s.db.NamedQuery(q, images)
	if err != nil {
		return nil, model.NewAppErr("PgProductImageStore.BulkInsertImages", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertImages, http.StatusInternalServerError, nil)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, model.NewAppErr("PgProductImageStore.BulkInsertImages", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertImages, http.StatusInternalServerError, nil)
	}
	return ids, nil
}

// Get gets single image by id
func (s PgProductImageStore) Get(id int64) (*model.ProductImage, *model.AppErr) {
	q := `SELECT img.id AS img_id, img.product_id AS img_product_id, img.url AS img_url, img.created_at AS img_created_at, img.updated_at AS img_updated_at FROM public.product_image img WHERE img.id = $1`
	var img model.ProductImage
	if err := s.db.Get(&img, q, id); err != nil {
		return nil, model.NewAppErr("PgProductImageStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProductImage, http.StatusInternalServerError, nil)
	}
	return &img, nil
}

// GetAll gets all product's images
func (s PgProductImageStore) GetAll(pid int64) ([]*model.ProductImage, *model.AppErr) {
	q := `SELECT img.id AS img_id, img.product_id AS img_product_id, img.url AS img_url, img.created_at AS img_created_at, img.updated_at AS img_updated_at FROM public.product_image img WHERE img.product_id = $1`
	imgs := make([]*model.ProductImage, 0)
	if err := s.db.Select(&imgs, q, pid); err != nil {
		return nil, model.NewAppErr("PgProductImageStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProductImages, http.StatusInternalServerError, nil)
	}
	return imgs, nil
}

// Update updates the image
func (s PgProductImageStore) Update(id int64, pi *model.ProductImage) (*model.ProductImage, *model.AppErr) {
	q := `UPDATE public.product_image SET url=:img_url, created_at=:img_created_at, updated_at=:img_updated_at WHERE id = :img_id`
	if _, err := s.db.NamedExec(q, pi); err != nil {
		return nil, model.NewAppErr("PgProductImageStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateProductImage, http.StatusInternalServerError, nil)
	}
	return pi, nil
}

// Delete deletes the image
func (s PgProductImageStore) Delete(id int64) *model.AppErr {
	if _, err := s.db.NamedExec(`DELETE FROM public.product_image WHERE id=:id`, map[string]interface{}{"id": id}); err != nil {
		return model.NewAppErr("PgProductImageStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteProductImage, http.StatusInternalServerError, nil)
	}
	return nil
}
