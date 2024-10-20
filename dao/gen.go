// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q               = new(Query)
	Comment         *comment
	CommentRelation *commentRelation
	CommentVote     *commentVote
	Community       *community
	ContentVote     *contentVote
	Post            *post
	User            *user
	Vote            *vote
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Comment = &Q.Comment
	CommentRelation = &Q.CommentRelation
	CommentVote = &Q.CommentVote
	Community = &Q.Community
	ContentVote = &Q.ContentVote
	Post = &Q.Post
	User = &Q.User
	Vote = &Q.Vote
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:              db,
		Comment:         newComment(db, opts...),
		CommentRelation: newCommentRelation(db, opts...),
		CommentVote:     newCommentVote(db, opts...),
		Community:       newCommunity(db, opts...),
		ContentVote:     newContentVote(db, opts...),
		Post:            newPost(db, opts...),
		User:            newUser(db, opts...),
		Vote:            newVote(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Comment         comment
	CommentRelation commentRelation
	CommentVote     commentVote
	Community       community
	ContentVote     contentVote
	Post            post
	User            user
	Vote            vote
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:              db,
		Comment:         q.Comment.clone(db),
		CommentRelation: q.CommentRelation.clone(db),
		CommentVote:     q.CommentVote.clone(db),
		Community:       q.Community.clone(db),
		ContentVote:     q.ContentVote.clone(db),
		Post:            q.Post.clone(db),
		User:            q.User.clone(db),
		Vote:            q.Vote.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:              db,
		Comment:         q.Comment.replaceDB(db),
		CommentRelation: q.CommentRelation.replaceDB(db),
		CommentVote:     q.CommentVote.replaceDB(db),
		Community:       q.Community.replaceDB(db),
		ContentVote:     q.ContentVote.replaceDB(db),
		Post:            q.Post.replaceDB(db),
		User:            q.User.replaceDB(db),
		Vote:            q.Vote.replaceDB(db),
	}
}

type queryCtx struct {
	Comment         ICommentDo
	CommentRelation ICommentRelationDo
	CommentVote     ICommentVoteDo
	Community       ICommunityDo
	ContentVote     IContentVoteDo
	Post            IPostDo
	User            IUserDo
	Vote            IVoteDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Comment:         q.Comment.WithContext(ctx),
		CommentRelation: q.CommentRelation.WithContext(ctx),
		CommentVote:     q.CommentVote.WithContext(ctx),
		Community:       q.Community.WithContext(ctx),
		ContentVote:     q.ContentVote.WithContext(ctx),
		Post:            q.Post.WithContext(ctx),
		User:            q.User.WithContext(ctx),
		Vote:            q.Vote.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
