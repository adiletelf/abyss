package object

import (
	"fmt"
)

const (
	NULL_OBJ    = "NULL"
	INTEGET_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGET_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
