package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	pb "forum-chat/proto"
	"forum/log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	RedisPrefixKey = "TeaHouse"
	Chat           = "chat"
)

func GetKey(id uint32) string {
	return RedisPrefixKey + ":" + Chat + ":" + strconv.Itoa(int(id))
}

func (d *Dao) Create(data *ChatData) error {
	//序列化为string后推入数列左侧(越靠左的消息越新)
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return d.Redis.LPush(GetKey(data.Receiver), msg).Err()
}

func (d *Dao) GetList(id uint32, expiration time.Duration, wait bool) ([]string, error) {
	t := time.Now()
	defer func() {
		fmt.Println(time.Now().Sub(t))
	}()
	// 使用用户的id创建一个key
	key := GetKey(id)
	// 如果数列里面为空的话则阻塞等待
	if d.Redis.LLen(key).Val() == 0 {
		if !wait {
			return nil, nil
		}
		msg, err := d.Redis.BRPop(expiration, key).Result() // 阻塞
		if err != nil {
			if errors.Is(err, redis.Nil) {
				// 超时但没拿到值，是正常现象
				return nil, nil
			}
			return nil, err
		}

		return msg[1:], nil // first ele is key
	}

	// 逐个读取整个队列,从右侧取(先取旧的消息)
	var list []string
	for d.Redis.LLen(key).Val() != 0 {
		msg, err := d.Redis.RPop(key).Result()
		if err != nil {
			return nil, err
		}

		list = append(list, msg)
	}

	return list, nil
}

// Rewrite 未成功发送的消息逆序放回list的Right
func (d *Dao) Rewrite(id uint32, list []string) error {
	log.Info("Rewrite")

	//从右侧将消息回写,这个地方为了保持消息的顺序与原来保持一致需要倒序写入
	for i := len(list); i > 0; i-- {
		if err := d.Redis.RPush(GetKey(id), list[i-1]).Err(); err != nil {
			return err
		}
	}

	return nil
}

//func (d *Dao) ListHistory(userId, otherUserId, offset, limit uint32, pagination bool) ([]*pb.Message, error) {
//	//这个地方为了保证顺序要进行交换
//	if otherUserId < userId {
//		otherUserId, userId = userId, otherUserId
//	}
//
//	//读取这两个人的历史消息
//	key := "history:" + strconv.Itoa(int(userId)) + "-" + strconv.Itoa(int(otherUserId)) // history:min_id-max_id
//
//	var start int64 = 0
//	var end int64 = -1
//
//	if pagination {
//		start = int64(offset)
//		end = int64(offset + limit)
//	}
//
//	//批量读取消息,这个地方要注意索引从小到大是从新到旧
//	list, err := d.Redis.LRange(key, start, end).Result() // DESC
//	if err != nil {
//		return nil, err
//	}
//
//	histories := make([]*pb.Message, len(list))
//	for i, history := range list {
//		var msg pb.Message
//		if err := json.Unmarshal([]byte(history), &msg); err != nil {
//			return nil, err
//		}
//		histories[i] = &msg
//	}
//
//	return histories, nil
//}

func (d *Dao) ListHistory(userId, otherUserId, offset, limit uint32, pagination bool) ([]*pb.Message, error) {
	var history []*pb.Message

	if pagination {
		err := d.DB.Table("messages").Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userId, otherUserId, otherUserId, userId).Order("time DESC").Offset(int(offset)).Limit(int(limit)).Find(&history).Error
		if err != nil {
			return nil, err
		}
	} else {
		if err := d.DB.Table("messages").Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userId, otherUserId, otherUserId, userId).Order("time DESC").Find(&history).Error; err != nil {
			return nil, err
		}
	}

	return history, nil
}

//func (d *Dao) CreateHistory(userId uint32, list []string) error {
//	log.Info("CreateHistory")
//
//	for i := len(list); i > 0; i-- {
//		var msg ChatData
//		if err := json.Unmarshal([]byte(list[i-1]), &msg); err != nil {
//			return err
//		}
//
//		//调整为统一顺序,为什么要这么做呢?因为这么做首先可以保证一致性,其次可以减半存储空间
//		minId := userId
//		if minId > msg.Sender {
//			minId, msg.Sender = msg.Sender, minId
//		}
//
//		//推送到这两个人的历史消息里面去
//		if err := d.Redis.LPush("history:"+strconv.Itoa(int(minId))+"-"+strconv.Itoa(int(msg.Sender)), list[i-1]).Err(); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

func (d *Dao) CreateHistory(userId uint32, list []string) error {
	for i := len(list); i > 0; i-- {
		var msg ChatData
		if err := json.Unmarshal([]byte(list[i-1]), &msg); err != nil {
			return err
		}

		data := DBdata{
			ReceiverID: userId,
			SenderID:   msg.Sender,
			Content:    msg.Content,
			TypeName:   msg.TypeName,
			Time:       msg.Time,
		}

		if err := d.DB.Table("messages").Create(&data).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *Dao) GetUserList(userId uint32, limit, page int) ([]uint32, error) {
	var userList []uint32

	// order by id 让最新聊天的用户排在最前面
	err := d.DB.Table("messages").
		Select("CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END AS other_id", userId).
		Where("sender_id = ? OR receiver_id = ?", userId, userId).
		Group("other_id").
		Order("MAX(time) DESC").
		Offset(page * limit).
		Limit(limit).
		Scan(&userList).Error

	if err != nil {
		return nil, err
	}

	return userList, nil
}

func (d *Dao) GetUserById(userIds []uint32) ([]*pb.UserStatus, error) {
	var usersStatus []*pb.UserStatus
	for _, id := range userIds {
		var userStatus pb.UserStatus
		err := d.DB.Table("users").Where("id = ?", id).Find(&userStatus).Error
		if err != nil {
			return nil, err
		}

		usersStatus = append(usersStatus, &userStatus)
	}

	return usersStatus, nil
}
