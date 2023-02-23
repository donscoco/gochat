package model

/*
{"code":200,"message":"成功","data":[{"urls":"stun:www.boxim.online:3478","username":"admin","credential":"admin123"},{"urls":"turn:www.boxim.online:3478","username":"admin","credential":"admin123"}]}
*/

type WebRTCOffer struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"` // 请求call的值是 "offer" ，接受accept的值是 "answer"
}

type WebRTCOfferOutput struct {
	Content  string `json:"content"`
	Id       int    `json:"id"`
	RecvId   int    `json:"recvId"`
	SendId   int    `json:"sendId"`
	SendTime *int64 `json:"sendTime"` //null
	Type     int    `json:"type"`     //101
}

const (
	RTC_CALL      = 101
	RTC_ACCEPT    = 102
	RTC_REJECT    = 103
	RTC_CANCEL    = 104
	RTC_FAILED    = 105
	RTC_HANDUP    = 106
	RTC_CANDIDATE = 107
)

type WebRTCCandidate struct {
	Candidate     string `json:"candidate"`
	SDPMid        string `json:"sdpMid"`
	SDPMLineIndex int    `json:"sdpMLineIndex"`
}

type WebRTCCandidateOutput struct {
	Content  string `json:"content"`
	Id       int    `json:"id"`
	RecvId   int    `json:"recvId"`
	SendId   int    `json:"sendId"`
	SendTime *int64 `json:"sendTime"` //null
	Type     int    `json:"type"`     //101
}
