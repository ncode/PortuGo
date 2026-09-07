package runtime

import "errors"

// MaxTextBytes is the project safeguard for one decoded text value.
const MaxTextBytes = 16 << 20

// ErrTextSize reports a rejected allocation, before growing a text value.
var ErrTextSize = errors.New("text size limit exceeded")
