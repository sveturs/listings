#!/usr/bin/env python3
"""
OpenSearch reindex via docker exec (workaround for pg_hba.conf auth issues)
"""
import json
import subprocess
import sys
from datetime import datetime
from typing import Any, Dict, List

import requests
from rich.console import Console
from rich.progress import Progress, SpinnerColumn, TextColumn

console = Console()


def convert_timestamp_to_iso(ts: str) -> str:
    """Convert PostgreSQL timestamp to ISO8601 format for OpenSearch"""
    if not ts or ts == '\\N':
        return None

    try:
        # Parse PostgreSQL timestamp: "2025-10-11 17:38:47.121733+00"
        # Remove timezone offset and parse
        if '+' in ts:
            ts_clean = ts.split('+')[0]
        elif '-' in ts and ts.count('-') > 2:  # Has timezone offset
            ts_clean = ts.rsplit('-', 1)[0]
        else:
            ts_clean = ts

        # Parse datetime
        dt = datetime.fromisoformat(ts_clean.replace(' ', 'T'))

        # Return ISO8601 with Z suffix
        return dt.isoformat() + 'Z'

    except (ValueError, AttributeError) as e:
        console.print(f"[yellow]Warning: Could not parse timestamp '{ts}': {e}[/yellow]")
        return None


def get_listings_from_postgres() -> List[Dict[str, Any]]:
    """Fetch listings from PostgreSQL via docker exec"""
    console.print("[cyan]Fetching listings from PostgreSQL via docker exec...[/cyan]")

    sql_query = """
    SELECT
        l.id, l.uuid, l.user_id, l.category_id, l.storefront_id,
        l.title, l.description, l.price, l.currency, l.status,
        l.visibility, l.quantity, l.sku, l.views_count, l.favorites_count,
        l.published_at, l.is_deleted, l.created_at, l.updated_at,
        COALESCE(
            json_agg(
                json_build_object(
                    'id', i.id,
                    'listing_id', i.listing_id,
                    'url', i.url,
                    'display_order', i.display_order,
                    'is_primary', i.is_primary
                )
            ) FILTER (WHERE i.id IS NOT NULL),
            '[]'::json
        ) as images
    FROM listings l
    LEFT JOIN listing_images i ON l.id = i.listing_id
    GROUP BY l.id
    ORDER BY l.id;
    """

    cmd = [
        "docker", "exec", "listings_postgres",
        "psql", "-U", "listings_user", "-d", "listings_db",
        "-t", "-A", "-F", "\t",
        "-c", sql_query
    ]

    try:
        result = subprocess.run(cmd, capture_output=True, text=True, check=True)
        lines = result.stdout.strip().split('\n')

        listings = []
        for line in lines:
            if not line.strip():
                continue

            fields = line.split('\t')
            if len(fields) < 20:
                continue

            # Parse images JSON
            images_json = fields[19] if len(fields) > 19 else '[]'
            try:
                images = json.loads(images_json)
            except json.JSONDecodeError:
                images = []

            listing = {
                'id': int(fields[0]) if fields[0] else None,
                'uuid': fields[1],
                'user_id': int(fields[2]) if fields[2] else None,
                'category_id': int(fields[3]) if fields[3] else None,
                'storefront_id': int(fields[4]) if fields[4] and fields[4] != '\\N' else None,
                'title': fields[5],
                'description': fields[6] if fields[6] != '\\N' else None,
                'price': float(fields[7]) if fields[7] else 0.0,
                'currency': fields[8],
                'status': fields[9],
                'visibility': fields[10],
                'quantity': int(fields[11]) if fields[11] else 1,
                'sku': fields[12] if fields[12] != '\\N' else None,
                'views_count': int(fields[13]) if fields[13] else 0,
                'favorites_count': int(fields[14]) if fields[14] else 0,
                'published_at': fields[15] if fields[15] != '\\N' else None,
                'is_deleted': fields[16] == 't',
                'created_at': fields[17],
                'updated_at': fields[18],
                'images': images
            }
            listings.append(listing)

        console.print(f"[green]✓ Fetched {len(listings)} listings from PostgreSQL[/green]")
        return listings

    except subprocess.CalledProcessError as e:
        console.print(f"[red]✗ Failed to fetch data from PostgreSQL: {e}[/red]")
        console.print(f"[red]STDERR: {e.stderr}[/red]")
        sys.exit(1)


