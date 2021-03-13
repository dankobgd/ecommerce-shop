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
	msgTagFromJSON         = &i18n.Message{ID: "model.tag.from_json.app_error", Other: "could not decode tag json"}
	msgTagCreateErr        = &i18n.Message{ID: "api.tag.create_tag.app_error", Other: "could not create tag"}
	msgTagsGetErr          = &i18n.Message{ID: "api.tag.get_tags.app_error", Other: "could not get tags"}
	msgTagGetErr           = &i18n.Message{ID: "api.tag.get_tag.app_error", Other: "could not get tag"}
	msgTagPatchErr         = &i18n.Message{ID: "api.tag.patch_tag.app_error", Other: "could not update tag"}
	msgTagDeleteerr        = &i18n.Message{ID: "api.tag.delete_tag.app_error", Other: "could not delete tag"}
	msgTagMultipartErr     = &i18n.Message{ID: "api.tag.create_tag.multipart.app_error", Other: "could not decode tag multipart data"}
	msgTagURLParamErr      = &i18n.Message{ID: "api.tag.url.params.app_error", Other: "could not parse URL params"}
	msgTagPatchFromJSONErr = &i18n.Message{ID: "api.tag.patch_tag.app_error", Other: "could not decode tag patch data"}
)

// InitTags inits the tag routes
func InitTags(a *API) {
	a.Routes.Tags.Post("/", a.AdminSessionRequired(a.createTag))
	a.Routes.Tags.Delete("/bulk", a.AdminSessionRequired(a.deleteTags))
	a.Routes.Tags.Get("/count", a.getTagsCount)
	a.Routes.Tags.Get("/", a.getTags)
	a.Routes.Tag.Get("/", a.getTag)
	a.Routes.Tag.Patch("/", a.AdminSessionRequired(a.patchTag))
	a.Routes.Tag.Delete("/", a.AdminSessionRequired(a.deleteTag))
}

func (a *API) createTag(w http.ResponseWriter, r *http.Request) {
	t, e := model.TagFromJSON(r.Body)
	if e != nil {
		respondError(w, model.NewAppErr("createTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgTagFromJSON, http.StatusInternalServerError, nil))
		return
	}

	tag, err := a.app.CreateTag(t)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusCreated, tag)
}

func (a *API) getTag(w http.ResponseWriter, r *http.Request) {
	tid, e := strconv.ParseInt(chi.URLParam(r, "tag_id"), 10, 64)
	if e != nil {
		respondError(w, model.NewAppErr("getTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgTagURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	t, err := a.app.GetTag(tid)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, t)
}

func (a *API) getTagsCount(w http.ResponseWriter, r *http.Request) {
	c := a.app.GetTagsCount()
	respondJSON(w, http.StatusOK, map[string]int{"count": c})
}

func (a *API) getTags(w http.ResponseWriter, r *http.Request) {
	pages := pagination.NewFromRequest(r)
	tags, err := a.app.GetTags(pages.Limit(), pages.Offset())
	if err != nil {
		respondError(w, err)
		return
	}

	totalCount := -1
	if len(tags) > 0 {
		totalCount = tags[0].TotalCount
	}
	pages.SetData(tags, totalCount)

	respondJSON(w, http.StatusOK, pages)
}

func (a *API) patchTag(w http.ResponseWriter, r *http.Request) {
	tid, err := strconv.ParseInt(chi.URLParam(r, "tag_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("patchTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgTagURLParamErr, http.StatusInternalServerError, nil))
		return
	}

	patch, err := model.TagPatchFromJSON(r.Body)
	if err != nil {
		respondError(w, model.NewAppErr("patchTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgTagPatchFromJSONErr, http.StatusInternalServerError, nil))
		return
	}

	utag, tErr := a.app.PatchTag(tid, patch)
	if err != nil {
		respondError(w, tErr)
		return
	}

	respondJSON(w, http.StatusOK, utag)
}

func (a *API) deleteTag(w http.ResponseWriter, r *http.Request) {
	tid, err := strconv.ParseInt(chi.URLParam(r, "tag_id"), 10, 64)
	if err != nil {
		respondError(w, model.NewAppErr("deleteTag", model.ErrInternal, locale.GetUserLocalizer("en"), msgTagURLParamErr, http.StatusInternalServerError, nil))
		return
	}
	if err := a.app.DeleteTag(tid); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}

func (a *API) deleteTags(w http.ResponseWriter, r *http.Request) {
	ids := model.IntSliceFromJSON(r.Body)

	if err := a.app.DeleteTags(ids); err != nil {
		respondError(w, err)
		return
	}

	respondOK(w)
}
