package hammer

import (
	"fmt"
	"io"
	"os"

	"github.com/curtisnewbie/gocommon/client"
	"github.com/curtisnewbie/gocommon/common"
)

type GenFileTempTokenReq struct {
	Filekeys    []string `json:"fileKeys"`
	ExpireInMin int      `json:"expireInMin"`
}

type GenFileTempTokenResp struct {
	common.Resp
	Data map[string]string `json:"data"`
}

type FstoreFile struct {
	Id         int64         `json:"id"`
	FileId     string        `json:"fileId"`
	Name       string        `json:"name"`
	Status     string        `json:"status"`
	Size       int64         `json:"size"`
	Md5        string        `json:"md5"`
	UplTime    common.ETime  `json:"uplTime"`
	LogDelTime *common.ETime `json:"logDelTime"`
	PhyDelTime *common.ETime `json:"phyDelTime"`
}

func GetFstoreTmpToken(rail common.Rail, fileId string) (string /* tmpToken */, error) {
	r := client.NewDynTClient(rail, "/file/key", "fstore").
		EnableTracing().
		AddQueryParams("fileId", fileId).
		Get()
	if r.Err != nil {
		return "", r.Err
	}
	defer r.Close()

	var res common.GnResp[string]
	if e := r.ReadJson(&res); e != nil {
		return "", e
	}

	if res.Error {
		return "", res.Err()
	}
	return res.Data, nil
}

func DownloadFstoreFile(rail common.Rail, tmpToken string, absPath string) error {
	r := client.NewDynTClient(rail, "/file/raw", "fstore").
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

func UploadFstoreFile(rail common.Rail, filename string, file string) (string /* uploadFileId */, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("failed to open file, %v", err)
	}
	defer f.Close()

	r := client.NewDynTClient(rail, "/file", "fstore").
		EnableTracing().
		AddHeaders(map[string]string{"filename": filename}).
		Put(f)
	if r.Err != nil {
		return "", r.Err
	}
	defer r.Close()

	var res common.GnResp[string]
	if e := r.ReadJson(&res); e != nil {
		return "", e
	}

	if res.Error {
		return "", res.Err()
	}
	return res.Data, nil
}

func FetchFstoreFileInfo(rail common.Rail, fileId string, uploadFileId string) (FstoreFile, error) {
	r := client.NewDynTClient(rail, "/file/info", "fstore").
		EnableTracing().
		AddQueryParams("fileId", fileId).
		AddQueryParams("uploadFileId", uploadFileId).
		Get()
	if r.Err != nil {
		return FstoreFile{}, r.Err
	}
	defer r.Close()

	var res common.GnResp[FstoreFile]
	if e := r.ReadJson(&res); e != nil {
		return FstoreFile{}, e
	}
	return res.Data, res.Err()
}
