package proxy

import (
	"io"
	"net/http"
)

type flushWriter struct {
	flusher http.Flusher
	writer  io.Writer
}

func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.writer.Write(p)
	if fw.flusher != nil {
		fw.flusher.Flush()
	}
	return
}

func MakeFlushWriter(writer io.Writer) io.Writer {
	if flusher, ok := writer.(http.Flusher); ok {
		return &flushWriter{
			writer:  writer,
			flusher: flusher,
		}
	}

	return writer
}
