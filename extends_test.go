package templatex_test

import (
	. "github.com/chakrit/templatex"
	a "github.com/stretchr/testify/assert"
	"testing"
	"bytes"
	"io/ioutil"
)

const (
	OutputFilename    = "test_output.txt"
	BaseTemplateName  = "templates/base.template"
	ChildTemplateName = "templates/subfolder/main.template"
)

func Test(t *testing.T) {
	template, e := ParseFile("", ChildTemplateName)
	a.NoError(t, e)

	out := &bytes.Buffer{}
	e = template.ExecuteTemplate(out, BaseTemplateName, nil)
	a.NoError(t, e)

	expected, e := ioutil.ReadFile(OutputFilename)
	a.NoError(t, e)
	a.Equal(t, string(expected), string(out.Bytes()))
}
