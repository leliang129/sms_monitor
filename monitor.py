#!/usr/bin/env python3
"""HeroSMS 号码价格监控脚本。

监控指定服务在指定价格以下的可用号码数量。
"""

import hashlib
import hmac
import base64
import json
import subprocess
import time
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, List, Optional
from urllib.parse import urlencode


CONFIG_PATH = Path(__file__).parent / "config.json"


def load_config() -> Dict[str, Any]:
    """加载配置文件。"""
    with open(CONFIG_PATH, "r", encoding="utf-8") as f:
        return json.load(f)


def curl_get(url: str, timeout: int = 30) -> str:
    """使用 curl 发送 GET 请求。"""
    result = subprocess.run(
        ["curl", "-s", "-L", "--max-time", str(timeout), url],
        capture_output=True, text=True, timeout=timeout + 5,
    )
    if result.returncode != 0:
        raise RuntimeError(f"curl 失败: {result.stderr}")
    return result.stdout


def curl_post(url: str, payload: str, timeout: int = 10) -> str:
    """使用 curl 发送 POST 请求。"""
    result = subprocess.run(
        ["curl", "-s", "-L", "--max-time", str(timeout),
         "-X", "POST", "-H", "Content-Type: application/json",
         "-d", payload, url],
        capture_output=True, text=True, timeout=timeout + 5,
    )
    if result.returncode != 0:
        raise RuntimeError(f"curl 失败: {result.stderr}")
    return result.stdout


def fetch_country_names(api_key: str, base_url: str) -> Dict[int, str]:
    """从 API 获取国家 ID -> 中文名称映射。"""
    params = urlencode({"action": "getCountries", "api_key": api_key})
    url = f"{base_url}?{params}"
    names: Dict[int, str] = {}
    try:
        data = json.loads(curl_get(url))
        for cid, info in data.items():
            names[int(cid)] = info.get("chn", info.get("eng", f"国家#{cid}"))
    except Exception:
        pass
    return names


def fetch_offers(api_key: str, base_url: str) -> Dict[str, Any]:
    """获取所有服务的报价数据。"""
    params = urlencode({"action": "getPrices", "api_key": api_key})
    url = f"{base_url}?{params}"
    return json.loads(curl_get(url))


def fetch_balance(api_key: str, base_url: str) -> Optional[float]:
    """查询账户余额。"""
    params = urlencode({"action": "getBalance", "api_key": api_key})
    url = f"{base_url}?{params}"
    try:
        text = curl_get(url, timeout=15)
        if text.startswith("ACCESS_BALANCE:"):
            return float(text.split(":")[1])
        return None
    except Exception:
        return None


def check_hits(
    offers: Dict[str, Any],
    watchlist: List[Dict[str, Any]],
    country_names: Dict[int, str],
) -> List[Dict[str, Any]]:
    """检查监控列表，返回命中项。

    offers 格式: {country_id: {service_id: {cost, count, physicalCount}}}
    """
    hits: List[Dict[str, Any]] = []

    for item in watchlist:
        service = item["service"]
        country = str(item["country"])
        max_price = item["max_price"]
        service_name = item.get("service_name", service)
        country_name = item.get("country_name", country_names.get(int(country), f"国家#{country}"))

        # 获取该国家该服务的数据
        country_data = offers.get(country)
        if not country_data:
            continue

        service_data = country_data.get(service)
        if not service_data:
            continue

        cost = service_data.get("cost", 0)
        count = service_data.get("count", 0)
        physical_count = service_data.get("physicalCount", 0)

        # 检查是否低于目标价格且有号码
        if cost <= max_price and count > 0:
            hits.append({
                "service": service,
                "service_name": service_name,
                "country": int(country),
                "country_name": country_name,
                "cost": cost,
                "count": count,
                "physical_count": physical_count,
                "max_price": max_price,
            })

    return hits


def format_hit(hit: Dict[str, Any], checked_at: str) -> List[str]:
    """格式化单条命中消息。"""
    return [
        f"🎯 {hit['service_name']} 号码命中",
        "",
        hit["country_name"],
        f"💰 当前价格：${hit['cost']}",
        f"📦 可用数量：{hit['count']}",
        f"🎚 目标价格：≤ ${hit['max_price']}",
        "",
        f"✅ 可用号码数：{hit['count']}（物理卡 {hit['physical_count']}）",
        f"🕒 检查时间：{checked_at}",
    ]


