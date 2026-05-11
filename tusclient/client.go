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
	// HTTPClient 为 nil 时使用 http.DefaultClient。可设置超时等策略。
	HTTPClient *http.Client
}

// NewTusClient 实例化一个tus客户端对象
func NewTusClient(endpoint string) *TusClient {
	return &TusClient{Endpoint: endpoint}
}

func (s *TusClient) client() *http.Client {
	if s != nil && s.HTTPClient != nil {
		return s.HTTPClient
	}
	return http.DefaultClient
}

func drainAndClose(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}

func parseUploadOffset(resp *http.Response) (int64, error) {
	v := resp.Header.Get("Upload-Offset")
	if v == "" {
		return 0, errors.New("missing Upload-Offset header in response")
	}
	return strconv.ParseInt(v, 10, 64)
}

// CreateFile 创建文件
func (s *TusClient) CreateFile(filename string, fileSize int64) (fileID string, err error) {
	req, err := http.NewRequest(http.MethodPost, s.Endpoint, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Upload-Length", strconv.FormatInt(fileSize, 10))
	req.Header.Set("Tus-Resumable", "1.0.0")
	req.Header.Set("Upload-Metadata", "filename "+base64.StdEncoding.EncodeToString([]byte(filename)))

	response, err := s.client().Do(req)
	if err != nil {
		return "", err
	}
	defer drainAndClose(response)

	if response.StatusCode == http.StatusCreated {
		return response.Header.Get("Location"), nil
	}
	return "", fmt.Errorf("CreateFileInServer failed: %s", response.Status)
}

// GetUploadPace 获取已经上传的偏移
func (s *TusClient) GetUploadPace(fileURL string) (offset int64, err error) {
	req, err := http.NewRequest(http.MethodHead, fileURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Tus-Resumable", "1.0.0")

	response, err := s.client().Do(req)
	if err != nil {
		return 0, err
	}
	defer drainAndClose(response)

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		return 0, fmt.Errorf("HEAD failed: %s", response.Status)
	}

	return strconv.ParseInt(response.Header.Get("Upload-Offset"), 10, 64)
}

// PatchDataBlock 分块上传，返回服务端确认的 Upload-Offset。
func (s *TusClient) PatchDataBlock(fileURL string, breakIndex int64, dataBytes []byte) (newOffset int64, err error) {
	req, err := http.NewRequest(http.MethodPatch, fileURL, bytes.NewReader(dataBytes))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/offset+octet-stream")
	req.Header.Set("Upload-Offset", strconv.FormatInt(breakIndex, 10))
	req.Header.Set("Tus-Resumable", "1.0.0")

	response, err := s.client().Do(req)
	if err != nil {
		return 0, err
	}
	defer drainAndClose(response)

	if response.StatusCode != http.StatusNoContent {
		return 0, fmt.Errorf("PATCH failed: %s", response.Status)
	}
	return parseUploadOffset(response)
}

// DeleteFile 删除文件
func (s *TusClient) DeleteFile(fileURL string) error {
	req, err := http.NewRequest(http.MethodDelete, fileURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Tus-Resumable", "1.0.0")
	response, err := s.client().Do(req)
	if err != nil {
		return err
	}
	defer drainAndClose(response)

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("DELETE failed: %s", response.Status)
	}

	return nil
}

// GetOssKey 获取文件的key（相对 Endpoint 的路径前缀）
func (s *TusClient) GetOssKey(fileURL string) string {
	return strings.TrimPrefix(fileURL, s.Endpoint)
}

// WriteFile 写文件，返回服务端最终 Upload-Offset。
func (s *TusClient) WriteFile(r io.ReadSeeker, fileID string) (finalOffset int64, err error) {
	offset, err := s.GetUploadPace(fileID)
	if err != nil {
		return 0, err
	}
	_, err = r.Seek(offset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	buff := make([]byte, 32*1024)
	for {
		n, readErr := r.Read(buff)
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			return 0, readErr
		}
		d := buff[:n]

		req, err := http.NewRequest(http.MethodPatch, fileID, bytes.NewReader(d))
		if err != nil {
			return 0, err
		}
		req.Header.Set("Content-Type", "application/offset+octet-stream")
		req.Header.Set("Upload-Offset", strconv.FormatInt(offset, 10))
		req.Header.Set("Tus-Resumable", "1.0.0")

		response, err := s.client().Do(req)
		if err != nil {
			return 0, err
		}

		if response.StatusCode != http.StatusNoContent {
			drainAndClose(response)
			return 0, fmt.Errorf("PATCH failed: %s", response.Status)
		}

		newOff, perr := parseUploadOffset(response)
		drainAndClose(response)
		if perr != nil {
			return 0, perr
		}
		offset = newOff
	}
	return offset, nil
}

// GetFileSize 获取文件的大小
func (s *TusClient) GetFileSize(fileID string) (offset int64, err error) {
	return s.GetUploadPace(fileID)
}

// DownloadFile 使用 HTTP Range 下载（与 tusd ServeContent 行为一致）；若服务端忽略 Range 则按整包回退处理。
func (s *TusClient) DownloadFile(fileID string, w io.Writer) error {
	fileSize, err := s.GetFileSize(fileID)
	if err != nil {
		return err
	}
	if fileSize == 0 {
		return nil
	}

	var offset int64
	for offset < fileSize {
		req, err := http.NewRequest(http.MethodGet, fileID, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))

		resp, err := s.client().Do(req)
		if err != nil {
			return err
		}

		switch resp.StatusCode {
		case http.StatusPartialContent:
			writ, err := io.Copy(w, resp.Body)
			drainAndClose(resp)
			if err != nil {
				return err
			}
			offset += writ

		case http.StatusOK:
			if offset == 0 {
				writ, err := io.Copy(w, resp.Body)
				drainAndClose(resp)
				if err != nil {
					return err
				}
				offset += writ
			} else {
				if _, err := io.CopyN(io.Discard, resp.Body, offset); err != nil {
					drainAndClose(resp)
					return err
				}
				writ, err := io.Copy(w, resp.Body)
				drainAndClose(resp)
				if err != nil {
					return err
				}
				offset += writ
			}

		case http.StatusRequestedRangeNotSatisfiable:
			drainAndClose(resp)
			return fmt.Errorf("download range not satisfiable: %s", resp.Status)

		default:
			drainAndClose(resp)
			return fmt.Errorf("download failed: %s", resp.Status)
		}
	}
	return nil
}
