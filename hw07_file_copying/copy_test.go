package main

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		out       string
		reference string
		limit     int64
		offset    int64
	}{
		{
			"offset 0 limit 0",
			"testdata/input.txt",
			"out.txt",
			"testdata/out_offset0_limit0.txt",
			0,
			0,
		},
		{
			"limit is greater than the copied file",
			"testdata/input.txt",
			"out2.txt",
			"testdata/out_offset0_limit0.txt",
			6617 + 10000,
			0,
		},
		{
			"offset 0 limit 10",
			"testdata/input.txt",
			"out3.txt",
			"testdata/out_offset0_limit10.txt",
			10,
			0,
		},
		{
			"offset 0 limit 1000",
			"testdata/input.txt",
			"out4.txt",
			"testdata/out_offset0_limit1000.txt",
			1000,
			0,
		},
		{
			"offset 0 limit 10000",
			"testdata/input.txt",
			"out5.txt",
			"testdata/out_offset0_limit10000.txt",
			10000,
			0,
		},
		{
			"offset 100 limit 1000",
			"testdata/input.txt",
			"out6.txt",
			"testdata/out_offset100_limit1000.txt",
			1000,
			100,
		},
		{
			"offset 6000 limit 1000",
			"testdata/input.txt",
			"out7.txt",
			"testdata/out_offset6000_limit1000.txt",
			1000,
			6000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.in, tt.out, tt.offset, tt.limit)
			require.NoError(t, err)

			defer os.Remove(tt.out)

			equalFiles(t, tt.reference, tt.out)
		})
	}

	t.Run("offset > fileSize", func(t *testing.T) {
		fromPath := "testdata/input.txt"

		fromFile, err := os.Open(fromPath)
		require.NoError(t, err)

		info, err := fromFile.Stat()
		require.NoError(t, err)

		err = Copy(fromPath, "out.txt", info.Size()+1000, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	})

	t.Run("file length unknown", func(t *testing.T) {
		fromPath := "/dev/urandom"

		err := Copy(fromPath, "out.txt", 0, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})
}

func equalFiles(t *testing.T, in, out string) {
	t.Helper()

	hash1, err := hashFileMD5(in)
	require.NoError(t, err)

	hash2, err := hashFileMD5(out)
	require.NoError(t, err)

	require.True(t, bytes.Equal(hash1, hash2))
}

func hashFileMD5(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}
