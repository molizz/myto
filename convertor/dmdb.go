package convertor

import (
	"fmt"
	"io"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/xwb1989/sqlparser"
)

var _ Element = (*dmdbDropTableIfExists)(nil)
var _ Element = (*dmdbCreateTable)(nil)
var _ Element = (*dmdbTableColumn)(nil)
var _ Element = (*dmdbColumnComment)(nil)

var mysqlWithDMDatatypeMapping = map[string]string{
	"varchar":   "varchar2",
	"varbinary": "longvarbinary",
	"char":      "char",
	"binary":    "binary",

	"int":       "int",
	"integer":   "int",
	"bigint":    "bigint",
	"bit":       "bit",
	"tinyint":   "int",
	"smallint":  "smallint",
	"mediumint": "int",
	"decimal":   "numeric",
	"dec":       "dec",
	"float":     "float",
	"double":    "double",

	"text":       "text",
	"longtext":   "text",
	"tinyblob":   "blob",
	"tinytext":   "varchar2",
	"blob":       "blob",
	"mediumblob": "blob",
	"mediumtext": "text",
	"longblob":   "blob",
	"bool":       "boolean",
	"boolean":    "boolean",

	"date":     "datetime",
	"datetime": "datetime",

	"enum": "enum",
	"set":  "set",
}

type DMDB struct {
	sqlTokenizer *sqlparser.Tokenizer
}

func NewDMDB(sqlTokenizer *sqlparser.Tokenizer) *DMDB {
	return &DMDB{sqlTokenizer: sqlTokenizer}
}

func (o *DMDB) Exec() (string, error) {
	var container = NewContainerWithSuffix("\n", true)

	for {
		st, err := sqlparser.ParseNext(o.sqlTokenizer)
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		// TODO view table
		switch ddl := st.(type) {
		case *sqlparser.DDL:
			switch ddl.Action {
			case sqlparser.DropStr:
				container.Append(&dmdbDropTableIfExists{DDL: ddl})
			case sqlparser.CreateStr:
				container.Append(&dmdbCreateTable{
					DDL:                     ddl,
					columnContainer:         NewContainerWithSuffix(",\n", true),
					columnCommentsContainer: NewContainerWithSuffix("\n", false),
					indexContainer:          NewContainerWithSuffix("\n", false),
				})
			}
		}
	}
	return container.Render(), nil
}

type dmdbCreateTable struct {
	*sqlparser.DDL
	columnContainer         *Container // 列
	columnCommentsContainer *Container // 列注释
	indexContainer          *Container
	sb                      strings.Builder
}

func (o *dmdbCreateTable) Format() string {
	tableName := o.NewName.Name.String()

	for _, column := range o.DDL.TableSpec.Columns {
		o.columnContainer.Append(&dmdbTableColumn{ColumnDefinition: column})
		// 生成表中的字段注释
		if column.Type.Comment != nil {
			o.columnCommentsContainer.Append(&dmdbColumnComment{
				tableName:        tableName,
				ColumnDefinition: column,
			})
		}
	}
	for _, index := range o.DDL.TableSpec.Indexes {
		o.indexContainer.Append(&dmdbTableIndex{
			tableName:       tableName,
			IndexDefinition: index,
		})
	}

	o.sb.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tableName))
	o.sb.WriteString(o.columnContainer.Render())
	o.sb.WriteString(");\n")

	// table index
	o.sb.WriteString(o.indexContainer.Render())

	// table comment
	opt := parseMysqlTableOptions(o.DDL.TableSpec.Options)
	if comment, found := opt.options["comment"]; found {
		o.sb.WriteString(fmt.Sprintf("COMMENT ON TABLE %v IS '%v';\n", tableName, comment))
	}

	// table column comment
	o.sb.WriteString(o.columnCommentsContainer.Render())
	return o.sb.String()
}

type dmdbTableIndex struct {
	tableName string
	*sqlparser.IndexDefinition
}

