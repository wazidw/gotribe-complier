package main

import (
	"fmt"
	"gotribe/compiler/lib"
)

func main() {
	lang := "rust"
	code := "fn main() {\n    println!(\"Hello, world hhh!\");\n}\n"
	tpl := lib.Run(lang)
	output := lib.DockerRun(tpl.Image, code, tpl.File, tpl.Cmd)
	fmt.Println(tpl)
	fmt.Println(output)

}
