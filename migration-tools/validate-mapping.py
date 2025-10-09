#!/usr/bin/env python3
"""
Валидирует полноту маппинга и находит потенциальные пропуски.
"""
import json
import os
import re
import subprocess
from pathlib import Path

REPO_ROOT = Path("/data/hostel-booking-system")
MAPPING_FILE = REPO_ROOT / "migration-tools/naming-map.json"

def load_mapping():
    with open(MAPPING_FILE) as f:
        return json.load(f)

def find_old_references(pattern, paths):
    """Ищет упоминания старых имён в коде."""
    cmd = ["grep", "-r", pattern, "--include=*.go", "--include=*.ts", "--include=*.tsx"] + paths
    result = subprocess.run(cmd, capture_output=True, text=True)
    return result.stdout.splitlines()

def validate_database_tables(mapping):
    """Проверяет упоминания старых таблиц в SQL и Go коде."""
    print("🔍 Checking database table references...")

    old_tables = list(mapping["database_tables"].keys())
    issues = []

    for old_table in old_tables:
        refs = find_old_references(old_table, [
            str(REPO_ROOT / "backend/migrations"),
            str(REPO_ROOT / "backend/internal")
        ])
        if refs:
            issues.append({
                "category": "database",
                "old_name": old_table,
                "new_name": mapping["database_tables"][old_table],
                "references": len(refs),
                "files": list(set([r.split(":")[0] for r in refs]))
            })

    return issues

def validate_go_types(mapping):
    """Проверяет использование старых Go типов."""
    print("🔍 Checking Go type references...")

    old_types = list(mapping["go_types"].keys())
    issues = []

    for old_type in old_types:
        pattern = f"\\b{old_type}\\b"
        refs = find_old_references(pattern, [str(REPO_ROOT / "backend")])
        if refs:
            issues.append({
                "category": "go_types",
                "old_name": old_type,
                "new_name": mapping["go_types"][old_type],
                "references": len(refs)
            })

    return issues

def main():
    print("=" * 60)
    print("🧪 MIGRATION MAPPING VALIDATOR")
    print("=" * 60)

    mapping = load_mapping()

    all_issues = []
    all_issues.extend(validate_database_tables(mapping))
    all_issues.extend(validate_go_types(mapping))

    if not all_issues:
        print("✅ No old references found - mapping is complete!")
        return 0

    print(f"\n⚠️  Found {len(all_issues)} categories with old references:")
    for issue in all_issues:
        print(f"\n  📌 {issue['old_name']} → {issue['new_name']}")
        print(f"     References: {issue['references']}")
        if 'files' in issue:
            print(f"     Files affected: {len(issue['files'])}")

    return 1

if __name__ == "__main__":
    exit(main())
