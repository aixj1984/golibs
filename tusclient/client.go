// Package tusclient provides a small client for tus upload servers.
package tusclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const tusResumableVersion = "1.0.0"

// TusClient is a client for tus upload operations.
type TusClient struct {
	Endpoint string
	// HTTPClient uses http.DefaultClient when nil.
	HTTPClient *http.Client
}

// UploadInfo describes the current state returned by a tus HEAD request.
type UploadInfo struct {
	Offset int64
	Length int64
}

// NewTusClient creates a tus client.
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

func parseUploadLength(resp *http.Response) (int64, error) {
	v := resp.Header.Get("Upload-Length")
	if v == "" {
		return 0, errors.New("missing Upload-Length header in response")
	}
	return strconv.ParseInt(v, 10, 64)
}

// FileURL returns the full upload URL for a fileID.
func (s *TusClient) FileURL(fileID string) string {
	if strings.HasPrefix(fileID, "http://") || strings.HasPrefix(fileID, "https://") {
		return fileID
	}
	return strings.TrimRight(s.Endpoint, "/") + "/" + strings.TrimLeft(fileID, "/")
}

func (s *TusClient) fileIDFromLocation(location string) string {
	return strings.TrimPrefix(location, strings.TrimRight(s.Endpoint, "/")+"/")
}

// CreateFile creates a tus upload resource.
func (s *TusClient) CreateFile(filename string, fileSize int64) (fileID string, err error) {
	return s.CreateFileContext(context.Background(), filename, fileSize)
}

// CreateFileContext creates a tus upload resource.
func (s *TusClient) CreateFileContext(ctx context.Context, filename string, fileSize int64) (fileID string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.Endpoint, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Upload-Length", strconv.FormatInt(fileSize, 10))
	req.Header.Set("Tus-Resumable", tusResumableVersion)
	req.Header.Set("Upload-Metadata", "filename "+base64.StdEncoding.EncodeToString([]byte(filename)))

	response, err := s.client().Do(req)
	if err != nil {
		return "", err
	}
	defer drainAndClose(response)

	if response.StatusCode == http.StatusCreated {
		return s.fileIDFromLocation(response.Header.Get("Location")), nil
	}
	return "", fmt.Errorf("CreateFileInServer failed: %s", response.Status)
}

// HeadUpload returns upload offset and length from a tus HEAD request.
func (s *TusClient) HeadUpload(ctx context.Context, fileID string) (*UploadInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, s.FileURL(fileID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Tus-Resumable", tusResumableVersion)

	response, err := s.client().Do(req)
	if err != nil {
		return nil, err
	}
	defer drainAndClose(response)

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("HEAD failed: %s", response.Status)
	}

	offset, err := parseUploadOffset(response)
	if err != nil {
		return nil, err
	}
	length, err := parseUploadLength(response)
	if err != nil {
		return nil, err
	}
	return &UploadInfo{Offset: offset, Length: length}, nil
}

// GetUploadOffset returns the uploaded offset for a tus resource.
func (s *TusClient) GetUploadOffset(ctx context.Context, fileID string) (offset int64, err error) {
	info, err := s.HeadUpload(ctx, fileID)
	if err != nil {
		return 0, err
	}
	return info.Offset, nil
}

// GetUploadPace returns the uploaded offset for a tus resource.
//
// Deprecated: use HeadUpload or GetUploadOffset.
func (s *TusClient) GetUploadPace(fileID string) (offset int64, err error) {
	return s.GetUploadOffset(context.Background(), fileID)
}

// PatchDataBlock uploads one byte slice and returns the server-confirmed Upload-Offset.
func (s *TusClient) PatchDataBlock(fileID string, offset int64, dataBytes []byte) (newOffset int64, err error) {
	return s.PatchDataBlockContext(context.Background(), fileID, offset, dataBytes)
}

// PatchDataBlockContext uploads one byte slice and returns the server-confirmed Upload-Offset.
func (s *TusClient) PatchDataBlockContext(ctx context.Context, fileID string, offset int64, dataBytes []byte) (newOffset int64, err error) {
	return s.PatchDataBlockReader(ctx, fileID, offset, bytes.NewReader(dataBytes))
}

// PatchDataBlockReader uploads one data block from r and returns the server-confirmed Upload-Offset.
func (s *TusClient) PatchDataBlockReader(ctx context.Context, fileID string, offset int64, r io.Reader) (newOffset int64, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, s.FileURL(fileID), r)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/offset+octet-stream")
	req.Header.Set("Upload-Offset", strconv.FormatInt(offset, 10))
	req.Header.Set("Tus-Resumable", tusResumableVersion)

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

// DeleteFile deletes a tus resource.
func (s *TusClient) DeleteFile(fileID string) error {
	return s.DeleteFileContext(context.Background(), fileID)
}

// DeleteFileContext deletes a tus resource.
func (s *TusClient) DeleteFileContext(ctx context.Context, fileID string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, s.FileURL(fileID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Tus-Resumable", tusResumableVersion)
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

// GetOssKey extracts a fileID from a full upload URL.
func (s *TusClient) GetOssKey(fileURL string) string {
	return s.fileIDFromLocation(fileURL)
}

// WriteFile writes r to an existing tus resource and returns the final Upload-Offset.
func (s *TusClient) WriteFile(r io.ReadSeeker, fileID string) (finalOffset int64, err error) {
	return s.WriteFileContext(context.Background(), r, fileID)
}

// WriteFileContext writes r to an existing tus resource and returns the final Upload-Offset.
func (s *TusClient) WriteFileContext(ctx context.Context, r io.ReadSeeker, fileID string) (finalOffset int64, err error) {
	offset, err := s.GetUploadOffset(ctx, fileID)
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

		newOffset, err := s.PatchDataBlockReader(ctx, fileID, offset, bytes.NewReader(buff[:n]))
		if err != nil {
			return 0, err
		}
		offset = newOffset
	}
	return offset, nil
}

// GetFileSize returns Upload-Length from a tus HEAD response.
func (s *TusClient) GetFileSize(fileID string) (length int64, err error) {
	return s.GetFileSizeContext(context.Background(), fileID)
}

// GetFileSizeContext returns Upload-Length from a tus HEAD response.
func (s *TusClient) GetFileSizeContext(ctx context.Context, fileID string) (length int64, err error) {
	info, err := s.HeadUpload(ctx, fileID)
	if err != nil {
		return 0, err
	}
	return info.Length, nil
}

// DownloadFile downloads a tus resource using HTTP Range requests.
func (s *TusClient) DownloadFile(fileID string, w io.Writer) error {
	return s.DownloadFileContext(context.Background(), fileID, w)
}

// DownloadFileContext downloads a tus resource using HTTP Range requests.
func (s *TusClient) DownloadFileContext(ctx context.Context, fileID string, w io.Writer) error {
	fileSize, err := s.GetFileSizeContext(ctx, fileID)
	if err != nil {
		return err
	}
	if fileSize == 0 {
		return nil
	}

	var offset int64
	for offset < fileSize {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.FileURL(fileID), nil)
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
