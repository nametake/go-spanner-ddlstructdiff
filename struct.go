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

func (f *Field) OriginalName() string {
	return f.name
}

type Struct struct {
	pos  token.Pos
	name string
	s    []*Field
	m    map[string]*Field
}

func NewStruct(name string, pos token.Pos) *Struct {
	return &Struct{
		name: name,
		pos:  pos,
		s:    []*Field{},
		m:    map[string]*Field{},
	}
}

func (s *Struct) Name() string {
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
