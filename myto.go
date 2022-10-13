package myto

import (
	"github.com/molizz/myto/convertor"
	"github.com/xwb1989/sqlparser"
)

type Convertor interface {
	Exec() (string, error)
}

type Myto struct {
	isDDL        bool
	sql          string
	sqlTokenizer *sqlparser.Tokenizer
}

func New(sql string, isDDL bool) *Myto {
	sqlTokenizer := sqlparser.NewStringTokenizer(sql)
	return &Myto{
		sql:          sql,
		isDDL:        isDDL,
		sqlTokenizer: sqlTokenizer,
	}
}

// ToDMDB 达梦数据库
func (m *Myto) ToDMDB() (string, error) {
	var conv Convertor = convertor.NewDMDB(m.sqlTokenizer)
	return conv.Exec()
}
