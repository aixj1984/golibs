package tusclient

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	END_POINT = "http://xxxx:1800/files/"
	FILE_PATH = "./data/"
	FILE_NAME = "xxxxxxx.zip"
)
var fileID string

var client *TusClient

func TestCreate(t *testing.T) {
	filePath := FILE_PATH + FILE_NAME
	// file info
	fd, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	client = NewTusClient(END_POINT)
	fileID, err = client.CreateFile(FILE_NAME, fd.Size())
	if err != nil {
		t.Error(err)
	}
	fmt.Println("new file id is : " + fileID)
}

func TestGetUploadPace(t *testing.T) {
	offset, err := client.GetUploadPace(fileID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, offset, int64(0))
}

func TestPatch(t *testing.T) {
	filePath := FILE_PATH + FILE_NAME

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 指定要读取的长度（这里设为 5 表示读取 5 个字节）
	length := 20
	buf := make([]byte, length)

	// 读取文件内容
	_, err = io.ReadAtLeast(file, buf, length)
	if err != nil {
		log.Fatal(err)
	}
	_, offset := client.PatchDataBlock(fileID, 0, buf)
	assert.Equal(t, offset, int64(length))

	// 设置读取的起始位置（这里设为 10 表示从文件的第 11 个字节开始读取）
	offset = int64(20)
	_, err = file.Seek(offset, 0)
	if err != nil {
		log.Fatal(err)
	}
	// 读取文件内容
	_, err = io.ReadAtLeast(file, buf, length)
	if err != nil {
		log.Fatal(err)
	}
	err, offset = client.PatchDataBlock(fileID, offset, buf)
	if err != nil {
		fmt.Println("client.PatchDataBlock err: " + err.Error())
	}
	assert.Equal(t, offset, int64(length))

	offset, err = client.GetUploadPace(fileID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, offset, int64(length+20))

	// file info
	fd, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}

	// 设置读取的起始位置（这里设为 10 表示从文件的第 11 个字节开始读取）
	offset = int64(40)
	_, err = file.Seek(offset, 0)
	if err != nil {
		log.Fatal(err)
	}

	buf = make([]byte, int(fd.Size())-40)
	// 读取文件内容
	_, err = io.ReadAtLeast(file, buf, int(fd.Size())-40)
	if err != nil {
		log.Fatal(err)
	}
	err, offset = client.PatchDataBlock(fileID, 40, buf)
	if err != nil {
		fmt.Println("client.PatchDataBlock err: " + err.Error())
	}
	assert.Equal(t, offset, int64(int64(fd.Size())-40))

	offset, err = client.GetUploadPace(fileID)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, offset, int64(fd.Size()))
}

func TestGetFileSize(t *testing.T) {
	offset, err := client.GetUploadPace(fileID)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, offset, int64(946801))
}

func TestGetKey(t *testing.T) {
	ossKey := client.GetOssKey(fileID)

	assert.Equal(t, ossKey, strings.Replace(fileID, END_POINT, "", 1))
}

func TestDownload(t *testing.T) {
	outFile, err := os.Create("./data/file.zip")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	err = client.DownloadFile(fileID, outFile)
	if err != nil {
		log.Fatal(err)
	}
	outFile.Close()
	os.Remove("./data/file.zip")
}

func TestDeleteFile(t *testing.T) {
	err := client.DeleteFile(fileID)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestCreate2(t *testing.T) {
	filePath := FILE_PATH + FILE_NAME
	// file info
	fd, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	client = NewTusClient(END_POINT)
	fileID, err = client.CreateFile(FILE_NAME, fd.Size())
	if err != nil {
		t.Error(err)
	}
	outFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println("new file id is : " + fileID)
	err, _ = client.WriteFile(outFile, fileID)
	if err != nil {
		panic(err)
	}
	client.DeleteFile(fileID)
}
