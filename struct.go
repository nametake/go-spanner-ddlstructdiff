package ddlstructdiff

type Field struct{}

type Struct map[string]*Field

func NewEmptyStruct() Struct {
	return make(Struct)
}

func (s Struct) Field(field string) (*Field, bool) {
	f, ok := s[field]
	return f, ok
}

func (s Struct) AddField(field string, f *Field) {
	s[field] = f
}

type Structs map[string]Struct

func NewEmptyStructs() Structs {
	return make(Structs)
}

func (s Structs) AddStruct(name string, st Struct) {
	s[name] = st
}

func (s Structs) Struct(name string) (Struct, bool) {
	st, ok := s[name]
	return st, ok
}
