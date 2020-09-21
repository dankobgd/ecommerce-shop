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
	msgReviewFromJSON         = &i18n.Message{ID: "model.review.from_json.app_error", Other: "could not decode review json"}
	msgReviewCreateErr        = &i18n.Message{ID: "api.review.create_review.app_error", Other: "could not create review"}
	msgReviewsGetErr          = &i18n.Message{ID: "api.review.get_reviews.app_error", Other: "could not get reviews"}
	msgReviewGetErr           = &i18n.Message{ID: "api.review.get_review.app_error", Other: "could not get review"}
	msgReviewPatchErr         = &i18n.Message{ID: "api.review.patch_review.app_error", Other: "could not update review"}
	msgReviewDeleteerr        = &i18n.Message{ID: "api.review.delete_review.app_error", Other: "could not delete review"}
	msgReviewURLParamErr      = &i18n.Message{ID: "api.review.url.params.app_error", Other: "could not parse URL params"}
	msgReviewPatchFromJSONErr = &i18n.Message{ID: "api.review.patch_review.app_error", Other: "could not decode review patch data"}
)

// InitReviews inits the review routes
func InitReviews(a *API) {
	a.Routes.Reviews.Post("/", a.AdminSessionRequired(a.createReview))
	a.Routes.Reviews.Get("/", a.getReviews)
	a.Routes.Review.Get("/", a.getReview)
	a.Routes.Review.Patch("/", a.AdminSessionRequired(a.patchReview))
	a.Routes.Review.Delete("/", a.AdminSessionRequired(a.deleteReview))
}

func (a *API) createReview(w http.ResponseWriter, r *http.Request) {
	rev, e := model.ReviewFromJSON(r.Body)
	if e != nil {
		respondError(w, model.NewAppErr("createReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgReviewFromJSON, http.StatusInternalServerError, nil))
		return
	}

	review, err := a.app.CreateReview(rev)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, review)
}

func (a *API) getReview(w http.ResponseWriter, r *http.Request) {
	rid, e := strconv.ParseInt(chi.URLParam(r, "review_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgReviewURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	rev, err := a.app.GetReview(rid)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, rev)
}

func (a *API) getReviews(w http.ResponseWriter, r *http.Request) {
	pages := pagination.NewFromRequest(r)
	reviews, err := a.app.GetReviews(pages.Limit(), pages.Offset())
	if err != nil {
		respondError(w, err)
		return
	}

	totalCount := -1
	if len(reviews) > 0 {
		totalCount = reviews[0].TotalCount
	}
	pages.SetData(reviews, totalCount)

	respondJSON(w, http.StatusOK, pages)
}

func (a *API) patchReview(w http.ResponseWriter, r *http.Request) {
	rid, err := strconv.ParseInt(chi.URLParam(r, "review_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("patchReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgReviewURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	patch, err := model.ReviewPatchFromJSON(r.Body)
	if err != nil {
		respondError(w, model.NewAppErr("patchReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgReviewPatchFromJSONErr, http.StatusInternalServerError, nil))
		return
	}

	urev, rErr := a.app.PatchReview(rid, patch)
	if err != nil {
		respondError(w, rErr)
		return
	}

	respondJSON(w, http.StatusOK, urev)
}

func (a *API) deleteReview(w http.ResponseWriter, r *http.Request) {
	rid, err := strconv.ParseInt(chi.URLParam(r, "review_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("deleteReview", model.ErrInternal, locale.GetUserLocalizer("en"), msgReviewURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	if err := a.app.DeleteReview(rid); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}
