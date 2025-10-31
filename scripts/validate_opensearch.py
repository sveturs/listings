#!/usr/bin/env python3
"""
Validate OpenSearch index after reindex

This script performs comprehensive validation of the OpenSearch index:
- Document count matches PostgreSQL
- All required fields are indexed
- Search queries work correctly
- Facets and aggregations work
- Performance is acceptable

Usage:
    python3 scripts/validate_opensearch.py [options]

Options:
    --target-host HOST       Target database host (default: localhost)
    --target-port PORT       Target database port (default: 35433)
    --target-user USER       Target database user (default: listings_user)
    --target-password PASS   Target database password (required)
    --target-db DB           Target database name (default: listings_db)
    --opensearch-url URL     OpenSearch URL (default: http://localhost:9200)
    --opensearch-index NAME  OpenSearch index name (default: listings_microservice)
    --verbose                Enable verbose output
"""

import sys
import argparse
from typing import Dict, Any, List
import time

try:
    import psycopg2
    import requests
    from rich.console import Console
    from rich.table import Table
    from rich.panel import Panel
    from rich.progress import Progress, SpinnerColumn, TextColumn
except ImportError:
    print("Error: Required packages not installed.")
    print("Install with: pip3 install psycopg2-binary requests rich")
    sys.exit(1)

console = Console()


class ValidationResult:
    """Stores validation results"""

    def __init__(self):
        self.passed = []
        self.failed = []
        self.warnings = []

    def add_pass(self, test: str, message: str = ""):
        self.passed.append((test, message))

    def add_fail(self, test: str, message: str):
        self.failed.append((test, message))

    def add_warning(self, test: str, message: str):
        self.warnings.append((test, message))

    def is_success(self) -> bool:
        return len(self.failed) == 0

    def print_results(self):
        """Print validation results"""
        console.print("\n" + "=" * 70)
        console.print("[bold cyan]Validation Results[/bold cyan]")
        console.print("=" * 70)

        if self.passed:
            console.print(f"\n[green]✓ Passed ({len(self.passed)})[/green]")
            for test, msg in self.passed:
                console.print(f"  • {test}")
                if msg:
                    console.print(f"    {msg}")

        if self.warnings:
            console.print(f"\n[yellow]⚠ Warnings ({len(self.warnings)})[/yellow]")
            for test, msg in self.warnings:
                console.print(f"  • {test}")
                if msg:
                    console.print(f"    {msg}")

        if self.failed:
            console.print(f"\n[red]✗ Failed ({len(self.failed)})[/red]")
            for test, msg in self.failed:
                console.print(f"  • {test}")
                if msg:
                    console.print(f"    {msg}")

        console.print("\n" + "=" * 70)
        if self.is_success():
            console.print("[bold green]✓ All validations passed![/bold green]")
        else:
            console.print("[bold red]✗ Some validations failed[/bold red]")
        console.print("=" * 70 + "\n")


