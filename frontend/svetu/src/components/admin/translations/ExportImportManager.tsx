'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import {
  translationAdminApi,
  ExportRequest,
  ImportRequest,
} from '@/services/translationAdminApi';
import {
  DocumentArrowDownIcon,
  DocumentArrowUpIcon,
  DocumentTextIcon,
  TableCellsIcon,
  CodeBracketIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
} from '@heroicons/react/24/outline';

export default function ExportImportManager() {
  const _t = useTranslations('admin.translations');

  // Export state
  const [exportFormat, setExportFormat] = useState<'json' | 'csv' | 'xliff'>(
    'json'
  );
  const [exportFilters, setExportFilters] = useState({
    entity_type: '',
    language: '',
    module: '',
    only_verified: false,
    include_metadata: true,
  });
  const [isExporting, setIsExporting] = useState(false);

  // Import state
  const [importFormat, setImportFormat] = useState<'json' | 'csv' | 'xliff'>(
    'json'
  );
  const [importData, setImportData] = useState<string>('');
  const [_importFile, _setImportFile] = useState<File | null>(null);
  const [importOptions, setImportOptions] = useState({
    overwrite_existing: false,
    validate_only: false,
  });
  const [isImporting, setIsImporting] = useState(false);
  const [importResult, setImportResult] = useState<any>(null);

  const LANGUAGES = [
    { code: '', name: 'Все языки' },
    { code: 'sr', name: 'Srpski 🇷🇸' },
    { code: 'en', name: 'English 🇺🇸' },
    { code: 'ru', name: 'Русский 🇷🇺' },
  ];

  const ENTITY_TYPES = [
    { code: '', name: 'Все типы' },
    { code: 'category', name: 'Категории' },
    { code: 'attribute', name: 'Атрибуты' },
    { code: 'listing', name: 'Объявления' },
  ];

  const FORMATS = [
    {
      code: 'json' as const,
      name: 'JSON',
      icon: DocumentTextIcon,
      description: 'Структурированные данные для программной обработки',
    },
    {
      code: 'csv' as const,
      name: 'CSV',
      icon: TableCellsIcon,
      description: 'Табличный формат для Excel и Google Sheets',
    },
    {
      code: 'xliff' as const,
      name: 'XLIFF',
      icon: CodeBracketIcon,
      description: 'Индустриальный стандарт для переводов',
    },
  ];

  const handleExport = async () => {
    try {
      setIsExporting(true);

      const request: ExportRequest = {
        format: exportFormat,
        entity_type: exportFilters.entity_type || undefined,
        language: exportFilters.language || undefined,
        module: exportFilters.module || undefined,
        only_verified: exportFilters.only_verified,
        include_metadata: exportFilters.include_metadata,
      };

      const result = await translationAdminApi.export(request);

      // Handle different export formats
      if (exportFormat === 'json') {
        // For JSON, create a downloadable file
        const blob = new Blob([JSON.stringify(result, null, 2)], {
          type: 'application/json',
        });
        downloadFile(blob, `translations_export.json`);
      } else {
        // For CSV and XLIFF, result is already a Blob
        const extension = exportFormat === 'csv' ? 'csv' : 'xlf';
        downloadFile(result, `translations_export.${extension}`);
      }
    } catch (error) {
      console.error('Export failed:', error);
      alert('Ошибка при экспорте данных');
    } finally {
      setIsExporting(false);
    }
  };

  const downloadFile = (blob: Blob, filename: string) => {
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  };

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    _setImportFile(file);

    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target?.result as string;
      setImportData(content);
    };
    reader.readAsText(file);
  };

  const handleImport = async () => {
    if (!importData.trim()) {
      alert('Пожалуйста, выберите файл или введите данные');
      return;
    }

    try {
      setIsImporting(true);

      let data: any;
      try {
        if (importFormat === 'json') {
          data = JSON.parse(importData);
        } else if (importFormat === 'csv') {
          // For CSV, we'll send the raw text and let backend parse it
          data = importData;
        } else if (importFormat === 'xliff') {
          // For XLIFF, we'll send the raw XML and let backend parse it
          data = importData;
        }
      } catch {
        alert('Ошибка при парсинге данных. Проверьте формат файла.');
        return;
      }

      const request: ImportRequest = {
        format: importFormat,
        data,
        overwrite_existing: importOptions.overwrite_existing,
        validate_only: importOptions.validate_only,
        metadata: {
          imported_at: new Date().toISOString(),
          imported_by: 'admin', // TODO: get real user
        },
      };

      const result = await translationAdminApi.import(request);
      setImportResult(result);
    } catch (error) {
      console.error('Import failed:', error);
      alert('Ошибка при импорте данных');
    } finally {
      setIsImporting(false);
    }
  };

  const clearImport = () => {
    setImportData('');
    _setImportFile(null);
    setImportResult(null);

    // Clear file input
    const fileInput = document.getElementById(
      'import-file'
    ) as HTMLInputElement;
    if (fileInput) fileInput.value = '';
  };

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold">Экспорт и Импорт</h2>
        <p className="text-base-content/60 mt-1">
          Управление переводами через файлы различных форматов
        </p>
      </div>

      <div className="grid lg:grid-cols-2 gap-8">
        {/* Export Section */}
        <div className="space-y-6">
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <h3 className="card-title text-lg flex items-center gap-2">
                <DocumentArrowDownIcon className="h-5 w-5" />
                Экспорт переводов
              </h3>

              {/* Format Selection */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    Формат экспорта
                  </span>
                </label>
                <div className="grid grid-cols-1 gap-2">
                  {FORMATS.map((format) => (
                    <label
                      key={format.code}
                      className="label cursor-pointer justify-start"
                    >
                      <input
                        type="radio"
                        name="export-format"
                        className="radio mr-3"
                        checked={exportFormat === format.code}
                        onChange={() => setExportFormat(format.code)}
                      />
                      <div>
                        <div className="flex items-center gap-2">
                          <format.icon className="h-4 w-4" />
                          <span className="font-medium">{format.name}</span>
                        </div>
                        <div className="text-xs text-base-content/60">
                          {format.description}
                        </div>
                      </div>
                    </label>
                  ))}
                </div>
              </div>

              {/* Export Filters */}
              <div className="divider">Фильтры</div>

              <div className="grid grid-cols-1 gap-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Тип сущности</span>
                  </label>
                  <select
                    className="select select-bordered"
                    value={exportFilters.entity_type}
                    onChange={(e) =>
                      setExportFilters((prev) => ({
                        ...prev,
                        entity_type: e.target.value,
                      }))
                    }
                  >
                    {ENTITY_TYPES.map((type) => (
                      <option key={type.code} value={type.code}>
                        {type.name}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Язык</span>
                  </label>
                  <select
                    className="select select-bordered"
                    value={exportFilters.language}
                    onChange={(e) =>
                      setExportFilters((prev) => ({
                        ...prev,
                        language: e.target.value,
                      }))
                    }
                  >
                    {LANGUAGES.map((lang) => (
                      <option key={lang.code} value={lang.code}>
                        {lang.name}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Модуль</span>
                  </label>
                  <input
                    type="text"
                    placeholder="например: marketplace, common"
                    className="input input-bordered"
                    value={exportFilters.module}
                    onChange={(e) =>
                      setExportFilters((prev) => ({
                        ...prev,
                        module: e.target.value,
                      }))
                    }
                  />
                </div>
              </div>

              {/* Export Options */}
              <div className="divider">Опции</div>

              <div className="form-control">
                <label className="label cursor-pointer">
                  <span className="label-text">
                    Только проверенные переводы
                  </span>
                  <input
                    type="checkbox"
                    className="toggle"
                    checked={exportFilters.only_verified}
                    onChange={(e) =>
                      setExportFilters((prev) => ({
                        ...prev,
                        only_verified: e.target.checked,
                      }))
                    }
                  />
                </label>
                <label className="label cursor-pointer">
                  <span className="label-text">Включить метаданные</span>
                  <input
                    type="checkbox"
                    className="toggle"
                    checked={exportFilters.include_metadata}
                    onChange={(e) =>
                      setExportFilters((prev) => ({
                        ...prev,
                        include_metadata: e.target.checked,
                      }))
                    }
                  />
                </label>
              </div>

              {/* Export Button */}
              <button
                className="btn btn-primary w-full mt-4"
                onClick={handleExport}
                disabled={isExporting}
              >
                {isExporting ? (
                  <>
                    <span className="loading loading-spinner"></span>
                    Экспортируется...
                  </>
                ) : (
                  <>
                    <DocumentArrowDownIcon className="h-4 w-4" />
                    Экспортировать {exportFormat.toUpperCase()}
                  </>
                )}
              </button>
            </div>
          </div>
        </div>

        {/* Import Section */}
        <div className="space-y-6">
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <h3 className="card-title text-lg flex items-center gap-2">
                <DocumentArrowUpIcon className="h-5 w-5" />
                Импорт переводов
              </h3>

              {/* Format Selection */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">Формат импорта</span>
                </label>
                <div className="grid grid-cols-1 gap-2">
                  {FORMATS.map((format) => (
                    <label
                      key={format.code}
                      className="label cursor-pointer justify-start"
                    >
                      <input
                        type="radio"
                        name="import-format"
                        className="radio mr-3"
                        checked={importFormat === format.code}
                        onChange={() => setImportFormat(format.code)}
                      />
                      <div>
                        <div className="flex items-center gap-2">
                          <format.icon className="h-4 w-4" />
                          <span className="font-medium">{format.name}</span>
                        </div>
                      </div>
                    </label>
                  ))}
                </div>
              </div>

              {/* File Upload */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">Выбор файла</span>
                </label>
                <input
                  id="import-file"
                  type="file"
                  accept={
                    importFormat === 'json'
                      ? '.json'
                      : importFormat === 'csv'
                        ? '.csv'
                        : '.xlf,.xliff'
                  }
                  className="file-input file-input-bordered w-full"
                  onChange={handleFileUpload}
                />
              </div>

              {/* Manual Input */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    Или введите данные вручную
                  </span>
                </label>
                <textarea
                  className="textarea textarea-bordered h-32"
                  placeholder={`Вставьте ${importFormat.toUpperCase()} данные здесь...`}
                  value={importData}
                  onChange={(e) => setImportData(e.target.value)}
                />
              </div>

              {/* Import Options */}
              <div className="divider">Опции импорта</div>

              <div className="form-control">
                <label className="label cursor-pointer">
                  <span className="label-text">
                    Перезаписать существующие переводы
                  </span>
                  <input
                    type="checkbox"
                    className="toggle"
                    checked={importOptions.overwrite_existing}
                    onChange={(e) =>
                      setImportOptions((prev) => ({
                        ...prev,
                        overwrite_existing: e.target.checked,
                      }))
                    }
                  />
                </label>
                <label className="label cursor-pointer">
                  <span className="label-text">
                    Только валидация (без сохранения)
                  </span>
                  <input
                    type="checkbox"
                    className="toggle"
                    checked={importOptions.validate_only}
                    onChange={(e) =>
                      setImportOptions((prev) => ({
                        ...prev,
                        validate_only: e.target.checked,
                      }))
                    }
                  />
                </label>
              </div>

              {/* Import Buttons */}
              <div className="flex gap-2 mt-4">
                <button
                  className="btn btn-primary flex-1"
                  onClick={handleImport}
                  disabled={isImporting || !importData.trim()}
                >
                  {isImporting ? (
                    <>
                      <span className="loading loading-spinner"></span>
                      {importOptions.validate_only
                        ? 'Валидация...'
                        : 'Импортируется...'}
                    </>
                  ) : (
                    <>
                      <DocumentArrowUpIcon className="h-4 w-4" />
                      {importOptions.validate_only
                        ? 'Валидировать'
                        : 'Импортировать'}
                    </>
                  )}
                </button>

                <button className="btn btn-ghost" onClick={clearImport}>
                  Очистить
                </button>
              </div>

              {/* Import Results */}
              {importResult && (
                <div className="mt-6">
                  <div className="divider">Результат импорта</div>

                  <div className="grid grid-cols-3 gap-4 mb-4">
                    <div className="stat bg-success/10 rounded-lg">
                      <div className="stat-figure text-success">
                        <CheckCircleIcon className="h-8 w-8" />
                      </div>
                      <div className="stat-title">Успешно</div>
                      <div className="stat-value text-success">
                        {importResult.success}
                      </div>
                    </div>

                    <div className="stat bg-error/10 rounded-lg">
                      <div className="stat-figure text-error">
                        <ExclamationTriangleIcon className="h-8 w-8" />
                      </div>
                      <div className="stat-title">Ошибки</div>
                      <div className="stat-value text-error">
                        {importResult.failed}
                      </div>
                    </div>

                    <div className="stat bg-warning/10 rounded-lg">
                      <div className="stat-figure text-warning">
                        <ExclamationTriangleIcon className="h-8 w-8" />
                      </div>
                      <div className="stat-title">Пропущено</div>
                      <div className="stat-value text-warning">
                        {importResult.skipped}
                      </div>
                    </div>
                  </div>

                  {importResult.errors && importResult.errors.length > 0 && (
                    <div className="bg-error/10 p-4 rounded-lg">
                      <h4 className="font-semibold text-error mb-2">Ошибки:</h4>
                      <ul className="list-disc list-inside text-sm space-y-1">
                        {importResult.errors.map(
                          (error: string, index: number) => (
                            <li key={index} className="text-error">
                              {error}
                            </li>
                          )
                        )}
                      </ul>
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
