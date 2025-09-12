#!/usr/bin/env python3
"""
Анализатор системы переводов для Sve Tu проекта
Проверяет консистентность, находит ошибки и генерирует отчет
"""

import json
import os
import re
from pathlib import Path
from typing import Dict, List, Set, Any
from collections import defaultdict

class TranslationAnalyzer:
    def __init__(self, project_root: str):
        self.project_root = Path(project_root)
        self.messages_dir = self.project_root / 'frontend/svetu/src/messages'
        self.src_dir = self.project_root / 'frontend/svetu/src'
        
        self.languages = ['en', 'ru', 'sr']
        self.report = {
            'summary': {},
            'critical_errors': [],
            'missing_translations': [],
            'incorrect_paths': [],
            'duplicate_keys': [],
            'unused_keys': [],
            'structure_issues': [],
            'recommendations': []
        }
        
    def analyze(self) -> Dict:
        """Основная функция анализа"""
        print("🔍 Начинаю анализ системы переводов...")
        
        # 1. Анализ структуры файлов
        self._analyze_file_structure()
        
        # 2. Проверка консистентности ключей между языками
        self._check_key_consistency()
        
        # 3. Анализ использования переводов в коде
        self._analyze_code_usage()
        
        # 4. Поиск структурных проблем
        self._check_structure_issues()
        
        # 5. Генерация рекомендаций
        self._generate_recommendations()
        
        return self.report
        
    def _analyze_file_structure(self):
        """Анализ структуры файлов переводов"""
        print("📁 Анализирую структуру файлов...")
        
        files_by_lang = {}
        for lang in self.languages:
            lang_dir = self.messages_dir / lang
            if not lang_dir.exists():
                self.report['critical_errors'].append(f"Отсутствует директория для языка: {lang}")
                continue
                
            files = [f.stem for f in lang_dir.glob('*.json') if f.stem != 'index']
            files_by_lang[lang] = set(files)
        
        # Проверяем одинаковость файлов между языками
        all_modules = set()
        for files in files_by_lang.values():
            all_modules.update(files)
            
        for lang in self.languages:
            lang_files = files_by_lang.get(lang, set())
            missing = all_modules - lang_files
            extra = lang_files - all_modules
            
            if missing:
                for module in missing:
                    self.report['missing_translations'].append({
                        'type': 'missing_module_file',
                        'language': lang,
                        'module': module,
                        'severity': 'critical'
                    })
                    
            if extra:
                for module in extra:
                    self.report['structure_issues'].append({
                        'type': 'extra_module_file',
                        'language': lang,
                        'module': module,
                        'severity': 'warning'
                    })
        
        self.report['summary']['total_modules'] = len(all_modules)
        self.report['summary']['languages'] = len(self.languages)
        
    def _check_key_consistency(self):
        """Проверка консистентности ключей между языками"""
        print("🔑 Проверяю консистентность ключей...")
        
        modules_data = defaultdict(dict)
        
        # Загружаем все модули
        for lang in self.languages:
            lang_dir = self.messages_dir / lang
            if not lang_dir.exists():
                continue
                
            for json_file in lang_dir.glob('*.json'):
                if json_file.stem == 'index':
                    continue
                    
                try:
                    with open(json_file, 'r', encoding='utf-8') as f:
                        data = json.load(f)
                        modules_data[json_file.stem][lang] = data
                except json.JSONDecodeError as e:
                    self.report['critical_errors'].append({
                        'type': 'json_parse_error',
                        'file': str(json_file),
                        'error': str(e)
                    })
                except Exception as e:
                    self.report['critical_errors'].append({
                        'type': 'file_read_error',
                        'file': str(json_file),
                        'error': str(e)
                    })
        
        # Проверяем консистентность ключей
        total_missing = 0
        for module_name, lang_data in modules_data.items():
            base_lang = 'en'  # Используем английский как базовый
            if base_lang not in lang_data:
                continue
                
            base_keys = self._get_all_keys(lang_data[base_lang])
            
            for lang in self.languages:
                if lang == base_lang or lang not in lang_data:
                    continue
                    
                lang_keys = self._get_all_keys(lang_data[lang])
                
                missing = base_keys - lang_keys
                extra = lang_keys - base_keys
                
                for key in missing:
                    self.report['missing_translations'].append({
                        'type': 'missing_key',
                        'module': module_name,
                        'language': lang,
                        'key': key,
                        'severity': 'high'
                    })
                    total_missing += 1
                    
                for key in extra:
                    self.report['structure_issues'].append({
                        'type': 'extra_key',
                        'module': module_name,
                        'language': lang,
                        'key': key,
                        'severity': 'low'
                    })
        
        self.report['summary']['missing_keys_total'] = total_missing
        
    def _get_all_keys(self, data: Dict, prefix: str = '') -> Set[str]:
        """Рекурсивно получает все ключи из вложенного объекта"""
        keys = set()
        
        for key, value in data.items():
            current_key = f"{prefix}.{key}" if prefix else key
            keys.add(current_key)
            
            if isinstance(value, dict):
                keys.update(self._get_all_keys(value, current_key))
                
        return keys
        
    def _analyze_code_usage(self):
        """Анализирует использование переводов в коде"""
        print("💻 Анализирую использование переводов в коде...")
        
        # Паттерны для поиска
        use_translations_pattern = re.compile(r"useTranslations\(['\"]([^'\"]+)['\"]\)")
        translation_call_pattern = re.compile(r"t\(['\"]([^'\"]+)['\"]")
        
        used_modules = set()
        used_keys = set()
        incorrect_paths = []
        
        # Проходим по всем TypeScript/JavaScript файлам
        for file_path in self.src_dir.rglob('*.{tsx,ts,jsx,js}'):
            if file_path.is_file():
                try:
                    with open(file_path, 'r', encoding='utf-8') as f:
                        content = f.read()
                        
                    # Ищем useTranslations
                    for match in use_translations_pattern.finditer(content):
                        module_path = match.group(1)
                        
                        # Проверяем на вложенные пути (неправильные)
                        if '.' in module_path:
                            incorrect_paths.append({
                                'file': str(file_path.relative_to(self.project_root)),
                                'line': content[:match.start()].count('\n') + 1,
                                'incorrect_path': module_path,
                                'suggested_fix': f"useTranslations('{module_path.split('.')[0]}')"
                            })
                        else:
                            used_modules.add(module_path)
                    
                    # Ищем вызовы t()
                    for match in translation_call_pattern.finditer(content):
                        key = match.group(1)
                        used_keys.add(key)
                        
                except Exception as e:
                    self.report['structure_issues'].append({
                        'type': 'file_analysis_error',
                        'file': str(file_path),
                        'error': str(e)
                    })
        
        self.report['incorrect_paths'] = incorrect_paths
        self.report['summary']['used_modules'] = len(used_modules)
        self.report['summary']['used_keys'] = len(used_keys)
        self.report['summary']['incorrect_paths'] = len(incorrect_paths)
        
    def _check_structure_issues(self):
        """Проверяет структурные проблемы в JSON файлах"""
        print("🏗️ Проверяю структурные проблемы...")
        
        structure_issues = 0
        
        for lang in self.languages:
            lang_dir = self.messages_dir / lang
            if not lang_dir.exists():
                continue
                
            for json_file in lang_dir.glob('*.json'):
                if json_file.stem == 'index':
                    continue
                    
                try:
                    with open(json_file, 'r', encoding='utf-8') as f:
                        data = json.load(f)
                        
                    # Проверяем структуру
                    issues = self._validate_json_structure(data, json_file.stem, lang)
                    structure_issues += len(issues)
                    self.report['structure_issues'].extend(issues)
                    
                except Exception:
                    # Ошибки уже добавлены в _check_key_consistency
                    pass
                    
        self.report['summary']['structure_issues'] = structure_issues
        
    def _validate_json_structure(self, data: Dict, module: str, lang: str) -> List[Dict]:
        """Валидирует структуру JSON"""
        issues = []
        
        def check_recursive(obj, path=""):
            if isinstance(obj, dict):
                for key, value in obj.items():
                    current_path = f"{path}.{key}" if path else key
                    
                    if isinstance(value, dict):
                        # Проверяем пустые объекты
                        if not value:
                            issues.append({
                                'type': 'empty_object',
                                'module': module,
                                'language': lang,
                                'path': current_path
                            })
                        else:
                            check_recursive(value, current_path)
                    elif isinstance(value, str):
                        # Проверяем пустые строки
                        if not value.strip():
                            issues.append({
                                'type': 'empty_string',
                                'module': module,
                                'language': lang,
                                'path': current_path
                            })
                    
        check_recursive(data)
        return issues
        
    def _generate_recommendations(self):
        """Генерирует рекомендации по улучшению"""
        print("💡 Генерирую рекомендации...")
        
        recommendations = []
        
        # Критические ошибки
        if self.report['critical_errors']:
            recommendations.append({
                'priority': 'critical',
                'category': 'Критические ошибки',
                'description': f"Найдено {len(self.report['critical_errors'])} критических ошибок, которые блокируют работу системы переводов",
                'action': "Немедленно исправить все критические ошибки"
            })
        
        # Недостающие переводы
        if self.report['missing_translations']:
            high_priority = [t for t in self.report['missing_translations'] if t.get('severity') == 'high']
            recommendations.append({
                'priority': 'high',
                'category': 'Недостающие переводы',
                'description': f"Найдено {len(high_priority)} отсутствующих переводов высокого приоритета",
                'action': "Добавить недостающие переводы во все языковые файлы"
            })
        
        # Неправильные пути
        if self.report['incorrect_paths']:
            recommendations.append({
                'priority': 'medium',
                'category': 'Неправильные пути переводов',
                'description': f"Найдено {len(self.report['incorrect_paths'])} случаев использования вложенных путей в useTranslations",
                'action': "Запустить автоматический скрипт исправления или исправить вручную"
            })
        
        # Структурные проблемы
        if self.report['structure_issues']:
            recommendations.append({
                'priority': 'low',
                'category': 'Структурные проблемы',
                'description': f"Найдено {len(self.report['structure_issues'])} структурных проблем",
                'action': "Очистить пустые объекты и строки в файлах переводов"
            })
        
        self.report['recommendations'] = recommendations

def main():
    project_root = '/data/hostel-booking-system'
    analyzer = TranslationAnalyzer(project_root)
    report = analyzer.analyze()
    
    # Сохраняем отчет
    output_file = Path(project_root) / 'spec-kit/translation-audit/reports/analysis_results.json'
    output_file.parent.mkdir(parents=True, exist_ok=True)
    
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(report, f, ensure_ascii=False, indent=2)
    
    print(f"\n✅ Анализ завершен. Результаты сохранены в {output_file}")
    print(f"📊 Найдено проблем:")
    print(f"   - Критические ошибки: {len(report['critical_errors'])}")
    print(f"   - Недостающие переводы: {len(report['missing_translations'])}")
    print(f"   - Неправильные пути: {len(report['incorrect_paths'])}")
    print(f"   - Структурные проблемы: {len(report['structure_issues'])}")

if __name__ == '__main__':
    main()