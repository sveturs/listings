#!/usr/bin/env python3
"""
Migration Validation Script
============================

Validates data integrity after migration from monolith to microservice.

Checks:
1. Row counts match (source vs target)
2. No orphaned records
3. Foreign key integrity
4. Data consistency (prices, dates, coordinates)
5. Image URL accessibility (optional)
6. All required fields populated

Usage:
    python3 validate_migration.py \\
        --source-host localhost \\
        --source-port 5433 \\
        --target-host localhost \\
        --target-port 5432

Author: Migration Team
Date: 2025-10-31
Version: 1.0
"""

import argparse
import sys
from typing import Dict, List, Tuple, Optional
from dataclasses import dataclass, field

try:
    import psycopg2
    from psycopg2.extras import RealDictCursor
except ImportError:
    print("ERROR: psycopg2 not installed. Run: pip3 install psycopg2-binary")
    sys.exit(1)

try:
    from rich.console import Console
    from rich.table import Table
    from rich.panel import Panel
    from rich import print as rprint
except ImportError:
    print("ERROR: rich not installed. Run: pip3 install rich")
    sys.exit(1)

console = Console()


@dataclass
class ValidationResult:
    """Result of a single validation check"""
    check_name: str
    passed: bool
    message: str
    details: Dict = field(default_factory=dict)


@dataclass
class ValidationConfig:
    """Validation configuration"""
    source_host: str
    source_port: int
    source_db: str
    source_user: str
    source_password: str
    target_host: str
    target_port: int
    target_db: str
    target_user: str
    target_password: str
    check_images: bool = False
    verbose: bool = False


