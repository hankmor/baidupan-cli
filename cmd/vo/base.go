package vo

type BaseVo struct {
	Errno     int    `json:"errno,omitempty"`
	GuidInfo  string `json:"guid_info,omitempty"`
	RequestId uint64 `json:"request_id,omitempty"`
	Guid      int    `json:"guid,omitempty"`
}
