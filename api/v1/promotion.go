package apiv1

import (
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/utils/pagination"
	"github.com/go-chi/chi"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgPromotionCreateErr        = &i18n.Message{ID: "api.promotion.create_promotion.app_error", Other: "could not create promotion"}
	msgPromotionsGetErr          = &i18n.Message{ID: "api.promotion.get_promotion.app_error", Other: "could not get promotions"}
	msgPromotionGetErr           = &i18n.Message{ID: "api.promotion.get_promotion.app_error", Other: "could not get promotion"}
	msgPromotionPatchErr         = &i18n.Message{ID: "api.promotion.patch_promotion.app_error", Other: "could not update promotion"}
	msgPromotionDeleteerr        = &i18n.Message{ID: "api.promotion.delete_promotion.app_error", Other: "could not delete promotion"}
	msgPromotionURLParamErr      = &i18n.Message{ID: "api.promotion.url.params.app_error", Other: "could not parse URL params"}
	msgPromotionPatchFromJSONErr = &i18n.Message{ID: "api.promotion.patch_product.app_error", Other: "could not decode promotion patch data"}
)

// InitPromotions inits the promotion routes
func InitPromotions(a *API) {
	a.Routes.Promotions.Post("/", a.AdminSessionRequired(a.createPromotion))
	a.Routes.Promotions.Get("/", a.AdminSessionRequired(a.getPromotions))
	a.Routes.Promotions.Delete("/bulk", a.AdminSessionRequired(a.deletePromotions))
	a.Routes.Promotion.Get("/", a.AdminSessionRequired(a.getPromotion))
	a.Routes.Promotion.Patch("/", a.AdminSessionRequired(a.patchPromotion))
	a.Routes.Promotion.Delete("/", a.AdminSessionRequired(a.deletePromotion))
	a.Routes.Promotion.Get("/valid", a.SessionRequired(a.getPromotionIsValid))
	a.Routes.Promotion.Get("/used", a.SessionRequired(a.getPromotionIsUsed))
	a.Routes.Promotion.Get("/status", a.SessionRequired(a.getPromotionStatus))
}

func (a *API) createPromotion(w http.ResponseWriter, r *http.Request) {
	p, e := model.PromotionFromJSON(r.Body)
	if e != nil {
		respondError(w, model.NewAppErr("createPromotion", model.ErrInternal, locale.GetUserLocalizer("en"), msgPromotionCreateErr, http.StatusInternalServerError, nil))
		return
	}

	promotion, err := a.app.CreatePromotion(p)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, promotion)
}

func (a *API) getPromotion(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "promo_code")
	promotion, err := a.app.GetPromotion(code)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, promotion)
}

func (a *API) getPromotions(w http.ResponseWriter, r *http.Request) {
	pages := pagination.NewFromRequest(r)
	promotions, err := a.app.GetPromotions(pages.Limit(), pages.Offset())
	if err != nil {
		respondError(w, err)
		return
	}

	totalCount := -1
	if len(promotions) > 0 {
		totalCount = promotions[0].TotalCount
	}
	pages.SetData(promotions, totalCount)

	respondJSON(w, http.StatusOK, pages)
}

func (a *API) patchPromotion(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "promo_code")
	patch, err := model.PromotionPatchFromJSON(r.Body)
	if err != nil {
		respondError(w, model.NewAppErr("patchPromotion", model.ErrInternal, locale.GetUserLocalizer("en"), msgPromotionPatchErr, http.StatusInternalServerError, nil))
		return
	}

	up, pErr := a.app.PatchPromotion(code, patch)
	if err != nil {
		respondError(w, pErr)
		return
	}

	respondJSON(w, http.StatusOK, up)
}

func (a *API) deletePromotion(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "promo_code")
	if err := a.app.DeletePromotion(code); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}

func (a *API) getPromotionIsValid(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "promo_code")

	if err := a.app.IsValidPromotion(code); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}
func (a *API) getPromotionIsUsed(w http.ResponseWriter, r *http.Request) {
	uid := a.app.GetUserIDFromContext(r.Context())
	code := chi.URLParam(r, "promo_code")

	if err := a.app.IsUsedPromotion(code, uid); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}
func (a *API) getPromotionStatus(w http.ResponseWriter, r *http.Request) {
	uid := a.app.GetUserIDFromContext(r.Context())
	code := chi.URLParam(r, "promo_code")

	if err := a.app.GetPromotionStatus(code, uid); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}

func (a *API) deletePromotions(w http.ResponseWriter, r *http.Request) {
	codes := model.StrSliceFromJSON(r.Body)

	if err := a.app.DeletePromotions(codes); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}
