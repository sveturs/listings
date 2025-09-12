#!/usr/bin/env python3
"""
Real Translation Audit Tool using spec-kit methodology
Finds ALL missing translations and untranslated keys
"""

import json
import os
from pathlib import Path
from collections import defaultdict
import re

def load_json_file(filepath):
    """Load and parse JSON file"""
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            return json.load(f)
    except Exception as e:
        print(f"Error loading {filepath}: {e}")
        return {}

def extract_all_keys(obj, prefix=''):
    """Extract all keys from nested JSON structure"""
    keys = set()
    if isinstance(obj, dict):
        for key, value in obj.items():
            full_key = f"{prefix}.{key}" if prefix else key
            keys.add(full_key)
            if isinstance(value, dict):
                keys.update(extract_all_keys(value, full_key))
    return keys

def find_untranslated_values(obj, locale, prefix=''):
    """Find values that are still in English or are placeholders"""
    untranslated = []
    
    if isinstance(obj, dict):
        for key, value in obj.items():
            full_key = f"{prefix}.{key}" if prefix else key
            
            if isinstance(value, str):
                # Check if value looks like a key (contains dots or is camelCase)
                if '.' in value or re.match(r'^[a-z]+[A-Z]', value):
                    untranslated.append({
                        'key': full_key,
                        'value': value,
                        'type': 'placeholder'
                    })
                # Check if value is in English for non-English locales
                elif locale != 'en' and value and all(ord(c) < 128 for c in value if c.isalpha()):
                    # Check if it's likely English (all ASCII letters)
                    if len(value) > 3 and not value.isupper():  # Skip abbreviations
                        untranslated.append({
                            'key': full_key,
                            'value': value,
                            'type': 'english_in_non_english'
                        })
            elif isinstance(value, dict):
                untranslated.extend(find_untranslated_values(value, locale, full_key))
    
    return untranslated

def analyze_translations():
    """Main analysis function"""
    base_path = Path('/data/hostel-booking-system/frontend/svetu/src/messages')
    
    if not base_path.exists():
        print(f"Error: Path {base_path} does not exist")
        return
    
    # Get all locales
    locales = [d.name for d in base_path.iterdir() if d.is_dir()]
    print(f"Found locales: {locales}")
    
    # Collect all modules and their keys
    all_modules = set()
    locale_modules = defaultdict(set)
    module_keys = defaultdict(lambda: defaultdict(set))
    
    for locale in locales:
        locale_path = base_path / locale
        for json_file in locale_path.glob('*.json'):
            module_name = json_file.stem
            all_modules.add(module_name)
            locale_modules[locale].add(module_name)
            
            data = load_json_file(json_file)
            keys = extract_all_keys(data)
            module_keys[locale][module_name] = keys
    
    print(f"\nFound {len(all_modules)} modules: {sorted(all_modules)}")
    
    # Find missing modules per locale
    missing_modules = defaultdict(list)
    for locale in locales:
        for module in all_modules:
            if module not in locale_modules[locale]:
                missing_modules[locale].append(module)
    
    # Find missing keys
    missing_keys = defaultdict(lambda: defaultdict(list))
    
    # Use English as reference
    if 'en' in module_keys:
        for module in module_keys['en']:
            en_keys = module_keys['en'][module]
            
            for locale in locales:
                if locale != 'en' and module in module_keys[locale]:
                    locale_keys = module_keys[locale][module]
                    missing = en_keys - locale_keys
                    if missing:
                        missing_keys[locale][module] = sorted(list(missing))
    
    # Find untranslated values (placeholders and English text)
    untranslated_values = defaultdict(lambda: defaultdict(list))
    
    for locale in locales:
        if locale != 'en':  # Skip English
            locale_path = base_path / locale
            for json_file in locale_path.glob('*.json'):
                module_name = json_file.stem
                data = load_json_file(json_file)
                untranslated = find_untranslated_values(data, locale)
                if untranslated:
                    untranslated_values[locale][module_name] = untranslated
    
    # Generate report
    report = {
        'summary': {
            'total_modules': len(all_modules),
            'locales': locales,
            'issues_found': 0
        },
        'missing_modules': missing_modules,
        'missing_keys': {},
        'untranslated_values': {},
        'critical_issues': []
    }
    
    # Count missing keys
    total_missing = 0
    for locale, modules in missing_keys.items():
        report['missing_keys'][locale] = {}
        for module, keys in modules.items():
            if keys:
                report['missing_keys'][locale][module] = keys
                total_missing += len(keys)
                
                # Mark as critical if it's a major module
                if module in ['marketplace', 'common', 'auth', 'storefront']:
                    report['critical_issues'].append({
                        'type': 'missing_keys',
                        'locale': locale,
                        'module': module,
                        'count': len(keys),
                        'keys': keys[:10]  # First 10 as sample
                    })
    
    # Count untranslated values
    total_untranslated = 0
    for locale, modules in untranslated_values.items():
        report['untranslated_values'][locale] = {}
        for module, items in modules.items():
            if items:
                report['untranslated_values'][locale][module] = items
                total_untranslated += len(items)
                
                # Mark placeholder issues as critical
                placeholders = [i for i in items if i['type'] == 'placeholder']
                if placeholders:
                    report['critical_issues'].append({
                        'type': 'untranslated_placeholders',
                        'locale': locale,
                        'module': module,
                        'count': len(placeholders),
                        'examples': placeholders[:5]
                    })
    
    report['summary']['total_missing_keys'] = total_missing
    report['summary']['total_untranslated_values'] = total_untranslated
    report['summary']['issues_found'] = total_missing + total_untranslated
    
    # Calculate coverage
    total_keys = sum(len(keys) for locale_dict in module_keys['en'].values() for keys in [locale_dict])
    if total_keys > 0:
        coverage = ((total_keys * (len(locales) - 1) - total_missing - total_untranslated) / 
                   (total_keys * (len(locales) - 1))) * 100
        report['summary']['coverage_percentage'] = round(coverage, 2)
    
    # Save detailed report
    output_path = Path('/data/hostel-booking-system/spec-kit-translations/reports')
    output_path.mkdir(exist_ok=True)
    
    with open(output_path / 'real_audit_results.json', 'w', encoding='utf-8') as f:
        json.dump(report, f, ensure_ascii=False, indent=2)
    
    # Print summary
    print("\n" + "="*60)
    print("REAL TRANSLATION AUDIT RESULTS")
    print("="*60)
    print(f"\n📊 SUMMARY:")
    print(f"  Total Modules: {report['summary']['total_modules']}")
    print(f"  Locales: {', '.join(report['summary']['locales'])}")
    print(f"  Missing Keys: {total_missing}")
    print(f"  Untranslated Values: {total_untranslated}")
    print(f"  Total Issues: {report['summary']['issues_found']}")
    if 'coverage_percentage' in report['summary']:
        print(f"  Coverage: {report['summary']['coverage_percentage']}%")
    
    print(f"\n⚠️  CRITICAL ISSUES: {len(report['critical_issues'])}")
    for issue in report['critical_issues'][:5]:  # Show first 5
        print(f"\n  • {issue['type'].upper()} in {issue['locale']}/{issue['module']}")
        print(f"    Count: {issue['count']}")
        if 'examples' in issue:
            print(f"    Examples:")
            for ex in issue['examples'][:3]:
                print(f"      - {ex['key']}: '{ex['value']}'")
        elif 'keys' in issue:
            print(f"    Sample keys: {', '.join(issue['keys'][:5])}")
    
    print(f"\n💾 Full report saved to: {output_path / 'real_audit_results.json'}")
    
    return report

if __name__ == '__main__':
    analyze_translations()