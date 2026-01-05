package cmd

import (
	"baidupan-cli/app"
	openapi "baidupan-cli/openxpanapi"
	"baidupan-cli/util"
	"context"
	"fmt"
	pathpkg "path"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/desertbit/grumble"
	"github.com/liushuochen/gotable"
)

// =============================================
// 文件列表查询
// =============================================

const (
	orderByTime = "time"
	orderByName = "name"
	orderBySize = "size"

	Yes = "Y"
	No  = "N"
)

var fileListCmd = &grumble.Command{
	Name:    "ls",
	Aliases: []string{"list"},
	Help:    "show file lists",
	Usage:   "ls [OPTIONS]",
	Flags: func(f *grumble.Flags) {
		f.String("d", "dir", "", "the directory to show files in it (default: current directory)")
		f.String("o", "order", "name", `order type, support 'time','name' and 'size', default is 'name':
				1. time: sort files by file type first, then sort by modification time
				2. name: sort files by file type first, then sort by file name
				3. size: sort files by file type first, then sort by file size`)
		f.Bool("r", "recurse", false, "whether to list files recursively")
		f.Bool("D", "desc", false, "whether to sort in descending order")
		f.Bool("f", "only-folder", false, "whether only to query folders, so files will be filtered")
		f.Bool("F", "only-files", false, "whether only to query files, so folders will be filtered")
		f.Bool("E", "show-empty", false, "whether to show empty folder info, ONLY SUPPORT when `recurse` is false")
		f.Int("l", "limit", 1000, "the number of queries, default is 1000, and it is recommended that the maximum number not exceed 1000")
		f.Bool("v", "verbose", false, "whether to show verbose info of files")
		f.Bool("H", "human-readable", false, "whether to show files info as human-readable")
		f.Bool("g", "show-form", false, "whether to show files info as form, ONLY SUPPORT when `verbose` is true")
		// f.StringL("ctime", "", "creation time to filter, when the creation time of the file is greater than it will be list, ONLY SUPPORTED when `recurse` is true")
		// f.StringL("mtime", "", "update time to filter, when the modification time of the file is greater than it will be list, ONLY SUPPORTED when `recurse` is true")
	},
	Run: func(ctx *grumble.Context) error {
		if err := checkAuthorized(ctx); err != nil {
			return err
		}

		options := NewFileListOptions()
		verbose := ctx.Flags.Bool("verbose")
		showForm := ctx.Flags.Bool("show-form")
		dir := ctx.Flags.String("dir")
		dir = ResolvePath(dir)
		recurse := ctx.Flags.Bool("recurse")
		desc := ctx.Flags.Bool("desc")
		humanReadable := ctx.Flags.Bool("human-readable")
		if desc {
			options.Desc()
		}
		limit := ctx.Flags.Int("limit")
		if limit <= 0 {
			return fmt.Errorf("invalid parameter: limit must be great than 0, but found %d", limit)
		}
		options.Limit(int32(limit))
		order := ctx.Flags.String("order")
		switch order {
		case orderByTime:
			options.OrderByTime()
		case orderByName:
			options.OrderByName()
		case orderBySize:
			options.OrderBySize()
		default:
			return fmt.Errorf("invalid parameter: unsupported order type %s", order)
		}
		onlyFolder := ctx.Flags.Bool("only-folder")
		onlyFiles := ctx.Flags.Bool("only-files")
		if onlyFolder && onlyFiles {
			return fmt.Errorf("both giving \"only-folder\" and \"only-files\" tags are ambiguous and not allowed")
		}
		if onlyFolder {
			options.OnlyDir()
		}
		if onlyFiles {
			options.OnlyFiles()
		}
		showEmpty := ctx.Flags.Bool("show-empty")
		if showEmpty {
			options.ShowEmpty()
		}

		var fileLister FileLister
		if recurse {
			fileLister = &RecursionFileLister{}
		} else {
			fileLister = &SimpleFileLister{}
		}
		files, err := fileLister.List(dir, *options)
		if err != nil {
			return err
		}
		return fileLister.Print(dir, files, FilePrinterOption{Verbose: verbose, HumanReadable: humanReadable, ShowForm: showForm})
	},
}

