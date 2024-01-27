package hammer

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/curtisnewbie/hammer/api"
	fstore "github.com/curtisnewbie/mini-fstore/client"
	"github.com/curtisnewbie/miso/miso"
)

var (
	compressedImageFileIdCache = miso.NewTTLCache[string](15*time.Minute, 200)
)

func BootstrapServer(args []string) {
	miso.PreServerBootstrap(func(rail miso.Rail) error {
		miso.NewEventBus(api.CompressImageTriggerEventBus)
		miso.SubEventBus(api.CompressImageTriggerEventBus, 2, ListenCompressImageEvent)
		return nil
	})
	miso.BootstrapServer(os.Args)
}

func ListenCompressImageEvent(rail miso.Rail, evt api.ImageCompressTriggerEvent) error {
	rail.Infof("Received CompressImageEvent: %+v", evt)

	// wrap the compression logic with cache
	// it's possible that we compressed the image but somehow failed to reply to rabbitmq
	// using a tiny time-based cache can avoid compressing the same images again and again
	var err error
	generatedFileId, ok := compressedImageFileIdCache.Get(evt.Identifier, func() (string, bool) {
		generatedFileId, er := CompressImage(rail, evt)
		err = er
		return generatedFileId, er == nil && generatedFileId != ""
	})
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// reply to the specified event bus
	err = miso.PubEventBus(rail,
		api.ImageCompressReplyEvent{Identifier: evt.Identifier, FileId: generatedFileId},
		evt.ReplyTo)
	if err == nil {
		compressedImageFileIdCache.Del(evt.Identifier)
	}
	return err
}

func CompressImage(rail miso.Rail, evt api.ImageCompressTriggerEvent) (string, error) {
	rail.Infof("Received CompressImageEvent: %+v", evt)

	originFile, err := fstore.FetchFileInfo(rail, fstore.FetchFileInfoReq{FileId: evt.FileId})
	if err != nil {
		if errors.Is(err, fstore.ErrFileDeleted) || errors.Is(err, fstore.ErrFileNotFound) {
			rail.Warnf("File %v is not found or deleted, %v", evt.FileId, evt.Identifier)
			return "", nil
		}
		return "", fmt.Errorf("failed to fetch fstore file info: %v, %v", evt.FileId, err)
	}

	// generate temp token for downloading file from mini-fstore
	originFileToken, e := fstore.GenTempFileKey(rail, evt.FileId, "")
	if e != nil {
		rail.Errorf("Failed to GetFstoreTmpToken, %v", e)
		return "", nil
	}
	rail.Infof("tkn: %v", originFileToken)

	// download the origin file from mini-fstore
	downloadPath := "/tmp/" + miso.RandNum(20)
	downloadFile, err := os.Create(downloadPath)
	if err != nil {
		return "", err
	}
	defer downloadFile.Close()

	if e := fstore.DownloadFile(rail, originFileToken, downloadFile); e != nil {
		rail.Errorf("Failed to DownloadFstoreFile, %v", e)
		return "", nil
	}
	rail.Infof("File downloaded to %v", downloadPath)
	defer os.Remove(downloadPath)

	// compress the origin image, if the compression failed, we just give up
	compressPath := downloadPath + "_compressed"
	if e := GiftCompressImage(rail, downloadPath, compressPath); e != nil {
		rail.Errorf("Failed to compress image, giving up, %v", e)
		return "", nil // don't retry
	}
	defer os.Remove(compressPath)
	rail.Infof("Image %v compressed to %v", evt.Identifier, compressPath)

	// upload the compressed image to mini-fstore
	compressedFile, err := os.Open(compressPath)
	if err != nil {
		return "", fmt.Errorf("failed to open compressed file, file: %v, %v", compressPath, err)
	}
	defer compressedFile.Close()

	uploadFileId, e := fstore.UploadFile(rail, originFile.Name+"_thumbnail", compressedFile)
	if e != nil {
		rail.Errorf("Failed to UploadFstoreFile, %v", e)
		return "", fmt.Errorf("failed to upload fstore file, %v", e)
	}

	// exchange the uploadFileId with the real fileId
	thumbnailFile, e := fstore.FetchFileInfo(rail, fstore.FetchFileInfoReq{UploadFileId: uploadFileId})
	if e != nil {
		rail.Errorf("Failed to FetchFstoreFileInfo, %v", e)
		return "", fmt.Errorf("failed to fetch fstore file info, %v", e)
	}

	return thumbnailFile.FileId, nil
}
