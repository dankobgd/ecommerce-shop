package apiv1

import (
	"net/http"
	"strconv"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/utils/pagination"
	"github.com/go-chi/chi"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgBrandCreateErr        = &i18n.Message{ID: "api.brand.create_brand.app_error", Other: "could not create brand"}
	msgBrandsGetErr          = &i18n.Message{ID: "api.brand.get_brands.app_error", Other: "could not get brands"}
	msgBrandGetErr           = &i18n.Message{ID: "api.brand.get_brand.app_error", Other: "could not get brand"}
	msgBrandPatchErr         = &i18n.Message{ID: "api.brand.patch_brand.app_error", Other: "could not update brand"}
	msgBrandDeleteerr        = &i18n.Message{ID: "api.brand.delete_brand.app_error", Other: "could not delete brand"}
	msgBrandMultipartErr     = &i18n.Message{ID: "api.brand.create_brand.multipart.app_error", Other: "could not decode brand multipart data"}
	msgBrandURLParamErr      = &i18n.Message{ID: "api.brand.url.params.app_error", Other: "could not parse URL params"}
	msgBrandPatchFromJSONErr = &i18n.Message{ID: "api.brand.patch_brand.app_error", Other: "could not decode brand patch data"}
)

// InitBrands inits the brand routes
func InitBrands(a *API) {
	a.Routes.Brands.Post("/", a.AdminSessionRequired(a.createBrand))
	a.Routes.Brands.Get("/", a.AdminSessionRequired(a.getBrands))
	a.Routes.Brand.Get("/", a.AdminSessionRequired(a.getBrand))
	a.Routes.Brand.Patch("/", a.AdminSessionRequired(a.patchBrand))
	a.Routes.Brand.Delete("/", a.AdminSessionRequired(a.deleteBrand))
}

func (a *API) createBrand(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(model.FileUploadSizeLimit); err != nil {
		respondError(w, model.NewAppErr("createBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgBrandMultipartErr, http.StatusInternalServerError, nil))
		return
	}

	mpf := r.MultipartForm
	model.SchemaDecoder.IgnoreUnknownKeys(true)

	var b model.Brand
	if err := model.SchemaDecoder.Decode(&b, mpf.Value); err != nil {
		respondError(w, model.NewAppErr("createBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgBrandMultipartErr, http.StatusInternalServerError, nil))
		return
	}

	fh := mpf.File["logo"][0]

	brand, bErr := a.app.CreateBrand(&b, fh)
	if bErr != nil {
		respondError(w, bErr)
		return
	}
	respondJSON(w, http.StatusCreated, brand)
}

func (a *API) getBrand(w http.ResponseWriter, r *http.Request) {
	bid, e := strconv.ParseInt(chi.URLParam(r, "brand_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgBrandURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	b, err := a.app.GetBrand(bid)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, b)
}

func (a *API) getBrands(w http.ResponseWriter, r *http.Request) {
	pages := pagination.NewFromRequest(r)
	brands, err := a.app.GetBrands(pages.Limit(), pages.Offset())
	if err != nil {
		respondError(w, err)
		return
	}

	totalCount := -1
	if len(brands) > 0 {
		totalCount = brands[0].TotalCount
	}
	pages.SetData(brands, totalCount)

	respondJSON(w, http.StatusOK, pages)
}

func (a *API) patchBrand(w http.ResponseWriter, r *http.Request) {
	bid, err := strconv.ParseInt(chi.URLParam(r, "brand_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("patchBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgBrandURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	patch, err := model.BrandPatchFromJSON(r.Body)
	if err != nil {
		respondError(w, model.NewAppErr("patchBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgBrandPatchFromJSONErr, http.StatusInternalServerError, nil))
		return
	}

	ubrand, bErr := a.app.PatchBrand(bid, patch)
	if err != nil {
		respondError(w, bErr)
		return
	}

	respondJSON(w, http.StatusOK, ubrand)
}

func (a *API) deleteBrand(w http.ResponseWriter, r *http.Request) {
	bid, err := strconv.ParseInt(chi.URLParam(r, "brand_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("deleteBrand", model.ErrInternal, locale.GetUserLocalizer("en"), msgBrandURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	if err := a.app.DeleteBrand(bid); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}
