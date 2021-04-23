package oppopush

import (
	"errors"
	"fmt"
	"io/ioutil"
	"testing"
	"net/http"
	"time"
)


var (
	appKey = ""
	masterSecret = ""

	targetId = ""

	oc = NewClient(appKey, masterSecret)
)

func TestIconPushBuffer(t *testing.T) {
	buf,err := httpGet("https://i.ibb.co/tPcQnJr/picx.png")
	if err != nil {
		t.Log(err)
		return;
	}

	res,err := oc.UploadIcon("pickkfwe.png", buf, 86400)
	if err != nil {
		t.Log(err)
		return;
	}
	
	OppoPushIcon("CN_8995c8e604617e01e0c2a845df06d2d3", res.Data.SmallPicId)
}




func httpGet(durl string) ([]byte, error) {

	client := http.DefaultClient
	client.Timeout = time.Second * 60
	resp, err := client.Get(durl)

	if err != nil {
		return nil, err
	}
	if resp.ContentLength <= 0 {
		return nil, errors.New("http Get rsp length <=0")
	}

	raw := resp.Body
	defer raw.Close()

	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return rspBody, nil
}



func TestIconPush(t *testing.T) {
	res,err := oc.UploadIcon("F:\\pic3_120_120.png", nil, 86400)
	if err != nil {
		t.Log(err)
		return;
	}
	
	OppoPushIcon("CN_8995c8e604617e01e0c2a845df06d2d3", res.Data.SmallPicId)
}


func OppoPushIcon(targetId, picId string) {
	//保存通知栏消息内容体
	pp := "67878"
	msg := NewSaveMessageContent(pp+"4532453453", pp+"hrehertfgd33")
	msg.SetID(pp)

	msg.SetOffLine(true)
	msg.SetOffLineTtl(86400)
	msg.SetChannelId("seewo")
	msg.SetSmallPictureId(picId)

	result, err := oc.SaveMessageContent(msg)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Printf("save result:%+v\n", result)
		//广播推送-通知栏消息
		broadcast := NewBroadcast(result.Data.MessageID).
			SetTargetType(2).
			SetTargetValue(targetId)
		_, res, err := oc.Broadcast(broadcast)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("push result data: ", res)
	}
}







func TestOppoPush_Broadcast(t *testing.T) {
	//保存通知栏消息内容体
	msg0 := NewSaveMessageContent("您好，我是配钥匙的", "您配几把").
		SetSubTitle("您配吗？")

	//根据appmessageId做了消息去重
	// msg0.SetID("12345646546")

	msg0.SetCallBackParameter("jiushihuizhi nengshoudao ?15646854546")
	msg0.SetOffLine(true)
	msg0.SetOffLineTtl(1800)

	msg0.SetActionParameters(`{"key1":"value1","key2":"value2"}`)

	result, err := oc.SaveMessageContent(msg0)
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println("hahah MessageID:", result.Data)
		//广播推送-通知栏消息
		broadcast := NewBroadcast(result.Data.MessageID).
			SetTargetType(2).
			SetTargetValue(targetId)
		_, res, err := oc.Broadcast(broadcast)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("push result data: ", res)
		}
	}
}


