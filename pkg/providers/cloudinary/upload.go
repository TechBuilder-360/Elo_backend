package cloudinary

import (
	"context"
	"errors"
	"path"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func (c *Cloud) upload(ctx context.Context, file string, fileName string, fileDestination FileDestination, entityType string) (string, error) {
	filePath := c.generateFilePath(fileDestination, entityType, fileName)

	upload, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: filePath,
	})
	if err != nil {
		return "", err
	}

	if upload.SecureURL == "" {
		return "", errors.New(upload.Error.Message)
	}

	return upload.SecureURL, nil
}

// UploadBusinessFile uploads a file to the business folder structure
func (c *Cloud) UploadBusinessFile(ctx context.Context, file string, fileName string, businessID string, uploadType string) (string, error) {
	destination := FileDestination{
		Type: uploadType,
		ID:   businessID,
	}
	return c.upload(ctx, file, fileName, destination, BusinessEntity)
}

// UploadUserFile uploads a file to the user folder structure
func (c *Cloud) UploadUserFile(ctx context.Context, file string, fileName string, userID string, uploadType string) (string, error) {
	destination := FileDestination{
		Type: uploadType,
		ID:   userID,
	}
	return c.upload(ctx, file, fileName, destination, UserEntity)
}

// UploadURLFile uploads a file to the folder structure
func (c *Cloud) UploadURLFile(ctx context.Context, fileURL string, fileName string, id string, uploadType string) (string, error) {
	destination := FileDestination{
		Type: uploadType,
		ID:   id,
	}
	return c.upload(ctx, fileURL, fileName, destination, UserEntity)
}

// generateFilePath creates a file path with proper directory structure
// Path format: elo/{entityType}/{entityID}/{uploadType}/{fileName}
func (c *Cloud) generateFilePath(destination FileDestination, entityType string, fileName string) string {
	return path.Join(
		"elo",
		entityType,
		destination.ID,
		destination.Type,
		fileName,
	)
}
