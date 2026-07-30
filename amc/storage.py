from __future__ import annotations

import json
import shutil
from datetime import datetime
from pathlib import Path
from typing import Any

from amc.models import GachaData, GachaRecord


def default_data_dir() -> Path:
    return Path.home() / ".amc" / "data"


def player_data_path(data_dir: Path, player_id: str) -> Path:
    return data_dir / player_id / "gacha_data.json"


def backup_existing(path: Path) -> None:
    if not path.exists():
        return
    backup_dir = path.parent / "backup"
    backup_dir.mkdir(parents=True, exist_ok=True)
    stamp = datetime.now().strftime("%Y-%m-%d_%H%M%S")
    shutil.copy2(path, backup_dir / f"gacha_data_{stamp}.json")


def merge_pool_records(
    existing: list[dict[str, Any]],
    incoming: list[dict[str, Any]],
) -> list[dict[str, Any]]:
    """合并新旧记录，保留 API 已过期（超过约 6 个月）的本地历史。"""
    incoming_keys = {
        GachaRecord(
            card_pool_type=str(item.get("cardPoolType", "")),
            resource_id=int(item.get("resourceId", 0)),
            quality_level=int(item.get("qualityLevel", 0)),
            resource_type=str(item.get("resourceType", "")),
            name=str(item.get("name", "")),
            count=int(item.get("count", 1)),
            time=str(item.get("time", "")),
        ).unique_key()
        for item in incoming
    }

    merged = list(incoming)
    for item in existing:
        record = GachaRecord(
            card_pool_type=str(item.get("cardPoolType", "")),
            resource_id=int(item.get("resourceId", 0)),
            quality_level=int(item.get("qualityLevel", 0)),
            resource_type=str(item.get("resourceType", "")),
            name=str(item.get("name", "")),
            count=int(item.get("count", 1)),
            time=str(item.get("time", "")),
        )
        if record.unique_key() not in incoming_keys:
            merged.append(item)

    merged.sort(key=lambda item: item.get("time", ""), reverse=True)
    return merged


def merge_gacha_data(existing: GachaData, incoming_pools: dict[str, list[dict[str, Any]]]) -> GachaData:
    merged_pools: dict[str, list[dict[str, Any]]] = {}

    all_keys = set(existing.pools) | set(incoming_pools)
    for key in all_keys:
        old_records = existing.pools.get(key, [])
        new_records = incoming_pools.get(key, [])
        if new_records:
            merged_pools[key] = merge_pool_records(old_records, new_records)
        elif old_records:
            merged_pools[key] = old_records

    return GachaData(
        player_id=existing.player_id or "",
        svr_area=existing.svr_area,
        fetched_at=datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
        pools=merged_pools,
    )


def load_gacha_data(path: Path) -> GachaData | None:
    if not path.exists():
        return None
    with path.open("r", encoding="utf-8") as file:
        return GachaData.from_dict(json.load(file))


def save_gacha_data(path: Path, data: GachaData) -> Path:
    path.parent.mkdir(parents=True, exist_ok=True)
    backup_existing(path)
    with path.open("w", encoding="utf-8") as file:
        json.dump(data.to_dict(), file, ensure_ascii=False, indent=2)
    return path


def list_player_ids(data_dir: Path) -> list[str]:
    if not data_dir.exists():
        return []
    return sorted(
        entry.name
        for entry in data_dir.iterdir()
        if entry.is_dir() and (entry / "gacha_data.json").exists()
    )
