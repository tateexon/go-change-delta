package utils

import (
	"bytes"
	"io"
	"os"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func Ptr[T any](value T) *T {
	return &value
}

func StringToBytesBuffer(s string) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString(s)
	return buffer
}

func BytesToBytesBuffer(b []byte) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.Write(b)
	return buffer
}

func FileToBytes(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader, err := EnsureConsistentEncoding(file)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func FileToBytesBuffer(filepath string) (bytes.Buffer, error) {
	b, err := FileToBytes(filepath)
	if err != nil {
		var bb bytes.Buffer
		return bb, err
	}
	return BytesToBytesBuffer(b), nil
}

func FileToString(filepath string) (string, error) {
	content, err := FileToBytes(filepath)
	return string(content), err
}

func EnsureConsistentEncoding(file *os.File) (io.Reader, error) {
	bom := make([]byte, 2)
	_, err := file.Read(bom)
	if err != nil {
		return nil, err
	}

	// Reset the file pointer to the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	var reader io.Reader
	if bom[0] == 0xFF && bom[1] == 0xFE {
		// UTF-16 LE
		reader = transform.NewReader(file, unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder())
	} else if bom[0] == 0xFE && bom[1] == 0xFF {
		// UTF-16 BE
		reader = transform.NewReader(file, unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder())
	} else {
		// Assume UTF-8 or no BOM
		reader = file
	}

	return reader, nil
}
