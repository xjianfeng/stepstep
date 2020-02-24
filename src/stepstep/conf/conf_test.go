package conf

import (
	"testing"
)

func TestInitData(t *testing.T) {
	t.Logf("%+v, %+v, %+v", CfgSever, CfgDb, CfgRedis)
}