func (t *dmdbTableIndex) Format() string {
	var info = t.IndexDefinition.Info
	var indexName = t.IndexDefinition.Info.Name.String()
	var sb strings.Builder

	if info.Primary {
		// 主键索引
		_, _ = fmt.Fprintf(&sb, "ALTER TABLE %s ADD CONSTRAINT %s PRIMARY KEY (%s);",
			t.tableName, buildPKName(t.IndexDefinition.Columns), buildIndexColumns(t.IndexDefinition.Columns))
	} else if info.Unique {
		// 唯一索引
		_, _ = fmt.Fprintf(&sb, "CREATE UNIQUE INDEX %s ON %s(%s);",
			indexName, t.tableName, buildIndexColumns(t.IndexDefinition.Columns))
	} else {
		// 普通索引
		_, _ = fmt.Fprintf(&sb, "CREATE INDEX %s ON %s(%s);",
			indexName, t.tableName, buildIndexColumns(t.IndexDefinition.Columns))
	}
	return sb.String()
}

type dmdbTableColumn struct {
	*sqlparser.ColumnDefinition
}

func (o *dmdbTableColumn) Format() string {
	var sb = &strings.Builder{}

	columnName := o.ColumnDefinition.Name.String()
	columnType := o.ColumnDefinition.Type

	// column name
	if IsDMKeyword(columnName) {
		sb.WriteByte('"')
		sb.WriteString(columnName)
		sb.WriteByte('"')
	} else {
		sb.WriteString(columnName)
	}
	sb.WriteByte(' ')

	// column type name
	if t, found := mysqlWithDMDatatypeMapping[columnType.Type]; found {
		sb.WriteString(fmt.Sprintf("%s", t))
	} else {
		log.Fatalf("the mysql column mapping was not found")
		return ""
	}

	// column type
	o.formatColumnType(sb, columnType)
	sb.WriteByte(' ')

	// column default(NULL or NOT NULL)
	if columnType.NotNull {
		sb.WriteString("NOT NULL")
		sb.WriteByte(' ')
	}
	return sb.String()
}

func (o *dmdbTableColumn) formatColumnType(sb *strings.Builder, columnType sqlparser.ColumnType) {
	switch columnType.Type {
	case "varchar", "varbinary", "char", "binary",
		"tinytext":
		if columnType.Length != nil {
			num, err := strconv.ParseInt(string(columnType.Length.Val), 0, 64)
			if err != nil {
				log.Fatalf("invalid length val: %v %v", columnType.Length.Type, columnType.Length.Val)
			}
			sb.WriteString(fmt.Sprintf("(%d)", num))
		}
	case "blob", "tinyblob":
		sb.WriteString("(255)")
	case "mediumblob":
		sb.WriteString("(16777215)")
	case "longblob":
		sb.WriteString(fmt.Sprintf("(%d)", math.MaxInt32))
	case "decimal", "dec", "float", "double":
		if columnType.Length != nil && columnType.Scale != nil {
			sb.WriteString(fmt.Sprintf("(%v,%v)", columnType.Length, columnType.Scale))
		} else if columnType.Length != nil {
			sb.WriteString(fmt.Sprintf("(%v,0)", columnType.Length))
		}
	case "enum", "set":
		sb.WriteString(fmt.Sprintf("(%s)", strings.Join(columnType.EnumValues, ", ")))
	case "text", "mediumtext", "longtext",
		"boolean", "bool",
		"date", "datetime",
		"int", "integer", "bigint", "bit", "tinyint", "smallint", "mediumint":
		// ignore
	default:
		log.Fatalf("undeliverable date type '%v'", columnType)
	}
}

type dmdbColumnComment struct {
	tableName string
	*sqlparser.ColumnDefinition
}

func (d *dmdbColumnComment) Format() string {
	if d.ColumnDefinition.Type.Comment != nil {
		return fmt.Sprintf(`COMMENT ON COLUMN %s.%s IS '%s';`,
			d.tableName, d.ColumnDefinition.Name, d.ColumnDefinition.Type.Comment.Val)
	}
	return ""
}

type dmdbDropTableIfExists struct {
	*sqlparser.DDL
}

func (d *dmdbDropTableIfExists) Format() string {
	if d.IfExists {
		return fmt.Sprintf(`BEGIN
   EXECUTE IMMEDIATE 'DROP TABLE %s';
EXCEPTION
   WHEN OTHERS THEN NULL;
END;`, d.Table.Name.String())
	}
	return ""
}
