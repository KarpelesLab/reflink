package reflink

import "errors"

// https://github.com/torvalds/linux/blob/v5.2/include/uapi/linux/fs.h#L195
// #define FICLONE              _IOW(0x94, 9, int)
//
// printf("%lx", FICLONE) → 0x40049409

const FICLONE = 0x40049409

var ErrReflinkUnsupported = errors.New("reflink is not supported on this OS")
