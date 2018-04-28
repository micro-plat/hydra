package mqtt

import (
	"io"
)

// Payload is the interface for Publish payloads. Typically the BytesPayload
// implementation will be sufficient for small payloads whose full contents
// will exist in memory. However, other implementations can read or write
// payloads requiring them holding their complete contents in memory.
type Payload interface {
	// Size returns the number of bytes that WritePayload will write.
	Size() int

	// WritePayload writes the payload data to w. Implementations must write
	// Size() bytes of data, but it is *not* required to do so prior to
	// returning. Size() bytes must have been written to w prior to another
	// message being encoded to the underlying connection.
	WritePayload(w io.Writer) error

	// ReadPayload reads the payload data from r (r will EOF at the end of the
	// payload). It is *not* required for r to have been consumed prior to this
	// returning. r must have been consumed completely prior to another message
	// being decoded from the underlying connection.
	ReadPayload(r io.Reader) error
}

// BytesPayload reads/writes a plain slice of bytes.
type BytesPayload []byte

func (p BytesPayload) Size() int {
	return len(p)
}

func (p BytesPayload) WritePayload(w io.Writer) error {
	_, err := w.Write(p)
	return err
}

func (p BytesPayload) ReadPayload(r io.Reader) error {
	_, err := io.ReadFull(r, p)
	return err
}

// StreamedPayload writes payload data from reader, or reads payload data into a writer.
type StreamedPayload struct {
	// N indicates payload size to the encoder. This many bytes will be read from
	// the reader when encoding. The number of bytes in the payload will be
	// stored here when decoding.
	N int

	// EncodingSource is used to copy data from when encoding a Publish message
	// onto the wire. This can be
	EncodingSource io.Reader

	// DecodingSink is used to copy data to when decoding a Publish message from
	// the wire. This can be nil if the payload is only being used for encoding.
	DecodingSink io.Writer
}

func (p *StreamedPayload) Size() int {
	return p.N
}

func (p *StreamedPayload) WritePayload(w io.Writer) error {
	_, err := io.CopyN(w, p.EncodingSource, int64(p.N))
	return err
}

func (p *StreamedPayload) ReadPayload(r io.Reader) error {
	n, err := io.Copy(p.DecodingSink, r)
	p.N = int(n)
	return err
}
