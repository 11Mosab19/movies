package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
)

type CloudinaryService struct {
	Client *cloudinary.Cloudinary
}

type UploadedImage struct {
	URL      string
	PublicID string
}

func NewCloudinaryService() (*CloudinaryService, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load environment configuration: %w", err)
	}

	cloudinaryURL := strings.TrimSpace(os.Getenv("CLOUDINARY_URL"))
	if cloudinaryURL == "" {
		return nil, errors.New("CLOUDINARY_URL is not configured")
	}

	client, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, fmt.Errorf("initialize Cloudinary client: %w", err)
	}

	client.Config.URL.Secure = true

	return &CloudinaryService{Client: client}, nil
}

func (CS *CloudinaryService) UploadImage(
	ctx context.Context,
	image io.Reader,
	folder string,
) (UploadedImage, error) {
	if image == nil {
		return UploadedImage{}, errors.New("image content is required")
	}

	result, err := CS.Client.Upload.Upload(ctx, image, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return UploadedImage{}, fmt.Errorf("upload image to Cloudinary: %w", err)
	}

	return UploadedImage{
		URL:      result.SecureURL,
		PublicID: result.PublicID,
	}, nil
}

func (CS *CloudinaryService) DeleteImage(ctx context.Context, publicID string) error {
	if strings.TrimSpace(publicID) == "" {
		return errors.New("Cloudinary public ID is required")
	}

	result, err := CS.Client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
	})
	if err != nil {
		return fmt.Errorf("delete image from Cloudinary: %w", err)
	}

	if result.Result != "ok" {
		return fmt.Errorf(
			"delete image from Cloudinary: unexpected result %q",
			result.Result,
		)
	}

	return nil
}
