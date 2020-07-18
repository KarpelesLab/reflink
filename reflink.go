//+build !linux

package reflink

import "os"

func reflinkInternal(d, s *os.File) error {
	return ErrReflinkUnsupported
}