// 打印输出

type FilePrinterOption struct {
	Verbose       bool
	HumanReadable bool
	ShowForm      bool
}

type FileListPrinter interface {
	Print(root string, files []*File, options FilePrinterOption) error
}

// 查询目录下的文件列表

type FileListResp struct {
	BaseVo
	GuidInfo  string  `json:"guid_info,omitempty"`
	RequestId uint64  `json:"request_id,omitempty"`
	Guid      int     `json:"guid,omitempty"`
	Files     []*File `json:"list,omitempty"`
}

type File struct {
	FsId           uint64            `json:"fs_id,omitempty"`           // 文件在云端的唯一标识ID
	Path           string            `json:"path,omitempty"`            // 文件的绝对路径
	ServerFilename string            `json:"server_filename,omitempty"` // 文件名称
	Size           uint              `json:"size,omitempty"`            // 文件大小，单位B
	ServerMtime    uint              `json:"server_mtime,omitempty"`    // 文件在服务器修改时间
	ServerCtime    uint              `json:"server_ctime,omitempty"`    // 文件在服务器创建时间
	LocalMtime     uint              `json:"local_mtime,omitempty"`     // 文件在客户端修改时间
	LocalCtime     uint              `json:"local_ctime,omitempty"`     // 文件在客户端创建时间
	IsDir          uint              `json:"isdir,omitempty"`           // 是否为目录，0 文件、1 目录
	Category       uint              `json:"category,omitempty"`        // 文件类型，1 视频、2 音频、3 图片、4 文档、5 应用、6 其他、7 种子
	Md5            string            `json:"md5,omitempty"`             // 云端哈希（非文件真实MD5），只有是文件类型时，该字段才存在
	DirEmpty       int               `json:"dir_empty,omitempty"`       // 该目录是否存在子目录，只有请求参数web=1且该条目为目录时，该字段才存在， 0为存在， 1为不存在
	Thumbs         map[string]string `json:"thumbs,omitempty"`          // 缩略图地址
}

func fileParentDir(p, name string) string {
	p = strings.TrimSpace(p)
	name = strings.TrimSpace(name)
	if p == "" {
		return "/"
	}
	// prefer stripping suffix if possible
	if name != "" && strings.HasSuffix(p, name) {
		p = strings.TrimSuffix(p, name)
	} else {
		// fallback
		p = pathpkg.Dir(strings.TrimRight(p, "/"))
	}
	if p == "" {
		p = "/"
	}
	if !strings.HasSuffix(p, "/") {
		p += "/"
	}
	return p
}

type FileLister interface {
	FileListPrinter
	List(path string, options FileListOptions) ([]*File, error)
}

type SimpleFileLister struct{}

func (sfl *SimpleFileLister) List(Path string, options FileListOptions) ([]*File, error) {
	req := app.APIClient.FileinfoApi.Xpanfilelist(RootContext)
	reqptr := &req

	sfl.applyOptions(reqptr, options)

	resp, _, err := reqptr.Dir(Path).AccessToken(*TokenResp.AccessToken).Execute()
	if err != nil {
		return nil, err
	}

	var fileListResp FileListResp
	err = sonic.UnmarshalString(resp, &fileListResp)
	if err != nil {
		return nil, err
	}
	if fileListResp.Success() {
		// 只显示文件
		if options.onlyFile {
			var fs []*File
			for _, f := range fileListResp.Files {
				if f.IsDir == 0 {
					fs = append(fs, f)
				}
			}
			return fs, nil
		}
		return fileListResp.Files, nil
	}
	return nil, fmt.Errorf("error code: %d", fileListResp.Errno)
}

func (sfl *SimpleFileLister) applyOptions(req *openapi.ApiXpanfilelistRequest, options FileListOptions) {
	if options.desc {
		req.Desc(1) // 降序排序
	}
	if options.order != "" {
		req.Order(options.order) // 排序属性
	}
	if options.onlyDir {
		req.Folder("1") // 只显示文件夹
	}
	if options.showEmpty {
		req.Showempty(1) // 显示空文件夹信息
	}
	if options.limit > 0 {
		req.Limit(options.limit) // 文件数量限制
	}
}

