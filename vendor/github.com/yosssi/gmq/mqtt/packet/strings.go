package packet

// Maximum length of the UTF-8 encoded strings
const maxStringsLen = 65535

// appendLenStr appends the length of the strings
// and the strings to the byte slice.
func appendLenStr(b []byte, s []byte) []byte {
	b = append(b, encodeUint16(uint16(len(s)))...)
	b = append(b, s...)
	return b
}
