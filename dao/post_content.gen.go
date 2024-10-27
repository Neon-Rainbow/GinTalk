// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"GinTalk/model"
)

func newPostContent(db *gorm.DB, opts ...gen.DOOption) postContent {
	_postContent := postContent{}

	_postContent.postContentDo.UseDB(db, opts...)
	_postContent.postContentDo.UseModel(&model.PostContent{})

	tableName := _postContent.postContentDo.TableName()
	_postContent.ALL = field.NewAsterisk(tableName)
	_postContent.PostID = field.NewInt64(tableName, "post_id")
	_postContent.Content = field.NewString(tableName, "content")
	_postContent.CreateTime = field.NewTime(tableName, "create_time")
	_postContent.UpdateTime = field.NewTime(tableName, "update_time")
	_postContent.DeleteTime = field.NewInt(tableName, "delete_time")

	_postContent.fillFieldMap()

	return _postContent
}

// postContent 帖子内容表：存储帖子的详细内容
type postContent struct {
	postContentDo postContentDo

	ALL        field.Asterisk
	PostID     field.Int64  // 帖子ID
	Content    field.String // 帖子内容
	CreateTime field.Time   // 帖子内容创建时间，默认当前时间
	UpdateTime field.Time   // 帖子内容更新时间，每次更新时自动修改
	DeleteTime field.Int    // 逻辑删除时间，NULL表示未删除

	fieldMap map[string]field.Expr
}

func (p postContent) Table(newTableName string) *postContent {
	p.postContentDo.UseTable(newTableName)
	return p.updateTableName(newTableName)
}

func (p postContent) As(alias string) *postContent {
	p.postContentDo.DO = *(p.postContentDo.As(alias).(*gen.DO))
	return p.updateTableName(alias)
}

func (p *postContent) updateTableName(table string) *postContent {
	p.ALL = field.NewAsterisk(table)
	p.PostID = field.NewInt64(table, "post_id")
	p.Content = field.NewString(table, "content")
	p.CreateTime = field.NewTime(table, "create_time")
	p.UpdateTime = field.NewTime(table, "update_time")
	p.DeleteTime = field.NewInt(table, "delete_time")

	p.fillFieldMap()

	return p
}

func (p *postContent) WithContext(ctx context.Context) IPostContentDo {
	return p.postContentDo.WithContext(ctx)
}

func (p postContent) TableName() string { return p.postContentDo.TableName() }

func (p postContent) Alias() string { return p.postContentDo.Alias() }

func (p postContent) Columns(cols ...field.Expr) gen.Columns { return p.postContentDo.Columns(cols...) }

func (p *postContent) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := p.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (p *postContent) fillFieldMap() {
	p.fieldMap = make(map[string]field.Expr, 5)
	p.fieldMap["post_id"] = p.PostID
	p.fieldMap["content"] = p.Content
	p.fieldMap["create_time"] = p.CreateTime
	p.fieldMap["update_time"] = p.UpdateTime
	p.fieldMap["delete_time"] = p.DeleteTime
}

func (p postContent) clone(db *gorm.DB) postContent {
	p.postContentDo.ReplaceConnPool(db.Statement.ConnPool)
	return p
}

func (p postContent) replaceDB(db *gorm.DB) postContent {
	p.postContentDo.ReplaceDB(db)
	return p
}

type postContentDo struct{ gen.DO }

