package security

type IEncrypt interface {
	Encrypt(input []byte) string
}

type ConfEncrypt struct {
	EnableEncryption bool
}

func (c ConfEncrypt) Encrypt(input []byte) string {
	if c.EnableEncryption {
		return Encrypt(input)
	}
	return string(input)
}
