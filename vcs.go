package apicompat

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/Masterminds/vcs"
)

const revisionFileSystem = "."

type VCS struct {
	vcs.Repo
	cachedVersion string
}

func NewLocalVCS(local string) (*VCS, error) {
	repo, err := vcs.NewRepo("", local)
	if err != nil {
		return nil, err
	}
	if !repo.CheckLocal() {
		return nil, errors.New("directory is not a repository")
	}
	if !repo.Ping() {
		return nil, errors.New("cannot ping remote repository")
	}
	return &VCS{
		Repo: repo,
	}, nil
}

func NewVCS(remoteURL string) (*VCS, error) {
	local, err := ioutil.TempDir("", "apicompat")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Creating new repo %s at %s\n", remoteURL, local)
	repo, err := vcs.NewRepo(remoteURL, local)
	if err != nil {
		return nil, err
	}
	fmt.Println("Pinging repository...")
	if !repo.Ping() {
		return nil, errors.New("Cannot ping remote repository")
	}

	fmt.Println("Cloning repository...")
	if err := repo.Get(); err != nil {
		return nil, err
	}
	return &VCS{
		Repo: repo,
	}, nil
}

func (v *VCS) ReadDir(revision, path string) ([]os.FileInfo, error) {
	if err := v.ensureVersion(revision); err != nil {
		return nil, err
	}
	return ioutil.ReadDir(path)
}

func (v *VCS) OpenFile(revision, path string) (io.ReadCloser, error) {
	if err := v.ensureVersion(revision); err != nil {
		return nil, err
	}
	return os.Open(path)
}

func (v *VCS) DefaultRevision() (before, after string, err error) {
	// Before should be the 'latest' semver tag
	// After defaults to the special 'filesystem' value to check the local working copy.
	before, err = v.latestSemVerTag()
	if err != nil {
		return
	}
	after = revisionFileSystem
	return
}

func (v *VCS) latestSemVerTag() (string, error) {
	tags, err := v.Repo.Tags()
	if err != nil {
		return "", err
	}
	if len(tags) == 0 {
		return "", errors.New("0 tags detected")
	}
	var vs []*semver.Version
	for _, t := range tags {
		// discard errors as we don't expect all tags to be semvers
		v, err := semver.NewVersion(t)
		if err == nil {
			vs = append(vs, v)
		}
	}
	sort.Sort(semver.Collection(vs))

	// Return original string since we'll need this to checkout
	return vs[0].Original(), nil
}

func (v *VCS) ensureVersion(ver string) error {
	if v.cachedVersion == ver {
		return nil
	}
	repoVer, err := v.Repo.Version()
	if err != nil {
		return err
	}
	if repoVer != ver {
		if err := v.Repo.UpdateVersion(ver); err != nil {
			return err
		}
	}

	v.cachedVersion = repoVer
	return nil
}
