package registry

import "fmt"

type IDGenerator struct {
	namespace string
	cursor    int
}

func NewIDGenerator(namespace string) *IDGenerator {
	return &IDGenerator{namespace: namespace, cursor: 1}
}

func (g *IDGenerator) Next(kind string) string {
	value := fmt.Sprintf("%s-%s-%03d", g.namespace, kind, g.cursor)
	g.cursor++
	return value
}

func (g *IDGenerator) Reset() { g.cursor = 1 }
