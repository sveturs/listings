#!/usr/bin/env python3
"""
Memory Leak Detection Tool

Analyzes memory monitoring CSV data to detect potential leaks.
"""

import csv
import sys
from typing import Dict, List, Tuple


def detect_memory_leak(csv_file: str, threshold_mb: float = 50) -> bool:
    """
    Detect memory leaks by analyzing memory growth over time.

    Args:
        csv_file: Path to CSV file with memory monitoring data
        threshold_mb: Minimum heap growth (MB) to consider as leak

    Returns:
        True if leak detected, False otherwise
    """

    # Read CSV data
    try:
        with open(csv_file) as f:
            reader = csv.DictReader(f)
            data = list(reader)
    except FileNotFoundError:
        print(f"❌ File not found: {csv_file}")
        return False
    except Exception as e:
        print(f"❌ Error reading CSV: {e}")
        return False

    if len(data) < 10:
        print("❌ Insufficient data points (need at least 10)")
        print(f"   Found: {len(data)} data points")
        return False

    # Parse metrics
    try:
        heap_start = float(data[0]['heap_alloc_mb'])
        heap_end = float(data[-1]['heap_alloc_mb'])

        goroutines_start = int(float(data[0]['goroutines']))
        goroutines_end = int(float(data[-1]['goroutines']))

        db_conn_start = int(float(data[0]['db_connections_open']))
        db_conn_end = int(float(data[-1]['db_connections_open']))

        db_inuse_start = int(float(data[0]['db_connections_in_use']))
        db_inuse_end = int(float(data[-1]['db_connections_in_use']))

        gc_start = int(float(data[0]['num_gc']))
        gc_end = int(float(data[-1]['num_gc']))

    except (KeyError, ValueError) as e:
        print(f"❌ Error parsing data: {e}")
        return False

    # Calculate growth
    heap_growth = heap_end - heap_start
    goroutine_growth = goroutines_end - goroutines_start
    db_conn_growth = db_conn_end - db_conn_start

    # Calculate duration
    duration_seconds = int(data[-1]['timestamp']) - int(data[0]['timestamp'])
    duration_minutes = duration_seconds / 60

    # Calculate growth rates
    heap_growth_rate = heap_growth / duration_minutes if duration_minutes > 0 else 0
    goroutine_growth_rate = goroutine_growth / duration_minutes if duration_minutes > 0 else 0

    # Calculate max values
    heap_max = max(float(row['heap_alloc_mb']) for row in data)
    goroutines_max = max(int(float(row['goroutines'])) for row in data)

    # Calculate GC activity
    gc_count = gc_end - gc_start
    gc_rate = gc_count / duration_minutes if duration_minutes > 0 else 0

    # Print report
    print("=" * 60)
    print("Memory Leak Detection Report")
    print("=" * 60)
    print(f"Duration: {duration_minutes:.1f} minutes ({duration_seconds}s)")
    print(f"Data points: {len(data)}")
    print()

    # Heap analysis
    print("Heap Allocation:")
    print(f"  Start:        {heap_start:7.1f} MB")
    print(f"  End:          {heap_end:7.1f} MB")
    print(f"  Max:          {heap_max:7.1f} MB")
    print(f"  Growth:       {heap_growth:+7.1f} MB ({heap_growth/heap_start*100:+.1f}%)")
    print(f"  Growth rate:  {heap_growth_rate:7.2f} MB/min")

    leak_detected = False

    if heap_growth > threshold_mb and heap_growth_rate > 1.0:
        print(f"  ❌ LEAK DETECTED: Heap growing at {heap_growth_rate:.2f} MB/min")
        leak_detected = True
    elif heap_growth > threshold_mb:
        print(f"  ⚠️  WARNING: Significant heap growth but slow rate")
    else:
        print(f"  ✅ No heap leak detected")

    print()

    # Goroutine analysis
    print("Goroutines:")
    print(f"  Start:        {goroutines_start:5d}")
    print(f"  End:          {goroutines_end:5d}")
    print(f"  Max:          {goroutines_max:5d}")
    print(f"  Growth:       {goroutine_growth:+5d}")
    print(f"  Growth rate:  {goroutine_growth_rate:7.2f} goroutines/min")

    if goroutine_growth > 100:
        print(f"  ❌ LEAK DETECTED: {goroutine_growth} goroutines not terminated")
        leak_detected = True
    elif goroutine_growth > 50:
        print(f"  ⚠️  WARNING: Moderate goroutine growth")
    else:
        print(f"  ✅ No goroutine leak detected")

    print()

    # DB connection analysis
    print("DB Connections:")
    print(f"  Open (start): {db_conn_start:3d}")
    print(f"  Open (end):   {db_conn_end:3d}")
    print(f"  Growth:       {db_conn_growth:+3d}")
    print(f"  In-use (end): {db_inuse_end:3d}")

    if db_conn_growth > 10:
        print(f"  ❌ LEAK DETECTED: {db_conn_growth} connections not released")
        leak_detected = True
    elif db_conn_growth > 5:
        print(f"  ⚠️  WARNING: Moderate connection growth")
    else:
        print(f"  ✅ No connection leak detected")

    print()

    # GC analysis
    print("Garbage Collection:")
    print(f"  GC count:     {gc_count}")
    print(f"  GC rate:      {gc_rate:.2f} GC/min")

    if gc_rate > 10:
        print(f"  ⚠️  High GC activity (potential memory pressure)")
    else:
        print(f"  ✅ Normal GC activity")

    print()
    print("=" * 60)

    if leak_detected:
        print("🔴 MEMORY LEAK DETECTED!")
        print()
        print("Action Required:")
        print("  1. Run: ./scripts/profile_memory.sh")
        print("  2. Review pprof profiles in /tmp/memory_profiles_*")
        print("  3. Check for unclosed resources (DB, HTTP, contexts)")
        print("  4. Review goroutine creation patterns")
        print("  5. Analyze heap diff: go tool pprof -base baseline.pprof postload.pprof")
        print()
        print("Common leak patterns:")
        print("  - Missing defer rows.Close() after DB queries")
        print("  - Goroutines without exit conditions")
        print("  - Contexts not cancelled (missing defer cancel())")
        print("  - Unbounded caches or slices")
        print("  - HTTP clients without timeouts")
    else:
        print("✅ NO MEMORY LEAKS DETECTED")
        print()
        print("Service memory usage is stable.")
        print(f"  - Heap grew {heap_growth:.1f} MB over {duration_minutes:.1f} minutes")
        print(f"  - {goroutine_growth:+d} goroutine change")
        print(f"  - {db_conn_growth:+d} connection change")

    print("=" * 60)
    return leak_detected


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: detect_leaks.py <csv_file> [threshold_mb]")
        print()
        print("Example:")
        print("  ./scripts/detect_leaks.py /tmp/memory_monitoring_20250104_120000.csv")
        print("  ./scripts/detect_leaks.py /tmp/memory_monitoring_20250104_120000.csv 100")
        sys.exit(1)

    csv_file = sys.argv[1]
    threshold = float(sys.argv[2]) if len(sys.argv) > 2 else 50.0

    leak_found = detect_memory_leak(csv_file, threshold)
    sys.exit(1 if leak_found else 0)
