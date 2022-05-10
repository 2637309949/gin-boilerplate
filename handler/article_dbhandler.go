package handler

import (
	"context"
	"gin-boilerplate/comm/errors"
	"gin-boilerplate/models"

	"github.com/jinzhu/gorm"
)

func (s *Handler) QueryArticleDB(ctx context.Context, session *gorm.DB, where *models.Article, list *[]models.Article, count ...*int32) error {
	session = session.Table(where.TableName()).Where(where).Find(list)
	if len(count) > 0 {
		session = session.Offset(0).Count(count[0])
	}

	if errs := session.GetErrors(); len(errs) != 0 {
		return errors.New(errors.ERecordFindFailed, "查询Article失败. [%v]", errs)
	}
	return nil
}

func (s *Handler) QueryArticleDetailDB(ctx context.Context, session *gorm.DB, where *models.Article, data *models.Article) error {
	var err error
	var lst []models.Article
	s.QueryArticleDB(ctx, session, where, &lst)
	if err != nil {
		return err
	}
	if len(lst) == 0 {
		return errors.New(errors.ERecordNotFound, "查询Article为空")
	}
	*data = lst[0]
	return nil
}

func (s *Handler) InsertArticleDB(ctx context.Context, session *gorm.DB, data *models.Article) error {
	if err := session.Create(data).Error; err != nil {
		return errors.New(errors.ERecordCreateFailed, "新增Article失败. [%s]", err.Error())
	}
	return nil
}

func (s *Handler) UpdateArticleDB(ctx context.Context, session *gorm.DB, data *models.Article) error {
	if err := session.Table(data.TableName()).Model(&data).Updates(&data).Error; err != nil {
		return errors.New(errors.ERecordUpdateFailed, "更新Article失败. [%s]", err.Error())
	}
	return nil
}

//SaveArticleDB 存在primarykey update 否则 insert
func (s *Handler) SaveArticleDB(ctx context.Context, session *gorm.DB, data *models.Article) error {
	if err := session.Save(data).Error; err != nil {
		return errors.New(errors.ERecordSaveFailed, "保存Article失败. [%s]", err.Error())
	}
	return nil
}

//DeleteArticleDB 删除
func (s *Handler) DeleteArticleDB(ctx context.Context, session *gorm.DB, data *models.Article) error {
	if err := session.Where(data).Delete(&data).Error; err != nil {
		return errors.New(errors.ERecordDeleteFailed, "删除Article失败. [%s]", err.Error())
	}
	return nil
}
