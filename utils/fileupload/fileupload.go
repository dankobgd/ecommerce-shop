package fileupload

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/dankobgd/ecommerce-shop/gocloudinary"
	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgCloudinaryDial           = &i18n.Message{ID: "cloudinary.dial.app_error", Other: "could not connect to cloudinary service"}
	msgCloudinaryUploadImage    = &i18n.Message{ID: "cloudinary.upload.image.app_error", Other: "could not upload image"}
	msgCloudinaryRecieveDetails = &i18n.Message{ID: "cloudinary.resource.details.app_error", Other: "could not get resource details"}
)

// UploadImageToCloudinary uploads the image and returns the preview url
func UploadImageToCloudinary(fileBytes []byte, fh *multipart.FileHeader, cloudEnvURI string) (string, *model.AppErr) {
	cloudinary, err := gocloudinary.Dial(cloudEnvURI)
	if err != nil {
		return "", model.NewAppErr("UploadImageToCloudinary", model.ErrInternal, locale.GetUserLocalizer("en"), msgCloudinaryDial, http.StatusInternalServerError, "")
	}

	imageName := strconv.FormatInt(time.Now().Unix(), 10) + "-" + fh.Filename
	publicID, err := cloudinary.UploadImage(imageName, bytes.NewBuffer(fileBytes), "ecommerce")
	if err != nil {
		fmt.Printf("upload err: %v\n", err)
		return "", model.NewAppErr("UploadImageToCloudinary", model.ErrInternal, locale.GetUserLocalizer("en"), msgCloudinaryUploadImage, http.StatusInternalServerError, "")
	}

	details, err := cloudinary.ResourceDetails(publicID)
	if err != nil {
		fmt.Printf("resource info err: %v\n", err)
		return "", model.NewAppErr("UploadImageToCloudinary", model.ErrInternal, locale.GetUserLocalizer("en"), msgCloudinaryRecieveDetails, http.StatusInternalServerError, nil)
	}

	return details.SecureURL, nil
}
