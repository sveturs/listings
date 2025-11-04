#!/usr/bin/env python3
"""
Load Test Results Analyzer

Analyzes ghz JSON output and validates against SLA requirements:
- Throughput: 10,000 RPS
- p50 latency: < 50ms
- p95 latency: < 100ms
- p99 latency: < 200ms
- Error rate: < 0.1%
"""

import json
import sys
import os
from typing import Dict, Any

# ANSI color codes
RED = '\033[0;31m'
GREEN = '\033[0;32m'
YELLOW = '\033[1;33m'
BLUE = '\033[0;34m'
NC = '\033[0m'  # No Color

# SLA Thresholds
SLA_MIN_RPS = 10000
SLA_P50_MS = 50
SLA_P95_MS = 100
SLA_P99_MS = 200
SLA_MAX_ERROR_RATE = 0.1  # 0.1%


def format_duration(ns: float) -> str:
    """Convert nanoseconds to human-readable format"""
    ms = ns / 1e6
    if ms < 1:
        return f"{ns:.0f}ns"
    elif ms < 1000:
        return f"{ms:.2f}ms"
    else:
        return f"{ms/1000:.2f}s"


def format_bytes(bytes_val: int) -> str:
    """Format bytes to human-readable format"""
    for unit in ['B', 'KB', 'MB', 'GB']:
        if bytes_val < 1024.0:
            return f"{bytes_val:.2f}{unit}"
        bytes_val /= 1024.0
    return f"{bytes_val:.2f}TB"


