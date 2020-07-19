package reflink

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Always will perform a reflink operation and fail on error
func Always(src, dst string) error {
	return reflinkFile(src, dst, false)
}

// Auto will attempt to perform a reflink operation and fallback to normal data
// copy if reflink is not supported.
func Auto(src, dst string) error {
	return reflinkFile(src, dst, true)
}

// reflinkFile perform the reflink operation in order to copy src into dst using
// the underlying filesystem's copy-on-write reflink system. If this fails (for
// example the filesystem does not support reflink) and fallback is true, then
// io.Copy will be used to copy the data.
func reflinkFile(src, dst string, fallback bool) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	// generate temporary file for output
	tmp, err := ioutil.TempFile(filepath.Dir(dst), "")
	if err != nil {
		return err
	}

	// copy to temp file
	err = reflinkInternal(tmp, s)
	if (err != nil) && fallback {
		_, err = io.Copy(tmp, s)
	}
	tmp.Close() // we're not writing to this anymore

	// if an error happened, remove temp file and signal error
	if err != nil {
		os.Remove(tmp.Name())
		return err
	}

	// keep src file mode if possible
	if st, err := s.Stat(); err == nil {
		tmp.Chmod(st.Mode())
	}

	// replace dst file
	err = os.Rename(tmp.Name(), dst)
	if err != nil {
		// failed to rename (dst is not writable?)
		os.Remove(tmp.Name())
		return err
	}

	return nil
}

// FReflink performs the reflink operation on the passed files, replacing
// dst's contents with src. If fallback is true and reflink fails, io.Copy will
// be used to copy the data.
//
// In case of fallback, seek position in src and dst will be affected.
func Reflink(dst, src *os.File, fallback bool) error {
	err := reflinkInternal(dst, src)
	if (err != nil) && fallback {
		// seek both src & dst at beginning
		src.Seek(0, io.SeekStart)
		dst.Seek(0, io.SeekStart)
		dst.Truncate(0) // assuming any error in trucate will result in copy error
		_, err = io.Copy(dst, src)
	}
	return err
}

// Partial performs a range reflink operation on the passed files, replacing
// part of dst's contents with data from src. If fallback is true and reflink
// fails, io.CopyN will be used to copy the data.
//
// In case of fallback, seek position in src and dst will be affected.
func Partial(dst, src *os.File, dstOffset, srcOffset, n int64, fallback bool) error {
	err := reflinkRangeInternal(dst, src, dstOffset, srcOffset, n)
	if (err != nil) && fallback {
		// seek both src & dst
		src.Seek(srcOffset, io.SeekStart)
		dst.Seek(dstOffset, io.SeekStart)
		_, err = io.CopyN(dst, src, n)
	}
	return err
}
