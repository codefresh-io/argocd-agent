package git

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/google/go-github/github"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestExtractRepoAndOwnerFromUrl(t *testing.T) {
	urls := []string{
		"ssh://user@host.xz/path/owner/repo.git/",
		"ssh://user@host.xz:8080/path/owner/repo.git/",
		"git://host.xz/path/owner/repo.git/",
		"git://host.xz:8080/path/owner/repo.git/",
		"http://host.xz/path/owner/repo.git/",
		"http://host.xz:8080/path/owner/repo.git/",
		"ftp://host.xz/path/owner/repo.git/",
		"ftp://host.xz:8080/path/owner/repo.git/",
	}

	for _, url := range urls {
		err, owner, repo := extractRepoAndOwnerFromUrl(url)
		if err != nil {
			t.Errorf("'ExtractRepoAndOwnerFromUrl' check error failed, error: %v", err.Error())
		}
		if owner != "owner" {
			t.Errorf("'ExtractRepoAndOwnerFromUrl check owner' failed, expected '%v', got '%v'", "owner", owner)
		}
		if repo != "repo" {
			t.Errorf("'ExtractRepoAndOwnerFromUrl check repo' failed, expected '%v', got '%v'", "repo", repo)
		}
	}
}

func TestGetInstance(t *testing.T) {
	auth := struct {
		Type           string `json:"type"`
		Username       string `json:"username"`
		Password       string `json:"password"`
		ApiHost        string `json:"apiHost"`
		ApiPathPrefix  string `json:"apiPathPrefix"`
		SshPrivateKey  string `json:"sshPrivateKey"`
		AppId          string `json:"appId"`
		InstallationId string `json:"installationId"`
		PrivateKey     string `json:"privateKey"`
	}{Type: "git.github", Username: "u", Password: "p", ApiHost: "h", ApiPathPrefix: "pr", SshPrivateKey: "pk", AppId: "123", InstallationId: "123", PrivateKey: "test"}
	data := struct {
		Auth struct {
			Type           string `json:"type"`
			Username       string `json:"username"`
			Password       string `json:"password"`
			ApiHost        string `json:"apiHost"`
			ApiPathPrefix  string `json:"apiPathPrefix"`
			SshPrivateKey  string `json:"sshPrivateKey"`
			AppId          string `json:"appId"`
			InstallationId string `json:"installationId"`
			PrivateKey     string `json:"privateKey"`
		} `json:"auth"`
	}{
		Auth: auth,
	}
	context := codefresh.ContextPayload{
		Metadata: struct {
			Name string `json:"name"`
		}{},
		Spec: struct {
			Type string `json:"type"`
			Data struct {
				Auth struct {
					Type           string `json:"type"`
					Username       string `json:"username"`
					Password       string `json:"password"`
					ApiHost        string `json:"apiHost"`
					ApiPathPrefix  string `json:"apiPathPrefix"`
					SshPrivateKey  string `json:"sshPrivateKey"`
					AppId          string `json:"appId"`
					InstallationId string `json:"installationId"`
					PrivateKey     string `json:"privateKey"`
				} `json:"auth"`
			} `json:"data"`
		}{
			Data: data,
			Type: "git.github",
		},
	}
	store.SetGitContext(context)
	err, api := GetInstance("https://github.com/owner/repo")
	_ = api
	if err != nil {
		t.Errorf("'GetInstance' check error failed, error: %v", err.Error())
	}
}

