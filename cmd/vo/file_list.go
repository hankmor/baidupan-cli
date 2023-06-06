package vo

type FileListResp struct {
	BaseVo
	Files []FileList `json:"list,omitempty"`
}

type FileList struct {
	FsId           uint64 `json:"fs_id"`           // 文件在云端的唯一标识ID
	Path           string `json:"path"`            // 文件的绝对路径
	ServerFilename string `json:"server_filename"` // 文件名称
	Size           uint   `json:"size"`            // 文件大小，单位B
	ServerMtime    uint   `json:"server_mtime"`    // 文件在服务器修改时间
	ServerCtime    uint   `json:"server_ctime"`    // 文件在服务器创建时间
	LocalMtime     uint   `json:"local_mtime"`     // 文件在客户端修改时间
	LocalCtime     uint   `json:"local_ctime"`     // 文件在客户端创建时间
	IsDir          uint   `json:"isdir"`           // 是否为目录，0 文件、1 目录
	Category       uint   `json:"category"`        // 文件类型，1 视频、2 音频、3 图片、4 文档、5 应用、6 其他、7 种子
	Md5            string `json:"md5"`             // 云端哈希（非文件真实MD5），只有是文件类型时，该字段才存在
	DirEmpty       int    `json:"dir_empty"`       // 该目录是否存在子目录，只有请求参数web=1且该条目为目录时，该字段才存在， 0为存在， 1为不存在
}
