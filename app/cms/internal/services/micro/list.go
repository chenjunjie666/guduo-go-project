package micro

import (
	"encoding/json"
	"guduo/app/internal/model_clean/micro_drama_rank_tmp_model"
)

func List(dayAt int) []map[string]interface{} {
	var res []*micro_drama_rank_tmp_model.Table

	mdl := micro_drama_rank_tmp_model.Model().Select("name", "num", "platform_id")
	mdl.Order("num desc").
		Where("day_at", dayAt).
		Find(&res)
	ret := make([]map[string]interface{}, len(res))

	for k, row := range res {
		pidsStr := row.PlatformId
		pids := make([]interface{}, 0, 5)
		json.Unmarshal([]byte(pidsStr), &pids)
		ret[k] = map[string]interface{}{
			"name": row.Name,
			"platform": pids,
			"num": row.Num,
		}
	}

	return ret
}

func SaveRank(data []*micro_drama_rank_tmp_model.Table, dayAt uint) {
	micro_drama_rank_tmp_model.Model().Where("day_at", dayAt).Delete(nil)
	micro_drama_rank_tmp_model.Model().Create(&data)
}