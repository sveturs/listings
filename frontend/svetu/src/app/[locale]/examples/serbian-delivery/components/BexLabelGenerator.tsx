'use client';

import React, { useState } from 'react';
import {
  DocumentTextIcon,
  PrinterIcon,
  QrCodeIcon,
  ArrowDownTrayIcon,
  CheckCircleIcon,
  ClockIcon,
  DocumentDuplicateIcon,
  Squares2X2Icon,
} from '@heroicons/react/24/outline';

interface LabelSettings {
  shipmentId: string;
  pageSize: 'A4' | 'A6';
  position: 1 | 2 | 3 | 4;
  parcelNo: number;
  format: 'pdf' | 'png';
}

interface BulkLabel {
  id: string;
  shipmentId: string;
  sender: string;
  receiver: string;
  status: 'pending' | 'generating' | 'completed' | 'error';
}

export default function BexLabelGenerator() {
  const [labelSettings, setLabelSettings] = useState<LabelSettings>({
    shipmentId: '180923112',
    pageSize: 'A4',
    position: 1,
    parcelNo: 0,
    format: 'pdf',
  });

  const [bulkLabels, setBulkLabels] = useState<BulkLabel[]>([
    {
      id: '1',
      shipmentId: '170123456',
      sender: 'Марко Петровић',
      receiver: 'Ана Јовановић',
      status: 'completed',
    },
    {
      id: '2',
      shipmentId: '170123457',
      sender: 'Милан Николић',
      receiver: 'Јелена Стојановић',
      status: 'completed',
    },
    {
      id: '3',
      shipmentId: '170123458',
      sender: 'Петар Марковић',
      receiver: 'Милица Ђорђевић',
      status: 'generating',
    },
    {
      id: '4',
      shipmentId: '170123459',
      sender: 'Стефан Јовановић',
      receiver: 'Тамара Павловић',
      status: 'pending',
    },
    {
      id: '5',
      shipmentId: '170123460',
      sender: 'Никола Милић',
      receiver: 'Марија Савић',
      status: 'pending',
    },
  ]);

  const [showLabel, setShowLabel] = useState(false);
  const [selectedTemplate, setSelectedTemplate] = useState('standard');

  const positions = [
    { value: 1, label: 'Горе лево', icon: '↖️' },
    { value: 2, label: 'Горе десно', icon: '↗️' },
    { value: 3, label: 'Доле лево', icon: '↙️' },
    { value: 4, label: 'Доле десно', icon: '↘️' },
  ];

  const templates = [
    {
      id: 'standard',
      name: 'Стандардна',
      description: 'Основна адресница са баркодом',
    },
    {
      id: 'express',
      name: 'Express',
      description: 'Са великим Express ознакама',
    },
    {
      id: 'fragile',
      name: 'Ломљиво',
      description: 'Са упозорењима за ломљиво',
    },
    { id: 'cod', name: 'Откупнина', description: 'Истакнут износ откупнине' },
  ];

  const handleGenerateLabel = () => {
    setShowLabel(true);
    setTimeout(() => {
      setShowLabel(false);
    }, 3000);
  };

  const handleBulkGenerate = () => {
    setBulkLabels((prev) =>
      prev.map((label) => ({
        ...label,
        status: label.status === 'pending' ? 'generating' : label.status,
      }))
    );

    setTimeout(() => {
      setBulkLabels((prev) =>
        prev.map((label) => ({
          ...label,
          status: label.status === 'generating' ? 'completed' : label.status,
        }))
      );
    }, 2000);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'text-success';
      case 'generating':
        return 'text-warning';
      case 'error':
        return 'text-error';
      default:
        return 'text-base-content/60';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <CheckCircleIcon className="w-5 h-5" />;
      case 'generating':
        return <ClockIcon className="w-5 h-5 animate-spin" />;
      default:
        return <DocumentTextIcon className="w-5 h-5" />;
    }
  };

  return (
    <div className="space-y-6">
      {/* Single Label Generation */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h3 className="card-title">
            <DocumentTextIcon className="w-6 h-6" />
            Појединачна адресница
          </h3>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">Број пошиљке</span>
              </label>
              <input
                type="text"
                className="input input-bordered"
                value={labelSettings.shipmentId}
                onChange={(e) =>
                  setLabelSettings({
                    ...labelSettings,
                    shipmentId: e.target.value,
                  })
                }
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Величина папира</span>
              </label>
              <select
                className="select select-bordered"
                value={labelSettings.pageSize}
                onChange={(e) =>
                  setLabelSettings({
                    ...labelSettings,
                    pageSize: e.target.value as 'A4' | 'A6',
                  })
                }
              >
                <option value="A4">A4 (210 × 297 mm)</option>
                <option value="A6">A6 (105 × 148 mm)</option>
              </select>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Позиција на страници</span>
              </label>
              <div className="grid grid-cols-2 gap-2">
                {positions.map((pos) => (
                  <button
                    key={pos.value}
                    onClick={() =>
                      setLabelSettings({
                        ...labelSettings,
                        position: pos.value as 1 | 2 | 3 | 4,
                      })
                    }
                    className={`btn ${
                      labelSettings.position === pos.value
                        ? 'btn-primary'
                        : 'btn-outline'
                    }`}
                  >
                    {pos.icon} {pos.label}
                  </button>
                ))}
              </div>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Број пакета</span>
                <span className="label-text-alt">0 за све</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                min="0"
                value={labelSettings.parcelNo}
                onChange={(e) =>
                  setLabelSettings({
                    ...labelSettings,
                    parcelNo: Number(e.target.value),
                  })
                }
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Формат</span>
              </label>
              <select
                className="select select-bordered"
                value={labelSettings.format}
                onChange={(e) =>
                  setLabelSettings({
                    ...labelSettings,
                    format: e.target.value as 'pdf' | 'png',
                  })
                }
              >
                <option value="pdf">PDF</option>
                <option value="png">PNG слика</option>
              </select>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Шаблон</span>
              </label>
              <select
                className="select select-bordered"
                value={selectedTemplate}
                onChange={(e) => setSelectedTemplate(e.target.value)}
              >
                {templates.map((template) => (
                  <option key={template.id} value={template.id}>
                    {template.name} - {template.description}
                  </option>
                ))}
              </select>
            </div>
          </div>

          {showLabel && (
            <div className="alert alert-success mt-4">
              <CheckCircleIcon className="w-6 h-6" />
              <div>
                <h3 className="font-bold">Адресница генерисана!</h3>
                <div className="text-xs">Спремна за преузимање</div>
              </div>
            </div>
          )}

          <div className="card-actions justify-end mt-4">
            <button onClick={handleGenerateLabel} className="btn btn-primary">
              <QrCodeIcon className="w-5 h-5" />
              Генериши адресницу
            </button>
            <button className="btn btn-outline">
              <ArrowDownTrayIcon className="w-5 h-5" />
              Преузми
            </button>
            <button className="btn btn-outline">
              <PrinterIcon className="w-5 h-5" />
              Штампај
            </button>
          </div>
        </div>
      </div>

      {/* Label Preview */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h3 className="card-title mb-4">
            <Squares2X2Icon className="w-6 h-6" />
            Преглед адреснице
          </h3>

          <div className="bg-white rounded-lg p-6 border-2 border-dashed border-base-300">
            <div className="max-w-md mx-auto space-y-4">
              {/* Barcode */}
              <div className="flex justify-center">
                <div className="space-y-2">
                  <div className="bg-black h-16 w-64 rounded flex items-center justify-center">
                    <div className="grid grid-cols-12 gap-0.5">
                      {[...Array(12)].map((_, i) => (
                        <div
                          key={i}
                          className={`${
                            i % 2 === 0 ? 'bg-white' : 'bg-transparent'
                          } ${i % 3 === 0 ? 'w-1' : 'w-0.5'} h-12`}
                        />
                      ))}
                    </div>
                  </div>
                  <div className="text-center font-mono text-lg font-bold">
                    {labelSettings.shipmentId}
                  </div>
                </div>
              </div>

              {/* Express Label (if selected) */}
              {selectedTemplate === 'express' && (
                <div className="bg-red-600 text-white text-center py-2 rounded font-bold text-xl">
                  EXPRESS ДОСТАВА
                </div>
              )}

              {/* Fragile Warning (if selected) */}
              {selectedTemplate === 'fragile' && (
                <div className="bg-orange-500 text-white text-center py-2 rounded font-bold">
                  ⚠️ ПАЖЉИВО - ЛОМЉИВО ⚠️
                </div>
              )}

              <div className="border-t-2 border-black pt-4">
                <div className="grid grid-cols-2 gap-4">
                  {/* Sender */}
                  <div>
                    <p className="text-xs font-bold mb-1">ПОШИЉАЛАЦ:</p>
                    <p className="text-sm font-semibold">Марко Петровић</p>
                    <p className="text-xs">БЕОГРАД</p>
                    <p className="text-xs">Кнез Михаилова 25/5</p>
                    <p className="text-xs">📞 064/851-6928</p>
                  </div>

                  {/* Service Info */}
                  <div className="text-right">
                    <div className="inline-block bg-black text-white px-2 py-1 rounded text-xs">
                      Стандардна достава
                    </div>
                    <p className="text-xs mt-2">
                      Датум: {new Date().toLocaleDateString('sr-RS')}
                    </p>
                    <p className="text-xs">Пакета: 1/1</p>
                  </div>
                </div>
              </div>

              <div className="border-t-2 border-black pt-4">
                {/* Receiver */}
                <div className="bg-gray-100 p-3 rounded">
                  <p className="text-xs font-bold mb-1">ПРИМАЛАЦ:</p>
                  <p className="text-lg font-bold">Ана Јовановић</p>
                  <p className="text-sm">НОВИ САД</p>
                  <p className="text-sm">Дунавска 10</p>
                  <p className="text-sm">📞 065/123-4567</p>
                </div>
              </div>

              {/* COD Amount (if selected) */}
              {selectedTemplate === 'cod' && (
                <div className="bg-green-600 text-white p-3 rounded text-center">
                  <p className="text-xs">ОТКУПНИНА</p>
                  <p className="text-2xl font-bold">2,500 РСД</p>
                </div>
              )}

              {/* Comments */}
              <div className="border-t pt-3">
                <p className="text-xs font-bold">НАПОМЕНА:</p>
                <p className="text-sm">Позвати пре доставе</p>
              </div>

              {/* Footer */}
              <div className="grid grid-cols-2 gap-2 text-xs border-t pt-3">
                <div>
                  <span className="font-bold">Осигурање:</span> 5,000 РСД
                </div>
                <div>
                  <span className="font-bold">Тежина:</span> 0.5 kg
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Bulk Label Generation */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <div className="flex justify-between items-center mb-4">
            <h3 className="card-title">
              <DocumentDuplicateIcon className="w-6 h-6" />
              Масовно генерисање адресница
            </h3>
            <div className="badge badge-primary badge-lg">
              {bulkLabels.filter((l) => l.status === 'completed').length} /{' '}
              {bulkLabels.length} завршено
            </div>
          </div>

          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th>
                    <label>
                      <input type="checkbox" className="checkbox" />
                    </label>
                  </th>
                  <th>Број пошиљке</th>
                  <th>Пошиљалац</th>
                  <th>Прималац</th>
                  <th>Статус</th>
                  <th>Акције</th>
                </tr>
              </thead>
              <tbody>
                {bulkLabels.map((label) => (
                  <tr key={label.id}>
                    <th>
                      <label>
                        <input type="checkbox" className="checkbox" />
                      </label>
                    </th>
                    <td className="font-mono">{label.shipmentId}</td>
                    <td>{label.sender}</td>
                    <td>{label.receiver}</td>
                    <td>
                      <div
                        className={`flex items-center gap-2 ${getStatusColor(label.status)}`}
                      >
                        {getStatusIcon(label.status)}
                        <span className="capitalize">{label.status}</span>
                      </div>
                    </td>
                    <td>
                      {label.status === 'completed' && (
                        <button className="btn btn-ghost btn-xs">
                          <ArrowDownTrayIcon className="w-4 h-4" />
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Progress */}
          <div className="mt-4">
            <div className="flex justify-between text-sm mb-2">
              <span>Прогрес генерисања</span>
              <span>
                {Math.round(
                  (bulkLabels.filter((l) => l.status === 'completed').length /
                    bulkLabels.length) *
                    100
                )}
                %
              </span>
            </div>
            <progress
              className="progress progress-primary w-full"
              value={bulkLabels.filter((l) => l.status === 'completed').length}
              max={bulkLabels.length}
            />
          </div>

          <div className="card-actions justify-between mt-4">
            <div className="flex gap-2">
              <select className="select select-bordered select-sm">
                <option>A4 формат</option>
                <option>A6 формат</option>
              </select>
              <select className="select select-bordered select-sm">
                <option>PDF</option>
                <option>PNG</option>
              </select>
            </div>
            <div className="flex gap-2">
              <button onClick={handleBulkGenerate} className="btn btn-primary">
                <DocumentDuplicateIcon className="w-5 h-5" />
                Генериши све
              </button>
              <button className="btn btn-outline">
                <ArrowDownTrayIcon className="w-5 h-5" />
                Преузми ZIP
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* API Info */}
      <div className="collapse collapse-arrow bg-base-200">
        <input type="checkbox" />
        <div className="collapse-title text-xl font-medium">
          📡 API Информације
        </div>
        <div className="collapse-content">
          <div className="mockup-code">
            <pre data-prefix="$">
              <code className="text-xs">
                GET https://api.bex.rs:62502/getLabelWithProperties
              </code>
            </pre>
            <pre data-prefix="?">
              <code className="text-xs">
                pageSize={labelSettings.pageSize === 'A4' ? '4' : '6'}
              </code>
            </pre>
            <pre data-prefix="&">
              <code className="text-xs">
                pagePosition={labelSettings.position}
              </code>
            </pre>
            <pre data-prefix="&">
              <code className="text-xs">
                shipmentId={labelSettings.shipmentId}
              </code>
            </pre>
            <pre data-prefix="&">
              <code className="text-xs">parcelNo={labelSettings.parcelNo}</code>
            </pre>
          </div>

          <div className="mt-4">
            <h4 className="font-bold mb-2">Response:</h4>
            <div className="mockup-code">
              <pre>
                <code className="text-xs">
                  {JSON.stringify(
                    {
                      state: true,
                      shipmentId: labelSettings.shipmentId,
                      parcelNo: labelSettings.parcelNo,
                      parcelLabel: 'JVBERi0xLjQKJeLjz9M...', // base64
                      err: '',
                    },
                    null,
                    2
                  )}
                </code>
              </pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
