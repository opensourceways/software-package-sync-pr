package syncprimpl

import (
	"fmt"
	"strings"
)

const repoEndpointPrefix = "https://"

type Config struct {
	WorkDir       string          `json:"work_dir"         required:"true"`
	RobotRepo     robotRepoConfig `json:"robot_repo"       required:"true"`
	SyncRepoShell string          `json:"sync_repo_shell"  required:"true"`
}

func (cfg *Config) Validate() error {
	return cfg.RobotRepo.validate()
}

// robotRepoConfig
type robotRepoConfig struct {
	HomeLink   string     `json:"home_link"   required:"true"`
	Credential credential `json:"credential"  required:"true"`
}

func (t *robotRepoConfig) validate() error {
	if !strings.HasPrefix(t.HomeLink, repoEndpointPrefix) {
		return fmt.Errorf("unsupported protocol")
	}

	return nil
}

func (t *robotRepoConfig) remoteURL() string {
	e := strings.TrimSuffix(t.HomeLink, "/")

	return fmt.Sprintf(
		"%s%s:%sxi@%s/",
		repoEndpointPrefix,
		t.Credential.UserName,
		t.Credential.Token,
		strings.TrimPrefix(e, repoEndpointPrefix),
	)
}

// credential
type credential struct {
	UserName string `json:"user_name"  required:"true"`
	Token    string `json:"token"      required:"true"`
}
