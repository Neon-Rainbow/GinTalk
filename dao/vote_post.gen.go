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

func newVotePost(db *gorm.DB, opts ...gen.DOOption) votePost {
	_votePost := votePost{}

	_votePost.votePostDo.UseDB(db, opts...)
	_votePost.votePostDo.UseModel(&model.VotePost{})

	tableName := _votePost.votePostDo.TableName()
	_votePost.ALL = field.NewAsterisk(tableName)
	_votePost.ID = field.NewInt64(tableName, "id")
	_votePost.PostID = field.NewInt64(tableName, "post_id")
	_votePost.UserID = field.NewInt64(tableName, "user_id")
	_votePost.Vote = field.NewInt32(tableName, "vote")
	_votePost.CreateTime = field.NewTime(tableName, "create_time")
	_votePost.UpdateTime = field.NewTime(tableName, "update_time")
	_votePost.DeleteTime = field.NewInt(tableName, "delete_time")

	_votePost.fillFieldMap()

	return _votePost
}

// votePost 帖子投票表：存储用户对帖子的投票记录
type votePost struct {
	votePostDo votePostDo

	ALL        field.Asterisk
	ID         field.Int64 // 自增主键，唯一标识每条投票记录
	PostID     field.Int64 // 投票所属的帖子ID
	UserID     field.Int64 // 投票用户的用户ID
	Vote       field.Int32 // 投票类型：1-赞
	CreateTime field.Time  // 投票创建时间，默认当前时间
	UpdateTime field.Time  // 投票更新时间，每次更新时自动修改
	DeleteTime field.Int   // 逻辑删除时间，NULL表示未删除

	fieldMap map[string]field.Expr
}

func (v votePost) Table(newTableName string) *votePost {
	v.votePostDo.UseTable(newTableName)
	return v.updateTableName(newTableName)
}

func (v votePost) As(alias string) *votePost {
	v.votePostDo.DO = *(v.votePostDo.As(alias).(*gen.DO))
	return v.updateTableName(alias)
}

func (v *votePost) updateTableName(table string) *votePost {
	v.ALL = field.NewAsterisk(table)
	v.ID = field.NewInt64(table, "id")
	v.PostID = field.NewInt64(table, "post_id")
	v.UserID = field.NewInt64(table, "user_id")
	v.Vote = field.NewInt32(table, "vote")
	v.CreateTime = field.NewTime(table, "create_time")
	v.UpdateTime = field.NewTime(table, "update_time")
	v.DeleteTime = field.NewInt(table, "delete_time")

	v.fillFieldMap()

	return v
}

func (v *votePost) WithContext(ctx context.Context) IVotePostDo { return v.votePostDo.WithContext(ctx) }

func (v votePost) TableName() string { return v.votePostDo.TableName() }

func (v votePost) Alias() string { return v.votePostDo.Alias() }

func (v votePost) Columns(cols ...field.Expr) gen.Columns { return v.votePostDo.Columns(cols...) }

func (v *votePost) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := v.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (v *votePost) fillFieldMap() {
	v.fieldMap = make(map[string]field.Expr, 7)
	v.fieldMap["id"] = v.ID
	v.fieldMap["post_id"] = v.PostID
	v.fieldMap["user_id"] = v.UserID
	v.fieldMap["vote"] = v.Vote
	v.fieldMap["create_time"] = v.CreateTime
	v.fieldMap["update_time"] = v.UpdateTime
	v.fieldMap["delete_time"] = v.DeleteTime
}

func (v votePost) clone(db *gorm.DB) votePost {
	v.votePostDo.ReplaceConnPool(db.Statement.ConnPool)
	return v
}

func (v votePost) replaceDB(db *gorm.DB) votePost {
	v.votePostDo.ReplaceDB(db)
	return v
}

type votePostDo struct{ gen.DO }

