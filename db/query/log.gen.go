// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/Fonzeca/Trackin/db/model"
)

func newLog(db *gorm.DB) log {
	_log := log{}

	_log.logDo.UseDB(db)
	_log.logDo.UseModel(&model.Log{})

	tableName := _log.logDo.TableName()
	_log.ALL = field.NewField(tableName, "*")
	_log.ID = field.NewInt32(tableName, "id")
	_log.Imei = field.NewString(tableName, "imei")
	_log.ProtocolType = field.NewInt32(tableName, "protocol_type")
	_log.Latitud = field.NewFloat64(tableName, "latitud")
	_log.Longitud = field.NewFloat64(tableName, "longitud")
	_log.Date = field.NewTime(tableName, "date")
	_log.Speed = field.NewFloat32(tableName, "speed")
	_log.AnalogInput1 = field.NewFloat32(tableName, "analog_input_1")
	_log.DeviceTemp = field.NewInt32(tableName, "device_temp")
	_log.Mileage = field.NewInt32(tableName, "mileage")
	_log.IsGps = field.NewBool(tableName, "is_gps")
	_log.IsHistory = field.NewBool(tableName, "is_history")
	_log.EngineStatus = field.NewBool(tableName, "engine_status")
	_log.Azimuth = field.NewInt32(tableName, "azimuth")
	_log.Payload = field.NewString(tableName, "payload")

	_log.fillFieldMap()

	return _log
}

type log struct {
	logDo

	ALL          field.Field
	ID           field.Int32
	Imei         field.String
	ProtocolType field.Int32
	Latitud      field.Float64
	Longitud     field.Float64
	Date         field.Time
	Speed        field.Float32
	AnalogInput1 field.Float32
	DeviceTemp   field.Int32
	Mileage      field.Int32
	IsGps        field.Bool
	IsHistory    field.Bool
	EngineStatus field.Bool
	Azimuth      field.Int32
	Payload      field.String

	fieldMap map[string]field.Expr
}

func (l log) Table(newTableName string) *log {
	l.logDo.UseTable(newTableName)
	return l.updateTableName(newTableName)
}

func (l log) As(alias string) *log {
	l.logDo.DO = *(l.logDo.As(alias).(*gen.DO))
	return l.updateTableName(alias)
}

func (l *log) updateTableName(table string) *log {
	l.ALL = field.NewField(table, "*")
	l.ID = field.NewInt32(table, "id")
	l.Imei = field.NewString(table, "imei")
	l.ProtocolType = field.NewInt32(table, "protocol_type")
	l.Latitud = field.NewFloat64(table, "latitud")
	l.Longitud = field.NewFloat64(table, "longitud")
	l.Date = field.NewTime(table, "date")
	l.Speed = field.NewFloat32(table, "speed")
	l.AnalogInput1 = field.NewFloat32(table, "analog_input_1")
	l.DeviceTemp = field.NewInt32(table, "device_temp")
	l.Mileage = field.NewInt32(table, "mileage")
	l.IsGps = field.NewBool(table, "is_gps")
	l.IsHistory = field.NewBool(table, "is_history")
	l.EngineStatus = field.NewBool(table, "engine_status")
	l.Azimuth = field.NewInt32(table, "azimuth")
	l.Payload = field.NewString(table, "payload")

	l.fillFieldMap()

	return l
}

func (l *log) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := l.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (l *log) fillFieldMap() {
	l.fieldMap = make(map[string]field.Expr, 15)
	l.fieldMap["id"] = l.ID
	l.fieldMap["imei"] = l.Imei
	l.fieldMap["protocol_type"] = l.ProtocolType
	l.fieldMap["latitud"] = l.Latitud
	l.fieldMap["longitud"] = l.Longitud
	l.fieldMap["date"] = l.Date
	l.fieldMap["speed"] = l.Speed
	l.fieldMap["analog_input_1"] = l.AnalogInput1
	l.fieldMap["device_temp"] = l.DeviceTemp
	l.fieldMap["mileage"] = l.Mileage
	l.fieldMap["is_gps"] = l.IsGps
	l.fieldMap["is_history"] = l.IsHistory
	l.fieldMap["engine_status"] = l.EngineStatus
	l.fieldMap["azimuth"] = l.Azimuth
	l.fieldMap["payload"] = l.Payload
}

func (l log) clone(db *gorm.DB) log {
	l.logDo.ReplaceDB(db)
	return l
}

type logDo struct{ gen.DO }

func (l logDo) Debug() *logDo {
	return l.withDO(l.DO.Debug())
}

func (l logDo) WithContext(ctx context.Context) *logDo {
	return l.withDO(l.DO.WithContext(ctx))
}

func (l logDo) ReadDB(ctx context.Context) *logDo {
	return l.WithContext(ctx).Clauses(dbresolver.Read)
}

