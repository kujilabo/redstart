migrate -database 'mysql://user:password@tcp(127.0.0.1:3307)/testdb' -source file://sqls/mysql/ drop
migrate -database 'mysql://user:password@tcp(127.0.0.1:3307)/testdb' -source file://sqls/mysql/ up