package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/armon/go-socks5"
	"github.com/kost/vornex/v2/internal/testutils"
	"github.com/kost/vornex/v2/pkg/privileges"
	"github.com/kost/vornex/v2/pkg/result"
	"github.com/kost/vornex/v2/pkg/runner"
)

var libraryTestcases = map[string]testutils.TestCase{
	"sdk - one passive execution":         &vornexPassiveSingleLibrary{},
	"sdk - one execution - connect":       &vornexSingleLibrary{scanType: "c"},
	"sdk - multiple executions - connect": &vornexMultipleExecLibrary{scanType: "c"},
	"sdk - one execution - syn":           &vornexSingleLibrary{scanType: "s"},
	"sdk - multiple executions - syn":     &vornexMultipleExecLibrary{scanType: "s"},
	"sdk - connect with proxy":            &vornexWithSocks5{},
}

type vornexPassiveSingleLibrary struct {
}

func (h *vornexPassiveSingleLibrary) Execute() error {
	testFile := "test.txt"
	err := os.WriteFile(testFile, []byte("scanme.sh"), 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(testFile); err != nil {
			log.Printf("could not remove test file: %s\n", err)
		}
	}()

	options := runner.Options{
		HostsFile: testFile,
		Ports:     "80",
		Passive:   true,
		OnResult:  func(hr *result.HostResult) {},
	}

	vornexRunner, err := runner.NewRunner(&options)
	if err != nil {
		return err
	}
	defer func() {
		if err := vornexRunner.Close(); err != nil {
			log.Printf("could not close vornex runner: %s\n", err)
		}
	}()

	return vornexRunner.RunEnumeration(context.TODO())
}

type vornexSingleLibrary struct {
	scanType string
}

func (h *vornexSingleLibrary) Execute() error {
	if h.scanType == "s" && !privileges.IsPrivileged {
		usr, _ := user.Current()
		return errors.New("invalid user" + usr.Name)
	}

	testFile := "test.txt"
	err := os.WriteFile(testFile, []byte("scanme.sh"), 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(testFile); err != nil {
			log.Printf("could not remove test file: %s\n", err)
		}
	}()

	var got bool

	options := runner.Options{
		HostsFile: testFile,
		Ports:     "80",
		ScanType:  h.scanType,
		OnResult: func(hr *result.HostResult) {
			got = true
		},
		WarmUpTime: 2,
	}

	vornexRunner, err := runner.NewRunner(&options)
	if err != nil {
		return err
	}
	defer func() {
		if err := vornexRunner.Close(); err != nil {
			log.Printf("could not close vornex runner: %s\n", err)
		}
	}()

	if err = vornexRunner.RunEnumeration(context.TODO()); err != nil {
		return err
	}
	if !got {
		return errors.New("no results found")
	}

	return nil
}

type vornexMultipleExecLibrary struct {
	scanType string
}

func (h *vornexMultipleExecLibrary) Execute() error {
	if h.scanType == "s" && !privileges.IsPrivileged {
		usr, _ := user.Current()
		return errors.New("invalid user" + usr.Name)
	}

	testFile := "test.txt"
	err := os.WriteFile(testFile, []byte("scanme.sh"), 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(testFile); err != nil {
			log.Printf("could not remove test file: %s\n", err)
		}
	}()

	var got bool

	options := runner.Options{
		HostsFile: testFile,
		Ports:     "80",
		ScanType:  h.scanType,
		OnResult: func(hr *result.HostResult) {
			got = true
		},
		WarmUpTime: 2,
	}

	for i := 0; i < 3; i++ {
		vornexRunner, err := runner.NewRunner(&options)
		if err != nil {
			return err
		}

		if err = vornexRunner.RunEnumeration(context.TODO()); err != nil {
			return err
		}
		if !got {
			return errors.New("no results found")
		}
		if err := vornexRunner.Close(); err != nil {
			log.Printf("could not close vornex runner: %s\n", err)
		}
	}
	return nil
}

type vornexWithSocks5 struct{}

func (h *vornexWithSocks5) Execute() error {
	// Start local SOCKS5 proxy server with test:test credentials
	conf := &socks5.Config{
		Credentials: socks5.StaticCredentials{
			"test": "test",
		},
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}
	go func() {
		if err = server.ListenAndServe("tcp", "127.0.0.1:38401"); err != nil {
			panic(err)
		}
	}()

	testFile := "test.txt"
	err = os.WriteFile(testFile, []byte("scanme.sh"), 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(testFile); err != nil {
			log.Printf("could not remove test file: %s\n", err)
		}
	}()

	var got bool

	options := runner.Options{
		HostsFile: testFile,
		Ports:     "80",
		ScanType:  "c",
		Proxy:     "127.0.0.1:38401",
		ProxyAuth: "test:test",
		OnResult: func(hr *result.HostResult) {
			got = true
		},
		WarmUpTime: 2,
		Timeout:    10 * time.Second,
	}

	vornexRunner, err := runner.NewRunner(&options)
	if err != nil {
		return err
	}
	defer func() {
		if err := vornexRunner.Close(); err != nil {
			log.Printf("could not close vornex runner: %s\n", err)
		}
	}()

	if err = vornexRunner.RunEnumeration(context.TODO()); err != nil {
		return err
	}
	if !got {
		return errors.New("no results found")
	}

	return nil
}
