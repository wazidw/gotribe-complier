package lib

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"log"
)

type runTpl struct {
	Image   string `json:"image"`
	File    string `json:"file"`
	Cmd     string `json:"cmd"`
	Timeout int    `json:"timeout"`
	Memory  string `json:"memory"`
	CpuSet  string `json:"cpuset"`
}

func Run(lang string) runTpl {
	var tpl runTpl
	lang = fmt.Sprintf("lib/lang/%s.json", lang)

	file, err := ioutil.ReadFile(lang)
	if err != nil {
		log.Fatalf("Some error occured while reading file. Error: %s", err)
	}
	err = json.Unmarshal(file, &tpl)
	if err != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
	}
	fmt.Println(tpl.Image)
	fmt.Printf("tpl Struct: %#v\n", tpl)
	return tpl
}
