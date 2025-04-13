package dto

type OutgoingHook struct {
	DocumentID []string `json:"document_id"`
	Auth       struct {
		Domain         string `json:"domain"`
		ClientEndpoint string `json:"client_endpoint"`
		ServerEndpoint string `json:"server_endpoint"`
		MemberID       string `json:"member_id"`
	} `json:"auth"`
}
