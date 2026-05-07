package msgpraw

import (
	"bytes"
	"math/bits"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// IMsgpWriter is satisfied by *MsgpWriter.
var _ IMsgpWriter = (*MsgpWriter)(nil)

func TestWriter_WritePosFixInt_Range(t *testing.T) {
	w := &MsgpWriter{}
	require.NoError(t, w.WritePosFixInt(0))
	require.NoError(t, w.WritePosFixInt(0x7f))
	assert.Equal(t, []byte{0x00, 0x7f}, w.Buff)

	require.ErrorIs(t, w.WritePosFixInt(0x80), ErrPosFixIntRange)
}

func TestWriter_WriteString_AutoSize(t *testing.T) {
	cases := []struct {
		name    string
		length  int
		wantTag byte
	}{
		{"empty_fixstr", 0, byte(FixStr)},
		{"fixstr_max", 31, byte(FixStr) | 31},
		{"str8_min", 32, byte(Str8)},
		{"str8_max", 255, byte(Str8)},
		{"str16_min", 256, byte(Str16)},
		{"str16_max", 65535, byte(Str16)},
		{"str32_min", 65536, byte(Str32)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := &MsgpWriter{}
			require.NoError(t, w.WriteString(strings.Repeat("x", tc.length)))
			assert.Equal(t, tc.wantTag, w.Buff[0])

			// Roundtrip.
			r := &MsgpReader{Buff: w.Buff}
			_, _, data, err := r.Read()
			require.NoError(t, err)
			assert.Equal(t, tc.length, len(data))
		})
	}
}

func TestWriter_WriteBytes_AutoSize(t *testing.T) {
	cases := []struct {
		name    string
		length  int
		wantTag byte
	}{
		{"bin8_empty", 0, byte(Bin8)},
		{"bin8_max", 255, byte(Bin8)},
		{"bin16_min", 256, byte(Bin16)},
		{"bin16_max", 65535, byte(Bin16)},
		{"bin32_min", 65536, byte(Bin32)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := &MsgpWriter{}
			require.NoError(t, w.WriteBytes(make([]byte, tc.length)))
			assert.Equal(t, tc.wantTag, w.Buff[0])

			r := &MsgpReader{Buff: w.Buff}
			_, _, data, err := r.Read()
			require.NoError(t, err)
			assert.Equal(t, tc.length, len(data))
		})
	}
}

func TestWriter_WriteArray_TooLong_64bit(t *testing.T) {
	if bits.UintSize < 64 {
		t.Skip("requires 64-bit int")
	}
	w := &MsgpWriter{}
	require.ErrorIs(t, w.WriteArray(int(maxUint32)+1), ErrArrayTooLong)
	assert.Empty(t, w.Buff)
}

func TestWriter_WriteMap_TooLong_64bit(t *testing.T) {
	if bits.UintSize < 64 {
		t.Skip("requires 64-bit int")
	}
	w := &MsgpWriter{}
	require.ErrorIs(t, w.WriteMap(int(maxUint32)+1), ErrMapTooLong)
}

func TestWriter_WriteArray_AutoSize(t *testing.T) {
	cases := []struct {
		name    string
		n       int
		wantTag byte
	}{
		{"fixarray_zero", 0, byte(FixArray)},
		{"fixarray_max", 15, byte(FixArray) | 15},
		{"array16_min", 16, byte(Array16)},
		{"array16_max", 65535, byte(Array16)},
		{"array32_min", 65536, byte(Array32)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := &MsgpWriter{}
			require.NoError(t, w.WriteArray(tc.n))
			assert.Equal(t, tc.wantTag, w.Buff[0])

			r := &MsgpReader{Buff: w.Buff}
			_, n, _, err := r.Read()
			require.NoError(t, err)
			assert.Equal(t, tc.n, n)
		})
	}
	t.Run("negative", func(t *testing.T) {
		w := &MsgpWriter{}
		require.ErrorIs(t, w.WriteArray(-1), ErrArray32Range)
	})
}

