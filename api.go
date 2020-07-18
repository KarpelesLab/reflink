package reflink

import (
	"io"
	"os"
)

// Always will perform a reflink operation and fail on error
func Always(src, dst string) error {
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

	// keep file mode if possible
	if st, err := s.Stat(); err == nil {
		d.Chmod(st.Mode())
	}

	return reflinkInternal(d, s)
}

// Auto will attempt to perform a reflink operation and fallback to normal data
// copy if reflink is not supported.
func Auto(src, dst string) error {
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

	// keep file mode if possible
	if st, err := s.Stat(); err == nil {
		d.Chmod(st.Mode())
	}

	err = reflinkInternal(d, s)
	if err != nil {
		_, err = io.Copy(d, s)
	}

	return err
}
