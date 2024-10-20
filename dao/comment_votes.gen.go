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

func newCommentVote(db *gorm.DB, opts ...gen.DOOption) commentVote {
	_commentVote := commentVote{}

	_commentVote.commentVoteDo.UseDB(db, opts...)
	_commentVote.commentVoteDo.UseModel(&model.CommentVote{})

	tableName := _commentVote.commentVoteDo.TableName()
	_commentVote.ALL = field.NewAsterisk(tableName)
	_commentVote.CommentID = field.NewInt64(tableName, "comment_id")
	_commentVote.Up = field.NewInt32(tableName, "up")
	_commentVote.Down = field.NewInt32(tableName, "down")
	_commentVote.CreateTime = field.NewTime(tableName, "create_time")
	_commentVote.UpdateTime = field.NewTime(tableName, "update_time")
	_commentVote.DeleteTime = field.NewField(tableName, "delete_time")

	_commentVote.fillFieldMap()

	return _commentVote
}

type commentVote struct {
	commentVoteDo commentVoteDo

	ALL        field.Asterisk
	CommentID  field.Int64 // 投票所属的评论ID
	Up         field.Int32 // 赞数
	Down       field.Int32 // 踩数
	CreateTime field.Time  // 投票创建时间，默认当前时间
	UpdateTime field.Time  // 投票更新时间，每次更新时自动修改
	DeleteTime field.Field // 逻辑删除时间，NULL表示未删除

	fieldMap map[string]field.Expr
}

func (c commentVote) Table(newTableName string) *commentVote {
	c.commentVoteDo.UseTable(newTableName)
	return c.updateTableName(newTableName)
}

func (c commentVote) As(alias string) *commentVote {
	c.commentVoteDo.DO = *(c.commentVoteDo.As(alias).(*gen.DO))
	return c.updateTableName(alias)
}

func (c *commentVote) updateTableName(table string) *commentVote {
	c.ALL = field.NewAsterisk(table)
	c.CommentID = field.NewInt64(table, "comment_id")
	c.Up = field.NewInt32(table, "up")
	c.Down = field.NewInt32(table, "down")
	c.CreateTime = field.NewTime(table, "create_time")
	c.UpdateTime = field.NewTime(table, "update_time")
	c.DeleteTime = field.NewField(table, "delete_time")

	c.fillFieldMap()

	return c
}

func (c *commentVote) WithContext(ctx context.Context) ICommentVoteDo {
	return c.commentVoteDo.WithContext(ctx)
}

func (c commentVote) TableName() string { return c.commentVoteDo.TableName() }

func (c commentVote) Alias() string { return c.commentVoteDo.Alias() }

func (c commentVote) Columns(cols ...field.Expr) gen.Columns { return c.commentVoteDo.Columns(cols...) }

func (c *commentVote) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := c.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (c *commentVote) fillFieldMap() {
	c.fieldMap = make(map[string]field.Expr, 6)
	c.fieldMap["comment_id"] = c.CommentID
	c.fieldMap["up"] = c.Up
	c.fieldMap["down"] = c.Down
	c.fieldMap["create_time"] = c.CreateTime
	c.fieldMap["update_time"] = c.UpdateTime
	c.fieldMap["delete_time"] = c.DeleteTime
}

func (c commentVote) clone(db *gorm.DB) commentVote {
	c.commentVoteDo.ReplaceConnPool(db.Statement.ConnPool)
	return c
}

func (c commentVote) replaceDB(db *gorm.DB) commentVote {
	c.commentVoteDo.ReplaceDB(db)
	return c
}

type commentVoteDo struct{ gen.DO }

