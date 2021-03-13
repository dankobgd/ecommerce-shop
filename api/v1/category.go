package apiv1

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/utils/pagination"
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
	a.Routes.Categories.Get("/count", a.getCategoriesCount)
	a.Routes.Categories.Get("/", a.getCategories)
	a.Routes.Categories.Get("/featured", a.getFeaturedCategories)
	a.Routes.Categories.Delete("/bulk", a.deleteCategories)
	a.Routes.Category.Get("/", a.getCategory)
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

	c := &model.Category{}
	if err := model.SchemaDecoder.Decode(c, mpf.Value); err != nil {
		respondError(w, model.NewAppErr("createCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgCategoryMultipartErr, http.StatusInternalServerError, nil))
		return
	}

	var fh *multipart.FileHeader
	if len(mpf.File["logo"]) > 0 {
		fh = mpf.File["logo"][0]
	}

	category, cErr := a.app.CreateCategory(c, fh)
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

func (a *API) getCategoriesCount(w http.ResponseWriter, r *http.Request) {
	c := a.app.GetCategoriesCount()
	respondJSON(w, http.StatusOK, map[string]int{"count": c})
}

func (a *API) getCategories(w http.ResponseWriter, r *http.Request) {
	pages := pagination.NewFromRequest(r)
	categories, err := a.app.GetCategories(pages.Limit(), pages.Offset())
	if err != nil {
		respondError(w, err)
		return
	}

	totalCount := -1
	if len(categories) > 0 {
		totalCount = categories[0].TotalCount
	}
	pages.SetData(categories, totalCount)

	respondJSON(w, http.StatusOK, pages)
}

func (a *API) patchCategory(w http.ResponseWriter, r *http.Request) {
	cid, err := strconv.ParseInt(chi.URLParam(r, "category_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("patchCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgCategoryURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	if err := r.ParseMultipartForm(model.FileUploadSizeLimit); err != nil {
		respondError(w, model.NewAppErr("patchCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgCategoryMultipartErr, http.StatusInternalServerError, nil))
		return
	}

	mpf := r.MultipartForm
	model.SchemaDecoder.IgnoreUnknownKeys(true)

	patch := &model.CategoryPatch{}
	if err := model.SchemaDecoder.Decode(patch, mpf.Value); err != nil {
		respondError(w, model.NewAppErr("patchCategory", model.ErrInternal, locale.GetUserLocalizer("en"), msgCategoryMultipartErr, http.StatusInternalServerError, nil))
		return
	}
	patch.SetProperties(patch.PropertiesText)

	var image *multipart.FileHeader
	if len(mpf.File["logo"]) > 0 {
		image = mpf.File["logo"][0]
	}

	ucat, cErr := a.app.PatchCategory(cid, patch, image)
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

func (a *API) getFeaturedCategories(w http.ResponseWriter, r *http.Request) {
	pages := pagination.NewFromRequest(r)
	featured, err := a.app.GetFeaturedCategories(pages.Limit(), pages.Offset())
	if err != nil {
		respondError(w, err)
		return
	}

	totalCount := -1
	if len(featured) > 0 {
		totalCount = featured[0].TotalCount
	}
	pages.SetData(featured, totalCount)

	respondJSON(w, http.StatusOK, pages)
}

func (a *API) deleteCategories(w http.ResponseWriter, r *http.Request) {
	ids := model.IntSliceFromJSON(r.Body)

	if err := a.app.DeleteCategories(ids); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}
