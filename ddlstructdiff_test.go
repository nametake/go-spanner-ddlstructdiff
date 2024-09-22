package ddlstructdiff_test

import (
	"path/filepath"
	"testing"

	ddlstructdiff "github.com/nametake/go-spanner-ddlstructdiff"

	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	tests := []struct {
		name     string
		ddl      string
		strict   bool
		patterns []string
	}{
		{
			name:     "nofield",
			ddl:      "testdata/src/nofield/ddl.sql",
			strict:   false,
			patterns: []string{"nofield"},
		},
		{
			name:     "withtags",
			ddl:      "testdata/src/withtags/ddl.sql",
			strict:   false,
			patterns: []string{"withtags"},
		},
		{
			name:     "nocolumn",
			ddl:      "testdata/src/nocolumn/ddl.sql",
			strict:   false,
			patterns: []string{"nocolumn"},
		},
		{
			name:     "lowercase",
			ddl:      "testdata/src/lowercase/ddl.sql",
			strict:   false,
			patterns: []string{"lowercase"},
		},
		{
			name:     "lowercasestrict",
			ddl:      "testdata/src/lowercasestrict/ddl.sql",
			strict:   true,
			patterns: []string{"lowercasestrict"},
		},
		// TODO: token.NoPos is not supported in analysistest
		// {
		// 	name:     "notable",
		// 	ddl:      "testdata/src/notable/ddl.sql",
		// 	patterns: []string{"notable"},
		// },
	}
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultPath, err := filepath.Abs(tt.ddl)
			if err != nil {
				t.Error(err)
				return
			}
			if err := ddlstructdiff.Analyzer.Flags.Set("ddl", defaultPath); err != nil {
				t.Error(err)
				return
			}
			if tt.strict {
				if err := ddlstructdiff.Analyzer.Flags.Set("strict", "true"); err != nil {
					t.Error(err)
					return
				}
			}

			analysistest.Run(t, testdata, ddlstructdiff.Analyzer, tt.patterns...)
		})
	}
}
