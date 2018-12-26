package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

const INFOID = "5c1f824eb034c620bf7b0d15"

var templateSeg = `
db.segment.insert({
	"infoID": ObjectId("%s"),
	"title": "%s",
    "no": %d,
	"labels": %s,
    "content": {
		"image2":"/5c1f824eb034c620bf7b0d15/5c21cdf0c1fb2bb20d6a8a7c/88e93d1f8f3aa55c52384b7a5c7fb7cd.jpg",
		"image1":"/5c1f824eb034c620bf7b0d15/5c21cdf0c1fb2bb20d6a8a7c/c7135d5883d6e1104d2bc075e4947339.jpg"
	},
    watchCount: 0
});

`

var titleList = []string{
	"看漫画",
	"阿达姆松",
	"爱漫画",
	"恭喜发财",
	"梦幻蜡笔王国",
	"哆啦A梦映画版",
	"RPG哆啦A梦游戏书",
	"老夫子魔界梦战记",
	"千奇百怪",
	"变身娃娃",
	"济公Q传",
	"男女生对对碰",
	"暴走派对",
	"哆啦A梦深入导览",
	"北方酱的日常",
	"西游记",
	"终极米迷口袋书",
	"壹周漫画",
	"彩色世界童话全集",
	"超级漫画素描技法",
	"画书大王",
	"整人大夫",
	"爆笑王国",
	"摩登蕃仔",
	"幽默大师",
	"亚洲黄龙传奇",
	"情趣花生",
	"宇宙旗袍",
	"开喜阿婆",
	"淘漫画",
	"漫王",
	"回到明朝当王爷",
	"艳势番",
	"少年P",
	"最漫画",
	"步步惊心",
	"我御齐天",
	"校园宠物阿汤猫",
	"天漫",
	"倒数5秒",
	"中国漫画",
	"漫友",
	"漫画王",
}

func getTitle() string {
	return titleList[rand.Intn(len(titleList))]
}

func main() {
	result := "db = db.getSiblingDB('teddy');\n"
	for i := 0; i < 1500; i++ {
		result += fmt.Sprintf(templateSeg, INFOID, getTitle(), i, "[\"normal\"]")
	}
	ioutil.WriteFile("./info_segment_init.js", []byte(result), 0666)
}
