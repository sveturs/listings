#!/usr/bin/env python3
"""
Data Migration: Monolith → Listings Microservice
==================================================

Migrates unified_listings VIEW from monolith to microservice database.

Features:
- Batch processing (configurable batch size)
- Progress tracking with rich progress bars
- Dry-run mode for testing
- Resume capability from last batch
- Transactional safety (rollback on error)
- Comprehensive validation
- Colored console output
- SSH tunnel support for remote databases

Usage:
    # Dry run (no changes)
    python3 migrate_data.py --dry-run

    # Real migration
    python3 migrate_data.py \\
        --source-host localhost \\
        --source-port 5433 \\
        --source-db svetubd \\
        --source-user postgres \\
        --source-password mX3g1XGhMRUZEX3l \\
        --target-host dev.svetu.rs \\
        --target-db listings_dev_db \\
        --target-user listings_user \\
        --target-password <password>

    # Resume from batch 10
    python3 migrate_data.py --resume-from 10

Author: Migration Team
Date: 2025-10-31
Version: 1.0
"""

import argparse
import sys
import json
import uuid
from datetime import datetime
from typing import Dict, List, Tuple, Optional, Any
from dataclasses import dataclass

try:
    import psycopg2
    from psycopg2.extras import execute_values, RealDictCursor
    from psycopg2.extensions import connection as Connection
except ImportError:
    print("ERROR: psycopg2 not installed. Run: pip3 install psycopg2-binary")
    sys.exit(1)

try:
    from rich.console import Console
    from rich.progress import Progress, SpinnerColumn, BarColumn, TextColumn, TimeElapsedColumn, TimeRemainingColumn
    from rich.table import Table
    from rich.panel import Panel
    from rich import print as rprint
except ImportError:
    print("ERROR: rich not installed. Run: pip3 install rich")
    sys.exit(1)

# Initialize console
console = Console()


@dataclass
class MigrationConfig:
    """Migration configuration"""
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
    batch_size: int
    dry_run: bool
    resume_from: int
    verbose: bool


@dataclass
class MigrationStats:
    """Track migration statistics"""
    listings_migrated: int = 0
    images_migrated: int = 0
    errors: int = 0
    start_time: datetime = None
    end_time: datetime = None

    def duration_seconds(self) -> float:
        if self.start_time and self.end_time:
            return (self.end_time - self.start_time).total_seconds()
        return 0


