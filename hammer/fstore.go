package hammer

import (
	"fmt"
	"io"
	"os"

	"github.com/curtisnewbie/miso/miso"
)

type GenFileTempTokenReq struct {
	Filekeys    []string `json:"fileKeys"`
	ExpireInMin int      `json:"expireInMin"`
}

type GenFileTempTokenResp struct {
	miso.Resp
	Data map[string]string `json:"data"`
}

type FstoreFile struct {
	Id         int64       `json:"id"`
	FileId     string      `json:"fileId"`
	Name       string      `json:"name"`
	Status     string      `json:"status"`
	Size       int64       `json:"size"`
	Md5        string      `json:"md5"`
	UplTime    miso.ETime  `json:"uplTime"`
	LogDelTime *miso.ETime `json:"logDelTime"`
	PhyDelTime *miso.ETime `json:"phyDelTime"`
}

func GetFstoreTmpToken(rail miso.Rail, fileId string) (string /* tmpToken */, error) {
	var res miso.GnResp[string]
	err := miso.NewDynTClient(rail, "/file/key", "fstore").
		EnableTracing().
		AddQueryParams("fileId", fileId).
		Get().
		Json(&res)
	if err != nil {
		return "", fmt.Errorf("failed to GetFstoreTmpToken, fileId: %v, %v", fileId, err)
	}
	return res.Res()
}

func DownloadFstoreFile(rail miso.Rail, tmpToken string, absPath string) error {
	r := miso.NewDynTClient(rail, "/file/raw", "fstore").
		EnableTracing().
		AddQueryParams("key", tmpToken).
		Get()
	if r.Err != nil {
		return r.Err
	}
	defer r.Close()

	out, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r.Resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func UploadFstoreFile(rail miso.Rail, filename string, file string) (string /* uploadFileId */, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("failed to open file, %v", err)
	}
	defer f.Close()

	var res miso.GnResp[string]
	err = miso.NewDynTClient(rail, "/file", "fstore").
		EnableTracing().
		AddHeaders(map[string]string{"filename": filename}).
		Put(f).
		Json(&res)
	if err != nil {
		return "", fmt.Errorf("failed to UploadFstoreFile, filename: %v, file: %v, %v", filename, file, err)
	}
	return res.Res()
}

func FetchFstoreFileInfo(rail miso.Rail, fileId string, uploadFileId string) (*FstoreFile, error) {
	var res miso.GnResp[FstoreFile]
	err := miso.NewDynTClient(rail, "/file/info", "fstore").
		EnableTracing().
		AddQueryParams("fileId", fileId).
		AddQueryParams("uploadFileId", uploadFileId).
		Get().
		Json(&res)
	if err != nil {
		return nil, fmt.Errorf("failed FetchFstoreFileInfo, fileId: %v, uploadFileId: %v, %v", fileId,
			uploadFileId, err)
	}

	if res.Error && res.ErrorCode == "FILE_NOT_FOUND" {
		return nil, nil
	}

	ff, err := res.Res()
	if err != nil {
		return nil, err
	}
	return &ff, nil
}
