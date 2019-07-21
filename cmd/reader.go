package cmd

import (
	"os"
	"sync/atomic"

	"github.com/cheggaaa/pb"
)

type customReader struct {
	fp       *os.File
	size     int64
	read     int64
	progress *pb.ProgressBar
}

func (r *customReader) Read(p []byte) (int, error) {
	return r.fp.Read(p)
}

func (r *customReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := r.fp.ReadAt(p, off)
	check(err, "There was a problem while reading the file")

	atomic.AddInt64(&r.read, int64(n))

	readInMbs := convertBytesToMb(int(r.read / 2))

	r.progress.Set(readInMbs)

	return n, err
}

func (r *customReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}
