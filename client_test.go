package oppopush

import (
	"fmt"
	"testing"
)

var (
	appKey = ""
	masterSecret = ""

	targetId = ""
)

func TestOppoPush_Broadcast(t *testing.T) {
	var client = NewClient(appKey, masterSecret)
	//保存通知栏消息内容体
	msg0 := NewSaveMessageContent("您好，我是配钥匙的", "您配几把").
		SetSubTitle("您配吗？")

	//根据appmessageId做了消息去重
	// msg0.SetID("12345646546")

	msg0.SetCallBackParameter("jiushihuizhi nengshoudao ?15646854546")
	msg0.SetOffLine(true)
	msg0.SetOffLineTtl(1800)

	msg0.SetActionParameters(`{"key1":"value1","key2":"value2"}`)

	result, err := client.SaveMessageContent(msg0)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println("hahah MessageID:", result.Data)
		//广播推送-通知栏消息
		broadcast := NewBroadcast(result.Data.MessageID).
			SetTargetType(2).
			SetTargetValue(targetId)
		_, res, err := client.Broadcast(broadcast)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("push result data: ", res)
		}
	}
}