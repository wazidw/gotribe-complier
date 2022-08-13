package lib

import (
	"testing"
)

func TestDockerRun(t *testing.T) {
	type args struct {
		lang string
		code string
		dest string
		cmd string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"python3",
			args{lang: "python:3",dest:"test.py",cmd:"python3 test.py",code:"def print_welcome(name):\n    print(\"Welcome\", name)\n \nprint_welcome(\"gotribe\")"},
			"\x01\x00\x00\x00\x00\x00\x00\x10Welcome gotribe\n",
		},
		{
			"rust",
			args{lang: "rust",dest:"test.rs",cmd:"rustc test.rs -o test\nif test -f \"./test\"; then\n./test\nfi",code:"fn main() {\n    println!(\"Hello, gotribe!\");}"},
			"\x01\x00\x00\x00\x00\x00\x00\x10Hello, gotribe!\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DockerRun(tt.args.lang, tt.args.code, tt.args.dest,tt.args.cmd); got != tt.want {
				t.Errorf("DockerRun() = %v, want %v", got, tt.want)
			}
		})
	}
}