func (sfl *SimpleFileLister) Print(root string, files []*File, options FilePrinterOption) error {
	// 表格输出详细信息
	if options.Verbose {
		if options.ShowForm {
			table, err := gotable.Create("NAME", "VALUE")
			if err != nil {
				return err
			}
			var sizestr, lctime, lmtime, sctime, smtime string
			var i int
			for _, f := range files {
				if options.HumanReadable {
					sizestr = util.ConvReadableSize(int64(f.Size))
					lctime = util.ConvTimestamp(int64(f.LocalCtime))
					lmtime = util.ConvTimestamp(int64(f.LocalMtime))
					sctime = util.ConvTimestamp(int64(f.ServerCtime))
					smtime = util.ConvTimestamp(int64(f.ServerMtime))
				} else {
					sizestr = util.Int64ToStr(int64(f.Size))
					lctime = util.Int64ToStr(int64(f.LocalCtime))
					lmtime = util.Int64ToStr(int64(f.LocalMtime))
					sctime = util.Int64ToStr(int64(f.ServerCtime))
					smtime = util.Int64ToStr(int64(f.ServerMtime))
				}
				_ = table.AddRow([]string{"FsId", util.Int64ToStr(int64(f.FsId))})
				_ = table.AddRow([]string{"Path", fileParentDir(f.Path, f.ServerFilename)})
				_ = table.AddRow([]string{"Name", f.ServerFilename})
				_ = table.AddRow([]string{"Dir", getDirLabel(int(f.IsDir))})
				_ = table.AddRow([]string{"Category", getCategoryLabel(int(f.Category))})
				_ = table.AddRow([]string{"Md5", f.Md5})
				_ = table.AddRow([]string{"Size", sizestr})
				_ = table.AddRow([]string{"Local CTime", lctime})
				_ = table.AddRow([]string{"Local MTime", lmtime})
				_ = table.AddRow([]string{"Server CTime", sctime})
				_ = table.AddRow([]string{"Server MTime", smtime})
				if i < len(files)-1 {
					_ = table.AddRow([]string{"-", "-"})
				}
				i++
			}
			fmt.Println(table)
		} else {
			table, err := gotable.Create("FsId", "Path", "Name", "Dir", "Category", "MD5", "Size", "Local CTime", "Local MTime", "Server CTime", "Server MTime")
			if err != nil {
				return err
			}
			var sizestr, lctime, lmtime, sctime, smtime string
			for _, f := range files {
				if options.HumanReadable {
					sizestr = util.ConvReadableSize(int64(f.Size))
					lctime = util.ConvTimestamp(int64(f.LocalCtime))
					lmtime = util.ConvTimestamp(int64(f.LocalMtime))
					sctime = util.ConvTimestamp(int64(f.ServerCtime))
					smtime = util.ConvTimestamp(int64(f.ServerMtime))
				} else {
					sizestr = util.Int64ToStr(int64(f.Size))
					lctime = util.Int64ToStr(int64(f.LocalCtime))
					lmtime = util.Int64ToStr(int64(f.LocalMtime))
					sctime = util.Int64ToStr(int64(f.ServerCtime))
					smtime = util.Int64ToStr(int64(f.ServerMtime))
				}
				_ = table.AddRow([]string{util.Int64ToStr(int64(f.FsId)), fileParentDir(f.Path, f.ServerFilename), f.ServerFilename, getDirLabel(int(f.IsDir)), getCategoryLabel(int(f.Category)), f.Md5, sizestr, lctime, lmtime, sctime, smtime})
			}
			fmt.Println(table)
		}
	} else {
		for _, f := range files {
			fmt.Println(f.ServerFilename)
		}
	}
	return nil
}

// 递归查询文件列表

type RecursionFileResp struct {
	BaseVo
	RequestId string  `json:"request_id,omitempty"`
	HasMore   int     `json:"has_more"` // 是否还有下一页，0表示无，1表示有
	Cursor    int     `json:"cursor"`   // 当还有下一页时，为下一次查询的起点
	Files     []*File `json:"list"`     // 文件列表
}

