package ddlstructdiff

import "go/token"

type Field struct{}

func NewField() *Field {
	return &Field{}
}

type Struct struct {
	Pos    token.Pos
	Fields map[string]*Field
}

func NewStruct(pos token.Pos) *Struct {
	return &Struct{
		Pos:    pos,
		Fields: map[string]*Field{},
	}
}

func (s *Struct) Field(field string) (*Field, bool) {
	f, ok := s.Fields[field]
	return f, ok
}

func (s *Struct) AddField(field string, f *Field) {
	s.Fields[field] = f
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
