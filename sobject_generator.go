package main

import (
	"fmt"
	"strings"
)

type SObjectGenerator struct{}

type SObjectMeta struct {
	Name   string
	Fields []*Field
}

type Field struct {
	Name    string
	Type    string
	Default string
}

func (m *SObjectMeta) GetFileName() string {
	return m.Name
}

func parseMetadata() []*SObjectMeta {
	sobjects := []*SObjectMeta{
		{
			Name: "Account",
			Fields: []*Field{
				{
					Name:    "Name",
					Type:    "String",
					Default: "hoge",
				},
			},
		},
	}
	return sobjects
}

func (g *SObjectGenerator) generate() error {
	sobjects := parseMetadata()
	for _, sobject := range sobjects {
		err := g.generateSObjectFile(sobject)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *SObjectGenerator) generateSObjectFile(meta *SObjectMeta) error {
	fields := make([]string, len(meta.Fields))
	for i, f := range meta.Fields {
		if f.Type == "String" {
			fields[i] = fmt.Sprintf("public %s %s = \"%s\";", f.Type, f.Name, f.Default)
		} else {
			fields[i] = fmt.Sprintf("public %s %s = %s;", f.Type, f.Name, f.Default)
		}
	}
	body := fmt.Sprintf(`public class %s {
%s
}`, meta.Name, strings.Join(fields, "\n"))
	fmt.Println(body)
	//err := ioutil.WriteFile(meta.GetFileName(), []byte(body), 0777)
	//if err != nil {
	//	return err
	//}
	return nil
}
