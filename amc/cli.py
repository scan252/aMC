from __future__ import annotations

import json
import subprocess
import sys
from pathlib import Path
from typing import Annotated, Optional

import typer
from rich.console import Console
from rich.table import Table

from amc import __version__
from amc.api_client import fetch_all_pools_sync
from amc.log_finder import discover_log_files, is_macos, resolve_log_path
from amc.log_parser import extract_credentials, extract_latest_url
from amc.models import GachaData, POOL_TYPE_NAMES
from amc.storage import (
    default_data_dir,
    list_player_ids,
    load_gacha_data,
    merge_gacha_data,
    player_data_path,
    save_gacha_data,
)

app = typer.Typer(
    name="amc",
    help="aMC — 鸣潮 Mac 版抽卡记录助手",
    no_args_is_help=True,
)
console = Console()


def _copy_to_clipboard(text: str) -> bool:
    if not is_macos():
        return False
    try:
        subprocess.run(["pbcopy"], input=text.encode("utf-8"), check=True)
        return True
    except (OSError, subprocess.CalledProcessError):
        return False


@app.callback()
def main(
    version: Annotated[
        Optional[bool],
        typer.Option("--version", "-V", help="显示版本号"),
    ] = None,
) -> None:
    if version:
        console.print(f"amc {__version__}")
        raise typer.Exit()


@app.command("status")
def status_cmd() -> None:
    """检查 Mac 鸣潮日志文件是否存在。"""
    if not is_macos():
        console.print("[yellow]警告: 当前不在 macOS 上运行，日志路径可能不可用。[/yellow]")

    locations = discover_log_files()
    if not locations:
        console.print("[red]未找到任何鸣潮 Mac 版日志文件。[/red]")
        console.print("\n请确认已安装 Mac 原生鸣潮客户端，并在游戏中打开「唤取 → 唤取记录」。")
        console.print("\n支持的日志路径:")
        from amc.log_finder import MAC_LOG_LOCATIONS

        for loc in MAC_LOG_LOCATIONS:
            console.print(f"  [{loc.label}] {loc.log_path}")
        raise typer.Exit(1)

    table = Table(title="已发现的鸣潮日志")
    table.add_column("服务器", style="cyan")
    table.add_column("路径")
    table.add_column("修改时间")

    for loc in locations:
        mtime = loc.log_path.stat().st_mtime
        from datetime import datetime

        table.add_row(
            loc.label,
            str(loc.log_path),
            datetime.fromtimestamp(mtime).strftime("%Y-%m-%d %H:%M:%S"),
        )

    console.print(table)

    players = list_player_ids(default_data_dir())
    if players:
        console.print(f"\n已保存的玩家 UID: {', '.join(players)}")


@app.command("url")
def url_cmd(
    log_path: Annotated[
        Optional[Path],
        typer.Option("--log", "-l", help="自定义 Client.log 路径"),
    ] = None,
    player_id: Annotated[
        Optional[str],
        typer.Option("--player", "-p", help="指定玩家 UID"),
    ] = None,
    copy: Annotated[
        bool,
        typer.Option("--copy", "-c", help="复制 URL 到剪贴板"),
    ] = True,
) -> None:
    """从游戏日志提取唤取记录 URL。"""
    try:
        path = resolve_log_path(log_path)
        url = extract_latest_url(path, player_id)
    except (FileNotFoundError, ValueError) as exc:
        console.print(f"[red]{exc}[/red]")
        raise typer.Exit(1) from exc

    console.print("[green]唤取记录 URL:[/green]")
    console.print(url)

    if copy and _copy_to_clipboard(url):
        console.print("\n[dim]已复制到剪贴板[/dim]")


