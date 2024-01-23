package svc

import (
	"management_system/service/api/internal/config"
	"management_system/service/models"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                config.Config // Config
	RedisClient           *redis.Redis
	UserModel             models.UserModel
	UserAuthModel         models.UserAuthModel
	ApplicationFormModel  models.ApplicationFormModel
	RoleChangeModel       models.RoleChangeModel
	PendingPersonnelModel models.PendingPersonnelModel
	ObjectiveExamModel    models.ObjectiveExamModel
	SubjectiveExamModel   models.SubjectiveExamModel
	RescueInfoModel       models.RescueInfoModel
	RescueTargetModel     models.RescueTargetModel
	SignatureModel        models.SignatureModel
	RescueProcessModel    models.RescueProcessModel
	FileModel             models.FileModel
	FolderModel           models.FolderModel
	ViewingRecordModel    models.ViewingRecordModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
		UserModel:             models.NewUserModel(sqlx.NewMysql(c.DB.DataSource)),
		UserAuthModel:         models.NewUserAuthModel(sqlx.NewMysql(c.DB.DataSource)),
		ApplicationFormModel:  models.NewApplicationFormModel(sqlx.NewMysql(c.DB.DataSource)),
		RoleChangeModel:       models.NewRoleChangeModel(sqlx.NewMysql(c.DB.DataSource)),
		PendingPersonnelModel: models.NewPendingPersonnelModel(sqlx.NewMysql(c.DB.DataSource)),
		ObjectiveExamModel:    models.NewObjectiveExamModel(sqlx.NewMysql(c.DB.DataSource)),
		SubjectiveExamModel:   models.NewSubjectiveExamModel(sqlx.NewMysql(c.DB.DataSource)),
		RescueInfoModel:       models.NewRescueInfoModel(sqlx.NewMysql(c.DB.DataSource)),
		RescueTargetModel:     models.NewRescueTargetModel(sqlx.NewMysql(c.DB.DataSource)),
		SignatureModel:        models.NewSignatureModel(sqlx.NewMysql(c.DB.DataSource)),
		RescueProcessModel:    models.NewRescueProcessModel(sqlx.NewMysql(c.DB.DataSource)),
		FileModel:             models.NewFileModel(sqlx.NewMysql(c.DB.DataSource)),
		FolderModel:           models.NewFolderModel(sqlx.NewMysql(c.DB.DataSource)),
		ViewingRecordModel:    models.NewViewingRecordModel(sqlx.NewMysql(c.DB.DataSource)),
	}
}
