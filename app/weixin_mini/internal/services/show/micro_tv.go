package show

import (
	"encoding/json"
	"fmt"
	"guduo/app/internal/model_clean/micro_drama_rank_tmp_model"
)

func GetMicroRank(dayAt, pid int) []map[string]interface{} {
	var res []*micro_drama_rank_tmp_model.Table

	mdl := micro_drama_rank_tmp_model.Model().Select("name", "num", "platform_id")

	if pid > 0 {
		mdl = mdl.Where(fmt.Sprintf("JSON_CONTAINS(`platform_id`, '%d')", pid))
	}

	mdl.Order("num desc").
		Where("day_at", dayAt).
		Limit(50).
		Find(&res)
	ret := make([]map[string]interface{}, len(res))

	for k, row := range res {
		pidsStr := row.PlatformId
		pids := make([]interface{}, 0, 5)
		json.Unmarshal([]byte(pidsStr), &pids)
		ret[k] = map[string]interface{}{
			"name": row.Name,
			"platform": pids,
			"hot": row.Num,
			"trend": 0,
		}
	}

	return ret
}