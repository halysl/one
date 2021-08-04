package routers

import (
	"encoding/json"
	"fmt"
	"github.com/halysl/one/db"
	. "github.com/halysl/one/log"
	"strings"
)

var OrphanBlockAPI = "https://api.filscout.com/api/v1/lonelyblock"

type OrphanBlockInfo struct {
	Code      int    `json:"code"`
	Total     int    `json:"total"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
	Message   string `json:"Message"`
	Data      []struct {
		Height  int    `json:"height"`
		Time    string `json:"time"`
		TimeAgo int    `json:"timeAgo"`
		Blocks  []struct {
			Height     int    `json:"height"`
			Cid        string `json:"cid"`
			MineTime   string `json:"mineTime"`
			Miner      string `json:"miner"`
			MinerTag   string `json:"minerTag"`
			IsVerified int    `json:"isVerified"`
		} `json:"blocks"`
	} `json:"data"`
}

type OrphanBlockAPIParm struct {
	Miner     string `json:"miner"`
	Type      int    `json:"type"`
	Height    int    `json:"height"`
	PageIndex int    `json:"pageindex"`
	PageSize  int    `json:"pagesize"`
}

func orphanBlockView(pageIndex, PageSize int) []OrphanBlockInfo {
	res := []OrphanBlockInfo{}
	database, err := db.GetDBCursor()
	if err != nil {
		ErrorLogger.Println(err)
	}
	defer database.Close()
	data, err := database.Get([]byte("miner"), nil)
	if err != nil {
		ErrorLogger.Println(err)
	}
	var minerList []db.Miner
	err = json.Unmarshal(data, &minerList)
	for _, val := range minerList {
		res = append(res, getOrphanBlock(val.MinerID,pageIndex,PageSize))
	}
	return res
}

func getOrphanBlock(minerID string, pageIndex, PageSize int) (obi OrphanBlockInfo) {
	postBody := &OrphanBlockAPIParm{
		Miner:     minerID,
		Type:      0,
		Height:    0,
		PageIndex: pageIndex,
		PageSize:  PageSize,
	}
	b, _ := json.Marshal(postBody)
	data, err := requestCommon("POST", OrphanBlockAPI, strings.NewReader(string(b)))
	if err != nil {
		ErrorLogger.Println(err)
	}
	if err := json.Unmarshal(data, &obi); err != nil {
		ErrorLogger.Println(err)
	}
	return
}

func humanReadOrphanBlock(obiList []OrphanBlockInfo) string {
	type orphanBlock struct {
		minerID string
		mineTime string
		height int
	}
	var obList []orphanBlock

	for _, obi := range obiList {
		for _, block := range obi.Data {
			for _, b := range block.Blocks {
				obList = append(obList, orphanBlock{
					minerID:  b.Miner,
					mineTime: b.MineTime,
					height:   b.Height,
				})
			}
		}
	}

	database, err := db.GetDBCursor()
	if err != nil {
		ErrorLogger.Println(err)
	}
	defer database.Close()
	data, err := database.Get([]byte("miner"), nil)
	if err != nil {
		ErrorLogger.Println(err)
	}
	var minerList []db.Miner
	err = json.Unmarshal(data, &minerList)

	tileMiner := make(map[string]db.Miner)
	for _, miner := range minerList {
		tileMiner[miner.MinerID] = miner
	}

	minerTitle := "ID：%s 代号：%s"
	orphanBlockMsg := "  - 时间：%s，高度：%d"
    var msgList []string
	minerFlag := true
	defaultMinerID := "f0000"
	for _, ob := range obList {
		if ob.minerID !=defaultMinerID {
			minerFlag = true
			defaultMinerID = ob.minerID
		} else {
			minerFlag = false
		}
		if minerFlag {
			msgList = append(msgList, fmt.Sprintf(minerTitle, ob.minerID, tileMiner[ob.minerID].Cname))
		}
		msgList = append(msgList, fmt.Sprintf(orphanBlockMsg, ob.mineTime, ob.height))
	}
	return strings.Join(msgList, "\n")
}
