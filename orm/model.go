package orm

type model struct {
	tableName string
	// 表名
	fields map[string]field
}

type field struct {
	// 列名
	colName string
}

func parseModel(entity any) (*model, error) {

}
