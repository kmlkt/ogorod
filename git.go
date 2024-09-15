package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func (c Config) Download() {
	for _, site := range c.Sites {
		site.Download()
	}
}

func (s Site) Download() {
	if PathExists(LocalPath(s)) {
		s.GitReset()
	} else {
		s.GitClone()
	}
}

func (s Site) GitReset() {
	repo, err := git.PlainOpen(LocalPath(s))
	StupidHandle(err)
	repo.Fetch(&git.FetchOptions{})
	worktree, err := repo.Worktree()
	StupidHandle(err)
	hash, err := repo.ResolveRevision(plumbing.Revision("origin/" + s.Branch))
	StupidHandle(err)
	worktree.Reset(&git.ResetOptions{
		Commit: *hash,
		Mode:   git.HardReset,
	})
}

func (s Site) GitClone() {
	_, err := git.PlainClone(LocalPath(s), false, &git.CloneOptions{
		URL:           s.Repository,
		ReferenceName: plumbing.ReferenceName(s.Branch),
	})
	StupidHandle(err)
}
