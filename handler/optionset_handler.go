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
// @Summary QueryOptionset
// @Description get optionset by give params
// @Tags optionsets
// @Accept  json
// @Produce  json
// @Param  page_no  query  int  true  "page no"
// @Param  page_size  query  int  true  "page size"
// @Param  order_type  query  int  true  "order type"
// @Param  order_col  query  int  true  "order col"
// @Router /api/v1/queryOptionset [get]
func (h *Handler) QueryOptionset(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "QueryOptionset")()

	var articleFilter types.OptionsetFilter
	if err := ctx.Bind(&articleFilter); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(articleFilter.Filter(err)))
		return
	}
	var session = db.GetDB()
	session = db.SetLimit(ctx, session, &articleFilter)
	session = db.SetOrder(ctx, session, &articleFilter)
	timemark.Mark("InitDb")

	var totalCount int64
	var lst []models.Optionset
	where := models.Optionset{
		Name: articleFilter.Name,
	}
	err := h.QueryOptionsetDB(ctx.Request.Context(), session, &where, &lst, &totalCount)
	timemark.Mark("QueryOptionsetDB")
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption("QueryOptionsetDB failed. [%s]", err.Error()))
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
// @Summary QueryOptionsetDetail
// @Description get optionset by gived id
// @Tags optionsets
// @Accept  json
// @Produce  json
// @Router /api/v1/queryOptionsetDetail [get]
func (h *Handler) QueryOptionsetDetail(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "QueryOptionsetDetail")()

	var session = db.GetDB()
	timemark.Mark("InitDb")

	inOptionset := models.Optionset{}
	where := models.Optionset{}
	where.ID = util.MustUInt(ctx.Param("id"))
	if where.ID == 0 {
		logger.Errorf(ctx.Request.Context(), "ID is not set")
		http.Fail(ctx, http.MsgOption("ID is not set"))
		return
	}

	err := h.QueryOptionsetDetailDB(ctx.Request.Context(), session, &where, &inOptionset)
	timemark.Mark("QueryOptionsetDetailDB")
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption("QueryOptionsetDetailDB failed. [%s]", err.Error()))
		return
	}

	http.Success(ctx, http.FlatOption(inOptionset))
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary InsertOptionset
// @Description new article
// @Tags optionsets
// @Accept  json
// @Produce  json
// @Router /api/v1/insertOptionset [post]
func (h *Handler) InsertOptionset(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "InsertOptionset")()

	var session = db.GetDB()
	timemark.Mark("InitDb")

	var articleForm types.OptionsetForm
	if err := ctx.ShouldBindJSON(&articleForm); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(articleForm.Insert(err)))
		return
	}

	where, inOptionset := models.Optionset{
		Name: articleForm.Name,
	}, models.Optionset{}
	copier.Copy(&inOptionset, &articleForm)
	err := h.QueryOptionsetDetailDB(ctx.Request.Context(), session, &where, &inOptionset)
	timemark.Mark("QueryOptionsetDetailDB")
	if err == nil {
		logger.Errorf(ctx.Request.Context(), "Record already exists")
		http.Fail(ctx, http.MsgOption("Record already exists"))
		return
	}

	err = h.InsertOptionsetDB(ctx.Request.Context(), session, &inOptionset)
	timemark.Mark("InsertOptionsetDB")
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption("InsertSchedulePositionDB failed. [%v]", err))
		return
	}

	http.Success(ctx, http.FlatOption(inOptionset))
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary UpdateOptionset
// @Description update optionset
// @Tags optionsets
// @Accept  json
// @Produce  json
// @Router /api/v1/updateOptionset [post]
func (h *Handler) UpdateOptionset(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "UpdateOptionset")()

	var session = db.GetDB()
	timemark.Mark("InitDb")

	inOptionset := models.Optionset{}
	inOptionset.ID = util.MustUInt(ctx.Param("id"))
	if inOptionset.ID == 0 {
		logger.Errorf(ctx.Request.Context(), "ID is not set")
		http.Fail(ctx, http.MsgOption("ID is not set"))
		return
	}

	var articleForm types.OptionsetForm
	if err := ctx.ShouldBindJSON(&articleForm); err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption(articleForm.Insert(err)))
		return
	}
	copier.Copy(&inOptionset, &articleForm)

	err := h.UpdateOptionsetDB(ctx.Request.Context(), session, &inOptionset)
	timemark.Mark("UpdateOptionsetDB")
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption("UpdateOptionsetDB failed. [%v]", err))
		return
	}

	http.Success(ctx, http.FlatOption(inOptionset))
}

// ref: https://swaggo.github.io/swaggo.io/declarative_comments_format/api_operation.html
// @Summary DeleteOptionset
// @Description delete optionset
// @Tags optionsets
// @Accept  json
// @Produce  json
// @Router /api/v1/deleteOptionset [post]
func (h *Handler) DeleteOptionset(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "DeleteOptionset")()

	var session = db.GetDB()
	timemark.Mark("InitDb")
	where := models.Optionset{}
	where.ID = util.MustUInt(ctx.Param("id"))
	if where.ID == 0 {
		logger.Errorf(ctx.Request.Context(), "ID is not set")
		http.Fail(ctx, http.MsgOption("ID is not set"))
		return
	}

	err := h.DeleteOptionsetDB(ctx.Request.Context(), session, &where)
	timemark.Mark("DeleteOptionsetDB")
	if err != nil {
		logger.Error(ctx.Request.Context(), err)
		http.Fail(ctx, http.MsgOption("DeleteOptionsetDB failed. [%v]", err))
		return
	}

	http.Success(ctx, http.MsgOption("Delete successfully"))
}
