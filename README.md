# msgpraw

[![Go Reference](https://pkg.go.dev/badge/github.com/marino39/msgpraw.svg)](https://pkg.go.dev/github.com/marino39/msgpraw)
[![Go Report Card](https://goreportcard.com/badge/github.com/marino39/msgpraw)](https://goreportcard.com/report/github.com/marino39/msgpraw)

A low-level [MessagePack](https://msgpack.org/) reader and writer for Go, designed for **zero-allocation decoding**. `msgpraw` does not unmarshal values into Go types — instead it returns the raw payload bytes alongside the format tag, leaving big-endian decoding to the caller. That keeps the reader allocation-free and lets callers decode lazily, only paying for what they actually use.

## Features

- Zero allocations on the reader's success path (verified by `TestReader_NoAllocs` and benchmarks).
- Minimal allocations on the writer (zero with a pre-sized buffer).
- Full MessagePack format coverage: `nil`, `bool`, all int/uint widths, `float32`/`float64`, `fixstr`/`str8`/`str16`/`str32`, `bin8`/`bin16`/`bin32`, `fixarray`/`array16`/`array32`, `fixmap`/`map16`/`map32`, `fixext1..16`, `ext8`/`ext16`/`ext32`, `posfixint`, `negfixint`.
- Both auto-sized writers (smallest format that fits) and explicit fixed-format writers.
- Pre-allocated sentinel errors — error paths don't allocate either.

## Install

```sh
go get github.com/marino39/msgpraw
```

Requires Go 1.20 or later.

## Reader

```go
import "github.com/marino39/msgpraw"

r := &msgpraw.MsgpReader{Buff: payload}
for {
    tag, count, data, err := r.Read()
    if err == msgpraw.EOF {
        break
    }
    if err != nil {
        // ErrTruncated or ErrUnknownType
        return err
    }
    switch tag {
    case msgpraw.Int64:
        // data is the 8-byte big-endian payload
        v := int64(binary.BigEndian.Uint64(data))
        _ = v
    case msgpraw.Array16, msgpraw.Array32:
        // count is the number of elements; call Read again `count` times
        _ = count
    case msgpraw.Map16, msgpraw.Map32:
        // count is the pair count; call Read 2*count times
        _ = count
    }
}
```

### Read return values

`Read()` returns `(Type, int, []byte, error)`:

| Tag family                          | `int`           | `[]byte`                                 |
|-------------------------------------|-----------------|------------------------------------------|
| `Nil`, `True`, `False`, `PosFixInt`, `NegFixInt` | `0` | `nil` (value carried by tag byte)        |
| Scalar (`Int*`, `Uint*`, `Float*`)  | `0`             | raw big-endian payload bytes             |
| String/binary (`FixStr`, `Str*`, `Bin*`) | `0`        | raw bytes                                |
| Ext (`FixExt*`, `Ext*`)             | `0`             | `[extType, dataBytes...]` — caller does `data[0]`, `data[1:]` |
| Collections (`FixArray`/`Array*`, `FixMap`/`Map*`) | element/pair count | remainder of buffer; caller calls `Read()` repeatedly |

The returned `[]byte` is a sub-slice of the input buffer — **not a copy**. Don't write to it; copy if the caller needs ownership.

### Skip

```go
err := r.Skip()
```

`Skip()` advances past the next value without inspecting it (still O(1) per scalar; for collections you must call `Skip` once per element to fully consume them).

### Errors

| Error            | When                                                         |
|------------------|--------------------------------------------------------------|
| `EOF` (`io.EOF`) | The buffer is exhausted between values.                      |
| `ErrTruncated`   | The buffer ends mid-value (truncated length prefix or data). |
| `ErrUnknownType` | The leading byte is not a defined MessagePack format.        |

All errors are pre-allocated package-level sentinels. Compare with `errors.Is`.

## Writer

```go
w := &msgpraw.MsgpWriter{Buff: make([]byte, 0, 1024)}

// Auto-sized: pick the smallest spec format that fits.
_ = w.WriteString("hello")        // FixStr
_ = w.WriteArray(2)               // FixArray
_ = w.WriteInt(42)                // Int64
_ = w.WriteFloat64(3.14159)
_ = w.WriteExt(0x05, []byte{1,2}) // picks FixExt2

// Explicit: pin the wire format.
_ = w.WriteStr16("longer payload")
_ = w.WriteArray32(1_000_000)
_ = w.WriteFixExt4(0x07, []byte{1, 2, 3, 4})

send(w.Buff)
```

### Auto-sized vs. explicit

| Auto-sized         | Explicit variants                                                          |
|--------------------|----------------------------------------------------------------------------|
| `WriteString(s)`   | `WriteFixStr` / `WriteStr8` / `WriteStr16` / `WriteStr32`                  |
| `WriteBytes(b)`    | `WriteBin8` / `WriteBin16` / `WriteBin32`                                  |
| `WriteArray(n)`    | `WriteFixArray` / `WriteArray16` / `WriteArray32`                          |
| `WriteMap(n)`      | `WriteFixMap` / `WriteMap16` / `WriteMap32`                                |
| `WriteExt(t, d)`   | `WriteFixExt1/2/4/8/16` / `WriteExt8` / `WriteExt16` / `WriteExt32`        |

Auto-sized methods choose the smallest format that fits the input and return a sentinel error if the input exceeds the largest variant. Explicit methods return a range error when the input doesn't fit the named format.

`WriteInt` / `WriteUint` always emit `Int64` / `Uint64` — use `WritePosFixInt`, `WriteNegFixInt`, or `WriteInt8..64` for compact forms. `WriteNegFixInt` accepts only `-32..-1` per spec.

### Zero-alloc writes

For zero-allocation writes, pre-size `Buff` so it doesn't have to grow:

```go
buf := make([]byte, 0, 4096)
w := msgpraw.MsgpWriter{Buff: buf}
// ... writes ...
buf = w.Buff // grown if you wrote more than 4096 bytes total
```

The numeric writers use `binary.BigEndian.AppendUint*` (Go 1.19+) directly on `Buff` — no temporary slices.

## Benchmarks

On an Apple M4 Max (`go test -bench . -benchmem -run=^$`):

```
BenchmarkReader_AllTags-16    7,376,100    139.6 ns/op    0 B/op    0 allocs/op
BenchmarkWriter_AllTags-16   65,632,200     18.4 ns/op    0 B/op    0 allocs/op
```

`BenchmarkReader_AllTags` exercises every supported tag in a single pass. `BenchmarkWriter_AllTags` writes a representative mix of tags to a pre-sized buffer.

## Testing

```sh
go test ./...                                # unit tests
go test -race -cover ./...                   # race detector + coverage
go test -run TestReader_NoAllocs ./...       # zero-alloc assertion
go test -bench . -benchmem -run=^$ ./...     # benchmarks
```

## Non-goals

- **Marshalling/unmarshalling Go types.** This is a raw codec — bring your own `binary.BigEndian.Uint*` calls, or layer a typed codec on top.
- **Streaming `io.Reader`/`io.Writer`.** The buffer-based API is what makes zero-allocation reads possible. Streaming would require an internal buffer.
- **Predefined extension types** (Timestamp, etc.). Easy to layer on top of `WriteExt` / the ext payload format.

## License

[MIT](LICENSE.md)
