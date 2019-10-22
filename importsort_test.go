// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package main

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

const (
	goldFile = "testdata/test.go.gold"
	inFile   = "testdata/test.go.in"
)

func TestImportSort(t *testing.T) {
	in, err := ioutil.ReadFile(inFile)
	if err != nil {
		t.Fatal(err)
	}
	gold, err := ioutil.ReadFile(goldFile)
	if err != nil {
		t.Fatal(err)
	}
	sections := []string{"foobar", "cvshub.com/foobar"}
	if out, err := genFile(gold, sections); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(out, gold) {
		t.Errorf("importsort on %s file produced a change", goldFile)
		t.Log(string(out))
	}
	if out, err := genFile(in, sections); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(out, gold) {
		t.Errorf("importsort on %s different than gold", inFile)
		t.Log(string(out))
	}
}

func Test_sortImports(t *testing.T) {
	type args struct {
		in       []byte
		sections []string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "no changes",
			args: args{
				in:       []byte(""),
				sections: []string{},
			},
			want: []byte(""),
		},
		{
			name: "regroup",
			args: args{
				in:       []byte("\"gopkg.in/foo/bar\"\n\"encoding/json\"\n\"github.com/foo/bar\""),
				sections: []string{"github.com", "gopkg.in"},
			},
			want: []byte("\"encoding/json\"\n\n\"github.com/foo/bar\"\n\n\"gopkg.in/foo/bar\"\n"),
		},
		{
			name: "regroup with 3rd party",
			args: args{
				in:       []byte("\"gopkg.in/foo/bar\"\n\"encoding/json\"\n\"github.com/foo/bar\""),
				sections: []string{"github.com"},
			},
			want: []byte("\"encoding/json\"\n\n\"gopkg.in/foo/bar\"\n\n\"github.com/foo/bar\"\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortImports(tt.args.in, tt.args.sections); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortImports() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
