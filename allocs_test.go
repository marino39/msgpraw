package msgpraw

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// allTagsFixture builds a buffer that exercises every msgp tag the reader
// handles, so a single Read loop touches every code path.
func allTagsFixture(t require.TestingT) []byte {
	w := &MsgpWriter{Buff: make([]byte, 0, 1024)}

	require.NoError(t, w.WritePosFixInt(7))
	require.NoError(t, w.WriteNegFixInt(-7))
	require.NoError(t, w.WriteNil())
	require.NoError(t, w.WriteBool(true))
	require.NoError(t, w.WriteBool(false))
	require.NoError(t, w.WriteInt8(-1))
	require.NoError(t, w.WriteInt16(-300))
	require.NoError(t, w.WriteInt32(-70000))
	require.NoError(t, w.WriteInt64(-int64(1)<<33))
	require.NoError(t, w.WriteUint8(200))
	require.NoError(t, w.WriteUint16(40000))
	require.NoError(t, w.WriteUint32(4_000_000_000))
	require.NoError(t, w.WriteUint64(uint64(1)<<40))
	require.NoError(t, w.WriteFloat32(1.5))
	require.NoError(t, w.WriteFloat64(2.71828))

	require.NoError(t, w.WriteFixStr("hi"))
	require.NoError(t, w.WriteStr8(string(make([]byte, 100))))
	require.NoError(t, w.WriteStr16(string(make([]byte, 300))))
	require.NoError(t, w.WriteStr32(string(make([]byte, 70000))))

	require.NoError(t, w.WriteBin8(make([]byte, 10)))
	require.NoError(t, w.WriteBin16(make([]byte, 300)))
	require.NoError(t, w.WriteBin32(make([]byte, 70000)))

	require.NoError(t, w.WriteFixArray(0))
	require.NoError(t, w.WriteArray16(0))
	require.NoError(t, w.WriteArray32(0))
	require.NoError(t, w.WriteFixMap(0))
	require.NoError(t, w.WriteMap16(0))
	require.NoError(t, w.WriteMap32(0))

	require.NoError(t, w.WriteFixExt1(1, []byte{0xaa}))
	require.NoError(t, w.WriteFixExt2(1, []byte{0xaa, 0xbb}))
	require.NoError(t, w.WriteFixExt4(1, make([]byte, 4)))
	require.NoError(t, w.WriteFixExt8(1, make([]byte, 8)))
	require.NoError(t, w.WriteFixExt16(1, make([]byte, 16)))
	require.NoError(t, w.WriteExt8(1, make([]byte, 3)))
	require.NoError(t, w.WriteExt16(1, make([]byte, 300)))
	require.NoError(t, w.WriteExt32(1, make([]byte, 70000)))

	return w.Buff
}

// readAllOnce drains a reader without allocating; helper used by both
// the allocation assertion and the benchmark.
func readAllOnce(buf []byte) {
	r := MsgpReader{Buff: buf}
	for {
		_, _, _, err := r.Read()
		if err != nil {
			return
		}
	}
}

func TestReader_NoAllocs(t *testing.T) {
	buf := allTagsFixture(t)

	allocs := testing.AllocsPerRun(100, func() {
		readAllOnce(buf)
	})
	require.Zero(t, allocs, "MsgpReader.Read must not allocate on success path")
}

func TestReader_NoAllocs_ErrorPaths(t *testing.T) {
	cases := [][]byte{
		nil,                      // EOF
		{byte(Int8)},             // truncated
		{0xc1},                   // unknown
		{byte(Bin8), 0x05, 'a'},  // truncated payload
	}
	for i, buf := range cases {
		buf := buf
		allocs := testing.AllocsPerRun(100, func() {
			r := MsgpReader{Buff: buf}
			_, _, _, _ = r.Read()
		})
		require.Zerof(t, allocs, "case %d (buf=%x): error path allocated", i, buf)
	}
}

func BenchmarkReader_AllTags(b *testing.B) {
	buf := allTagsFixture(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readAllOnce(buf)
	}
}

func BenchmarkWriter_AllTags(b *testing.B) {
	// Pre-sized buffer so growth doesn't dominate.
	scratch := make([]byte, 0, 200_000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := MsgpWriter{Buff: scratch[:0]}
		_ = w.WritePosFixInt(7)
		_ = w.WriteNegFixInt(-7)
		_ = w.WriteInt64(1 << 40)
		_ = w.WriteUint64(1 << 40)
		_ = w.WriteFloat64(2.71828)
		_ = w.WriteFixStr("hi")
		_ = w.WriteFixArray(0)
		_ = w.WriteFixMap(0)
		_ = w.WriteFixExt4(1, []byte{1, 2, 3, 4})
		_ = w.WriteExt8(1, []byte{1, 2, 3})
	}
}

// Sanity: the all-tags fixture round-trips cleanly to EOF without an unknown
// tag. If a future tag is added to types.go but not the reader, this fails.
func TestReader_AllTags_ReadToEOF(t *testing.T) {
	buf := allTagsFixture(t)
	r := MsgpReader{Buff: buf}
	for {
		_, _, _, err := r.Read()
		if err != nil {
			require.True(t, errors.Is(err, io.EOF))
			return
		}
	}
}
