/*
Package drive provides the definitions for the low-level pfs storage drivers.
*/
package drive

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/pachyderm/pachyderm/src/client/pfs"
)

// ValidateRepoName determines if a repo name is valid
func ValidateRepoName(name string) error {
	match, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", name)

	if !match {
		return fmt.Errorf("repo name (%v) invalid: only alphanumeric characters, underscores, and dashes are allowed", name)
	}

	return nil
}

// ListFileMode specifies how ListFile executes.
type ListFileMode int

const (
	// ListFileNORMAL computes sizes for files but not for directories
	ListFileNORMAL ListFileMode = iota
	// ListFileFAST does not compute sizes for files or directories
	ListFileFAST
	// ListFileRECURSE computes sizes for files and directories
	ListFileRECURSE
)

// IsPermissionError returns true if a given error is a permission error.
func IsPermissionError(err error) bool {
	return strings.Contains(err.Error(), "has already finished")
}

// CommitEvent is an event that contains a CommitInfo or an error
type CommitEvent struct {
	Err   error
	Value *pfs.CommitInfo
}

// CommitStream is a stream of CommitInfos
type CommitStream interface {
	Stream() <-chan CommitEvent
	Close()
}

// Driver represents a low-level pfs storage driver.
type Driver interface {
	CreateRepo(ctx context.Context, repo *pfs.Repo, provenance []*pfs.Repo) error
	InspectRepo(ctx context.Context, repo *pfs.Repo) (*pfs.RepoInfo, error)
	ListRepo(ctx context.Context, provenance []*pfs.Repo) ([]*pfs.RepoInfo, error)
	DeleteRepo(ctx context.Context, repo *pfs.Repo, force bool) error

	StartCommit(ctx context.Context, parent *pfs.Commit, branch string, provenance []*pfs.Commit) (*pfs.Commit, error)
	BuildCommit(ctx context.Context, parent *pfs.Commit, branch string, provenance []*pfs.Commit, tree *pfs.Object) (*pfs.Commit, error)
	FinishCommit(ctx context.Context, commit *pfs.Commit) error
	InspectCommit(ctx context.Context, commit *pfs.Commit) (*pfs.CommitInfo, error)

	ListCommit(ctx context.Context, repo *pfs.Repo, from *pfs.Commit, to *pfs.Commit, number uint64) ([]*pfs.CommitInfo, error)
	SubscribeCommit(ctx context.Context, repo *pfs.Repo, branch string, from *pfs.Commit) (CommitStream, error)
	FlushCommit(ctx context.Context, fromCommits []*pfs.Commit, toRepos []*pfs.Repo) (CommitStream, error)
	DeleteCommit(ctx context.Context, commit *pfs.Commit) error

	ListBranch(ctx context.Context, repo *pfs.Repo) ([]*pfs.Branch, error)
	SetBranch(ctx context.Context, commit *pfs.Commit, name string) error
	DeleteBranch(ctx context.Context, repo *pfs.Repo, name string) error

	PutFile(ctx context.Context, file *pfs.File, delimiter pfs.Delimiter,
		targetFileDatums int64, targetFileBytes int64, reader io.Reader) error
	MakeDirectory(ctx context.Context, file *pfs.File) error
	GetFile(ctx context.Context, file *pfs.File, offset int64, size int64) (io.Reader, error)
	InspectFile(ctx context.Context, file *pfs.File) (*pfs.FileInfo, error)
	ListFile(ctx context.Context, file *pfs.File) ([]*pfs.FileInfo, error)
	GlobFile(ctx context.Context, commit *pfs.Commit, pattern string) ([]*pfs.FileInfo, error)
	DeleteFile(ctx context.Context, file *pfs.File) error

	DeleteAll(ctx context.Context) error
	Dump(ctx context.Context)
}
