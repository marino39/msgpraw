package msgpraw

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

// isTagEncodedTag reports whether typ carries its value in the tag byte
// itself, so Read returns a nil payload for it.
func isTagEncodedTag(typ Type) bool {
	return typ == Nil || typ == True || typ == False || typ <= PosFixIntMax || typ >= NegFixInt
}

// isCollectionTag reports whether typ is an array or map header, whose
// payload is the remainder of the buffer rather than the value's own bytes.
func isCollectionTag(typ Type) bool {
	return (typ >= FixArray && typ <= FixArrayMax) || typ == Array16 || typ == Array32 ||
		(typ >= FixMap && typ <= FixMapMax) || typ == Map16 || typ == Map32
}

// checkReadInvariants drains data through MsgpReader.Read and asserts that
// the reader never panics, never returns an error outside the documented
// set, and never hands back a payload slice that reaches outside data. It
// then repeats the pass with Skip and requires an identical outcome.
func checkReadInvariants(t *testing.T, data []byte) {
	t.Helper()

	// cap-pin: a fuzzer-supplied []byte can have cap > len, which would
	// make the off := len(buf) - cap(p) arithmetic below overshoot.
	buf := data[:len(data):len(data)]
	r := MsgpReader{Buff: buf}

	var terminalErr error
	iterations := 0
	for {
		prev := r.Idx
		typ, count, p, err := r.Read()
		iterations++
		if iterations > len(buf)+1 {
			t.Fatalf("Read failed to terminate")
		}

		if err != nil {
			require.Truef(t, errors.Is(err, EOF) || errors.Is(err, ErrTruncated) || errors.Is(err, ErrUnknownType),
				"unexpected error %v", err)
			require.LessOrEqual(t, r.Idx, len(buf))
			if errors.Is(err, ErrUnknownType) {
				require.Less(t, prev, len(buf))
				require.Equal(t, byte(0xc1), buf[prev])
			}
			if errors.Is(err, EOF) {
				require.Equal(t, prev, r.Idx)
				require.Equal(t, len(buf), r.Idx)
			}
			terminalErr = err
			break
		}

		require.Greater(t, r.Idx, prev)
		require.LessOrEqual(t, r.Idx, len(buf))
		require.Equal(t, typ, Type(buf[prev]))

		off := len(buf) - cap(p)
		require.GreaterOrEqual(t, off, 0)
		require.LessOrEqual(t, off, len(buf))
		require.LessOrEqual(t, off+len(p), len(buf))
		require.LessOrEqual(t, cap(p), len(buf))

		switch {
		case isTagEncodedTag(typ):
			require.Nil(t, p)
			require.Equal(t, 0, count)
		case isCollectionTag(typ):
			require.GreaterOrEqual(t, count, 0)
			require.Equal(t, r.Idx, off)
			require.Equal(t, len(buf), off+len(p))
		default:
			require.Equal(t, 0, count)
			require.Equal(t, r.Idx, off+len(p))
		}
	}

	finalIdx := r.Idx

	r2 := MsgpReader{Buff: buf}
	var skipErr error
	iterations = 0
	for {
		skipErr = r2.Skip()
		iterations++
		if iterations > len(buf)+1 {
			t.Fatalf("Skip failed to terminate")
		}
		if skipErr != nil {
			break
		}
	}
	require.Equal(t, terminalErr, skipErr)
	require.Equal(t, finalIdx, r2.Idx)
}

func FuzzReader(f *testing.F) {
	f.Add(allTagsFixture(f))
	f.Add([]byte{})
	for b := 0; b <= 0xff; b++ {
		f.Add([]byte{byte(b)})
	}
	f.Add([]byte{byte(Bin8), 0x05, 'a'})
	f.Add([]byte{byte(Str16), 0x01})
	f.Add([]byte{byte(Ext8), 0x03, 0x01})
	f.Add([]byte{byte(Array16), 0x00})
	f.Add([]byte{byte(Map32), 0x00, 0x00, 0x00})
	f.Add([]byte{byte(FixStr) | 0x03, 'a'})
	f.Add(bytes.Repeat([]byte{byte(FixArray) | 0x01}, 64))
	f.Add(bytes.Repeat([]byte{byte(FixMap) | 0x01}, 64))

	f.Fuzz(func(t *testing.T, data []byte) {
		checkReadInvariants(t, data)
	})
}

// readBack reads exactly one value from buf and requires the reader to
// then report EOF, proving the writer emitted nothing beyond that value.
func readBack(t *testing.T, buf []byte) (Type, int, []byte) {
	t.Helper()

	r := MsgpReader{Buff: buf}
	typ, count, p, err := r.Read()
	require.NoError(t, err)

	_, _, _, err = r.Read()
	require.ErrorIs(t, err, EOF)

	return typ, count, p
}

func expectedStrType(n int) Type {
	switch {
	case n <= maxFixStr:
		return Type(byte(FixStr) | byte(n))
	case n <= maxUint8:
		return Str8
	case n <= maxUint16:
		return Str16
	default:
		return Str32
	}
}

func expectedBytesType(n int) Type {
	switch {
	case n <= maxUint8:
		return Bin8
	case n <= maxUint16:
		return Bin16
	default:
		return Bin32
	}
}