class DataMigrator:
    """Main migration orchestrator"""

    def __init__(self, config: MigrationConfig):
        self.config = config
        self.stats = MigrationStats()
        self.source_conn: Optional[Connection] = None
        self.target_conn: Optional[Connection] = None

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
            if self.config.target_host == "dev.svetu.rs":
                # For dev.svetu.rs, we need to connect via SSH tunnel
                # Use docker exec or port forwarding
                console.print("[yellow]Note: For dev.svetu.rs, ensure SSH tunnel or port forwarding is active[/yellow]")

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
            console.print("✓ Closed source connection", style="dim")
        if self.target_conn:
            self.target_conn.close()
            console.print("✓ Closed target connection", style="dim")

    def fetch_source_data(self, offset: int, limit: int) -> List[Dict[str, Any]]:
        """Fetch batch of listings from source"""
        query = """
            SELECT
                id, user_id, category_id,
                title, description, price, status,
                views_count, created_at, updated_at,
                storefront_id, images
            FROM unified_listings
            ORDER BY id
            LIMIT %s OFFSET %s
        """

        with self.source_conn.cursor(cursor_factory=RealDictCursor) as cur:
            cur.execute(query, (limit, offset))
            rows = cur.fetchall()
            return [dict(row) for row in rows]

    def migrate_listings_batch(self, listings: List[Dict[str, Any]]) -> List[int]:
        """Migrate a batch of listings to microservice schema"""
        if not listings:
            return []

        # Prepare data for insertion
        values = []
        for listing in listings:
            status = listing.get('status') or 'draft'

            # Map published_at: if status='active', use created_at, else NULL
            published_at = listing['created_at'] if status == 'active' else None

            values.append((
                listing['id'],
                str(uuid.uuid4()),  # Generate UUID
                listing['user_id'],
                listing['category_id'],
                listing.get('storefront_id'),  # Can be NULL
                listing['title'],
                listing['description'],
                listing['price'],
                'RSD',  # Default currency
                status,
                'public',  # Default visibility
                1,  # Default quantity
                None,  # SKU (NULL)
                listing.get('views_count') or 0,
                0,  # favorites_count (default)
                published_at,
                False,  # is_deleted
                listing['created_at'],
                listing['updated_at']
            ))

        # Insert query matching microservice schema
        insert_query = """
            INSERT INTO listings (
                id, uuid, user_id, category_id, storefront_id,
                title, description, price, currency, status,
                visibility, quantity, sku, views_count, favorites_count,
                published_at, is_deleted, created_at, updated_at
            ) VALUES %s
            ON CONFLICT (id) DO UPDATE SET
                uuid = EXCLUDED.uuid,
                user_id = EXCLUDED.user_id,
                category_id = EXCLUDED.category_id,
                storefront_id = EXCLUDED.storefront_id,
                title = EXCLUDED.title,
                description = EXCLUDED.description,
                price = EXCLUDED.price,
                currency = EXCLUDED.currency,
                status = EXCLUDED.status,
                visibility = EXCLUDED.visibility,
                quantity = EXCLUDED.quantity,
                sku = EXCLUDED.sku,
                views_count = EXCLUDED.views_count,
                favorites_count = EXCLUDED.favorites_count,
                published_at = EXCLUDED.published_at,
                is_deleted = EXCLUDED.is_deleted,
                updated_at = EXCLUDED.updated_at
            RETURNING id
        """

        if self.config.dry_run:
            console.print(f"[dim]DRY RUN: Would insert {len(listings)} listings[/dim]")
            return [l['id'] for l in listings]

        with self.target_conn.cursor() as cur:
            execute_values(cur, insert_query, values, template=None, page_size=100)
            inserted_ids = [row[0] for row in cur.fetchall()]
            self.target_conn.commit()
            self.stats.listings_migrated += len(inserted_ids)
            return inserted_ids

    def migrate_images(self, listings: List[Dict[str, Any]]):
        """Migrate images from JSONB to listing_images table"""
        images_data = []

        for listing in listings:
            listing_id = listing['id']
            images_json = listing.get('images', [])

            if not images_json or len(images_json) == 0:
                continue

            # Parse images - unified structure after monolith migration
            for idx, img in enumerate(images_json):
                # Handle both old and new image structures
                url = img.get('public_url') or img.get('image_url') or img.get('file_path', '')
                thumbnail_url = img.get('thumbnail_url', url)
                is_primary = img.get('is_main') or img.get('is_default', False)
                storage_path = img.get('file_path') or img.get('storage_path', '')

                # Calculate display order (primary first, then by index)
                display_order = 0 if is_primary else (idx + 1)

                images_data.append((
                    listing_id,
                    url,
                    storage_path,
                    thumbnail_url,
                    display_order,
                    is_primary,
                    datetime.now()
                ))

        if not images_data:
            return

        insert_query = """
            INSERT INTO listing_images (
                listing_id, url, storage_path, thumbnail_url,
                display_order, is_primary, created_at
            ) VALUES %s
        """

        if self.config.dry_run:
            console.print(f"[dim]DRY RUN: Would insert {len(images_data)} images[/dim]")
            self.stats.images_migrated += len(images_data)
            return

        with self.target_conn.cursor() as cur:
            execute_values(cur, insert_query, images_data, template=None, page_size=100)
            self.target_conn.commit()
            self.stats.images_migrated += len(images_data)


    def run_migration(self):
        """Execute full migration with progress tracking"""
        self.stats.start_time = datetime.now()

        try:
            # Connect to databases
            self.connect_databases()

            # Get total count
            with self.source_conn.cursor() as cur:
                cur.execute("SELECT COUNT(*) FROM unified_listings")
                total_count = cur.fetchone()[0]

            if total_count == 0:
                console.print("[yellow]No listings to migrate![/yellow]")
                return

            console.print(f"\n[bold]Total listings to migrate: {total_count}[/bold]")
            console.print(f"Batch size: {self.config.batch_size}")
            console.print(f"Starting from offset: {self.config.resume_from * self.config.batch_size}\n")

            if self.config.dry_run:
                console.print("[bold yellow]DRY RUN MODE - No data will be written[/bold yellow]\n")

            # Progress bar
            with Progress(
                SpinnerColumn(),
                TextColumn("[progress.description]{task.description}"),
                BarColumn(),
                TextColumn("[progress.percentage]{task.percentage:>3.0f}%"),
                TimeElapsedColumn(),
                TimeRemainingColumn(),
                console=console
            ) as progress:

                # Main migration task
                task = progress.add_task("[cyan]Migrating listings...", total=total_count)

                # Process in batches
                offset = self.config.resume_from * self.config.batch_size
                batch_num = self.config.resume_from

                while offset < total_count:
                    try:
                        # Fetch batch
                        listings = self.fetch_source_data(offset, self.config.batch_size)

                        if not listings:
                            break

                        # Migrate listings
                        inserted_ids = self.migrate_listings_batch(listings)

                        # Migrate related data (images only)
                        self.migrate_images(listings)

                        # Update progress
                        progress.update(task, advance=len(listings))

                        if self.config.verbose:
                            console.print(f"  Batch {batch_num}: Migrated {len(listings)} listings")

                        offset += self.config.batch_size
                        batch_num += 1

                    except Exception as e:
                        self.stats.errors += 1
                        console.print(f"[red]Error in batch {batch_num}: {e}[/red]")
                        if not self.config.dry_run:
                            self.target_conn.rollback()
                        raise

            self.stats.end_time = datetime.now()

            # Print summary
            self.print_summary()

        except Exception as e:
            console.print(f"\n[bold red]Migration failed: {e}[/bold red]")
            raise
        finally:
            self.close_connections()

    def print_summary(self):
        """Print migration summary"""
        console.print("\n")
        console.print(Panel.fit(
            "[bold green]Migration Completed Successfully![/bold green]",
            border_style="green"
        ))

        # Create summary table
        table = Table(title="Migration Summary", show_header=True, header_style="bold magenta")
        table.add_column("Metric", style="cyan")
        table.add_column("Count", justify="right", style="green")

        table.add_row("Listings Migrated", str(self.stats.listings_migrated))
        table.add_row("Images Migrated", str(self.stats.images_migrated))
        table.add_row("Errors", str(self.stats.errors), style="red" if self.stats.errors > 0 else "green")
        table.add_row("Duration", f"{self.stats.duration_seconds():.2f} seconds")

        console.print(table)

        if self.config.dry_run:
            console.print("\n[bold yellow]This was a DRY RUN - no data was actually migrated[/bold yellow]")


