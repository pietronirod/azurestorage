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

func handleErrors(err error, message string) {
	if err != nil {
		log.Fatal(message)
	}
}

func main() {
	accountName := flag.String("accountname", "", "Azure Storage Account")
	accountKey := flag.String("accountkey", "", "Azure Account Key")
	containerName := flag.String("containername", "", "Container Name")
	fileName := flag.String("filename", "", "File name")
	flag.Parse()

	credential, err := azblob.NewSharedKeyCredential(*accountName, *accountKey)
	handleErrors(err, fmt.Sprintf("%s - account name: %s, account key: %s", sharedKey, *accountName, *accountKey))

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", *accountName, *containerName))

	containerURL := azblob.NewContainerURL(*URL, p)
	ctx := context.Background()

	file, err := os.Open(*fileName)
	handleErrors(err, fmt.Sprintf("%s: %s", localFile, *fileName))
	blobURL := containerURL.NewBlockBlobURL(filepath.Base(*fileName))

	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	handleErrors(err, fmt.Sprintf("%s: %s", upload, blobURL.BlobURL))

	//log.Print(blobURL.BlobURL)
	fmt.Println(blobURL.BlobURL)
}
