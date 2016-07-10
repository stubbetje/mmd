package mmd

import(
	"io/ioutil"
	"github.com/ghodss/yaml"
)

type MetaData struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

type FileData struct {
	File        string `json:"file"`
	Directory   string `json:"dir"`
	Target      string `json:"as"`
}

type ContentData []FileData

type Definition struct {
	Meta      MetaData    `json:"meta"`
	Content   ContentData `json:"content"`
}

func LoadFromYamlFile( filename string ) (*Definition, error) {
	content, err := ioutil.ReadFile( filename )

	if err != nil {
		return nil,err
	}

	var definition Definition

	err = yaml.Unmarshal( content, & definition )
	if err != nil {
		return nil,err
	}

	return & definition, nil
}
