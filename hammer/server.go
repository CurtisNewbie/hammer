package hammer

import (
	"fmt"
	"os"

	"github.com/curtisnewbie/gocommon/bus"
	"github.com/curtisnewbie/gocommon/common"
)

const (
	comprImgProcBus   = "hammer.image.compress.processing"
	comprImgNotifyBus = "hammer.image.compress.notification"
)

type CompressImageEvent struct {
	FileKey string // file key from vfm
	FileId  string // file id from mini-fstore
}

func PrepareServer(c common.ExecContext) {
	if e := bus.DeclareEventBus(comprImgProcBus); e != nil {
		c.Log.Fatalf("failed to declare event bus, %v", e)
	}

	if e := bus.DeclareEventBus(comprImgNotifyBus); e != nil {
		c.Log.Fatalf("failed to declare event bus, %v", e)
	}

	bus.SubscribeEventBus(comprImgProcBus, 2, ListenCompressImageEvent)
}

func ListenCompressImageEvent(evt CompressImageEvent) error {
	c := common.EmptyExecContext()
	c.Log.Infof("Received CompressImageEvent: %+v", evt)

	// generate temp token for downloading file from mini-fstore
	tkn, e := GetFstoreTmpToken(c, evt.FileId)
	if e != nil {
		c.Log.Errorf("Failed to GetFstoreTmpToken, %v", e)
		return fmt.Errorf("failed to generate fstore temp token, %v", e)
	}

	// temp path for the downloaded file
	tmpFile := "/tmp/" + common.RandNum(20)

	// download the file from mini-fstore
	if e := DownloadFstoreFile(c, tkn, tmpFile); e != nil {
		c.Log.Errorf("Failed to DownloadFstoreFile, %v", e)
		return fmt.Errorf("failed to download fstore file, %v", e)
	}
	c.Log.Infof("File downloaded to %v", tmpFile)
	defer os.Remove(tmpFile)

	// compress the image
	compressed := tmpFile + "_compressed"
	if e := CompressImage(tmpFile, compressed); e != nil {
		c.Log.Errorf("Failed to compress image, %v", e)
		return nil // if the compression failed, there is no need to retry
	}

	// upload the image back to mini-fstore
	uploadFileId, e := UploadFstoreFile(c, evt.FileKey+"_thumbnail", compressed)
	if e != nil {
		c.Log.Errorf("Failed to UploadFstoreFile, %v", e)
		return fmt.Errorf("failed to upload fstore file, %v", e)
	}

	// exchange the uploadFileId with the real fileId
	thumbnailFile, e := FetchFstoreFileInfo(c, "", uploadFileId)
	if e != nil {
		c.Log.Errorf("Failed to FetchFstoreFileInfo, %v", e)
		return fmt.Errorf("failed to fetch fstore file info, %v", e)
	}

	// record exists, dispatch the event to the oubound event bus (notify vfm about the fileId of the thumbnail)
	if thumbnailFile.Id > 0 {
		outboundEvent := CompressImageEvent{FileKey: evt.FileKey, FileId: thumbnailFile.FileId}
		bus.SendToEventBus(outboundEvent, comprImgNotifyBus)
	}
	return nil
}
