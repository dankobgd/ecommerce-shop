package apiv1

import "net/http"

// AuthRequired middleware
func (a *API) AuthRequired(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := a.app.TokenValid(r)
		if err != nil {
			respondError(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
