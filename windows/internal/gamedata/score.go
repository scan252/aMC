package gamedata

// SubstatValue 一条副词条：名称与数值。
type SubstatValue struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// ScoredSubstat 评分明细。
type ScoredSubstat struct {
	Name   string  `json:"name"`
	Value  float64 `json:"value"`
	Rolls  float64 `json:"rolls"`  // 等效满-roll 词条数
	Weight float64 `json:"weight"` // 方案权重（0 表示不计分）
}

// ScoreResult 声骸评分结果。
type ScoreResult struct {
	// EquivalentRolls 权重内副词条的等效满-roll 词条数合计
	EquivalentRolls float64 `json:"equivalentRolls"`
	// Percent 评分百分比 = EquivalentRolls / 10 × 100（5 槽位 × 2 roll 满值）
	Percent float64         `json:"percent"`
	Grade   string          `json:"grade"`
	Detail  []ScoredSubstat `json:"detail"`
	Unknown []string        `json:"unknown,omitempty"`
	Scheme  string          `json:"scheme"`
}

// ScoreEcho 按指定权重方案为声骸副词条评分。
// 词条满值表来自种子库 scoreTables.substats；未知词条被忽略并列入 Unknown。
func ScoreEcho(t ScoreTables, schemeName string, subs []SubstatValue) ScoreResult {
	res := ScoreResult{Scheme: schemeName, Detail: []ScoredSubstat{}}

	maxOf := map[string]float64{}
	for _, s := range t.Substats {
		maxOf[s.Name] = s.Max
	}
	scheme := t.Schemes[0]
	for _, sc := range t.Schemes {
		if sc.Name == schemeName {
			scheme = sc
			break
		}
	}

	for _, sub := range subs {
		if sub.Name == "" {
			continue
		}
		max, ok := maxOf[sub.Name]
		rolls := 0.0
		if ok && max > 0 {
			rolls = sub.Value / max
		}
		weight := scheme.Ratio[sub.Name]
		if !ok {
			res.Unknown = append(res.Unknown, sub.Name)
			continue
		}
		if weight > 0 {
			res.EquivalentRolls += rolls
		}
		res.Detail = append(res.Detail, ScoredSubstat{
			Name:   sub.Name,
			Value:  sub.Value,
			Rolls:  rolls,
			Weight: weight,
		})
	}

	res.Percent = res.EquivalentRolls / 10.0 * 100
	res.Grade = gradeOf(res.Percent)
	return res
}

// gradeOf 百分比 → 评级。基准：5 条有效词条 = 100%。
func gradeOf(pct float64) string {
	switch {
	case pct >= 70:
		return "S"
	case pct >= 55:
		return "A"
	case pct >= 40:
		return "B"
	case pct >= 25:
		return "C"
	default:
		return "D"
	}
}
