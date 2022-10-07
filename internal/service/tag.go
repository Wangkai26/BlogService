package service

import (
	"github.com/go-programming-tour-book/blog-service/internal/model"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
)

// form 和 binding 分别代表着表单的映射字段名和入参校验的规则内容
//其主要功能是实现参数绑定和参数校验
type CountTagRequest struct {
	Name  string `form:"name" binding:"max=100"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

//查看标签请求
type TagListRequest struct {
	Name  string `form:"name" binding:"max=100"`
	State uint8 `form:"state,default=1" binding:"oneof=0 1"`
}

// 增加标签请求
type CreateTagRequest struct {
	Name     string `form:"name" binding:"required,min=3,max=100"`
	CreatedBy string `form:"created_by" binding:"required,min=3,max=100"`
	State    uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

//更新标签请求
type UpdateTagRequest struct {
	ID 		uint32 `form:"id" binding:"required,gte=1"`
	Name 	string `form:"name" binding:"max=100"`
	State 	uint8  `form:"state" binding:"oneof=0 1"`
	ModifiedBy string `form:"modified_by" binding:"required,min=3,max=100"`
}

//删除标签请求
type DeleteTagRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}


//下述代码中，定义了Request结构体作为接口入参的基准，而本项目由于不会太复杂，所以直接放在了service层中便于使用
//若后续业务不断增长，程序越来越复杂，service也冗杂起来了
//可以考虑抽离-层接口校验层，便于解耦逻辑

func (svc *Service) CountTag(param *CountTagRequest) (int, error) {
	return svc.dao.CountTag(param.Name, param.State)
}

func (svc *Service) GetTagList(param *TagListRequest, pager *app.Pager) ([]*model.Tag, error) {
	return svc.dao.GetTagList(param.Name, param.State, pager.Page, pager.PageSize)
}

func (svc *Service) CreateTag(param *CreateTagRequest) error {
	return svc.dao.CreateTag(param.Name, param.State, param.CreatedBy)
}

func (svc *Service) UpdateTag(param *UpdateTagRequest) error {
	return svc.dao.UpdateTag(param.ID, param.Name, param.State, param.ModifiedBy)
}

func (svc *Service) DeleteTag(param *DeleteTagRequest) error {
	return svc.dao.DeleteTag(param.ID)
}