package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/jmoiron/sqlx"
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
	msgBulkInsertImages        = &i18n.Message{ID: "store.postgres.product_image.bulk_insert.app_error", Other: "could not bulk insert product images"}
	msgGetProductImage         = &i18n.Message{ID: "store.postgres.product_image.get.app_error", Other: "could not get product image"}
	msgGetProductImages        = &i18n.Message{ID: "store.postgres.product_image.get_all.app_error", Other: "could not get product images"}
	msgUpdateProductImage      = &i18n.Message{ID: "store.postgres.product_image.update.app_error", Other: "could not update product image"}
	msgDeleteProductImage      = &i18n.Message{ID: "store.postgres.product_image.delete.app_error", Other: "could not delete product image"}
	msgBulkDeleteProductImages = &i18n.Message{ID: "store.postgres.product_image.bulk_delete.app_error", Other: "could not bulk delete product images"}
)

// BulkInsert inserts multiple images in the db
func (s PgProductImageStore) BulkInsert(images []*model.ProductImage) *model.AppErr {
	q := `INSERT INTO public.product_image (product_id, url, public_id, created_at, updated_at) VALUES(:product_id, :url, :public_id, :created_at, :updated_at)`
	if _, err := s.db.NamedExec(q, images); err != nil {
		return model.NewAppErr("PgProductImageStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertImages, http.StatusInternalServerError, nil)
	}
	return nil
}

// Save inserts image in the product image table
func (s PgProductImageStore) Save(pid int64, img *model.ProductImage) (*model.ProductImage, *model.AppErr) {
	q := `INSERT INTO public.product_image (product_id, url, public_id, created_at, updated_at) VALUES(:product_id, :url, :public_id, :created_at, :updated_at) RETURNING id`

	var id int64
	rows, err := s.db.NamedQuery(q, map[string]interface{}{"product_id": pid, "url": img.URL, "public_id": img.PublicID, "created_at": img.CreatedAt, "updated_at": img.UpdatedAt})
	if err != nil {
		return nil, model.NewAppErr("PgProductImageStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveBrand, http.StatusInternalServerError, nil)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
	}
	if err := rows.Err(); err != nil {
		if IsUniqueConstraintViolationError(err) {
			return nil, model.NewAppErr("PgProductImageStore.Save", model.ErrConflict, locale.GetUserLocalizer("en"), msgUniqueConstraintBrand, http.StatusInternalServerError, nil)
		}
		return nil, model.NewAppErr("PgProductImageStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveBrand, http.StatusInternalServerError, nil)
	}

	img.ID = &id
	img.ProductID = &pid
	return img, nil
}

// Get gets single image by id
func (s PgProductImageStore) Get(pid, id int64) (*model.ProductImage, *model.AppErr) {
	q := `SELECT img.id AS id, img.product_id AS product_id, img.url AS url, img.public_id AS public_id, img.created_at AS created_at, img.updated_at AS updated_at FROM public.product_image img WHERE img.id = $1 AND product_id= $2`
	var img model.ProductImage
	if err := s.db.Get(&img, q, id, pid); err != nil {
		return nil, model.NewAppErr("PgProductImageStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProductImage, http.StatusInternalServerError, nil)
	}
	return &img, nil
}

// GetAll gets all product's images
func (s PgProductImageStore) GetAll(pid int64) ([]*model.ProductImage, *model.AppErr) {
	q := `SELECT img.id AS id, img.product_id AS product_id, img.url AS url, img.public_id AS public_id, img.created_at AS created_at, img.updated_at AS updated_at FROM public.product_image img WHERE img.product_id = $1`
	imgs := make([]*model.ProductImage, 0)
	if err := s.db.Select(&imgs, q, pid); err != nil {
		return nil, model.NewAppErr("PgProductImageStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetProductImages, http.StatusInternalServerError, nil)
	}
	return imgs, nil
}

// Update updates the image
func (s PgProductImageStore) Update(pid, id int64, pi *model.ProductImage) (*model.ProductImage, *model.AppErr) {
	q := `UPDATE public.product_image SET url=:url, public_id=:public_id, created_at=:created_at, updated_at=:updated_at WHERE id=:id AND product_id=:product_id`
	if _, err := s.db.NamedExec(q, map[string]interface{}{"product_id": pid, "id": id, "url": pi.URL, "public_id": pi.PublicID, "created_at": pi.CreatedAt, "updated_at": pi.UpdatedAt}); err != nil {
		return nil, model.NewAppErr("PgProductImageStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateProductImage, http.StatusInternalServerError, nil)
	}
	return pi, nil
}

// Delete deletes the image
func (s PgProductImageStore) Delete(pid, id int64) *model.AppErr {
	if _, err := s.db.NamedExec(`DELETE FROM public.product_image WHERE product_id=:product_id AND id=:id`, map[string]interface{}{"product_id": pid, "id": id}); err != nil {
		return model.NewAppErr("PgProductImageStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteProductImage, http.StatusInternalServerError, nil)
	}
	return nil
}

// BulkDelete deletes images with given ids
func (s PgProductImageStore) BulkDelete(pid int64, ids []int) *model.AppErr {
	q, args, err := sqlx.In(`DELETE FROM public.product_image WHERE product_id = ? AND id IN (?)`, pid, ids)
	if err != nil {
		return model.NewAppErr("PgProductImageStore.BulkDelete", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkDeleteProductImages, http.StatusInternalServerError, nil)
	}

	if _, err := s.db.Exec(s.db.Rebind(q), args...); err != nil {
		return model.NewAppErr("PgProductImageStore.BulkDelete", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkDeleteProductImages, http.StatusInternalServerError, nil)
	}

	return nil
}
