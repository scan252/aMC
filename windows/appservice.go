package main

// AppService 暴露给前端的应用级基础服务（版本信息等）。
type AppService struct{}

func (s *AppService) Version() string {
	return "0.1.0"
}

func (s *AppService) Platform() string {
	return "windows"
}
