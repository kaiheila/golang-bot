package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/bytedance/sonic"
	event2 "github.com/kaiheila/golang-bot/api/base/event"
	"github.com/kaiheila/golang-bot/api/helper"
	"github.com/kaiheila/golang-bot/example/conf"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"testing"
)

func TestSendChallenge(t *testing.T) {
	signal := event2.NewChallengeEventSignal("xxxx", "yyyy")
	data, err := sonic.Marshal(signal)
	var buf bytes.Buffer
	g := zlib.NewWriter(&buf)
	if _, err = g.Write(data); err != nil {
		log.Error(err)
		return
	}
	g.Close()
	client := http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080", &buf)
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Content-Encoding", "compress/zlib")
	resp, err := client.Do(req)
	respData, err := io.ReadAll(resp.Body)

	t.Log("respData", string(respData))
}

func TestMessageUpdate(t *testing.T) {

	for i := 0; i < 250; i++ {
		bodyData, _ := sonic.Marshal(map[string]interface{}{
			"msg_id": "8090a08b-0869-4e08-9fdb-36323520c9e6",
			//"content": "[\n  {\n    \"type\": \"card\",\n    \"theme\": \"secondary\",\n    \"size\": \"lg\",\n    \"modules\": [\n      {\n        \"type\": \"section\",\n        \"text\": {\n          \"type\": \"plain-text\",\n          \"content\": \"KOOK：专属游戏玩家的文字、语音与组队工具。\"\n        }\n      },\n      {\n        \"type\": \"section\",\n        \"text\": {\n          \"type\": \"kmarkdown\",\n          \"content\": \"(font)安全免费(font)[success] (font)没有广告(font)[purple] (font)低资源占用(font)[warning] (font)高通话质量2(font)[pink]\"\n        }\n      }\n    ]\n  }\n]",
			"content": fmt.Sprintf("fff%d", i),
			//"quote":   "",
		})
		//quote:
		resp, err := helper.NewApiHelper("/v3/message/update", conf.Token, conf.BaseUrl, "", "").
			SetBody(bodyData).Post()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(string(resp))
	}

}

