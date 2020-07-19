//+build linux

package reflink

import (
	"os"
	"syscall"
	"unsafe"
)

// FICLONE is a constant from the Linux kernel include linux/fs.h, not found
// in the unix or syscall packages.
//
// https://github.com/torvalds/linux/blob/v5.2/include/uapi/linux/fs.h#L195
// #define FICLONE              _IOW(0x94, 9, int)
// #define FICLONERANGE	_IOW(0x94, 13, struct file_clone_range)
// #define FIDEDUPERANGE	_IOWR(0x94, 54, struct file_dedupe_range)
const (
	FICLONE       = 0x40049409
	FICLONERANGE  = 0x4020940d
	FIDEDUPERANGE = 0x40209436
)

// https://github.com/torvalds/linux/blob/v5.2/include/uapi/linux/fs.h#L51
type fsFileCloneRange struct {
	srcFd      int64
	srcOffset  uint64
	srcLength  uint64
	destOffset uint64
}

// reflinkInternal performs the actual reflink action without worrying about fallback
func reflinkInternal(d, s *os.File) error {
	ss, err := s.SyscallConn()
	if err != nil {
		return err
	}
	sd, err := d.SyscallConn()
	if err != nil {
		return err
	}

	var err2, err3 error

	err = sd.Control(func(dfd uintptr) {
		err2 = ss.Control(func(sfd uintptr) {
			// int ioctl(int dest_fd, FICLONE, int src_fd);
			_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, dfd, FICLONE, sfd)
			if errno != 0 {
				err3 = errno
			}
		})
	})

	if err != nil {
		// sd.Control failed
		return err
	}
	if err2 != nil {
		// ss.Control failed
		return err2
	}

	// err3 is ioctl() response
	return err3
}

func reflinkRangeInternal(dst, src *os.File, dstOffset, srcOffset, n int64) error {
	ss, err := src.SyscallConn()
	if err != nil {
		return err
	}
	sd, err := dst.SyscallConn()
	if err != nil {
		return err
	}

	var err2, err3 error

	err = sd.Control(func(dfd uintptr) {
		err2 = ss.Control(func(sfd uintptr) {
			req := &fsFileCloneRange{
				srcFd:      int64(sfd),
				srcOffset:  uint64(srcOffset),
				srcLength:  uint64(n),
				destOffset: uint64(dstOffset),
			}

			// int ioctl(int dest_fd, FICLONE, int src_fd);
			_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, dfd, FICLONERANGE, uintptr(unsafe.Pointer(&req)))
			if errno != 0 {
				err3 = errno
			}
		})
	})

	if err != nil {
		// sd.Control failed
		return err
	}
	if err2 != nil {
		// ss.Control failed
		return err2
	}

	// err3 is ioctl() response
	return err3
}
