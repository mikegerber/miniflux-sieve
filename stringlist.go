package main

import (
	"fmt"

	"github.com/goccy/go-yaml/ast"
)

// Read a single string or a list of strings from the YAML
type StringList []string

func (s *StringList) UnmarshalYAML(node ast.Node) error {
	switch node.Type() {
	case ast.StringType:
		stringNode := node.(*ast.StringNode)
		value := stringNode.GetValue().(string)
		*s = []string{value}
		return nil
	case ast.SequenceType:
		seq := node.(*ast.SequenceNode)
		values := make([]string, 0, len(seq.Values))
		for _, node := range seq.Values {
			stringNode, ok := node.(*ast.StringNode)
			if !ok {
				return fmt.Errorf("Expected string: %v", node)
			}
			values = append(values, stringNode.Value)
		}
		*s = values
		return nil
	default:
		return fmt.Errorf("Expected string or list of strings: %v", node)
	}
	return nil
}
