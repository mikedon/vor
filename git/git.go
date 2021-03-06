package git

import (
	"sync"

	"github.com/spf13/viper"
	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/system"
)

type NativeGit interface {
	Call(string) (string, error)
	Stash() bool
	UnStash(string)
}

type GitClient struct {
	Path string
}

var once sync.Once
var client GitClient

func New() GitClient {
	once.Do(func() {
		localGit := viper.GetString("git.path")
		exists, fsErr := system.Exists(localGit)
		if fsErr != nil || !exists {
			system.Exit("Could not find local git client at " + "\"" + localGit + "\"")
		}
		client.Path = localGit

		_, gitErr := system.Exec(client.Path + " status")
		if gitErr != nil {
			system.Exit("Invalid git repository")
		}
	})
	return client
}

// Call – call a Client command by name
// you can pass arguments as well E.G:
// Client.Call("checkout -b my-branch-name")
// returns the text output of the command and a standard error (if any)
func (git GitClient) Call(command string) (string, error) {
	logger.Debug("calling 'git " + command + "'")
	return system.Exec(git.Path + " " + command)
}

// Stash – stash changes if the working directory is unclean
func (git GitClient) Stash() {
	_, err := git.Call("stash")
	if err != nil {
		system.Exit("error stashing changes")
	}
}

// UnStash – unstash the top most stash (called after a Stash())
func (git GitClient) UnStash() {
	_, err := git.Call("stash apply")
	if err != nil {
		system.Exit("error unstashing changes")
	}
}
