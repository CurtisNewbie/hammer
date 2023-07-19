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
	bus.DeclareEventBus(comprImgNotifyBus)
	bus.SubscribeEventBus(comprImgProcBus, 1, ListenCompressImageEvent)
}

func ListenCompressImageEvent(evt CompressImageEvent) error {
	c := common.EmptyExecContext()
	c.Log.Infof("Received CompressImageEvent: %+v", evt)

	// generate temp token for downloading file from mini-fstore
	tkn, e := GetFstoreTmpToken(c, evt.FileId)
	if e != nil {
		c.Log.Errorf("Failed to GetFstoreTmpToken, %v", e)
		return nil
	}
	c.Log.Infof("tkn: %v", tkn)

	// temp path for the downloaded file
	downloaded := "/tmp/" + common.RandNum(20)

	// download the file from mini-fstore
	if e := DownloadFstoreFile(c, tkn, downloaded); e != nil {
		c.Log.Errorf("Failed to DownloadFstoreFile, %v", e)
		return nil
	}
	c.Log.Infof("File downloaded to %v", downloaded)
	defer os.Remove(downloaded)

	// compress the image
	compressed := downloaded + "_compressed"
	if e := CompressImage(downloaded, compressed); e != nil {
		c.Log.Errorf("Failed to compress image, %v", e)
		return nil // if the compression failed, there is no need to retry
	}
	defer os.Remove(compressed)
	c.Log.Infof("Image %v compressed to %v", evt.FileKey, compressed)

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
		bus.SendToEventBus(c, outboundEvent, comprImgNotifyBus)
	}
	return nil
}
