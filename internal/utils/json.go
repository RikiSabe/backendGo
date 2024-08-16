package utils

import (
	"encoding/json"
)

func JsonDaoToModel(dao any, model interface{}) error {
	v, err := json.Marshal(dao)
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(string(v)), &model)
	return err
}
