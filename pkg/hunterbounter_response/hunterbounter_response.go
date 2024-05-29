package hunterbounter_response

import "hunterbounter.com/web-panel/pkg/hunterbounter_json"

/*
	success : bool
	message : string
	data : interface{}
*/

type HunterResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HunterBounterResponse(success bool, message string, data interface{}) *HunterResponse {
	return &HunterResponse{
		Success: success,
		Message: message,
		Data:    data,
	}
}

func (r *HunterResponse) SetSuccess(success bool) {
	r.Success = success
}

func (r *HunterResponse) SetMessage(message string) {
	r.Message = message
}

func (r *HunterResponse) SetData(data interface{}) {
	r.Data = data
}

func (r *HunterResponse) ToString() string {
	return hunterbounter_json.ToString(r)
}

func (r *HunterResponse) ToJSON() []byte {
	return hunterbounter_json.ToJson(r)
}
