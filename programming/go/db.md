---
title: DB
---

## DB

### Prepared statement と database/sql
* DB レベルでは、prepared statement はある connection に bound されている
    * The typical flow is that the client sends a SQL statement with placeholders to the server for preparation, the server responds with a statement ID, and then the client executes the statement by sending its ID and parameters
* Go(database/sql)では、そもそも connection が直では expose されていない
    * DB or Tx に対しての操作となる
        * 逆に driver レベルの詳細は隠蔽されているということになる
* 仕組みは以下
    1. When you prepare a statement, it’s prepared on a connection in the pool.
    2. The Stmt object remembers which connection was used.
    3. When you execute the Stmt, it tries to use the connection. 
    If it’s not available because it’s closed or busy doing something else, 
    it gets another connection from the pool and re-prepares the statement with the database on another connection.      
    
## Reference
* [To ORM or not to ORM](https://eli.thegreenplace.net/2019/to-orm-or-not-to-orm/)
* [Go database/sql tutorial](http://go-database-sql.org/index.html)
* [Illustrated guide to SQLX](http://jmoiron.github.io/sqlx/)