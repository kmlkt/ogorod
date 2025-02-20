package sys

import (
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"strings"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

const AppTag = "ogorod"

var SiteDir string

func init() {
	userDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	sudoUID, sudo := os.LookupEnv("SUDO_UID")
	if sudo {
		sudoUser, err := user.LookupId(sudoUID)
		if err != nil {
			panic(err)
		}
		userDir = sudoUser.HomeDir
	}
	SiteDir = filepath.Join(userDir, AppTag)
	os.MkdirAll(SiteDir, 0750)
}

func (s Service) download() (bool, string, error) {
	lastLocal, currentPath, err := s.lastLocalCommit()
	slog.Debug("Last local commit", "id", s.ID, "repo", s.Repository, "commit", lastLocal)
	if err != nil {
		return false, "", err
	}
	lastRemote, err := s.lastRemoteCommit()
	slog.Debug("Last remote commit", "err", err, "repo", s.Repository, "commit", lastRemote)
	if err != nil {
		return false, "", err
	}

	if lastLocal == lastRemote {
		return false, currentPath, nil
	}

	nextPath := filepath.Join(SiteDir, s.ID.String()+"_"+s.ID.Next().String())
	_, err = git.PlainClone(nextPath, false, &git.CloneOptions{
		URL:           s.repositoryURL(),
		ReferenceName: plumbing.ReferenceName(s.Branch),
	})
	if err != nil {
		return false, "", err
	}
	slog.Debug("Downloaded repository", "id", s.ID, "repo", s.Repository, "path", nextPath)
	return true, nextPath, nil
}

func (s Service) lastLocalCommit() (string, string, error) {
	currentPath, err := s.currentPath()
	if currentPath == "" {
		return "", "", nil
	}
	if err != nil {
		return "", "", err
	}
	currentRepo, err := git.PlainOpen(currentPath + "/.git")
	if err != nil {
		return "", "", err
	}
	head, err := currentRepo.Head()
	if err != nil {
		return "", "", err
	}
	return head.Hash().String(), currentPath, nil
}

func (s Service) currentPath() (string, error) {
	dirs, err := os.ReadDir(SiteDir)
	if err != nil {
		return "", err
	}
	prefix := s.ID.String()
	matching := make([]string, 0)
	for _, dir := range dirs {
		if dir.IsDir() && strings.HasPrefix(dir.Name(), prefix) {
			matching = append(matching, dir.Name())
		}
	}
	if len(matching) == 0 {
		return "", nil
	}
	slices.Sort(matching)
	return filepath.Join(SiteDir, matching[len(matching)-1]), nil
}

func (s Service) lastRemoteCommit() (string, error) {
	storer := memory.NewStorage()
	fs := memfs.New()

	repo, err := git.Init(storer, fs)
	if err != nil {
		return "", err
	}
	remoteName := "origin"
	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{s.repositoryURL()},
	})
	if err != nil {
		return "", err
	}
	err = remote.Fetch(&git.FetchOptions{})
	if err != nil {
		return "", err
	}
	slog.Debug("Remote", "name", remote.String())
	branches, err := remote.List(&git.ListOptions{})
	if err != nil {
		return "", err
	}
	slog.Debug("Remote", "list", branches)
	branchName := "refs/heads/" + s.Branch
	for _, branch := range branches {
		slog.Debug("Branch", "name", branch.Name().String())
		if branch.Name().String() == branchName {
			return branch.Hash().String(), nil
		}
	}
	return "", nil
}

func (s Service) repositoryURL() string {
	return fmt.Sprintf("https://%s.git", s.Repository)
}

func ClearStorage(pm ProcessMap) error {
	dirs, err := os.ReadDir(SiteDir)
	if err != nil {
		return err
	}
	toKeep := make(map[string]struct{})
	for _, process := range pm {
		toKeep[process.Path] = struct{}{}
	}
	for _, dir := range dirs {
		dirPath := filepath.Join(SiteDir, dir.Name())
		_, keep := toKeep[dirPath]
		if !keep {
			err := os.RemoveAll(dirPath)
			if err != nil {
				return err
			}
			slog.Debug("Removed folder", "path", dirPath)
		}
	}
	return nil
}
