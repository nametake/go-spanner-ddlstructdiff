package ddlstructdiff

import (
	"fmt"
	"io"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
)

type Column struct{}

type Table map[string]*Column

func (t Table) Column(column string) (*Column, bool) {
	c, ok := t[column]
	return c, ok
}

type DDL map[string]Table

func (d DDL) Add(table, column string) {
	if _, ok := d[table]; !ok {
		d[table] = make(map[string]*Column)
	}
	d[table][column] = &Column{}
}

func (d DDL) Table(table string) (Table, bool) {
	t, ok := d[table]
	return t, ok
}

func (d DDL) HasTable(table string) bool {
	_, ok := d[table]
	return ok
}

func (d DDL) HasColumn(table, column string) bool {
	if _, ok := d[table]; !ok {
		return false
	}
	_, ok := d[table][column]
	return ok
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
		for _, c := range ct.Columns {
			ddl.Add(ct.Name.Name, c.Name.Name)
		}
	}

	return ddl, nil
}
