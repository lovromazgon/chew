package chew

import (
	"io"
	"os"
	"sync"
)

// Writer extends the io.Writer and adds the option to set the output filename.
//
// When processing Chewable the SetOut method will be called before each template execution
// informing the writer of the new target file. When writing to standard output this information
// is meaningless.
type Writer interface {
	io.Writer
	// SetOut sets the output filename for the current template
	SetOut(filename string)
}


// WriterWrapper is a convenience object to allow wrapping an io.Writer and implement the Writer interface.
// The SetOut method is empty and does nothing.
type WriterWrapper struct {
	io.Writer
}
// SetOut is empty and does nothing.
func (WriterWrapper) SetOut(filename string) {}


// MultiFileWriter is a Writer which writes everything to files in the folder Out. If the folder defined in Out
// doesn't exist it is created before writing to the first file. SetOut has to be called before starting to
// write to this Writer, so that the file is created or truncated before writing to it.
type MultiFileWriter struct {
	*os.File
	sync.Once
	Out string
}

// SetOut sets the filename of the output file into which the succeeding calls to Write will output the content.
// A file with the provided filename will be created in the folder defined in Out. If the file exists, it will
// be truncated.
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