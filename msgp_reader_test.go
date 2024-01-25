package msgpraw

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader_Read_Int(t *testing.T) {
	testCase := []struct {
		name    string
		msgType Type
		msg     []byte
	}{
		{
			name:    "int8_min",
			msgType: Int8,
			msg: []byte{
				0xd0, // int8
				0x00, // 0
			},
		},
		{
			name:    "int8_max",
			msgType: Int8,
			msg: []byte{
				0xd0, // int8
				0xff, // 255
			},
		},
		{
			name:    "int16_min",
			msgType: Int16,
			msg: []byte{
				0xd1, // int16
				0x00, // 0
				0x00, // 0
			},
		},
		{
			name:    "int16_max",
			msgType: Int16,
			msg: []byte{
				0xd1, // int16
				0xff, // 255
				0xff, // 255
			},
		},
		{
			name:    "int32_min",
			msgType: Int32,
			msg: []byte{
				0xd2, // int32
				0x00, // 0
				0x00, // 0
				0x00, // 0
				0x00, // 0
			},
		},
		{
			name:    "int32_max",
			msgType: Int32,
			msg: []byte{
				0xd2, // int32
				0xff, // 255
				0xff, // 255
				0xff, // 255
				0xff, // 255
			},
		},
		{
			name:    "int64_min",
			msgType: Int64,
			msg: []byte{
				0xd3, // int64
				0x00, // 0
				0x00, // 0
				0x00, // 0
				0x00, // 0
				0x00, // 0
				0x00, // 0
				0x00, // 0
				0x00, // 0
			},
		},
		{
			name:    "int64_max",
			msgType: Int64,
			msg: []byte{
				0xd3, // int64
				0xff, // 255
				0xff, // 255
				0xff, // 255
				0xff, // 255
				0xff, // 255
				0xff, // 255
				0xff, // 255
				0xff, // 255
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			reader := &MsgpReader{Buffer: tc.msg}

			msgType, _, value, err := reader.Read()
			require.NoError(t, err)
			assert.Equal(t, tc.msgType, msgType)

			if tc.msgType == Int8 {
				assert.Equal(t, []byte{tc.msg[1]}, value)
			} else if tc.msgType == Int16 {
				assert.Equal(t, tc.msg[1:3], value)
			} else if tc.msgType == Int32 {
				assert.Equal(t, tc.msg[1:5], value)
			} else if tc.msgType == Int64 {
				assert.Equal(t, tc.msg[1:9], value)
			} else {
				assert.Failf(t, "unexpected msgType", "msgType: %d", tc.msgType)
			}
		})
	}
}