type ICommentVoteDo interface {
	gen.SubQuery
	Debug() ICommentVoteDo
	WithContext(ctx context.Context) ICommentVoteDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ICommentVoteDo
	WriteDB() ICommentVoteDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ICommentVoteDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ICommentVoteDo
	Not(conds ...gen.Condition) ICommentVoteDo
	Or(conds ...gen.Condition) ICommentVoteDo
	Select(conds ...field.Expr) ICommentVoteDo
	Where(conds ...gen.Condition) ICommentVoteDo
	Order(conds ...field.Expr) ICommentVoteDo
	Distinct(cols ...field.Expr) ICommentVoteDo
	Omit(cols ...field.Expr) ICommentVoteDo
	Join(table schema.Tabler, on ...field.Expr) ICommentVoteDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ICommentVoteDo
	RightJoin(table schema.Tabler, on ...field.Expr) ICommentVoteDo
	Group(cols ...field.Expr) ICommentVoteDo
	Having(conds ...gen.Condition) ICommentVoteDo
	Limit(limit int) ICommentVoteDo
	Offset(offset int) ICommentVoteDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ICommentVoteDo
	Unscoped() ICommentVoteDo
	Create(values ...*model.CommentVote) error
	CreateInBatches(values []*model.CommentVote, batchSize int) error
	Save(values ...*model.CommentVote) error
	First() (*model.CommentVote, error)
	Take() (*model.CommentVote, error)
	Last() (*model.CommentVote, error)
	Find() ([]*model.CommentVote, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.CommentVote, err error)
	FindInBatches(result *[]*model.CommentVote, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.CommentVote) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ICommentVoteDo
	Assign(attrs ...field.AssignExpr) ICommentVoteDo
	Joins(fields ...field.RelationField) ICommentVoteDo
	Preload(fields ...field.RelationField) ICommentVoteDo
	FirstOrInit() (*model.CommentVote, error)
	FirstOrCreate() (*model.CommentVote, error)
	FindByPage(offset int, limit int) (result []*model.CommentVote, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ICommentVoteDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (c commentVoteDo) Debug() ICommentVoteDo {
	return c.withDO(c.DO.Debug())
}

func (c commentVoteDo) WithContext(ctx context.Context) ICommentVoteDo {
	return c.withDO(c.DO.WithContext(ctx))
}

func (c commentVoteDo) ReadDB() ICommentVoteDo {
	return c.Clauses(dbresolver.Read)
}

func (c commentVoteDo) WriteDB() ICommentVoteDo {
	return c.Clauses(dbresolver.Write)
}

func (c commentVoteDo) Session(config *gorm.Session) ICommentVoteDo {
	return c.withDO(c.DO.Session(config))
}

func (c commentVoteDo) Clauses(conds ...clause.Expression) ICommentVoteDo {
	return c.withDO(c.DO.Clauses(conds...))
}

func (c commentVoteDo) Returning(value interface{}, columns ...string) ICommentVoteDo {
	return c.withDO(c.DO.Returning(value, columns...))
}

func (c commentVoteDo) Not(conds ...gen.Condition) ICommentVoteDo {
	return c.withDO(c.DO.Not(conds...))
}

func (c commentVoteDo) Or(conds ...gen.Condition) ICommentVoteDo {
	return c.withDO(c.DO.Or(conds...))
}

func (c commentVoteDo) Select(conds ...field.Expr) ICommentVoteDo {
	return c.withDO(c.DO.Select(conds...))
}

func (c commentVoteDo) Where(conds ...gen.Condition) ICommentVoteDo {
	return c.withDO(c.DO.Where(conds...))
}

func (c commentVoteDo) Order(conds ...field.Expr) ICommentVoteDo {
	return c.withDO(c.DO.Order(conds...))
}

func (c commentVoteDo) Distinct(cols ...field.Expr) ICommentVoteDo {
	return c.withDO(c.DO.Distinct(cols...))
}

func (c commentVoteDo) Omit(cols ...field.Expr) ICommentVoteDo {
	return c.withDO(c.DO.Omit(cols...))
}

func (c commentVoteDo) Join(table schema.Tabler, on ...field.Expr) ICommentVoteDo {
	return c.withDO(c.DO.Join(table, on...))
}

func (c commentVoteDo) LeftJoin(table schema.Tabler, on ...field.Expr) ICommentVoteDo {
	return c.withDO(c.DO.LeftJoin(table, on...))
}

func (c commentVoteDo) RightJoin(table schema.Tabler, on ...field.Expr) ICommentVoteDo {
	return c.withDO(c.DO.RightJoin(table, on...))
}

func (c commentVoteDo) Group(cols ...field.Expr) ICommentVoteDo {
	return c.withDO(c.DO.Group(cols...))
}

func (c commentVoteDo) Having(conds ...gen.Condition) ICommentVoteDo {
	return c.withDO(c.DO.Having(conds...))
}

func (c commentVoteDo) Limit(limit int) ICommentVoteDo {
	return c.withDO(c.DO.Limit(limit))
}

func (c commentVoteDo) Offset(offset int) ICommentVoteDo {
	return c.withDO(c.DO.Offset(offset))
}

func (c commentVoteDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ICommentVoteDo {
	return c.withDO(c.DO.Scopes(funcs...))
}

func (c commentVoteDo) Unscoped() ICommentVoteDo {
	return c.withDO(c.DO.Unscoped())
}

func (c commentVoteDo) Create(values ...*model.CommentVote) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Create(values)
}

func (c commentVoteDo) CreateInBatches(values []*model.CommentVote, batchSize int) error {
	return c.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (c commentVoteDo) Save(values ...*model.CommentVote) error {
	if len(values) == 0 {
		return nil
	}
	return c.DO.Save(values)
}

func (c commentVoteDo) First() (*model.CommentVote, error) {
	if result, err := c.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.CommentVote), nil
	}
}

func (c commentVoteDo) Take() (*model.CommentVote, error) {
	if result, err := c.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.CommentVote), nil
	}
}

