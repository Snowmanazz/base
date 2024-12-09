package main

type Model struct {
	AppId   string `json:"app_id"`
	EntryId string `json:"entry_id"`
}

type DataValue struct {
	Value any `json:"value"`
}

type ModelOther struct {
	IsStartWorkflow bool   `json:"is_start_workflow"` //启用流程
	TransactionId   string `json:"transaction_id"`    //事务ID
}

type DataAddParam struct {
	Model
	ModelOther
	DataCreator    string               `json:"data_creator"`     //数据创建人
	IsStartTrigger bool                 `json:"is_start_trigger"` //触发智能助手
	Data           map[string]DataValue `json:"data"`             //数据
}

type DataBatchAddParam struct {
	Model
	ModelOther
	DataList []map[string]DataValue `json:"data_list"`
}

type DataBatchAddResponse struct {
	Status       string   `json:"status"`
	SuccessCount int      `json:"success_count"`
	SuccessIds   []string `json:"success_ids"`
}

type DataGetParam struct {
	Model
	DataId string `json:"data_id"`
}

type DataListParam struct {
	Model
	DataId string         `json:"data_id"`
	Limit  int            `json:"limit"`
	Fields []string       `json:"fields"`
	Filter DataListFilter `json:"filter"`
}
type DataListFilter struct {
	Rel  string               `json:"rel"`
	Cond []DataListFilterCond `json:"cond"`
}
type DataListFilterCond struct {
	Field  string `json:"field"`
	Type   string `json:"type"`
	Method string `json:"method"`
	Value  []any  `json:"value,omitempty"`
}
