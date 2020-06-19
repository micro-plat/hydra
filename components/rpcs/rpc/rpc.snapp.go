package rpc

import (
	"io"
	"sync"

	"github.com/golang/snappy"
	"google.golang.org/grpc/encoding"
)

//Snappy is the name registered for the snappy Compressor.
const Snappy = "snappy"

func init() {
	encoding.RegisterCompressor(&Compressor{})
}

var (
	// cmpPool stores writers
	cmpPool sync.Pool
	// dcmpPool stores readers
	dcmpPool sync.Pool
)

type Compressor struct {
}

func (c *Compressor) Name() string {
	return Snappy
}

func (c *Compressor) Compress(w io.Writer) (io.WriteCloser, error) {
	wr, inPool := cmpPool.Get().(*writeCloser)
	if !inPool {
		return &writeCloser{Writer: snappy.NewBufferedWriter(w)}, nil
	}
	wr.Reset(w)

	return wr, nil
}

func (c *Compressor) Decompress(r io.Reader) (io.Reader, error) {
	dr, inPool := dcmpPool.Get().(*reader)
	if !inPool {
		return &reader{Reader: snappy.NewReader(r)}, nil
	}
	dr.Reset(r)

	return dr, nil
}

type writeCloser struct {
	*snappy.Writer
}

func (w *writeCloser) Close() error {
	defer func() {
		cmpPool.Put(w)
	}()

	return w.Writer.Close()
}

type reader struct {
	*snappy.Reader
}

func (r *reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	if err == io.EOF {
		dcmpPool.Put(r)
	}

	return n, err
}
