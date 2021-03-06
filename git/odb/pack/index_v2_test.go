package pack

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	V2IndexHeader = []byte{
		0xff, 0x74, 0x4f, 0x63,
		0x00, 0x00, 0x00, 0x02,
	}
	V2IndexFanout = make([]uint32, indexFanoutEntries)

	V2IndexNames = []byte{
		0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,
		0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1,

		0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2,
		0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2,

		0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3,
		0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3,
	}
	V2IndexSmallSha  = V2IndexNames[0:20]
	V2IndexMediumSha = V2IndexNames[20:40]
	V2IndexLargeSha  = V2IndexNames[40:60]

	V2IndexCRCs = []byte{
		0x0, 0x0, 0x0, 0x0,
		0x1, 0x1, 0x1, 0x1,
		0x2, 0x2, 0x2, 0x2,
	}

	V2IndexOffsets = []byte{
		0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x02,
		0x80, 0x00, 0x04, 0x5c,

		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
	}

	V2Index = &Index{
		fanout:  V2IndexFanout,
		version: new(V2),
	}
)

func TestIndexV2EntryExact(t *testing.T) {
	e, err := new(V2).Entry(V2Index, 1)

	assert.NoError(t, err)
	assert.EqualValues(t, 2, e.PackOffset)
}

func TestIndexV2EntryExtendedOffset(t *testing.T) {
	e, err := new(V2).Entry(V2Index, 2)

	assert.NoError(t, err)
	assert.EqualValues(t, 3, e.PackOffset)
}

func TestIndexVersionWidthV2(t *testing.T) {
	assert.EqualValues(t, 8, new(V2).Width())
}

func init() {
	V2IndexFanout[1] = 1
	V2IndexFanout[2] = 2
	V2IndexFanout[3] = 3

	for i := 3; i < len(V2IndexFanout); i++ {
		V2IndexFanout[i] = 3
	}

	fanout := make([]byte, indexFanoutWidth)
	for i, n := range V2IndexFanout {
		binary.BigEndian.PutUint32(fanout[i*indexFanoutEntryWidth:], n)
	}

	buf := make([]byte, 0, indexOffsetV2Start+3*(indexObjectEntryV2Width)+indexObjectLargeOffsetWidth)
	buf = append(buf, V2IndexHeader...)
	buf = append(buf, fanout...)
	buf = append(buf, V2IndexNames...)
	buf = append(buf, V2IndexCRCs...)
	buf = append(buf, V2IndexOffsets...)

	V2Index.r = bytes.NewReader(buf)
}
