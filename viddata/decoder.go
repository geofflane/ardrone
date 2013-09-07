package viddata

import (
  "net"
  "bytes"
  "encoding/binary"
)

type PaVE struct {
  Header PaVEHeader
  Payload Payload
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
  DisplayHeight         uint16
  FrameNumber           uint32
  Timestamp             uint32
  TotalChunks           uint8
  ChunkIndex            uint8
  FrameType             uint8
  Control               uint8
  StreamBytePositionLw  uint32
  StreamBytePositionUw  uint32
  StreamId              uint16
  TotalSlices           uint8
  SliceIndex            uint8
  Header1Size           uint8
  Header2Size           uint8
  Reserved2             uint16 // 2 byte
  AdvertisedSize        uint32
  Reserved3_1           uint32 // 12 bytes
  Reserved3_2           uint32
  Reserved3_3           uint32
}

type Payload []byte

// Decode blocks until the next video packets becomes available, which
// it then parses and returns.
func Decode(con net.Conn) (paveData *PaVE, err error) {
	// readOrPanic() panics, while not expected, should not propagate to the
	// caller, so we return them like regular errors instead.
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

  buf := make([]byte, 68)
  _, err = con.Read(buf)
  if nil != err {
    panic("at the disco")
  }

  paveData = &PaVE{}
	reader := bytes.NewReader(buf)
	readOrPanic(reader, &paveData.Header)

  buf2 := make([]byte, paveData.Header.PayloadSize)
  _, err = con.Read(buf2)
  if nil != err {
    panic("at the disco 2")
  }
  paveData.Payload = Payload(buf2)
	return
}

func readOrPanic(r *bytes.Reader, value interface{}) {
	if err := binary.Read(r, binary.LittleEndian, value); err != nil {
		panic(err)
	}
}