@app.command("fetch")
def fetch_cmd(
    log_path: Annotated[
        Optional[Path],
        typer.Option("--log", "-l", help="自定义 Client.log 路径"),
    ] = None,
    player_id: Annotated[
        Optional[str],
        typer.Option("--player", "-p", help="指定玩家 UID"),
    ] = None,
    data_dir: Annotated[
        Optional[Path],
        typer.Option("--data-dir", help="数据存储目录"),
    ] = None,
    url: Annotated[
        Optional[str],
        typer.Option("--url", "-u", help="直接提供唤取记录 URL（跳过日志读取）"),
    ] = None,
) -> None:
    """抓取全部卡池的唤取历史记录并保存到本地。"""
    from amc.models import GachaCredentials

    storage_dir = data_dir or default_data_dir()

    try:
        if url:
            gacha_url = url
            creds = GachaCredentials.from_url(url)
        else:
            path = resolve_log_path(log_path)
            gacha_url, creds = extract_credentials(path, player_id)
    except (FileNotFoundError, ValueError) as exc:
        console.print(f"[red]{exc}[/red]")
        raise typer.Exit(1) from exc

    console.print(f"[cyan]玩家 UID:[/cyan] {creds.player_id}")
    console.print(f"[cyan]服务器:[/cyan] {'国际服' if creds.is_global else '国服'}")
    console.print(f"[dim]URL: {gacha_url[:80]}...[/dim]\n")

    def on_progress(current: int, total: int, pool_name: str, error: str | None = None) -> None:
        if error:
            console.print(f"  [{current}/{total}] {pool_name} [red]失败: {error}[/red]")
        else:
            console.print(f"  [{current}/{total}] 正在获取 {pool_name} ...")

    console.print("[bold]开始抓取唤取记录...[/bold]")
    pools = fetch_all_pools_sync(creds, on_progress=on_progress)

    total_records = sum(len(records) for records in pools.values())
    if total_records == 0:
        console.print(
            "\n[red]未获取到任何记录。[/red] "
            "唤取 URL 可能已过期（约 1 小时有效），请重新在游戏中打开唤取记录页面。"
        )
        raise typer.Exit(1)

    output_path = player_data_path(storage_dir, creds.player_id)
    existing = load_gacha_data(output_path)
    new_data = GachaData(
        player_id=creds.player_id,
        svr_area=creds.svr_area,
        fetched_at="",
        pools=pools,
    )

    if existing:
        final_data = merge_gacha_data(existing, pools)
        console.print("\n[dim]已与本地历史数据合并[/dim]")
    else:
        from datetime import datetime

        final_data = GachaData(
            player_id=creds.player_id,
            svr_area=creds.svr_area,
            fetched_at=datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
            pools=pools,
        )

    saved = save_gacha_data(output_path, final_data)
    console.print(f"\n[green]已保存 {total_records} 条新记录到:[/green] {saved}")

    table = Table(title="抽卡统计摘要")
    table.add_column("卡池", style="cyan")
    table.add_column("总数", justify="right")
    table.add_column("5★", justify="right", style="yellow")
    table.add_column("4★", justify="right", style="magenta")
    table.add_column("3★", justify="right")

    for pool_type, records in final_data.pools.items():
        if not records:
            continue
        pool_name = POOL_TYPE_NAMES.get(int(pool_type), pool_type)
        star5 = sum(1 for r in records if r.get("qualityLevel") == 5)
        star4 = sum(1 for r in records if r.get("qualityLevel") == 4)
        star3 = sum(1 for r in records if r.get("qualityLevel") == 3)
        table.add_row(pool_name, str(len(records)), str(star5), str(star4), str(star3))

    console.print(table)


@app.command("export")
def export_cmd(
    player_id: Annotated[str, typer.Argument(help="玩家 UID")],
    output: Annotated[
        Optional[Path],
        typer.Option("--output", "-o", help="导出文件路径"),
    ] = None,
    data_dir: Annotated[
        Optional[Path],
        typer.Option("--data-dir", help="数据存储目录"),
    ] = None,
) -> None:
    """导出已保存的抽卡数据为 JSON 文件。"""
    storage_dir = data_dir or default_data_dir()
    source = player_data_path(storage_dir, player_id)

    if not source.exists():
        console.print(f"[red]未找到玩家 {player_id} 的数据[/red]")
        raise typer.Exit(1)

    data = load_gacha_data(source)
    if data is None:
        console.print("[red]数据文件损坏[/red]")
        raise typer.Exit(1)

    dest = output or Path.cwd() / f"wuwa_pulls_{player_id}.json"
    with dest.open("w", encoding="utf-8") as file:
        json.dump(data.to_dict(), file, ensure_ascii=False, indent=2)

    console.print(f"[green]已导出到:[/green] {dest}")


def run() -> None:
    app()


if __name__ == "__main__":
    run()
