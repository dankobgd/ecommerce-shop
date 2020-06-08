package apiv1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dankobgd/ecommerce-shop/model"
)

func respondJSON(w http.ResponseWriter, code int, obj interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("could not encode json response: %v - %v", obj, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(b)
	return err
}

func respondError(w http.ResponseWriter, appErr *model.AppErr) error {
	b, err := json.Marshal(appErr)
	if err != nil {
		return fmt.Errorf("could not encode json response: %v - %v", appErr, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)
	_, err = w.Write(b)

	log.Println(appErr)

	return err
}

func respondOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
