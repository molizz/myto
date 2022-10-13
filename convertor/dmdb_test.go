package convertor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xwb1989/sqlparser"
)

func TestOracle_Exec(t *testing.T) {
	type fields struct {
		sqlTokenizer sqlparser.Statement
	}
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "create & drop",
			sql: "\nDROP TABLE IF EXISTS `draft`;\n/*!40101 SET @saved_cs_client     = @@character_set_client */;\n/*!40101 SET character_set_client = utf8 */;\n" +
				"CREATE TABLE `draft` (\n  " +
				"`uuid` VARCHAR(8) COLLATE latin1_bin NOT NULL COMMENT 'uuid',\n  " +
				"`number` int(11) NOT NULL COMMENT '编号',\n  " +
				"`type` tinyint(4) NOT NULL COMMENT '关联数据类型',\n  " +
				"`name` varchar(128) CHARACTER SET utf8mb4 NOT NULL COMMENT '名称',\n  " +
				"`name_pinyin` varchar(1024) CHARACTER SET utf8mb4 NOT NULL COMMENT '名称拼音',\n  " +
				"`desc` longtext CHARACTER SET utf8mb4 COMMENT '描述',\n  " +
				"`owner` varchar(8) COLLATE latin1_bin NOT NULL COMMENT ' uuid',\n  " +
				"`progress` int(11) NOT NULL COMMENT '进度',\n  " +
				"`position` bigint(20) NOT NULL COMMENT '显示位置',\n " +
				"`create_time` bigint(20) NOT NULL COMMENT '创建时间',\n  " +
				"`update_time` bigint(20) NOT NULL COMMENT '更新时间',\n  " +
				"`data` text CHARACTER SET utf8mb4 NOT NULL COMMENT '详细信息',\n  " +
				"`config_type` enum('a1', 'a2') CHARACTER SET utf8 NOT NULL,\n  " +
				"`config_type2` set('e1', 'e2') CHARACTER SET utf8 NOT NULL, \n" +
				"PRIMARY KEY (`uuid`),\n  UNIQUE KEY `team_number` (`team_uuid`,`number`),\n  " +
				"KEY `idx_team_chart_path` (`team_uuid`,`chart_uuid`,`path`)\n" +
				") " +
				"ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin COMMENT='标签库';\n" +
				"/*!40101 SET character_set_client = @saved_cs_client */;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := sqlparser.NewStringTokenizer(tt.sql)

			o := NewDMDB(st)
			got, err := o.Exec()
			assert.Nil(t, err)
			assert.Equal(t, "", got)
		})
	}
}
