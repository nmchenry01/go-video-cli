package upload

import (
	"os"
	"sync/atomic"

	"github.com/cheggaaa/pb"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Custom reader encountered an error")
	}
	atomic.AddInt64(&r.read, int64(n))

	readInMbs := convertBytesToMb(int(r.read / 2))

	r.progress.Set(readInMbs)

	return n, err
}

func (r *customReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}

func convertBytesToMb(bytes int) int {
	return bytes / 1024 / 1024
}
