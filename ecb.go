package utools

import "crypto/cipher"

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypt ecb

// NewECBEncrypted returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypted(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypt)(newECB(b))
}
func (x *ecbEncrypt) BlockSize() int { return x.blockSize }
func (x *ecbEncrypt) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypt ecb

// NewECBDecrypted returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypted(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypt)(newECB(b))
}
func (x *ecbDecrypt) BlockSize() int { return x.blockSize }
func (x *ecbDecrypt) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
