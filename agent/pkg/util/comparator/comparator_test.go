package comparator

import (
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestEnvironmentComparatorWithSameEnv(t *testing.T) {

	envComparator := EnvComparator{}

	env1 := codefreshSdk.Environment{
		Gitops:       codefreshSdk.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   nil,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	env2 := codefreshSdk.Environment{
		Gitops:       codefreshSdk.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   nil,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	if !envComparator.Compare(&env1, &env2) {
		t.Errorf("'EnvComparator' comparation failed")
	}

}

func TestEnvironmentComparatorWithSameEnvAndActivities(t *testing.T) {

	envComparator := EnvComparator{}

	act1 := codefreshSdk.EnvironmentActivity{
		Name:         "test",
		TargetImages: nil,
		Status:       "test",
		LiveImages:   nil,
	}

	act2 := codefreshSdk.EnvironmentActivity{
		Name:         "test2",
		TargetImages: nil,
		Status:       "test2",
		LiveImages:   nil,
	}

	activities1 := make([]codefreshSdk.EnvironmentActivity, 0)

	activities2 := make([]codefreshSdk.EnvironmentActivity, 0)

	activities1 = append(activities1, act1)
	activities1 = append(activities1, act2)

	activities2 = append(activities2, act2)
	activities2 = append(activities2, act1)

	env1 := codefreshSdk.Environment{
		Gitops:       codefreshSdk.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   activities1,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	env2 := codefreshSdk.Environment{
		Gitops:       codefreshSdk.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   activities2,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	if !envComparator.Compare(&env1, &env2) {
		t.Errorf("'EnvComparator' comparation failed")
	}

}

func TestEnvironmentComparatorWithSameEnvAndDifferentActivities(t *testing.T) {

	envComparator := EnvComparator{}

	act1 := codefreshSdk.EnvironmentActivity{
		Name:         "test",
		TargetImages: nil,
		Status:       "test",
		LiveImages:   nil,
	}

	act2 := codefreshSdk.EnvironmentActivity{
		Name:         "test",
		TargetImages: nil,
		Status:       "test4",
		LiveImages:   nil,
	}

	activities1 := make([]codefreshSdk.EnvironmentActivity, 0)

	activities2 := make([]codefreshSdk.EnvironmentActivity, 0)

	activities1 = append(activities1, act1)
	activities1 = append(activities1, act2)

	activities2 = append(activities2, act1)

	env1 := codefreshSdk.Environment{
		Gitops:       codefreshSdk.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   activities1,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	env2 := codefreshSdk.Environment{
		Gitops:       codefreshSdk.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   activities2,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	if envComparator.Compare(&env1, &env2) {
		t.Errorf("'EnvComparator' comparation failed")
	}

}

func TestEnvironmentComparatorWithDiffEnv(t *testing.T) {

	envComparator := EnvComparator{}

	env1 := codefreshSdk.Environment{
		Gitops:       codefreshSdk.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    123,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   nil,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	env2 := codefreshSdk.Environment{
		Gitops:       codefreshSdk.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    12,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   nil,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	if envComparator.Compare(&env1, &env2) {
		t.Errorf("'EnvComparator' comparation failed")
	}

}

func TestEnvironmentComparatorWithDifferentGitops(t *testing.T) {

	envComparator := EnvComparator{}

	env1 := codefreshSdk.Environment{
		Gitops:       codefreshSdk.Gitops{},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    123,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   nil,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	env2 := codefreshSdk.Environment{
		Gitops: codefreshSdk.Gitops{
			Prs: []codefreshSdk.Annotation{
				{
					Key:   "test",
					Value: "test",
				},
			},
		},
		FinishedAt:   "",
		HealthStatus: "HEALTH",
		SyncStatus:   "OUT_OF_SYNC",
		HistoryId:    123,
		SyncRevision: "123",
		Name:         "Test",
		Activities:   nil,
		Resources:    nil,
		RepoUrl:      "https://google.com",
	}

	if !envComparator.Compare(&env1, &env2) {
		t.Errorf("'EnvComparator' comparation failed")
	}

}
