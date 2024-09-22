package ddlstructdiff

import (
	"go/token"
	"strings"
)

type Field struct {
	name   string
	strict bool
}

func NewField(name string, strict bool) *Field {
	return &Field{
		name:   name,
		strict: strict,
	}
}

func (f *Field) Name() string {
	if f.strict {
		return f.name
	}
	return strings.ToLower(f.name)
}

func (f *Field) OriginalName() string {
	return f.name
}

type Struct struct {
	name   string
	pos    token.Pos
	strict bool
	s      []*Field
	m      map[string]*Field
}

func NewStruct(name string, pos token.Pos, strict bool) *Struct {
	return &Struct{
		name:   name,
		pos:    pos,
		strict: strict,
		s:      []*Field{},
		m:      map[string]*Field{},
	}
}

func (s *Struct) Name() string {
	if s.strict {
		return s.name
	}
	return strings.ToLower(s.name)
}

func (s *Struct) OriginalName() string {
	return s.name
}

func (s *Struct) Pos() token.Pos {
	return s.pos
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

type Structs struct {
	s []*Struct
	m map[string]*Struct
}

func NewStructs() *Structs {
	return &Structs{
		s: []*Struct{},
		m: map[string]*Struct{},
	}
}

func (s Structs) AddStruct(st *Struct) {
	s.s = append(s.s, st)
	s.m[st.Name()] = st
}

func (s Structs) Struct(name string) (*Struct, bool) {
	st, ok := s.m[name]
	return st, ok
}

func (s Structs) Structs() []*Struct {
	return s.s
}
