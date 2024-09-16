package ddlstructdiff_test

import (
	"path/filepath"
	"testing"

	"github.com/nametake/ddlstructdiff"

	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	tests := []struct {
		name     string
		ddl      string
		patterns []string
	}{
		{
			name:     "singer",
			ddl:      "testdata/src/singer/ddl.sql",
			patterns: []string{"singer"},
		},
		{
			name:     "withtags",
			ddl:      "testdata/src/withtags/ddl.sql",
			patterns: []string{"withtags"},
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
			analysistest.Run(t, testdata, ddlstructdiff.Analyzer, tt.patterns...)
		})
	}
}
