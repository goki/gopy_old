// Copyright 2015 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func init() {
	var (
		py2   = "python2"
		py3   = "python3"
		pypy2 = "pypy"
		pypy3 = "pypy3"
	)

	if os.Getenv("GOPY_APPVEYOR_CI") == "1" {
		log.Printf("Running in appveyor CI")
		var (
			cpy2dir  = os.Getenv("CPYTHON2DIR")
			cpy3dir  = os.Getenv("CPYTHON3DIR")
			pypy2dir = os.Getenv("PYPY2DIR")
			pypy3dir = os.Getenv("PYPY3DIR")
		)
		py2 = path.Join(cpy2dir, "python")
		py3 = path.Join(cpy3dir, "python")
		pypy2 = path.Join(pypy2dir, "pypy")
		pypy3 = path.Join(pypy3dir, "pypy")
	}

	var (
		disabled []string
		missing  int
	)
	for _, be := range []struct {
		name      string
		vm        string
		module    string
		mandatory bool
	}{
		//		{"py2", py2, "", true},
		{"py2-cffi", py2, "cffi", true},
		//		{"py3", py3, "", true},
		{"py3-cffi", py3, "cffi", true},
		{"pypy2-cffi", pypy2, "cffi", false},
		{"pypy3-cffi", pypy3, "cffi", false},
	} {
		args := []string{"-c", ""}
		if be.module != "" {
			args[1] = "import " + be.module
		}
		log.Printf("checking testbackend: %q...", be.name)
		cmd := exec.Command(be.vm, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("disabling testbackend: %q, error: '%s'", be.name, err.Error())
			testBackends[be.name] = ""
			disabled = append(disabled, be.name)
			if be.mandatory {
				missing++
			}
		} else {
			log.Printf("enabling testbackend: %q", be.name)
			testBackends[be.name] = be.vm
		}
	}

	if len(disabled) > 0 {
		log.Printf("The following test backends are not available: %s",
			strings.Join(disabled, ", "))
		if os.Getenv("GOPY_APPVEYOR_CI") == "1" && missing > 0 {
			log.Fatalf("Not all backends available in appveyor CI")
		}
	}
}