func TestGetGitHttpClientBaseOnGithubApp(t *testing.T) {
	key := "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcGdJQkFBS0NBUUVBeEN4NzE5VDZBM3dIMXUxWS9WVlN1WmRJNGJrVDdGaWFUU2dHcmxXdG84SDhTL2lZCmFCM2dRdTRKVk4wQlJYei8reHVSMGw3ejVTeVp1cGd2WkY4RkhHYjAva29TQ3JHSytnck1oTUV4NU82NE0yd1oKWkN6QnJucDNYdUR6STI1RGJzZDhXYllTSm43d0N6cmlFMEFWMFFTcjV5cEFScHBpTDZ4QjQ2d1FBSkJzZXl4Ngp0cG12UjM3V3NUWVdpOFFVZ3JieFFYZ1o1VVBpVGdQNnlsMExEbXJjbTRDZ1kwNlpXRFhnaE1wKzc0V2pYRWt4CmdwQlJFL2E2UHZlTHZtZnUrYzdMRHY3VWNsSHRTcVV1RjVITnE5OFRISlFhd2dRSGJYMnN6VWJNYjhaRjcxRzYKK05Lb05VcFpnQUJMeHJhcGsrR28rakw4NW03T3JmWGsxdEVyQ1FJREFRQUJBb0lCQVFDbC9HMmRKWnVWcnpDQwo3cmpKUVpTSmJEUkNxWEx1RzlvVFJyYkFjOFpFTlRMZ3BTdHZqVGZmNmNFRGlTdzJPNW5zUWx1VUFMdWxRYU9oCmVucy9GaGNnL1F4MnpQMlBCc0pzNXc0OWxhbzk1cTc4ODQ1WWNIWkF4MmFSWlF6VkFjc1V4TDIydXBPSTl3YnMKdVpub1orVU53a0loaW1Kd1d0aVJOZE5hYkkvdHFNZGVENWtiS3NxcWZmVmdudGR3RTl3Q2N6RysrTE1udzFvUAprZWg0MVRiUnRoeENXdldGZm4zWjNaRkNxRjJRRHcvd0NXUzdaSkt0M2pKUTZiQ3JyNXBkaU4yUWNIbnFlVjgxCjZubFZXeEF6SDUwNEdEdExpbHpBekRva1krazZJRzhXcmNRNHdtdWhGZEtzK3JZbCtMSXR2WEJQWGlVRWNJRk8KUVkzZS9ibkJBb0dCQVArUXRyajh2OEd2UkZtM0JncE9hcGtFWVNlcUNiT2R3cm9jazdNUml3NWZLMWt1bVFScwp6UEVqNURyMHJBdTNIbUJ6eVB4cHM2amw2TVFnQ2VxNW03ajM2bExYUG41MTFUUDQ0NFFDUWZuZERUVGdnQ0gvClRadVEreVNjSEtITDVudGhCOStaNTY2OVRSSHFLNEMwVm5JQ1U4aFlVU2lGMldTQmhwdklscGU5QW9HQkFNU0IKNkdzcFBiL0dCYlkyRjN5NVBuY3kyUm50emszbXlNZWpIamNXNnRndWNNWFBYd3owOUlvb2FTaVV5SS9ON0UvYgplZENCdzM4QnlWYmUyWWU3V3dDWTVWd3l2TXhFditXZm85end5STI4ZHd4MUhUVXZEQ1FmaWRsL01iK3hFZXdGClpKbnpkOURyTWRZUEt1eTI1Z1lvM3ExVU0zZkdVNXhlKzRtSnJiODlBb0dCQU1WSlVjVThXRXVNb1pjZ1V1bGgKMzZpQVdQL2xvOWVrMGM0YWdXcWJBRjMzMmQ3ZXVnRlFmR1VxNytVVFBEMU8vNFExM2RIOVIxUDdKOVUvWm1odApJR21KK0xvNnIyT3dVd1hyL0xiTGgyTDc0bFlQZU5yRjI0TmNTSVBhZjcvblIrVzI0ZjBiTWw4U2c0eHcyV1JoCjB0bndNZjFYTUUrNEJEb3lRMWUvWVlHQkFvR0JBTGFSV3pYMFl3SkJJQjFodEFDVXVveFVHWkFWZUk2MzArSm0Ka2pQc2Z0UEtrY3UyRmtFYmMvYklCS3RIVCs4TENucEhGcTI1WWNBbUVNRTgyaTFZeS91S0VjM085Y2x5TmpkSQpVaDE3TjFrM3VBTkM2NWYxMWZuWnMyRDI0Mm1OUVhGZXNWQzIrcUtIWVEzWG1iSERXNEp0aGpUUy9kNVJ6R3lECmNuOGVBdWFoQW9HQkFQM2tobE45VEZnUXZEOHZzVDRSVlhBbzNRa0pFa2t4MDhrcEVHbkF2N0UyWWdBZ0hxT0gKWjlWTjhnMlhhUjVpZklNd1Rwa1RJVlVKaEYycDlLcFVsNm1SekFndjNqcStBMm9qdDRESnJ1eDU0M2pyamR5aAo0QUJkSXp6QW9aVnM2dlFJSXF2cGI0NWJwOXpVOEs3b3BoZ2p5SGVvQWcvaTFXUTNoVUtuQko4eQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="
	auth := struct {
		Type           string `json:"type"`
		Username       string `json:"username"`
		Password       string `json:"password"`
		ApiHost        string `json:"apiHost"`
		ApiPathPrefix  string `json:"apiPathPrefix"`
		SshPrivateKey  string `json:"sshPrivateKey"`
		AppId          string `json:"appId"`
		InstallationId string `json:"installationId"`
		PrivateKey     string `json:"privateKey"`
	}{Type: "git.github", Username: "u", Password: "p", ApiHost: "h", ApiPathPrefix: "pr", SshPrivateKey: "pk", AppId: "64781", InstallationId: "9646502", PrivateKey: key}
	data := struct {
		Auth struct {
			Type           string `json:"type"`
			Username       string `json:"username"`
			Password       string `json:"password"`
			ApiHost        string `json:"apiHost"`
			ApiPathPrefix  string `json:"apiPathPrefix"`
			SshPrivateKey  string `json:"sshPrivateKey"`
			AppId          string `json:"appId"`
			InstallationId string `json:"installationId"`
			PrivateKey     string `json:"privateKey"`
		} `json:"auth"`
	}{
		Auth: auth,
	}
	context := codefresh.ContextPayload{
		Metadata: struct {
			Name string `json:"name"`
		}{},
		Spec: struct {
			Type string `json:"type"`
			Data struct {
				Auth struct {
					Type           string `json:"type"`
					Username       string `json:"username"`
					Password       string `json:"password"`
					ApiHost        string `json:"apiHost"`
					ApiPathPrefix  string `json:"apiPathPrefix"`
					SshPrivateKey  string `json:"sshPrivateKey"`
					AppId          string `json:"appId"`
					InstallationId string `json:"installationId"`
					PrivateKey     string `json:"privateKey"`
				} `json:"auth"`
			} `json:"data"`
		}{
			Data: data,
			Type: "git.github-app",
		},
	}

	err, client := getGitHttpClient(&context)
	if err != nil {
		t.Error("Get github client base on github app should fail without error")
	}

	if client == nil {
		t.Error("Client transport should be created")
	}

}

