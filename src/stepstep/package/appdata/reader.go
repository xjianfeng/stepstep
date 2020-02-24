package appdata

import (
	"encoding/json"
	log "github.com/xjianfeng/gocomm/logger"
	"io/ioutil"
	"stepstep/conf"
	"sync"
)

type RwDataInfo struct {
	sync.RWMutex
	data map[string][]byte
}

type CacheInfo struct {
	sync.RWMutex
	data map[string]map[string]bool
}

type XlsInfo struct {
	sync.RWMutex
	data map[string]SheetData
}

var (
	xlsCacheKey = CacheInfo{
		data: make(map[string]map[string]bool),
	}
	dataInfo = RwDataInfo{
		data: make(map[string][]byte),
	}
	xlsInfo = XlsInfo{
		data: make(map[string]SheetData),
	}
	reloadFunc = make(map[string]func())
	dataPath   = conf.CfgServer.DataPath + "data/"
)

type RowData map[string]interface{}
type SheetData map[string]RowData
type XlsSheel map[string]SheetData

type DecodeArgs struct {
	FileName  string
	SheetName string
	UseNumber bool
}

func GetDataInfo(key string) ([]byte, bool) {
	dataInfo.RLock()
	defer dataInfo.RUnlock()

	data, ok := dataInfo.data[key]
	return data, ok
}

func SetDataInfo(key string, value []byte) {
	dataInfo.Lock()
	defer dataInfo.Unlock()

	dataInfo.data[key] = value
}

func DelDataInfo(key string) {
	dataInfo.Lock()
	defer dataInfo.Unlock()

	delete(dataInfo.data, key)
}

func GetCacheData(key string) (map[string]bool, bool) {
	xlsCacheKey.RLock()
	defer xlsCacheKey.RUnlock()

	data, ok := xlsCacheKey.data[key]
	return data, ok
}

func SetCacheData(key1, key2 string, value bool) {
	xlsCacheKey.Lock()
	defer xlsCacheKey.Unlock()

	if _, ok := xlsCacheKey.data[key1]; !ok {
		xlsCacheKey.data[key1] = make(map[string]bool)
	}
	xlsCacheKey.data[key1][key2] = value
}

func DelteCache(fileName string) {
	xlsCacheKey.Lock()
	defer xlsCacheKey.Unlock()

	data, ok := xlsCacheKey.data[fileName]
	if !ok {
		return
	}
	DelXlsInfo(data)
}

func DelXlsInfo(kInfo map[string]bool) {
	xlsInfo.Lock()
	defer xlsInfo.Unlock()

	if len(kInfo) == 0 {
		return
	}
	for k, _ := range kInfo {
		delete(xlsInfo.data, k)
	}
}

func GetXlsInfo(key string) (SheetData, bool) {
	xlsInfo.RLock()
	defer xlsInfo.RUnlock()

	data, ok := xlsInfo.data[key]
	return data, ok
}

func SetXlsInfo(key string, value SheetData) {
	xlsInfo.Lock()
	defer xlsInfo.Unlock()

	xlsInfo.data[key] = value
}

func UpDateGameData(fileName string) {
	filePath := dataPath + fileName
	log.LogInfo("UpDateGameData %s", filePath)
	DelDataInfo(filePath)
	DelteCache(fileName)

	GetGameData(fileName)
	if fun, ok := reloadFunc[filePath]; ok {
		fun()
	}
}

func AddToReload(fileName string, fun func()) {
	filePath := dataPath + fileName
	reloadFunc[filePath] = fun
}

func GetGameData(fileName string) ([]byte, bool) {
	filePath := dataPath + fileName
	data, ok := GetDataInfo(filePath)
	if ok {
		return data, false
	}

	readData, err := ioutil.ReadFile(filePath)
	log.LogInfo("GetGameData Read File %s", filePath)
	if err != nil {
		log.LogError("Read gamedata LogError %v", err.Error())
		return nil, false
	}

	SetDataInfo(filePath, readData)
	return readData, true
}

func RemoveCache(fileName, sheetName string) {
	cacheKey := fileName + sheetName
	xlsInfo.Lock()
	defer xlsInfo.Unlock()

	delete(xlsInfo.data, cacheKey)
}

//获取xls json
func GetXlsJsonData(fileName string) (SheetData, error) {
	var err error
	var data SheetData

	cacheKey := fileName
	readData, _ := GetGameData(fileName)
	cache, ok := GetXlsInfo(cacheKey)
	if ok {
		return cache, nil
	}
	err = json.Unmarshal(readData, &data)
	if err != nil {
		log.LogError("GetXlsData fileName %s, LogError %s", fileName, err.Error())
		return nil, err
	}
	SetXlsInfo(cacheKey, data)
	SetCacheData(fileName, cacheKey, true)
	return data, nil
}

func GetXlsData(fileName, sheetName string) (SheetData, error) {
	var err error
	var data XlsSheel

	readData, _ := GetGameData(fileName)
	cacheKey := fileName + sheetName
	cache, ok := GetXlsInfo(cacheKey)
	if ok {
		return cache, nil
	}
	err = json.Unmarshal(readData, &data)

	if err != nil {
		log.LogError("GetXlsData fileName %s, sheetName %s, LogError %s", fileName, sheetName, err.Error())
		return nil, err
	}
	sheetData := data[sheetName]
	SetXlsInfo(cacheKey, sheetData)
	SetCacheData(fileName, cacheKey, true)
	return sheetData, nil
}
