package postgres

import (
	"net/http"
	"time"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/store"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// PgUserStore is the postgres implementation
type PgUserStore struct {
	PgStore
}

// NewPgUserStore creates the new user store
func NewPgUserStore(pgst *PgStore) store.UserStore {
	return &PgUserStore{*pgst}
}

var (
	msgUniqueConstraintUser = &i18n.Message{ID: "store.postgres.user.save.unique_constraint.app_error", Other: "invalid credentials"}
	msgSaveUser             = &i18n.Message{ID: "store.postgres.user.save.app_error", Other: "could not save user"}
	msgUpdateUserProfile    = &i18n.Message{ID: "store.postgres.user.update.app_error", Other: "could not update user"}
	msgBulkInsertUsers      = &i18n.Message{ID: "store.postgres.user.bulk.insert.app_error", Other: "could not bulk insert users"}
	msgGetUser              = &i18n.Message{ID: "store.postgres.user.get.app_error", Other: "could not get the user"}
	msgVerifyEmail          = &i18n.Message{ID: "store.postgres.user.verify_email.app_error", Other: "could not verify email"}
	msgDeleteToken          = &i18n.Message{ID: "store.postgres.user.verify_email.delete_token.app_error", Other: "could not delete verify token"}
	msgUpdatePassword       = &i18n.Message{ID: "store.postgres.user.update_password.app_error", Other: "could not update password"}
	msgDeleteUser           = &i18n.Message{ID: "store.postgres.user.delete.app_error", Other: "could not delete user"}
	msgUpdateUserAvatar     = &i18n.Message{ID: "store.postgres.user.update_avatar.app_error", Other: "could not delete user avatar"}
	msgDeleteUserAvatar     = &i18n.Message{ID: "store.postgres.user.delete_avatar.app_error", Other: "could not delete user avatar"}

	msgCreateWishlist = &i18n.Message{ID: "store.postgres.user.create_wishlist.app_error", Other: "could not add product to wishlist"}
	msgGetWishlist    = &i18n.Message{ID: "store.postgres.user.get_wishlist.app_error", Other: "could not get wishlist"}
	msgDeleteWishlist = &i18n.Message{ID: "store.postgres.user.delete_wishlist.app_error", Other: "could not delete product from wishlist"}
	msgClearWishlist  = &i18n.Message{ID: "store.postgres.user.clear_wishlist.app_error", Other: "could not delete all products from wishlist"}
)

// BulkInsert inserts multiple users in the db
func (s PgUserStore) BulkInsert(users []*model.User) *model.AppErr {
	q := `INSERT INTO public.user(first_name, last_name, username, email, password, role, gender, locale, avatar_url, active, email_verified, failed_attempts, last_login_at, created_at, updated_at, deleted_at) 
	VALUES(:first_name, :last_name, :username, :email, :password, :role, :gender, :locale, :avatar_url, :active, :email_verified, :failed_attempts, :last_login_at, :created_at, :updated_at, :deleted_at) RETURNING id`

	if _, err := s.db.NamedExec(q, users); err != nil {
		return model.NewAppErr("PgUserStore.BulkInsert", model.ErrInternal, locale.GetUserLocalizer("en"), msgBulkInsertUsers, http.StatusInternalServerError, nil)
	}
	return nil
}

// Save inserts the new user in the db
func (s PgUserStore) Save(user *model.User) (*model.User, *model.AppErr) {
	q := `INSERT INTO public.user (first_name, last_name, username, email, password, role, gender, locale, avatar_url, active, email_verified, failed_attempts, last_login_at, created_at, updated_at, deleted_at) 
	VALUES (:first_name, :last_name, :username, :email, :password, :role, :gender, :locale, :avatar_url, :active, :email_verified, :failed_attempts, :last_login_at, :created_at, :updated_at, :deleted_at) RETURNING id`

	var id int64
	rows, err := s.db.NamedQuery(q, user)
	if err != nil {
		return nil, model.NewAppErr("PgUserStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveUser, http.StatusInternalServerError, nil)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
	}
	if err := rows.Err(); err != nil {
		if IsUniqueConstraintViolationError(err) {
			return nil, model.NewAppErr("PgUserStore.Save", model.ErrConflict, locale.GetUserLocalizer("en"), msgUniqueConstraintUser, http.StatusInternalServerError, nil)
		}
		return nil, model.NewAppErr("PgUserStore.Save", model.ErrInternal, locale.GetUserLocalizer("en"), msgSaveUser, http.StatusInternalServerError, nil)
	}
	user.ID = id
	return user, nil
}

// Update updates the user profile
func (s PgUserStore) Update(id int64, u *model.User) (*model.User, *model.AppErr) {
	q := `UPDATE public.user SET first_name=:first_name, last_name=:last_name, username=:username, email=:email, gender=:gender, locale=:locale, updated_at=:updated_at WHERE id=:id`
	if _, err := s.db.NamedExec(q, u); err != nil {
		return nil, model.NewAppErr("PgUserStore.Update", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateUserProfile, http.StatusInternalServerError, nil)
	}
	return u, nil
}

// Get gets one user by id
func (s PgUserStore) Get(id int64) (*model.User, *model.AppErr) {
	var user model.User
	if err := s.db.Get(&user, "SELECT * FROM public.user WHERE id = $1 AND deleted_at IS NULL", id); err != nil {
		return nil, model.NewAppErr("PgUserStore.Get", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetUser, http.StatusInternalServerError, nil)
	}
	return &user, nil
}

// GetByEmail gets one user by email
func (s PgUserStore) GetByEmail(email string) (*model.User, *model.AppErr) {
	var user model.User
	if err := s.db.Get(&user, "SELECT * FROM public.user WHERE email = $1 AND deleted_at IS NULL", email); err != nil {
		return nil, model.NewAppErr("PgUserStore.GetByEmail", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetUser, http.StatusInternalServerError, nil)
	}
	return &user, nil
}

// GetAll returns all users
func (s PgUserStore) GetAll(limit, offset int) ([]*model.User, *model.AppErr) {
	return []*model.User{}, nil
}

// VerifyEmail updates the email_verified field
func (s PgUserStore) VerifyEmail(id int64) *model.AppErr {
	m := map[string]interface{}{"updated_at": time.Now(), "id": id}
	if _, err := s.db.NamedExec("UPDATE public.user SET updated_at = :updated_at, email_verified = true WHERE id = :id", m); err != nil {
		return model.NewAppErr("PgUserStore.VerifyEmail", model.ErrInternal, locale.GetUserLocalizer("en"), msgVerifyEmail, http.StatusInternalServerError, nil)
	}
	return nil
}

// UpdatePassword updates the user's password
func (s PgUserStore) UpdatePassword(userID int64, hashedPassword string) *model.AppErr {
	m := map[string]interface{}{"id": userID, "password": hashedPassword, "updated_at": time.Now()}
	if _, err := s.db.NamedExec("UPDATE public.user SET password = :password, updated_at = :updated_at WHERE id = :id", m); err != nil {
		return model.NewAppErr("PgUserStore.UpdatePassword", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdatePassword, http.StatusInternalServerError, nil)
	}
	return nil
}

// Delete soft deletes the user
func (s PgUserStore) Delete(id int64) *model.AppErr {
	m := map[string]interface{}{"id": id, "deleted_at": time.Now()}
	if _, err := s.db.NamedExec("UPDATE public.user SET deleted_at = :deleted_at WHERE id = :id", m); err != nil {
		return model.NewAppErr("PgUserStore.Delete", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteUser, http.StatusInternalServerError, nil)
	}
	return nil
}

// UpdateAvatar updates the user avatar image
func (s PgUserStore) UpdateAvatar(id int64, url *string, publicID *string) (*string, *string, *model.AppErr) {
	m := map[string]interface{}{"id": id, "avatar_url": url, "avatar_public_id": publicID, "updated_at": time.Now()}
	if _, err := s.db.NamedExec("UPDATE public.user SET avatar_url = :avatar_url, avatar_public_id = :avatar_public_id, updated_at = :updated_at WHERE id = :id", m); err != nil {
		return model.NewString(""), model.NewString(""), model.NewAppErr("PgUserStore.UpdateAvatar", model.ErrInternal, locale.GetUserLocalizer("en"), msgUpdateUserAvatar, http.StatusInternalServerError, nil)
	}
	return url, publicID, nil
}

// DeleteAvatar deletes the user avatar image
func (s PgUserStore) DeleteAvatar(id int64) *model.AppErr {
	m := map[string]interface{}{"id": id, "updated_at": time.Now()}
	if _, err := s.db.NamedExec("UPDATE public.user SET avatar_url = NULL, avatar_public_id = NULL, updated_at = :updated_at WHERE id = :id", m); err != nil {
		return model.NewAppErr("PgUserStore.DeleteAvatar", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteUserAvatar, http.StatusInternalServerError, nil)
	}
	return nil
}

// GetAllOrders returns all orders for the user
func (s PgUserStore) GetAllOrders(uid int64, limit, offset int) ([]*model.Order, *model.AppErr) {
	var orders = make([]*model.Order, 0)
	if err := s.db.Select(&orders, `SELECT COUNT(*) OVER() AS total_count, * FROM public.order WHERE user_id = $1 LIMIT $2 OFFSET $3`, uid, limit, offset); err != nil {
		return nil, model.NewAppErr("PgUserStore.GetAllOrders", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetOrders, http.StatusInternalServerError, nil)
	}

	return orders, nil
}

// CreateWishlist adds product to wishlist
func (s PgUserStore) CreateWishlist(userID, productID int64) *model.AppErr {
	q := `INSERT INTO public.product_wishlist (user_id, product_id) VALUES (:user_id, :product_id)`

	if _, err := s.db.NamedExec(q, map[string]interface{}{"user_id": userID, "product_id": productID}); err != nil {
		return model.NewAppErr("PgUserStore.CreateWishlist", model.ErrInternal, locale.GetUserLocalizer("en"), msgCreateWishlist, http.StatusInternalServerError, nil)
	}
	return nil
}

// GetWishlist gets users wishlist products
func (s PgUserStore) GetWishlist(userID int64) ([]*model.Product, *model.AppErr) {
	q := `SELECT 
	p.*,
	b.name AS brand_name,
	b.slug AS brand_slug,
	b.type AS brand_type,
	b.description AS brand_description,
	b.email AS brand_email,
	b.logo AS brand_logo,
	b.website_url AS brand_website_url,
	b.created_at AS brand_created_at,
	b.updated_at AS brand_updated_at,
	c.name AS category_name,
	c.slug AS category_slug,
	c.description AS category_description,
	c.logo AS category_logo,
	c.created_at AS category_created_at,
	c.updated_at AS category_updated_at
	FROM public.product p
	LEFT JOIN brand b ON p.brand_id = b.id
	LEFT JOIN category c ON p.category_id = c.id	
	LEFT JOIN product_wishlist w ON p.id = w.product_id
	LEFT JOIN public.user u ON u.id = w.user_id
	WHERE u.id = $1`

	var pj []productJoin
	if err := s.db.Select(&pj, q, userID); err != nil {
		return nil, model.NewAppErr("PgUserStore.GetWishlist", model.ErrInternal, locale.GetUserLocalizer("en"), msgGetWishlist, http.StatusInternalServerError, nil)
	}

	wishlist := make([]*model.Product, 0)
	for _, x := range pj {
		wishlist = append(wishlist, x.ToProduct())
	}

	return wishlist, nil
}

// DeleteWishlist deletes product from wishlist
func (s PgUserStore) DeleteWishlist(userID, productID int64) *model.AppErr {
	m := map[string]interface{}{"user_id": userID, "product_id": productID}
	if _, err := s.db.NamedExec("DELETE from public.product_wishlist WHERE user_id = :user_id AND product_id = :product_id", m); err != nil {
		return model.NewAppErr("PgUserStore.DeleteWishlist", model.ErrInternal, locale.GetUserLocalizer("en"), msgDeleteWishlist, http.StatusInternalServerError, nil)
	}
	return nil
}

// ClearWishlist clears whole wishlist
func (s PgUserStore) ClearWishlist(userID int64) *model.AppErr {
	if _, err := s.db.NamedExec("DELETE from public.product_wishlist WHERE user_id = :user_id", map[string]interface{}{"user_id": userID}); err != nil {
		return model.NewAppErr("PgUserStore.ClearWishlist", model.ErrInternal, locale.GetUserLocalizer("en"), msgClearWishlist, http.StatusInternalServerError, nil)
	}
	return nil
}
