package msgpraw

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader_Read_PosFixInt(t *testing.T) {
	for _, b := range []byte{0x00, 0x01, 0x40, 0x7f} {
		r := &MsgpReader{Buff: []byte{b}}
		ty, n, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Type(b), ty)
		assert.Equal(t, 0, n)
		assert.Nil(t, data)
		assert.Equal(t, 1, r.Idx)
	}
}

func TestReader_Read_NegFixInt(t *testing.T) {
	for _, b := range []byte{0xe0, 0xf0, 0xff} {
		r := &MsgpReader{Buff: []byte{b}}
		ty, n, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Type(b), ty)
		assert.Equal(t, 0, n)
		assert.Nil(t, data)
		assert.Equal(t, 1, r.Idx)
	}
}

func TestReader_Read_NilBool(t *testing.T) {
	cases := []struct {
		name string
		b    byte
		want Type
	}{
		{"nil", 0xc0, Nil},
		{"true", 0xc3, True},
		{"false", 0xc2, False},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &MsgpReader{Buff: []byte{tc.b}}
			ty, n, data, err := r.Read()
			require.NoError(t, err)
			assert.Equal(t, tc.want, ty)
			assert.Equal(t, 0, n)
			assert.Nil(t, data)
		})
	}
}

func TestReader_Read_Uint(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		r := &MsgpReader{Buff: []byte{0xcc, 0xab}}
		ty, _, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Uint8, ty)
		assert.Equal(t, byte(0xab), data[0])
	})
	t.Run("uint16", func(t *testing.T) {
		r := &MsgpReader{Buff: []byte{0xcd, 0x12, 0x34}}
		ty, _, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Uint16, ty)
		assert.Equal(t, uint16(0x1234), binary.BigEndian.Uint16(data))
	})
	t.Run("uint32", func(t *testing.T) {
		r := &MsgpReader{Buff: []byte{0xce, 0xde, 0xad, 0xbe, 0xef}}
		ty, _, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Uint32, ty)
		assert.Equal(t, uint32(0xdeadbeef), binary.BigEndian.Uint32(data))
	})
	t.Run("uint64", func(t *testing.T) {
		buf := []byte{0xcf, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
		r := &MsgpReader{Buff: buf}
		ty, _, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Uint64, ty)
		assert.Equal(t, uint64(0x0102030405060708), binary.BigEndian.Uint64(data))
	})
}

func TestReader_Read_Float(t *testing.T) {
	t.Run("float32", func(t *testing.T) {
		buf := make([]byte, 5)
		buf[0] = byte(Float32)
		binary.BigEndian.PutUint32(buf[1:], math.Float32bits(1.5))
		r := &MsgpReader{Buff: buf}
		ty, _, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Float32, ty)
		assert.Equal(t, float32(1.5), math.Float32frombits(binary.BigEndian.Uint32(data)))
	})
	t.Run("float64", func(t *testing.T) {
		buf := make([]byte, 9)
		buf[0] = byte(Float64)
		binary.BigEndian.PutUint64(buf[1:], math.Float64bits(-3.14))
		r := &MsgpReader{Buff: buf}
		ty, _, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Float64, ty)
		assert.Equal(t, -3.14, math.Float64frombits(binary.BigEndian.Uint64(data)))
	})
}

func TestReader_Read_FixStr(t *testing.T) {
	cases := []struct {
		name string
		s    string
	}{
		{"empty", ""},
		{"one", "a"},
		{"max", "abcdefghijklmnopqrstuvwxyz12345"}, // len 31
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			buf := append([]byte{byte(FixStr) | byte(len(tc.s))}, tc.s...)
			r := &MsgpReader{Buff: buf}
			ty, n, data, err := r.Read()
			require.NoError(t, err)
			assert.Equal(t, Type(byte(FixStr)|byte(len(tc.s))), ty)
			assert.Equal(t, 0, n)
			assert.Equal(t, tc.s, string(data))
		})
	}
}

