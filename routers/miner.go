package routers

import (
	"encoding/json"
	"fmt"
	"github.com/halysl/one/db"
	. "github.com/halysl/one/log"
	"strings"
	"math/big"
)

var (
	filscoutOverview   = "https://api.filscout.com/api/v1/network/overview/info"
	filscoutMinerView  = "https://www.filscout.com/zh/miner/%s"
	filscoutMinerInfo  = "https://api.filscout.com/api/v1/miners/%s/info"
	filscoutMinerStats = "https://api.filscout.com/api/v1/miners/%s/miningstats"
)

type FilscoutOverviewInfo struct {
	Data struct {
		Height             int     `json:"tipSetHeight"`
		TotalPower         float64 `json:"totalPower"`
		ActiveMiners       int     `json:"activeMiners"`
		AvgBlocksReword    float64 `json:"avgBlocksReword"`
		BlockRewardIn24h float64 `json:"blockRewardIn24h"`
		OutPutFil24H       float64 `json:"newlyFilIn24h"`
		MessageCount24H    int     `json:"oneDayMessages"`
		CostPerTB          float64 `json:"currentPledgeCollateralTB"`
	} `json:"data"`
}

type MinerInfo struct {
	Data struct {
		Miner        string `json:"miner"`
		QualityPower int64  `json:"qualityPower"`
		RawPower     int64  `json:"rawPower"`
		Blocks       int    `json:"blocks"`
		WinCount     int    `json:"winCount"`
		BlockReward  big.Int  `json:"blockReward"`
		PowerRank    int    `json:"powerRank"`
		Sectors      struct {
			SectorSize    int64 `json:"sectorSize"`
			SectorCount   int   `json:"sectorCount"`
			ActiveCount   int   `json:"activeCount"`
			FaultCount    int   `json:"faultCount"`
			RecoveryCount int   `json:"recoveryCount"`
		} `json:"sectors"`
		Balance struct {
			Balance       big.Int `json:"balance"`
			Available     int64 `json:"available"`
			SectorsPledge int64 `json:"sectorsPledge"`
			LockedFunds   int64 `json:"lockedFunds"`
		} `json:"balance"`
	} `json:"data"`
}

type MinerStats struct {
	Data struct {
		Miner                 string  `json:"miner"`
		QualityPowerGrowth    int64   `json:"qualityPowerGrowth"`
		ProvingPower          int     `json:"provingPower"`
		MiningEfficiencyFloat float64 `json:"miningEfficiencyFloat"`
		MiningEfficiency      string  `json:"miningEfficiency"`
		MachinesNum           float64     `json:"machinesNum"`
		Blocks                int     `json:"blocks"`
		BlockReward           big.Int   `json:"blockReward"`
		LuckyValueFloat       float64 `json:"luckyValueFloat"`
		StatsType             string  `json:"statsType"`
	} `json:"data"`
}

type MinerResultInfo struct {
	MinerID    string
	CName      string
	SectorSize int
	MinerURL   string
	BlockCount struct {
		BlockCount24h      int
		BlockCount7d       int
		BlockCount30d      int
		Lucky24h           float64
		Lucky7d            float64
		Lucky30d           float64
		TotalBlockCount    int
		TotalWinBlockCount int
	}
	SectorCount struct {
		TotalSectorCount  int
		ActiveCount       int
		FaultCount        int
		RecoveryCount     int
		UnPostSectorCount int
	}
	Power struct {
		PowerRank      int
		TotalPower     int64
		PowerGrowth24h int64
		PowerGrowth7d  int64
		PowerGrowth30d int64
	}
	Msg string
}

func GetTotalNetInfo() (FilscoutOverviewInfo, error) {
	var filOverview FilscoutOverviewInfo
	method := "GET"
	data, err := requestCommon(method, filscoutOverview, nil)
	if err != nil {
		ErrorLogger.Println(err)
	}
	if err := json.Unmarshal(data, &filOverview); err != nil {
		ErrorLogger.Println(err)
	}
	return filOverview, nil
}

func requestFilscoutByAddress(miner db.Miner) (MinerResultInfo, error) {
	var (
		info     MinerInfo
		stats24h MinerStats
		stats7d  MinerStats
		stats30d MinerStats
	)
	mr := MinerResultInfo{MinerID: miner.MinerID, CName: miner.Cname, SectorSize: miner.Size}

	infoURL := fmt.Sprintf(filscoutMinerInfo, miner.MinerID)
	statsURL := fmt.Sprintf(filscoutMinerStats, miner.MinerID)

	method := "POST"
	data, _ := requestCommon(method, infoURL, nil)
	err := json.Unmarshal(data, &info)
	if err != nil {
		ErrorLogger.Println(err)
	}
	mr.MinerURL = fmt.Sprintf(filscoutMinerView, mr.MinerID)
	mr.SectorCount.TotalSectorCount = info.Data.Sectors.SectorCount
	mr.SectorCount.ActiveCount = info.Data.Sectors.ActiveCount
	mr.SectorCount.FaultCount = info.Data.Sectors.FaultCount
	mr.SectorCount.RecoveryCount = info.Data.Sectors.RecoveryCount
	mr.SectorCount.UnPostSectorCount = mr.SectorCount.TotalSectorCount - mr.SectorCount.ActiveCount
	mr.BlockCount.TotalBlockCount = info.Data.Blocks
	mr.BlockCount.TotalWinBlockCount = info.Data.WinCount
	mr.Power.TotalPower = info.Data.QualityPower
	mr.Power.PowerRank = info.Data.PowerRank

	data, _ = requestCommon(method, statsURL, strings.NewReader(`{"statsType": "24h"}`))
	err = json.Unmarshal(data, &stats24h)
	if err != nil {
		ErrorLogger.Println(err)
	}
	data, _ = requestCommon(method, statsURL, strings.NewReader(`{"statsType": "7d"}`))
	err = json.Unmarshal(data, &stats7d)
	if err != nil {
		ErrorLogger.Println(err)
	}
	data, _ = requestCommon(method, statsURL, strings.NewReader(`{"statsType": "30d"}`))
	err = json.Unmarshal(data, &stats30d)
	if err != nil {
		ErrorLogger.Println(err)
	}
	mr.BlockCount.BlockCount24h = stats24h.Data.Blocks
	mr.BlockCount.BlockCount7d = stats7d.Data.Blocks
	mr.BlockCount.BlockCount30d = stats30d.Data.Blocks
	mr.BlockCount.Lucky24h = stats24h.Data.LuckyValueFloat
	mr.BlockCount.Lucky7d = stats7d.Data.LuckyValueFloat
	mr.BlockCount.Lucky30d = stats30d.Data.LuckyValueFloat
	mr.Power.PowerGrowth24h = stats24h.Data.QualityPowerGrowth
	mr.Power.PowerGrowth7d = stats7d.Data.QualityPowerGrowth
	mr.Power.PowerGrowth30d = stats30d.Data.QualityPowerGrowth
	return mr, nil
}

func GetfilScountData() []MinerResultInfo {
	var minerList []db.Miner
	var result []MinerResultInfo
	database, err := db.GetDBCursor()
	if err != nil {
		ErrorLogger.Println(err)
	}
	defer database.Close()

	data, err := database.Get([]byte("miner"), nil)
	if err != nil {
		ErrorLogger.Println(err)
	}
	err = json.Unmarshal(data, &minerList)
	if err != nil {
		ErrorLogger.Println(err)
	}
	for _, v := range minerList {
		mr, err := requestFilscoutByAddress(v)
		if err != nil {
			ErrorLogger.Println(err)
		}
		result = append(result, mr)
	}
	return result
}
