#!/usr/bin/env python3
"""
Дополнительный анализатор компонентов для системы переводов
Анализирует конкретные компоненты и их использование переводов
"""

import re
import os
from pathlib import Path
from typing import Dict, List, Set

def analyze_component_translations(project_root: str) -> Dict:
    """Анализ использования переводов в компонентах"""
    
    src_dir = Path(project_root) / 'frontend/svetu/src'
    
    # Файлы с useTranslations - берем из предыдущего анализа
    components_with_translations = [
        'components/Header.tsx',
        'app/[locale]/create-listing/CreateListingClient.tsx', 
        'components/marketplace/HomePage.tsx',
        'components/AuthButton.tsx',
        'components/search/CategorySelector.tsx',
        'components/cart/ShoppingCartModal.tsx',
        'components/Chat/ChatWindow.tsx',
        'app/[locale]/admin/layout-client.tsx',
        'components/storefronts/create/steps/BasicInfoStep.tsx',
        'components/GIS/LocationPicker.tsx'
    ]
    
    analysis = {
        'components_analyzed': [],
        'translation_patterns': {},
        'potential_issues': [],
        'namespace_usage': {}
    }
    
    namespace_counter = {}
    
    for component_path in components_with_translations:
        full_path = src_dir / component_path
        
        if not full_path.exists():
            analysis['potential_issues'].append({
                'type': 'file_not_found',
                'file': component_path,
                'severity': 'low'
            })
            continue
            
        try:
            with open(full_path, 'r', encoding='utf-8') as f:
                content = f.read()
                
            component_analysis = analyze_single_component(content, component_path)
            analysis['components_analyzed'].append(component_analysis)
            
            # Считаем использование namespace
            for namespace in component_analysis['namespaces']:
                namespace_counter[namespace] = namespace_counter.get(namespace, 0) + 1
                
        except Exception as e:
            analysis['potential_issues'].append({
                'type': 'analysis_error',
                'file': component_path,
                'error': str(e),
                'severity': 'medium'
            })
    
    analysis['namespace_usage'] = dict(sorted(namespace_counter.items(), key=lambda x: x[1], reverse=True))
    
    return analysis

def analyze_single_component(content: str, file_path: str) -> Dict:
    """Анализ одного компонента"""
    
    # Паттерны для поиска
    use_translations_pattern = re.compile(r"useTranslations\(['\"]([^'\"]+)['\"]\)")
    translation_calls_pattern = re.compile(r"(?:^|\s)t\(['\"]([^'\"]+)['\"]")
    const_t_pattern = re.compile(r"const\s+(\w+)\s*=\s*useTranslations\(['\"]([^'\"]+)['\"]\)")
    
    analysis = {
        'file': file_path,
        'namespaces': [],
        'translation_calls': [],
        'const_names': [],
        'potential_issues': []
    }
    
    # Ищем useTranslations
    for match in use_translations_pattern.finditer(content):
        namespace = match.group(1)
        analysis['namespaces'].append(namespace)
        
        # Проверяем на вложенные пути
        if '.' in namespace:
            analysis['potential_issues'].append({
                'type': 'nested_namespace',
                'namespace': namespace,
                'line': content[:match.start()].count('\n') + 1
            })
    
    # Ищем константы с переводами 
    for match in const_t_pattern.finditer(content):
        const_name = match.group(1)
        namespace = match.group(2)
        analysis['const_names'].append({'name': const_name, 'namespace': namespace})
    
    # Ищем вызовы t()
    translation_calls = []
    for match in translation_calls_pattern.finditer(content):
        key = match.group(1)
        line = content[:match.start()].count('\n') + 1
        translation_calls.append({'key': key, 'line': line})
    
    analysis['translation_calls'] = translation_calls[:20]  # Ограничиваем вывод
    
    return analysis

def main():
    project_root = '/data/hostel-booking-system'
    
    print("🔍 Анализирую выборочные компоненты с переводами...")
    
    analysis = analyze_component_translations(project_root)
    
    print(f"\n📊 РЕЗУЛЬТАТЫ АНАЛИЗА КОМПОНЕНТОВ:")
    print(f"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
    print(f"📁 Проанализировано компонентов: {len(analysis['components_analyzed'])}")
    print(f"⚠️  Потенциальных проблем: {len(analysis['potential_issues'])}")
    print()
    
    # Топ используемых namespace
    print("🏷️  ТОП-10 ИСПОЛЬЗУЕМЫХ NAMESPACE:")
    for i, (namespace, count) in enumerate(list(analysis['namespace_usage'].items())[:10], 1):
        print(f"   {i:2}. {namespace:<20} ({count} раз)")
    print()
    
    # Детальный анализ каждого компонента
    print("🔍 ДЕТАЛЬНЫЙ АНАЛИЗ КОМПОНЕНТОВ:")
    for comp in analysis['components_analyzed'][:10]:  # Показываем первые 10
        print(f"\n📄 {comp['file']}")
        print(f"   Namespace: {', '.join(comp['namespaces']) if comp['namespaces'] else 'НЕТ'}")
        print(f"   Вызовы t(): {len(comp['translation_calls'])}")
        
        if comp['potential_issues']:
            print(f"   ⚠️  Проблемы: {len(comp['potential_issues'])}")
            for issue in comp['potential_issues']:
                print(f"      - {issue['type']}: {issue.get('namespace', 'N/A')}")
    
    # Сохраняем результаты
    output_file = Path(project_root) / 'spec-kit/translation-audit/reports/component_analysis.txt'
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write("АНАЛИЗ ИСПОЛЬЗОВАНИЯ ПЕРЕВОДОВ В КОМПОНЕНТАХ\n")
        f.write("=" * 50 + "\n\n")
        
        f.write(f"Проанализировано компонентов: {len(analysis['components_analyzed'])}\n")
        f.write(f"Потенциальных проблем: {len(analysis['potential_issues'])}\n\n")
        
        f.write("ТОП ИСПОЛЬЗУЕМЫХ NAMESPACE:\n")
        for namespace, count in analysis['namespace_usage'].items():
            f.write(f"  {namespace}: {count}\n")
        f.write("\n")
        
        f.write("ДЕТАЛЬНЫЙ АНАЛИЗ КОМПОНЕНТОВ:\n")
        for comp in analysis['components_analyzed']:
            f.write(f"\nКомпонент: {comp['file']}\n")
            f.write(f"  Namespace: {', '.join(comp['namespaces'])}\n")
            f.write(f"  Вызовы переводов: {len(comp['translation_calls'])}\n")
            if comp['potential_issues']:
                f.write(f"  Проблемы: {comp['potential_issues']}\n")
    
    print(f"\n✅ Результаты сохранены в {output_file}")

if __name__ == '__main__':
    main()