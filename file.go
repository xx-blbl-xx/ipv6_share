package main

import (
	"crypto/sha256"
	"errors"
	"io"
	xlog "ipv6_share/log"
	"os"
)

type shareFileInfo struct {
	*os.File
	Name     string `json:"name"`
	Hash     []byte `json:"hash"`
	Size     int64  `json:"size"`
}

func newFileInfo(file *os.File, fileInfo os.FileInfo) (*shareFileInfo, error) {
	if file == nil {
		return nil, errors.New("no file")
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		xlog.Error("copy err", err)
	}

	return &shareFileInfo{
		File:     file,
		Name:     fileInfo.Name(),
		Hash:     hash.Sum(nil),
		Size:     fileInfo.Size(),
	}, nil
}
