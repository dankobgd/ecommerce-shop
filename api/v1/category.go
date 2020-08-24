package apiv1

import (
	"net/http"
	"strconv"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/go-chi/chi"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgCategoryCreateErr        = &i18n.Message{ID: "api.category.create_category.app_error", Other: "could not create category"}
	msgCategoriesGetErr         = &i18n.Message{ID: "api.category.get_categories.app_error", Other: "could not get categories"}
	msgCategoryGetErr           = &i18n.Message{ID: "api.category.get_category.app_error", Other: "could not get category"}
	msgCategoryPatchErr         = &i18n.Message{ID: "api.category.patch_category.app_error", Other: "could not update category"}
	msgCategoryDeleteerr        = &i18n.Message{ID: "api.category.delete_category.app_error", Other: "could not delete category"}
	msgCategoryMultipartErr     = &i18n.Message{ID: "api.category.create_category.multipart.app_error", Other: "could not decode category multipart data"}
	msgCategoryURLParamErr      = &i18n.Message{ID: "api.category.url.params.app_error", Other: "could not parse URL params"}
	msgCategoryPatchFromJSONErr = &i18n.Message{ID: "api.category.patch_product.app_error", Other: "could not decode product patch data"}
)

// InitCategories inits the category routes
func InitCategories(a *API) {
	a.Routes.Categories.Post("/", a.AdminSessionRequired(a.createCategory))
	a.Routes.Categories.Get("/", a.AdminSessionRequired(a.getCategories))
	a.Routes.Category.Get("/", a.AdminSessionRequired(a.getCategory))
	a.Routes.Category.Patch("/", a.AdminSessionRequired(a.patchCategory))
	a.Routes.Category.Delete("/", a.AdminSessionRequired(a.deleteCategory))
}

func (a *API) createCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(model.FileUploadSizeLimit); err != nil {
		respondError(w, model.NewAppErr("createCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgCategoryMultipartErr, http.StatusInternalServerError, nil))
		return
	}

	mpf := r.MultipartForm
	model.SchemaDecoder.IgnoreUnknownKeys(true)

	var c model.Category
	if err := model.SchemaDecoder.Decode(&c, mpf.Value); err != nil {
		respondError(w, model.NewAppErr("createCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgCategoryMultipartErr, http.StatusInternalServerError, nil))
		return
	}

	fh := mpf.File["logo"][0]

	category, cErr := a.app.CreateCategory(&c, fh)
	if cErr != nil {
		respondError(w, cErr)
		return
	}
	respondJSON(w, http.StatusCreated, category)
}

func (a *API) getCategory(w http.ResponseWriter, r *http.Request) {
	cid, e := strconv.ParseInt(chi.URLParam(r, "category_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	c, err := a.app.GetCategory(cid)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, c)
}

func (a *API) getCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := a.app.GetCategories()
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, categories)
}

func (a *API) patchCategory(w http.ResponseWriter, r *http.Request) {
	cid, err := strconv.ParseInt(chi.URLParam(r, "category_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("patchCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgCategoryURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	patch, err := model.CategoryPatchFromJSON(r.Body)
	if err != nil {
		respondError(w, model.NewAppErr("patchCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgCategoryPatchFromJSONErr, http.StatusInternalServerError, nil))
		return
	}

	ucat, cErr := a.app.PatchCategory(cid, patch)
	if err != nil {
		respondError(w, cErr)
		return
	}

	respondJSON(w, http.StatusOK, ucat)
}

func (a *API) deleteCategory(w http.ResponseWriter, r *http.Request) {
	cid, err := strconv.ParseInt(chi.URLParam(r, "category_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("deleteCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgCategoryURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	if err := a.app.DeleteCategory(cid); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}
