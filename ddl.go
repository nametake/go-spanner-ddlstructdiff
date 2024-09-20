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
	Name string
}

func NewColumn(name string) *Column {
	return &Column{
		Name: name,
	}
}

func (c *Column) LowerName() string {
	return strings.ToLower(c.Name)
}

type Table struct {
	s []*Column
	m map[string]*Column
}

func NewTable() Table {
	return Table{
		s: []*Column{},
		m: map[string]*Column{},
	}
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
	t.m[c.LowerName()] = c
}

type DDL map[string]Table

func (d DDL) Table(table string) (Table, bool) {
	t, ok := d[table]
	return t, ok
}

func (d DDL) AddTable(k string, t Table) {
	d[k] = t
}

func loadDDL(r io.Reader) (DDL, error) {
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

	ddl := DDL{}
	for _, s := range stmt {
		ct, ok := s.(*ast.CreateTable)
		if !ok {
			continue
		}
		table := NewTable()
		for _, c := range ct.Columns {
			table.AddColumn(NewColumn(c.Name.Name))
		}
		ddl.AddTable(ct.Name.Name, table)
	}

	return ddl, nil
}
