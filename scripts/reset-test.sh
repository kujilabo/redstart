migrate -database 'postgres://user:password@127.0.0.1:5433/postgres?sslmode=disable' -source file://sqls/postgres/ drop
migrate -database 'postgres://user:password@127.0.0.1:5433/postgres?sslmode=disable' -source file://sqls/postgres/ up

migrate -database 'mysql://user:password@tcp(127.0.0.1:3307)/testdb' -source file://sqls/mysql/ drop
migrate -database 'mysql://user:password@tcp(127.0.0.1:3307)/testdb' -source file://sqls/mysql/ up