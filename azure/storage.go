package azure

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/joho/godotenv"
)

const containerName = "user-files"

type AzureStorageClient struct {
	serviceClient *service.Client
}

// Initializes the Azure Blob Storage client
func NewAzureStorageClient() (*AzureStorageClient, error) {
	_ = godotenv.Load()

	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	accountKey := os.Getenv("AZURE_STORAGE_KEY")

	if accountName == "" || accountKey == "" {
		return nil, fmt.Errorf("AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_KEY is not set")
	}

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create shared key credential: %w", err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	client, err := service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create service client: %w", err)
	}

	storageClient := &AzureStorageClient{serviceClient: client}

	// Ensure the shared container exists
	if err := storageClient.ensureContainerExists(); err != nil {
		return nil, err
	}

	return storageClient, nil
}

// Ensures that the shared container exists
func (c *AzureStorageClient) ensureContainerExists() error {
	containerClient := c.serviceClient.NewContainerClient(containerName)

	_, err := containerClient.Create(context.TODO(), nil)
	if err != nil && !strings.Contains(err.Error(), "ContainerAlreadyExists") {
		return fmt.Errorf("failed to create or access container '%s': %w", containerName, err)
	}

	return nil
}

// Uploads a file into the "folder" for a specific user
func (c *AzureStorageClient) UploadFile(userID, filename string, file io.Reader) error {
	if c.serviceClient == nil {
		return fmt.Errorf("Azure service client is not initialized")
	}

	blobPath := fmt.Sprintf("user-%s/%s", userID, filename)
	containerClient := c.serviceClient.NewContainerClient(containerName)
	blobClient := containerClient.NewBlockBlobClient(blobPath)

	log.Printf("[Azure] Uploading file '%s' to '%s'", filename, blobPath)

	_, err := blobClient.UploadStream(context.TODO(), file, nil)
	if err != nil {
		return fmt.Errorf("failed to upload blob '%s': %w", filename, err)
	}

	log.Printf("[Azure] File '%s' uploaded successfully", blobPath)
	return nil
}

// Lists all files for a specific user
func (c *AzureStorageClient) ListFiles(userID string) ([]string, error) {
	if c.serviceClient == nil {
		return nil, fmt.Errorf("Azure service client is not initialized")
	}

	prefix := fmt.Sprintf("user-%s/", userID)
	containerClient := c.serviceClient.NewContainerClient(containerName)
	pager := containerClient.NewListBlobsFlatPager(&azblob.ListBlobsFlatOptions{
		Prefix: &prefix,
	})

	var files []string

	log.Printf("[Azure] Listing files for user '%s'", userID)

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to list blobs: %w", err)
		}

		for _, blob := range resp.Segment.BlobItems {
			if blob.Name != nil {
				// Strip the user prefix to return only the filename
				files = append(files, strings.TrimPrefix(*blob.Name, prefix))
			}
		}
	}

	log.Printf("[Azure] Found %d files for user '%s'", len(files), userID)
	return files, nil
}

// Deletes a specific file for a user
func (c *AzureStorageClient) DeleteFile(userID, filename string) error {
	if c.serviceClient == nil {
		return fmt.Errorf("Azure service client is not initialized")
	}

	blobPath := fmt.Sprintf("user-%s/%s", userID, filename)
	containerClient := c.serviceClient.NewContainerClient(containerName)
	blobClient := containerClient.NewBlockBlobClient(blobPath)

	log.Printf("[Azure] Deleting blob '%s'", blobPath)

	_, err := blobClient.Delete(context.TODO(), nil)
	if err != nil {
		return fmt.Errorf("failed to delete blob '%s': %w", blobPath, err)
	}

	log.Printf("[Azure] Blob '%s' deleted successfully", blobPath)
	return nil
}
