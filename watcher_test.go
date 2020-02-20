package main

import (
	"github.com/urfave/cli/v2"
	"testing"
)

func Test_execShell(t *testing.T) {
	type args struct {
		in0     *cli.Context
		command []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test1", args: args{in0: nil, command: []string{"echo", "test"}}, want: "test\n"},
		{name: "test2", args: args{in0: nil, command: []string{"echo test"}}, want: ""},
		{name: "test3", args: args{in0: nil, command: []string{"echo", "test", "test2", "test3"}}, want: "test test2 test3\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := execShell(tt.args.in0, tt.args.command)
			if string(out) !=  tt.want {
				t.Errorf("run() out = %s, want %s", out, tt.want)
			}
		})
	}
}
