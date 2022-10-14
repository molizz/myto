package convertor

import (
	"strings"
	"unicode"

	"github.com/xwb1989/sqlparser"
)

type TableOptions struct {
	options map[string]string
}

func parseMysqlTableOptions(options string) *TableOptions {
	options = strings.TrimSpace(options)

	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)
		}
	}
	result := map[string]string{}
	items := strings.FieldsFunc(options, f)
	for _, item := range items {
		raw := strings.SplitN(item, "=", 2)
		if len(raw) == 1 {
			result[raw[0]] = ""
		} else if len(raw) == 2 {
			value := raw[1]
			value = strings.Trim(value, "\"")
			value = strings.Trim(value, "'")
			result[raw[0]] = value
		}
	}
	return &TableOptions{result}
}

func buildPKName(tableName string, indexColumns []*sqlparser.IndexColumn) string {
	var sb strings.Builder
	sb.WriteString("pk_")
	sb.WriteString(tableName)
	sb.WriteString("_")
	sb.WriteString(joinIndexColumns(indexColumns, "_", nil))
	return sb.String()
}

func buildIdxName(prefix, tableName, oldIndexName string) string {
	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteString(tableName)
	sb.WriteString("_")
	sb.WriteString(oldIndexName)
	return sb.String()
}

func buildIndexColumns(indexColumns []*sqlparser.IndexColumn, keyFn func(columnName string) string) string {
	return joinIndexColumns(indexColumns, ", ", keyFn)
}

func joinIndexColumns(indexColumns []*sqlparser.IndexColumn, sep string, keyFn func(columnName string) string) string {
	var sb strings.Builder
	for i, col := range indexColumns {
		columnName := col.Column.String()
		if keyFn != nil {
			columnName = keyFn(columnName)
		}
		sb.WriteString(columnName)
		if i != (len(indexColumns) - 1) {
			sb.WriteString(sep)
		}
	}
	return sb.String()
}
