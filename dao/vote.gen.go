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

func newVote(db *gorm.DB, opts ...gen.DOOption) vote {
	_vote := vote{}

	_vote.voteDo.UseDB(db, opts...)
	_vote.voteDo.UseModel(&model.Vote{})

	tableName := _vote.voteDo.TableName()
	_vote.ALL = field.NewAsterisk(tableName)
	_vote.ID = field.NewInt64(tableName, "id")
	_vote.PostID = field.NewInt64(tableName, "post_id")
	_vote.CommentID = field.NewInt64(tableName, "comment_id")
	_vote.UserID = field.NewInt64(tableName, "user_id")
	_vote.Vote = field.NewInt32(tableName, "vote")
	_vote.CreateTime = field.NewTime(tableName, "create_time")
	_vote.UpdateTime = field.NewTime(tableName, "update_time")
	_vote.DeleteTime = field.NewInt(tableName, "delete_time")

	_vote.fillFieldMap()

	return _vote
}

type vote struct {
	voteDo voteDo

	ALL        field.Asterisk
	ID         field.Int64 // 自增主键，唯一标识每条投票记录
	PostID     field.Int64 // 投票所属的帖子ID
	CommentID  field.Int64 // 投票所属的评论ID
	UserID     field.Int64 // 投票用户的用户ID
	Vote       field.Int32 // 投票类型：1-赞，-1-踩
	CreateTime field.Time  // 投票创建时间，默认当前时间
	UpdateTime field.Time  // 投票更新时间，每次更新时自动修改
	DeleteTime field.Int   // 逻辑删除时间，NULL表示未删除

	fieldMap map[string]field.Expr
}

func (v vote) Table(newTableName string) *vote {
	v.voteDo.UseTable(newTableName)
	return v.updateTableName(newTableName)
}

func (v vote) As(alias string) *vote {
	v.voteDo.DO = *(v.voteDo.As(alias).(*gen.DO))
	return v.updateTableName(alias)
}

func (v *vote) updateTableName(table string) *vote {
	v.ALL = field.NewAsterisk(table)
	v.ID = field.NewInt64(table, "id")
	v.PostID = field.NewInt64(table, "post_id")
	v.CommentID = field.NewInt64(table, "comment_id")
	v.UserID = field.NewInt64(table, "user_id")
	v.Vote = field.NewInt32(table, "vote")
	v.CreateTime = field.NewTime(table, "create_time")
	v.UpdateTime = field.NewTime(table, "update_time")
	v.DeleteTime = field.NewInt(table, "delete_time")

	v.fillFieldMap()

	return v
}

func (v *vote) WithContext(ctx context.Context) IVoteDo { return v.voteDo.WithContext(ctx) }

func (v vote) TableName() string { return v.voteDo.TableName() }

func (v vote) Alias() string { return v.voteDo.Alias() }

func (v vote) Columns(cols ...field.Expr) gen.Columns { return v.voteDo.Columns(cols...) }

func (v *vote) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := v.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (v *vote) fillFieldMap() {
	v.fieldMap = make(map[string]field.Expr, 8)
	v.fieldMap["id"] = v.ID
	v.fieldMap["post_id"] = v.PostID
	v.fieldMap["comment_id"] = v.CommentID
	v.fieldMap["user_id"] = v.UserID
	v.fieldMap["vote"] = v.Vote
	v.fieldMap["create_time"] = v.CreateTime
	v.fieldMap["update_time"] = v.UpdateTime
	v.fieldMap["delete_time"] = v.DeleteTime
}

func (v vote) clone(db *gorm.DB) vote {
	v.voteDo.ReplaceConnPool(db.Statement.ConnPool)
	return v
}

func (v vote) replaceDB(db *gorm.DB) vote {
	v.voteDo.ReplaceDB(db)
	return v
}

type voteDo struct{ gen.DO }

