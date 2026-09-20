package main

import (
	"log/slog"
	"os"

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
				slog.Error("Expected string", "node", node)
				os.Exit(1)
			}
			values = append(values, stringNode.Value)
		}
		*s = values
		return nil
	default:
		slog.Error("Expected string or list of strings", "node", node)
		os.Exit(1)
	}
	return nil
}
