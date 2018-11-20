package readonly_test

import (
	"testing"

	"github.com/tenntenn/readonly"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, readonly.Analyzer, "a")
}
