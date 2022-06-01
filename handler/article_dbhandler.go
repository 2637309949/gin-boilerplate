package handler

import (
	"context"
	"gin-boilerplate/comm/errors"
	"gin-boilerplate/models"

	"gorm.io/gorm"
)

func (h *Handler) QueryArticleDB(ctx context.Context, session *gorm.DB, where *models.Article, list *[]models.Article, count ...*int64) error {
	session = session.Table(where.TableName()).Where(where).Find(list)
	if len(count) > 0 {
		session = session.Offset(0).Count(count[0])
	}

	if err := session.Error; err != nil {
		return errors.New(errors.ERecordFindFailed, "查询Article失败. [%v]", err)
	}
	return nil
}

func (h *Handler) QueryArticleDetailDB(ctx context.Context, session *gorm.DB, where *models.Article, data *models.Article) error {
	var err error
	var lst []models.Article
	h.QueryArticleDB(ctx, session, where, &lst)
	if err != nil {
		return err
	}
	if len(lst) == 0 {
		return errors.New(errors.ERecordNotFound, "查询Article为空")
	}
	*data = lst[0]
	return nil
}

func (h *Handler) InsertArticleDB(ctx context.Context, session *gorm.DB, data *models.Article) error {
	if err := session.Create(data).Error; err != nil {
		return errors.New(errors.ERecordCreateFailed, "新增Article失败. [%s]", err.Error())
	}
	return nil
}

func (h *Handler) UpdateArticleDB(ctx context.Context, session *gorm.DB, data *models.Article) error {
	if err := session.Table(data.TableName()).Model(&data).Updates(&data).Error; err != nil {
		return errors.New(errors.ERecordUpdateFailed, "更新Article失败. [%s]", err.Error())
	}
	return nil
}

//SaveArticleDB 存在primarykey update 否则 insert
func (h *Handler) SaveArticleDB(ctx context.Context, session *gorm.DB, data *models.Article) error {
	if err := session.Save(data).Error; err != nil {
		return errors.New(errors.ERecordSaveFailed, "保存Article失败. [%s]", err.Error())
	}
	return nil
}

//DeleteArticleDB 删除
func (h *Handler) DeleteArticleDB(ctx context.Context, session *gorm.DB, data *models.Article) error {
	if err := session.Where(data).Delete(&data).Error; err != nil {
		return errors.New(errors.ERecordDeleteFailed, "删除Article失败. [%s]", err.Error())
	}
	return nil
}