type IPostContentDo interface {
	gen.SubQuery
	Debug() IPostContentDo
	WithContext(ctx context.Context) IPostContentDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IPostContentDo
	WriteDB() IPostContentDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IPostContentDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IPostContentDo
	Not(conds ...gen.Condition) IPostContentDo
	Or(conds ...gen.Condition) IPostContentDo
	Select(conds ...field.Expr) IPostContentDo
	Where(conds ...gen.Condition) IPostContentDo
	Order(conds ...field.Expr) IPostContentDo
	Distinct(cols ...field.Expr) IPostContentDo
	Omit(cols ...field.Expr) IPostContentDo
	Join(table schema.Tabler, on ...field.Expr) IPostContentDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IPostContentDo
	RightJoin(table schema.Tabler, on ...field.Expr) IPostContentDo
	Group(cols ...field.Expr) IPostContentDo
	Having(conds ...gen.Condition) IPostContentDo
	Limit(limit int) IPostContentDo
	Offset(offset int) IPostContentDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IPostContentDo
	Unscoped() IPostContentDo
	Create(values ...*model.PostContent) error
	CreateInBatches(values []*model.PostContent, batchSize int) error
	Save(values ...*model.PostContent) error
	First() (*model.PostContent, error)
	Take() (*model.PostContent, error)
	Last() (*model.PostContent, error)
	Find() ([]*model.PostContent, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.PostContent, err error)
	FindInBatches(result *[]*model.PostContent, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.PostContent) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IPostContentDo
	Assign(attrs ...field.AssignExpr) IPostContentDo
	Joins(fields ...field.RelationField) IPostContentDo
	Preload(fields ...field.RelationField) IPostContentDo
	FirstOrInit() (*model.PostContent, error)
	FirstOrCreate() (*model.PostContent, error)
	FindByPage(offset int, limit int) (result []*model.PostContent, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IPostContentDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (p postContentDo) Debug() IPostContentDo {
	return p.withDO(p.DO.Debug())
}

func (p postContentDo) WithContext(ctx context.Context) IPostContentDo {
	return p.withDO(p.DO.WithContext(ctx))
}

func (p postContentDo) ReadDB() IPostContentDo {
	return p.Clauses(dbresolver.Read)
}

func (p postContentDo) WriteDB() IPostContentDo {
	return p.Clauses(dbresolver.Write)
}

func (p postContentDo) Session(config *gorm.Session) IPostContentDo {
	return p.withDO(p.DO.Session(config))
}

func (p postContentDo) Clauses(conds ...clause.Expression) IPostContentDo {
	return p.withDO(p.DO.Clauses(conds...))
}

func (p postContentDo) Returning(value interface{}, columns ...string) IPostContentDo {
	return p.withDO(p.DO.Returning(value, columns...))
}

func (p postContentDo) Not(conds ...gen.Condition) IPostContentDo {
	return p.withDO(p.DO.Not(conds...))
}

func (p postContentDo) Or(conds ...gen.Condition) IPostContentDo {
	return p.withDO(p.DO.Or(conds...))
}

func (p postContentDo) Select(conds ...field.Expr) IPostContentDo {
	return p.withDO(p.DO.Select(conds...))
}

func (p postContentDo) Where(conds ...gen.Condition) IPostContentDo {
	return p.withDO(p.DO.Where(conds...))
}

func (p postContentDo) Order(conds ...field.Expr) IPostContentDo {
	return p.withDO(p.DO.Order(conds...))
}

func (p postContentDo) Distinct(cols ...field.Expr) IPostContentDo {
	return p.withDO(p.DO.Distinct(cols...))
}

func (p postContentDo) Omit(cols ...field.Expr) IPostContentDo {
	return p.withDO(p.DO.Omit(cols...))
}

func (p postContentDo) Join(table schema.Tabler, on ...field.Expr) IPostContentDo {
	return p.withDO(p.DO.Join(table, on...))
}

func (p postContentDo) LeftJoin(table schema.Tabler, on ...field.Expr) IPostContentDo {
	return p.withDO(p.DO.LeftJoin(table, on...))
}

func (p postContentDo) RightJoin(table schema.Tabler, on ...field.Expr) IPostContentDo {
	return p.withDO(p.DO.RightJoin(table, on...))
}

func (p postContentDo) Group(cols ...field.Expr) IPostContentDo {
	return p.withDO(p.DO.Group(cols...))
}

func (p postContentDo) Having(conds ...gen.Condition) IPostContentDo {
	return p.withDO(p.DO.Having(conds...))
}

func (p postContentDo) Limit(limit int) IPostContentDo {
	return p.withDO(p.DO.Limit(limit))
}

func (p postContentDo) Offset(offset int) IPostContentDo {
	return p.withDO(p.DO.Offset(offset))
}

func (p postContentDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IPostContentDo {
	return p.withDO(p.DO.Scopes(funcs...))
}

func (p postContentDo) Unscoped() IPostContentDo {
	return p.withDO(p.DO.Unscoped())
}

func (p postContentDo) Create(values ...*model.PostContent) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Create(values)
}

func (p postContentDo) CreateInBatches(values []*model.PostContent, batchSize int) error {
	return p.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (p postContentDo) Save(values ...*model.PostContent) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Save(values)
}

func (p postContentDo) First() (*model.PostContent, error) {
	if result, err := p.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.PostContent), nil
	}
}

func (p postContentDo) Take() (*model.PostContent, error) {
	if result, err := p.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.PostContent), nil
	}
}

func (p postContentDo) Last() (*model.PostContent, error) {
	if result, err := p.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.PostContent), nil
	}
}

func (p postContentDo) Find() ([]*model.PostContent, error) {
	result, err := p.DO.Find()
	return result.([]*model.PostContent), err
}

func (p postContentDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.PostContent, err error) {
	buf := make([]*model.PostContent, 0, batchSize)
	err = p.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (p postContentDo) FindInBatches(result *[]*model.PostContent, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return p.DO.FindInBatches(result, batchSize, fc)
}

func (p postContentDo) Attrs(attrs ...field.AssignExpr) IPostContentDo {
	return p.withDO(p.DO.Attrs(attrs...))
}

func (p postContentDo) Assign(attrs ...field.AssignExpr) IPostContentDo {
	return p.withDO(p.DO.Assign(attrs...))
}

func (p postContentDo) Joins(fields ...field.RelationField) IPostContentDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Joins(_f))
	}
	return &p
}

func (p postContentDo) Preload(fields ...field.RelationField) IPostContentDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Preload(_f))
	}
	return &p
}

func (p postContentDo) FirstOrInit() (*model.PostContent, error) {
	if result, err := p.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.PostContent), nil
	}
}

func (p postContentDo) FirstOrCreate() (*model.PostContent, error) {
	if result, err := p.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.PostContent), nil
	}
}

func (p postContentDo) FindByPage(offset int, limit int) (result []*model.PostContent, count int64, err error) {
	result, err = p.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = p.Offset(-1).Limit(-1).Count()
	return
}

func (p postContentDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = p.Count()
	if err != nil {
		return
	}

	err = p.Offset(offset).Limit(limit).Scan(result)
	return
}

func (p postContentDo) Scan(result interface{}) (err error) {
	return p.DO.Scan(result)
}

func (p postContentDo) Delete(models ...*model.PostContent) (result gen.ResultInfo, err error) {
	return p.DO.Delete(models)
}

func (p *postContentDo) withDO(do gen.Dao) *postContentDo {
	p.DO = *do.(*gen.DO)
	return p
}