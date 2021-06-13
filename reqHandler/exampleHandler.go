package reqHandler

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"nonolive/model/cms"
	nonomongo "nonolive/nonoutils/nonomongo_v2"
	. "nonolive/servers/cms-api-server/utils/logs"
	. "nonolive/servers/cms-api-server/websrv"
	"time"
)

// 分页列表查询参数
type ExampleParams struct {
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Row1  string         `json:"row1"` // 字段一
	Row2  string         `json:"row2"` // 字段二
	Row3  *time.Time     `json:"row3"` // 字段三
	Row4  *time.Time     `json:"row4"` // 字段四
	Row5  []cms.Example2 `json:"row5"` // 对象五

}

type ExampleListRespBody struct {
	Models    interface{} `json:"models"`
	TotalRows int         `json:"total_rows"`
}
type ExampleParams struct {
	cms.Example
	Id string `json:"_id"`
}

func GetExampleList(ctx *RequestResponseContext, r *http.Request) {
	rpPtr := ctx.ResponseBodyPtr()
	params := ExampleParams{}
	if code, message := CommonCheckAndGetPostBody(&params, r); code != 0 {
		rpPtr.Code = code
		rpPtr.Message = message
		return
	}
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 20
	}
	skip := (params.Page - 1) * params.Limit
	var result []*cms.Example
	var mongo nonomongo.MongoContext
	defer mongo.Close()
	var where = bson.M{}

	if params.Row1 != nil {
		where["string"] = params.Row1
	}

	if params.Row2 != nil {
		where["string"] = params.Row2
	}

	if params.Row3 != nil {
		where["Date"] = params.Row3
	}

	if params.Row4 != nil {
		where["Date"] = params.Row4
	}

	if params.Row5 != nil {
		where["cms.Example2"] = params.Row5
	}

	err := mongo.FindAll(new(cms.Example), where).Sort("-_id").Skip(skip).Limit(params.Limit).Retry(3).List(&result)
	if err != nil {
		rpPtr.ErrorCodeAndMessage(1, err.Error())
		return
	}
	count, err := mongo.Count(new(cms.Example), where)
	if err != nil {
		rpPtr.ErrorCodeAndMessage(2, err.Error())
		return
	}
	body := ExampleRespBody{
		Models:    result,
		TotalRows: count,
	}
	rpPtr.SetBody(body)
}
func FindExample(ctx *RequestResponseContext, r *http.Request) {
	rpPtr := ctx.ResponseBodyPtr()
	_id := r.FormValue("_id")
	var result = &cms.Example{}
	var mongo nonomongo.MongoContext
	defer mongo.Close()
	if !bson.IsObjectIdHex(_id) {
		rpPtr.ErrorCodeAndMessage(1, "invalid object id")
		return
	}
	var query = bson.M{"_id": bson.ObjectIdHex(_id)}
	err := mongo.FindOne(result, query)
	if err != nil && err != mgo.ErrNotFound {
		MainLogger.Errorf("find Example err, _id =%v, err: %v", _id, err.Error())
		rpPtr.ErrorCodeAndMessage(2, "mongo find failed")
		return
	}
	rpPtr.SetBody(result)
}
func UpdateExample(ctx *RequestResponseContext, r *http.Request) {
	rpPtr := ctx.ResponseBodyPtr()
	params := UpdateExampleParams{}
	if code, message := CommonCheckAndGetPostBody(&params, r); code != 0 {
		rpPtr.Code = code
		rpPtr.Message = message
		return
	}
	var where = bson.M{}
	if bson.IsObjectIdHex(params.Id) {
		where["_id"] = bson.ObjectIdHex(params.Id)
	} else {
		where["_id"] = bson.NewObjectId()
	}
	var update = bson.M{
		"$set": bson.M{
			"regions":      params.Regions,
			"countries":    params.Countries,
			"platform":     params.Platform,
			"method":       params.Method,
			"currencies":   params.Currencies,
			"top_up_rules": params.TopUpRules,
			"update_time":  time.Now(),
		},
		"$setOnInsert": bson.M{"create_time": time.Now()},
	}
	var change = mgo.Change{
		Update:    update,
		Upsert:    true,
		ReturnNew: true,
	}
	var mongo nonomongo.MongoContext
	defer mongo.Close()
	var result = &cms.Example{}
	_, err := mongo.FindAndModify(result, where, change)
	if err != nil {
		rpPtr.ErrorCodeAndMessage(1, err.Error())
		return
	}
	// 更新操作中规则全部被删除的话自动删除整个配置
	if bson.IsObjectIdHex(params.Id) && len(result.TopUpRules) == 0 {
		DeleteExample(params.Id)
	}
	rpPtr.SetBody("success")
}

func DeleteExampleById(_id string) bool {
	if !bson.IsObjectIdHex(_id) {
		return false
	}
	var mongo nonomongo.MongoContext
	defer mongo.Close()
	var query = bson.M{"_id": bson.ObjectIdHex(_id)}
	err := mongo.Remove(new(cms.Example), query)
	if err != nil {
		MainLogger.Errorf("delete Example err, _id =%v, err: %v", _id, err.Error())
		return false
	}
	return true
}
