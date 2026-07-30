import pytest

from amc.log_parser import decrypt_client_log, extract_gacha_urls, is_log_encrypted
from amc.models import GachaCredentials, GachaRecord
from amc.storage import merge_pool_records


SAMPLE_URL = (
    "https://aki-gm-resources.aki-game.com/aki/gacha/index.html#/record?"
    "svr_id=76402e5b20be2c39f095a152090afddc&"
    "player_id=106507910&"
    "lang=zh-Hans&"
    "gacha_id=4&"
    "gacha_type=6&"
    "svr_area=cn&"
    "record_id=b184031785f7f64c2f3917d7f1202dc2&"
    "resources_id=917dfa695d6c6634ee4e972bb9168f6a"
)

SAMPLE_GLOBAL_URL = (
    "https://aki-gm-resources-oversea.aki-game.net/aki/gacha/index.html#/record?"
    "svr_id=abc123&player_id=999888&lang=en&gacha_id=1&gacha_type=1&"
    "svr_area=global&record_id=rec123&resources_id=res456"
)

SAMPLE_LOG_LINE = (
    '{"title":"","url":"'
    + SAMPLE_URL
    + '","transparent":true,"titlebar":false}'
)


class TestLogParser:
    def test_extract_url_from_log_line(self):
        urls = extract_gacha_urls(SAMPLE_LOG_LINE)
        assert len(urls) == 1
        assert urls[0].startswith("https://aki-gm-resources")

    def test_extract_multiple_urls_returns_all(self):
        content = SAMPLE_LOG_LINE + "\n" + SAMPLE_LOG_LINE.replace("106507910", "200000001")
        urls = extract_gacha_urls(content)
        assert len(urls) == 2

    def test_decrypt_plain_text_unchanged(self):
        plain = b"Log file open, hello world"
        assert decrypt_client_log(plain) != plain or not is_log_encrypted(plain)


class TestCredentials:
    def test_parse_cn_url(self):
        creds = GachaCredentials.from_url(SAMPLE_URL)
        assert creds.player_id == "106507910"
        assert creds.server_id == "76402e5b20be2c39f095a152090afddc"
        assert creds.record_id == "b184031785f7f64c2f3917d7f1202dc2"
        assert creds.card_pool_id == "917dfa695d6c6634ee4e972bb9168f6a"
        assert creds.svr_area == "cn"
        assert creds.is_global is False

    def test_parse_global_url(self):
        creds = GachaCredentials.from_url(SAMPLE_GLOBAL_URL)
        assert creds.player_id == "999888"
        assert creds.is_global is True

    def test_api_payload(self):
        creds = GachaCredentials.from_url(SAMPLE_URL)
        payload = creds.to_api_payload(1)
        assert payload["cardPoolType"] == 1
        assert payload["playerId"] == "106507910"


class TestStorage:
    def test_merge_preserves_old_records(self):
        old = [
            {"time": "2024-01-01 00:00:00", "name": "旧角色", "qualityLevel": 5},
        ]
        new = [
            {"time": "2025-01-01 00:00:00", "name": "新角色", "qualityLevel": 5},
        ]
        merged = merge_pool_records(old, new)
        assert len(merged) == 2
        names = {item["name"] for item in merged}
        assert names == {"旧角色", "新角色"}

    def test_merge_deduplicates(self):
        record = {"time": "2025-01-01 00:00:00", "name": "角色A", "qualityLevel": 5}
        merged = merge_pool_records([record], [record])
        assert len(merged) == 1


class TestGachaRecord:
    def test_unique_key(self):
        record = GachaRecord(
            card_pool_type="1",
            resource_id=1,
            quality_level=5,
            resource_type="角色",
            name="测试",
            count=1,
            time="2025-01-01 00:00:00",
        )
        assert record.unique_key() == "2025-01-01 00:00:00|测试|5"