func expectedExtType(n int) Type {
	switch n {
	case 1:
		return FixExt1
	case 2:
		return FixExt2
	case 4:
		return FixExt4
	case 8:
		return FixExt8
	case 16:
		return FixExt16
	}
	switch {
	case n <= maxUint8:
		return Ext8
	case n <= maxUint16:
		return Ext16
	default:
		return Ext32
	}
}

func expectedArrayType(n int) Type {
	switch {
	case n <= maxFixArray:
		return Type(byte(FixArray) | byte(n))
	case n <= maxUint16:
		return Array16
	default:
		return Array32
	}
}

func expectedMapType(n int) Type {
	switch {
	case n <= maxFixMap:
		return Type(byte(FixMap) | byte(n))
	case n <= maxUint16:
		return Map16
	default:
		return Map32
	}
}

func FuzzRoundTrip(f *testing.F) {
	f.Add(uint8(0), int64(0), float64(0), "", []byte(nil), int8(0), uint32(0))
	f.Add(uint8(4), int64(1)<<40, -0.0, "hello world", []byte{1, 2, 3}, int8(-5), uint32(1000))
	f.Add(uint8(10), int64(-123456789), math.Inf(1), string(make([]byte, 40)), make([]byte, 300), int8(127), uint32(70000))
	f.Add(uint8(255), int64(math.MaxInt64), math.NaN(), "", []byte{}, int8(-128), uint32(4000000000))

	f.Fuzz(func(t *testing.T, sel uint8, i64 int64, f64 float64, s string, b []byte, extType int8, n uint32) {
		w := &MsgpWriter{}

		switch sel % 11 {
		case 0:
			require.NoError(t, w.WriteInt64(i64))
			typ, count, p := readBack(t, w.Buff)
			require.Equal(t, Int64, typ)
			require.Equal(t, 0, count)
			require.Len(t, p, 8)
			require.Equal(t, uint64(i64), binary.BigEndian.Uint64(p))

		case 1:
			require.NoError(t, w.WriteUint64(uint64(i64)))
			typ, count, p := readBack(t, w.Buff)
			require.Equal(t, Uint64, typ)
			require.Equal(t, 0, count)
			require.Len(t, p, 8)
			require.Equal(t, uint64(i64), binary.BigEndian.Uint64(p))

		case 2:
			require.NoError(t, w.WriteFloat64(f64))
			typ, count, p := readBack(t, w.Buff)
			require.Equal(t, Float64, typ)
			require.Equal(t, 0, count)
			require.Equal(t, math.Float64bits(f64), binary.BigEndian.Uint64(p))

		case 3:
			require.NoError(t, w.WriteFloat32(float32(f64)))
			typ, count, p := readBack(t, w.Buff)
			require.Equal(t, Float32, typ)
			require.Equal(t, 0, count)
			require.Equal(t, math.Float32bits(float32(f64)), binary.BigEndian.Uint32(p))

		case 4:
			v := uint8(n % 128)
			require.NoError(t, w.WritePosFixInt(v))
			typ, count, p := readBack(t, w.Buff)
			require.Equal(t, Type(v), typ)
			require.Equal(t, 0, count)
			require.Nil(t, p)

		case 5:
			v := int8(-1 - int8(n%32))
			require.NoError(t, w.WriteNegFixInt(v))
			typ, count, p := readBack(t, w.Buff)
			require.Equal(t, Type(byte(v)), typ)
			require.Equal(t, 0, count)
			require.Nil(t, p)

		case 6:
			if sel&0x10 != 0 {
				require.NoError(t, w.WriteNil())
				typ, count, p := readBack(t, w.Buff)
				require.Equal(t, Nil, typ)
				require.Equal(t, 0, count)
				require.Nil(t, p)
			} else {
				want := sel&1 == 0
				require.NoError(t, w.WriteBool(want))
				typ, count, p := readBack(t, w.Buff)
				if want {
					require.Equal(t, True, typ)
				} else {
					require.Equal(t, False, typ)
				}
				require.Equal(t, 0, count)
				require.Nil(t, p)
			}

		case 7:
			require.NoError(t, w.WriteString(s))
			typ, count, p := readBack(t, w.Buff)
			require.Equal(t, expectedStrType(len(s)), typ)
			require.Equal(t, 0, count)
			require.Equal(t, s, string(p))

		case 8:
			require.NoError(t, w.WriteBytes(b))
			typ, count, p := readBack(t, w.Buff)
			require.Equal(t, expectedBytesType(len(b)), typ)
			require.Equal(t, 0, count)
			require.True(t, bytes.Equal(p, b))

		case 9:
			require.NoError(t, w.WriteExt(extType, b))
			typ, count, p := readBack(t, w.Buff)
			require.Equal(t, expectedExtType(len(b)), typ)
			require.Equal(t, 0, count)
			require.Len(t, p, len(b)+1)
			require.Equal(t, byte(extType), p[0])
			require.True(t, bytes.Equal(p[1:], b))

		case 10:
			cnt := int(n % 100000)
			if sel&1 == 0 {
				require.NoError(t, w.WriteArray(cnt))
				typ, count, p := readBack(t, w.Buff)
				require.Equal(t, expectedArrayType(cnt), typ)
				require.Equal(t, cnt, count)
				require.Empty(t, p)
			} else {
				require.NoError(t, w.WriteMap(cnt))
				typ, count, p := readBack(t, w.Buff)
				require.Equal(t, expectedMapType(cnt), typ)
				require.Equal(t, cnt, count)
				require.Empty(t, p)
			}
		}

		checkReadInvariants(t, w.Buff)
	})
}
