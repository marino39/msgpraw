# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`github.com/marino39/msgpraw` is a low-level Go library for reading and writing the [MessagePack](https://msgpack.org/) wire format. It does not unmarshal values — the reader returns raw payload bytes alongside the type tag, leaving big-endian decoding to the caller. Targets Go 1.20.

Design goals:
- **Zero-allocation reads.** Every successful `Read()` returns sub-slices of the input buffer; error paths use pre-allocated sentinel errors. The `BenchmarkReader_AllTags` benchmark and `TestReader_NoAllocs` enforce this.
- **Minimal-allocation writes.** Numeric encoders use `binary.BigEndian.AppendUint*` (no temporary slices). With a pre-sized `Buff`, the writer also reports `0 allocs/op`.

## Commands

```sh
go build ./...                                    # build
go test ./...                                     # run all tests
go test -run TestReader_Read_Int ./...            # run a single test
go test -run TestReader_Read_Int/int16_max ./...  # run a single subtest
go test -race -cover ./...                        # race detector + coverage
go test -bench . -benchmem -run=^$ ./...          # benchmarks with alloc tracking
```

## Architecture

Three source files form the public API:

- **`types.go`** — `Type` (a `byte`) and all MessagePack format tag constants. Range constants like `PosFixIntMax`, `FixMapMax`, `FixArrayMax`, `FixStrMax`, `NegFixInt`/`NegFixIntMax` are used for *range checks* against the leading byte; do not treat them as standalone tags.
- **`msgp_reader.go`** — `MsgpReader{Buff, Idx}` and `IMsgpReader`. The single `Read()` returns `(Type, int, []byte, error)`:
  - **Scalars** (`Int*`/`Uint*`/`Float*`): payload is the raw big-endian sub-slice (a view, not a copy). `int` return is `0`.
  - **Tag-encoded values** (`Nil`, `True`, `False`, `PosFixInt`, `NegFixInt`): payload is `nil`. The tag byte itself carries the value.
  - **Strings** (`FixStr`, `Str8/16/32`) and **bin** (`Bin8/16/32`): payload is the raw bytes.
  - **Ext** (`FixExt1/2/4/8/16`, `Ext8/16/32`): payload is `[extType, dataBytes...]`. Caller does `data[0]` for the ext type byte and `data[1:]` for the payload.
  - **Collections** (`FixArray`/`Array16/32`, `FixMap`/`Map16/32`): the `int` return is the element count (pair count for maps); `[]byte` is the *remainder* of `Buff` after the header. The caller calls `Read()` repeatedly to consume the elements (2× count for maps). `Idx` is advanced past the header but not past the children — recursion is caller-driven and zero-alloc.
  - **Errors:** `io.EOF` (re-exported as `EOF`) on exhausted buffer; `ErrTruncated` when input ends mid-value; `ErrUnknownType` for reserved/invalid tag bytes. Sentinels are pre-allocated so the error path stays alloc-free.
  - `Skip()` is `Read()` with the return values discarded.
- **`msgp_writer.go`** — `MsgpWriter{Buff}` and `IMsgpWriter`. Two layers of write methods:
  - **Auto-sized** (`WriteString`/`WriteBytes`/`WriteArray`/`WriteMap`/`WriteExt`): pick the smallest spec format that fits the input. Return a sentinel error if the input exceeds the largest variant.
  - **Explicit** (`WriteFixStr`/`WriteStr8`/`WriteStr16`/`WriteStr32`, `WriteBin8/16/32`, `WriteFixArray`/`WriteArray16`/`WriteArray32`, `WriteFixMap`/`WriteMap16`/`WriteMap32`, `WriteFixExt1/2/4/8/16`, `WriteExt8/16/32`): pin the wire format. Return a range error if the input doesn't fit.
  - `WriteInt`/`WriteUint` always emit `Int64`/`Uint64` (no auto-compaction); use the explicit `WritePosFixInt`/`WriteNegFixInt`/`WriteInt8..64` for compact forms.
  - `WriteNegFixInt` accepts only `-32..-1` per spec.

### Conventions to preserve

- **Reader payloads alias `Buff`.** They are sub-slices, not copies. Don't write to them; copy if the caller needs ownership.
- **Big-endian everywhere.** All multi-byte integers and floats follow the MessagePack spec. Use `encoding/binary.BigEndian` and `math.Float{32,64}bits`/`Float{32,64}frombits`.
- **No `fmt.Errorf` in hot paths.** Errors are package-level sentinels (`errors.New`). When adding new error conditions, declare a sentinel; don't allocate.
- **No temporary slices in numeric writers.** Use `binary.BigEndian.AppendUint*` (Go 1.19+) directly on `Buff`.
- **Zero-alloc test harness.** `allocs_test.go` runs every reader path through `testing.AllocsPerRun` and asserts `0` allocations. Adding a new tag means extending `allTagsFixture`.
- **Tests are table-driven** with `testify/assert`/`require`. Reader/writer roundtrips live in `msgp_writer_test.go`; per-tag coverage in `msgp_*_coverage_test.go`.

### Files

- `msgp_reader.go` / `msgp_reader_test.go` / `msgp_reader_coverage_test.go`
- `msgp_writer.go` / `msgp_writer_test.go` / `msgp_writer_coverage_test.go`
- `types.go`
- `allocs_test.go` — zero-alloc assertion + cross-cutting benchmarks
