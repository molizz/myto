package convertor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseMysqlTableOptions(t *testing.T) {
	test := ` ENGINE="InnoDB" default charset=latin1 collate=latin1_bin comment='标签库'`

	result := parseMysqlTableOptions(test)
	assert.Equal(t, 5, len(result.options))
	assert.Equal(t, "InnoDB", result.options["ENGINE"])
	assert.Equal(t, "标签库", result.options["comment"])
}
