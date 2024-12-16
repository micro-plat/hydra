package nfs

import (
	"errors"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/pkgs/security"
)

const TypeNodeName = "nfs"

//NFS 网络文件系统配置
type NFS struct {
	security.ConfEncrypt
	Local  string `json:"local,omitempty" toml:"local,omitempty"`
	Domain string `json:"domain,omitempty" toml:"domain,omitempty"`
	Rename bool   `json:"rename,omitempty" toml:"rename,omitempty"`
	Watch  bool   `json:"watch,omitempty" toml:"watch,omitempty"`

	CloudNFS  string `json:"cloudNfs,omitempty" toml:"cloudNfs,omitempty"`
	Endpoint  string `json:"endpoint,omitempty" toml:"endpoint,omitempty"`
	AccessKey string `json:"ak,omitempty" toml:"ak,omitempty"`
	SecretKey string `json:"sk,omitempty" toml:"sk,omitempty"`

	AllowDownload   bool   `json:"allowDownload,omitempty" toml:"allowDownload,omitempty"`
	DownloadService string `json:"downloadService,omitempty" toml:"downloadService,omitempty"`

	UploadService string `json:"uploadService,omitempty" toml:"uploadService,omitempty"`
	DiableUpload  bool   `json:"diableUpload,omitempty" toml:"diableUpload,omitempty"`

	ListFileService string `json:"listFileService,omitempty" toml:"listFileService,omitempty"`
	AllowListFile   bool   `json:"allowListFile,omitempty" toml:"allowListFile,omitempty"`

	ListDirService string `json:"listDirService,omitempty" toml:"listDirService,omitempty"`
	AllowListDir   bool   `json:"allowListDir,omitempty" toml:"allowListDir,omitempty"`

	PreviewService string `json:"previewService,omitempty" toml:"previewService,omitempty"`
	AllowPreview   bool   `json:"allowPreview,omitempty" toml:"allowPreview,omitempty"`

	ScaleImageService string `json:"scaleImageService,omitempty" toml:"scaleImageService,omitempty"`
	AllowScaleImage   bool   `json:"allowScaleImage,omitempty" toml:"allowScaleImage,omitempty"`

	CreateDirService string `json:"createDirService,omitempty" toml:"createDirService,omitempty"`
	AllowCreateDir   bool   `json:"allowCreateDir,omitempty" toml:"allowCreateDir,omitempty"`

	RenameDirService string `json:"renameDirService,omitempty" toml:"renameDirService,omitempty"`
	AllowRenameDir   bool   `json:"allowRenameDir,omitempty" toml:"allowRenameDir,omitempty"`

	Excludes []string `json:"excludes,omitempty"  label:"排除目录"`

	Includes []string `json:"includes,omitempty"  label:"目录"`

	Disable bool `json:"disable,omitempty" toml:"disable,omitempty"`
}

//New 构建mqc NFS配置，默认为对等模式
func New(local string, opts ...Option) *NFS {
	s := &NFS{Local: local}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

//GetConf 获取主配置信息
func GetConf(cnf conf.IServerConf) (*NFS, error) {
	s := NFS{}
	_, err := cnf.GetSubObject(TypeNodeName, &s)
	if errors.Is(err, conf.ErrNoSetting) {
		s.Disable = true
		return &s, nil
	}
	if err != nil {
		return nil, err
	}

	if b, err := govalidator.ValidateStruct(&s); !b {
		return nil, fmt.Errorf("mqc服务器配置数据有误:%v %v", err, s)
	}
	return &s, nil
}
