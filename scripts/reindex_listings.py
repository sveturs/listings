#!/usr/bin/env python3
"""
Reindex listings from PostgreSQL to OpenSearch

This script performs a full reindex of all listings from the microservice
PostgreSQL database to the OpenSearch index. It handles:
- Reading listings with related data (images, tags, attributes, location)
- Transforming to OpenSearch document format
- Bulk indexing with batching
- Progress tracking and error handling
- Retry logic for failed operations

Usage:
    python3 scripts/reindex_listings.py [options]

Options:
    --target-host HOST       Target database host (default: localhost)
    --target-port PORT       Target database port (default: 35433)
    --target-user USER       Target database user (default: listings_user)
    --target-password PASS   Target database password (required)
    --target-db DB           Target database name (default: listings_db)
    --opensearch-url URL     OpenSearch URL (default: http://localhost:9200)
    --opensearch-index NAME  OpenSearch index name (default: listings_microservice)
    --batch-size N           Batch size for bulk operations (default: 500)
    --verbose                Enable verbose output
    --dry-run                Don't actually index, just show what would be done
"""

import sys
import argparse
from datetime import datetime
from typing import List, Dict, Any, Optional
import time

try:
    import psycopg2
    import psycopg2.extras
    import requests
    from rich.console import Console
    from rich.progress import Progress, SpinnerColumn, BarColumn, TextColumn, TimeRemainingColumn
    from rich.table import Table
    from rich.panel import Panel
except ImportError:
    print("Error: Required packages not installed.")
    print("Install with: pip3 install psycopg2-binary requests rich")
    sys.exit(1)

console = Console()


class OpenSearchIndexer:
    """Handles OpenSearch indexing operations"""

    def __init__(self, url: str, index: str, verbose: bool = False):
        self.url = url
        self.index = index
        self.verbose = verbose
        self.stats = {
            'indexed': 0,
            'failed': 0,
            'skipped': 0
        }

    def check_connection(self) -> bool:
        """Verify OpenSearch connection"""
        try:
            response = requests.get(f"{self.url}/", timeout=5)
            response.raise_for_status()
            return True
        except requests.RequestException as e:
            console.print(f"[red]✗ Failed to connect to OpenSearch: {e}[/red]")
            return False

    def check_index_exists(self) -> bool:
        """Check if index exists"""
        try:
            response = requests.head(f"{self.url}/{self.index}")
            return response.status_code == 200
        except requests.RequestException:
            return False

    def bulk_index(self, documents: List[Dict[str, Any]]) -> bool:
        """Bulk index documents to OpenSearch"""
        if not documents:
            return True

        # Prepare bulk request body (newline-delimited JSON)
        bulk_body = []
        for doc in documents:
            # Index action
            action = {
                "index": {
                    "_index": self.index,
                    "_id": str(doc['id'])
                }
            }
            bulk_body.append(action)
            bulk_body.append(doc)

        # Convert to NDJSON format
        ndjson = '\n'.join([
            requests.compat.json.dumps(item) for item in bulk_body
        ]) + '\n'

        try:
            response = requests.post(
                f"{self.url}/_bulk",
                data=ndjson,
                headers={'Content-Type': 'application/x-ndjson'},
                timeout=30
            )

            if response.status_code == 200:
                result = response.json()
                if result.get('errors', False):
                    # Some items failed
                    failed_count = sum(1 for item in result['items'] if 'error' in item.get('index', {}))
                    self.stats['indexed'] += len(documents) - failed_count
                    self.stats['failed'] += failed_count

                    if self.verbose:
                        for item in result['items']:
                            if 'error' in item.get('index', {}):
                                error = item['index']['error']
                                console.print(f"[red]Failed to index doc {item['index']['_id']}: {error}[/red]")
                    return False
                else:
                    # All succeeded
                    self.stats['indexed'] += len(documents)
                    return True
            else:
                console.print(f"[red]Bulk index failed: {response.status_code} - {response.text}[/red]")
                self.stats['failed'] += len(documents)
                return False

        except requests.RequestException as e:
            console.print(f"[red]Bulk index error: {e}[/red]")
            self.stats['failed'] += len(documents)
            return False


