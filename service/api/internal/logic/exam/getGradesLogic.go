package exam

import (
	"context"
	"encoding/json"

	"management_system/common/response"
	"management_system/service/api/internal/svc"
	"management_system/service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGradesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGradesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGradesLogic {
	return &GetGradesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGradesLogic) GetGrades() (resp *types.GetGradesResp, err error) {
	logx.Infof("userId: %v", l.ctx.Value("userId"))
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()

	objectiveGrades, err := l.svcCtx.ObjectiveExamModel.FindAllByUserId(l.ctx, l.svcCtx.ObjectiveExamModel.RowBuilder(), userId)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	subjectiveGrades, err := l.svcCtx.SubjectiveExamModel.FindAllByUserId(l.ctx, l.svcCtx.SubjectiveExamModel.RowBuilder(), userId)
	if err != nil {
		return nil, response.Error(500, err.Error())
	}

	var list []types.Grades
	if len(objectiveGrades) > 0 {
		for _, examModel := range objectiveGrades {
			if examModel.DelState == 0 {
				var typeExam types.Grades
				typeExam.Id = examModel.Id
				typeExam.UserId = examModel.UserId
				typeExam.Result = examModel.Result
				typeExam.Time = examModel.Time
				typeExam.Type = "客观题"
				list = append(list, typeExam)
			}
		}
	}

	if len(subjectiveGrades) > 0 {
		for _, examModel := range subjectiveGrades {
			if examModel.DelState == 0 {
				var typeExam types.Grades
				typeExam.Id = examModel.Id
				typeExam.UserId = examModel.UserId
				typeExam.Result = examModel.Result
				typeExam.Time = examModel.Time
				typeExam.Type = "主观题"
				list = append(list, typeExam)
			}
		}
	}

	return &types.GetGradesResp{
		List: list,
	}, nil
}