def format_no_hit(watchlist: List[Dict[str, Any]], checked_at: str) -> List[str]:
    """格式化无命中消息。"""
    lines = ["📭 本轮检查无命中", ""]
    for item in watchlist:
        lines.append(f"  {item.get('service_name', item['service'])} @ {item.get('country_name', str(item['country']))}：≤ ${item['max_price']}")
    lines.append("")
    lines.append(f"🕒 检查时间：{checked_at}")
    return lines


def print_lines(lines: List[str]) -> None:
    """打印消息行。"""
    for line in lines:
        print(line)


def gen_feishu_sign(secret: str, timestamp: int) -> str:
    """生成飞书 Webhook 签名。"""
    string_to_sign = f"{timestamp}\n{secret}"
    hmac_code = hmac.new(string_to_sign.encode("utf-8"), digestmod=hashlib.sha256).digest()
    return base64.b64encode(hmac_code).decode("utf-8")


def send_feishu_webhook(webhook_url: str, hits: List[Dict[str, Any]], checked_at: str, sign: Optional[str] = None) -> None:
    """通过飞书 Webhook 发送命中通知。"""
    for hit in hits:
        content = (
            f"🎯 {hit['service_name']} 号码命中\n"
            f"国家：{hit['country_name']}\n"
            f"💰 当前价格：${hit['cost']}\n"
            f"📦 可用数量：{hit['count']}\n"
            f"🎚 目标价格：≤ ${hit['max_price']}\n"
            f"✅ 可用号码数：{hit['count']}（物理卡 {hit['physical_count']}）\n"
            f"🕒 检查时间：{checked_at}"
        )
        body: Dict[str, Any] = {
            "msg_type": "text",
            "content": {"text": content},
        }

        # 加签
        if sign:
            timestamp = int(time.time())
            body["timestamp"] = str(timestamp)
            body["sign"] = gen_feishu_sign(sign, timestamp)

        payload = json.dumps(body)
        try:
            result = curl_post(webhook_url, payload)
            print(f"📨 飞书通知已发送：{result}")
        except Exception as e:
            print(f"❌ 飞书通知失败：{e}")


def run_once(config: Dict[str, Any], country_names: Dict[int, str]) -> List[Dict[str, Any]]:
    """执行一次检查，返回命中列表。"""
    api_key = config["api_key"]
    base_url = config.get("base_url", "https://hero-sms.com/stubs/handler_api.php")
    watchlist = config["watchlist"]
    checked_at = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    if not api_key:
        print("❌ 错误：未配置 api_key，请编辑 config.json")
        return []

    # 查询余额
    balance = fetch_balance(api_key, base_url)
    if balance is not None:
        print(f"💰 账户余额：${balance}")

    # 获取报价
    offers = fetch_offers(api_key, base_url)

    # 检查命中
    hits = check_hits(offers, watchlist, country_names)

    if hits:
        for hit in hits:
            print_lines(format_hit(hit, checked_at))
            print()

        # 飞书通知
        webhook = config.get("webhook", {})
        if webhook.get("enabled") and webhook.get("url"):
            send_feishu_webhook(webhook["url"], hits, checked_at, webhook.get("sign"))
    else:
        print_lines(format_no_hit(watchlist, checked_at))
        print()

    return hits


def main() -> None:
    """主函数：循环监控。"""
    config = load_config()
    api_key = config["api_key"]
    base_url = config.get("base_url", "https://hero-sms.com/stubs/handler_api.php")
    interval = config.get("interval_seconds", 60)

    # 启动时获取国家名称映射
    country_names = fetch_country_names(api_key, base_url)

    print("=" * 50)
    print("  HeroSMS 号码价格监控")
    print("=" * 50)
    print(f"监控间隔：{interval} 秒")
    print(f"监控项数：{len(config['watchlist'])}")
    print(f"国家库：{len(country_names)} 个国家")
    print("=" * 50)
    print()

    while True:
        try:
            run_once(config, country_names)
        except KeyboardInterrupt:
            print("\n监控已停止。")
            break
        except Exception as e:
            print(f"❌ 检查出错：{e}")

        try:
            time.sleep(interval)
        except KeyboardInterrupt:
            print("\n监控已停止。")
            break


if __name__ == "__main__":
    main()