type RecursionFileLister struct{}

func (rfl *RecursionFileLister) List(Path string, options FileListOptions) ([]*File, error) {
	req := app.APIClient.MultimediafileApi.Xpanfilelistall(context.Background())
	reqptr := &req

	rfl.applyOptions(reqptr, options)

	resp, _, err := reqptr.Path(Path).Recursion(1).AccessToken(*TokenResp.AccessToken).Execute()
	if err != nil {
		return nil, err
	}

	var recursionFileResp RecursionFileResp
	err = sonic.UnmarshalString(resp, &recursionFileResp)
	if err != nil {
		return nil, err
	}
	if recursionFileResp.Success() {
		// 进显示文件夹
		if options.onlyDir {
			var fs []*File
			for _, f := range recursionFileResp.Files {
				if f.IsDir == 1 {
					fs = append(fs, f)
				}
			}
			return fs, nil
		} else if options.onlyFile {
			var fs []*File
			for _, f := range recursionFileResp.Files {
				if f.IsDir == 0 {
					fs = append(fs, f)
				}
			}
			return fs, nil
		}
		return recursionFileResp.Files, nil
	}
	return nil, fmt.Errorf("error code: %d", recursionFileResp.Errno)
}

func (rfl *RecursionFileLister) applyOptions(req *openapi.ApiXpanfilelistallRequest, options FileListOptions) {
	if options.desc {
		req.Desc(1)
	}
	if options.order != "" {
		req.Order(options.order)
	}
	if options.limit > 0 {
		req.Limit(options.limit)
	}
}

func (rfl *RecursionFileLister) Print(root string, files []*File, options FilePrinterOption) error {
	// 表格输出详细信息
	if options.Verbose {
		if options.ShowForm {
			table, err := gotable.Create("NAME", "VALUE")
			if err != nil {
				return err
			}
			var sizestr, lctime, lmtime, sctime, smtime string
			var i int
			for _, f := range files {
				if options.HumanReadable {
					sizestr = util.ConvReadableSize(int64(f.Size))
					lctime = util.ConvTimestamp(int64(f.LocalCtime))
					lmtime = util.ConvTimestamp(int64(f.LocalMtime))
					sctime = util.ConvTimestamp(int64(f.ServerCtime))
					smtime = util.ConvTimestamp(int64(f.ServerMtime))
				} else {
					sizestr = util.Int64ToStr(int64(f.Size))
					lctime = util.Int64ToStr(int64(f.LocalCtime))
					lmtime = util.Int64ToStr(int64(f.LocalMtime))
					sctime = util.Int64ToStr(int64(f.ServerCtime))
					smtime = util.Int64ToStr(int64(f.ServerMtime))
				}
				_ = table.AddRow([]string{"FsId", util.Int64ToStr(int64(f.FsId))})
				_ = table.AddRow([]string{"Path", fileParentDir(f.Path, f.ServerFilename)})
				_ = table.AddRow([]string{"Name", getIndentedName(root, f.Path, f.ServerFilename)})
				_ = table.AddRow([]string{"Dir", getDirLabel(int(f.IsDir))})
				_ = table.AddRow([]string{"Category", getCategoryLabel(int(f.Category))})
				_ = table.AddRow([]string{"MD5", f.Md5})
				_ = table.AddRow([]string{"Size", sizestr})
				_ = table.AddRow([]string{"Local CTime", lctime})
				_ = table.AddRow([]string{"Local MTime", lmtime})
				_ = table.AddRow([]string{"Server CTime", sctime})
				_ = table.AddRow([]string{"Server MTime", smtime})
				if i < len(files)-1 {
					_ = table.AddRow([]string{"-", "-"})
				}
				i++
			}
			fmt.Println(table)
		} else {
			table, err := gotable.Create("FsId", "Path", "Name", "Dir", "Category", "MD5", "Size", "Local CTime", "Local MTime", "Server CTime", "Server MTime")
			if err != nil {
				return err
			}
			var sizestr, lctime, lmtime, sctime, smtime string
			for _, f := range files {
				if options.HumanReadable {
					sizestr = util.ConvReadableSize(int64(f.Size))
					lctime = util.ConvTimestamp(int64(f.LocalCtime))
					lmtime = util.ConvTimestamp(int64(f.LocalMtime))
					sctime = util.ConvTimestamp(int64(f.ServerCtime))
					smtime = util.ConvTimestamp(int64(f.ServerMtime))
				} else {
					sizestr = util.Int64ToStr(int64(f.Size))
					lctime = util.Int64ToStr(int64(f.LocalCtime))
					lmtime = util.Int64ToStr(int64(f.LocalMtime))
					sctime = util.Int64ToStr(int64(f.ServerCtime))
					smtime = util.Int64ToStr(int64(f.ServerMtime))
				}
				_ = table.AddRow([]string{util.Int64ToStr(int64(f.FsId)), fileParentDir(f.Path, f.ServerFilename), getIndentedName(root, f.Path, f.ServerFilename), getDirLabel(int(f.IsDir)), getCategoryLabel(int(f.Category)), f.Md5, sizestr, lctime, lmtime, sctime, smtime})
			}
			fmt.Println(table)
		}
	} else {
		for _, f := range files {
			fmt.Println(getIndentedName(root, f.Path, f.ServerFilename) + getTrailingSlash(int(f.IsDir)))
		}
	}
	return nil
}

