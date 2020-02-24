package appdata

import (
	"encoding/json"
	"fmt"
	"github.com/xjianfeng/gocomm/logger"
	"sync"
)

var (
	xlsAwardFile = "data.json"
	sheetData    = make(map[string]map[string]map[string]float64)
	xlsAwardData = AwardData{
		data: make(map[string]map[string]float64),
	}
)

type AwardData struct {
	sync.RWMutex
	data map[string]map[string]float64
}

func GetAwardValue(cnt int, awardIdx int) (float64, bool) {
	xlsAwardData.RLock()
	defer xlsAwardData.RUnlock()
	maxIdx := len(xlsAwardData.data)
	i := fmt.Sprintf("%d", cnt)
	if maxIdx <= cnt {
		i = fmt.Sprintf("%d", maxIdx)
	}
	idx := fmt.Sprintf("%d", awardIdx)

	info, ok := xlsAwardData.data[i]
	if !ok {
		return 0, false
	}
	value, ok := info[idx]
	if !ok {
		return 0, false
	}
	return value, true
}

func InitPollConfigByFile() {
	xlsByteData, ok := GetGameData(xlsAwardFile)
	if !ok {
		return
	}
	err := json.Unmarshal(xlsByteData, &sheetData)
	if err != nil {
		logger.LogError("ReadPollByFile Error:%s", err.Error())
	}

	xlsAwardData.Lock()
	defer xlsAwardData.Unlock()

	xlsAwardData.data = sheetData["Sheet1"]
}

func Reload() {
	InitPollConfigByFile()
}

func init() {
	InitPollConfigByFile()
	AddToReload(xlsAwardFile, Reload)
}
