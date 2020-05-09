package des

import (
	"crypto/cipher"
)

type Ecb struct {
	B         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *Ecb {
	return &Ecb{
		B:         b,
		blockSize: b.BlockSize(),
	}
}

type EcbEncrypter Ecb

func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*EcbEncrypter)(newECB(b))
}

func (e *EcbEncrypter) BlockSize() int {
	return e.blockSize
}

func (e *EcbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%e.blockSize != 0 {
		panic("aesecb: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("aesecb: output smaller than input")
	}
	for len(src) > 0 {
		e.B.Encrypt(dst, src[:e.blockSize])
		src = src[e.blockSize:]
		dst = dst[e.blockSize:]
	}
}

type EcbDecrypter Ecb

func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*EcbDecrypter)(newECB(b))
}

func (e *EcbDecrypter) BlockSize() int {
	return e.blockSize
}

func (e *EcbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%e.blockSize != 0 {
		panic("aesecb: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("aesecb: output smaller than input")
	}
	for len(src) > 0 {
		e.B.Decrypt(dst, src[:e.blockSize])
		src = src[e.blockSize:]
		dst = dst[e.blockSize:]
	}
}
