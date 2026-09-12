package runtime

import "errors"

// MaxTextBytes is the project safeguard for one decoded text value.
const MaxTextBytes = 16 << 20

// MaxVectorSlots bounds the scalar cells in one aggregate as a project safeguard.
const MaxVectorSlots = 1 << 20

// ErrTextSize reports a rejected allocation, before growing a text value.
var ErrTextSize = errors.New("text size limit exceeded")

// ErrVectorSize reports a rejected aggregate allocation before creating cells.
var ErrVectorSize = errors.New("vector scalar-slot limit exceeded")
