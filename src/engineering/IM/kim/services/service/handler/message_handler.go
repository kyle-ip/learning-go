package handler

import (
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/kataras/iris/v12"
	"github.com/klintcheng/kim/services/service/database"
	"github.com/klintcheng/kim/wire"
	"github.com/klintcheng/kim/wire/rpc"
	"gorm.io/gorm"
)

type ServiceHandler struct {
	BaseDb    *gorm.DB
	MessageDb *gorm.DB
	Cache     *redis.Client
	Idgen     *database.IDGenerator
}

func (h *ServiceHandler) InsertUserMessage(c iris.Context) {
	var req rpc.InsertMessageReq
	if err := c.ReadBody(&req); err != nil {
		c.StopWithError(iris.StatusBadRequest, err)
		return
	}
	messageId, err := h.insertUserMessage(&req)
	if err != nil {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	_, _ = c.Negotiate(&rpc.InsertMessageResp{
		MessageId: messageId,
	})
}

func (h *ServiceHandler) insertUserMessage(req *rpc.InsertMessageReq) (int64, error) {
	messageId := h.Idgen.Next().Int64()
	messageContent := database.MessageContent{
		ID:       messageId,
		Type:     byte(req.Message.Type),
		Body:     req.Message.Body,
		Extra:    req.Message.Extra,
		SendTime: req.SendTime,
	}
	// 扩散写
	idxs := make([]database.MessageIndex, 2)
	idxs[0] = database.MessageIndex{
		ID:        h.Idgen.Next().Int64(),
		MessageID: messageId,
		AccountA:  req.Dest,
		AccountB:  req.Sender,
		Direction: 0,
		SendTime:  req.SendTime,
	}
	idxs[1] = database.MessageIndex{
		ID:        h.Idgen.Next().Int64(),
		MessageID: messageId,
		AccountA:  req.Sender,
		AccountB:  req.Dest,
		Direction: 1,
		SendTime:  req.SendTime,
	}

	err := h.MessageDb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&messageContent).Error; err != nil {
			return err
		}
		if err := tx.Create(&idxs).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return messageId, nil
}

func (h *ServiceHandler) InsertGroupMessage(c iris.Context) {
	var req rpc.InsertMessageReq
	if err := c.ReadBody(&req); err != nil {
		c.StopWithError(iris.StatusBadRequest, err)
		return
	}
	messageId, err := h.insertGroupMessage(&req)
	if err != nil {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	_, _ = c.Negotiate(&rpc.InsertMessageResp{
		MessageId: messageId,
	})
}

func (h *ServiceHandler) insertGroupMessage(req *rpc.InsertMessageReq) (int64, error) {
	messageId := h.Idgen.Next().Int64()

	var members []database.GroupMember
	err := h.BaseDb.Where(&database.GroupMember{Group: req.Dest}).Find(&members).Error
	if err != nil {
		return 0, err
	}
	// 扩散写
	var idxs = make([]database.MessageIndex, len(members))
	for i, m := range members {
		idxs[i] = database.MessageIndex{
			ID:        h.Idgen.Next().Int64(),
			MessageID: messageId,
			AccountA:  m.Account,
			AccountB:  req.Sender,
			Direction: 0,
			Group:     m.Group,
			SendTime:  req.SendTime,
		}
		if m.Account == req.Sender {
			idxs[i].Direction = 1
		}
	}

	messageContent := database.MessageContent{
		ID:       messageId,
		Type:     byte(req.Message.Type),
		Body:     req.Message.Body,
		Extra:    req.Message.Extra,
		SendTime: req.SendTime,
	}

	err = h.MessageDb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&messageContent).Error; err != nil {
			return err
		}
		if err := tx.Create(&idxs).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return messageId, nil
}

func (h *ServiceHandler) MessageAck(c iris.Context) {
	var req rpc.AckMessageReq
	if err := c.ReadBody(&req); err != nil {
		c.StopWithError(iris.StatusBadRequest, err)
		return
	}
	// save in redis
	err := setMesssageAck(h.Cache, req.Account, req.MessageId)
	if err != nil {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}
}

func setMesssageAck(cache *redis.Client, account string, msgId int64) error {
	if msgId == 0 {
		return nil
	}
	key := database.KeyMessageAckIndex(account)
	return cache.Set(key, msgId, wire.OfflineReadIndexExpiresIn).Err()
}

func (h *ServiceHandler) GetOfflineMessageIndex(c iris.Context) {
	var req rpc.GetOfflineMessageIndexReq
	if err := c.ReadBody(&req); err != nil {
		c.StopWithError(iris.StatusBadRequest, err)
		return
	}
	msgId := req.MessageId
	start, err := h.getSentTime(req.Account, req.MessageId)
	if err != nil {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	var indexes []*rpc.MessageIndex
	tx := h.MessageDb.Model(&database.MessageIndex{}).Select("send_time", "account_b", "direction", "message_id", "group")
	err = tx.Where("account_a=? and send_time>? and direction=?", req.Account, start, 0).Order("send_time asc").Limit(wire.OfflineSyncIndexCount).Find(&indexes).Error
	if err != nil {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	err = setMesssageAck(h.Cache, req.Account, msgId)
	if err != nil {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	_, _ = c.Negotiate(&rpc.GetOfflineMessageIndexResp{
		List: indexes,
	})
}

func (h *ServiceHandler) getSentTime(account string, msgId int64) (int64, error) {
	// 1. 冷启动情况，从服务端拉取消息索引
	if msgId == 0 {
		key := database.KeyMessageAckIndex(account)
		msgId, _ = h.Cache.Get(key).Int64() // 如果一次都没有发ack包，这里就是0
	}
	var start int64
	if msgId > 0 {
		// 2.根据消息ID读取此条消息的发送时间。
		var content database.MessageContent
		err := h.MessageDb.Select("send_time").First(&content, msgId).Error
		if err != nil {
			//3.如果此条消息不存在，返回最近一天
			start = time.Now().AddDate(0, 0, -1).UnixNano()
		} else {
			start = content.SendTime
		}
	}
	// 4.返回默认的离线消息过期时间
	earliestKeepTime := time.Now().AddDate(0, 0, -1*wire.OfflineMessageExpiresIn).UnixNano()
	if start == 0 || start < earliestKeepTime {
		start = earliestKeepTime
	}
	return start, nil
}

func (h *ServiceHandler) GetOfflineMessageContent(c iris.Context) {
	var req rpc.GetOfflineMessageContentReq
	if err := c.ReadBody(&req); err != nil {
		c.StopWithError(iris.StatusBadRequest, err)
		return
	}
	mlen := len(req.MessageIds)
	if mlen > wire.MessageMaxCountPerPage {
		c.StopWithText(iris.StatusBadRequest, "too many MessageIds")
		return
	}
	var contents []*rpc.Message
	err := h.MessageDb.Model(&database.MessageContent{}).Where(req.MessageIds).Find(&contents).Error
	if err != nil {
		c.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	_, _ = c.Negotiate(&rpc.GetOfflineMessageContentResp{
		List: contents,
	})
}