func TestReader_Read_Str(t *testing.T) {
	t.Run("str8", func(t *testing.T) {
		s := make([]byte, 100)
		for i := range s {
			s[i] = byte(i)
		}
		buf := append([]byte{byte(Str8), 100}, s...)
		r := &MsgpReader{Buff: buf}
		ty, _, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Str8, ty)
		assert.Equal(t, s, data)
	})
	t.Run("str16", func(t *testing.T) {
		s := make([]byte, 300)
		buf := []byte{byte(Str16), 0x01, 0x2c} // 300
		buf = append(buf, s...)
		r := &MsgpReader{Buff: buf}
		ty, _, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Str16, ty)
		assert.Equal(t, 300, len(data))
	})
	t.Run("str32", func(t *testing.T) {
		s := make([]byte, 70000)
		buf := []byte{byte(Str32), 0x00, 0x01, 0x11, 0x70} // 70000
		buf = append(buf, s...)
		r := &MsgpReader{Buff: buf}
		ty, _, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Str32, ty)
		assert.Equal(t, 70000, len(data))
	})
}

func TestReader_Read_Bin(t *testing.T) {
	t.Run("bin8", func(t *testing.T) {
		buf := []byte{byte(Bin8), 0x03, 0xaa, 0xbb, 0xcc}
		r := &MsgpReader{Buff: buf}
		ty, _, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Bin8, ty)
		assert.Equal(t, []byte{0xaa, 0xbb, 0xcc}, data)
	})
	t.Run("bin16", func(t *testing.T) {
		buf := []byte{byte(Bin16), 0x00, 0x02, 0x01, 0x02}
		r := &MsgpReader{Buff: buf}
		ty, _, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Bin16, ty)
		assert.Equal(t, []byte{0x01, 0x02}, data)
	})
	t.Run("bin32", func(t *testing.T) {
		buf := []byte{byte(Bin32), 0x00, 0x00, 0x00, 0x01, 0xff}
		r := &MsgpReader{Buff: buf}
		ty, _, data, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Bin32, ty)
		assert.Equal(t, []byte{0xff}, data)
	})
}

func TestReader_Read_Array16_32(t *testing.T) {
	t.Run("array16", func(t *testing.T) {
		buf := []byte{byte(Array16), 0x00, 0x05}
		r := &MsgpReader{Buff: buf}
		ty, n, _, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Array16, ty)
		assert.Equal(t, 5, n)
		assert.Equal(t, 3, r.Idx)
	})
	t.Run("array32", func(t *testing.T) {
		buf := []byte{byte(Array32), 0x00, 0x00, 0x10, 0x00}
		r := &MsgpReader{Buff: buf}
		ty, n, _, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Array32, ty)
		assert.Equal(t, 4096, n)
		assert.Equal(t, 5, r.Idx)
	})
}

func TestReader_Read_Map16_32(t *testing.T) {
	t.Run("map16", func(t *testing.T) {
		buf := []byte{byte(Map16), 0x00, 0x07}
		r := &MsgpReader{Buff: buf}
		ty, n, _, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Map16, ty)
		assert.Equal(t, 7, n)
	})
	t.Run("map32", func(t *testing.T) {
		buf := []byte{byte(Map32), 0x00, 0x01, 0x00, 0x00}
		r := &MsgpReader{Buff: buf}
		ty, n, _, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Map32, ty)
		assert.Equal(t, 65536, n)
	})
}

func TestReader_Read_FixExt(t *testing.T) {
	cases := []struct {
		name    string
		tag     Type
		extType byte
		data    []byte
	}{
		{"fixext1", FixExt1, 0x05, []byte{0xaa}},
		{"fixext2", FixExt2, 0x06, []byte{0xaa, 0xbb}},
		{"fixext4", FixExt4, 0x07, []byte{0x01, 0x02, 0x03, 0x04}},
		{"fixext8", FixExt8, 0x08, []byte{1, 2, 3, 4, 5, 6, 7, 8}},
		{"fixext16", FixExt16, 0x09, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			buf := []byte{byte(tc.tag), tc.extType}
			buf = append(buf, tc.data...)
			r := &MsgpReader{Buff: buf}
			ty, n, payload, err := r.Read()
			require.NoError(t, err)
			assert.Equal(t, tc.tag, ty)
			assert.Equal(t, 0, n)
			require.Equal(t, 1+len(tc.data), len(payload))
			assert.Equal(t, tc.extType, payload[0])
			assert.Equal(t, tc.data, payload[1:])
			assert.Equal(t, len(buf), r.Idx)
		})
	}
}

