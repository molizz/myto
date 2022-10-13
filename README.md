# Myto

mysql sql convert to X sql.

- oracle / 达梦数据库

## 使用 

```golang
orcaleDDLSql, err := myto.New(sql, isDDL).ToDMDB()
fmt.Println(dmDDLSql)
```

#### cli
```shell
cat cli/test.sql | go run cli/main.go
```




## 参考文献

- https://dev.mysql.com/doc/refman/5.7/en/create-table.html
- https://dev.mysql.com/doc/refman/5.7/en/data-type-defaults.html
- https://docs.oracle.com/en/database/oracle/oracle-database/12.2/sqlrf/CREATE-INDEX.html#GUID-1F89BBC0-825F-4215-AF71-7588E31D8BFE
- https://docs.oracle.com/en/database/oracle/oracle-database/12.2/sqlrf/CREATE-TABLE.html#GUID-F9CE0CC3-13AE-4744-A43C-EAC7A71AAAB6
- https://docs.oracle.com/en/database/oracle/oracle-database/12.2/sqlrf/Data-Types.html#GUID-7B72E154-677A-4342-A1EA-C74C1EA928E6
- DM: https://eco.dameng.com/document/dm/zh-cn/faq/faq-sql-gramm#DM%20%E5%92%8C%E5%85%B6%E4%BB%96%E6%95%B0%E6%8D%AE%E5%BA%93%E7%9A%84%E5%85%BC%E5%AE%B9%E6%80%A7%E9%85%8D%E7%BD%AE

- test: https://www.sqlines.com/online 