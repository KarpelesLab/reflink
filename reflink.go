//+build !linux

package reflink

import "os"

func reflinkInternal(d, s *os.File) error {
	return ErrReflinkUnsupported
}

func reflinkRangeInternal(dst, src *os.File, dstOffset, srcOffset, n int64) error {
	return ErrReflinkUnsupported
}
