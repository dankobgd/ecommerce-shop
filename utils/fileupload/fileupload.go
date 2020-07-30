package fileupload

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/dankobgd/ecommerce-shop/gocloudinary"
	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const baseCloudinaryDir = "ecommerce"

var (
	msgCloudinaryDial           = &i18n.Message{ID: "cloudinary.dial.app_error", Other: "could not connect to cloudinary service"}
	msgCloudinaryUploadImage    = &i18n.Message{ID: "cloudinary.upload.image.app_error", Other: "could not upload image"}
	msgCloudinaryRecieveDetails = &i18n.Message{ID: "cloudinary.resource.details.app_error", Other: "could not get resource details"}
)

// UploadImageToCloudinary uploads the image and returns the preview url
func UploadImageToCloudinary(data io.Reader, filename string, cloudEnvURI string) (*gocloudinary.ResourceDetails, *model.AppErr) {
	cloudinary, err := gocloudinary.Dial(cloudEnvURI)
	if err != nil {
		return nil, model.NewAppErr("UploadImageToCloudinary", model.ErrInternal, locale.GetUserLocalizer("en"), msgCloudinaryDial, http.StatusInternalServerError, "")
	}

	imageName := strconv.FormatInt(time.Now().Unix(), 10) + "-" + filename
	publicID, err := cloudinary.UploadImage(imageName, data, baseCloudinaryDir)
	if err != nil {
		return nil, model.NewAppErr("UploadImageToCloudinary", model.ErrInternal, locale.GetUserLocalizer("en"), msgCloudinaryUploadImage, http.StatusInternalServerError, "")
	}

	details, err := cloudinary.ResourceDetails(publicID)
	if err != nil {
		return nil, model.NewAppErr("UploadImageToCloudinary", model.ErrInternal, locale.GetUserLocalizer("en"), msgCloudinaryRecieveDetails, http.StatusInternalServerError, nil)
	}

	return details, nil
}

// DeleteImageFromCloudinary deletes the image from cloudinary
func DeleteImageFromCloudinary(publicID string, cloudEnvURI string) *model.AppErr {
	cloudinary, err := gocloudinary.Dial(cloudEnvURI)
	if err != nil {
		return model.NewAppErr("DeleteImageFromCloudinary", model.ErrInternal, locale.GetUserLocalizer("en"), msgCloudinaryDial, http.StatusInternalServerError, "")
	}

	if err := cloudinary.Delete(publicID, "", gocloudinary.ImageType); err != nil {
		return model.NewAppErr("DeleteImageFromCloudinary", model.ErrInternal, locale.GetUserLocalizer("en"), msgCloudinaryUploadImage, http.StatusInternalServerError, "")
	}

	return nil
}