func TestGetGitHttpClientBaseOnUnknownType(t *testing.T) {
	key := "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcGdJQkFBS0NBUUVBeEN4NzE5VDZBM3dIMXUxWS9WVlN1WmRJNGJrVDdGaWFUU2dHcmxXdG84SDhTL2lZCmFCM2dRdTRKVk4wQlJYei8reHVSMGw3ejVTeVp1cGd2WkY4RkhHYjAva29TQ3JHSytnck1oTUV4NU82NE0yd1oKWkN6QnJucDNYdUR6STI1RGJzZDhXYllTSm43d0N6cmlFMEFWMFFTcjV5cEFScHBpTDZ4QjQ2d1FBSkJzZXl4Ngp0cG12UjM3V3NUWVdpOFFVZ3JieFFYZ1o1VVBpVGdQNnlsMExEbXJjbTRDZ1kwNlpXRFhnaE1wKzc0V2pYRWt4CmdwQlJFL2E2UHZlTHZtZnUrYzdMRHY3VWNsSHRTcVV1RjVITnE5OFRISlFhd2dRSGJYMnN6VWJNYjhaRjcxRzYKK05Lb05VcFpnQUJMeHJhcGsrR28rakw4NW03T3JmWGsxdEVyQ1FJREFRQUJBb0lCQVFDbC9HMmRKWnVWcnpDQwo3cmpKUVpTSmJEUkNxWEx1RzlvVFJyYkFjOFpFTlRMZ3BTdHZqVGZmNmNFRGlTdzJPNW5zUWx1VUFMdWxRYU9oCmVucy9GaGNnL1F4MnpQMlBCc0pzNXc0OWxhbzk1cTc4ODQ1WWNIWkF4MmFSWlF6VkFjc1V4TDIydXBPSTl3YnMKdVpub1orVU53a0loaW1Kd1d0aVJOZE5hYkkvdHFNZGVENWtiS3NxcWZmVmdudGR3RTl3Q2N6RysrTE1udzFvUAprZWg0MVRiUnRoeENXdldGZm4zWjNaRkNxRjJRRHcvd0NXUzdaSkt0M2pKUTZiQ3JyNXBkaU4yUWNIbnFlVjgxCjZubFZXeEF6SDUwNEdEdExpbHpBekRva1krazZJRzhXcmNRNHdtdWhGZEtzK3JZbCtMSXR2WEJQWGlVRWNJRk8KUVkzZS9ibkJBb0dCQVArUXRyajh2OEd2UkZtM0JncE9hcGtFWVNlcUNiT2R3cm9jazdNUml3NWZLMWt1bVFScwp6UEVqNURyMHJBdTNIbUJ6eVB4cHM2amw2TVFnQ2VxNW03ajM2bExYUG41MTFUUDQ0NFFDUWZuZERUVGdnQ0gvClRadVEreVNjSEtITDVudGhCOStaNTY2OVRSSHFLNEMwVm5JQ1U4aFlVU2lGMldTQmhwdklscGU5QW9HQkFNU0IKNkdzcFBiL0dCYlkyRjN5NVBuY3kyUm50emszbXlNZWpIamNXNnRndWNNWFBYd3owOUlvb2FTaVV5SS9ON0UvYgplZENCdzM4QnlWYmUyWWU3V3dDWTVWd3l2TXhFditXZm85end5STI4ZHd4MUhUVXZEQ1FmaWRsL01iK3hFZXdGClpKbnpkOURyTWRZUEt1eTI1Z1lvM3ExVU0zZkdVNXhlKzRtSnJiODlBb0dCQU1WSlVjVThXRXVNb1pjZ1V1bGgKMzZpQVdQL2xvOWVrMGM0YWdXcWJBRjMzMmQ3ZXVnRlFmR1VxNytVVFBEMU8vNFExM2RIOVIxUDdKOVUvWm1odApJR21KK0xvNnIyT3dVd1hyL0xiTGgyTDc0bFlQZU5yRjI0TmNTSVBhZjcvblIrVzI0ZjBiTWw4U2c0eHcyV1JoCjB0bndNZjFYTUUrNEJEb3lRMWUvWVlHQkFvR0JBTGFSV3pYMFl3SkJJQjFodEFDVXVveFVHWkFWZUk2MzArSm0Ka2pQc2Z0UEtrY3UyRmtFYmMvYklCS3RIVCs4TENucEhGcTI1WWNBbUVNRTgyaTFZeS91S0VjM085Y2x5TmpkSQpVaDE3TjFrM3VBTkM2NWYxMWZuWnMyRDI0Mm1OUVhGZXNWQzIrcUtIWVEzWG1iSERXNEp0aGpUUy9kNVJ6R3lECmNuOGVBdWFoQW9HQkFQM2tobE45VEZnUXZEOHZzVDRSVlhBbzNRa0pFa2t4MDhrcEVHbkF2N0UyWWdBZ0hxT0gKWjlWTjhnMlhhUjVpZklNd1Rwa1RJVlVKaEYycDlLcFVsNm1SekFndjNqcStBMm9qdDRESnJ1eDU0M2pyamR5aAo0QUJkSXp6QW9aVnM2dlFJSXF2cGI0NWJwOXpVOEs3b3BoZ2p5SGVvQWcvaTFXUTNoVUtuQko4eQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="
	auth := struct {
		Type           string `json:"type"`
		Username       string `json:"username"`
		Password       string `json:"password"`
		ApiHost        string `json:"apiHost"`
		ApiPathPrefix  string `json:"apiPathPrefix"`
		SshPrivateKey  string `json:"sshPrivateKey"`
		AppId          string `json:"appId"`
		InstallationId string `json:"installationId"`
		PrivateKey     string `json:"privateKey"`
	}{Type: "git.github", Username: "u", Password: "p", ApiHost: "h", ApiPathPrefix: "pr", SshPrivateKey: "pk", AppId: "64781", InstallationId: "9646502", PrivateKey: key}
	data := struct {
		Auth struct {
			Type           string `json:"type"`
			Username       string `json:"username"`
			Password       string `json:"password"`
			ApiHost        string `json:"apiHost"`
			ApiPathPrefix  string `json:"apiPathPrefix"`
			SshPrivateKey  string `json:"sshPrivateKey"`
			AppId          string `json:"appId"`
			InstallationId string `json:"installationId"`
			PrivateKey     string `json:"privateKey"`
		} `json:"auth"`
	}{
		Auth: auth,
	}
	context := codefresh.ContextPayload{
		Metadata: struct {
			Name string `json:"name"`
		}{},
		Spec: struct {
			Type string `json:"type"`
			Data struct {
				Auth struct {
					Type           string `json:"type"`
					Username       string `json:"username"`
					Password       string `json:"password"`
					ApiHost        string `json:"apiHost"`
					ApiPathPrefix  string `json:"apiPathPrefix"`
					SshPrivateKey  string `json:"sshPrivateKey"`
					AppId          string `json:"appId"`
					InstallationId string `json:"installationId"`
					PrivateKey     string `json:"privateKey"`
				} `json:"auth"`
			} `json:"data"`
		}{
			Data: data,
			Type: "git.github-app-unknown",
		},
	}

	err, _ := getGitHttpClient(&context)
	if err == nil || err.Error() != "Cant handle unknown git type" {
		t.Error("Get github client base on github app should fail without error")
	}
}

