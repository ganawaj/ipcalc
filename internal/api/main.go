package api

type IP struct {
	Address			string `json:"address"`
	Netmask			string `json:"mask"`
	Network 		string `json:"network"`
	Broadcast		string `json:"broadcast"`
	HostMin			string `json:"hostmin"`
	HostMax			string `json:"hostmax"`
	Error 			string `json:"error,omitempty"`
}