// Package gamedata 游戏数据服务：内置种子数据库 + 版本更新通道设计。
// 数据文件 data/seed.json 随应用打包；后续通过远端版本通道按版本号增量更新。
package gamedata

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed data/seed.json
var seedJSON []byte

// Character 角色条目（种子库最小字段集，待扩展）。
type Character struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Rarity  int    `json:"rarity"`
	Element string `json:"element"`
	Weapon  string `json:"weapon"`
}

// Echo 声骸条目：套装与主词条属性域。
type Echo struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Set  string   `json:"set"`
	Cost int      `json:"cost"`
	Main []string `json:"main"`
}

// EchoSubstat 声骸副词条定义（Max = 单词条满值）。
type EchoSubstat struct {
	Name string  `json:"name"`
	Max  float64 `json:"max"`
}

// WeightScheme 声骸评分权重方案。
type WeightScheme struct {
	Name  string             `json:"name"`
	Note  string             `json:"note"`
	Ratio map[string]float64 `json:"ratio"`
}

// ScoreTables 声骸评分数据表。
type ScoreTables struct {
	Substats []EchoSubstat  `json:"substats"`
	Schemes  []WeightScheme `json:"schemes"`
}

// SeedData 内置种子数据库。
type SeedData struct {
	Version     string      `json:"version"`
	Note        string      `json:"note"`
	Characters  []Character `json:"characters"`
	Echoes      []Echo      `json:"echoes"`
	ScoreTables ScoreTables `json:"scoreTables"`
}

// DB 数据库句柄。
type DB struct {
	seed SeedData
}

// Open 加载内置数据。
func Open() (*DB, error) {
	var s SeedData
	if err := json.Unmarshal(seedJSON, &s); err != nil {
		return nil, fmt.Errorf("内置种子库解析失败: %w", err)
	}
	return &DB{seed: s}, nil
}

// Version 数据库版本号。
func (db *DB) Version() string { return db.seed.Version }

// Note 数据说明。
func (db *DB) Note() string { return db.seed.Note }

// Characters 全部角色。
func (db *DB) Characters() []Character { return db.seed.Characters }

// SearchCharacters 按名称/属性/武器类型过滤。
func (db *DB) SearchCharacters(keyword, element, weapon string) []Character {
	var out []Character
	for _, c := range db.seed.Characters {
		if keyword != "" && !strings.Contains(c.Name, keyword) {
			continue
		}
		if element != "" && c.Element != element {
			continue
		}
		if weapon != "" && c.Weapon != weapon {
			continue
		}
		out = append(out, c)
	}
	return out
}

// Echoes 全部声骸。
func (db *DB) Echoes() []Echo { return db.seed.Echoes }

// ScoreTables 返回评分表。
func (db *DB) ScoreTables() ScoreTables { return db.seed.ScoreTables }