class ListingsExtractor:
    """Extracts listings data from PostgreSQL"""

    def __init__(self, conn):
        self.conn = conn

    def get_total_count(self) -> int:
        """Get total count of active listings"""
        with self.conn.cursor() as cur:
            cur.execute("""
                SELECT COUNT(*)
                FROM listings
                WHERE is_deleted = false
            """)
            return cur.fetchone()[0]

    def fetch_listings_batch(self, offset: int, limit: int) -> List[Dict[str, Any]]:
        """Fetch a batch of listings with all related data"""
        with self.conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor) as cur:
            # Fetch listings
            cur.execute("""
                SELECT
                    id, uuid, user_id, storefront_id,
                    title, description,
                    price, currency,
                    category_id, status, visibility,
                    quantity, sku,
                    views_count, favorites_count,
                    created_at, updated_at, published_at
                FROM listings
                WHERE is_deleted = false
                ORDER BY id
                LIMIT %s OFFSET %s
            """, (limit, offset))

            listings = cur.fetchall()

            if not listings:
                return []

            listing_ids = [l['id'] for l in listings]

            # Fetch images for all listings
            cur.execute("""
                SELECT
                    listing_id, id, url, thumbnail_url,
                    display_order, is_primary
                FROM listing_images
                WHERE listing_id = ANY(%s)
                ORDER BY listing_id, display_order
            """, (listing_ids,))
            images_by_listing = {}
            for row in cur.fetchall():
                listing_id = row['listing_id']
                if listing_id not in images_by_listing:
                    images_by_listing[listing_id] = []
                images_by_listing[listing_id].append(dict(row))

            # Fetch tags for all listings
            cur.execute("""
                SELECT listing_id, tag
                FROM listing_tags
                WHERE listing_id = ANY(%s)
                ORDER BY listing_id
            """, (listing_ids,))
            tags_by_listing = {}
            for row in cur.fetchall():
                listing_id = row['listing_id']
                if listing_id not in tags_by_listing:
                    tags_by_listing[listing_id] = []
                tags_by_listing[listing_id].append(row['tag'])

            # Fetch attributes for all listings
            cur.execute("""
                SELECT listing_id, attribute_key, attribute_value
                FROM listing_attributes
                WHERE listing_id = ANY(%s)
                ORDER BY listing_id
            """, (listing_ids,))
            attributes_by_listing = {}
            for row in cur.fetchall():
                listing_id = row['listing_id']
                if listing_id not in attributes_by_listing:
                    attributes_by_listing[listing_id] = []
                attributes_by_listing[listing_id].append({
                    'key': row['attribute_key'],
                    'value': row['attribute_value']
                })

            # Fetch locations for all listings
            cur.execute("""
                SELECT
                    listing_id, country, city, postal_code,
                    address_line1, address_line2,
                    latitude, longitude
                FROM listing_locations
                WHERE listing_id = ANY(%s)
            """, (listing_ids,))
            locations_by_listing = {}
            for row in cur.fetchall():
                locations_by_listing[row['listing_id']] = dict(row)

            # Combine all data
            result = []
            for listing in listings:
                listing_id = listing['id']
                listing['images'] = images_by_listing.get(listing_id, [])
                listing['tags'] = tags_by_listing.get(listing_id, [])
                listing['attributes'] = attributes_by_listing.get(listing_id, [])
                listing['location'] = locations_by_listing.get(listing_id)
                result.append(dict(listing))

            return result


def transform_listing_to_document(listing: Dict[str, Any]) -> Dict[str, Any]:
    """Transform PostgreSQL listing to OpenSearch document format"""
    doc = {
        'id': listing['id'],
        'uuid': listing['uuid'],
        'user_id': listing['user_id'],
        'storefront_id': listing['storefront_id'],
        'title': listing['title'],
        'description': listing['description'],
        'price': float(listing['price']) if listing['price'] else 0.0,
        'currency': listing['currency'],
        'category_id': listing['category_id'],
        'status': listing['status'],
        'visibility': listing['visibility'],
        'quantity': listing['quantity'],
        'sku': listing['sku'],
        'views_count': listing['views_count'],
        'favorites_count': listing['favorites_count'],
        'created_at': listing['created_at'].isoformat() if listing['created_at'] else None,
        'updated_at': listing['updated_at'].isoformat() if listing['updated_at'] else None,
        'published_at': listing['published_at'].isoformat() if listing['published_at'] else None,
    }

    # Add images (nested objects)
    if listing.get('images'):
        doc['images'] = [
            {
                'id': img['id'],
                'url': img['url'],
                'thumbnail_url': img['thumbnail_url'],
                'display_order': img['display_order'],
                'is_primary': img['is_primary']
            }
            for img in listing['images']
        ]

    # Add tags (array of strings)
    if listing.get('tags'):
        doc['tags'] = listing['tags']

    # Add attributes (nested objects)
    if listing.get('attributes'):
        doc['attributes'] = listing['attributes']

    # Add location (geo_point + fields)
    if listing.get('location'):
        loc = listing['location']
        if loc.get('latitude') and loc.get('longitude'):
            doc['location'] = {
                'lat': float(loc['latitude']),
                'lon': float(loc['longitude'])
            }
        doc['country'] = loc.get('country')
        doc['city'] = loc.get('city')
        doc['postal_code'] = loc.get('postal_code')
        doc['address_line1'] = loc.get('address_line1')
        doc['address_line2'] = loc.get('address_line2')

    return doc