func TestReader_Read_FixArray(t *testing.T) {
	testCase := []struct {
		name   string
		length int
		msg    []byte
	}{
		{
			name:   "fixarray(1)",
			length: 1,
			msg: []byte{
				0x91,      // fixarray(1)
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(2)",
			length: 2,
			msg: []byte{
				0x92,      // fixarray(2)
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(3)",
			length: 3,
			msg: []byte{
				0x93,      // fixarray(3)
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(4)",
			length: 4,
			msg: []byte{
				0x94,      // fixarray(4)
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(5)",
			length: 5,
			msg: []byte{
				0x95,      // fixarray(5)
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(6)",
			length: 6,
			msg: []byte{
				0x96,      // fixarray(6)
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(7)",
			length: 7,
			msg: []byte{
				0x97,      // fixarray(7)
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(8)",
			length: 8,
			msg: []byte{
				0x98,      // fixarray(8)
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(9)",
			length: 9,
			msg: []byte{
				0x99,      // fixarray(9)
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(10)",
			length: 10,
			msg: []byte{
				0x9a,      // fixarray(10)
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(11)",
			length: 11,
			msg: []byte{
				0x9b,      // fixarray(11)
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(12)",
			length: 12,
			msg: []byte{
				0x9c,      // fixarray(12)
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(13)",
			length: 13,
			msg: []byte{
				0x9d,      // fixarray(13)
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(14)",
			length: 14,
			msg: []byte{
				0x9e,      // fixarray(14)
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
		{
			name:   "fixarray(15)",
			length: 15,
			msg: []byte{
				0x9f,      // fixarray(15)
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
				byte(Nil), // nil
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			reader := &MsgpReader{Buffer: tc.msg}

			msgType, numberOfElements, _, err := reader.Read()
			require.NoError(t, err)
			assert.Equal(t, FixArray, msgType&0xf0)

			for i := 0; i < numberOfElements; i++ {
				msgType, _, value, err := reader.Read()
				if errors.Is(err, EOF) {
					break
				}

				require.NoError(t, err)
				assert.Equal(t, Nil, msgType)
				assert.Nil(t, value)
			}

			assert.Equal(t, tc.length, numberOfElements)
		})
	}
}

func TestReader_Read_FixMap(t *testing.T) {
	testCase := []struct {
		name   string
		length int
		msg    []byte
	}{
		{
			name:   "fixmap(1)",
			length: 1,
			msg: []byte{
				0x81,                // fixmap(1)
				byte(PosFixInt + 1), // nil
				byte(Nil),           // nil
			},
		},
		{
			name:   "fixmap(2)",
			length: 2,
			msg: []byte{
				0x82,                // fixmap(2)
				byte(PosFixInt + 1), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 2), // nil
				byte(Nil),           // nil
			},
		},
		{
			name:   "fixmap(3)",
			length: 3,
			msg: []byte{
				0x83,                // fixmap(3)
				byte(PosFixInt + 1), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 2), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 3), // nil
				byte(Nil),           // nil
			},
		},
		{
			name:   "fixmap(4)",
			length: 4,
			msg: []byte{
				0x84,                // fixmap(4)
				byte(PosFixInt + 1), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 2), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 3), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 4), // nil
				byte(Nil),           // nil
			},
		},
		{
			name:   "fixmap(5)",
			length: 5,
			msg: []byte{
				0x85,                // fixmap(5)
				byte(PosFixInt + 1), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 2), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 3), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 4), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 5), // nil
				byte(Nil),           // nil
			},
		},
		{
			name:   "fixmap(6)",
			length: 6,
			msg: []byte{
				0x86,                // fixmap(6)
				byte(PosFixInt + 1), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 2), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 3), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 4), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 5), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 6), // nil
				byte(Nil),           // nil
			},
		},
		{
			name:   "fixmap(7)",
			length: 7,
			msg: []byte{
				0x87,                // fixmap(7)
				byte(PosFixInt + 1), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 2), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 3), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 4), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 5), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 6), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 7), // nil
				byte(Nil),           // nil
			},
		},
		{
			name:   "fixmap(8)",
			length: 8,
			msg: []byte{
				0x88,                // fixmap(8)
				byte(PosFixInt + 1), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 2), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 3), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 4), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 5), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 6), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 7), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 8), // nil
				byte(Nil),           // nil
			},
		},
		{
			name:   "fixmap(9)",
			length: 9,
			msg: []byte{
				0x89,                // fixmap(9)
				byte(PosFixInt + 1), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 2), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 3), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 4), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 5), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 6), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 7), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 8), // nil
				byte(Nil),           // nil
				byte(PosFixInt + 9), // nil
				byte(Nil),           // nil
			},
		},
		{
			name:   "fixmap(10)",
			length: 10,
			msg: []byte{
				0x8a,                 // fixmap(10)
				byte(PosFixInt + 1),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 2),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 3),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 4),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 5),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 6),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 7),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 8),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 9),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 10), // nil
				byte(Nil),            // nil
			},
		},
		{
			name:   "fixmap(11)",
			length: 11,
			msg: []byte{
				0x8b,                 // fixmap(11)
				byte(PosFixInt + 1),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 2),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 3),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 4),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 5),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 6),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 7),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 8),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 9),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 10), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 11), // nil
				byte(Nil),            // nil
			},
		},
		{
			name:   "fixmap(12)",
			length: 12,
			msg: []byte{
				0x8c,                 // fixmap(12)
				byte(PosFixInt + 1),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 2),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 3),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 4),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 5),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 6),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 7),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 8),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 9),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 10), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 11), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 12), // nil
				byte(Nil),            // nil
			},
		},
		{
			name:   "fixmap(13)",
			length: 13,
			msg: []byte{
				0x8d,                 // fixmap(13)
				byte(PosFixInt + 1),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 2),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 3),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 4),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 5),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 6),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 7),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 8),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 9),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 10), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 11), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 12), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 13), // nil
				byte(Nil),            // nil
			},
		},
		{
			name:   "fixmap(14)",
			length: 14,
			msg: []byte{
				0x8e,                 // fixmap(14)
				byte(PosFixInt + 1),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 2),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 3),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 4),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 5),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 6),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 7),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 8),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 9),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 10), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 11), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 12), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 13), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 14), // nil
				byte(Nil),            // nil
			},
		},
		{
			name:   "fixmap(15)",
			length: 15,
			msg: []byte{
				0x8f,                 // fixmap(15)
				byte(PosFixInt + 1),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 2),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 3),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 4),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 5),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 6),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 7),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 8),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 9),  // nil
				byte(Nil),            // nil
				byte(PosFixInt + 10), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 11), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 12), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 13), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 14), // nil
				byte(Nil),            // nil
				byte(PosFixInt + 15), // nil
				byte(Nil),            // nil
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			reader := &MsgpReader{Buffer: tc.msg}

			msgType, numberOfElements, _, err := reader.Read()
			require.NoError(t, err)
			assert.Equal(t, FixMap, msgType&0xf0)

			for i := 0; i < numberOfElements; i++ {
				msgType, _, _, err = reader.Read()
				if errors.Is(err, EOF) {
					break
				}

				require.NoError(t, err)
				assert.Equal(t, PosFixInt, msgType&0xf0)
				assert.Equal(t, byte(i+1), byte(msgType))

				msgType, _, mapValue, err := reader.Read()
				if errors.Is(err, EOF) {
					break
				}

				require.NoError(t, err)
				assert.Equal(t, Nil, msgType)
				assert.Nil(t, mapValue)
			}

			assert.Equal(t, tc.length, numberOfElements)
		})
	}
}
