package steganography

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
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

const ENDCHUNKTYPE = "IEND" // finds EOF

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

// process each chunk sequence using a loop
func (mc *MetaChunk) ProcessImage(b *bytes.Reader, c *CmdLineOpts) {
	count := 1 // 0 is for header
	chunkType := ""
	for chunkType != ENDCHUNKTYPE {
		fmt.Println("---- Chunk # " + strconv.Itoa(count) + " ----")
		offset := chk.getOffset(b)
		fmt.Printf("Chunk Offset: %#02x\n", offset)
		chk.readChunk(b)
		chunkType = chk.chunkTypeToString()
		count++
	}
}

// gets current offset by using Seek() which returns 0 bytes from current position (represented by 1)
func (mc *MetaChunk) getOffset(b *bytes.Reader) {
	offset, _ := b.Seek(0, 1)
	mc.Offset = offset
}

// functions for reading chunks, broken down into parts of chunks
func (mc *MetaChunk) readChunk(b *bytes.Reader) {
	mc.readChunkSize(b)
	mc.readChunkType(b)
	mc.readChunkBytes(b, mc.Chk.Size)
	mc.readChunkCRC(b)
}
func (mc *MetaChunk) readChunkSize(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Size); err != nil {
		log.Fatal(err)
	}
}
func (mc *MetaChunk) readChunkType(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Type); err != nil {
		log.Fatal(err)
	}
}
func (mc *MetaChunk) readChunkBytes(b *bytes.Reader, cLen uint64) {
	mc.Chk.Data = make([]byte, cLen)
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Data); err != nil {
		log.Fatal(err)
	}
}
func (mc *MetaChunk) readChunkCRC(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.CRC); err != nil {
		log.Fatal(err)
	}
}

/* //CmdLineOpts represents the cli arguments
type CmdLineOpts struct {
	Input    string
	Output   string
	Meta     bool
	Suppress bool
	Offset   string
	Inject   bool
	Payload  string
	Type     string
	Encode   bool
	Decode   bool
	Key      string
} */
