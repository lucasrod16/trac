package commit

import "errors"

var (
	ErrEmptyCommitHash      = errors.New("commit hash is empty")
	ErrWorkingTreeClean     = errors.New("nothing to commit, working tree clean")
	ErrNothingAddedToCommit = errors.New(`nothing added to commit (use "trac add" to track)`)
)