func TestWriter_WriteMap_AutoSize(t *testing.T) {
	cases := []struct {
		name    string
		n       int
		wantTag byte
	}{
		{"fixmap_zero", 0, byte(FixMap)},
		{"fixmap_max", 15, byte(FixMap) | 15},
		{"map16_min", 16, byte(Map16)},
		{"map16_max", 65535, byte(Map16)},
		{"map32_min", 65536, byte(Map32)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := &MsgpWriter{}
			require.NoError(t, w.WriteMap(tc.n))
			assert.Equal(t, tc.wantTag, w.Buff[0])

			r := &MsgpReader{Buff: w.Buff}
			_, n, _, err := r.Read()
			require.NoError(t, err)
			assert.Equal(t, tc.n, n)
		})
	}
	t.Run("negative", func(t *testing.T) {
		w := &MsgpWriter{}
		require.ErrorIs(t, w.WriteMap(-1), ErrMap32Range)
	})
}

// --- Explicit sizing methods (range enforcement + correct tag) --------------

func TestWriter_WriteFixStr(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		w := &MsgpWriter{}
		require.NoError(t, w.WriteFixStr("abc"))
		assert.Equal(t, []byte{byte(FixStr) | 3, 'a', 'b', 'c'}, w.Buff)
	})
	t.Run("too_long", func(t *testing.T) {
		w := &MsgpWriter{}
		require.ErrorIs(t, w.WriteFixStr(strings.Repeat("x", 32)), ErrFixStrRange)
		assert.Empty(t, w.Buff)
	})
}

func TestWriter_WriteStr8(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		w := &MsgpWriter{}
		require.NoError(t, w.WriteStr8(strings.Repeat("a", 200)))
		assert.Equal(t, byte(Str8), w.Buff[0])
		assert.Equal(t, byte(200), w.Buff[1])
	})
	t.Run("too_long", func(t *testing.T) {
		w := &MsgpWriter{}
		require.ErrorIs(t, w.WriteStr8(strings.Repeat("x", 256)), ErrStr8Range)
	})
}

func TestWriter_WriteStr16(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		w := &MsgpWriter{}
		require.NoError(t, w.WriteStr16(strings.Repeat("z", 1000)))
		assert.Equal(t, byte(Str16), w.Buff[0])
		assert.Equal(t, byte(0x03), w.Buff[1])
		assert.Equal(t, byte(0xe8), w.Buff[2])
	})
	t.Run("too_long", func(t *testing.T) {
		w := &MsgpWriter{}
		require.ErrorIs(t, w.WriteStr16(strings.Repeat("z", 65536)), ErrStr16Range)
	})
}

func TestWriter_WriteStr32(t *testing.T) {
	w := &MsgpWriter{}
	require.NoError(t, w.WriteStr32(strings.Repeat("y", 70000)))
	assert.Equal(t, byte(Str32), w.Buff[0])
}

func TestWriter_WriteBin8(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		w := &MsgpWriter{}
		require.NoError(t, w.WriteBin8([]byte{0x01, 0x02}))
		assert.Equal(t, []byte{byte(Bin8), 0x02, 0x01, 0x02}, w.Buff)
	})
	t.Run("too_long", func(t *testing.T) {
		w := &MsgpWriter{}
		require.ErrorIs(t, w.WriteBin8(make([]byte, 256)), ErrBin8Range)
	})
}

func TestWriter_WriteBin16(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		w := &MsgpWriter{}
		require.NoError(t, w.WriteBin16(make([]byte, 1)))
		assert.Equal(t, byte(Bin16), w.Buff[0])
	})
	t.Run("too_long", func(t *testing.T) {
		w := &MsgpWriter{}
		require.ErrorIs(t, w.WriteBin16(make([]byte, 65536)), ErrBin16Range)
	})
}

