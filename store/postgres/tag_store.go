package postgres

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// PgTagStore is the postgres implementation
type PgTagStore struct {
	PgStore
}

// NewPgTagStore creates the new tag store
func NewPgTagStore(pgst *PgStore) store.TagStore {
	return &PgTagStore{*pgst}
}

var (
	msgUniqueConstraintTag = &i18n.Message{ID: "store.postgres.tag.save.unique_constraint.app_error", Other: "invalid tag foreign key"}
	msgSaveTag             = &i18n.Message{ID: "store.postgres.tag.save.app_error", Other: "could not save tag"}
	msgUpdateTag           = &i18n.Message{ID: "store.postgres.tag.update.app_error", Other: "could not update tag"}
	msgBulkInsertTags      = &i18n.Message{ID: "store.postgres.tag.bulk.insert.app_error", Other: "could not bulk insert tags"}
	msgGetTag              = &i18n.Message{ID: "store.postgres.tag.get.app_error", Other: "could not get the tag"}
	msgGetTags             = &i18n.Message{ID: "store.postgres.tag.get.app_error", Other: "could not get the tag"}
	msgDeleteTag           = &i18n.Message{ID: "store.postgres.tag.delete.app_error", Other: "could not delete tag"}
)

// BulkInsert inserts multiple tags in the db
func (s PgTagStore) BulkInsert(tags []*model.Tag) *model.AppErr {
	q := `INSERT INTO public.tag(name, slug, description, created_at, updated_at) VALUES(:name, :slug, :description, :created_at, :updated_at) RETURNING id`

	if _, err := s.db.NamedExec(q, tags); err != nil {
		return model.NewAppErr("PgTagStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertTags, http.StatusInternalServerError, nil)
	}
	return nil
}

// Save inserts the new tag in the db
func (s PgTagStore) Save(tag *model.Tag) (*model.Tag, *model.AppErr) {
	q := `INSERT INTO public.tag(name, slug, description, created_at, updated_at) VALUES(:name, :slug, :description, :created_at, :updated_at) RETURNING id`

	var id int64
	rows, err := s.db.NamedQuery(q, tag)
	if err != nil {
		return nil, model.NewAppErr("PgTagStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveTag, http.StatusInternalServerError, nil)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
	}
	if err := rows.Err(); err != nil {
		if IsUniqueConstraintViolationError(err) {
			return nil, model.NewAppErr("PgTagStore.Save", model.ErrConflict, locale.GetUserLocalizer("en"), msgUniqueConstraintTag, http.StatusInternalServerError, nil)
		}
		return nil, model.NewAppErr("PgTagStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveTag, http.StatusInternalServerError, nil)
	}

	tag.ID = id
	return tag, nil
}

// Update updates the tag
func (s PgTagStore) Update(id int64, tag *model.Tag) (*model.Tag, *model.AppErr) {
	q := `UPDATE public.tag SET name=:name, slug=:slug, description=:description, updated_at=:updated_at WHERE id=:id`
	if _, err := s.db.NamedExec(q, tag); err != nil {
		return nil, model.NewAppErr("PgTagStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateTag, http.StatusInternalServerError, nil)
	}
	return tag, nil
}

// Get gets one tag by id
func (s PgTagStore) Get(id int64) (*model.Tag, *model.AppErr) {
	var tag model.Tag
	if err := s.db.Get(&tag, "SELECT * FROM public.tag WHERE id = $1", id); err != nil {
		return nil, model.NewAppErr("PgTagStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetTag, http.StatusInternalServerError, nil)
	}
	return &tag, nil
}

// GetAll returns all tags
func (s PgTagStore) GetAll() ([]*model.Tag, *model.AppErr) {
	var tags = make([]*model.Tag, 0)
	if err := s.db.Select(&tags, `SELECT * FROM public.tag`); err != nil {
		return nil, model.NewAppErr("PgTagStore.GetAll", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetTags, http.StatusInternalServerError, nil)
	}

	return tags, nil
}

// Delete hard deletes the tag
func (s PgTagStore) Delete(id int64) *model.AppErr {
	if _, err := s.db.NamedExec("DELETE from public.tag WHERE id = :id", map[string]interface{}{"id": id}); err != nil {
		return model.NewAppErr("PgTagStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteTag, http.StatusInternalServerError, nil)
	}
	return nil
}