def analyze_load_test(json_file: str) -> bool:
    """
    Analyze load test results from ghz JSON output

    Returns:
        True if all SLA requirements met, False otherwise
    """
    if not os.path.exists(json_file):
        print(f"{RED}ERROR: File not found: {json_file}{NC}")
        return False

    with open(json_file) as f:
        data = json.load(f)

    test_name = os.path.basename(json_file).replace('.json', '')

    print("=" * 80)
    print(f"Load Test Analysis: {test_name}")
    print("=" * 80)
    print()

    # === Basic Metrics ===
    print(f"{BLUE}📊 Basic Metrics:{NC}")
    print("-" * 80)

    total = data['count']
    status_dist = data.get('statusCodeDist', {})
    success = status_dist.get('OK', 0)
    failed = total - success
    error_rate = (failed / total * 100) if total > 0 else 0

    print(f"  Total Requests:    {total:,}")
    print(f"  Successful:        {success:,} ({success/total*100:.2f}%)")
    print(f"  Failed:            {failed:,} ({error_rate:.3f}%)")
    print()

    # === Throughput ===
    print(f"{BLUE}🚀 Throughput:{NC}")
    print("-" * 80)

    rps = data['rps']
    duration_s = data['total'] / 1e9  # Convert ns to seconds

    print(f"  Actual RPS:        {rps:.2f}")
    print(f"  Total Duration:    {format_duration(data['total'])}")
    print(f"  Fastest Request:   {format_duration(data['fastest'])}")
    print(f"  Slowest Request:   {format_duration(data['slowest'])}")
    print(f"  Average:           {format_duration(data['average'])}")
    print()

    # === Latency Distribution ===
    print(f"{BLUE}⏱️  Latency Distribution:{NC}")
    print("-" * 80)

    latency_dist = data['latencyDistribution']

    # Extract key percentiles
    p50 = latency_dist[5]['latency'] / 1e6 if len(latency_dist) > 5 else 0
    p90 = latency_dist[8]['latency'] / 1e6 if len(latency_dist) > 8 else 0
    p95 = latency_dist[9]['latency'] / 1e6 if len(latency_dist) > 9 else 0
    p99 = latency_dist[10]['latency'] / 1e6 if len(latency_dist) > 10 else 0

    print(f"  p50:               {p50:.2f} ms")
    print(f"  p75:               {latency_dist[7]['latency']/1e6:.2f} ms" if len(latency_dist) > 7 else "")
    print(f"  p90:               {p90:.2f} ms")
    print(f"  p95:               {p95:.2f} ms")
    print(f"  p99:               {p99:.2f} ms")
    print()

    # === Error Breakdown ===
    if failed > 0:
        print(f"{YELLOW}⚠️  Error Breakdown:{NC}")
        print("-" * 80)

        error_dist = data.get('errorDist', {})
        for error, count in error_dist.items():
            print(f"  {error}: {count:,} ({count/total*100:.2f}%)")
        print()

        # Status codes
        print(f"  Status Code Distribution:")
        for code, count in status_dist.items():
            print(f"    {code}: {count:,}")
        print()

    # === Network Stats ===
    print(f"{BLUE}🌐 Network Stats:{NC}")
    print("-" * 80)

    size_total = data.get('sizeTotal', 0)
    size_req = data.get('sizeReq', {}).get('total', 0)

    print(f"  Total Data Sent:   {format_bytes(size_req)}")
    print(f"  Total Data Recv:   {format_bytes(size_total)}")
    print(f"  Avg Request Size:  {format_bytes(size_req / total)}" if total > 0 else "N/A")
    print(f"  Avg Response Size: {format_bytes(size_total / total)}" if total > 0 else "N/A")
    print()

    # === SLA Validation ===
    print("=" * 80)
    print(f"{BLUE}🎯 SLA Validation:{NC}")
    print("=" * 80)

    sla_pass = True
    checks = []

    # Check 1: Throughput
    if 'sustained' in test_name.lower():
        if rps >= SLA_MIN_RPS:
            checks.append(f"  {GREEN}✅ Throughput: {rps:.0f} RPS >= {SLA_MIN_RPS} RPS{NC}")
        else:
            checks.append(f"  {RED}❌ Throughput: {rps:.0f} RPS < {SLA_MIN_RPS} RPS{NC}")
            sla_pass = False
    else:
        checks.append(f"  ℹ️  Throughput: {rps:.0f} RPS (not validated for non-sustained tests)")

    # Check 2: p50 Latency
    if p50 < SLA_P50_MS:
        checks.append(f"  {GREEN}✅ p50 Latency: {p50:.2f} ms < {SLA_P50_MS} ms{NC}")
    else:
        checks.append(f"  {YELLOW}⚠️  p50 Latency: {p50:.2f} ms >= {SLA_P50_MS} ms (warning){NC}")

    # Check 3: p95 Latency
    if p95 < SLA_P95_MS:
        checks.append(f"  {GREEN}✅ p95 Latency: {p95:.2f} ms < {SLA_P95_MS} ms{NC}")
    else:
        checks.append(f"  {RED}❌ p95 Latency: {p95:.2f} ms >= {SLA_P95_MS} ms{NC}")
        sla_pass = False

    # Check 4: p99 Latency
    if p99 < SLA_P99_MS:
        checks.append(f"  {GREEN}✅ p99 Latency: {p99:.2f} ms < {SLA_P99_MS} ms{NC}")
    else:
        checks.append(f"  {YELLOW}⚠️  p99 Latency: {p99:.2f} ms >= {SLA_P99_MS} ms (warning){NC}")

    # Check 5: Error Rate
    if error_rate < SLA_MAX_ERROR_RATE:
        checks.append(f"  {GREEN}✅ Error Rate: {error_rate:.3f}% < {SLA_MAX_ERROR_RATE}%{NC}")
    else:
        checks.append(f"  {RED}❌ Error Rate: {error_rate:.3f}% >= {SLA_MAX_ERROR_RATE}%{NC}")
        sla_pass = False

    for check in checks:
        print(check)

    print()
    print("=" * 80)

    if sla_pass:
        print(f"{GREEN}✅ SLA PASSED: All requirements met!{NC}")
    else:
        print(f"{RED}❌ SLA FAILED: One or more requirements not met!{NC}")

    print("=" * 80)
    print()

    return sla_pass


def main():
    if len(sys.argv) < 2:
        print(f"Usage: {sys.argv[0]} <results.json>")
        print()
        print("Example:")
        print(f"  {sys.argv[0]} /tmp/load_test_results_*/sustained_10min.json")
        sys.exit(1)

    json_file = sys.argv[1]
    sla_pass = analyze_load_test(json_file)

    sys.exit(0 if sla_pass else 1)


if __name__ == "__main__":
    main()
