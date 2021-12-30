package infs

//notExcludes 服务路径中排除的
var NOTEXCLUDES = []string{"/**/_/nfs/**"}

const (
	FILENAME = "name"

	DIRNAME = "dir"

	NDIRNAME = "ndir"
)

const (
	//SVSUpload 用户端上传文件
	SVSUpload = "/nfs/upload"

	SVSPreview = "/nfs/preview"

	//SVSDonwload 用户端下载文件
	SVSDonwload = "/nfs/file/:dir/:name"

	//SVSList 文件列表
	SVSList = "/nfs/file/list"

	//SVSDir 目录列表
	SVSDir = "/nfs/dir/list"

	//SVSScalrImage 压缩文件
	SVSScalrImage = "/nfs/scale/:name"

	//SVSCreateDir 创建文件目录
	SVSCreateDir = "/nfs/create/:dir"

	//SVSRenameDir 重命名文件目录
	SVSRenameDir = "/nfs/create/:dir/:ndir"

	//获取远程文件的指纹信息
	RMT_FP_GET = "/_/nfs/fp/get"

	//推送指纹数据
	RMT_FP_NOTIFY = "/_/nfs/fp/notify"

	//拉取指纹列表
	RMT_FP_QUERY = "/_/nfs/fp/query"

	//获取远程文件数据
	RMT_FILE_DOWNLOAD = "/_/nfs/file/download"
)
