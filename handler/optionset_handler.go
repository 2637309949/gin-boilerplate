package handler

import (
	"gin-boilerplate/comm/db"
	"gin-boilerplate/comm/http"
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
// @Router /api/v1/QueryOptionset [get]
func (s *Handler) QueryOptionset(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "QueryOptionset")()

	var articleFilter types.OptionsetFilter
	if err := ctx.Bind(&articleFilter); err != nil {
		http.Fail(ctx, http.MsgOption(articleFilter.Filter(err)))
		return
	}
	var session = db.GetDB()
	session = db.SetLimit(ctx, session, &articleFilter)
	session = db.SetOrder(ctx, session, &articleFilter)
	timemark.Mark("InitDb")

	var totalCount int32
	var lst []models.Optionset
	where := models.Optionset{
		Name: articleFilter.Name,
	}
	err := s.QueryOptionsetDB(ctx, session, &where, &lst, &totalCount)
	timemark.Mark("QueryOptionsetDB")
	if err != nil {
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

//QueryOptionsetDetail...
func (s *Handler) QueryOptionsetDetail(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "QueryOptionsetDetail")()

	var session = db.GetDB()
	timemark.Mark("InitDb")

	inOptionset := models.Optionset{}
	where := models.Optionset{}
	where.ID = util.MustUInt(ctx.Param("id"))
	if where.ID == 0 {
		http.Fail(ctx, http.MsgOption("ID is not set"))
		return
	}

	err := s.QueryOptionsetDetailDB(ctx, session, &where, &inOptionset)
	timemark.Mark("QueryOptionsetDetailDB")
	if err != nil {
		http.Fail(ctx, http.MsgOption("QueryOptionsetDetailDB failed. [%s]", err.Error()))
		return
	}

	http.Success(ctx, http.FlatOption(inOptionset))
}

//InsertOptionset...
func (s *Handler) InsertOptionset(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "InsertOptionset")()

	var session = db.GetDB()
	timemark.Mark("InitDb")

	var articleForm types.OptionsetForm
	if err := ctx.ShouldBindJSON(&articleForm); err != nil {
		http.Fail(ctx, http.MsgOption(articleForm.Insert(err)))
		return
	}

	inOptionset := models.Optionset{}
	where := models.Optionset{
		Name: articleForm.Name,
	}
	copier.Copy(&inOptionset, &articleForm)
	err := s.QueryOptionsetDetailDB(ctx, session, &where, &inOptionset)
	timemark.Mark("QueryOptionsetDetailDB")
	if err == nil {
		http.Fail(ctx, http.MsgOption("Record already exists"))
		return
	}

	err = s.InsertOptionsetDB(ctx, session, &inOptionset)
	timemark.Mark("InsertOptionsetDB")
	if err != nil {
		http.Fail(ctx, http.MsgOption("InsertSchedulePositionDB failed. [%v]", err))
		return
	}

	http.Success(ctx, http.FlatOption(inOptionset))
}

//UpdateOptionset...
func (s *Handler) UpdateOptionset(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "UpdateOptionset")()

	var session = db.GetDB()
	timemark.Mark("InitDb")

	inOptionset := models.Optionset{}
	inOptionset.ID = util.MustUInt(ctx.Param("id"))
	if inOptionset.ID == 0 {
		http.Fail(ctx, http.MsgOption("ID is not set"))
		return
	}

	var articleForm types.OptionsetForm
	if err := ctx.ShouldBindJSON(&articleForm); err != nil {
		http.Fail(ctx, http.MsgOption(articleForm.Insert(err)))
		return
	}
	copier.Copy(&inOptionset, &articleForm)

	err := s.UpdateOptionsetDB(ctx, session, &inOptionset)
	timemark.Mark("UpdateOptionsetDB")
	if err != nil {
		http.Fail(ctx, http.MsgOption("UpdateOptionsetDB failed. [%v]", err))
		return
	}

	http.Success(ctx, http.FlatOption(inOptionset))
}

//DeleteOptionset...
func (s *Handler) DeleteOptionset(ctx *gin.Context) {
	var timemark mark.TimeMark
	defer timemark.Init(ctx.Request.Context(), "DeleteOptionset")()

	var session = db.GetDB()
	timemark.Mark("InitDb")
	where := models.Optionset{}
	where.ID = util.MustUInt(ctx.Param("id"))
	if where.ID == 0 {
		http.Fail(ctx, http.MsgOption("ID is not set"))
		return
	}

	err := s.DeleteOptionsetDB(ctx, session, &where)
	timemark.Mark("DeleteOptionsetDB")
	if err != nil {
		http.Fail(ctx, http.MsgOption("DeleteOptionsetDB failed. [%v]", err))
		return
	}

	http.Success(ctx, http.MsgOption("Delete successfully"))
}
