package routers

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	. "github.com/halysl/one/log"
	"time"
)

var title = []string{
	"矿工id", "节点名称", "扇区", "链接",
	"节点算力", "算力增长24h", "算力增长7d", "算力增长30d",
	"节点出块", "节点win块", "出块24h", "出块7d", "出块30d", "幸运值24h", "幸运值7d", "幸运值30d", "理论出块24h", "理论幸运值24h",
	"节点扇区", "有效扇区", "错误扇区", "恢复扇区", "未windpost扇区", "未windpost算力",
}

func renderReport(minerList []MinerResultInfo, overview FilscoutOverviewInfo) *excelize.File {
	InfoLogger.Printf("%+v", overview)
	baseDict := make(map[string]string)
	for index, value := range title {
		baseDict[value] = string(index+65)
	}

	table := createTable()
	sheet := "Sheet1"
	baseYAxis := 3
	for index, miner := range minerList {
		table.SetCellStr(sheet, getAxis("A", index+baseYAxis), miner.MinerID)
		table.SetCellStr(sheet, getAxis("B", index+baseYAxis), miner.CName)
		table.SetCellInt(sheet, getAxis("C", index+baseYAxis), miner.SectorSize)
		table.SetCellHyperLink(sheet, getAxis("D", index+baseYAxis), miner.MinerURL, "External")
		table.SetCellStr(sheet, getAxis("E", index+baseYAxis), HumanBytesLoaded(miner.Power.TotalPower))
		table.SetCellStr(sheet, getAxis("F", index+baseYAxis), HumanBytesLoaded(miner.Power.PowerGrowth24h))
		table.SetCellStr(sheet, getAxis("G", index+baseYAxis), HumanBytesLoaded(miner.Power.PowerGrowth7d))
		table.SetCellStr(sheet, getAxis("H", index+baseYAxis), HumanBytesLoaded(miner.Power.PowerGrowth30d))
		table.SetCellInt(sheet, getAxis("I", index+baseYAxis), miner.BlockCount.TotalBlockCount)
		table.SetCellInt(sheet, getAxis("J", index+baseYAxis), miner.BlockCount.TotalWinBlockCount)
		table.SetCellInt(sheet, getAxis("K", index+baseYAxis), miner.BlockCount.BlockCount24h)
		table.SetCellInt(sheet, getAxis("L", index+baseYAxis), miner.BlockCount.BlockCount7d)
		table.SetCellInt(sheet, getAxis("M", index+baseYAxis), miner.BlockCount.BlockCount30d)
		table.SetCellFloat(sheet, getAxis("N", index+baseYAxis), miner.BlockCount.Lucky24h, -1, 64)
		table.SetCellFloat(sheet, getAxis("O", index+baseYAxis), miner.BlockCount.Lucky7d, -1, 64)
		table.SetCellFloat(sheet, getAxis("P", index+baseYAxis), miner.BlockCount.Lucky30d, -1, 64)
		theoryBlock24h := overview.Data.BlockRewardIn24h * (float64(miner.Power.TotalPower)/float64(1 << 40)) / overview.Data.AvgBlocksReword
		theoryLucky24h := float64(miner.BlockCount.BlockCount24h)/theoryBlock24h
		table.SetCellFloat(sheet, getAxis("Q", index+baseYAxis), theoryBlock24h, -1, 64)
		table.SetCellFloat(sheet, getAxis("R", index+baseYAxis), theoryLucky24h, -1, 64)
		table.SetCellInt(sheet, getAxis("S", index+baseYAxis), miner.SectorCount.TotalSectorCount)
		table.SetCellInt(sheet, getAxis("T", index+baseYAxis), miner.SectorCount.ActiveCount)
		table.SetCellInt(sheet, getAxis("U", index+baseYAxis), miner.SectorCount.FaultCount)
		table.SetCellInt(sheet, getAxis("V", index+baseYAxis), miner.SectorCount.RecoveryCount)
		unwindpostSector := miner.SectorCount.TotalSectorCount - miner.SectorCount.ActiveCount
		unwindpostPower := unwindpostSector * miner.SectorSize
		table.SetCellInt(sheet, getAxis("W", index+baseYAxis), unwindpostSector)
		table.SetCellStr(sheet, getAxis("X", index+baseYAxis), HumanBytesLoaded(int64(unwindpostPower)*1024*1024*1024))
	}
	// set style
	luckyStyle, err := table.NewStyle(`{"number_format": 10, "decimal_places": 2}`)
	if err != nil {
		ErrorLogger.Println(err)
	}
	table.SetCellStyle(sheet, fmt.Sprintf("N%d", baseYAxis), fmt.Sprintf("P%d", baseYAxis+len(minerList)), luckyStyle)
	table.SetCellStyle(sheet, fmt.Sprintf("R%d", baseYAxis), fmt.Sprintf("R%d", baseYAxis+len(minerList)), luckyStyle)

	defaultFloatStyle, err := table.NewStyle(`{"number_format": 2, "decimal_places": 2}`)
	if err != nil {
		ErrorLogger.Println(err)
	}
	table.SetCellStyle(sheet, fmt.Sprintf("Q%d", baseYAxis), fmt.Sprintf("Q%d", baseYAxis+len(minerList)), defaultFloatStyle)
	// set column width
	table.SetColWidth(sheet, "A", "X", 10)
	table.SetColWidth(sheet, "C", "C", 5)
	return table
}

func getAxis(x string,y int) string {
	return fmt.Sprintf("%s%d", x, y)
}

func HumanBytesLoaded(num int64) string {
	suffix := ""
	b := float64(num)
	if num > (1 << 60) {
		suffix = "EiB"
		b = float64(num) / float64(1 << 60)
	} else if num > (1 << 50) {
		suffix = "PiB"
		b = float64(num) / float64(1 << 50)
	} else if num > (1 << 40) {
		suffix = "TiB"
		b = float64(num) / float64(1 << 40)
	} else if num > (1 << 30) {
		suffix = "GiB"
		b = float64(num) / float64(1 << 30)
	} else if num > (1 << 20) {
		suffix = "MiB"
		b = float64(num) / float64(1 << 20)
	} else if num > (1 << 10) {
		suffix = "KiB"
		b = float64(num) / float64(1 << 10)
	} else {
		suffix = "B"
		b = float64(num) / float64(1 << 0)
	}

	return fmt.Sprintf("%.2f%s", b, suffix)
}

func createTable() *excelize.File {
	f := excelize.NewFile()
	index := f.GetSheetIndex("Sheet1")
	f.SetCellStr("Sheet1", getAxis("A", 1), time.Now().Format(time.RFC3339))
	f.SetActiveSheet(index)
	for index, value := range title {
		if err := f.SetCellStr("Sheet1", fmt.Sprintf("%s%d", string(index+65), 2), value); err != nil {
			ErrorLogger.Println(err)
		}
	}
	return f
}
