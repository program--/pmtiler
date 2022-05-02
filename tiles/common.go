package tiles

import "encoding/binary"

func readUint24(b []byte) uint32 {
	return (uint32(binary.LittleEndian.Uint16(b[1:3])) << 8) + uint32(b[0])
}

func writeUint24(i uint32) []byte {
	result := make([]byte, 3)
	binary.LittleEndian.PutUint16(result[1:3], uint16(i>>8&0xFFFF))
	result[0] = uint8(i & 0xFF)
	return result
}

func readUint48(b []byte) uint64 {
	return (uint64(binary.LittleEndian.Uint32(b[2:6])) << 16) + uint64(uint32(binary.LittleEndian.Uint16(b[0:2])))
}

func writeUint48(i uint64) []byte {
	result := make([]byte, 6)
	binary.LittleEndian.PutUint32(result[2:6], uint32(i>>16&0xFFFFFFFF))
	binary.LittleEndian.PutUint16(result[0:2], uint16(i&0xFFFF))
	return result
}
