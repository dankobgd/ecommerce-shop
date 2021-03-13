package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/jmoiron/sqlx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// PgCategoryStore is the postgres implementation
type PgCategoryStore struct {
	PgStore
}

// NewPgCategoryStore creates the new category store
func NewPgCategoryStore(pgst *PgStore) store.CategoryStore {
	return &PgCategoryStore{*pgst}
}

var (
	msgUniqueConstraintCategory = &i18n.Message{ID: "store.postgres.category.save.unique_constraint.app_error", Other: "invalid category, it already exists"}
	msgSaveCategory             = &i18n.Message{ID: "store.postgres.category.save.app_error", Other: "could not save category"}
	msgUpdateCategory           = &i18n.Message{ID: "store.postgres.category.update.app_error", Other: "could not update category"}
	msgBulkInsertCategories     = &i18n.Message{ID: "store.postgres.category.bulk.insert.app_error", Other: "could not bulk insert categories"}
	msgGetCategory              = &i18n.Message{ID: "store.postgres.category.get.app_error", Other: "could not get the category"}
	msgGetCategories            = &i18n.Message{ID: "store.postgres.category.get.app_error", Other: "could not get categories"}
	msgDeleteCategory           = &i18n.Message{ID: "store.postgres.category.delete.app_error", Other: "could not delete category"}
	msgBulkDeleteCategories     = &i18n.Message{ID: "store.postgres.category.delete.app_error", Other: "could not bulk delete categories"}
)

// Count returns the total categories count
func (s PgCategoryStore) Count() int {
	var n int
	s.db.Get(&n, "SELECT COUNT(*) FROM public.category")
	return n
}

// BulkInsert inserts multiple categories in the db
func (s PgCategoryStore) BulkInsert(categories []*model.Category) *model.AppErr {
	q := `INSERT INTO public.category(name, slug, logo, logo_public_id, description, is_featured, properties, created_at, updated_at) VALUES(:name, :slug, :logo, :logo_public_id, :description, :is_featured, :properties, :created_at, :updated_at) RETURNING id`

	if _, err := s.db.NamedExec(q, categories); err != nil {
		return model.NewAppErr("PgCategoryStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertCategories, http.StatusInternalServerError, nil)
	}
	return nil
}

// Save inserts the new category in the db
func (s PgCategoryStore) Save(category *model.Category) (*model.Category, *model.AppErr) {
	q := `INSERT INTO public.category(name, slug, logo, logo_public_id, description, is_featured, properties, created_at, updated_at) VALUES(:name, :slug, :logo, :logo_public_id, :description, :is_featured, :properties, :created_at, :updated_at) RETURNING id`

	var id int64
	rows, err := s.db.NamedQuery(q, category)
	if err != nil {
		return nil, model.NewAppErr("PgCategoryStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveCategory, http.StatusInternalServerError, nil)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
	}
	if err := rows.Err(); err != nil {
		if IsUniqueConstraintViolationError(err) {
			return nil, model.NewAppErr("PgCategoryStore.Save", model.ErrConflict, locale.GetUserLocalizer("en"), msgUniqueConstraintCategory, http.StatusInternalServerError, nil)
		}
		return nil, model.NewAppErr("PgCategoryStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveCategory, http.StatusInternalServerError, nil)
	}

	category.ID = id
	return category, nil
}

// Update updates the category
func (s PgCategoryStore) Update(id int64, category *model.Category) (*model.Category, *model.AppErr) {
	q := `UPDATE public.category SET name=:name, slug=:slug, description=:description, is_featured=:is_featured, properties=:properties, logo=:logo, logo_public_id=:logo_public_id, updated_at=:updated_at WHERE id=:id`
	if _, err := s.db.NamedExec(q, category); err != nil {
		return nil, model.NewAppErr("PgCategoryStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateCategory, http.StatusInternalServerError, nil)
	}
	return category, nil
}

// Get gets one category by id
func (s PgCategoryStore) Get(id int64) (*model.Category, *model.AppErr) {
	var category model.Category
	if err := s.db.Get(&category, "SELECT * FROM public.category WHERE id = $1", id); err != nil {
		return nil, model.NewAppErr("PgCategoryStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetCategory, http.StatusInternalServerError, nil)
	}
	return &category, nil
}

// GetAll returns all categories
func (s PgCategoryStore) GetAll(limit, offset int) ([]*model.Category, *model.AppErr) {
	var categories = make([]*model.Category, 0)
	if err := s.db.Select(&categories, `SELECT COUNT(*) OVER() AS total_count, * FROM public.category ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		return nil, model.NewAppErr("PgCategoryStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetCategories, http.StatusInternalServerError, nil)
	}

	return categories, nil
}

// GetFeatured returns all featured categories
func (s PgCategoryStore) GetFeatured(limit, offset int) ([]*model.Category, *model.AppErr) {
	var categories = make([]*model.Category, 0)
	if err := s.db.Select(&categories, `SELECT COUNT(*) OVER() AS total_count, * FROM public.category WHERE is_featured = true ORDER BY id DESC LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		return nil, model.NewAppErr("PgCategoryStore.GetFeatured", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetCategories, http.StatusInternalServerError, nil)
	}

	return categories, nil
}

// Delete deletes the category
func (s PgCategoryStore) Delete(id int64) *model.AppErr {
	if _, err := s.db.NamedExec("DELETE from public.category WHERE id = :id", map[string]interface{}{"id": id}); err != nil {
		return model.NewAppErr("PgCategoryStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteUserAvatar, http.StatusInternalServerError, nil)
	}
	return nil
}

// BulkDelete deletes categories with given ids
func (s PgCategoryStore) BulkDelete(ids []int) *model.AppErr {
	q, args, err := sqlx.In(`DELETE FROM public.category WHERE id IN (?)`, ids)
	if err != nil {
		return model.NewAppErr("PgCategoryStore.BulkDelete", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkDeleteCategories, http.StatusInternalServerError, nil)
	}

	if _, err := s.db.Exec(s.db.Rebind(q), args...); err != nil {
		return model.NewAppErr("PgCategoryStore.BulkDelete", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkDeleteCategories, http.StatusInternalServerError, nil)
	}

	return nil
}
