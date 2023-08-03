package orm

import (
	"context"
	"database/sql"
)

// Querier 用于 Select 语句
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)

	// 这种设计形态也可以
	//Get(ctx context.Context) (T, error)
	//GetMulti(ctx context.Context) ([]T, error)
}

// Executor 用于 INSERT, DELEIE 和 UPDATE
type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

type QueryBuilder interface {
	Build() (*Query, error)

	// 这样也可以(顾虑 在AOP的时候修改这个Query 修改不了)
	//Build() (Query, error)
}

type Query struct {
	SQL  string
	Args []any
}

type TableNane interface {
	TableName() string
}

/*
	面试要点(一)
		这种面试一般面的都是 ORM 实现原理了:
			ORM 框架是怎么将一个结构体映射为一张表的(或者反过来)?
					核心就是依赖于元数据，元数据描述了两者之间的映射关系
			ORM 的元数据有什么用?
					在构造 SQL 的时候，用来将 G 类型映射为表，在处理结果集的时候，用来将表映射为 Go结构体。
			ORM 的元数据一般包含什么?
					一般包含表信息、列信息、索引信息。在支持关联关系的时候，还包含表之间的关联关系。
			ORM 的表信息包含什么?
					主要就是表级别上的配置，例如表名。如果 ORM 本身支持分库分表，那么还包含分库分表信息。
			ORM 的列信息包含什么?
					主要就是列名、类型 (和对应的 Go 类型)、索引、是否主键，以及关联关系。
			ORM 的索引信息包含什么?
					主要就是每一个索引的列，以及是否唯一。


	面试要点(二)
		ORM 如何获得模型信息?
			主要是利用反射来解析 Go 类型，同时可以利用 Tag，或者暴露编程接口，允许用户额外定制模型(例如指定表名)。
		Go 字段上的 Tag (标签)有什么用?
			用来描述字段本身的额外信息，例如使用 jsn 来指示转化 json之后的字段名字，或者如 GORM 使用 Tag 来指定列的名字、索引等。这种问题可能出在面试官问 Go语法上。
		GORM (Beego) 是如何实现的? 只要回答构造 SQL + 处理结果集 + 元数据就可以了。剩下的可能就是进一步问SQL怎么构造, 以及结果集是如何被处理的.

*/
