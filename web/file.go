package web

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type FileUploader struct {
	// 文件Key
	FileFiled string

	// 存放路径
	// 为什么要用户传?
	// 要考虑文件名冲突问题
	// 所以很多时候, 目标文件名字都是随机的
	DstPathFunc func(*multipart.FileHeader) string
}

// Builder模式 + Handle

func (u FileUploader) Handle() HandleFunc {
	// 可以做额外检测
	if u.FileFiled == "" {
		u.FileFiled = "file"
	}

	if u.DstPathFunc == nil {
		// 设置默认值
	}

	return func(ctx *Context) {
		// 上传文件的逻辑在这里

		// 第一步: 读到文件内容
		// 第二步: 计算出目标路径
		// 第三步: 保存文件
		// 第四步: 返回响应
		file, fileHeader, err := ctx.Req.FormFile(u.FileFiled)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败")
			return
		}

		defer file.Close()

		// 我怎么知道目标路径(第二步)
		// 这种做法就是, 将目标路径计算逻辑, 交给用户
		dst := u.DstPathFunc(fileHeader)

		// 可以尝试把 dst 上不存在的目录全部建立起来
		//os.MkdirAll()

		// (第三步)
		// O_WRONLY 写入数据
		// O_TRUNC 如果文件本身存在, 清空数据
		// O_CREATE 如果文件不存在, 创建一个新的
		dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败")
			return
		}
		defer dstFile.Close()

		// 复制数据
		// buf 会影响你的性能
		// 你要考虑复用
		_, err = io.CopyBuffer(dstFile, file, nil)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败")
			return
		}

		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("上传成功")
	}
}

// Option模式 + HandleFunc
//type FileUploaderOption func(uploader *FileUploader)
//
//func NewFileUplooader(opts ...FileUploaderOption) *FileUploader {
//	res := &FileUploader{
//		FileFiled: "file",
//		DstPathFunc: func(header *multipart.FileHeader) string {
//			return filepath.Join("testdata", "upload", uuid.New().String())
//		},
//	}
//
//	res = &FileUploader{}
//	for _, opt := range opts {
//		opt(res)
//	}
//
//	return res
//}
//
//func (u FileUploader) HandleFunc(ctx *Context) {
//	// 文件上传
//}

type FileDownloader struct {
	Dir string
}

func (f FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		// 用的是 xxx?file=xxx
		req, err := ctx.QueryValue("file")
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("找不到目标文件")
			return
		}
		req = filepath.Clean(req)
		dst := filepath.Join(f.Dir, req)
		// 做一个校验, 防止相对路径引起攻击者下载了你的系统文件
		//dst, err = filepath.Abs(dst)
		//if strings.Contains(dst, f.Dir){
		//
		//}

		fn := filepath.Base(dst)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn) // 最重要的
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")                     // 控制缓存的头 加上不会缓存
		header.Set("Cache-Control", "must-revalidate") // 控制缓存的头
		header.Set("Pragma", "public")

		// 没有缓存
		http.ServeFile(ctx.Resp, ctx.Req, dst)
	}
}

type StaticResourceHandlerOption func(handler *StaticResourceHandler)

// StaticResourceHandler 两个层面上
// 1. 大文件不缓存
// 2.控制住了缓存的文件的数量
// 所以, 最多消耗多少内存 size(cache) * maxSize
type StaticResourceHandler struct {
	dir                     string
	extensionContentTypeMap map[string]string
	cache                   *lru.Cache[string, any]
	// 大文件不缓存
	maxSize int
}

func NewStaticResourceHandler(dir string, opts ...StaticResourceHandlerOption) (*StaticResourceHandler, error) {
	// 总共缓存 key-value的数量
	c, err := lru.New[string, any](1000)
	if err != nil {
		return nil, err
	}
	res := &StaticResourceHandler{
		dir:   dir,
		cache: c,
		// 10 MB, 文件大小超过这个值, 就不会缓存
		maxSize: 1024 * 1024 * 10,
		extensionContentTypeMap: map[string]string{
			// 这里根据自己的需要不断添加
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
		},
	}

	for _, opt := range opts {
		opt(res)
	}

	return res, nil
}

func StaticWithMaxFileSize(maxSize int) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.maxSize = maxSize
	}
}

func StaticWithCache(c *lru.Cache[string, any]) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.cache = c
	}
}

func StaticWithMoreExtension(extMap map[string]string) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		for ext, contentType := range extMap {
			handler.extensionContentTypeMap[ext] = contentType
		}
	}
}

func (s *StaticResourceHandler) Handle(ctx *Context) {
	// 无缓存
	// 1. 拿到目标文件名
	// 2. 定位到目标文件, 并且读出来
	// 3.返回给前端

	// 有缓存
	file, err := ctx.PathValue("file")
	if err != nil {
		ctx.RespStatusCode = http.StatusBadRequest
		ctx.RespData = []byte("请求路径不对")
		return
	}

	dst := filepath.Join(s.dir, file)
	ext := filepath.Ext(dst)[1:]
	header := ctx.Resp.Header()

	if data, ok := s.cache.Get(file); ok {
		header := ctx.Resp.Header()
		// 可能的有 文本文件, 图片, 多媒体(视频, 音频)
		header.Set("Content-Type", s.extensionContentTypeMap[ext])
		header.Set("Content-Length", strconv.Itoa(len(data.([]byte))))
		ctx.RespData = data.([]byte)
		ctx.RespStatusCode = 200
		return
	}

	data, err := os.ReadFile(dst)

	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("服务器错误")
	}

	// 大文件不缓存
	if len(data) <= s.maxSize {
		s.cache.Add(file, data)
	}
	// 可能的有 文本文件, 图片, 多媒体(视频, 音频)
	header.Set("Content-Type", s.extensionContentTypeMap[ext])
	header.Set("Content-Length", strconv.Itoa(len(data)))
	ctx.RespData = data
	ctx.RespStatusCode = 200

}
