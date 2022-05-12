package handler

import (
	"context"
	"gin-boilerplate/comm/errors"
	"gin-boilerplate/models"

	"github.com/jinzhu/gorm"
)

func (s *Handler) QueryOptionsetDB(ctx context.Context, session *gorm.DB, where *models.Optionset, list *[]models.Optionset, count ...*int32) error {
	session = session.Table(where.TableName()).Where(where).Find(list)
	if len(count) > 0 {
		session = session.Offset(0).Count(count[0])
	}

	if errs := session.GetErrors(); len(errs) != 0 {
		return errors.New(errors.ERecordFindFailed, "查询Optionset失败. [%v]", errs)
	}
	return nil
}

func (s *Handler) QueryOptionsetDetailDB(ctx context.Context, session *gorm.DB, where *models.Optionset, data *models.Optionset) error {
	var err error
	var lst []models.Optionset
	s.QueryOptionsetDB(ctx, session, where, &lst)
	if err != nil {
		return err
	}
	if len(lst) == 0 {
		return errors.New(errors.ERecordNotFound, "查询Optionset为空")
	}
	*data = lst[0]
	return nil
}

func (s *Handler) InsertOptionsetDB(ctx context.Context, session *gorm.DB, data *models.Optionset) error {
	if err := session.Create(data).Error; err != nil {
		return errors.New(errors.ERecordCreateFailed, "新增Optionset失败. [%s]", err.Error())
	}
	return nil
}

func (s *Handler) UpdateOptionsetDB(ctx context.Context, session *gorm.DB, data *models.Optionset) error {
	if err := session.Table(data.TableName()).Model(&data).Updates(&data).Error; err != nil {
		return errors.New(errors.ERecordUpdateFailed, "更新Optionset失败. [%s]", err.Error())
	}
	return nil
}

//SaveOptionsetDB 存在primarykey update 否则 insert
func (s *Handler) SaveOptionsetDB(ctx context.Context, session *gorm.DB, data *models.Optionset) error {
	if err := session.Save(data).Error; err != nil {
		return errors.New(errors.ERecordSaveFailed, "保存Optionset失败. [%s]", err.Error())
	}
	return nil
}

//DeleteOptionsetDB 删除
func (s *Handler) DeleteOptionsetDB(ctx context.Context, session *gorm.DB, data *models.Optionset) error {
	if err := session.Where(data).Delete(&data).Error; err != nil {
		return errors.New(errors.ERecordDeleteFailed, "删除Optionset失败. [%s]", err.Error())
	}
	return nil
}
