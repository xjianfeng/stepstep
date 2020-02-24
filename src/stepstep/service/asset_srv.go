package service

import (
	"fmt"
	"stepstep/models"
)

type ReqAssetDetail struct {
	PageIndex int
	PageSize  int
}

type RespAssetDetail struct {
	NextPage bool          `json:"nextPage"`
	List     []interface{} `json:"list"`
}

func GetAssetDetail(unionId string, args *ReqAssetDetail) *RespAssetDetail {
	result := &RespAssetDetail{}
	assetList := models.GetAssetLog(unionId, args.PageIndex*args.PageSize, args.PageSize)

	result.NextPage = assetList.RecordNum > args.PageIndex*args.PageSize+args.PageSize
	result.List = []interface{}{}

	var total float64
	for _, info := range assetList.AssetLog {
		total = 0
		for _, data := range info.List {
			total += data.Value
		}
		tmpData := map[string]interface{}{
			"date":    info.Date,
			"list":    info.List,
			"summary": fmt.Sprintf("%.2f", total),
		}
		result.List = append(result.List, tmpData)
	}
	return result
}
