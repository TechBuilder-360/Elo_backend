package cloudinary

import (
	"context"
	"path"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func (c *cloud) Upload(ctx context.Context, file []byte, fileName string, fileDestination FileDestination, entityType string) (string, error) {
	filePath := c.generateFilePath(fileDestination, entityType, fileName)

	upload, err := c.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       filePath,
	})
	if err != nil {
		return "", err
	}

	return upload.SecureURL, nil
}

// UploadBusinessFile uploads a file to the business folder structure
func (c *cloud) UploadBusinessFile(ctx context.Context, file []byte, fileName string, businessID string, uploadType string) (string, error) {
	destination := FileDestination{
		Type: uploadType,
		ID:   businessID,
	}
	return c.Upload(ctx, file, fileName, destination, BusinessEntity)
}

// UploadUserFile uploads a file to the user folder structure
func (c *cloud) UploadUserFile(ctx context.Context, file []byte, fileName string, userID string, uploadType string) (string, error) {
	destination := FileDestination{
		Type: uploadType,
		ID:   userID,
	}
	return c.Upload(ctx, file, fileName, destination, UserEntity)
}

// generateFilePath creates a file path with proper directory structure
// Path format: elo/{entityType}/{entityID}/{uploadType}/{fileName}
func (c *cloud) generateFilePath(destination FileDestination, entityType string, fileName string) string {
	return path.Join(
		"elo",
		entityType,
		destination.ID,
		destination.Type,
		fileName,
	)
}