func TestMessageView(t *testing.T) {

	resp, err := helper.NewApiHelper("/v3/message/view", conf.Token, conf.BaseUrl, "", "").
		SetQuery(map[string]string{"msg_id": "e521e4ae-535c-4dfc-b833-72b4bfc3e7d4"}).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestMesssageCreate(t *testing.T) {
	for i := 0; i < 2; i++ {
		bodyData, _ := sonic.Marshal(map[string]interface{}{
			"target_id": "6762172889710477",
			"content":   "真棒",
			//"temp_target_id": "3086173823",
			// "quote":     "314c3a8c-0c01-489b-b2d2-2711b36a84c0",
		})
		//"quote":     "d4ab8788-2ba8-4a1f-966d-38fa20f178f0"
		resp, err := helper.NewApiHelper("/v3/message/create", conf.Token, conf.BaseUrl, "", "").
			SetBody(bodyData).Post()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(string(resp))
	}
}

func TestCardMessageCreate(t *testing.T) { //nolint:paralleltest
	bodyData, _ := sonic.Marshal(map[string]interface{}{
		"target_id": "1095267793744046",
		"content":   "[\n  {\n    \"type\": \"card\",\n    \"theme\": \"secondary\",\n    \"size\": \"lz\",\n    \"modules\": [\n      {\n        \"type\": \"section\",\n        \"text\": {\n          \"type\": \"plain-text\",\n          \"content\": \"KOOK：专属游戏玩家的文字、语音与组队工具。\"\n        }\n      },\n      {\n        \"type\": \"section\",\n        \"text\": {\n          \"type\": \"kmarkdown\",\n          \"content\": \"(font)安全免费(font)[success] (font)没有广告(font)[purple] (font)低资源占用(font)[warning] (font)高通话质量(font)[pink]\"\n        }\n      }\n    ]\n  }\n]",
		"type":      10,

		// "quote":     "314c3a8c-0c01-489b-b2d2-2711b36a84c0",
	})
	//"quote":     "d4ab8788-2ba8-4a1f-966d-38fa20f178f0"
	resp, err := helper.NewApiHelper("/v3/message/create", conf.Token, conf.BaseUrl, "", "").
		SetBody(bodyData).Post()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestChannelCreate(t *testing.T) {
	bodyData, _ := sonic.Marshal(map[string]interface{}{
		"guild_id":    "8933621933970340",
		"parent_id":   "7130450443597365",
		"name":        "文字6",
		"type":        1,
		"is_category": 0,
		//"quote":     "314c3a8c-0c01-489b-b2d2-2711b36a84c0",
	})
	//"quote":     "d4ab8788-2ba8-4a1f-966d-38fa20f178f0"
	resp, err := helper.NewApiHelper("/v3/channel/create", conf.Token, conf.BaseUrl, "", "").
		SetBody(bodyData).Post()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestChannelUpdate(t *testing.T) {
	bodyData, _ := sonic.Marshal(map[string]interface{}{
		"channel_id":    "8418843659211643",
		"voice_quality": "1",
		//"quote":     "314c3a8c-0c01-489b-b2d2-2711b36a84c0",
	})
	//"quote":     "d4ab8788-2ba8-4a1f-966d-38fa20f178f0"
	resp, err := helper.NewApiHelper("/v3/channel/update", conf.Token, conf.BaseUrl, "", "").
		SetBody(bodyData).Post()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestGuildView(t *testing.T) {
	resp, err := helper.NewApiHelper("/v3/guild/view", conf.Token, conf.BaseUrl, "", "").
		SetQuery(map[string]string{"guild_id": "2448831678506159", "channel_view": "1"}).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestChannelRoleView(t *testing.T) {
	resp, err := helper.NewApiHelper("/v3/channel-role/index", conf.Token, conf.BaseUrl, "", "").
		SetQuery(map[string]string{"channel_id": "3969409560197275"}).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestMessageList(t *testing.T) {
	resp, err := helper.NewApiHelper("/v3/message/list", conf.Token, conf.BaseUrl, "", "").
		SetQuery(map[string]string{"target_id": "7253269367245022", "pin": "0", "page_size": "10", "msg_id": "ffb2b8f5-3da8-4052-80cb-20ef7e109119"}).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestDirectMessageList(t *testing.T) {
	resp, err := helper.NewApiHelper("/v3/direct-message/list", conf.Token, conf.BaseUrl, "", "").
		SetQuery(map[string]string{"target_id": "3871634671", "page": "0"}).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestUserChatCreate(t *testing.T) {

}
func TestAssetCreate(t *testing.T) {

	resp, err := helper.NewApiHelper("/v3/asset/create", conf.Token, conf.BaseUrl, "", "").
		SetUploadFile("/Users/mikewei/Downloads/avatar_5.jpg").Post()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(resp))
	}

}

func TestChannelUpdate2(t *testing.T) {
	bodyData, _ := sonic.Marshal(map[string]interface{}{
		"channel_id":    "3301653046580657",
		"voice_quality": "3",
	})
	resp, err := helper.NewApiHelper("/v3/channel/update", conf.Token, conf.BaseUrl, "", "").
		SetBody(bodyData).Post()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestGameActivity(t *testing.T) {
	bodyData, _ := sonic.Marshal(map[string]interface{}{
		"id":        2,
		"data_type": 1,
	})
	resp, err := helper.NewApiHelper("/v3/game/activity", conf.Token, conf.BaseUrl, "", "").
		SetBody(bodyData).Post()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestGetGameActivity(t *testing.T) {

	resp, err := helper.NewApiHelper("/v3/game", conf.Token, conf.BaseUrl, "", "").Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestUserChatList(t *testing.T) {
	resp, err := helper.NewApiHelper("/v3/user-chat/list", conf.Token, conf.BaseUrl, "", "").
		SetQuery(map[string]string{"page_size": "50", "page": "1"}).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestGatewayVoice(t *testing.T) {
	bodyData, _ := sonic.Marshal(map[string]interface{}{
		"channel_id": 7617469941537816,
	})
	resp, err := helper.NewApiHelper("/v3/gateway/voice", conf.Token, conf.BaseUrl, "", "").SetBody(bodyData).Post()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestUserMe(t *testing.T) {
	resp, err := helper.NewApiHelper("/v3/user/me", conf.Token, conf.BaseUrl, "", "").Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestUserView(t *testing.T) {
	resp, err := helper.NewApiHelper("/v3/user/view", conf.Token, conf.BaseUrl, "", "").SetQuery(map[string]string{"user_id": "3917927897", "guild_id": "6961481406962448"}).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestGetChannelList(t *testing.T) {
	resp, err := helper.NewApiHelper("/v3/channel/list", conf.Token, conf.BaseUrl, "", "").SetQuery(map[string]string{"guild_id": "6961481406962448", "parent_id": "1"}).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}

func TestBlackListCreate(t *testing.T) {
	bodyData, _ := sonic.Marshal(map[string]interface{}{
		"guild_id":  "6961481406962448",
		"target_id": "406324848",
		"remark":    "test",
	})
	resp, err := helper.NewApiHelper("/v3/blacklist/create", conf.Token, conf.BaseUrl, "", "").SetBody(bodyData).Post()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(resp))
}
