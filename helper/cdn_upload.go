package helper

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImageToCDN(fileHeader *multipart.FileHeader, name string) (string, error) {
	var urlCloudinary = fmt.Sprintf("cloudinary://%s:%s@%s",
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
		os.Getenv("CLOUDINARY_CLOUD_NAME"))

	CDNService, err := cloudinary.NewFromURL(urlCloudinary)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	response, err := CDNService.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   os.Getenv("CLOUDINARY_UPLOAD_FOLDER_NAME"),
		PublicID: name,
	})
	if err != nil {
		return "", err
	}

	log.Println(response.SecureURL)
	return response.SecureURL, nil
}
