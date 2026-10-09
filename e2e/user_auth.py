#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
e2e：用户注册 / 登录接口测试

前置条件：
  1. MySQL 容器已启动（hazmat-mysql）
  2. hazmat-agent server 已监听（默认 http://127.0.0.1:8080）

用法：
  python3 e2e/user_auth.py
  BASE_URL=http://127.0.0.1:8080 python3 e2e/user_auth.py
"""

from __future__ import annotations

import json
import os
import sys
import time
import urllib.error
import urllib.request
from typing import Any

# ---------- 配置 ----------
BASE_URL = os.environ.get("BASE_URL", "http://127.0.0.1:8080").rstrip("/")
USERNAME = f"e2e_{int(time.time())}"
PASSWORD = "TestPass123!"

# 统计通过 / 失败数
pass_count = 0
fail_count = 0


def post_json(path: str, payload: Any) -> tuple[int, Any, str]:
    """发送 POST JSON 请求，返回 (状态码, 解析后的 body, 原始文本)。"""
    url = f"{BASE_URL}{path}"
    # payload 为 str 时按原样发送（用于测非法 JSON）；否则序列化为 JSON
    if isinstance(payload, str):
        body_bytes = payload.encode("utf-8")
    else:
        body_bytes = json.dumps(payload, ensure_ascii=False).encode("utf-8")

    req = urllib.request.Request(
        url,
        data=body_bytes,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    print(f"  → POST {url}")
    print(f"    请求体: {body_bytes.decode('utf-8')}")

    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            raw = resp.read().decode("utf-8")
            code = resp.getcode()
    except urllib.error.HTTPError as e:
        raw = e.read().decode("utf-8")
        code = e.code

    # 尽量解析 JSON，失败则保留原文
    try:
        data = json.loads(raw) if raw else None
    except json.JSONDecodeError:
        data = raw

    print(f"  ← HTTP {code}  响应: {raw}")
    return code, data, raw


def assert_case(
    name: str,
    expect_code: int,
    actual_code: int,
    raw: str,
    expect_substr: str | None = None,
) -> None:
    """断言 HTTP 状态码，以及响应体是否包含指定子串。"""
    global pass_count, fail_count

    if actual_code != expect_code:
        print(f"[失败] {name}: 期望 HTTP {expect_code}，实际 {actual_code}")
        print(f"       响应体: {raw}")
        fail_count += 1
        return

    if expect_substr and expect_substr not in raw:
        print(f"[失败] {name}: 响应体中未找到「{expect_substr}」")
        print(f"       响应体: {raw}")
        fail_count += 1
        return

    print(f"[通过] {name}（HTTP {actual_code}）")
    pass_count += 1


def main() -> int:
    global pass_count, fail_count

    print("=" * 50)
    print("e2e 用户认证测试开始")
    print(f"  服务地址 BASE_URL = {BASE_URL}")
    print(f"  测试用户 USERNAME = {USERNAME}")
    print("=" * 50)
    print()

    # ----- 用例 1：注册成功 -----
    print("【用例 1】正常注册，期望 200 且 msg=注册成功")
    code, data, raw = post_json(
        "/user/register",
        {"username": USERNAME, "password": PASSWORD},
    )
    assert_case("注册成功", 200, code, raw, "注册成功")
    print()

    # ----- 用例 2：重复注册 -----
    print("【用例 2】同一用户名再次注册，期望 400")
    code, data, raw = post_json(
        "/user/register",
        {"username": USERNAME, "password": PASSWORD},
    )
    assert_case("重复注册应失败", 400, code, raw, "注册失败")
    print()

    # ----- 用例 3：登录成功 -----
    print("【用例 3】正确账号密码登录，期望 200 且返回 token")
    code, data, raw = post_json(
        "/user/login",
        {"username": USERNAME, "password": PASSWORD},
    )
    assert_case("登录成功", 200, code, raw, "token")

    token = ""
    if isinstance(data, dict):
        token = data.get("token") or ""
    if not token:
        print("[失败] 登录成功但 token 为空")
        fail_count += 1
    else:
        print(f"[通过] 拿到 token（长度={len(token)}）")
        print(f"       token 前缀: {token[:24]}...")
        pass_count += 1
    print()

    # ----- 用例 4：错误密码 -----
    print("【用例 4】密码错误，期望 400")
    code, data, raw = post_json(
        "/user/login",
        {"username": USERNAME, "password": "wrong"},
    )
    assert_case("错误密码应失败", 400, code, raw, "登录失败")
    print()

    # ----- 用例 5：非法 JSON -----
    print("【用例 5】请求体非法 JSON，期望 400 且提示参数错误")
    code, data, raw = post_json("/user/register", "{not-json")
    assert_case("非法 JSON 应失败", 400, code, raw, "参数错误")
    print()

    # ----- 汇总 -----
    print("=" * 50)
    print(f"测试结束：通过={pass_count}  失败={fail_count}")
    print("=" * 50)

    if fail_count > 0:
        print("结果：存在失败用例")
        return 1
    print("结果：全部通过")
    return 0


if __name__ == "__main__":
    sys.exit(main())
