package ddlstructdiff

import (
	"go/token"
	"strings"
)

type Field struct {
	name string
}

func NewField(name string) *Field {
	return &Field{
		name: name,
	}
}

func (f *Field) Name() string {
	return strings.ToLower(f.name)
}

type Struct struct {
	Pos token.Pos
	s   []*Field
	m   map[string]*Field
}

func NewStruct(pos token.Pos) *Struct {
	return &Struct{
		Pos: pos,
		s:   []*Field{},
		m:   map[string]*Field{},
	}
}

func (s *Struct) Fields() []*Field {
	return s.s
}

func (s *Struct) Field(field string) (*Field, bool) {
	f, ok := s.m[field]
	return f, ok
}

func (s *Struct) AddField(f *Field) {
	s.s = append(s.s, f)
	s.m[f.Name()] = f
}

type Structs map[string]*Struct

func NewStructs() Structs {
	return make(Structs)
}

func (s Structs) AddStruct(name string, st *Struct) {
	s[name] = st
}

func (s Structs) Struct(name string) (*Struct, bool) {
	st, ok := s[name]
	return st, ok
}
