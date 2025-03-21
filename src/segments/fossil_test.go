package segments

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jandedobbeleer/oh-my-posh/src/properties"
	"github.com/jandedobbeleer/oh-my-posh/src/runtime/mock"

	"github.com/stretchr/testify/assert"
)

func TestFossilStatus(t *testing.T) {
	cases := []struct {
		OutputError      error
		Case             string
		Output           string
		ExpectedStatus   string
		ExpectedBranch   string
		HasCommand       bool
		ExpectedDisabled bool
	}{
		{
			Case:             "not installed",
			HasCommand:       false,
			ExpectedDisabled: true,
		},
		{
			Case:             "command error",
			HasCommand:       true,
			OutputError:      fmt.Errorf("error"),
			ExpectedDisabled: true,
		},
		{
			Case:       "default status",
			HasCommand: true,
			Output: `
			repository:   /Users/jan/Downloads/myclone.fossil
			local-root:   /Users/jan/Projects/fossil/
			config-db:    /Users/jan/.config/fossil.db
			checkout:     0fabc4f3566c7e7d9e528b17253de42e14dd5c7b 2022-06-05 04:06:17 UTC
			parent:       e8a051e6a943a26c9c33a30df8ceda069c06c174 2022-06-04 23:09:02 UTC
			tags:         trunk
			comment:      In the /setup_skin page, add a mention of/link to /skins, per request in the forum. (user: stephan)
			CONFLICT             test.tst
			DELETED	             test.tst
			MISSING              test.tst
			ADDED                test.tst
			ADDED_BY_INTEGRATE   test.tst
			ADDED_BY_MERGE       test.tst
			EDITED               auto.def
			UPDATED              test.tst
			UPDATED_BY_INTEGRATE test.tst
			UPDATED_BY_MERGE 	 test.tst
			CHANGED 	         test.tst
			RENAMED 		     test.tst
			`,
			ExpectedBranch: "trunk",
			ExpectedStatus: "+3 ~5 -2 >1 !1",
		},
	}
	for _, tc := range cases {
		env := new(mock.Environment)
		env.On("GOOS").Return("unix")
		env.On("IsWsl").Return(false)
		env.On("InWSLSharedDrive").Return(false)
		env.On("HasCommand", FOSSILCOMMAND).Return(tc.HasCommand)
		env.On("RunCommand", FOSSILCOMMAND, []string{"status"}).Return(strings.ReplaceAll(tc.Output, "\t", ""), tc.OutputError)

		f := &Fossil{}
		f.Init(properties.Map{}, env)

		got := f.Enabled()

		assert.Equal(t, !tc.ExpectedDisabled, got, tc.Case)
		if tc.ExpectedDisabled {
			continue
		}
		assert.Equal(t, tc.ExpectedStatus, f.Status.String(), tc.Case)
		assert.Equal(t, tc.ExpectedBranch, f.Branch, tc.Case)
	}
}
