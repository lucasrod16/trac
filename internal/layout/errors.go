package layout

import "errors"

var ErrNotTracRepository = errors.New("not a trac repository (or any of the parent directories): .trac")