func TestReader_Read_Ext(t *testing.T) {
	t.Run("ext8", func(t *testing.T) {
		buf := []byte{byte(Ext8), 0x03, 0x42, 0x01, 0x02, 0x03}
		r := &MsgpReader{Buff: buf}
		ty, _, payload, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Ext8, ty)
		assert.Equal(t, byte(0x42), payload[0])
		assert.Equal(t, []byte{0x01, 0x02, 0x03}, payload[1:])
	})
	t.Run("ext16", func(t *testing.T) {
		buf := []byte{byte(Ext16), 0x00, 0x02, 0x10, 0xaa, 0xbb}
		r := &MsgpReader{Buff: buf}
		ty, _, payload, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Ext16, ty)
		assert.Equal(t, byte(0x10), payload[0])
		assert.Equal(t, []byte{0xaa, 0xbb}, payload[1:])
	})
	t.Run("ext32", func(t *testing.T) {
		buf := []byte{byte(Ext32), 0x00, 0x00, 0x00, 0x01, 0x20, 0xff}
		r := &MsgpReader{Buff: buf}
		ty, _, payload, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Ext32, ty)
		assert.Equal(t, byte(0x20), payload[0])
		assert.Equal(t, []byte{0xff}, payload[1:])
	})
	t.Run("ext8_zero_length", func(t *testing.T) {
		buf := []byte{byte(Ext8), 0x00, 0x77}
		r := &MsgpReader{Buff: buf}
		ty, _, payload, err := r.Read()
		require.NoError(t, err)
		assert.Equal(t, Ext8, ty)
		assert.Equal(t, []byte{0x77}, payload)
	})
}

func TestReader_Read_EOF(t *testing.T) {
	r := &MsgpReader{Buff: nil}
	_, _, _, err := r.Read()
	assert.True(t, errors.Is(err, io.EOF))
}

func TestReader_Read_Unknown(t *testing.T) {
	r := &MsgpReader{Buff: []byte{0xc1}} // reserved per msgpack spec
	ty, _, _, err := r.Read()
	assert.Equal(t, Type(0xc1), ty)
	assert.True(t, errors.Is(err, ErrUnknownType))
}

