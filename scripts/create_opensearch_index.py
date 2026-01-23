#!/usr/bin/env python3
"""
Create OpenSearch index for listings microservice

This script creates a new OpenSearch index with the proper mapping
for the listings microservice. It can optionally delete an existing
index if it already exists.

Usage:
    python3 scripts/create_opensearch_index.py [--force]

Options:
    --force    Delete existing index without confirmation
"""

import sys
import json
import argparse
from pathlib import Path

try:
    import requests
    from rich.console import Console
    from rich.panel import Panel
    from rich.table import Table
except ImportError:
    print("Error: Required packages not installed.")
    print("Install with: pip3 install requests rich")
    sys.exit(1)

console = Console()

# Configuration
OPENSEARCH_URL = "http://localhost:9200"
INDEX_NAME = "listings_microservice"
SCHEMA_FILE = Path(__file__).parent / "opensearch_schema.json"


def load_schema():
    """Load OpenSearch schema from JSON file"""
    try:
        with open(SCHEMA_FILE) as f:
            return json.load(f)
    except FileNotFoundError:
        console.print(f"[red]Error: Schema file not found: {SCHEMA_FILE}[/red]")
        sys.exit(1)
    except json.JSONDecodeError as e:
        console.print(f"[red]Error: Invalid JSON in schema file: {e}[/red]")
        sys.exit(1)


def check_opensearch_connection():
    """Check if OpenSearch is accessible"""
    try:
        response = requests.get(f"{OPENSEARCH_URL}/", timeout=5)
        response.raise_for_status()
        info = response.json()

        console.print("\n[green]✓[/green] OpenSearch connection successful")
        console.print(f"  Version: {info.get('version', {}).get('number', 'unknown')}")
        console.print(f"  Cluster: {info.get('cluster_name', 'unknown')}")
        return True
    except requests.RequestException as e:
        console.print(f"\n[red]✗ Failed to connect to OpenSearch at {OPENSEARCH_URL}[/red]")
        console.print(f"  Error: {e}")
        return False


def index_exists():
    """Check if index already exists"""
    try:
        response = requests.head(f"{OPENSEARCH_URL}/{INDEX_NAME}")
        return response.status_code == 200
    except requests.RequestException:
        return False


def get_index_info():
    """Get information about existing index"""
    try:
        response = requests.get(f"{OPENSEARCH_URL}/{INDEX_NAME}")
        if response.status_code == 200:
            data = response.json()
            index_data = data.get(INDEX_NAME, {})

            # Get document count
            count_response = requests.get(f"{OPENSEARCH_URL}/{INDEX_NAME}/_count")
            count = count_response.json().get('count', 0) if count_response.status_code == 200 else 0

            return {
                'exists': True,
                'document_count': count,
                'settings': index_data.get('settings', {}),
                'mappings': index_data.get('mappings', {})
            }
    except requests.RequestException:
        pass
    return {'exists': False}


def delete_index():
    """Delete existing index"""
    try:
        response = requests.delete(f"{OPENSEARCH_URL}/{INDEX_NAME}")
        if response.status_code in [200, 404]:
            console.print(f"[green]✓[/green] Deleted existing index: {INDEX_NAME}")
            return True
        else:
            console.print(f"[red]✗ Failed to delete index: {response.text}[/red]")
            return False
    except requests.RequestException as e:
        console.print(f"[red]✗ Error deleting index: {e}[/red]")
        return False


def create_index(schema):
    """Create new index with schema"""
    try:
        response = requests.put(
            f"{OPENSEARCH_URL}/{INDEX_NAME}",
            json=schema,
            headers={'Content-Type': 'application/json'}
        )

        if response.status_code in [200, 201]:
            console.print(f"\n[green]✓ Index '{INDEX_NAME}' created successfully[/green]")

            # Display schema summary
            settings = schema.get('settings', {})
            mappings = schema.get('mappings', {})
            properties = mappings.get('properties', {})

            table = Table(title="Index Configuration")
            table.add_column("Setting", style="cyan")
            table.add_column("Value", style="green")

            table.add_row("Shards", str(settings.get('number_of_shards', 'N/A')))
            table.add_row("Replicas", str(settings.get('number_of_replicas', 'N/A')))
            table.add_row("Mapped Fields", str(len(properties)))
            table.add_row("Refresh Interval", settings.get('index', {}).get('refresh_interval', 'N/A'))

            console.print(table)
            return True
        else:
            console.print(f"\n[red]✗ Failed to create index[/red]")
            console.print(f"Status: {response.status_code}")
            console.print(f"Response: {response.text}")
            return False

    except requests.RequestException as e:
        console.print(f"[red]✗ Error creating index: {e}[/red]")
        return False


def verify_index():
    """Verify index was created correctly"""
    try:
        response = requests.get(f"{OPENSEARCH_URL}/{INDEX_NAME}")
        if response.status_code == 200:
            data = response.json()
            index_data = data.get(INDEX_NAME, {})
            mappings = index_data.get('mappings', {})
            properties = mappings.get('properties', {})

            console.print("\n[green]✓ Index verification passed[/green]")
            console.print(f"  Mapped fields: {len(properties)}")

            # List some key fields
            key_fields = ['id', 'uuid', 'title', 'price', 'status', 'location']
            missing_fields = [f for f in key_fields if f not in properties]

            if missing_fields:
                console.print(f"  [yellow]Warning: Missing expected fields: {', '.join(missing_fields)}[/yellow]")
            else:
                console.print(f"  [green]All key fields present[/green]")

            return True
        else:
            console.print("[red]✗ Index verification failed[/red]")
            return False

    except requests.RequestException as e:
        console.print(f"[red]✗ Error verifying index: {e}[/red]")
        return False


def main():
    parser = argparse.ArgumentParser(description='Create OpenSearch index for listings microservice')
    parser.add_argument('--force', action='store_true', help='Delete existing index without confirmation')
    args = parser.parse_args()

    console.print(Panel.fit(
        "[bold cyan]OpenSearch Index Creation[/bold cyan]\n"
        f"Index: {INDEX_NAME}\n"
        f"OpenSearch: {OPENSEARCH_URL}",
        border_style="cyan"
    ))

    # Check OpenSearch connection
    if not check_opensearch_connection():
        sys.exit(1)

    # Load schema
    console.print("\n[cyan]Loading schema...[/cyan]")
    schema = load_schema()
    console.print(f"[green]✓[/green] Schema loaded from {SCHEMA_FILE}")

    # Check if index exists
    info = get_index_info()
    if info['exists']:
        console.print(f"\n[yellow]⚠[/yellow]  Index '{INDEX_NAME}' already exists")
        console.print(f"  Document count: {info['document_count']}")

        if not args.force:
            response = input("\nDelete and recreate? [yes/no]: ").lower()
            if response != 'yes':
                console.print("[yellow]Operation cancelled[/yellow]")
                sys.exit(0)

        # Delete existing index
        if not delete_index():
            sys.exit(1)

    # Create index
    console.print("\n[cyan]Creating index...[/cyan]")
    if not create_index(schema):
        sys.exit(1)

    # Verify index
    if not verify_index():
        sys.exit(1)

    console.print("\n[bold green]✓ Index creation completed successfully![/bold green]")
    console.print(f"\nNext steps:")
    console.print(f"  1. Run reindex script: python3 scripts/reindex_listings.py")
    console.print(f"  2. Verify data: python3 scripts/validate_opensearch.py")
    console.print(f"  3. Update service config: VONDILISTINGS_OPENSEARCH_INDEX={INDEX_NAME}")


if __name__ == "__main__":
    main()