class OpenSearchValidator:
    """Validates OpenSearch index"""

    def __init__(self, url: str, index: str, verbose: bool = False):
        self.url = url
        self.index = index
        self.verbose = verbose

    def get_document_count(self) -> int:
        """Get total document count in index"""
        try:
            response = requests.get(f"{self.url}/{self.index}/_count")
            if response.status_code == 200:
                return response.json().get('count', 0)
        except requests.RequestException:
            pass
        return -1

    def get_index_stats(self) -> Dict[str, Any]:
        """Get index statistics"""
        try:
            response = requests.get(f"{self.url}/{self.index}/_stats")
            if response.status_code == 200:
                data = response.json()
                index_data = data.get('indices', {}).get(self.index, {})
                total = index_data.get('total', {})
                return {
                    'docs_count': total.get('docs', {}).get('count', 0),
                    'store_size': total.get('store', {}).get('size_in_bytes', 0),
                    'segments_count': total.get('segments', {}).get('count', 0)
                }
        except requests.RequestException:
            pass
        return {}

    def search_by_title(self, query: str) -> Dict[str, Any]:
        """Test search by title"""
        try:
            search_query = {
                "query": {
                    "match": {
                        "title": query
                    }
                },
                "size": 10
            }

            start = time.time()
            response = requests.post(
                f"{self.url}/{self.index}/_search",
                json=search_query,
                headers={'Content-Type': 'application/json'}
            )
            duration = (time.time() - start) * 1000  # ms

            if response.status_code == 200:
                data = response.json()
                return {
                    'success': True,
                    'hits': data['hits']['total']['value'],
                    'duration_ms': duration,
                    'results': data['hits']['hits']
                }
        except requests.RequestException as e:
            return {'success': False, 'error': str(e)}
        return {'success': False, 'error': 'Unknown error'}

    def test_aggregations(self) -> Dict[str, Any]:
        """Test aggregations (facets)"""
        try:
            agg_query = {
                "size": 0,
                "aggs": {
                    "categories": {
                        "terms": {"field": "category_id", "size": 10}
                    },
                    "statuses": {
                        "terms": {"field": "status", "size": 10}
                    },
                    "price_stats": {
                        "stats": {"field": "price"}
                    }
                }
            }

            start = time.time()
            response = requests.post(
                f"{self.url}/{self.index}/_search",
                json=agg_query,
                headers={'Content-Type': 'application/json'}
            )
            duration = (time.time() - start) * 1000  # ms

            if response.status_code == 200:
                data = response.json()
                aggs = data.get('aggregations', {})
                return {
                    'success': True,
                    'duration_ms': duration,
                    'categories': aggs.get('categories', {}),
                    'statuses': aggs.get('statuses', {}),
                    'price_stats': aggs.get('price_stats', {})
                }
        except requests.RequestException as e:
            return {'success': False, 'error': str(e)}
        return {'success': False, 'error': 'Unknown error'}

    def get_sample_document(self) -> Dict[str, Any]:
        """Get a sample document to check fields"""
        try:
            response = requests.post(
                f"{self.url}/{self.index}/_search",
                json={"size": 1},
                headers={'Content-Type': 'application/json'}
            )

            if response.status_code == 200:
                data = response.json()
                hits = data['hits']['hits']
                if hits:
                    return hits[0]['_source']
        except requests.RequestException:
            pass
        return {}

    def test_geo_query(self) -> Dict[str, Any]:
        """Test geo-location query"""
        try:
            # Search for documents with location within 50km of Belgrade (44.7866, 20.4489)
            geo_query = {
                "query": {
                    "bool": {
                        "filter": {
                            "geo_distance": {
                                "distance": "50km",
                                "location": {
                                    "lat": 44.7866,
                                    "lon": 20.4489
                                }
                            }
                        }
                    }
                },
                "size": 10
            }

            start = time.time()
            response = requests.post(
                f"{self.url}/{self.index}/_search",
                json=geo_query,
                headers={'Content-Type': 'application/json'}
            )
            duration = (time.time() - start) * 1000  # ms

            if response.status_code == 200:
                data = response.json()
                return {
                    'success': True,
                    'hits': data['hits']['total']['value'],
                    'duration_ms': duration
                }
        except requests.RequestException as e:
            return {'success': False, 'error': str(e)}
        return {'success': False, 'error': 'Unknown error'}


def validate_document_count(pg_conn, os_validator: OpenSearchValidator, results: ValidationResult):
    """Validate document count matches"""
    console.print("\n[cyan]Validating document count...[/cyan]")

    # Get PostgreSQL count
    with pg_conn.cursor() as cur:
        cur.execute("SELECT COUNT(*) FROM listings WHERE is_deleted = false")
        pg_count = cur.fetchone()[0]

    # Get OpenSearch count
    os_count = os_validator.get_document_count()

    if os_count == -1:
        results.add_fail("Document Count", "Failed to get OpenSearch count")
        return

    console.print(f"  PostgreSQL: {pg_count}")
    console.print(f"  OpenSearch: {os_count}")

    if pg_count == os_count:
        results.add_pass("Document Count", f"Matched: {os_count} documents")
    else:
        diff = abs(pg_count - os_count)
        results.add_fail("Document Count", f"Mismatch: PostgreSQL={pg_count}, OpenSearch={os_count} (diff={diff})")


def validate_required_fields(os_validator: OpenSearchValidator, results: ValidationResult):
    """Validate all required fields are present"""
    console.print("\n[cyan]Validating required fields...[/cyan]")

    sample = os_validator.get_sample_document()
    if not sample:
        results.add_fail("Required Fields", "No documents found to validate")
        return

    required_fields = [
        'id', 'uuid', 'user_id', 'title', 'description',
        'price', 'currency', 'category_id', 'status',
        'created_at', 'updated_at'
    ]

    missing = [field for field in required_fields if field not in sample]

    if not missing:
        results.add_pass("Required Fields", f"All {len(required_fields)} required fields present")
    else:
        results.add_fail("Required Fields", f"Missing fields: {', '.join(missing)}")


def validate_search_functionality(os_validator: OpenSearchValidator, results: ValidationResult):
    """Validate search works"""
    console.print("\n[cyan]Validating search functionality...[/cyan]")

    # Test simple search
    search_result = os_validator.search_by_title("listing")

    if not search_result.get('success'):
        results.add_fail("Search Functionality", f"Search failed: {search_result.get('error')}")
        return

    duration = search_result['duration_ms']
    hits = search_result['hits']

    console.print(f"  Search query: 'listing'")
    console.print(f"  Results: {hits} hits")
    console.print(f"  Duration: {duration:.2f}ms")

    if duration > 100:
        results.add_warning("Search Performance", f"Search took {duration:.2f}ms (>100ms threshold)")
    else:
        results.add_pass("Search Functionality", f"Search works ({duration:.2f}ms)")


