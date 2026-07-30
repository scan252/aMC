from __future__ import annotations

import re
from pathlib import Path

from amc.models import GachaCredentials

LOG_MAGIC = b"\xa5\xef\xa5"
GACHA_URL_RE = re.compile(
    r"https://aki-gm-resources(?:-oversea)?\.aki-game\.(?:net|com)"
    r"/aki/gacha/index\.html#/record[?=&\w\-%.]+"
)


def decrypt_client_log(data: bytes) -> bytes:
    """解密鸣潮 Client.log（Windows / Mac 新版均可能加密）。"""
    if len(data) < 3:
        return data
    if data[:3] == LOG_MAGIC:
        payload = data[3:]
    else:
        payload = data
    out = bytearray()
    for byte in payload:
        out.append(byte ^ (0xA5 if byte % 2 == 1 else 0xEF))
    return bytes(out)


def _decode_log_text(data: bytes) -> str:
    return data.decode("utf-8", errors="ignore").lstrip("\ufeff")


def is_log_encrypted(data: bytes) -> bool:
    if len(data) < 3:
        return False
    if data[:3] == LOG_MAGIC:
        return True
    if data[0] == 0:
        decoded = decrypt_client_log(data)
        return _decode_log_text(decoded).startswith("Log file open")
    return False


def read_log_content(log_path: Path) -> str:
    raw = log_path.read_bytes()
    if is_log_encrypted(raw):
        raw = decrypt_client_log(raw)
    return _decode_log_text(raw)


def extract_gacha_urls(content: str) -> list[str]:
    return GACHA_URL_RE.findall(content)


def extract_latest_url(
    log_path: Path,
    player_id: str | None = None,
) -> str:
    """从日志中提取最近一次唤取记录 URL。"""
    content = read_log_content(log_path)
    urls = extract_gacha_urls(content)
    if not urls:
        raise ValueError(
            f"在日志 {log_path.name} 中未找到唤取记录 URL。"
            "请先在游戏中打开「唤取 → 唤取记录」页面。"
        )

    if player_id:
        for url in reversed(urls):
            if player_id in url:
                return url

    return urls[-1]


def extract_credentials(
    log_path: Path,
    player_id: str | None = None,
) -> tuple[str, GachaCredentials]:
    url = extract_latest_url(log_path, player_id)
    return url, GachaCredentials.from_url(url)
