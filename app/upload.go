package app

import (
	"io"

	"github.com/dankobgd/ecommerce-shop/gocloudinary"
	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/fileupload"
)

func (a *App) uploadImageToCloudinary(data io.Reader, filename string) (*gocloudinary.ResourceDetails, *model.AppErr) {
	return fileupload.UploadImageToCloudinary(data, filename, a.Cfg().CloudinarySettings.EnvURI)
}

func (a *App) deleteImageFromCloudinary(publicID string) *model.AppErr {
	return fileupload.DeleteImageFromCloudinary(publicID, a.Cfg().CloudinarySettings.EnvURI)
}

// UploadImage uploads the image and returns the preview url
func (a *App) UploadImage(data io.Reader, filename string) (*gocloudinary.ResourceDetails, *model.AppErr) {
	return a.uploadImageToCloudinary(data, filename)
}

// DeleteImage deletes the image
func (a *App) DeleteImage(publicID string) *model.AppErr {
	return a.deleteImageFromCloudinary(publicID)
}