def bulk_index_to_opensearch(listings: List[Dict[str, Any]], opensearch_url: str, index_name: str) -> None:
    """Bulk index listings to OpenSearch"""
    console.print(f"[cyan]Indexing {len(listings)} listings to OpenSearch...[/cyan]")

    if not listings:
        console.print("[yellow]No listings to index[/yellow]")
        return

    # Build bulk request body
    bulk_body = []
    for listing in listings:
        # Index action
        bulk_body.append(json.dumps({
            "index": {
                "_index": index_name,
                "_id": str(listing['id'])
            }
        }))

        # Document
        doc = {
            'id': listing['id'],
            'uuid': listing['uuid'],
            'user_id': listing['user_id'],
            'category_id': listing['category_id'],
            'storefront_id': listing['storefront_id'],
            'title': listing['title'],
            'description': listing['description'],
            'price': listing['price'],
            'currency': listing['currency'],
            'status': listing['status'],
            'visibility': listing['visibility'],
            'quantity': listing['quantity'],
            'sku': listing['sku'],
            'views_count': listing['views_count'],
            'favorites_count': listing['favorites_count'],
            'published_at': convert_timestamp_to_iso(listing['published_at']),
            'is_deleted': listing['is_deleted'],
            'created_at': convert_timestamp_to_iso(listing['created_at']),
            'updated_at': convert_timestamp_to_iso(listing['updated_at']),
            'images': listing['images'],
            'indexed_at': datetime.utcnow().isoformat() + 'Z'
        }
        bulk_body.append(json.dumps(doc))

    bulk_data = '\n'.join(bulk_body) + '\n'

    # Send bulk request
    url = f"{opensearch_url}/{index_name}/_bulk"
    headers = {'Content-Type': 'application/x-ndjson'}

    try:
        response = requests.post(url, data=bulk_data, headers=headers)
        response.raise_for_status()

        result = response.json()

        if result.get('errors'):
            console.print("[red]✗ Some documents failed to index[/red]")
            for item in result['items']:
                if 'error' in item.get('index', {}):
                    console.print(f"[red]  - Doc ID {item['index']['_id']}: {item['index']['error']}[/red]")
        else:
            console.print(f"[green]✓ Successfully indexed {len(listings)} documents[/green]")

    except requests.exceptions.RequestException as e:
        console.print(f"[red]✗ Failed to index documents: {e}[/red]")
        sys.exit(1)


def main():
    opensearch_url = "http://localhost:9200"
    index_name = "listings_microservice"

    console.print("[bold cyan]OpenSearch Reindex (via Docker)[/bold cyan]\n")

    # Check OpenSearch connection
    console.print(f"[cyan]Checking OpenSearch connection ({opensearch_url})...[/cyan]")
    try:
        response = requests.get(opensearch_url)
        response.raise_for_status()
        console.print("[green]✓ Connected to OpenSearch[/green]\n")
    except requests.exceptions.RequestException as e:
        console.print(f"[red]✗ Failed to connect to OpenSearch: {e}[/red]")
        sys.exit(1)

    # Fetch listings from PostgreSQL
    listings = get_listings_from_postgres()

    # Index to OpenSearch
    bulk_index_to_opensearch(listings, opensearch_url, index_name)

    # Summary
    console.print(f"\n[bold green]✓ Reindex completed successfully![/bold green]")
    console.print(f"  - Total listings indexed: {len(listings)}")
    console.print(f"  - Index: {index_name}")
    console.print(f"  - OpenSearch URL: {opensearch_url}")


if __name__ == '__main__':
    main()