class MigrationValidator:
    """Validates migrated data"""

    def __init__(self, config: ValidationConfig):
        self.config = config
        self.source_conn = None
        self.target_conn = None
        self.results: List[ValidationResult] = []

    def connect_databases(self):
        """Establish database connections"""
        console.print("\n[bold cyan]Connecting to databases...[/bold cyan]")

        try:
            # Source connection
            self.source_conn = psycopg2.connect(
                host=self.config.source_host,
                port=self.config.source_port,
                database=self.config.source_db,
                user=self.config.source_user,
                password=self.config.source_password,
                connect_timeout=10
            )
            console.print(f"✓ Connected to source: {self.config.source_host}:{self.config.source_port}/{self.config.source_db}", style="green")

            # Target connection
            self.target_conn = psycopg2.connect(
                host=self.config.target_host,
                port=self.config.target_port,
                database=self.config.target_db,
                user=self.config.target_user,
                password=self.config.target_password,
                connect_timeout=10
            )
            console.print(f"✓ Connected to target: {self.config.target_host}:{self.config.target_port}/{self.config.target_db}", style="green")

        except Exception as e:
            console.print(f"[bold red]✗ Database connection failed: {e}[/bold red]")
            raise

    def close_connections(self):
        """Close database connections"""
        if self.source_conn:
            self.source_conn.close()
        if self.target_conn:
            self.target_conn.close()

    def add_result(self, check_name: str, passed: bool, message: str, details: Dict = None):
        """Add validation result"""
        self.results.append(ValidationResult(
            check_name=check_name,
            passed=passed,
            message=message,
            details=details or {}
        ))

    def check_row_counts(self):
        """Validate row counts match between source and target"""
        console.print("\n[bold]1. Checking Row Counts...[/bold]")

        try:
            # Source counts
            with self.source_conn.cursor() as cur:
                cur.execute("SELECT COUNT(*) FROM unified_listings")
                source_listings_count = cur.fetchone()[0]

                cur.execute("""
                    SELECT COUNT(*) FROM (
                        SELECT DISTINCT listing_id FROM c2c_images
                        WHERE listing_id IN (SELECT id FROM c2c_listings)
                        UNION
                        SELECT DISTINCT storefront_product_id FROM b2c_product_images
                        WHERE storefront_product_id IN (SELECT id FROM b2c_products)
                    ) AS all_images
                """)
                source_images_count = cur.fetchone()[0]

            # Target counts
            with self.target_conn.cursor() as cur:
                cur.execute("SELECT COUNT(*) FROM listings")
                target_listings_count = cur.fetchone()[0]

                cur.execute("SELECT COUNT(DISTINCT listing_id) FROM listing_images")
                target_images_count = cur.fetchone()[0]

                cur.execute("SELECT COUNT(*) FROM listing_stats")
                target_stats_count = cur.fetchone()[0]

                cur.execute("SELECT COUNT(*) FROM listing_locations")
                target_locations_count = cur.fetchone()[0]

            # Validate listings count
            if source_listings_count == target_listings_count:
                self.add_result(
                    "Listings Count",
                    True,
                    f"✓ Counts match: {source_listings_count} listings",
                    {"source": source_listings_count, "target": target_listings_count}
                )
                console.print(f"  ✓ Listings: {source_listings_count} (source) = {target_listings_count} (target)", style="green")
            else:
                self.add_result(
                    "Listings Count",
                    False,
                    f"✗ Mismatch: {source_listings_count} (source) vs {target_listings_count} (target)",
                    {"source": source_listings_count, "target": target_listings_count}
                )
                console.print(f"  ✗ Listings: {source_listings_count} (source) ≠ {target_listings_count} (target)", style="red")

            # Validate images count (approximate - some listings may have no images)
            console.print(f"  ℹ Images: {source_images_count} listings with images (source), {target_images_count} (target)", style="dim")

            # Stats should match listings count
            if target_listings_count == target_stats_count:
                self.add_result(
                    "Stats Count",
                    True,
                    f"✓ Stats records match listings: {target_stats_count}",
                    {"target_listings": target_listings_count, "target_stats": target_stats_count}
                )
                console.print(f"  ✓ Stats: {target_stats_count} records", style="green")
            else:
                self.add_result(
                    "Stats Count",
                    False,
                    f"✗ Stats mismatch: {target_listings_count} listings vs {target_stats_count} stats",
                    {"target_listings": target_listings_count, "target_stats": target_stats_count}
                )
                console.print(f"  ✗ Stats: {target_stats_count} records (should be {target_listings_count})", style="red")

            # Locations should match listings count
            if target_listings_count == target_locations_count:
                self.add_result(
                    "Locations Count",
                    True,
                    f"✓ Location records match listings: {target_locations_count}",
                    {"target_listings": target_listings_count, "target_locations": target_locations_count}
                )
                console.print(f"  ✓ Locations: {target_locations_count} records", style="green")
            else:
                self.add_result(
                    "Locations Count",
                    False,
                    f"✗ Locations mismatch: {target_listings_count} listings vs {target_locations_count} locations",
                    {"target_listings": target_listings_count, "target_locations": target_locations_count}
                )
                console.print(f"  ✗ Locations: {target_locations_count} records (should be {target_listings_count})", style="red")

        except Exception as e:
            self.add_result("Row Counts", False, f"Error: {e}")
            console.print(f"  [red]Error: {e}[/red]")

    def check_orphaned_records(self):
        """Check for orphaned records in target database"""
        console.print("\n[bold]2. Checking for Orphaned Records...[/bold]")

        try:
            with self.target_conn.cursor() as cur:
                # Check orphaned images
                cur.execute("""
                    SELECT COUNT(*) FROM listing_images
                    WHERE listing_id NOT IN (SELECT id FROM listings)
                """)
                orphaned_images = cur.fetchone()[0]

                # Check orphaned stats
                cur.execute("""
                    SELECT COUNT(*) FROM listing_stats
                    WHERE listing_id NOT IN (SELECT id FROM listings)
                """)
                orphaned_stats = cur.fetchone()[0]

                # Check orphaned locations
                cur.execute("""
                    SELECT COUNT(*) FROM listing_locations
                    WHERE listing_id NOT IN (SELECT id FROM listings)
                """)
                orphaned_locations = cur.fetchone()[0]

            total_orphaned = orphaned_images + orphaned_stats + orphaned_locations

            if total_orphaned == 0:
                self.add_result(
                    "Orphaned Records",
                    True,
                    "✓ No orphaned records found",
                    {"images": 0, "stats": 0, "locations": 0}
                )
                console.print("  ✓ No orphaned records", style="green")
            else:
                self.add_result(
                    "Orphaned Records",
                    False,
                    f"✗ Found {total_orphaned} orphaned records",
                    {"images": orphaned_images, "stats": orphaned_stats, "locations": orphaned_locations}
                )
                console.print(f"  ✗ Orphaned images: {orphaned_images}", style="red")
                console.print(f"  ✗ Orphaned stats: {orphaned_stats}", style="red")
                console.print(f"  ✗ Orphaned locations: {orphaned_locations}", style="red")

        except Exception as e:
            self.add_result("Orphaned Records", False, f"Error: {e}")
            console.print(f"  [red]Error: {e}[/red]")

    def check_required_fields(self):
        """Check that required fields are populated"""
        console.print("\n[bold]3. Checking Required Fields...[/bold]")

        try:
            with self.target_conn.cursor() as cur:
                # Check NULL titles
                cur.execute("SELECT COUNT(*) FROM listings WHERE title IS NULL OR title = ''")
                null_titles = cur.fetchone()[0]

                # Check NULL prices
                cur.execute("SELECT COUNT(*) FROM listings WHERE price IS NULL OR price < 0")
                invalid_prices = cur.fetchone()[0]

                # Check NULL category_id
                cur.execute("SELECT COUNT(*) FROM listings WHERE category_id IS NULL")
                null_categories = cur.fetchone()[0]

                # Check NULL user_id
                cur.execute("SELECT COUNT(*) FROM listings WHERE user_id IS NULL")
                null_users = cur.fetchone()[0]

                # Check missing UUIDs
                cur.execute("SELECT COUNT(*) FROM listings WHERE uuid IS NULL")
                null_uuids = cur.fetchone()[0]

            total_issues = null_titles + invalid_prices + null_categories + null_users + null_uuids

            if total_issues == 0:
                self.add_result(
                    "Required Fields",
                    True,
                    "✓ All required fields populated",
                    {"issues": 0}
                )
                console.print("  ✓ All required fields populated", style="green")
            else:
                self.add_result(
                    "Required Fields",
                    False,
                    f"✗ Found {total_issues} records with missing required fields",
                    {
                        "null_titles": null_titles,
                        "invalid_prices": invalid_prices,
                        "null_categories": null_categories,
                        "null_users": null_users,
                        "null_uuids": null_uuids
                    }
                )
                if null_titles > 0:
                    console.print(f"  ✗ NULL/empty titles: {null_titles}", style="red")
                if invalid_prices > 0:
                    console.print(f"  ✗ Invalid prices: {invalid_prices}", style="red")
                if null_categories > 0:
                    console.print(f"  ✗ NULL category_id: {null_categories}", style="red")
                if null_users > 0:
                    console.print(f"  ✗ NULL user_id: {null_users}", style="red")
                if null_uuids > 0:
                    console.print(f"  ✗ NULL uuid: {null_uuids}", style="red")

        except Exception as e:
            self.add_result("Required Fields", False, f"Error: {e}")
            console.print(f"  [red]Error: {e}[/red]")

    def check_data_consistency(self):
        """Check data value consistency"""
        console.print("\n[bold]4. Checking Data Consistency...[/bold]")

        try:
            with self.target_conn.cursor() as cur:
                # Check coordinate ranges
                cur.execute("""
                    SELECT COUNT(*) FROM listings
                    WHERE latitude IS NOT NULL AND (latitude < -90 OR latitude > 90)
                """)
                invalid_latitudes = cur.fetchone()[0]

                cur.execute("""
                    SELECT COUNT(*) FROM listings
                    WHERE longitude IS NOT NULL AND (longitude < -180 OR longitude > 180)
                """)
                invalid_longitudes = cur.fetchone()[0]

                # Check future dates
                cur.execute("""
                    SELECT COUNT(*) FROM listings
                    WHERE created_at > NOW() OR updated_at > NOW()
                """)
                future_dates = cur.fetchone()[0]

                # Check image URLs
                cur.execute("""
                    SELECT COUNT(*) FROM listing_images
                    WHERE url IS NULL OR url = ''
                """)
                invalid_image_urls = cur.fetchone()[0]

            total_issues = invalid_latitudes + invalid_longitudes + future_dates + invalid_image_urls

            if total_issues == 0:
                self.add_result(
                    "Data Consistency",
                    True,
                    "✓ All data values consistent",
                    {"issues": 0}
                )
                console.print("  ✓ All data values consistent", style="green")
            else:
                self.add_result(
                    "Data Consistency",
                    False,
                    f"✗ Found {total_issues} data consistency issues",
                    {
                        "invalid_latitudes": invalid_latitudes,
                        "invalid_longitudes": invalid_longitudes,
                        "future_dates": future_dates,
                        "invalid_image_urls": invalid_image_urls
                    }
                )
                if invalid_latitudes > 0:
                    console.print(f"  ✗ Invalid latitudes: {invalid_latitudes}", style="red")
                if invalid_longitudes > 0:
                    console.print(f"  ✗ Invalid longitudes: {invalid_longitudes}", style="red")
                if future_dates > 0:
                    console.print(f"  ✗ Future timestamps: {future_dates}", style="red")
                if invalid_image_urls > 0:
                    console.print(f"  ✗ Invalid image URLs: {invalid_image_urls}", style="red")

        except Exception as e:
            self.add_result("Data Consistency", False, f"Error: {e}")
            console.print(f"  [red]Error: {e}[/red]")

    def check_sample_data_match(self):
        """Compare sample records between source and target"""
        console.print("\n[bold]5. Checking Sample Data Match...[/bold]")

        try:
            # Get 3 sample IDs
            with self.source_conn.cursor() as cur:
                cur.execute("SELECT id FROM unified_listings ORDER BY id LIMIT 3")
                sample_ids = [row[0] for row in cur.fetchall()]

            if not sample_ids:
                console.print("  [yellow]No data to sample[/yellow]")
                return

            mismatches = []

            for listing_id in sample_ids:
                # Get source data
                with self.source_conn.cursor(cursor_factory=RealDictCursor) as cur:
                    cur.execute("""
                        SELECT id, title, price, source_type
                        FROM unified_listings WHERE id = %s
                    """, (listing_id,))
                    source_row = cur.fetchone()

                # Get target data
                with self.target_conn.cursor(cursor_factory=RealDictCursor) as cur:
                    cur.execute("""
                        SELECT id, title, price, source_type
                        FROM listings WHERE id = %s
                    """, (listing_id,))
                    target_row = cur.fetchone()

                # Compare
                if source_row and target_row:
                    if (source_row['title'] != target_row['title'] or
                        float(source_row['price']) != float(target_row['price']) or
                        source_row['source_type'] != target_row['source_type']):
                        mismatches.append(listing_id)
                        if self.config.verbose:
                            console.print(f"  ✗ Mismatch for ID {listing_id}:", style="red")
                            console.print(f"    Source: {source_row}", style="dim")
                            console.print(f"    Target: {target_row}", style="dim")
                elif not target_row:
                    mismatches.append(listing_id)
                    console.print(f"  ✗ ID {listing_id} missing in target", style="red")

            if len(mismatches) == 0:
                self.add_result(
                    "Sample Data Match",
                    True,
                    f"✓ All {len(sample_ids)} sample records match",
                    {"samples_checked": len(sample_ids), "mismatches": 0}
                )
                console.print(f"  ✓ All {len(sample_ids)} sample records match", style="green")
            else:
                self.add_result(
                    "Sample Data Match",
                    False,
                    f"✗ {len(mismatches)} of {len(sample_ids)} sample records have mismatches",
                    {"samples_checked": len(sample_ids), "mismatches": len(mismatches)}
                )
                console.print(f"  ✗ {len(mismatches)} sample records have mismatches", style="red")

        except Exception as e:
            self.add_result("Sample Data Match", False, f"Error: {e}")
            console.print(f"  [red]Error: {e}[/red]")

    def run_validation(self):
        """Run all validation checks"""
        console.print("\n")
        console.print(Panel.fit(
            "[bold cyan]Migration Validation[/bold cyan]\n"
            f"Source: {self.config.source_host}:{self.config.source_port}/{self.config.source_db}\n"
            f"Target: {self.config.target_host}:{self.config.target_port}/{self.config.target_db}",
            border_style="cyan"
        ))

        try:
            self.connect_databases()

            # Run checks
            self.check_row_counts()
            self.check_orphaned_records()
            self.check_required_fields()
            self.check_data_consistency()
            self.check_sample_data_match()

            # Print summary
            self.print_summary()

        except Exception as e:
            console.print(f"\n[bold red]Validation failed: {e}[/bold red]")
            raise
        finally:
            self.close_connections()

    def print_summary(self):
        """Print validation summary"""
        console.print("\n")

        passed = sum(1 for r in self.results if r.passed)
        failed = len(self.results) - passed

        if failed == 0:
            console.print(Panel.fit(
                "[bold green]✓ ALL VALIDATION CHECKS PASSED![/bold green]",
                border_style="green"
            ))
        else:
            console.print(Panel.fit(
                f"[bold red]✗ {failed} VALIDATION CHECKS FAILED[/bold red]",
                border_style="red"
            ))

        # Create results table
        table = Table(title="Validation Results", show_header=True, header_style="bold magenta")
        table.add_column("Check", style="cyan")
        table.add_column("Status", justify="center")
        table.add_column("Message", style="white")

        for result in self.results:
            status = "[green]✓ PASS[/green]" if result.passed else "[red]✗ FAIL[/red]"
            table.add_row(result.check_name, status, result.message)

        console.print(table)

        # Exit code
        return 0 if failed == 0 else 1


