from __future__ import annotations

import asyncio
from typing import Any, Callable

import httpx

from amc.models import ALL_POOL_TYPES, GachaCredentials, GachaRecord, POOL_TYPE_NAMES

REQUEST_INTERVAL = 0.5
REQUEST_TIMEOUT = 15.0

DEFAULT_HEADERS = {
    "Content-Type": "application/json",
    "User-Agent": (
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "
        "AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
    ),
}


class GachaApiError(Exception):
    pass


async def fetch_pool_records(
    client: httpx.AsyncClient,
    creds: GachaCredentials,
    pool_type: int,
) -> list[GachaRecord]:
    url = f"{creds.api_base}/record/query"
    payload = creds.to_api_payload(pool_type)

    response = await client.post(url, json=payload)
    response.raise_for_status()
    result: dict[str, Any] = response.json()

    code = result.get("code", -1)
    if code not in (0, 200) and result.get("message") not in ("成功", "success"):
        message = result.get("message", "未知错误")
        raise GachaApiError(f"卡池 {pool_type} API 错误: code={code}, message={message}")

    data = result.get("data", [])
    if not isinstance(data, list):
        return []
    return [GachaRecord.from_api(item, pool_type) for item in data]


async def fetch_all_pools(
    creds: GachaCredentials,
    on_progress: callable | None = None,
) -> dict[str, list[dict[str, Any]]]:
    pools: dict[str, list[dict[str, Any]]] = {}
    total = 0

    async with httpx.AsyncClient(headers=DEFAULT_HEADERS, timeout=REQUEST_TIMEOUT) as client:
        for index, pool_type in enumerate(ALL_POOL_TYPES, start=1):
            pool_name = POOL_TYPE_NAMES.get(pool_type, str(pool_type))
            if on_progress:
                on_progress(index, len(ALL_POOL_TYPES), pool_name)

            try:
                records = await fetch_pool_records(client, creds, pool_type)
            except (httpx.HTTPError, GachaApiError) as exc:
                if on_progress:
                    on_progress(index, len(ALL_POOL_TYPES), pool_name, error=str(exc))
                records = []

            pools[str(pool_type)] = [record.to_dict() for record in records]
            total += len(records)

            if pool_type != ALL_POOL_TYPES[-1]:
                await asyncio.sleep(REQUEST_INTERVAL)

    return pools


def fetch_all_pools_sync(
    creds: GachaCredentials,
    on_progress: Callable[..., None] | None = None,
) -> dict[str, list[dict[str, Any]]]:
    return asyncio.run(fetch_all_pools(creds, on_progress))


def fetch_pool_records_sync(
    creds: GachaCredentials,
    pool_type: int,
) -> list[GachaRecord]:
    async def _run() -> list[GachaRecord]:
        async with httpx.AsyncClient(headers=DEFAULT_HEADERS, timeout=REQUEST_TIMEOUT) as client:
            return await fetch_pool_records(client, creds, pool_type)

    return asyncio.run(_run())
