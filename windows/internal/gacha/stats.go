package gacha

import (
	"sort"
	"strconv"
)

// 概率常量（用于期望值与欧非指数计算）。
const (
	// ExpectedAvgPity 综合平均出金期望（社区通用近似值，含软保底效应）。
	ExpectedAvgPity = 62.0
	// HardPity 硬保底抽数（角色/武器/常驻均为 80）。
	HardPity = 80
	// SoftPityStart 软保底起始抽数。
	SoftPityStart = 65
)

// PoolStats 单卡池统计。
type PoolStats struct {
	PoolType   int    `json:"poolType"`
	PoolName   string `json:"poolName"`
	Total      int    `json:"total"`      // 本窗口内记录数
	Count5     int    `json:"count5"`     // 5★ 数量
	Count4     int    `json:"count4"`     // 4★ 数量
	Count3     int    `json:"count3"`     // 3★ 数量
	AvgPity    float64 `json:"avgPity"`   // 平均出金抽数（0 表示样本不足）
	Pity       int    `json:"pity"`       // 当前已垫抽数
	PityIsFloor bool  `json:"pityIsFloor"` // true 表示窗口内无 5★，pity 为下界
	Last5Name  string `json:"last5Name"`  // 最近 5★ 名称
	Last5Time  string `json:"last5Time"`
}

// OverallStats 全账号统计。
type OverallStats struct {
	Total      int         `json:"total"`
	Count5     int         `json:"count5"`
	Count4     int         `json:"count4"`
	AvgPity    float64     `json:"avgPity"`
	LuckIndex  float64     `json:"luckIndex"` // 期望/实际，>1 偏欧 <1 偏非
	Pools      []PoolStats `json:"pools"`
}

// ComputeStats 计算指定卡池记录的整体统计。
// 记录需为该 UID 全部已合并数据（pools 按 ID 分组）。
func ComputeStats(data *GachaData) OverallStats {
	overall := OverallStats{Pools: []PoolStats{}}

	for _, poolType := range AllPoolTypes {
		key := itoa(poolType)
		records := data.Pools[key]
		if len(records) == 0 {
			continue
		}
		ps := computePoolStats(poolType, records)
		overall.Pools = append(overall.Pools, ps)
		overall.Total += ps.Total
		overall.Count5 += ps.Count5
		overall.Count4 += ps.Count4
	}

	if overall.Total > 0 && overall.Count5 > 0 {
		overall.AvgPity = float64(overall.Total) / float64(overall.Count5)
		overall.LuckIndex = ExpectedAvgPity / overall.AvgPity
	}
	return overall
}

func computePoolStats(poolType int, records []Record) PoolStats {
	ps := PoolStats{PoolType: poolType, PoolName: PoolTypeNames[poolType], Total: len(records)}

	asc := make([]Record, len(records))
	copy(asc, records)
	sort.SliceStable(asc, func(i, j int) bool { return asc[i].Time < asc[j].Time })

	// 5★ 间隔样本与最近 5★
	sinceLast := 0
	found5 := false
	var gaps []int
	for _, r := range asc {
		switch r.QualityLevel {
		case 5:
			ps.Count5++
			if found5 {
				// 间隔 = 自上一个 5★ 以来的抽数（含本抽）
				gaps = append(gaps, sinceLast+1)
			}
			ps.Last5Name = r.Name
			ps.Last5Time = r.Time
			sinceLast = 0
			found5 = true
		case 4:
			ps.Count4++
			sinceLast++
		default:
			ps.Count3++
			sinceLast++
		}
	}
	if len(gaps) > 0 {
		sum := 0
		for _, g := range gaps {
			sum += g
		}
		ps.AvgPity = float64(sum) / float64(len(gaps))
	}
	// 平均出金 = 总抽数 / 5★ 数（含窗口边缘，与主界面指标口径一致）
	if ps.Count5 > 0 {
		ps.AvgPity = float64(ps.Total) / float64(ps.Count5)
	}

	if !found5 {
		ps.Pity = len(asc)
		ps.PityIsFloor = true
	} else {
		ps.Pity = sinceLast
	}
	return ps
}

// PityDist 返回出金间隔分布（桶：≤10,11-20,…,71-75,76-80）。
// 仅统计有 5★ 间隔样本的卡池，用于直方图展示。
func PityDist(data *GachaData) map[string][]int {
	// 桶边界与 demo 图表一致：9 个桶
	buckets := make(map[string][]int)
	for _, poolType := range AllPoolTypes {
		key := itoa(poolType)
		records := data.Pools[key]
		if len(records) == 0 {
			continue
		}
		asc := make([]Record, len(records))
		copy(asc, records)
		sort.SliceStable(asc, func(i, j int) bool { return asc[i].Time < asc[j].Time })

		sinceLast := 0
		found5 := false
		dist := make([]int, 9)
		for _, r := range asc {
			if r.QualityLevel == 5 {
				if found5 {
					dist[bucketIndex(sinceLast+1)]++
				}
				found5 = true
				sinceLast = 0
			} else {
				sinceLast++
			}
		}
		buckets[key] = dist
	}
	return buckets
}

func bucketIndex(n int) int {
	switch {
	case n <= 10:
		return 0
	case n <= 20:
		return 1
	case n <= 30:
		return 2
	case n <= 40:
		return 3
	case n <= 50:
		return 4
	case n <= 60:
		return 5
	case n <= 70:
		return 6
	case n <= 75:
		return 7
	default:
		return 8
	}
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
