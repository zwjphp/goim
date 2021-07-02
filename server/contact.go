package server

import "goim/model"

type ContactService struct {
}

// 获取群信息
func (service *ContactService) SearchComunityIds(userId int64) (comIds []int64) {
	// TODO 获取用户全部群ID
	conconts := make([]model.Contact, 0)
	comIds = make([]int64, 0)
	DbEngin.Where("ownerid = ? and cate = ?", userId, model.CONCAT_CATE_COMUNITY).Find(&conconts)
	for _, v := range conconts {
		comIds = append(comIds, v.Dstobj)
	}
	return comIds
}
