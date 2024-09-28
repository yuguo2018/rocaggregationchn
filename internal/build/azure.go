package build

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

// AzureBlobstoreConfig is an authentication and configuration struct containing
// the data needed by the Azure SDK to interact with a specific container in the
// blobstore.
type AzureBlobstoreConfig struct {
	Account   string // Account name to authorize API requests with
	Token     string // Access token for the above account
	Container string // Blob container to upload files into
}

// AzureBlobstoreUpload uploads a local file to the Azure Blob Storage. Note, this
// method assumes a max file size of 64MB (Azure limitation). Larger files will
// need a multi API call approach implemented.
//
// See: https://msdn.microsoft.com/en-us/library/azure/dd179451.aspx#Anchor_3
func AzureBlobstoreUpload(path string, name string, config AzureBlobstoreConfig) error {
	if *DryRunFlag {
		fmt.Printf("would upload %q to %s/%s/%s\n", path, config.Account, config.Container, name)
		return nil
	}
	// Create an authenticated client against the Azure cloud
	credential, err := azblob.NewSharedKeyCredential(config.Account, config.Token)
	if err != nil {
		return err
	}
	a := fmt.Sprintf("https://%s.blob.core.windows.net/", config.Account)
	client, err := azblob.NewClientWithSharedKeyCredential(a, credential, nil)
	if err != nil {
		return err
	}
	// Stream the file to upload into the designated blobstore container
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = client.UploadFile(context.Background(), config.Container, name, in, nil)
	return err
}

// AzureBlobstoreList lists all the files contained within an azure blobstore.
func AzureBlobstoreList(config AzureBlobstoreConfig) ([]*container.BlobItem, error) {
	// Create an authenticated client against the Azure cloud
	credential, err := azblob.NewSharedKeyCredential(config.Account, config.Token)
	if err != nil {
		return nil, err
	}
	a := fmt.Sprintf("https://%s.blob.core.windows.net/", config.Account)
	client, err := azblob.NewClientWithSharedKeyCredential(a, credential, nil)
	if err != nil {
		return nil, err
	}
	pager := client.NewListBlobsFlatPager(config.Container, nil)

	var blobs []*container.BlobItem
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		blobs = append(blobs, page.Segment.BlobItems...)
	}
	return blobs, nil
}

// AzureBlobstoreDelete iterates over a list of files to delete and removes them
// from the blobstore.
func AzureBlobstoreDelete(config AzureBlobstoreConfig, blobs []*container.BlobItem) error {
	if *DryRunFlag {
		for _, blob := range blobs {
			fmt.Printf("would delete %s (%s) from %s/%s\n", *blob.Name, blob.Properties.LastModified, config.Account, config.Container)
		}
		return nil
	}
	// Create an authenticated client against the Azure cloud
	credential, err := azblob.NewSharedKeyCredential(config.Account, config.Token)
	if err != nil {
		return err
	}
	a := fmt.Sprintf("https://%s.blob.core.windows.net/", config.Account)
	client, err := azblob.NewClientWithSharedKeyCredential(a, credential, nil)
	if err != nil {
		return err
	}
	// Iterate over the blobs and delete them
	for _, blob := range blobs {
		if _, err := client.DeleteBlob(context.Background(), config.Container, *blob.Name, nil); err != nil {
			return err
		}
		fmt.Printf("deleted  %s (%s)\n", *blob.Name, blob.Properties.LastModified)
	}
	return nil
}
