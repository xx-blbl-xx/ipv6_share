package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	xlog "ipv6_share/log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

type shareFileInfo struct {
	*os.File
	Name     string `json:"name"`
	Hash     string `json:"hash"`
	Size     int64  `json:"size"`
	FileType string `json:"file_type"`
}

func newFileInfo(file *os.File, fileInfo os.FileInfo) (*shareFileInfo, error) {
	if file == nil {
		return nil, errors.New("no file")
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		xlog.Error("copy err", err)
	}

	fileType := mime.TypeByExtension(filepath.Ext(fileInfo.Name()))
	if fileType == "" {
		// read a chunk to decide between utf-8 text and binary
		var buf [512]byte
		n, err := io.ReadFull(file, buf[:])
		if err != nil {
			xlog.Error("ReadFull err", err)
			return nil, err
		}
		fileType = http.DetectContentType(buf[:n])
		_, err = file.Seek(0, io.SeekStart) // rewind to output whole file
		if err != nil {
			xlog.Error("seeker can't seek", err)
			return nil, err
		}
	}

	return &shareFileInfo{
		File:     file,
		Name:     fileInfo.Name(),
		Hash:     base64.StdEncoding.EncodeToString(hash.Sum(nil)),
		Size:     fileInfo.Size(),
		FileType: fileType,
	}, nil
}
