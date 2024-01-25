package msgpraw

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgpWriter_WriteInt(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteInt(123)
	require.NoError(t, err)

	expected := make([]byte, 9)
	expected[0] = byte(Int64)
	binary.BigEndian.PutUint64(expected[1:], uint64(123))

	require.Equal(t, len(expected), len(w.Buff))
	for i, b := range w.Buff {
		assert.Equal(t, expected[i], b)
	}
}

func TestWriteReadInt(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteInt(123)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)

	assert.Equal(t, Int64, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 8, len(data))
	assert.Equal(t, uint64(123), binary.BigEndian.Uint64(data))
}

func TestWriteReadInt8(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteInt8(123)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)

	assert.Equal(t, Int8, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 1, len(data))
	assert.Equal(t, byte(123), data[0])
}

func TestWriteReadInt16(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteInt16(123)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)

	assert.Equal(t, Int16, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 2, len(data))
	assert.Equal(t, uint16(123), binary.BigEndian.Uint16(data))
}

func TestWriteReadInt32(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteInt32(123)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)

	assert.Equal(t, Int32, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 4, len(data))
	assert.Equal(t, uint32(123), binary.BigEndian.Uint32(data))
}

func TestWriteReadInt64(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteInt64(123)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)

	assert.Equal(t, Int64, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 8, len(data))
	assert.Equal(t, uint64(123), binary.BigEndian.Uint64(data))
}

func TestWriteReadUint(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteUint(123)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)

	assert.Equal(t, Uint64, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 8, len(data))
	assert.Equal(t, uint64(123), binary.BigEndian.Uint64(data))
}

func TestWriteReadUint8(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteUint8(123)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)

	assert.Equal(t, Uint8, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 1, len(data))
	assert.Equal(t, byte(123), data[0])
}

func TestWriteReadUint16(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteUint16(123)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)

	assert.Equal(t, Uint16, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 2, len(data))
	assert.Equal(t, uint16(123), binary.BigEndian.Uint16(data))
}

func TestWriteReadUint32(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteUint32(123)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)

	assert.Equal(t, Uint32, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 4, len(data))
	assert.Equal(t, uint32(123), binary.BigEndian.Uint32(data))
}

func TestWriteReadUint64(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteUint64(123)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)
	assert.Equal(t, Uint64, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 8, len(data))
	assert.Equal(t, uint64(123), binary.BigEndian.Uint64(data))
}

func TestWriteReadFloat32(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteFloat32(123.456)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)
	assert.Equal(t, Float32, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 4, len(data))
	assert.Equal(t, math.Float32bits(123.456), binary.BigEndian.Uint32(data))
}

func TestWriteReadFloat64(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteFloat64(123.456)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)
	assert.Equal(t, Float64, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 8, len(data))
	assert.Equal(t, math.Float64bits(123.456), binary.BigEndian.Uint64(data))
}

func TestWriteReadString(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteString("test")
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)
	assert.Equal(t, Str8, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 4, len(data))
	assert.Equal(t, "test", string(data))
}

func TestWriteReadBytes(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteBytes([]byte("test"))
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)
	assert.Equal(t, Bin8, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 4, len(data))
	assert.Equal(t, "test", string(data))
}

func TestWriteReadNil(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteNil()
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)
	assert.Equal(t, Nil, fType)
	assert.Equal(t, 0, num)
	assert.Nil(t, data)
}

func TestWriteReadBool(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteBool(true)
	require.NoError(t, err)

	err = w.WriteBool(false)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)
	assert.Equal(t, True, fType)
	assert.Equal(t, 0, num)
	assert.Nil(t, data)

	fType, num, data, err = r.Read()
	require.NoError(t, err)
	assert.Equal(t, False, fType)
	assert.Equal(t, 0, num)
	assert.Nil(t, data)
}

func TestWriteReadArray(t *testing.T) {
	w := &MsgpWriter{Buff: make([]byte, 0)}

	err := w.WriteArray(2)
	require.NoError(t, err)

	err = w.WriteInt(123)
	require.NoError(t, err)

	err = w.WriteInt(456)
	require.NoError(t, err)

	r := &MsgpReader{Buff: w.Buff}
	fType, num, data, err := r.Read()
	require.NoError(t, err)
	assert.Equal(t, Array16, fType)
	assert.Equal(t, 2, num)

	fType, num, data, err = r.Read()
	require.NoError(t, err)
	assert.Equal(t, Int64, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 8, len(data))
	assert.Equal(t, uint64(123), binary.BigEndian.Uint64(data))

	fType, num, data, err = r.Read()
	require.NoError(t, err)
	assert.Equal(t, Int64, fType)
	assert.Equal(t, 0, num)
	assert.Equal(t, 8, len(data))
	assert.Equal(t, uint64(456), binary.BigEndian.Uint64(data))
}

func BenchmarkMsgpWriter_WriteArray(b *testing.B) {
	buff := make([]byte, 255)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := &MsgpWriter{Buff: buff[:0]}

		_ = w.WriteMap(15)
		for j := 0; j < 15; j++ {
			_ = w.WriteInt8(int8(j))
			_ = w.WriteNil()
		}
	}
}