func TestWriter_WriteBin32(t *testing.T) {
	w := &MsgpWriter{}
	require.NoError(t, w.WriteBin32(make([]byte, 1)))
	assert.Equal(t, byte(Bin32), w.Buff[0])
}

func TestWriter_WriteFixArray(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		w := &MsgpWriter{}
		require.NoError(t, w.WriteFixArray(7))
		assert.Equal(t, []byte{byte(FixArray) | 7}, w.Buff)
	})
	t.Run("too_big", func(t *testing.T) {
		w := &MsgpWriter{}
		require.ErrorIs(t, w.WriteFixArray(16), ErrFixArrayRange)
	})
	t.Run("negative", func(t *testing.T) {
		w := &MsgpWriter{}
		require.ErrorIs(t, w.WriteFixArray(-1), ErrFixArrayRange)
	})
}

func TestWriter_WriteArray16(t *testing.T) {
	w := &MsgpWriter{}
	require.NoError(t, w.WriteArray16(0xabcd))
	assert.Equal(t, []byte{byte(Array16), 0xab, 0xcd}, w.Buff)

	require.ErrorIs(t, w.WriteArray16(65536), ErrArray16Range)
	require.ErrorIs(t, w.WriteArray16(-1), ErrArray16Range)
}

func TestWriter_WriteArray32(t *testing.T) {
	w := &MsgpWriter{}
	require.NoError(t, w.WriteArray32(0x10000))
	assert.Equal(t, []byte{byte(Array32), 0x00, 0x01, 0x00, 0x00}, w.Buff)

	require.ErrorIs(t, w.WriteArray32(-1), ErrArray32Range)
}

func TestWriter_WriteFixMap(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		w := &MsgpWriter{}
		require.NoError(t, w.WriteFixMap(3))
		assert.Equal(t, []byte{byte(FixMap) | 3}, w.Buff)
	})
	t.Run("too_big", func(t *testing.T) {
		w := &MsgpWriter{}
		require.ErrorIs(t, w.WriteFixMap(16), ErrFixMapRange)
	})
	t.Run("negative", func(t *testing.T) {
		w := &MsgpWriter{}
		require.ErrorIs(t, w.WriteFixMap(-1), ErrFixMapRange)
	})
}

func TestWriter_WriteMap16(t *testing.T) {
	w := &MsgpWriter{}
	require.NoError(t, w.WriteMap16(0x1234))
	assert.Equal(t, []byte{byte(Map16), 0x12, 0x34}, w.Buff)

	require.ErrorIs(t, w.WriteMap16(65536), ErrMap16Range)
	require.ErrorIs(t, w.WriteMap16(-1), ErrMap16Range)
}

func TestWriter_WriteMap32(t *testing.T) {
	w := &MsgpWriter{}
	require.NoError(t, w.WriteMap32(0x10000))
	assert.Equal(t, []byte{byte(Map32), 0x00, 0x01, 0x00, 0x00}, w.Buff)

	require.ErrorIs(t, w.WriteMap32(-1), ErrMap32Range)
}

// --- Ext writers ------------------------------------------------------------

func TestWriter_WriteFixExt(t *testing.T) {
	cases := []struct {
		name string
		size int
		fn   func(*MsgpWriter, int8, []byte) error
		tag  Type
	}{
		{"fixext1", 1, (*MsgpWriter).WriteFixExt1, FixExt1},
		{"fixext2", 2, (*MsgpWriter).WriteFixExt2, FixExt2},
		{"fixext4", 4, (*MsgpWriter).WriteFixExt4, FixExt4},
		{"fixext8", 8, (*MsgpWriter).WriteFixExt8, FixExt8},
		{"fixext16", 16, (*MsgpWriter).WriteFixExt16, FixExt16},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := bytes.Repeat([]byte{0xab}, tc.size)
			w := &MsgpWriter{}
			require.NoError(t, tc.fn(w, 0x42, data))
			assert.Equal(t, byte(tc.tag), w.Buff[0])
			assert.Equal(t, byte(0x42), w.Buff[1])
			assert.Equal(t, data, w.Buff[2:])

			// wrong size errors
			w = &MsgpWriter{}
			require.ErrorIs(t, tc.fn(w, 0x42, make([]byte, tc.size+1)), ErrFixExtSize)
			require.ErrorIs(t, tc.fn(w, 0x42, nil), ErrFixExtSize)
		})
	}
}

