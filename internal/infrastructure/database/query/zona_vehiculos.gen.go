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

	"github.com/Fonzeca/Trackin/internal/infrastructure/database/model"
)

func newZonaVehiculo(db *gorm.DB, opts ...gen.DOOption) zonaVehiculo {
	_zonaVehiculo := zonaVehiculo{}

	_zonaVehiculo.zonaVehiculoDo.UseDB(db, opts...)
	_zonaVehiculo.zonaVehiculoDo.UseModel(&model.ZonaVehiculo{})

	tableName := _zonaVehiculo.zonaVehiculoDo.TableName()
	_zonaVehiculo.ALL = field.NewAsterisk(tableName)
	_zonaVehiculo.ID = field.NewInt32(tableName, "id")
	_zonaVehiculo.ZonaID = field.NewInt32(tableName, "zona_id")
	_zonaVehiculo.Imei = field.NewString(tableName, "imei")
	_zonaVehiculo.AvisarEntrada = field.NewBool(tableName, "avisar_entrada")
	_zonaVehiculo.AvisarSalida = field.NewBool(tableName, "avisar_salida")

	_zonaVehiculo.fillFieldMap()

	return _zonaVehiculo
}

type zonaVehiculo struct {
	zonaVehiculoDo

	ALL           field.Asterisk
	ID            field.Int32
	ZonaID        field.Int32
	Imei          field.String
	AvisarEntrada field.Bool
	AvisarSalida  field.Bool

	fieldMap map[string]field.Expr
}

func (z zonaVehiculo) Table(newTableName string) *zonaVehiculo {
	z.zonaVehiculoDo.UseTable(newTableName)
	return z.updateTableName(newTableName)
}

func (z zonaVehiculo) As(alias string) *zonaVehiculo {
	z.zonaVehiculoDo.DO = *(z.zonaVehiculoDo.As(alias).(*gen.DO))
	return z.updateTableName(alias)
}

func (z *zonaVehiculo) updateTableName(table string) *zonaVehiculo {
	z.ALL = field.NewAsterisk(table)
	z.ID = field.NewInt32(table, "id")
	z.ZonaID = field.NewInt32(table, "zona_id")
	z.Imei = field.NewString(table, "imei")
	z.AvisarEntrada = field.NewBool(table, "avisar_entrada")
	z.AvisarSalida = field.NewBool(table, "avisar_salida")

	z.fillFieldMap()

	return z
}

func (z *zonaVehiculo) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := z.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (z *zonaVehiculo) fillFieldMap() {
	z.fieldMap = make(map[string]field.Expr, 5)
	z.fieldMap["id"] = z.ID
	z.fieldMap["zona_id"] = z.ZonaID
	z.fieldMap["imei"] = z.Imei
	z.fieldMap["avisar_entrada"] = z.AvisarEntrada
	z.fieldMap["avisar_salida"] = z.AvisarSalida
}

func (z zonaVehiculo) clone(db *gorm.DB) zonaVehiculo {
	z.zonaVehiculoDo.ReplaceConnPool(db.Statement.ConnPool)
	return z
}

func (z zonaVehiculo) replaceDB(db *gorm.DB) zonaVehiculo {
	z.zonaVehiculoDo.ReplaceDB(db)
	return z
}

type zonaVehiculoDo struct{ gen.DO }

func (z zonaVehiculoDo) Debug() *zonaVehiculoDo {
	return z.withDO(z.DO.Debug())
}

func (z zonaVehiculoDo) WithContext(ctx context.Context) *zonaVehiculoDo {
	return z.withDO(z.DO.WithContext(ctx))
}

func (z zonaVehiculoDo) ReadDB() *zonaVehiculoDo {
	return z.Clauses(dbresolver.Read)
}

func (z zonaVehiculoDo) WriteDB() *zonaVehiculoDo {
	return z.Clauses(dbresolver.Write)
}

func (z zonaVehiculoDo) Session(config *gorm.Session) *zonaVehiculoDo {
	return z.withDO(z.DO.Session(config))
}

func (z zonaVehiculoDo) Clauses(conds ...clause.Expression) *zonaVehiculoDo {
	return z.withDO(z.DO.Clauses(conds...))
}

func (z zonaVehiculoDo) Returning(value interface{}, columns ...string) *zonaVehiculoDo {
	return z.withDO(z.DO.Returning(value, columns...))
}

func (z zonaVehiculoDo) Not(conds ...gen.Condition) *zonaVehiculoDo {
	return z.withDO(z.DO.Not(conds...))
}

func (z zonaVehiculoDo) Or(conds ...gen.Condition) *zonaVehiculoDo {
	return z.withDO(z.DO.Or(conds...))
}

func (z zonaVehiculoDo) Select(conds ...field.Expr) *zonaVehiculoDo {
	return z.withDO(z.DO.Select(conds...))
}

func (z zonaVehiculoDo) Where(conds ...gen.Condition) *zonaVehiculoDo {
	return z.withDO(z.DO.Where(conds...))
}

