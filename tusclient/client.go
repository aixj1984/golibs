// Package tusclient 是TUS的客户端操作包
package tusclient

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// TusClient Tus对接客户端类
type TusClient struct {
	Endpoint string
}

// NewTusClient 实例化一个tus客户端对象
func NewTusClient(endpoint string) *TusClient {
	return &TusClient{endpoint}
}

// CreateFile 创建文件
func (s *TusClient) CreateFile(filename string, fileSize int64) (fileID string, err error) {
	req, err := http.NewRequest(http.MethodPost, s.Endpoint, nil)
	if err != nil {
		return
	}

	// fmt.Println("FileName:", filename)
	// fmt.Println("FileSize:", fileSize)

	req.Header.Add("Upload-Length", strconv.FormatInt(fileSize, 10))
	req.Header.Add("Tus-Resumable", "1.0.0")
	req.Header.Add("Upload-Metadata", "filename "+base64.StdEncoding.EncodeToString([]byte(filename)))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		req = nil
		return
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			fmt.Println("response.Body.Close error ", err.Error())
		}
	}()

	// fmt.Println("response.StatusCode: ", response.StatusCode)
	// fmt.Println("http.StatusCreated: ", http.StatusCreated)

	if response.StatusCode == http.StatusCreated {
		fileID = response.Header.Get("Location")
	} else {
		err = errors.New("CreateFileInServer failed. " + response.Status)
	}

	return
}

// GetUploadPace 获取已经上传的偏移
func (s *TusClient) GetUploadPace(fileURL string) (offset int64, err error) {
	req, err := http.NewRequest(http.MethodHead, fileURL, nil)
	if err != nil {
		return
	}
	req.Header.Add("Tus-Resumable", "1.0.0")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			fmt.Println("response.Body.Close error ", err.Error())
		}
	}()
	offset, err = strconv.ParseInt(response.Header.Get("Upload-Offset"), 10, 64)
	return
}

// PatchDataBlock 分块上传
func (s *TusClient) PatchDataBlock(fileURL string, breakIndex int64, dataBytes []byte) (error, int64) {
	req, err := http.NewRequest(http.MethodPatch, fileURL, bytes.NewBuffer(dataBytes))
	if err != nil {
		return err, 0
	}
	req.Header.Add("Content-Type", "application/offset+octet-stream")
	req.Header.Add("Upload-Offset", strconv.FormatInt(breakIndex, 10))
	req.Header.Add("Tus-Resumable", "1.0.0")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err, 0
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			fmt.Println("response.Body.Close error ", err.Error())
		}
	}()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("response, err: %s", response.Status), 0
	}

	return nil, int64(len(dataBytes))
}

// DeleteFile 删除文件
func (s *TusClient) DeleteFile(fileURL string) error {
	req, err := http.NewRequest(http.MethodDelete, fileURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Tus-Resumable", "1.0.0")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			fmt.Println("response.Body.Close error ", err.Error())
		}
	}()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("response, err: %s", response.Status)
	}

	return nil
}

// GetOssKey 获取文件的key
func (s *TusClient) GetOssKey(fileURL string) string {
	return strings.Replace(fileURL, s.Endpoint, "", 1)
}

// WriteFile 写文件
func (s *TusClient) WriteFile(r io.ReadSeeker, fileID string) (error, int64) {
	offset, err := s.GetUploadPace(fileID)
	if err != nil {
		return err, 0
	}
	offset, err = r.Seek(offset, 0)
	if err != nil {
		return err, 0
	}

	buff := make([]byte, 32*1024)
	for {
		n, err := r.Read(buff)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err, 0
		}
		d := buff[:n]

		req, err := http.NewRequest(http.MethodPatch, fileID, bytes.NewBuffer(d))
		if err != nil {
			return err, 0
		}
		req.Header.Add("Content-Type", "application/offset+octet-stream")
		req.Header.Add("Upload-Offset", strconv.FormatInt(offset, 10))
		req.Header.Add("Tus-Resumable", "1.0.0")
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			return err, 0
		}

		defer func() {
			err = response.Body.Close()
			if err != nil {
				fmt.Println("response.Body.Close error ", err.Error())
			}
		}()

		if response.StatusCode != http.StatusNoContent {
			return fmt.Errorf("response, err: %s", response.Status), 0
		}
		offset += int64(n)
	}
	return nil, offset
}

// GetFileSize 获取文件的大小
func (s *TusClient) GetFileSize(fileID string) (offset int64, err error) {
	return s.GetUploadPace(fileID)
}

// DownloadFile 下载文件
func (s *TusClient) DownloadFile(fileID string, w io.Writer) error {
	var offset int64
	fileSize, err := s.GetFileSize(fileID)
	if err != nil {
		return err
	}
	for offset < fileSize {
		req, err := http.NewRequest(http.MethodGet, fileID, nil)
		if err != nil {
			return err
		}
		req.Header.Add("Offset", strconv.FormatInt(offset, 10))
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		d, err := io.ReadAll(response.Body)
		if err != nil {
			errClose := response.Body.Close()
			if errClose != nil {
				fmt.Printf("response body close error %s\n", errClose.Error())
			}
			return err
		}
		errClose := response.Body.Close()
		if errClose != nil {
			fmt.Printf("response body close error %s\n", errClose.Error())
		}
		_, err = w.Write(d)
		if err != nil {
			return err
		}
		offset += int64(len(d))
	}
	return nil
}