func TestWriter_WriteExt8(t *testing.T) {
	w := &MsgpWriter{}
	require.NoError(t, w.WriteExt8(0x05, []byte{0xaa, 0xbb, 0xcc}))
	assert.Equal(t, []byte{byte(Ext8), 0x03, 0x05, 0xaa, 0xbb, 0xcc}, w.Buff)

	require.ErrorIs(t, w.WriteExt8(0x05, make([]byte, 256)), ErrExt8Range)
}

func TestWriter_WriteExt16(t *testing.T) {
	w := &MsgpWriter{}
	require.NoError(t, w.WriteExt16(0x06, []byte{0xff}))
	assert.Equal(t, []byte{byte(Ext16), 0x00, 0x01, 0x06, 0xff}, w.Buff)

	require.ErrorIs(t, w.WriteExt16(0x06, make([]byte, 65536)), ErrExt16Range)
}

func TestWriter_WriteExt32(t *testing.T) {
	w := &MsgpWriter{}
	require.NoError(t, w.WriteExt32(0x07, []byte{0x11}))
	assert.Equal(t, []byte{byte(Ext32), 0x00, 0x00, 0x00, 0x01, 0x07, 0x11}, w.Buff)
}

func TestWriter_WriteExt_AutoSize(t *testing.T) {
	cases := []struct {
		name    string
		dataLen int
		wantTag byte
	}{
		{"empty_uses_ext8", 0, byte(Ext8)},
		{"len1_uses_fixext1", 1, byte(FixExt1)},
		{"len2_uses_fixext2", 2, byte(FixExt2)},
		{"len3_uses_ext8", 3, byte(Ext8)},
		{"len4_uses_fixext4", 4, byte(FixExt4)},
		{"len5_uses_ext8", 5, byte(Ext8)},
		{"len8_uses_fixext8", 8, byte(FixExt8)},
		{"len9_uses_ext8", 9, byte(Ext8)},
		{"len16_uses_fixext16", 16, byte(FixExt16)},
		{"len17_uses_ext8", 17, byte(Ext8)},
		{"len255_uses_ext8", 255, byte(Ext8)},
		{"len256_uses_ext16", 256, byte(Ext16)},
		{"len65535_uses_ext16", 65535, byte(Ext16)},
		{"len65536_uses_ext32", 65536, byte(Ext32)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := &MsgpWriter{}
			require.NoError(t, w.WriteExt(0x33, make([]byte, tc.dataLen)))
			assert.Equal(t, tc.wantTag, w.Buff[0])

			// Roundtrip via reader.
			r := &MsgpReader{Buff: w.Buff}
			_, _, payload, err := r.Read()
			require.NoError(t, err)
			require.Equal(t, 1+tc.dataLen, len(payload))
			assert.Equal(t, byte(0x33), payload[0])
		})
	}
}

// --- WriteInt promotion behavior (current: always Int64) --------------------

func TestWriter_WriteInt_AlwaysInt64(t *testing.T) {
	w := &MsgpWriter{}
	require.NoError(t, w.WriteInt(0))
	assert.Equal(t, byte(Int64), w.Buff[0])

	w = &MsgpWriter{}
	require.NoError(t, w.WriteInt(-1))
	assert.Equal(t, byte(Int64), w.Buff[0])
}

func TestWriter_WriteUint_AlwaysUint64(t *testing.T) {
	w := &MsgpWriter{}
	require.NoError(t, w.WriteUint(0))
	assert.Equal(t, byte(Uint64), w.Buff[0])
}