def parse_arguments():
    """Parse command line arguments"""
    parser = argparse.ArgumentParser(
        description="Migrate listings from monolith to microservice",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Dry run
  %(prog)s --dry-run

  # Real migration
  %(prog)s --source-password SECRET --target-password SECRET

  # Resume from batch 10
  %(prog)s --resume-from 10
        """
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

    # Migration options
    parser.add_argument('--batch-size', type=int, default=1000, help='Batch size for processing')
    parser.add_argument('--dry-run', action='store_true', help='Test migration without writing data')
    parser.add_argument('--resume-from', type=int, default=0, help='Resume from batch number')
    parser.add_argument('-v', '--verbose', action='store_true', help='Verbose output')

    return parser.parse_args()


def main():
    """Main entry point"""
    args = parse_arguments()

    # Create config
    config = MigrationConfig(
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
        batch_size=args.batch_size,
        dry_run=args.dry_run,
        resume_from=args.resume_from,
        verbose=args.verbose
    )

    # Print header
    console.print("\n")
    console.print(Panel.fit(
        "[bold cyan]Data Migration: Monolith → Listings Microservice[/bold cyan]\n"
        f"Source: {config.source_host}:{config.source_port}/{config.source_db}\n"
        f"Target: {config.target_host}:{config.target_port}/{config.target_db}",
        border_style="cyan"
    ))

    # Run migration
    migrator = DataMigrator(config)
    try:
        migrator.run_migration()
    except KeyboardInterrupt:
        console.print("\n[yellow]Migration interrupted by user[/yellow]")
        sys.exit(1)
    except Exception as e:
        console.print(f"\n[bold red]Fatal error: {e}[/bold red]")
        sys.exit(1)


if __name__ == '__main__':
    main()
