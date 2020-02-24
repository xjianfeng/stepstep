package service

import (
	"encoding/json"
	"fmt"
	"stepstep/define"
	"stepstep/service"
)

func TestAssetDetail() {
	u, _ := service.GetUser("opYS45Yyx-o-UUxiKG_M_dDLP6Lk")
	u.AddMoney(1.2, define.SYSTEM_TASK, define.DESC_TYPE_TASK)
	ret := service.GetAssetDetail("opYS45Yyx-o-UUxiKG_M_dDLP6Lk", &service.ReqAssetDetail{PageIndex: 0, PageSize: 10})
	strdata, err := json.Marshal(ret)
	fmt.Printf("%v, %v", string(strdata), err)
}
