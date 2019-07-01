package crc32

import "hash/crc32"

func Encrypt(buffer []byte) uint32 {
	return crc32.ChecksumIEEE(buffer)
}
