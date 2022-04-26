package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

// first 8 bytes of any valid PNG will be the header
type Header struct {
	Header uint64
}

// remaining data is in chunks of set size
type Chunk struct {
	Size uint64
	Type uint64
	Data []byte
	CRC  uint64
}

//MetaChunk inherits a Chunk struct, contains chunk and offset
type MetaChunk struct {
	Chk    Chunk
	Offset int64
}

// opens and preprocesses file to be read
func PreProcessImage(dat *os.File) (*bytes.Reader, error) {
	stats, err := dat.Stat()
	if err != nil {
		return nil, err
	}

	var size = stats.Size()
	b := make([]byte, size)

	bufR := bufio.NewReader(dat)
	_, err = bufR.Read(b)
	bReader := bytes.NewReader(b)

	return bReader, err
}

// validates the file by checking the header
func (mc *MetaChunk) validate(b *bytes.Reader) {
	var header Header
	if err := binary.Read(b, binary.BigEndian, &header.Header); err != nil {
		log.Fatal(err)
	}

	bArr := make([]byte, 8)
	binary.BigEndian.PutUint64(bArr, header.Header)

	if string(bArr[1:4]) != "PNG" {
		log.Fatal("File is not in valid PNG format. :(")
	} else {
		fmt.Println("File is a valid PNG :)")
	}
}
