package main

import (
	"reflect"
	"testing"
)

func Test_convertEvents(t *testing.T) {
	type args struct {
		eventsString string
	}
	tests := []struct {
		name       string
		args       args
		wantEvents int
	}{
		{"test1", args{"CREATE"}, 1},
		{"test2", args{"WRITE"}, 2},
		{"test3", args{"REMOVE"}, 4},
		{"test4", args{"RENAME"}, 8},
		{"test5", args{"test"}, 0},
		{"test5", args{"CREATE WRITE"}, 3},
		{"test5", args{"CREATE WRITE RENAME"}, 11},
		{"test5", args{"CREATE WRITE RENAME REMOVE"}, 15},
		{"test6", args{"CREATE WRITE RENAME REMOVE TEST"}, 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEvents := convertEvents(tt.args.eventsString); gotEvents != tt.wantEvents {
				t.Errorf("convertEvents() = %v, want %v", gotEvents, tt.wantEvents)
			}
		})
	}
}

func Test_deduplicate(t *testing.T) {
	type args struct {
		watchList []string
	}
	tests := []struct {
		name        string
		args        args
		wantUniques []string
	}{
		{"test1", args{[]string{"watcher.go", "../watcher/watcher.go"}}, []string{"../watcher/watcher.go"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotUniques := deduplicate(tt.args.watchList); !reflect.DeepEqual(gotUniques, tt.wantUniques) {
				t.Errorf("deduplicate() = %v, want %v", gotUniques, tt.wantUniques)
			}
		})
	}
}

func Test_parseFiles(t *testing.T) {
	type args struct {
		files string
	}
	tests := []struct {
		name          string
		args          args
		wantWatchList []string
	}{
		{name: "test1", args: args{files: "watcher.go"}, wantWatchList: []string{"watcher.go"}},
		{name: "test2", args: args{files: "*.go"}, wantWatchList: []string{"utils.go", "utils_test.go", "watcher.go", "watcher_test.go"}},
		{name: "test3", args: args{files: "tester.go"}, wantWatchList: nil},
		{name: "test4", args: args{files: "**/*.go"}, wantWatchList: nil},
		{name: "test5", args: args{files: "*!/*.go"}, wantWatchList: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWatchList := parseFiles(tt.args.files); !reflect.DeepEqual(gotWatchList, tt.wantWatchList) {
				t.Errorf("parseFiles() = %v, want %v", gotWatchList, tt.wantWatchList)
			}
		})
	}
}
