package provider

import (
	"fmt"
	"testing"
)

func TestGithub_GetCommitByRevision(t *testing.T) {

	gl := NewGitlabProvider()
	_, commit := gl.GetCommitByRevision("p.kostohrys/test", "c50c435e5ef05050a48e465a51df321202c6a6cf")
	fmt.Println(commit)
}
