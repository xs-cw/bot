package dingding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	botSecret    = ""
	defaultReply = "nothing to do"
	dingBot      *DingRobot
)

type DingRobotMsgHeader struct {
	TimeStamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}

func RobotInit(token, robotSecret string) {
	botSecret = robotSecret
	dingBot = NewRobot(token, robotSecret)
}

type DingRobotMessage struct {
	AtUsers []struct {
		DingTalkID string `json:"dingtalkId"`
		StaffID    string `json:"staffId"`
	} `json:"atUsers"`
	ChatBotUserID     string `json:"chatbotUserId"`
	ConversationID    string `json:"conversationId"`
	ConversationTitle string `json:"conversationTitle"`
	ConversationType  string `json:"conversationType"`
	CreateAt          int    `json:"createAt"`
	MsgID             string `json:"msgId"`
	MsgType           string `json:"msgtype"`
	SenderCorpID      string `json:"senderCorpId"`
	SenderID          string `json:"senderId"`
	SenderNick        string `json:"senderNick"`
	SenderStaffID     string `json:"senderStaffId"`
	Text              struct {
		Content string `json:"content"`
	} `json:"text"`
}

func Handle(timestamp string, sign string, m *DingRobotMessage) error {
	if !verifySign(&DingRobotMsgHeader{TimeStamp: timestamp, Sign: sign}) {
		return errors.New("sign not match")
	}
	if m == nil {
		return errors.New("invalid msg")
	}

	content := strings.TrimSpace(m.Text.Content)
	// TODO 命令动态配置
	switch content {
	case "":
	default:
		err := replyText(defaultReply)
		if err != nil {
			log.Println(err.Error())
		}
	}
	return nil
}

var (
	replyFuncs = make(map[string]func(msg string, msgType string) error)
)

func replyText(msg string) error {
	if len(msg) == 0 {
		return errors.New("no message to reply")
	}
	return dingBot.SendMessage(NewMessageBuilder(TypeText).Text(msg).Build())
}

func verifySign(m *DingRobotMsgHeader) bool {
	t, err := strconv.Atoi(m.TimeStamp)
	if err != nil {
		return false
	}
	limit := time.Now().Add(-time.Hour).Unix() * 1000
	if t < int(limit) {
		return false
	}
	message := m.TimeStamp + "\n" + botSecret
	h := hmac.New(sha256.New, []byte(botSecret))
	_, _ = h.Write([]byte(message))
	s := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return m.Sign == s
}
