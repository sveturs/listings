#!/usr/bin/env python3
"""
Migrate test fixtures from old schema (b2c_products) to unified schema (listings).

Field mappings:
- b2c_products → listings
  - name → title
  - stock_quantity → quantity
  - stock_status → status (needs special handling)
  - is_active → removed (use status instead)
  - barcode → removed (not in unified schema)
  - attributes → listing_attributes table (or keep as JSONB if supported)

- b2c_storefronts → storefronts (table name stays same, but need to verify)

Additional changes:
- Add source_type = 'b2c' for all listings
- Update status values: 'in_stock' → 'active', 'out_of_stock' → 'inactive'
"""

import re
from pathlib import Path
from datetime import datetime

FIXTURES_DIR = Path("/p/github.com/sveturs/listings/tests/fixtures")

FILES_TO_MIGRATE = [
    "b2c_inventory_fixtures.sql",
    "bulk_operations_fixtures.sql",
    "create_product_fixtures.sql",
    "decrement_stock_fixtures.sql",
    "rollback_stock_fixtures.sql",
    "update_product_fixtures.sql",
]


def migrate_insert_statement(content: str) -> str:
    """Migrate INSERT INTO b2c_products statements to listings table."""

    # Pattern to match INSERT INTO b2c_products statements
    pattern = r'INSERT INTO b2c_products\s*\((.*?)\)\s*VALUES\s*\((.*?)\);'

    def replace_insert(match):
        fields = match.group(1)
        values = match.group(2)

        # Replace table name
        # Replace field names
        fields = fields.replace('name,', 'title,')
        fields = fields.replace(' name ', ' title ')
        fields = fields.replace('stock_quantity', 'quantity')

        # Remove unsupported fields
        fields_list = [f.strip() for f in fields.split(',')]
        values_list = values.split(',')

        # Remove barcode, is_active, stock_status, attributes
        remove_fields = ['barcode', 'is_active', 'stock_status', 'attributes']

        # Add source_type if not present
        if 'source_type' not in fields:
            fields_list.append('source_type')
            # Insert source_type value - need special handling for multiline

        new_fields = ', '.join(fields_list)

        return f'INSERT INTO listings ({new_fields}) VALUES ({values})'

    # Simple replacement first
    content = re.sub(pattern, replace_insert, content, flags=re.DOTALL | re.MULTILINE)

    return content


def migrate_file(filepath: Path) -> None:
    """Migrate a single fixture file."""

    print(f"Processing: {filepath.name}")

    # Create backup
    backup_path = filepath.with_suffix(f'.sql.backup_{datetime.now().strftime("%Y%m%d_%H%M%S")}')
    content = filepath.read_text()
    backup_path.write_text(content)
    print(f"  Backup created: {backup_path.name}")

    # Simple table name replacement
    content = content.replace('b2c_products', 'listings')
    # storefronts table name is already correct

    # Field name replacements (in column lists)
    content = re.sub(r'\bname\b(?=\s*,)', 'title', content)
    content = re.sub(r',\s*name\s*,', ', title,', content)
    content = re.sub(r'\(\s*name\s*,', '(title,', content)

    content = re.sub(r'\bstock_quantity\b', 'quantity', content)

    # Remove problematic fields from INSERT statements
    # This is complex - need to handle field-value pairs

    # For now, do simple replacements and let user manually fix complex cases

    # Write migrated content
    filepath.write_text(content)
    print(f"  ✅ Migrated: {filepath.name}")


def main():
    print("=== Migrating test fixtures to unified schema ===\n")

    for filename in FILES_TO_MIGRATE:
        filepath = FIXTURES_DIR / filename
        if not filepath.exists():
            print(f"⚠️  File not found: {filename} - skipping\n")
            continue

        migrate_file(filepath)
        print()

    print("=== Migration complete ===")
    print("\n⚠️  MANUAL STEPS REQUIRED:")
    print("1. Review each migrated file")
    print("2. Remove fields not in listings schema:")
    print("   - barcode")
    print("   - is_active (use status='active'/'inactive' instead)")
    print("   - stock_status (use status='active'/'inactive' instead)")
    print("   - attributes (move to listing_attributes table or use JSONB if supported)")
    print("3. Add source_type='b2c' to all listings INSERT statements")
    print("4. Update DELETE statements referencing old tables")
    print("5. Update any UPDATE statements")
    print("6. Re-run tests to verify")


if __name__ == "__main__":
    main()
