package chew

import (
	"io"
	"os"
	"sync"
)

type Writer interface {
	io.Writer
	// sets the output filename
	SetOut(filename string)
}


type WriterWrapper struct {
	io.Writer
}
func (WriterWrapper) SetOut(filename string) {}


type MultiFileWriter struct {
	*os.File
	sync.Once
	Out string
}

func (w *MultiFileWriter) SetOut(filename string) {
	w.Once.Do(func() {
		// check if out folder exists and create it
		if _,err := os.Stat(w.Out); err != nil {
			if err := os.MkdirAll(w.Out, os.ModePerm); err != nil {
				panic(err)
			}
		}
	})

	var err error
	w.File, err = os.Create(w.Out + "/" + filename)
	if err != nil {
		panic(err)
	}
}