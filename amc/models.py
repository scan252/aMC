from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any


POOL_TYPE_NAMES: dict[int, str] = {
    1: "角色活动唤取",
    2: "武器活动唤取",
    3: "角色常驻唤取",
    4: "武器常驻唤取",
    5: "新手唤取",
    6: "新手自选唤取",
    7: "感恩定向唤取",
    8: "角色新旅唤取",
    9: "武器新旅唤取",
    10: "角色联动唤取",
    11: "武器联动唤取",
    12: "角色忆旅唤取",
    13: "武器忆旅唤取",
}

ALL_POOL_TYPES = list(POOL_TYPE_NAMES.keys())

API_BASE_CN = "https://gmserver-api.aki-game2.com/gacha"
API_BASE_GLOBAL = "https://gmserver-api.aki-game2.net/gacha"


@dataclass
class GachaCredentials:
    record_id: str
    player_id: str
    server_id: str
    card_pool_id: str
    language_code: str = "zh-Hans"
    svr_area: str = "cn"
    is_global: bool = False

    @classmethod
    def from_url(cls, url: str) -> GachaCredentials:
        from urllib.parse import parse_qs, urlparse

        normalized = url.replace("#", "")
        parsed = urlparse(normalized)
        params = parse_qs(parsed.query)

        def get(key: str, default: str = "") -> str:
            values = params.get(key, [])
            return values[0] if values else default

        host = parsed.netloc
        is_global = host.endswith(".net") or "oversea" in host

        return cls(
            record_id=get("record_id"),
            player_id=get("player_id"),
            server_id=get("svr_id"),
            card_pool_id=get("resources_id") or get("gacha_id"),
            language_code=get("lang", "zh-Hans"),
            svr_area=get("svr_area", "global" if is_global else "cn"),
            is_global=is_global,
        )

    def to_api_payload(self, card_pool_type: int) -> dict[str, Any]:
        return {
            "recordId": self.record_id,
            "playerId": self.player_id,
            "serverId": self.server_id,
            "cardPoolId": self.card_pool_id,
            "cardPoolType": card_pool_type,
            "languageCode": self.language_code,
        }

    @property
    def api_base(self) -> str:
        return API_BASE_GLOBAL if self.is_global else API_BASE_CN


@dataclass
class GachaRecord:
    card_pool_type: str
    resource_id: int
    quality_level: int
    resource_type: str
    name: str
    count: int
    time: str

    @classmethod
    def from_api(cls, raw: dict[str, Any], pool_type: int) -> GachaRecord:
        return cls(
            card_pool_type=str(raw.get("cardPoolType", pool_type)),
            resource_id=int(raw.get("resourceId", 0)),
            quality_level=int(raw.get("qualityLevel", 0)),
            resource_type=str(raw.get("resourceType", "")),
            name=str(raw.get("name", "")),
            count=int(raw.get("count", 1)),
            time=str(raw.get("time", "")),
        )

    def to_dict(self) -> dict[str, Any]:
        return {
            "cardPoolType": self.card_pool_type,
            "resourceId": self.resource_id,
            "qualityLevel": self.quality_level,
            "resourceType": self.resource_type,
            "name": self.name,
            "count": self.count,
            "time": self.time,
        }

    def unique_key(self) -> str:
        return f"{self.time}|{self.name}|{self.quality_level}"


@dataclass
class GachaData:
    player_id: str
    svr_area: str
    fetched_at: str
    pools: dict[str, list[dict[str, Any]]] = field(default_factory=dict)

    def to_dict(self) -> dict[str, Any]:
        return {
            "player_id": self.player_id,
            "svr_area": self.svr_area,
            "fetched_at": self.fetched_at,
            "pools": self.pools,
        }

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> GachaData:
        pools = data.get("pools") or data.get("pulls") or {}
        return cls(
            player_id=str(data.get("player_id", data.get("uid", ""))),
            svr_area=str(data.get("svr_area", "cn")),
            fetched_at=str(data.get("fetched_at", "")),
            pools={str(k): v for k, v in pools.items()},
        )