func (c commentVoteDo) Last() (*model.CommentVote, error) {
	if result, err := c.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.CommentVote), nil
	}
}

func (c commentVoteDo) Find() ([]*model.CommentVote, error) {
	result, err := c.DO.Find()
	return result.([]*model.CommentVote), err
}

func (c commentVoteDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.CommentVote, err error) {
	buf := make([]*model.CommentVote, 0, batchSize)
	err = c.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (c commentVoteDo) FindInBatches(result *[]*model.CommentVote, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return c.DO.FindInBatches(result, batchSize, fc)
}

func (c commentVoteDo) Attrs(attrs ...field.AssignExpr) ICommentVoteDo {
	return c.withDO(c.DO.Attrs(attrs...))
}

func (c commentVoteDo) Assign(attrs ...field.AssignExpr) ICommentVoteDo {
	return c.withDO(c.DO.Assign(attrs...))
}

func (c commentVoteDo) Joins(fields ...field.RelationField) ICommentVoteDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Joins(_f))
	}
	return &c
}

func (c commentVoteDo) Preload(fields ...field.RelationField) ICommentVoteDo {
	for _, _f := range fields {
		c = *c.withDO(c.DO.Preload(_f))
	}
	return &c
}

func (c commentVoteDo) FirstOrInit() (*model.CommentVote, error) {
	if result, err := c.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.CommentVote), nil
	}
}

func (c commentVoteDo) FirstOrCreate() (*model.CommentVote, error) {
	if result, err := c.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.CommentVote), nil
	}
}

func (c commentVoteDo) FindByPage(offset int, limit int) (result []*model.CommentVote, count int64, err error) {
	result, err = c.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = c.Offset(-1).Limit(-1).Count()
	return
}

func (c commentVoteDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = c.Count()
	if err != nil {
		return
	}

	err = c.Offset(offset).Limit(limit).Scan(result)
	return
}

func (c commentVoteDo) Scan(result interface{}) (err error) {
	return c.DO.Scan(result)
}

func (c commentVoteDo) Delete(models ...*model.CommentVote) (result gen.ResultInfo, err error) {
	return c.DO.Delete(models)
}

func (c *commentVoteDo) withDO(do gen.Dao) *commentVoteDo {
	c.DO = *do.(*gen.DO)
	return c
}