
DROP TABLE IF EXISTS `draft`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `draft` (
  `uuid` varchar(8) COLLATE latin1_bin NOT NULL COMMENT 'uuid',
  `team_uuid` varchar(8) COLLATE latin1_bin NOT NULL COMMENT 'uuid',
  `type` tinyint(4) NOT NULL COMMENT '关联数据类型',
  `name` varchar(128) CHARACTER SET utf8mb4 NOT NULL COMMENT '名称',
  `name_pinyin` varchar(1024) CHARACTER SET utf8mb4 NOT NULL COMMENT '名称拼音',
  `desc` longtext CHARACTER SET utf8mb4 COMMENT '任务描述',
  `progress` int(11) NOT NULL COMMENT '进度',
  `position` bigint(20) NOT NULL COMMENT '显示位置',
  `data` text CHARACTER SET utf8mb4 NOT NULL COMMENT '详细信息',
  `activity_uuid` varchar(8) COLLATE latin1_bin NOT NULL COMMENT ' uuid',
  `status` tinyint(4) NOT NULL COMMENT '1',
  `config_type` enum('a1', 'a2') CHARACTER SET utf8 NOT NULL,
  `config_type2` set('e1', 'e2') CHARACTER SET utf8 NOT NULL,
  PRIMARY KEY (`uuid`, `team_uuid`),
  UNIQUE KEY `team_number` (`team_uuid`,`number`),
  KEY `idx_team_chart_path` (`team_uuid`,`path`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_bin  COMMENT='标签库';
/*!40101 SET character_set_client = @saved_cs_client */;