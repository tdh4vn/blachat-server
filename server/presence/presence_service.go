package presence

import (
	"blachat-server/repo"
	"blachat-server/services"
)

type PresenceService interface {
	SendUserOnlineEvent(userId string) error
	SendUserOfflineEvent(userId string) error
}

type _PresenceServiceImpl struct {
	ContactRepo repo.ContactsRepo
	PresenceRepo repo.PresenceRepo
}

func (this _PresenceServiceImpl) SendUserOnlineEvent(userId string) error {
	//lấy tất cả user mà userid này có trong danh bạ của họ
	if ids, err := this.ContactRepo.GetUsersRelated(userId); err != nil {
		return err
	} else {
		var userIds []string
		for _, id := range ids {
			userIds = append(userIds, id.String())
		}
		//lấy tất cả những thằng online trong đó
		usersOnline := this.PresenceRepo.CheckUsersOnline(userIds)

		//gửi về cho danh sách user này
		if len(usersOnline) > 0 {
			services.SendUserOnline(userId, usersOnline)
		}
		return nil
	}
}


func (this _PresenceServiceImpl) SendUserOfflineEvent(userId string) error {
	//lấy tất cả user mà userid này có trong danh bạ của họ
	if ids, err := this.ContactRepo.GetUsersRelated(userId); err != nil {
		return err
	} else {
		var userIds []string
		for _, id := range ids {
			userIds = append(userIds, id.String())
		}
		//lấy tất cả những thằng online trong đó
		usersOnline := this.PresenceRepo.CheckUsersOnline(userIds)

		//gửi về cho danh sách user này
		if len(usersOnline) > 0 {
			services.SendUserOnline(userId, usersOnline)
		}
		return nil
	}
}

func NewPresenceService() PresenceService {
	return _PresenceServiceImpl{}
}

func ServePresence() {

}