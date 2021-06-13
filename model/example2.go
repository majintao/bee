// @BeeOverwrite YES
// @BeeGenerateTime 20210613_192934
package model

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"reflect"
	"strings"
)

type Example2 struct {
	Row1 string     ` orm:"string" json:"row1" form:"row1"` // 字段一
	Row2 string     ` orm:"string" json:"row2" form:"row2"` // 字段二
	Row3 *time.Time ` orm:"Date" json:"row3" form:"row3"`   // 字段三
	Row4 *time.Time ` orm:"Date" json:"row4" form:"row4"`   // 字段四

}

func (t *Example2) TableName() string {
	return "example2"
}
