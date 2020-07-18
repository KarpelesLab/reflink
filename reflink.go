package reflink

import (
	"os"
	"syscall"
)

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
