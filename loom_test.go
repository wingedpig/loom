package loom

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_Run(t *testing.T) {

	// These tests assume you can ssh log into localhost. Adjust config.User and such accordingly.
	var config Config
	config.Host = "127.0.0.1:22"
	localout, err := config.Local("/bin/ls -l /var/tmp")
	if err != nil {
		t.Errorf("Local error, %s", err)
	}
	remoteout, err := config.Run("/bin/ls -l /var/tmp")
	if err != nil {
		t.Errorf("Run error, %s", err)
	}
	remoteout = strings.Replace(remoteout, "\r", "", -1)
	if localout != remoteout {
		t.Errorf("Local and remote differ")
	}
	dir, err := ioutil.TempDir("/var/tmp", "loom")
	if err != nil {
		t.Errorf("Cannot create temp directory, %s", err)
	}
	defer os.RemoveAll(dir)

	file1 := dir + "/a"
	file2 := dir + "/b"
	testText := `Lorem ipsum Minim voluptate aliquip commodo.`

	err = config.PutString(testText, file1)
	if err != nil {
		t.Errorf("PutString error, %s", err)
	}
	err = config.Get(file1, file2)
	if err != nil {
		t.Errorf("Get error, %s", err)
	}
	// read file2, verify that it matches testText
	verifyText, err := ioutil.ReadFile(file2)
	if err != nil {
		t.Errorf("Get failed, can't read %s, %s", file2, err)
	}
	if string(verifyText) != testText {
		t.Errorf("Text mismatch, expected '%s', got '%s'", testText, verifyText)
	}
}

func TestInjectSudoPasswordIfNecessary(t *testing.T) {
	config := Config{
		User:          "user",
		Password:      "user",
		Host:          "127.0.0.1:22",
		DisplayOutput: true,
		AbortOnError:  false,
	}

	val, err := config.Sudo("ls -l")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(val)
}

func TestEnsureSudoMatcherShouldMatch(t *testing.T) {
	totalPayload := `Command 1
	[sudo] password for user:
	Command2`

	sm := newSudoMatcher("user")

	match := sm.Match([]byte(totalPayload))
	if !match {
		t.Errorf("Expected a match")
	}
}

func TestEnsureSudoMatcherShouldNotMatch(t *testing.T) {
	totalPayload := `some othercommand
	foo
	bar`

	sm := newSudoMatcher("user")

	match := sm.Match([]byte(totalPayload))
	if match {
		t.Errorf("No match was expected")
	}
}

func TestEnsureSudoMatcherShouldMatchAfterMultipleMatchTries(t *testing.T) {
	sm := newSudoMatcher("user")

	match := sm.Match([]byte(`Command 1
	[sudo] password`))
	if match {
		t.Errorf("Expected no match")
	}

	match = sm.Match([]byte(` for user:`))
	if !match {
		t.Errorf("Expected a match")
	}
}
