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

func Copy(fromPath, toPath string, offset, limit int64) error {
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

	bytesToCopy := size - offset
	if limit > 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	const bufSize = 4096
	buf := make([]byte, bufSize)
	var copied int64

	for copied < bytesToCopy {
		remaining := bytesToCopy - copied
		readSize := bufSize
		if remaining < int64(bufSize) {
			readSize = int(remaining)
		}

		n, err := src.Read(buf[:readSize])
		if n > 0 {
			_, writeErr := dst.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			copied += int64(n)

			percent := float64(copied) / float64(bytesToCopy) * 100
			fmt.Printf("\rProgress: %.2f%%", percent)
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