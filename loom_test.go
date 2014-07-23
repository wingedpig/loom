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