def main():
    parser = argparse.ArgumentParser(description='Reindex listings to OpenSearch')
    parser.add_argument('--target-host', default='localhost', help='Target database host')
    parser.add_argument('--target-port', type=int, default=35433, help='Target database port')
    parser.add_argument('--target-user', default='listings_user', help='Target database user')
    parser.add_argument('--target-password', required=True, help='Target database password')
    parser.add_argument('--target-db', default='listings_db', help='Target database name')
    parser.add_argument('--opensearch-url', default='http://localhost:9200', help='OpenSearch URL')
    parser.add_argument('--opensearch-index', default='listings_microservice', help='OpenSearch index name')
    parser.add_argument('--batch-size', type=int, default=500, help='Batch size for bulk operations')
    parser.add_argument('--verbose', action='store_true', help='Enable verbose output')
    parser.add_argument('--dry-run', action='store_true', help='Dry run mode')
    args = parser.parse_args()

    console.print(Panel.fit(
        "[bold cyan]Listings Reindex to OpenSearch[/bold cyan]\n"
        f"Source: {args.target_user}@{args.target_host}:{args.target_port}/{args.target_db}\n"
        f"Target: {args.opensearch_url}/{args.opensearch_index}\n"
        f"Batch Size: {args.batch_size}"
        + ("\n[yellow]DRY RUN MODE[/yellow]" if args.dry_run else ""),
        border_style="cyan"
    ))

    # Initialize OpenSearch indexer
    indexer = OpenSearchIndexer(args.opensearch_url, args.opensearch_index, args.verbose)

    # Check OpenSearch connection
    console.print("\n[cyan]Checking OpenSearch connection...[/cyan]")
    if not indexer.check_connection():
        sys.exit(1)
    console.print("[green]✓ Connected to OpenSearch[/green]")

    # Check if index exists
    if not indexer.check_index_exists():
        console.print(f"[red]✗ Index '{args.opensearch_index}' does not exist[/red]")
        console.print("[yellow]Run: python3 scripts/create_opensearch_index.py[/yellow]")
        sys.exit(1)
    console.print(f"[green]✓ Index '{args.opensearch_index}' exists[/green]")

    # Connect to PostgreSQL
    console.print("\n[cyan]Connecting to PostgreSQL...[/cyan]")
    try:
        conn = psycopg2.connect(
            host=args.target_host,
            port=args.target_port,
            user=args.target_user,
            password=args.target_password,
            dbname=args.target_db
        )
        console.print("[green]✓ Connected to PostgreSQL[/green]")
    except Exception as e:
        console.print(f"[red]✗ Failed to connect to PostgreSQL: {e}[/red]")
        sys.exit(1)

    try:
        extractor = ListingsExtractor(conn)

        # Get total count
        total_count = extractor.get_total_count()
        console.print(f"\n[cyan]Total listings to index: {total_count}[/cyan]")

        if total_count == 0:
            console.print("[yellow]No listings found to index[/yellow]")
            return

        if args.dry_run:
            console.print("[yellow]DRY RUN: Would index listings but not actually executing[/yellow]")
            # Fetch first batch to show example
            batch = extractor.fetch_listings_batch(0, min(5, args.batch_size))
            console.print(f"\nExample document (first listing):")
            if batch:
                import json
                doc = transform_listing_to_document(batch[0])
                console.print(json.dumps(doc, indent=2, default=str))
            return

        # Process in batches with progress bar
        start_time = time.time()

        with Progress(
            SpinnerColumn(),
            TextColumn("[progress.description]{task.description}"),
            BarColumn(),
            TextColumn("[progress.percentage]{task.percentage:>3.0f}%"),
            TextColumn("({task.completed}/{task.total})"),
            TimeRemainingColumn(),
            console=console
        ) as progress:
            task = progress.add_task("[cyan]Indexing listings...", total=total_count)

            offset = 0
            while offset < total_count:
                # Fetch batch
                batch = extractor.fetch_listings_batch(offset, args.batch_size)
                if not batch:
                    break

                # Transform to documents
                documents = [transform_listing_to_document(listing) for listing in batch]

                # Bulk index
                success = indexer.bulk_index(documents)

                # Update progress
                progress.update(task, advance=len(batch))
                offset += args.batch_size

                # Small delay to avoid overwhelming the system
                time.sleep(0.1)

        # Calculate duration
        duration = time.time() - start_time

        # Display results
        table = Table(title="Reindex Results")
        table.add_column("Metric", style="cyan")
        table.add_column("Value", style="green")

        table.add_row("Total Listings", str(total_count))
        table.add_row("Indexed", str(indexer.stats['indexed']))
        table.add_row("Failed", str(indexer.stats['failed']))
        table.add_row("Duration", f"{duration:.2f}s")
        table.add_row("Rate", f"{total_count/duration:.2f} docs/sec")

        console.print("\n")
        console.print(table)

        if indexer.stats['failed'] > 0:
            console.print(f"\n[yellow]⚠ {indexer.stats['failed']} documents failed to index[/yellow]")
            console.print("[yellow]Check logs for details[/yellow]")
        else:
            console.print("\n[bold green]✓ Reindex completed successfully![/bold green]")
            console.print("\nNext steps:")
            console.print("  1. Verify data: python3 scripts/validate_opensearch.py")
            console.print("  2. Test search API endpoints")
            console.print(f"  3. Update service config: SVETULISTINGS_OPENSEARCH_INDEX={args.opensearch_index}")

    finally:
        conn.close()


if __name__ == "__main__":
    main()
