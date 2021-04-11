package kube

import (
	"os"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestGetConfigWithCustomCreds(t *testing.T) {

	_ = os.Setenv("IN_CLUSTER", "false")
	_ = os.Setenv("MASTERURL", "test")
	_ = os.Setenv("BEARERTOKEN", "token")

	conf, _ := BuildConfig()

	if conf.BearerToken != "token" {
		t.Errorf("'TestGetConfigWithCustomCreds' failed, expected '%v', got '%v'", "token", conf.BearerToken)
	}

	if conf.Host != "test" {
		t.Errorf("'TestGetConfigWithCustomCreds' failed, expected '%v', got '%v'", "test", conf.Host)
	}

}
