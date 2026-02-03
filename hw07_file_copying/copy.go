package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, limit, offset int64) error {
	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()

	stat, err := src.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	size := stat.Size()

	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	if offset > 0 {
		_, err = src.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	maxAvailable := size - offset
	bytesToCopy := maxAvailable

	if limit > 0 && limit < maxAvailable {
		bytesToCopy = limit
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	const bufSize = 1024
	var copied int64

	for copied < bytesToCopy {
		remaining := bytesToCopy - copied
		readSize := bufSize
		if remaining < int64(bufSize) {
			readSize = int(remaining)
		}

		buf := make([]byte, readSize)

		n, err := src.Read(buf)

		if n > 0 {
			_, writeErr := dst.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			copied += int64(n)

			percent := int(float64(copied) / float64(bytesToCopy) * 100)
			fmt.Printf("\rProgress: %d%%", percent)
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	fmt.Println()
	return nil
}
