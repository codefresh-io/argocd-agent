package provider

import (
	"github.com/xanzy/go-gitlab"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

type MockGitlabApi struct {
}

func (gl *MockGitlabApi) ListProjects(search string) (error, []*gitlab.Project) {
	return nil, []*gitlab.Project{
		&gitlab.Project{
			ID:                               0,
			Description:                      "",
			DefaultBranch:                    "",
			Public:                           false,
			Visibility:                       "",
			SSHURLToRepo:                     "",
			HTTPURLToRepo:                    "",
			WebURL:                           "",
			ReadmeURL:                        "",
			TagList:                          nil,
			Owner:                            nil,
			Name:                             "",
			NameWithNamespace:                "",
			Path:                             "",
			PathWithNamespace:                "p.kostohrys/test",
			IssuesEnabled:                    false,
			OpenIssuesCount:                  0,
			MergeRequestsEnabled:             false,
			ApprovalsBeforeMerge:             0,
			JobsEnabled:                      false,
			WikiEnabled:                      false,
			SnippetsEnabled:                  false,
			ResolveOutdatedDiffDiscussions:   false,
			ContainerExpirationPolicy:        nil,
			ContainerRegistryEnabled:         false,
			CreatedAt:                        nil,
			LastActivityAt:                   nil,
			CreatorID:                        0,
			Namespace:                        nil,
			ImportStatus:                     "",
			ImportError:                      "",
			Permissions:                      nil,
			MarkedForDeletionAt:              nil,
			EmptyRepo:                        false,
			Archived:                         false,
			AvatarURL:                        "",
			LicenseURL:                       "",
			License:                          nil,
			SharedRunnersEnabled:             false,
			ForksCount:                       0,
			StarCount:                        0,
			RunnersToken:                     "",
			PublicBuilds:                     false,
			AllowMergeOnSkippedPipeline:      false,
			OnlyAllowMergeIfPipelineSucceeds: false,
			OnlyAllowMergeIfAllDiscussionsAreResolved: false,
			RemoveSourceBranchAfterMerge:              false,
			LFSEnabled:                                false,
			RequestAccessEnabled:                      false,
			MergeMethod:                               "",
			ForkedFromProject:                         nil,
			Mirror:                                    false,
			MirrorUserID:                              0,
			MirrorTriggerBuilds:                       false,
			OnlyMirrorProtectedBranches:               false,
			MirrorOverwritesDivergedBranches:          false,
			PackagesEnabled:                           false,
			ServiceDeskEnabled:                        false,
			ServiceDeskAddress:                        "",
			IssuesAccessLevel:                         "",
			RepositoryAccessLevel:                     "",
			MergeRequestsAccessLevel:                  "",
			ForkingAccessLevel:                        "",
			WikiAccessLevel:                           "",
			BuildsAccessLevel:                         "",
			SnippetsAccessLevel:                       "",
			PagesAccessLevel:                          "",
			OperationsAccessLevel:                     "",
			AutocloseReferencedIssues:                 false,
			SuggestionCommitMessage:                   "",
			CIForwardDeploymentEnabled:                false,
			SharedWithGroups:                          nil,
			Statistics:                                nil,
			Links:                                     nil,
			CIConfigPath:                              "",
			CIDefaultGitDepth:                         0,
			CustomAttributes:                          nil,
			ComplianceFrameworks:                      nil,
			BuildCoverageRegex:                        "",
			IssuesTemplate:                            "",
			MergeRequestsTemplate:                     "",
		},
	}
}

func (gl *MockGitlabApi) RetrieveAvatar(email string) (error, string) {
	return nil, ""
}

func (gl *MockGitlabApi) GetCommit(projectId int, revision string) (error, *gitlab.Commit) {
	return nil, &gitlab.Commit{
		ID:             "",
		ShortID:        "",
		Title:          "",
		AuthorName:     "",
		AuthorEmail:    "",
		AuthoredDate:   nil,
		CommitterName:  "",
		CommitterEmail: "",
		CommittedDate:  nil,
		CreatedAt:      nil,
		Message:        "Test",
		ParentIDs:      nil,
		Stats:          nil,
		Status:         nil,
		LastPipeline:   nil,
		ProjectID:      0,
		WebURL:         "",
	}
}

func TestGetCommitByRevisionWithWrongAmountOfArguments(t *testing.T) {
	gl := &Gitlab{api: nil}
	err, _ := gl.GetCommitByRevision("p.kostohrys/test/test", "revision")
	if err == nil || err.Error() != "wrong amount of arguments" {
		t.Error("Should be failed with wrong amount of arguments")
	}
}

func TestGetCommitByRevision(t *testing.T) {
	gl := &Gitlab{api: &MockGitlabApi{}}
	err, commit := gl.GetCommitByRevision("p.kostohrys/test", "revision")
	if err != nil {
		t.Error("SHould be executed without error")
	}
	if *commit.Message != "Test" {
		t.Error("Wrong commit message")
	}
}

func TestInitGitlab(t *testing.T) {
	gl := NewGitlabProvider()
	if gl == nil {
		t.Error("Should be inited without error")
	}
}

func TestGetManifestDetails(t *testing.T) {
	gl := &Gitlab{api: &MockGitlabApi{}}
	_, gitops := gl.GetManifestRepoInfo("test", "123")
	if gitops == nil {
		t.Error("Should be able retrieve manifest details")
	}
}
