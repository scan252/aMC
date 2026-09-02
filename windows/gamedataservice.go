package main

import (
	"github.com/scan252/aMC/windows/internal/gamedata"
	"github.com/scan252/aMC/windows/internal/kurobbs"
)

// GameDataService 图鉴与声骸评分对前端的入口。
type GameDataService struct {
	db *gamedata.DB
}

func NewGameDataService() (*GameDataService, error) {
	db, err := gamedata.Open()
	if err != nil {
		return nil, err
	}
	return &GameDataService{db: db}, nil
}

// Characters 角色列表（支持名称/属性/武器过滤）。
func (s *GameDataService) Characters(keyword, element, weapon string) ([]gamedata.Character, error) {
	return s.db.SearchCharacters(keyword, element, weapon), nil
}

// EchoScoreMeta 声骸评分元数据：词条满值表 + 权重方案 + 声骸列表。
func (s *GameDataService) EchoScoreMeta() (gamedata.ScoreTables, []gamedata.Echo, error) {
	return s.db.ScoreTables(), s.db.Echoes(), nil
}

// ScoreEcho 声骸评分。
func (s *GameDataService) ScoreEcho(scheme string, subs []gamedata.SubstatValue) (gamedata.ScoreResult, error) {
	tables := s.db.ScoreTables()
	return gamedata.ScoreEcho(tables, scheme, subs), nil
}

// NewsService 资讯与兑换码对前端的入口。
type NewsService struct {
	client kurobbs.Client
}

func NewNewsService() *NewsService {
	return &NewsService{client: kurobbs.NewMockClient()}
}

// News 资讯列表（已绑定账号且真实模式时走官方接口，否则演示数据）。
func (s *NewsService) News(limit int) ([]kurobbs.ForumPost, error) {
	return s.client.ForumList(nil, limit)
}

// Codes 兑换码列表。
func (s *NewsService) Codes() ([]kurobbs.RedeemCode, error) {
	return s.client.RedeemCodes()
}
