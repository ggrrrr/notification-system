package sms

type To struct {
	To string `json:"to"`
}

type Message struct {
	From string `json:"from"`
	Dest []To   `json:"destinations"`
	Text string `json:"text"`
}

type Request struct {
	Messages []Message `json:"messages"`
}

type Response struct {
	BulkId string       `json:"bulkId"`
	Error  RequestError `json:"requestError"`
}

type RequestError struct {
	Exception ServiceException `json:"serviceException"`
}

type ServiceException struct {
	MessageId string `json:"messageId"`
}
