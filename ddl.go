package ddlstructdiff

import (
	"fmt"
	"io"
	"strings"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
)

type Column struct {
	name string
}

func NewColumn(name string) *Column {
	return &Column{
		name: name,
	}
}

func (c *Column) Name() string {
	return strings.ToLower(c.name)
}

type Table struct {
	name string
	s    []*Column
	m    map[string]*Column
}

func NewTable(name string) Table {
	return Table{
		name: name,
		s:    []*Column{},
		m:    map[string]*Column{},
	}
}

func (t *Table) Name() string {
	return strings.ToLower(t.name)
}

func (t *Table) Columns() []*Column {
	return t.s
}

func (t *Table) Column(column string) (*Column, bool) {
	c, ok := t.m[column]
	return c, ok
}

func (t *Table) AddColumn(c *Column) {
	t.s = append(t.s, c)
	t.m[c.Name()] = c
}

type DDL struct {
	s []*Table
	m map[string]Table
}

func NewDDL() *DDL {
	return &DDL{
		s: []*Table{},
		m: map[string]Table{},
	}
}

func (d *DDL) Table(table string) (Table, bool) {
	t, ok := d.m[table]
	return t, ok
}

func (d DDL) AddTable(t Table) {
	d.m[t.Name()] = t
}

func loadDDL(r io.Reader) (*DDL, error) {
	ddlReader, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read SQL file: %w", err)
	}

	file := &token.File{
		Buffer: string(ddlReader),
	}

	p := memefish.Parser{
		Lexer: &memefish.Lexer{File: file},
	}

	stmt, err := p.ParseDDLs()
	if err != nil {
		return nil, fmt.Errorf("failed to parse DDL: %w", err)
	}

	ddl := NewDDL()
	for _, s := range stmt {
		ct, ok := s.(*ast.CreateTable)
		if !ok {
			continue
		}
		table := NewTable(ct.Name.Name)
		for _, c := range ct.Columns {
			table.AddColumn(NewColumn(c.Name.Name))
		}
		ddl.AddTable(table)
	}

	return ddl, nil
}
