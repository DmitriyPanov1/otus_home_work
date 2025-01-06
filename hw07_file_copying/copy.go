package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	defer fromFile.Close()

	fileInfo, err := fromFile.Stat()
	if err != nil {
		return errors.New("could not get file information")
	}

	fileSize := fileInfo.Size()

	if fileSize == 0 {
		return ErrUnsupportedFile
	}

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	newFile, err := os.Create(toPath)
	if err != nil {
		return errors.New("failed to create copy")
	}

	defer newFile.Close()

	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return errors.New("error copying")
	}

	copyBytes := fileSize - offset

	if limit != 0 && copyBytes > limit {
		copyBytes = limit
	}

	bar := pb.Start64(copyBytes)
	barReader := bar.NewProxyReader(fromFile)
	defer bar.Finish()

	_, err = io.CopyN(newFile, barReader, copyBytes)
	if err != nil {
		return errors.New("error copying")
	}

	return nil
}