type IVoteDo interface {
	gen.SubQuery
	Debug() IVoteDo
	WithContext(ctx context.Context) IVoteDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IVoteDo
	WriteDB() IVoteDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IVoteDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IVoteDo
	Not(conds ...gen.Condition) IVoteDo
	Or(conds ...gen.Condition) IVoteDo
	Select(conds ...field.Expr) IVoteDo
	Where(conds ...gen.Condition) IVoteDo
	Order(conds ...field.Expr) IVoteDo
	Distinct(cols ...field.Expr) IVoteDo
	Omit(cols ...field.Expr) IVoteDo
	Join(table schema.Tabler, on ...field.Expr) IVoteDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IVoteDo
	RightJoin(table schema.Tabler, on ...field.Expr) IVoteDo
	Group(cols ...field.Expr) IVoteDo
	Having(conds ...gen.Condition) IVoteDo
	Limit(limit int) IVoteDo
	Offset(offset int) IVoteDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IVoteDo
	Unscoped() IVoteDo
	Create(values ...*model.Vote) error
	CreateInBatches(values []*model.Vote, batchSize int) error
	Save(values ...*model.Vote) error
	First() (*model.Vote, error)
	Take() (*model.Vote, error)
	Last() (*model.Vote, error)
	Find() ([]*model.Vote, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Vote, err error)
	FindInBatches(result *[]*model.Vote, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Vote) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IVoteDo
	Assign(attrs ...field.AssignExpr) IVoteDo
	Joins(fields ...field.RelationField) IVoteDo
	Preload(fields ...field.RelationField) IVoteDo
	FirstOrInit() (*model.Vote, error)
	FirstOrCreate() (*model.Vote, error)
	FindByPage(offset int, limit int) (result []*model.Vote, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IVoteDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (v voteDo) Debug() IVoteDo {
	return v.withDO(v.DO.Debug())
}

func (v voteDo) WithContext(ctx context.Context) IVoteDo {
	return v.withDO(v.DO.WithContext(ctx))
}

func (v voteDo) ReadDB() IVoteDo {
	return v.Clauses(dbresolver.Read)
}

func (v voteDo) WriteDB() IVoteDo {
	return v.Clauses(dbresolver.Write)
}

func (v voteDo) Session(config *gorm.Session) IVoteDo {
	return v.withDO(v.DO.Session(config))
}

func (v voteDo) Clauses(conds ...clause.Expression) IVoteDo {
	return v.withDO(v.DO.Clauses(conds...))
}

func (v voteDo) Returning(value interface{}, columns ...string) IVoteDo {
	return v.withDO(v.DO.Returning(value, columns...))
}

func (v voteDo) Not(conds ...gen.Condition) IVoteDo {
	return v.withDO(v.DO.Not(conds...))
}

func (v voteDo) Or(conds ...gen.Condition) IVoteDo {
	return v.withDO(v.DO.Or(conds...))
}

func (v voteDo) Select(conds ...field.Expr) IVoteDo {
	return v.withDO(v.DO.Select(conds...))
}

func (v voteDo) Where(conds ...gen.Condition) IVoteDo {
	return v.withDO(v.DO.Where(conds...))
}

func (v voteDo) Order(conds ...field.Expr) IVoteDo {
	return v.withDO(v.DO.Order(conds...))
}

func (v voteDo) Distinct(cols ...field.Expr) IVoteDo {
	return v.withDO(v.DO.Distinct(cols...))
}

func (v voteDo) Omit(cols ...field.Expr) IVoteDo {
	return v.withDO(v.DO.Omit(cols...))
}

func (v voteDo) Join(table schema.Tabler, on ...field.Expr) IVoteDo {
	return v.withDO(v.DO.Join(table, on...))
}

func (v voteDo) LeftJoin(table schema.Tabler, on ...field.Expr) IVoteDo {
	return v.withDO(v.DO.LeftJoin(table, on...))
}

func (v voteDo) RightJoin(table schema.Tabler, on ...field.Expr) IVoteDo {
	return v.withDO(v.DO.RightJoin(table, on...))
}

func (v voteDo) Group(cols ...field.Expr) IVoteDo {
	return v.withDO(v.DO.Group(cols...))
}

func (v voteDo) Having(conds ...gen.Condition) IVoteDo {
	return v.withDO(v.DO.Having(conds...))
}

func (v voteDo) Limit(limit int) IVoteDo {
	return v.withDO(v.DO.Limit(limit))
}

func (v voteDo) Offset(offset int) IVoteDo {
	return v.withDO(v.DO.Offset(offset))
}

func (v voteDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IVoteDo {
	return v.withDO(v.DO.Scopes(funcs...))
}

func (v voteDo) Unscoped() IVoteDo {
	return v.withDO(v.DO.Unscoped())
}

func (v voteDo) Create(values ...*model.Vote) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Create(values)
}

func (v voteDo) CreateInBatches(values []*model.Vote, batchSize int) error {
	return v.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (v voteDo) Save(values ...*model.Vote) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Save(values)
}

func (v voteDo) First() (*model.Vote, error) {
	if result, err := v.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Vote), nil
	}
}

func (v voteDo) Take() (*model.Vote, error) {
	if result, err := v.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Vote), nil
	}
}

func (v voteDo) Last() (*model.Vote, error) {
	if result, err := v.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Vote), nil
	}
}

func (v voteDo) Find() ([]*model.Vote, error) {
	result, err := v.DO.Find()
	return result.([]*model.Vote), err
}

func (v voteDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Vote, err error) {
	buf := make([]*model.Vote, 0, batchSize)
	err = v.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (v voteDo) FindInBatches(result *[]*model.Vote, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return v.DO.FindInBatches(result, batchSize, fc)
}

func (v voteDo) Attrs(attrs ...field.AssignExpr) IVoteDo {
	return v.withDO(v.DO.Attrs(attrs...))
}

func (v voteDo) Assign(attrs ...field.AssignExpr) IVoteDo {
	return v.withDO(v.DO.Assign(attrs...))
}

func (v voteDo) Joins(fields ...field.RelationField) IVoteDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Joins(_f))
	}
	return &v
}

func (v voteDo) Preload(fields ...field.RelationField) IVoteDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Preload(_f))
	}
	return &v
}

func (v voteDo) FirstOrInit() (*model.Vote, error) {
	if result, err := v.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Vote), nil
	}
}

func (v voteDo) FirstOrCreate() (*model.Vote, error) {
	if result, err := v.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Vote), nil
	}
}

func (v voteDo) FindByPage(offset int, limit int) (result []*model.Vote, count int64, err error) {
	result, err = v.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = v.Offset(-1).Limit(-1).Count()
	return
}

func (v voteDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = v.Count()
	if err != nil {
		return
	}

	err = v.Offset(offset).Limit(limit).Scan(result)
	return
}

func (v voteDo) Scan(result interface{}) (err error) {
	return v.DO.Scan(result)
}

func (v voteDo) Delete(models ...*model.Vote) (result gen.ResultInfo, err error) {
	return v.DO.Delete(models)
}

func (v *voteDo) withDO(do gen.Dao) *voteDo {
	v.DO = *do.(*gen.DO)
	return v
}
