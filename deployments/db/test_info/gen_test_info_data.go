package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

const UID = "7791850604"

type typeAndTag struct {
	Type string `json:"type"`
	Tag  string `json:"tag"`
}

var templateInfo = `
db.info.insert({
    "uid": "%s",
    "author": "%s",
    "title": "%s",
    "summary": "%s",
    "country": "japan",
    "contentTime": ISODate("2018-12-21T00:00:00Z"),
    "coverResources": {
        "image1": "/resource/image1.png",
        "image2": "/resource/image2.jpg",
        "image3": "/resource/image3.gif"
    },
    "publishTime": ISODate("2018-12-22T16:16:16.207Z"),
    "lastReviewTime": ISODate("0001-01-01T00:00:00.000Z"),
    "valid": true,
    "watchCount": 0,
    "tags": %s,
    "lastModifyTime": ISODate("2018-12-22T16:16:16.207Z"),
    "canReview": true,
    "thumbUp": 0,
    "thumbDown": 0,
    "favorites": 0,
    "archived": false,
    "latestSegmentID": objectid.ObjectID("000000000000000000000000"),
    "segmentCount": 0
});

`

var templateSeg = `
db.segment.insert({
	"infoID": "%s",
	"title": "%s",
    "no": %d,
	"labels": "%s",
    "content": {
		"image1":"/xxx/xxx/xxx.jpg",
		"image1":"/xxx/xxx/xxx.jpg",
		"image1":"/xxx/xxx/xxx.jpg",
		"image1":"/xxx/xxx/xxx.jpg",
		"image1":"/xxx/xxx/xxx.jpg"
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

var authorList = []string{"新海诚", "宫崎骏", "波多野结衣"}

func getAuthor() string {
	return authorList[rand.Intn(len(authorList))]
}

var summaryList = []string{"根据福令永三的儿童文学作品“蜡笔王国系列”改编而成。银公主是一位长有长而银色的头发，而且笑容十分可爱动人的女孩子。是很多王子心目中的女神：例如饭团王国的古拉度王子，和汉堡饱王国的菲里昂王子就是她的拥趸。", "又有魅力又有品味的你，希望找到一本书，既能给自己带来快乐，又能帮助你了解别人和自己的心灵世界，恭喜你，你的选择绝对正确！"}

func getSummary() string {
	return summaryList[rand.Intn(len(summaryList))]
}

var tagList = []*typeAndTag{&typeAndTag{
	Type: "剧情",
	Tag:  "热血",
},
	&typeAndTag{
		Type: "剧情",
		Tag:  "冒险",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "魔幻",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "神鬼",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "搞笑",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "萌系",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "爱情",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "科幻",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "魔法",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "格斗",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "武侠",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "机战",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "战争",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "竞技",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "体育",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "校园",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "生活",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "励志",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "历史",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "伪娘",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "宅男",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "腐女",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "耽美",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "百合",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "后宫",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "治愈",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "美食",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "推理",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "悬疑",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "恐怖",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "四格",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "职场",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "侦探",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "社会",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "音乐",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "舞蹈",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "杂志",
	},
	&typeAndTag{
		Type: "剧情",
		Tag:  "黑道",
	}, &typeAndTag{
		Type: "受众",
		Tag:  "少女",
	},
	&typeAndTag{
		Type: "受众",
		Tag:  "少年",
	},
	&typeAndTag{
		Type: "受众",
		Tag:  "青年",
	},
	&typeAndTag{
		Type: "受众",
		Tag:  "儿童",
	},
	&typeAndTag{
		Type: "受众",
		Tag:  "通用",
	},
}

func getTags(n int) string {
	tmp := make([]*typeAndTag, 0, n)
	for i := 0; i < n; i++ {
		tmp = append(tmp, tagList[rand.Intn(len(tagList))])
	}
	result, _ := json.Marshal(tmp)
	return string(result)
}

func main() {
	result := "db = db.getSiblingDB('teddy');\n"
	for i := 0; i < 10000; i++ {
		result += fmt.Sprintf(templateInfo, UID, getAuthor(), getTitle(), getSummary(), getTags(rand.Intn(3)+1))
	}
	ioutil.WriteFile("./info_init.js", []byte(result), 0666)
}