func (z zonaVehiculoDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *zonaVehiculoDo {
	return z.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (z zonaVehiculoDo) Order(conds ...field.Expr) *zonaVehiculoDo {
	return z.withDO(z.DO.Order(conds...))
}

func (z zonaVehiculoDo) Distinct(cols ...field.Expr) *zonaVehiculoDo {
	return z.withDO(z.DO.Distinct(cols...))
}

func (z zonaVehiculoDo) Omit(cols ...field.Expr) *zonaVehiculoDo {
	return z.withDO(z.DO.Omit(cols...))
}

func (z zonaVehiculoDo) Join(table schema.Tabler, on ...field.Expr) *zonaVehiculoDo {
	return z.withDO(z.DO.Join(table, on...))
}

func (z zonaVehiculoDo) LeftJoin(table schema.Tabler, on ...field.Expr) *zonaVehiculoDo {
	return z.withDO(z.DO.LeftJoin(table, on...))
}

func (z zonaVehiculoDo) RightJoin(table schema.Tabler, on ...field.Expr) *zonaVehiculoDo {
	return z.withDO(z.DO.RightJoin(table, on...))
}

func (z zonaVehiculoDo) Group(cols ...field.Expr) *zonaVehiculoDo {
	return z.withDO(z.DO.Group(cols...))
}

func (z zonaVehiculoDo) Having(conds ...gen.Condition) *zonaVehiculoDo {
	return z.withDO(z.DO.Having(conds...))
}

func (z zonaVehiculoDo) Limit(limit int) *zonaVehiculoDo {
	return z.withDO(z.DO.Limit(limit))
}

func (z zonaVehiculoDo) Offset(offset int) *zonaVehiculoDo {
	return z.withDO(z.DO.Offset(offset))
}

func (z zonaVehiculoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *zonaVehiculoDo {
	return z.withDO(z.DO.Scopes(funcs...))
}

func (z zonaVehiculoDo) Unscoped() *zonaVehiculoDo {
	return z.withDO(z.DO.Unscoped())
}

func (z zonaVehiculoDo) Create(values ...*model.ZonaVehiculo) error {
	if len(values) == 0 {
		return nil
	}
	return z.DO.Create(values)
}

func (z zonaVehiculoDo) CreateInBatches(values []*model.ZonaVehiculo, batchSize int) error {
	return z.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (z zonaVehiculoDo) Save(values ...*model.ZonaVehiculo) error {
	if len(values) == 0 {
		return nil
	}
	return z.DO.Save(values)
}

func (z zonaVehiculoDo) First() (*model.ZonaVehiculo, error) {
	if result, err := z.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.ZonaVehiculo), nil
	}
}

func (z zonaVehiculoDo) Take() (*model.ZonaVehiculo, error) {
	if result, err := z.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.ZonaVehiculo), nil
	}
}

func (z zonaVehiculoDo) Last() (*model.ZonaVehiculo, error) {
	if result, err := z.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.ZonaVehiculo), nil
	}
}

func (z zonaVehiculoDo) Find() ([]*model.ZonaVehiculo, error) {
	result, err := z.DO.Find()
	return result.([]*model.ZonaVehiculo), err
}

func (z zonaVehiculoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ZonaVehiculo, err error) {
	buf := make([]*model.ZonaVehiculo, 0, batchSize)
	err = z.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (z zonaVehiculoDo) FindInBatches(result *[]*model.ZonaVehiculo, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return z.DO.FindInBatches(result, batchSize, fc)
}

func (z zonaVehiculoDo) Attrs(attrs ...field.AssignExpr) *zonaVehiculoDo {
	return z.withDO(z.DO.Attrs(attrs...))
}

func (z zonaVehiculoDo) Assign(attrs ...field.AssignExpr) *zonaVehiculoDo {
	return z.withDO(z.DO.Assign(attrs...))
}

func (z zonaVehiculoDo) Joins(fields ...field.RelationField) *zonaVehiculoDo {
	for _, _f := range fields {
		z = *z.withDO(z.DO.Joins(_f))
	}
	return &z
}

func (z zonaVehiculoDo) Preload(fields ...field.RelationField) *zonaVehiculoDo {
	for _, _f := range fields {
		z = *z.withDO(z.DO.Preload(_f))
	}
	return &z
}

func (z zonaVehiculoDo) FirstOrInit() (*model.ZonaVehiculo, error) {
	if result, err := z.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.ZonaVehiculo), nil
	}
}

func (z zonaVehiculoDo) FirstOrCreate() (*model.ZonaVehiculo, error) {
	if result, err := z.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.ZonaVehiculo), nil
	}
}

func (z zonaVehiculoDo) FindByPage(offset int, limit int) (result []*model.ZonaVehiculo, count int64, err error) {
	result, err = z.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = z.Offset(-1).Limit(-1).Count()
	return
}

func (z zonaVehiculoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = z.Count()
	if err != nil {
		return
	}

	err = z.Offset(offset).Limit(limit).Scan(result)
	return
}

func (z zonaVehiculoDo) Scan(result interface{}) (err error) {
	return z.DO.Scan(result)
}

func (z zonaVehiculoDo) Delete(models ...*model.ZonaVehiculo) (result gen.ResultInfo, err error) {
	return z.DO.Delete(models)
}

func (z *zonaVehiculoDo) withDO(do gen.Dao) *zonaVehiculoDo {
	z.DO = *do.(*gen.DO)
	return z
}