func TestGetComittersByCommits(t *testing.T) {
	api := api{
		Client: nil,
		Owner:  "",
		Repo:   "",
		Ctx:    nil,
	}

	login := "Ivan"
	avatar := "Link"

	commits := []*github.RepositoryCommit{
		{
			SHA:    nil,
			Commit: nil,
			Author: &github.User{
				Login:             &login,
				ID:                nil,
				NodeID:            nil,
				AvatarURL:         &avatar,
				HTMLURL:           nil,
				GravatarID:        nil,
				Name:              nil,
				Company:           nil,
				Blog:              nil,
				Location:          nil,
				Email:             nil,
				Hireable:          nil,
				Bio:               nil,
				PublicRepos:       nil,
				PublicGists:       nil,
				Followers:         nil,
				Following:         nil,
				CreatedAt:         nil,
				UpdatedAt:         nil,
				SuspendedAt:       nil,
				Type:              nil,
				SiteAdmin:         nil,
				TotalPrivateRepos: nil,
				OwnedPrivateRepos: nil,
				PrivateGists:      nil,
				DiskUsage:         nil,
				Collaborators:     nil,
				Plan:              nil,
				URL:               nil,
				EventsURL:         nil,
				FollowingURL:      nil,
				FollowersURL:      nil,
				GistsURL:          nil,
				OrganizationsURL:  nil,
				ReceivedEventsURL: nil,
				ReposURL:          nil,
				StarredURL:        nil,
				SubscriptionsURL:  nil,
				TextMatches:       nil,
				Permissions:       nil,
			},
			Committer:   nil,
			Parents:     nil,
			HTMLURL:     nil,
			URL:         nil,
			CommentsURL: nil,
			Stats:       nil,
			Files:       nil,
		},
	}

	err, result := api.GetComittersByCommits(commits)
	if err != nil {
		t.Error("Shouldnt fail during retrieve committers from commit")
	}

	if len(result) != 1 {
		t.Error("Should retrieve only one comitter from commits")
	}

	if result[0].Name != login {
		t.Error("Wrong login name")
	}

	if result[0].Avatar != avatar {
		t.Error("Wrong avatar")
	}

}
