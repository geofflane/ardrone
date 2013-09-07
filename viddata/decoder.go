package viddata

import (
  "net"
  "bytes"
  "encoding/binary"
  "fmt"
  "os"
  "io"
)
type PaVE struct {
  Header PaVEHeader
  Payload []byte
}

type PaVEHeader struct {
  Signature             uint32 // 4 chars
  Version               uint8
  VideoCodec            uint8
  HeaderSize            uint16
  PayloadSize           uint32
  EncodedStreamWidth    uint16
  EncodedStreamHeight   uint16
  DisplayWidth          uint16
  DisplayHeight         uint16 // 20
  FrameNumber           uint32
  Timestamp             uint32
  TotalChunks           uint8
  ChunkIndex            uint8
  FrameType             uint8
  Control               uint8
  StreamBytePositionLw  uint32
  StreamBytePositionUw  uint32 // 40
  StreamId              uint16
  TotalSlices           uint8
  SliceIndex            uint8
  Header1Size           uint8
  Header2Size           uint8
  Reserved2             uint16 // 2 byte
  AdvertisedSize        uint32
  Reserved3_1           uint32 // 12 bytes
  Reserved3_2           uint32
  Reserved3_3           uint32 // 64
  Reserved3_4           uint32
  Reserved3_5           uint32
  Reserved3_6           uint32
}

// Decode blocks until the next video packets becomes available, which
// it then parses and returns.
func Decode(con net.Conn) (pave *PaVE, err error) {
	// readOrPanic() panics, while not expected, should not propagate to the
	// caller, so we return them like regular errors instead.
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

  buf := make([]byte, 76)
  _, err = con.Read(buf)
  if nil != err {
    panic(err)
  }

  pave = &PaVE{}
	reader := bytes.NewReader(buf)
	readOrPanic(reader, &pave.Header)

  fmt.Fprintf(os.Stderr, "Header: %v, Payload: %v\n", pave.Header.HeaderSize, pave.Header.PayloadSize)
  buf2 := make([]byte, pave.Header.PayloadSize)
  _, err = io.ReadFull(con, buf2)
  if nil != err {
    panic(err)
  }
  pave.Payload = buf2

	return
}

func readOrPanic(r *bytes.Reader, value interface{}) {
	if err := binary.Read(r, binary.LittleEndian, value); err != nil {
		panic(err)
	}
}

