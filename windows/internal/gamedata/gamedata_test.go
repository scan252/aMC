package gamedata

import (
	"math"
	"testing"
)

func TestOpenSeed(t *testing.T) {
	db, err := Open()
	if err != nil {
		t.Fatal(err)
	}
	if len(db.Characters()) < 20 {
		t.Fatalf("种子角色数量异常: %d", len(db.Characters()))
	}
	if len(db.Echoes()) < 5 {
		t.Fatalf("种子声骸数量异常: %d", len(db.Echoes()))
	}
	tables := db.ScoreTables()
	if len(tables.Substats) < 10 || len(tables.Schemes) < 2 {
		t.Fatalf("评分表异常: %+v", tables)
	}
}

func TestSearchCharacters(t *testing.T) {
	db, _ := Open()
	if got := db.SearchCharacters("今汐", "", ""); len(got) != 1 {
		t.Fatalf("按名称搜索应命中今汐: %v", got)
	}
	if got := db.SearchCharacters("", "冷凝", ""); len(got) < 2 {
		t.Fatalf("按属性搜索冷凝: %v", got)
	}
	if got := db.SearchCharacters("", "不存在", ""); len(got) != 0 {
		t.Fatalf("无效属性应无结果: %v", got)
	}
}

func TestScoreEcho(t *testing.T) {
	db, _ := Open()
	tables := db.ScoreTables()

	// 双爆通用方案，满分副词条（每词条均取 1 个满-roll 等效）
	sub := []SubstatValue{
		{Name: "暴击", Value: 6.3},
		{Name: "暴击伤害", Value: 12.6},
		{Name: "攻击力%", Value: 5.8},
	}
	res := ScoreEcho(tables, "双爆通用", sub)

	// 三条均等效 1 词条，权重均 >0 → 等效词条 = 3
	if math.Abs(res.EquivalentRolls-3.0) > 0.001 {
		t.Fatalf("等效词条应为 3，got %v", res.EquivalentRolls)
	}
	// 有权重的词条合计 3/5（+25 声骸最多 5 词条），百分比 = 60
	if math.Abs(res.Percent-30.0) > 0.1 {
		t.Fatalf("评分应为 30%%，got %v", res.Percent)
	}
	if res.Grade != "C" {
		t.Fatalf("30%% 应评为 C，got %s", res.Grade)
	}

	// 零权重词条不应计入
	res2 := ScoreEcho(tables, "双爆通用", []SubstatValue{
		{Name: "小防御", Value: 34.0},
	})
	if res2.EquivalentRolls != 0 || res2.Percent != 0 {
		t.Fatalf("零权重词条评分为 0: %+v", res2)
	}

	// 未知词条被忽略且报告
	res3 := ScoreEcho(tables, "双爆通用", []SubstatValue{{Name: "不存在词条", Value: 1}})
	if len(res3.Unknown) != 1 {
		t.Fatalf("未知词条应被报告: %+v", res3)
	}
}

func TestGradeThresholds(t *testing.T) {
	cases := map[string]struct {
		pct  float64
		want string
	}{
		"S": {80, "S"},
		"A": {60, "A"},
		"B": {45, "B"},
		"C": {30, "C"},
		"D": {10, "D"},
	}
	for _, c := range cases {
		if got := gradeOf(c.pct); got != c.want {
			t.Fatalf("pct=%v want=%s got=%s", c.pct, c.want, got)
		}
	}
}
