'use client';

import { useState } from 'react';
import {
  DocumentDuplicateIcon,
  PrinterIcon,
  TruckIcon,
  CheckIcon,
  XMarkIcon,
  ArrowDownTrayIcon,
  CalendarIcon,
  ExclamationTriangleIcon,
} from '@heroicons/react/24/outline';

interface Order {
  id: string;
  customerName: string;
  city: string;
  amount: string;
  status: 'pending' | 'processing' | 'ready';
  selected: boolean;
}

export default function BulkShipmentManager() {
  const [orders, setOrders] = useState<Order[]>([
    {
      id: '12345',
      customerName: 'Милан Йованович',
      city: 'Нови Сад',
      amount: '4,500 RSD',
      status: 'pending',
      selected: false,
    },
    {
      id: '12346',
      customerName: 'Ана Петрович',
      city: 'Крагуевац',
      amount: '7,200 RSD',
      status: 'pending',
      selected: false,
    },
    {
      id: '12347',
      customerName: 'Стефан Николич',
      city: 'Ниш',
      amount: '3,100 RSD',
      status: 'pending',
      selected: false,
    },
    {
      id: '12348',
      customerName: 'Елена Марич',
      city: 'Суботица',
      amount: '5,800 RSD',
      status: 'ready',
      selected: false,
    },
    {
      id: '12349',
      customerName: 'Петар Стоянович',
      city: 'Белград',
      amount: '9,200 RSD',
      status: 'processing',
      selected: false,
    },
  ]);

  const [isProcessing, setIsProcessing] = useState(false);
  const [selectAll, setSelectAll] = useState(false);

  const handleSelectAll = () => {
    const newSelectAll = !selectAll;
    setSelectAll(newSelectAll);
    setOrders(orders.map((order) => ({ ...order, selected: newSelectAll })));
  };

  const handleSelectOrder = (orderId: string) => {
    setOrders(
      orders.map((order) =>
        order.id === orderId ? { ...order, selected: !order.selected } : order
      )
    );
  };

  const selectedCount = orders.filter((o) => o.selected).length;

  const processSelectedOrders = () => {
    setIsProcessing(true);
    setTimeout(() => {
      setOrders(
        orders.map((order) =>
          order.selected
            ? { ...order, status: 'ready', selected: false }
            : order
        )
      );
      setIsProcessing(false);
      setSelectAll(false);
    }, 2000);
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'pending':
        return <span className="badge badge-warning badge-sm">Ожидает</span>;
      case 'processing':
        return <span className="badge badge-info badge-sm">Обработка</span>;
      case 'ready':
        return <span className="badge badge-success badge-sm">Готово</span>;
      default:
        return null;
    }
  };

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body p-4 sm:p-6">
        <div className="flex items-center justify-between mb-3 sm:mb-4">
          <h3 className="card-title text-base sm:text-lg">
            Массовая обработка
          </h3>
          <div className="badge badge-primary badge-sm">
            {orders.length} заказов
          </div>
        </div>

        {/* Bulk Actions */}
        {selectedCount > 0 && (
          <div className="alert alert-info mb-3 sm:mb-4">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 w-full">
              <span className="text-xs sm:text-sm">
                Выбрано заказов: {selectedCount}
              </span>
              <div className="flex gap-2">
                <button
                  className={`btn btn-xs sm:btn-sm btn-primary ${isProcessing ? 'loading' : ''}`}
                  onClick={processSelectedOrders}
                  disabled={isProcessing}
                >
                  {!isProcessing && (
                    <DocumentDuplicateIcon className="w-3 h-3 sm:w-4 sm:h-4" />
                  )}
                  <span className="hidden sm:inline">Создать отправления</span>
                  <span className="sm:hidden">Создать</span>
                </button>
                <button className="btn btn-xs sm:btn-sm btn-ghost">
                  <PrinterIcon className="w-3 h-3 sm:w-4 sm:h-4" />
                  <span className="hidden sm:inline">Печать этикеток</span>
                  <span className="sm:hidden">Печать</span>
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Orders Table - Mobile Cards */}
        <div className="block sm:hidden space-y-3">
          {orders.map((order) => (
            <div
              key={order.id}
              className={`card ${order.selected ? 'ring-2 ring-primary' : 'bg-base-200'}`}
            >
              <div className="card-body p-3">
                <div className="flex items-start justify-between">
                  <label className="flex items-start gap-2">
                    <input
                      type="checkbox"
                      className="checkbox checkbox-sm mt-0.5"
                      checked={order.selected}
                      onChange={() => handleSelectOrder(order.id)}
                    />
                    <div>
                      <div className="font-mono text-sm">#{order.id}</div>
                      <div className="font-semibold text-sm">
                        {order.customerName}
                      </div>
                      <div className="text-xs text-base-content/60">
                        {order.city}
                      </div>
                    </div>
                  </label>
                  <div className="text-right">
                    <div className="font-semibold text-sm">{order.amount}</div>
                    {getStatusBadge(order.status)}
                  </div>
                </div>
                <div className="flex justify-end mt-2">
                  {order.status === 'ready' ? (
                    <button className="btn btn-ghost btn-xs">
                      <PrinterIcon className="w-3 h-3" />
                    </button>
                  ) : (
                    <button className="btn btn-ghost btn-xs">
                      <TruckIcon className="w-3 h-3" />
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Orders Table - Desktop */}
        <div className="hidden sm:block overflow-x-auto">
          <table className="table table-zebra table-sm">
            <thead>
              <tr>
                <th>
                  <label>
                    <input
                      type="checkbox"
                      className="checkbox checkbox-sm"
                      checked={selectAll}
                      onChange={handleSelectAll}
                    />
                  </label>
                </th>
                <th>Заказ</th>
                <th>Клиент</th>
                <th>Город</th>
                <th>Сумма</th>
                <th>Статус</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {orders.map((order) => (
                <tr
                  key={order.id}
                  className={order.selected ? 'bg-primary/5' : ''}
                >
                  <td>
                    <label>
                      <input
                        type="checkbox"
                        className="checkbox checkbox-sm"
                        checked={order.selected}
                        onChange={() => handleSelectOrder(order.id)}
                      />
                    </label>
                  </td>
                  <td className="font-mono">#{order.id}</td>
                  <td>{order.customerName}</td>
                  <td>{order.city}</td>
                  <td className="font-semibold">{order.amount}</td>
                  <td>{getStatusBadge(order.status)}</td>
                  <td>
                    <div className="flex gap-1">
                      {order.status === 'ready' ? (
                        <button className="btn btn-ghost btn-xs">
                          <PrinterIcon className="w-4 h-4" />
                        </button>
                      ) : (
                        <button className="btn btn-ghost btn-xs">
                          <TruckIcon className="w-4 h-4" />
                        </button>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Quick Stats */}
        <div className="grid grid-cols-3 gap-2 sm:gap-4 mt-4 sm:mt-6 pt-4 sm:pt-6 border-t">
          <div className="text-center">
            <div className="text-lg sm:text-2xl font-bold text-warning">
              {orders.filter((o) => o.status === 'pending').length}
            </div>
            <div className="text-xs text-base-content/60">Ожидают</div>
          </div>
          <div className="text-center">
            <div className="text-lg sm:text-2xl font-bold text-info">
              {orders.filter((o) => o.status === 'processing').length}
            </div>
            <div className="text-xs text-base-content/60">В обработке</div>
          </div>
          <div className="text-center">
            <div className="text-lg sm:text-2xl font-bold text-success">
              {orders.filter((o) => o.status === 'ready').length}
            </div>
            <div className="text-xs text-base-content/60">Готово</div>
          </div>
        </div>

        {/* Additional Actions */}
        <div className="flex flex-col sm:flex-row gap-2 sm:gap-3 mt-4 sm:mt-6">
          <button className="btn btn-outline btn-sm">
            <CalendarIcon className="w-3 h-3 sm:w-4 sm:h-4" />
            <span className="text-xs sm:text-sm">Заказать курьера</span>
          </button>
          <button className="btn btn-outline btn-sm">
            <ArrowDownTrayIcon className="w-3 h-3 sm:w-4 sm:h-4" />
            <span className="text-xs sm:text-sm">Экспорт в Excel</span>
          </button>
        </div>

        {/* Info */}
        <div className="alert mt-3 sm:mt-4">
          <ExclamationTriangleIcon className="w-4 h-4 sm:w-5 sm:h-5 flex-shrink-0" />
          <div>
            <div className="text-xs sm:text-sm font-semibold">
              Совет по оптимизации
            </div>
            <div className="text-xs">
              Обрабатывайте заказы партиями по 10-20 штук для максимальной
              эффективности. Курьер может забрать до 50 посылок за один приезд.
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
