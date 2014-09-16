package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RepoConfig struct {
	RealPath string `json:"real_path"`

	SshCloneUrl  string `json:"ssh_clone_url"`
	HttpCloneUrl string `json:"http_clone_url"`
	GitCloneUrl  string `json:"git_clone_url"`

	CustomPreReceivePath  string `json:"custom_pre_receive_path"`
	CustomPostReceivePath string `json:"custom_post_receive_path"`
	CustomUpdatePath      string `json:"custom_update_path"`
}

type InternalApi interface {
	GetRepoConfig(string, string) (*RepoConfig, error)
}

type GitoriousInternalApi struct {
	ApiUrl string
}

func (a *GitoriousInternalApi) GetRepoConfig(repoPath, username string) (*RepoConfig, error) {
	path := fmt.Sprintf("/repo-config?repo_path=%v&username=%v", repoPath, username)
	var repoConfig RepoConfig

	if err := a.getJson(path, &repoConfig); err != nil {
		return nil, err
	}

	return &repoConfig, nil
}

func (a *GitoriousInternalApi) getJson(path string, target interface{}) error {
	url := a.ApiUrl + path

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	request.Header.Add("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("got status %v from %v", response.StatusCode, url))
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(target)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", target)

	return nil
}