func getDirLabel(i int) string {
	if i == 1 {
		return "d"
	}
	return "f"
}

func getTrailingSlash(isDir int) string {
	if isDir == 1 {
		return "/"
	}
	return ""
}

func getIndentedName(root, path, name string) string {
	// Clean paths to ensure consistent format
	root = strings.TrimSuffix(pathpkg.Clean(root), "/")
	path = pathpkg.Clean(path)
	
	// If path is just root/name...
	if !strings.HasPrefix(path, root) {
		return name
	}
	
	rel := strings.TrimPrefix(path, root)
    rel = strings.TrimPrefix(rel, "/") // Remove leading slash of relative part
    
    // Parent dir of file relative to root
    // rel: c/d.txt -> depth 1
    // rel: c -> depth 0
    
    depth := strings.Count(rel, "/")
    if depth > 0 {
    	return strings.Repeat("  ", depth) + name
    }
    return name
}

func getCategoryLabel(i int) string {
	switch i {
	case 1:
		return "视频"
	case 2:
		return "音乐"
	case 3:
		return "图片"
	case 4:
		return "文档"
	case 5:
		return "应用"
	case 6:
		return "其他"
	case 7:
		return "种子"
	default:
		return "unknown"
	}
}

// 选项

type FileListOptions struct {
	desc      bool
	order     string
	limit     int32
	onlyDir   bool
	onlyFile  bool
	showEmpty bool
	// createTime string
	// updateTime string
}

func NewFileListOptions() *FileListOptions {
	return &FileListOptions{}
}

func (options *FileListOptions) Desc() *FileListOptions {
	options.desc = true
	return options
}

func (options *FileListOptions) OrderByTime() *FileListOptions {
	options.order = orderByTime
	return options
}

func (options *FileListOptions) OrderByName() *FileListOptions {
	options.order = orderByName
	return options
}

func (options *FileListOptions) OrderBySize() *FileListOptions {
	options.order = orderBySize
	return options
}

func (options *FileListOptions) Limit(limit int32) *FileListOptions {
	options.limit = limit
	return options
}

func (options *FileListOptions) OnlyDir() *FileListOptions {
	options.onlyDir = true
	return options
}

func (options *FileListOptions) OnlyFiles() *FileListOptions {
	options.onlyFile = true
	return options
}

func (options *FileListOptions) ShowEmpty() *FileListOptions {
	options.showEmpty = true
	return options
}

// func (options *FileListOptions) CreateTimeFrom(ctime string) *FileListOptions {
// 	options.createTime = ctime
// 	return options
// }
//
// func (options *FileListOptions) UpdateTimeFrom(mtime string) *FileListOptions {
// 	options.updateTime = mtime
// 	return options
// }
