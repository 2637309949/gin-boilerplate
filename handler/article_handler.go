package handler

import (
	"gin-boilerplate/comm/db"
	"gin-boilerplate/comm/http"
	"gin-boilerplate/comm/util"
	"gin-boilerplate/models"
	"gin-boilerplate/types"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

//InsertArticle...
func (s *Handler) InsertArticle(ctx *gin.Context) {
	var session = db.GetDB()
	var articleForm types.ArticleForm
	if err := ctx.ShouldBindJSON(&articleForm); err != nil {
		http.Fail(ctx, http.MsgOption(articleForm.Insert(err)))
		return
	}

	inArticle := models.Article{}
	where := models.Article{
		Title: articleForm.Title,
	}
	copier.Copy(&inArticle, &articleForm)
	err := s.QueryArticleDetailDB(ctx, session, &where, &inArticle)
	if err == nil {
		http.Fail(ctx, http.MsgOption("记录已存在"))
		return
	}

	err = s.InsertArticleDB(ctx, session, &inArticle)
	if err != nil {
		http.Fail(ctx, http.MsgOption("InsertSchedulePositionDB failed. [%v]", err))
		return
	}

	http.Success(ctx, http.FlatOption(inArticle))
}

//QueryArticle...
func (s *Handler) QueryArticle(ctx *gin.Context) {
	var articleFilter types.ArticleFilter
	if err := ctx.Bind(&articleFilter); err != nil {
		http.Fail(ctx, http.MsgOption(articleFilter.Filter(err)))
		return
	}
	var session = db.GetDB()
	session = db.SetLimit(ctx, session, &articleFilter)
	session = db.SetOrder(ctx, session, &articleFilter)

	var totalCount int32
	var lst []models.Article
	where := models.Article{
		Title: articleFilter.Title,
	}
	err := s.QueryArticleDB(ctx, session, &where, &lst, &totalCount)
	if err != nil {
		http.Fail(ctx, http.MsgOption("QueryArticleDB failed. [%s]", err.Error()))
		return
	}

	var pager = http.Pager{}
	pager.Data = lst
	pager.TotalCount = totalCount
	pager.CurPage = articleFilter.GetPageNo()
	pager.TotalPage = totalCount / articleFilter.GetPageSize()
	if totalCount%articleFilter.GetPageSize() != 0 {
		pager.TotalPage += 1
	}

	http.Success(ctx, http.FlatOption(pager))
}

//QueryArticleDetail...
func (s *Handler) QueryArticleDetail(ctx *gin.Context) {
	var session = db.GetDB()

	inArticle := models.Article{}
	where := models.Article{}
	where.ID = util.MustUInt(ctx.Param("id"))
	if where.ID == 0 {
		http.Fail(ctx, http.MsgOption("ID 未设置"))
		return
	}

	err := s.QueryArticleDetailDB(ctx, session, &where, &inArticle)
	if err != nil {
		http.Fail(ctx, http.MsgOption("QueryArticleDetailDB failed. [%s]", err.Error()))
		return
	}

	http.Success(ctx, http.FlatOption(inArticle))
}

//UpdateArticle...
func (s *Handler) UpdateArticle(ctx *gin.Context) {
	var session = db.GetDB()

	inArticle := models.Article{}
	inArticle.ID = util.MustUInt(ctx.Param("id"))
	if inArticle.ID == 0 {
		http.Fail(ctx, http.MsgOption("ID 未设置"))
		return
	}

	var articleForm types.ArticleForm
	if err := ctx.ShouldBindJSON(&articleForm); err != nil {
		http.Fail(ctx, http.MsgOption(articleForm.Insert(err)))
		return
	}
	copier.Copy(&inArticle, &articleForm)

	err := s.UpdateArticleDB(ctx, session, &inArticle)
	if err != nil {
		http.Fail(ctx, http.MsgOption("UpdateArticleDB failed. [%v]", err))
		return
	}

	http.Success(ctx, http.FlatOption(inArticle))
}

//DeleteArticle...
func (s *Handler) DeleteArticle(ctx *gin.Context) {
	var session = db.GetDB()

	where := models.Article{}
	where.ID = util.MustUInt(ctx.Param("id"))
	if where.ID == 0 {
		http.Fail(ctx, http.MsgOption("ID 未设置"))
		return
	}

	err := s.DeleteArticleDB(ctx, session, &where)
	if err != nil {
		http.Fail(ctx, http.MsgOption("DeleteArticleDB failed. [%v]", err))
		return
	}

	http.Success(ctx, http.MsgOption("Delete successfully"))
}