func TestReader_Read_Truncated(t *testing.T) {
	cases := []struct {
		name string
		buf  []byte
	}{
		{"int8_no_data", []byte{byte(Int8)}},
		{"int16_partial", []byte{byte(Int16), 0x01}},
		{"int32_partial", []byte{byte(Int32), 0x01, 0x02}},
		{"int64_partial", []byte{byte(Int64), 0x01, 0x02, 0x03}},
		{"float32_partial", []byte{byte(Float32), 0x01}},
		{"float64_partial", []byte{byte(Float64), 0x01, 0x02}},
		{"fixstr_short", []byte{byte(FixStr) | 0x05, 'a', 'b'}},
		{"str8_no_len", []byte{byte(Str8)}},
		{"str8_short_data", []byte{byte(Str8), 0x05, 'a', 'b'}},
		{"str16_partial_len", []byte{byte(Str16), 0x00}},
		{"str16_short_data", []byte{byte(Str16), 0x00, 0x05, 'a'}},
		{"str32_partial_len", []byte{byte(Str32), 0x00, 0x00}},
		{"str32_short_data", []byte{byte(Str32), 0x00, 0x00, 0x00, 0x05, 'a'}},
		{"bin8_short_data", []byte{byte(Bin8), 0x05, 'a'}},
		{"bin16_short_data", []byte{byte(Bin16), 0x00, 0x05}},
		{"bin32_short_data", []byte{byte(Bin32), 0x00, 0x00, 0x00, 0x05}},
		{"ext8_no_typebyte", []byte{byte(Ext8), 0x01}},
		{"ext8_short_data", []byte{byte(Ext8), 0x05, 0x01, 'a'}},
		{"ext16_partial_len", []byte{byte(Ext16), 0x00}},
		{"ext16_short_data", []byte{byte(Ext16), 0x00, 0x05, 0x01}},
		{"ext32_partial_len", []byte{byte(Ext32), 0x00, 0x00}},
		{"ext32_short_data", []byte{byte(Ext32), 0x00, 0x00, 0x00, 0x05, 0x01}},
		{"fixext1_no_data", []byte{byte(FixExt1)}},
		{"fixext1_short", []byte{byte(FixExt1), 0x01}},
		{"fixext2_short", []byte{byte(FixExt2), 0x01, 0xaa}},
		{"fixext4_short", []byte{byte(FixExt4), 0x01, 0xaa, 0xbb}},
		{"fixext8_short", []byte{byte(FixExt8), 0x01, 0xaa, 0xbb, 0xcc, 0xdd}},
		{"fixext16_short", []byte{byte(FixExt16), 0x01, 0xaa}},
		{"array16_partial_len", []byte{byte(Array16), 0x00}},
		{"array32_partial_len", []byte{byte(Array32), 0x00, 0x00}},
		{"map16_partial_len", []byte{byte(Map16), 0x00}},
		{"map32_partial_len", []byte{byte(Map32), 0x00, 0x00}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &MsgpReader{Buff: tc.buf}
			_, _, _, err := r.Read()
			assert.True(t, errors.Is(err, ErrTruncated), "expected ErrTruncated, got %v", err)
		})
	}
}

func TestReader_Skip(t *testing.T) {
	w := &MsgpWriter{}
	require.NoError(t, w.WriteInt(123))
	require.NoError(t, w.WriteString("ab"))
	require.NoError(t, w.WriteFloat32(1.5))

	r := &MsgpReader{Buff: w.Buff}
	require.NoError(t, r.Skip())
	require.NoError(t, r.Skip())
	ty, _, data, err := r.Read()
	require.NoError(t, err)
	assert.Equal(t, Float32, ty)
	assert.Equal(t, float32(1.5), math.Float32frombits(binary.BigEndian.Uint32(data)))

	// next Skip → EOF
	assert.True(t, errors.Is(r.Skip(), io.EOF))
}

func TestReader_Read_Nested(t *testing.T) {
	// Build: Array16(2) -> [Map16(1) -> {"k": 42}, "tail"]
	w := &MsgpWriter{}
	require.NoError(t, w.WriteArray16(2))
	require.NoError(t, w.WriteMap16(1))
	require.NoError(t, w.WriteString("k"))
	require.NoError(t, w.WriteUint8(42))
	require.NoError(t, w.WriteString("tail"))

	r := &MsgpReader{Buff: w.Buff}

	ty, n, _, err := r.Read()
	require.NoError(t, err)
	assert.Equal(t, Array16, ty)
	assert.Equal(t, 2, n)

	ty, n, _, err = r.Read()
	require.NoError(t, err)
	assert.Equal(t, Map16, ty)
	assert.Equal(t, 1, n)

	ty, _, data, err := r.Read()
	require.NoError(t, err)
	assert.Equal(t, Type(byte(FixStr)|1), ty)
	assert.Equal(t, "k", string(data))

	ty, _, data, err = r.Read()
	require.NoError(t, err)
	assert.Equal(t, Uint8, ty)
	assert.Equal(t, byte(42), data[0])

	ty, _, data, err = r.Read()
	require.NoError(t, err)
	assert.Equal(t, Type(byte(FixStr)|4), ty)
	assert.Equal(t, "tail", string(data))

	_, _, _, err = r.Read()
	assert.True(t, errors.Is(err, io.EOF))
}

// IMsgpReader is satisfied by *MsgpReader.
var _ IMsgpReader = (*MsgpReader)(nil)
