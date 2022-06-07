package handler

import (
	"gin-boilerplate/comm/db"
	"gin-boilerplate/comm/http"
	"gin-boilerplate/comm/logger"
	"gin-boilerplate/comm/mark"
	"gin-boilerplate/comm/util"
	"gin-boilerplate/models"
	"gin-boilerplate/types"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary QueryArticle
// @Description get article by gived params
// @Tags articles
// @Accept  json
// @Produce  json
// @Param  page_no  query  int  true  "page no"
// @Param  page_size  query  int  true  "page size"
// @Param  order_type  query  int  true  "order type"
// @Param  order_col  query  int  true  "order col"
// @Router /api/v1/queryArticle [get]
func (h *Handler) QueryArticle(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "QueryArticle")()
	var articleFilter types.ArticleFilter
	if err := ctx.Bind(&articleFilter); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(articleFilter.Filter(err)))
		return
	}
	var session = db.GetDB()
	session = db.SetLimit(ctx, session, &articleFilter)
	session = db.SetOrder(ctx, session, &articleFilter)
	timemark.Mark("InitDb")
	where, lst, totalCount := models.Article{
		Title: articleFilter.Title,
	}, []models.Article{}, int64(0)
	err := h.QueryArticleDB(ctx, session, &where, &lst, &totalCount)
	timemark.Mark("QueryArticleDB")
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
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

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary QueryArticleDetail
// @Description get article by gived id
// @Tags articles
// @Accept  json
// @Produce  json
// @Router /api/v1/queryArticleDetail [get]
func (h *Handler) QueryArticleDetail(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "QueryArticleDetail")()

	var session = db.GetDB()
	timemark.Mark("InitDb")

	inArticle, where := models.Article{}, models.Article{}
	where.ID = util.MustUInt(ctx.Param("id"))
	if where.ID == 0 {
		logger.Errorf(ctx.Request.Context(), "ID is not set")
		http.Fail(ctx, http.MsgOption("ID is not set"))
		return
	}

	err := h.QueryArticleDetailDB(ctx, session, &where, &inArticle)
	timemark.Mark("QueryArticleDetailDB")
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption("QueryArticleDetailDB failed. [%s]", err.Error()))
		return
	}

	http.Success(ctx, http.FlatOption(inArticle))
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary InsertArticle
// @Description new article
// @Tags articles
// @Accept  json
// @Produce  json
// @Router /api/v1/insertArticle [post]
func (h *Handler) InsertArticle(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "InsertArticle")()

	var session = db.GetDB()
	timemark.Mark("InitDb")

	var articleForm types.ArticleForm
	if err := ctx.ShouldBindJSON(&articleForm); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(articleForm.Insert(err)))
		return
	}

	where, inArticle := models.Article{
		Title: articleForm.Title,
	}, models.Article{}
	copier.Copy(&inArticle, &articleForm)
	err := h.QueryArticleDetailDB(ctx, session, &where, &inArticle)
	timemark.Mark("QueryArticleDetailDB")
	if err == nil {
		logger.Errorf(ctx.Request.Context(), "Record already exists")
		http.Fail(ctx, http.MsgOption("Record already exists"))
		return
	}

	err = h.InsertArticleDB(ctx, session, &inArticle)
	timemark.Mark("InsertArticleDB")
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption("InsertSchedulePositionDB failed. [%v]", err))
		return
	}

	http.Success(ctx, http.FlatOption(inArticle))
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary UpdateArticle
// @Description update article
// @Tags articles
// @Accept  json
// @Produce  json
// @Router /api/v1/updateArticle [post]
func (h *Handler) UpdateArticle(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "UpdateArticle")()

	var session = db.GetDB()
	timemark.Mark("InitDb")

	inArticle := models.Article{}
	inArticle.ID = util.MustUInt(ctx.Param("id"))
	if inArticle.ID == 0 {
		logger.Errorf(ctx.Request.Context(), "ID is not set")
		http.Fail(ctx, http.MsgOption("ID is not set"))
		return
	}

	var articleForm types.ArticleForm
	if err := ctx.ShouldBindJSON(&articleForm); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(articleForm.Insert(err)))
		return
	}
	copier.Copy(&inArticle, &articleForm)

	err := h.UpdateArticleDB(ctx, session, &inArticle)
	timemark.Mark("UpdateArticleDB")
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption("UpdateArticleDB failed. [%v]", err))
		return
	}

	http.Success(ctx, http.FlatOption(inArticle))
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary DeleteArticle
// @Description delete article
// @Tags articles
// @Accept  json
// @Produce  json
// @Router /api/v1/deleteArticle [post]
func (h *Handler) DeleteArticle(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "DeleteArticle")()

	var session = db.GetDB()
	timemark.Mark("InitDb")
	where := models.Article{}
	where.ID = util.MustUInt(ctx.Param("id"))
	if where.ID == 0 {
		logger.Errorf(ctx.Request.Context(), "ID is not set")
		http.Fail(ctx, http.MsgOption("ID is not set"))
		return
	}

	err := h.DeleteArticleDB(ctx, session, &where)
	timemark.Mark("DeleteArticleDB")
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption("DeleteArticleDB failed. [%v]", err))
		return
	}

	http.Success(ctx, http.MsgOption("Delete successfully"))
}
