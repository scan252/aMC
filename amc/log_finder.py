from __future__ import annotations

import os
from dataclasses import dataclass
from pathlib import Path


@dataclass(frozen=True)
class MacLogLocation:
    bundle_id: str
    label: str
    log_path: Path

    @property
    def exists(self) -> bool:
        return self.log_path.is_file()


# Mac 原生鸣潮客户端的沙盒日志路径（与 Windows 的 Client/Saved/Logs 不同）
MAC_LOG_LOCATIONS: tuple[MacLogLocation, ...] = (
    MacLogLocation(
        bundle_id="com.kurogame.mingchao",
        label="国服",
        log_path=Path.home()
        / "Library/Containers/com.kurogame.mingchao/Data/Library/Logs/Client/Client.log",
    ),
    MacLogLocation(
        bundle_id="com.kurogame.wutheringwaves.global",
        label="国际服",
        log_path=Path.home()
        / "Library/Containers/com.kurogame.wutheringwaves.global/Data/Library/Logs/Client/Client.log",
    ),
)

GACHA_URL_PATTERN = (
    r"https://aki-gm-resources(?:-oversea)?\.aki-game\.(?:net|com)"
    r"/aki/gacha/index\.html#/record[^\s\"']*"
)


def discover_log_files() -> list[MacLogLocation]:
    """返回当前 Mac 上已存在的鸣潮日志文件。"""
    return [loc for loc in MAC_LOG_LOCATIONS if loc.exists]


def find_best_log_file() -> MacLogLocation | None:
    """选择最近修改的日志文件（多服同时安装时取最新）。"""
    candidates = discover_log_files()
    if not candidates:
        return None
    return max(candidates, key=lambda loc: loc.log_path.stat().st_mtime)


def resolve_log_path(custom_path: str | Path | None = None) -> Path:
    """解析日志路径：优先自定义路径，否则自动发现。"""
    if custom_path:
        path = Path(custom_path).expanduser()
        if not path.is_file():
            raise FileNotFoundError(f"指定的日志文件不存在: {path}")
        return path

    location = find_best_log_file()
    if location is None:
        searched = "\n".join(f"  - {loc.log_path}" for loc in MAC_LOG_LOCATIONS)
        raise FileNotFoundError(
            "未找到鸣潮 Mac 版日志文件。请先在游戏中打开「唤取 → 唤取记录」，"
            f"然后重试。\n已搜索路径:\n{searched}"
        )
    return location.log_path


def is_macos() -> bool:
    return os.uname().sysname == "Darwin"
