package db

import (
	"encoding/json"
	. "github.com/halysl/one/log"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"os"
)

var dbPath = "db/one.db"

func InitDB() {
	database, err := GetDBCursor()
	if err != nil {
		ErrorLogger.Println(err)
	}
	defer database.Close()

	fileDingDingBot, err := os.Open("db/dingdingbot.json")
	defer fileDingDingBot.Close()
	if err != nil {
		ErrorLogger.Println(err)
	}
	dataDingDingBot, _ := ioutil.ReadAll(fileDingDingBot)

	fileMinerList, err := os.OpenFile("db/miner.json", os.O_RDONLY, 0644)
	defer fileMinerList.Close()
	if err != nil {
		ErrorLogger.Println(err)
	}
	dataMinerList, _ := ioutil.ReadAll(fileMinerList)

	err = database.Put([]byte("dingdingbot"), dataDingDingBot, nil)
	err = database.Put([]byte("miner"), dataMinerList, nil)
}

func addMiner(miner Miner) error {
	database, err := GetDBCursor()
	if err != nil {
		ErrorLogger.Println(err)
	}
	defer database.Close()

	data, err := database.Get([]byte("miner"), nil)
	if err != nil {
		return err
	}
	var minerList []Miner
	err = json.Unmarshal(data, &minerList)
	if err != nil {
		return err
	}
	minerList = append(minerList, miner)
	res, err := json.Marshal(minerList)
	if err != nil {
		return err
	}
	database.Put([]byte("miner"), res, nil)
	return nil
}

func GetDBCursor() (*leveldb.DB, error) {
	return leveldb.OpenFile(dbPath, nil)
}
