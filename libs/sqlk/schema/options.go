package schema

type options struct {
	tableName        string
	columns          []string
	primaryKey       string
	autoIncrement    bool
	modelRef         interface{}
	as               string
	softDelete       bool
	softDeleteColumn string
	auditFields      bool
	createdAtColumn  string
	updatedAtColumn  string
	createdByColumn  string
	updatedByColumn  string
}

var defaultOptions = &options{
	tableName:        "",
	columns:          nil,
	primaryKey:       "id",
	autoIncrement:    true,
	modelRef:         nil,
	as:               "",
	softDelete:       false,
	softDeleteColumn: "deleted_at",
	auditFields:      false,
	createdAtColumn:  "created_at",
	updatedAtColumn:  "updated_at",
	createdByColumn:  "created_by",
	updatedByColumn:  "updated_by",
}

type OptionSetterFn func(*options)

func evaluateSchemaOptions(opts []OptionSetterFn) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// AutoIncrement set auto increment options that will affect insert query generator
func AutoIncrement(ai bool) OptionSetterFn {
	return func(o *options) {
		o.autoIncrement = ai
	}
}

// PrimaryKey set columns that will be primary key, otherwise it will use "id" column
func PrimaryKey(pk string) OptionSetterFn {
	return func(o *options) {
		o.primaryKey = pk
	}
}

// FromModelRef reflect referenced model as Schema
func FromModelRef(m interface{}) OptionSetterFn {
	return func(o *options) {
		o.modelRef = m
	}
}

// TableName set table name of schema
func TableName(s string) OptionSetterFn {
	return func(o *options) {
		o.tableName = s
	}
}

// Columns set columns in a schema
func Columns(cols ...string) OptionSetterFn {
	return func(o *options) {
		o.columns = cols
	}
}

// As set table alias to schema, will be use as reference if set
func As(as string) OptionSetterFn {
	return func(o *options) {
		o.as = as
	}
}

// SoftDelete enable/disable soft delete for schema
func SoftDelete(enabled bool) OptionSetterFn {
	return func(o *options) {
		o.softDelete = enabled
	}
}

// SoftDeleteColumn set custom column name for soft delete (default: "deleted_at")
func SoftDeleteColumn(col string) OptionSetterFn {
	return func(o *options) {
		o.softDeleteColumn = col
	}
}

// AuditFields enable/disable audit fields for schema
func AuditFields(enabled bool) OptionSetterFn {
	return func(o *options) {
		o.auditFields = enabled
	}
}

// CreatedAtColumn set custom column name for created_at (default: "created_at")
func CreatedAtColumn(col string) OptionSetterFn {
	return func(o *options) {
		o.createdAtColumn = col
	}
}

// UpdatedAtColumn set custom column name for updated_at (default: "updated_at")
func UpdatedAtColumn(col string) OptionSetterFn {
	return func(o *options) {
		o.updatedAtColumn = col
	}
}

// CreatedByColumn set custom column name for created_by (default: "created_by")
func CreatedByColumn(col string) OptionSetterFn {
	return func(o *options) {
		o.createdByColumn = col
	}
}

// UpdatedByColumn set custom column name for updated_by (default: "updated_by")
func UpdatedByColumn(col string) OptionSetterFn {
	return func(o *options) {
		o.updatedByColumn = col
	}
}
