package gacha

import "sort"

// MergePoolRecords 合并新旧记录：以新记录为主，保留 API 已过期（约 6 个月）的本地历史，
// 按 time 降序排序（移植自 storage.py merge_pool_records）。
func MergePoolRecords(existing, incoming []Record) []Record {
	incomingKeys := make(map[string]struct{}, len(incoming))
	for _, r := range incoming {
		incomingKeys[r.UniqueKey()] = struct{}{}
	}

	merged := make([]Record, 0, len(incoming)+len(existing))
	merged = append(merged, incoming...)
	for _, r := range existing {
		if _, ok := incomingKeys[r.UniqueKey()]; !ok {
			merged = append(merged, r)
		}
	}

	SortRecordsDesc(merged)
	return merged
}

// SortRecordsDesc 按 time 字段降序（字符串序即时间序，格式 YYYY-MM-DD HH:MM:SS）。
func SortRecordsDesc(records []Record) {
	sort.SliceStable(records, func(i, j int) bool {
		return records[i].Time > records[j].Time
	})
}