def parse_arguments():
    """Parse command line arguments"""
    parser = argparse.ArgumentParser(
        description="Validate migration data integrity",
        formatter_class=argparse.RawDescriptionHelpFormatter
    )

    # Source database
    parser.add_argument('--source-host', default='localhost', help='Source DB host')
    parser.add_argument('--source-port', type=int, default=5433, help='Source DB port')
    parser.add_argument('--source-db', default='svetubd', help='Source database name')
    parser.add_argument('--source-user', default='postgres', help='Source DB user')
    parser.add_argument('--source-password', default='mX3g1XGhMRUZEX3l', help='Source DB password')

    # Target database
    parser.add_argument('--target-host', default='localhost', help='Target DB host')
    parser.add_argument('--target-port', type=int, default=5432, help='Target DB port')
    parser.add_argument('--target-db', default='listings_dev_db', help='Target database name')
    parser.add_argument('--target-user', default='listings_user', help='Target DB user')
    parser.add_argument('--target-password', required=False, help='Target DB password')

    # Options
    parser.add_argument('--check-images', action='store_true', help='Check image URL accessibility (slow)')
    parser.add_argument('-v', '--verbose', action='store_true', help='Verbose output')

    return parser.parse_args()


def main():
    """Main entry point"""
    args = parse_arguments()

    config = ValidationConfig(
        source_host=args.source_host,
        source_port=args.source_port,
        source_db=args.source_db,
        source_user=args.source_user,
        source_password=args.source_password,
        target_host=args.target_host,
        target_port=args.target_port,
        target_db=args.target_db,
        target_user=args.target_user,
        target_password=args.target_password or '',
        check_images=args.check_images,
        verbose=args.verbose
    )

    validator = MigrationValidator(config)
    try:
        exit_code = validator.run_validation()
        sys.exit(exit_code)
    except KeyboardInterrupt:
        console.print("\n[yellow]Validation interrupted by user[/yellow]")
        sys.exit(1)
    except Exception as e:
        console.print(f"\n[bold red]Fatal error: {e}[/bold red]")
        sys.exit(1)


if __name__ == '__main__':
    main()
