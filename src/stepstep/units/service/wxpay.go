package service

import (
	"fmt"
	"stepstep/service"
)

func WxPayTest() {
	ret, err := service.WxPrePay("opYS45Yyx-o-UUxiKG_M_dDLP6Lk", 500)

	fmt.Printf("WxPayTest ret:%v, err:%v", ret, err)
}
