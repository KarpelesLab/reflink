package reflink

import (
	"io"
	"os"
)

// Always will perform a reflink operation and fail on error
func Always(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	return reflinkInternal(d, s)
}

// Auto will attempt to perform a reflink oepration and fallback to normal data
// copy if reflink is not supported.
func Auto(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	err = reflinkInternal(d, s)
	if err != nil {
		_, err = io.Copy(d, s)
	}

	return err
}