type IVotePostDo interface {
	gen.SubQuery
	Debug() IVotePostDo
	WithContext(ctx context.Context) IVotePostDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IVotePostDo
	WriteDB() IVotePostDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IVotePostDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IVotePostDo
	Not(conds ...gen.Condition) IVotePostDo
	Or(conds ...gen.Condition) IVotePostDo
	Select(conds ...field.Expr) IVotePostDo
	Where(conds ...gen.Condition) IVotePostDo
	Order(conds ...field.Expr) IVotePostDo
	Distinct(cols ...field.Expr) IVotePostDo
	Omit(cols ...field.Expr) IVotePostDo
	Join(table schema.Tabler, on ...field.Expr) IVotePostDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IVotePostDo
	RightJoin(table schema.Tabler, on ...field.Expr) IVotePostDo
	Group(cols ...field.Expr) IVotePostDo
	Having(conds ...gen.Condition) IVotePostDo
	Limit(limit int) IVotePostDo
	Offset(offset int) IVotePostDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IVotePostDo
	Unscoped() IVotePostDo
	Create(values ...*model.VotePost) error
	CreateInBatches(values []*model.VotePost, batchSize int) error
	Save(values ...*model.VotePost) error
	First() (*model.VotePost, error)
	Take() (*model.VotePost, error)
	Last() (*model.VotePost, error)
	Find() ([]*model.VotePost, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.VotePost, err error)
	FindInBatches(result *[]*model.VotePost, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.VotePost) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IVotePostDo
	Assign(attrs ...field.AssignExpr) IVotePostDo
	Joins(fields ...field.RelationField) IVotePostDo
	Preload(fields ...field.RelationField) IVotePostDo
	FirstOrInit() (*model.VotePost, error)
	FirstOrCreate() (*model.VotePost, error)
	FindByPage(offset int, limit int) (result []*model.VotePost, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IVotePostDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (v votePostDo) Debug() IVotePostDo {
	return v.withDO(v.DO.Debug())
}

func (v votePostDo) WithContext(ctx context.Context) IVotePostDo {
	return v.withDO(v.DO.WithContext(ctx))
}

func (v votePostDo) ReadDB() IVotePostDo {
	return v.Clauses(dbresolver.Read)
}

func (v votePostDo) WriteDB() IVotePostDo {
	return v.Clauses(dbresolver.Write)
}

func (v votePostDo) Session(config *gorm.Session) IVotePostDo {
	return v.withDO(v.DO.Session(config))
}

func (v votePostDo) Clauses(conds ...clause.Expression) IVotePostDo {
	return v.withDO(v.DO.Clauses(conds...))
}

func (v votePostDo) Returning(value interface{}, columns ...string) IVotePostDo {
	return v.withDO(v.DO.Returning(value, columns...))
}

func (v votePostDo) Not(conds ...gen.Condition) IVotePostDo {
	return v.withDO(v.DO.Not(conds...))
}

func (v votePostDo) Or(conds ...gen.Condition) IVotePostDo {
	return v.withDO(v.DO.Or(conds...))
}

func (v votePostDo) Select(conds ...field.Expr) IVotePostDo {
	return v.withDO(v.DO.Select(conds...))
}

func (v votePostDo) Where(conds ...gen.Condition) IVotePostDo {
	return v.withDO(v.DO.Where(conds...))
}

func (v votePostDo) Order(conds ...field.Expr) IVotePostDo {
	return v.withDO(v.DO.Order(conds...))
}

func (v votePostDo) Distinct(cols ...field.Expr) IVotePostDo {
	return v.withDO(v.DO.Distinct(cols...))
}

func (v votePostDo) Omit(cols ...field.Expr) IVotePostDo {
	return v.withDO(v.DO.Omit(cols...))
}

func (v votePostDo) Join(table schema.Tabler, on ...field.Expr) IVotePostDo {
	return v.withDO(v.DO.Join(table, on...))
}

func (v votePostDo) LeftJoin(table schema.Tabler, on ...field.Expr) IVotePostDo {
	return v.withDO(v.DO.LeftJoin(table, on...))
}

func (v votePostDo) RightJoin(table schema.Tabler, on ...field.Expr) IVotePostDo {
	return v.withDO(v.DO.RightJoin(table, on...))
}

func (v votePostDo) Group(cols ...field.Expr) IVotePostDo {
	return v.withDO(v.DO.Group(cols...))
}

func (v votePostDo) Having(conds ...gen.Condition) IVotePostDo {
	return v.withDO(v.DO.Having(conds...))
}

func (v votePostDo) Limit(limit int) IVotePostDo {
	return v.withDO(v.DO.Limit(limit))
}

func (v votePostDo) Offset(offset int) IVotePostDo {
	return v.withDO(v.DO.Offset(offset))
}

func (v votePostDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IVotePostDo {
	return v.withDO(v.DO.Scopes(funcs...))
}

func (v votePostDo) Unscoped() IVotePostDo {
	return v.withDO(v.DO.Unscoped())
}

func (v votePostDo) Create(values ...*model.VotePost) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Create(values)
}

func (v votePostDo) CreateInBatches(values []*model.VotePost, batchSize int) error {
	return v.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (v votePostDo) Save(values ...*model.VotePost) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Save(values)
}

func (v votePostDo) First() (*model.VotePost, error) {
	if result, err := v.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.VotePost), nil
	}
}

func (v votePostDo) Take() (*model.VotePost, error) {
	if result, err := v.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.VotePost), nil
	}
}

func (v votePostDo) Last() (*model.VotePost, error) {
	if result, err := v.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.VotePost), nil
	}
}

func (v votePostDo) Find() ([]*model.VotePost, error) {
	result, err := v.DO.Find()
	return result.([]*model.VotePost), err
}

func (v votePostDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.VotePost, err error) {
	buf := make([]*model.VotePost, 0, batchSize)
	err = v.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (v votePostDo) FindInBatches(result *[]*model.VotePost, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return v.DO.FindInBatches(result, batchSize, fc)
}

func (v votePostDo) Attrs(attrs ...field.AssignExpr) IVotePostDo {
	return v.withDO(v.DO.Attrs(attrs...))
}

func (v votePostDo) Assign(attrs ...field.AssignExpr) IVotePostDo {
	return v.withDO(v.DO.Assign(attrs...))
}

func (v votePostDo) Joins(fields ...field.RelationField) IVotePostDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Joins(_f))
	}
	return &v
}

func (v votePostDo) Preload(fields ...field.RelationField) IVotePostDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Preload(_f))
	}
	return &v
}

func (v votePostDo) FirstOrInit() (*model.VotePost, error) {
	if result, err := v.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.VotePost), nil
	}
}

func (v votePostDo) FirstOrCreate() (*model.VotePost, error) {
	if result, err := v.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.VotePost), nil
	}
}

func (v votePostDo) FindByPage(offset int, limit int) (result []*model.VotePost, count int64, err error) {
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

func (v votePostDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = v.Count()
	if err != nil {
		return
	}

	err = v.Offset(offset).Limit(limit).Scan(result)
	return
}

func (v votePostDo) Scan(result interface{}) (err error) {
	return v.DO.Scan(result)
}

func (v votePostDo) Delete(models ...*model.VotePost) (result gen.ResultInfo, err error) {
	return v.DO.Delete(models)
}

func (v *votePostDo) withDO(do gen.Dao) *votePostDo {
	v.DO = *do.(*gen.DO)
	return v
}