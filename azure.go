// Package main provides a function that copy local files to the Microsof Azure Storage
// i518034
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

const (
	sharedKey = "Error generating shared key credential"
	localFile = "Error opening local file"
	upload    = "Error during upload"
)

// handleErrors - In case of error finish the execution with an exception and a message
func HandleErrors(err error, message string) {
	if err != nil {
		log.Fatal(message)
	}
}

// main - The main code responsible to receive the command-line flags, and all steps to send the file to
// the Microsoft Azure Storage.
func main() {
	accountName := flag.String("accountname", "", "Azure Storage Account")
	accountKey := flag.String("accountkey", "", "Azure Account Key")
	containerName := flag.String("containername", "", "Container Name")
	fileName := flag.String("filename", "", "File name")
	flag.Parse()

	credential, err := azblob.NewSharedKeyCredential(*accountName, *accountKey)
	HandleErrors(err, fmt.Sprintf("%s - account name: %s, account key: %s", sharedKey, *accountName, *accountKey))

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", *accountName, *containerName))

	containerURL := azblob.NewContainerURL(*URL, p)
	ctx := context.Background()

	file, err := os.Open(*fileName)
	HandleErrors(err, fmt.Sprintf("%s: %s", localFile, *fileName))
	blobURL := containerURL.NewBlockBlobURL(filepath.Base(*fileName))

	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	HandleErrors(err, fmt.Sprintf("%s: %s", upload, blobURL.BlobURL))

	file.Close()
	//log.Print(blobURL.BlobURL)
	fmt.Println(blobURL.BlobURL)
}