func (l logDo) WriteDB(ctx context.Context) *logDo {
	return l.WithContext(ctx).Clauses(dbresolver.Write)
}

func (l logDo) Clauses(conds ...clause.Expression) *logDo {
	return l.withDO(l.DO.Clauses(conds...))
}

func (l logDo) Returning(value interface{}, columns ...string) *logDo {
	return l.withDO(l.DO.Returning(value, columns...))
}

func (l logDo) Not(conds ...gen.Condition) *logDo {
	return l.withDO(l.DO.Not(conds...))
}

func (l logDo) Or(conds ...gen.Condition) *logDo {
	return l.withDO(l.DO.Or(conds...))
}

func (l logDo) Select(conds ...field.Expr) *logDo {
	return l.withDO(l.DO.Select(conds...))
}

func (l logDo) Where(conds ...gen.Condition) *logDo {
	return l.withDO(l.DO.Where(conds...))
}

func (l logDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *logDo {
	return l.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (l logDo) Order(conds ...field.Expr) *logDo {
	return l.withDO(l.DO.Order(conds...))
}

func (l logDo) Distinct(cols ...field.Expr) *logDo {
	return l.withDO(l.DO.Distinct(cols...))
}

func (l logDo) Omit(cols ...field.Expr) *logDo {
	return l.withDO(l.DO.Omit(cols...))
}

func (l logDo) Join(table schema.Tabler, on ...field.Expr) *logDo {
	return l.withDO(l.DO.Join(table, on...))
}

func (l logDo) LeftJoin(table schema.Tabler, on ...field.Expr) *logDo {
	return l.withDO(l.DO.LeftJoin(table, on...))
}

func (l logDo) RightJoin(table schema.Tabler, on ...field.Expr) *logDo {
	return l.withDO(l.DO.RightJoin(table, on...))
}

func (l logDo) Group(cols ...field.Expr) *logDo {
	return l.withDO(l.DO.Group(cols...))
}

func (l logDo) Having(conds ...gen.Condition) *logDo {
	return l.withDO(l.DO.Having(conds...))
}

func (l logDo) Limit(limit int) *logDo {
	return l.withDO(l.DO.Limit(limit))
}

func (l logDo) Offset(offset int) *logDo {
	return l.withDO(l.DO.Offset(offset))
}

func (l logDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *logDo {
	return l.withDO(l.DO.Scopes(funcs...))
}

func (l logDo) Unscoped() *logDo {
	return l.withDO(l.DO.Unscoped())
}

func (l logDo) Create(values ...*model.Log) error {
	if len(values) == 0 {
		return nil
	}
	return l.DO.Create(values)
}

func (l logDo) CreateInBatches(values []*model.Log, batchSize int) error {
	return l.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (l logDo) Save(values ...*model.Log) error {
	if len(values) == 0 {
		return nil
	}
	return l.DO.Save(values)
}

func (l logDo) First() (*model.Log, error) {
	if result, err := l.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Log), nil
	}
}

func (l logDo) Take() (*model.Log, error) {
	if result, err := l.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Log), nil
	}
}

func (l logDo) Last() (*model.Log, error) {
	if result, err := l.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Log), nil
	}
}

func (l logDo) Find() ([]*model.Log, error) {
	result, err := l.DO.Find()
	return result.([]*model.Log), err
}

func (l logDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Log, err error) {
	buf := make([]*model.Log, 0, batchSize)
	err = l.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (l logDo) FindInBatches(result *[]*model.Log, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return l.DO.FindInBatches(result, batchSize, fc)
}

func (l logDo) Attrs(attrs ...field.AssignExpr) *logDo {
	return l.withDO(l.DO.Attrs(attrs...))
}

func (l logDo) Assign(attrs ...field.AssignExpr) *logDo {
	return l.withDO(l.DO.Assign(attrs...))
}

func (l logDo) Joins(fields ...field.RelationField) *logDo {
	for _, _f := range fields {
		l = *l.withDO(l.DO.Joins(_f))
	}
	return &l
}

func (l logDo) Preload(fields ...field.RelationField) *logDo {
	for _, _f := range fields {
		l = *l.withDO(l.DO.Preload(_f))
	}
	return &l
}

func (l logDo) FirstOrInit() (*model.Log, error) {
	if result, err := l.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Log), nil
	}
}

func (l logDo) FirstOrCreate() (*model.Log, error) {
	if result, err := l.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Log), nil
	}
}

func (l logDo) FindByPage(offset int, limit int) (result []*model.Log, count int64, err error) {
	result, err = l.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = l.Offset(-1).Limit(-1).Count()
	return
}

func (l logDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = l.Count()
	if err != nil {
		return
	}

	err = l.Offset(offset).Limit(limit).Scan(result)
	return
}

func (l *logDo) withDO(do gen.Dao) *logDo {
	l.DO = *do.(*gen.DO)
	return l
}