def validate_aggregations(os_validator: OpenSearchValidator, results: ValidationResult):
    """Validate aggregations work"""
    console.print("\n[cyan]Validating aggregations...[/cyan]")

    agg_result = os_validator.test_aggregations()

    if not agg_result.get('success'):
        results.add_fail("Aggregations", f"Aggregations failed: {agg_result.get('error')}")
        return

    duration = agg_result['duration_ms']
    console.print(f"  Duration: {duration:.2f}ms")

    # Check if we got aggregation results
    categories = agg_result.get('categories', {}).get('buckets', [])
    statuses = agg_result.get('statuses', {}).get('buckets', [])
    price_stats = agg_result.get('price_stats', {})

    console.print(f"  Categories: {len(categories)} buckets")
    console.print(f"  Statuses: {len(statuses)} buckets")
    console.print(f"  Price stats: min={price_stats.get('min', 0):.2f}, max={price_stats.get('max', 0):.2f}")

    if duration > 200:
        results.add_warning("Aggregation Performance", f"Aggregations took {duration:.2f}ms (>200ms threshold)")
    else:
        results.add_pass("Aggregations", f"Aggregations work ({duration:.2f}ms)")


def validate_geo_search(os_validator: OpenSearchValidator, results: ValidationResult):
    """Validate geo-location search"""
    console.print("\n[cyan]Validating geo-location search...[/cyan]")

    geo_result = os_validator.test_geo_query()

    if not geo_result.get('success'):
        results.add_warning("Geo Search", f"Geo search failed (may be expected if no location data): {geo_result.get('error')}")
        return

    hits = geo_result['hits']
    duration = geo_result['duration_ms']

    console.print(f"  Results: {hits} hits")
    console.print(f"  Duration: {duration:.2f}ms")

    if hits == 0:
        results.add_warning("Geo Search", "No documents with location data found")
    else:
        results.add_pass("Geo Search", f"Geo search works ({hits} hits, {duration:.2f}ms)")


def validate_index_stats(os_validator: OpenSearchValidator, results: ValidationResult):
    """Validate index statistics"""
    console.print("\n[cyan]Validating index statistics...[/cyan]")

    stats = os_validator.get_index_stats()
    if not stats:
        results.add_fail("Index Stats", "Failed to get index statistics")
        return

    docs_count = stats.get('docs_count', 0)
    store_size_mb = stats.get('store_size', 0) / (1024 * 1024)
    segments_count = stats.get('segments_count', 0)

    console.print(f"  Documents: {docs_count}")
    console.print(f"  Store size: {store_size_mb:.2f} MB")
    console.print(f"  Segments: {segments_count}")

    if segments_count > 10:
        results.add_warning("Index Segments", f"High segment count ({segments_count}), consider force merge")

    results.add_pass("Index Stats", f"{docs_count} docs, {store_size_mb:.2f} MB, {segments_count} segments")


def main():
    parser = argparse.ArgumentParser(description='Validate OpenSearch index')
    parser.add_argument('--target-host', default='localhost', help='Target database host')
    parser.add_argument('--target-port', type=int, default=35433, help='Target database port')
    parser.add_argument('--target-user', default='listings_user', help='Target database user')
    parser.add_argument('--target-password', required=True, help='Target database password')
    parser.add_argument('--target-db', default='listings_db', help='Target database name')
    parser.add_argument('--opensearch-url', default='http://localhost:9200', help='OpenSearch URL')
    parser.add_argument('--opensearch-index', default='listings_microservice', help='OpenSearch index name')
    parser.add_argument('--verbose', action='store_true', help='Enable verbose output')
    args = parser.parse_args()

    console.print(Panel.fit(
        "[bold cyan]OpenSearch Index Validation[/bold cyan]\n"
        f"Index: {args.opensearch_index}\n"
        f"OpenSearch: {args.opensearch_url}\n"
        f"Database: {args.target_user}@{args.target_host}:{args.target_port}/{args.target_db}",
        border_style="cyan"
    ))

    results = ValidationResult()

    # Initialize validator
    os_validator = OpenSearchValidator(args.opensearch_url, args.opensearch_index, args.verbose)

    # Connect to PostgreSQL
    try:
        console.print("\n[cyan]Connecting to PostgreSQL...[/cyan]")
        pg_conn = psycopg2.connect(
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
        # Run validations
        with Progress(
            SpinnerColumn(),
            TextColumn("[progress.description]{task.description}"),
            console=console
        ) as progress:
            task = progress.add_task("[cyan]Running validations...", total=None)

            validate_document_count(pg_conn, os_validator, results)
            validate_required_fields(os_validator, results)
            validate_search_functionality(os_validator, results)
            validate_aggregations(os_validator, results)
            validate_geo_search(os_validator, results)
            validate_index_stats(os_validator, results)

        # Print results
        results.print_results()

        # Exit code
        sys.exit(0 if results.is_success() else 1)

    finally:
        pg_conn.close()


if __name__ == "__main__":
    main()
