package handler

import (
	"context"
	"gin-boilerplate/comm/errors"
	"gin-boilerplate/models"

	"github.com/jinzhu/gorm"
)

func (h *Handler) QueryUserDB(ctx context.Context, session *gorm.DB, where *models.User, list *[]models.User, count ...*int32) error {
	session = session.Table(where.TableName()).Where(where).Find(list)
	if len(count) > 0 {
		session = session.Offset(0).Count(count[0])
	}

	if errs := session.GetErrors(); len(errs) != 0 {
		return errors.New(errors.ERecordFindFailed, "查询User失败. [%v]", errs)
	}
	return nil
}

func (h *Handler) QueryUserDetailDB(ctx context.Context, session *gorm.DB, where *models.User, data *models.User) error {
	var err error
	var lst []models.User
	h.QueryUserDB(ctx, session, where, &lst)
	if err != nil {
		return err
	}
	if len(lst) == 0 {
		return errors.New(errors.ERecordNotFound, "查询User为空")
	}
	*data = lst[0]
	return nil
}

func (h *Handler) InsertUserDB(ctx context.Context, session *gorm.DB, data *models.User) error {
	if err := session.Create(data).Error; err != nil {
		return errors.New(errors.ERecordCreateFailed, "新增User失败. [%s]", err.Error())
	}
	return nil
}

func (h *Handler) UpdateUserDB(ctx context.Context, session *gorm.DB, data *models.User) error {
	if err := session.Table(data.TableName()).Model(&data).Updates(&data).Error; err != nil {
		return errors.New(errors.ERecordUpdateFailed, "更新User失败. [%s]", err.Error())
	}
	return nil
}

//SaveUserDB 存在primarykey update 否则 insert
func SaveUserDB(ctx context.Context, session *gorm.DB, data *models.User) error {
	if err := session.Save(data).Error; err != nil {
		return errors.New(errors.ERecordSaveFailed, "保存User失败. [%s]", err.Error())
	}
	return nil
}

//DeleteUserDB 删除
func DeleteUserDB(ctx context.Context, session *gorm.DB, data *models.User) error {
	if err := session.Where(data).Delete(&data).Error; err != nil {
		return errors.New(errors.ERecordDeleteFailed, "删除User失败. [%s]", err.Error())
	}
	return nil
}
