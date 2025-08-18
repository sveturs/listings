'use client';

import { useState, useEffect, useRef } from 'react';
import {
  QrCodeIcon,
  ClipboardIcon,
  CheckIcon,
  MapPinIcon,
  ClockIcon,
  PhoneIcon,
  InformationCircleIcon,
  PrinterIcon,
  ShareIcon,
  CalendarDaysIcon,
} from '@heroicons/react/24/outline';
// import { useTranslations } from 'next-intl';

interface PickupOrder {
  id: number;
  pickup_code: string;
  qr_code?: string;
  status: string;
  created_at: string;
  expires_at: string;
  confirmed_at?: string;
  customer_name: string;
  customer_phone: string;
  items_count: number;
  total_amount: number;
  warehouse: {
    code: string;
    name: string;
    address: string;
    phone: string;
    working_hours: Record<string, string>;
  };
  notes?: string;
}

interface Props {
  pickupOrder: PickupOrder;
  onStatusUpdate?: (newStatus: string) => void;
  className?: string;
}

export default function PostExpressPickupCode({
  pickupOrder,
  onStatusUpdate: _onStatusUpdate,
  className = '',
}: Props) {
  // const t = useTranslations('delivery');
  const [copied, setCopied] = useState(false);
  const [qrCodeData, setQrCodeData] = useState<string | null>(null);
  const printRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Генерируем QR код данные (можно заменить на реальный QR код с сервера)
    if (pickupOrder.qr_code) {
      setQrCodeData(pickupOrder.qr_code);
    } else {
      // Генерируем простой QR код с данными заказа
      const qrData = JSON.stringify({
        code: pickupOrder.pickup_code,
        order_id: pickupOrder.id,
        warehouse: pickupOrder.warehouse.code,
      });
      setQrCodeData(qrData);
    }
  }, [pickupOrder]);

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(pickupOrder.pickup_code);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy to clipboard:', err);
    }
  };

  const handlePrint = () => {
    const printContent = printRef.current;
    if (!printContent) return;

    const printWindow = window.open('', '_blank');
    if (!printWindow) return;

    printWindow.document.write(`
      <html>
        <head>
          <title>Код самовывоза ${pickupOrder.pickup_code}</title>
          <style>
            body { font-family: Arial, sans-serif; margin: 20px; }
            .pickup-code { font-size: 24px; font-weight: bold; text-align: center; margin: 20px 0; }
            .qr-code { text-align: center; margin: 20px 0; }
            .details { margin: 10px 0; }
            @media print { .no-print { display: none; } }
          </style>
        </head>
        <body>
          ${printContent.innerHTML}
        </body>
      </html>
    `);

    printWindow.document.close();
    printWindow.print();
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'text-warning';
      case 'ready':
        return 'text-primary';
      case 'confirmed':
        return 'text-success';
      case 'expired':
        return 'text-error';
      default:
        return 'text-base-content';
    }
  };

  const getStatusText = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'Подготовка заказа';
      case 'ready':
        return 'Готов к выдаче';
      case 'confirmed':
        return 'Выдан';
      case 'expired':
        return 'Истек срок';
      default:
        return status;
    }
  };

  const isExpired = () => {
    return new Date() > new Date(pickupOrder.expires_at);
  };

  const getDaysLeft = () => {
    const now = new Date();
    const expires = new Date(pickupOrder.expires_at);
    const diffTime = expires.getTime() - now.getTime();
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return Math.max(0, diffDays);
  };

  const formatDateTime = (dateTime: string) => {
    return new Date(dateTime).toLocaleString('sr-RS', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getTodayWorkingHours = () => {
    const dayNames = [
      'sunday',
      'monday',
      'tuesday',
      'wednesday',
      'thursday',
      'friday',
      'saturday',
    ];
    const todayName = dayNames[new Date().getDay()];
    return pickupOrder.warehouse.working_hours[todayName] || 'Часы не указаны';
  };

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Статус заказа */}
      <div
        className={`alert ${pickupOrder.status === 'ready' ? 'alert-success' : 'alert-info'}`}
      >
        <InformationCircleIcon className="w-5 h-5" />
        <div>
          <h4 className="font-semibold">
            Статус заказа:{' '}
            <span className={getStatusColor(pickupOrder.status)}>
              {getStatusText(pickupOrder.status)}
            </span>
          </h4>
          <p className="text-sm mt-1">
            {pickupOrder.status === 'ready' &&
              'Ваш заказ готов к получению на складе!'}
            {pickupOrder.status === 'pending' &&
              'Заказ готовится к выдаче. Вы получите уведомление когда будет готов.'}
            {pickupOrder.status === 'confirmed' &&
              `Заказ был выдан ${formatDateTime(pickupOrder.confirmed_at!)}`}
            {pickupOrder.status === 'expired' &&
              'Срок действия кода истек. Обратитесь в службу поддержки.'}
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Код самовывоза */}
        <div className="card bg-gradient-to-r from-primary/5 to-secondary/5 shadow-lg">
          <div className="card-body p-6 text-center" ref={printRef}>
            <h3 className="text-xl font-bold mb-4">Код самовывоза</h3>

            {/* QR код */}
            {qrCodeData && (
              <div className="qr-code mb-6">
                <div className="w-32 h-32 mx-auto bg-white border-2 border-base-300 rounded-lg flex items-center justify-center">
                  <QrCodeIcon className="w-24 h-24 text-base-content/20" />
                  {/* Здесь должен быть реальный QR код */}
                </div>
                <div className="text-xs text-base-content/60 mt-2">
                  Покажите QR код на складе
                </div>
              </div>
            )}

            {/* Код */}
            <div className="pickup-code mb-6">
              <div className="text-3xl font-mono font-bold text-primary tracking-wider mb-2">
                {pickupOrder.pickup_code}
              </div>
              <div className="text-sm text-base-content/60">
                Назовите этот код сотруднику склада
              </div>
            </div>

            {/* Срок действия */}
            <div
              className={`p-3 rounded-lg mb-4 ${
                isExpired()
                  ? 'bg-error/10 text-error'
                  : getDaysLeft() <= 1
                    ? 'bg-warning/10 text-warning'
                    : 'bg-success/10 text-success'
              }`}
            >
              <div className="flex items-center justify-center gap-2">
                <CalendarDaysIcon className="w-5 h-5" />
                <span className="font-medium">
                  {isExpired()
                    ? 'Код истек'
                    : getDaysLeft() === 0
                      ? 'Последний день'
                      : `Осталось ${getDaysLeft()} дней`}
                </span>
              </div>
              <div className="text-sm mt-1">
                До {formatDateTime(pickupOrder.expires_at)}
              </div>
            </div>

            {/* Действия */}
            <div className="flex flex-wrap gap-2 justify-center no-print">
              <button
                className={`btn btn-sm ${copied ? 'btn-success' : 'btn-outline'}`}
                onClick={copyToClipboard}
              >
                {copied ? (
                  <CheckIcon className="w-4 h-4" />
                ) : (
                  <ClipboardIcon className="w-4 h-4" />
                )}
                {copied ? 'Скопировано' : 'Копировать'}
              </button>

              <button className="btn btn-sm btn-outline" onClick={handlePrint}>
                <PrinterIcon className="w-4 h-4" />
                Печать
              </button>

              <button className="btn btn-sm btn-outline">
                <ShareIcon className="w-4 h-4" />
                Поделиться
              </button>
            </div>
          </div>
        </div>

        {/* Информация о складе */}
        <div className="space-y-4">
          {/* Склад */}
          <div className="card bg-base-100 shadow-lg">
            <div className="card-body p-6">
              <h4 className="font-semibold text-lg mb-4 flex items-center gap-2">
                <MapPinIcon className="w-5 h-5 text-primary" />
                {pickupOrder.warehouse.name}
              </h4>

              <div className="space-y-3">
                <div className="flex items-start gap-2">
                  <MapPinIcon className="w-5 h-5 text-base-content/60 mt-0.5 flex-shrink-0" />
                  <div>
                    <div className="font-medium">
                      {pickupOrder.warehouse.address}
                    </div>
                    <div className="text-sm text-base-content/60">
                      Код склада: {pickupOrder.warehouse.code}
                    </div>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  <PhoneIcon className="w-5 h-5 text-base-content/60" />
                  <div>{pickupOrder.warehouse.phone}</div>
                </div>

                <div className="flex items-start gap-2">
                  <ClockIcon className="w-5 h-5 text-base-content/60 mt-0.5" />
                  <div>
                    <div className="font-medium">
                      Сегодня: {getTodayWorkingHours()}
                    </div>
                    <div className="text-sm text-base-content/60">
                      Время работы склада
                    </div>
                  </div>
                </div>
              </div>

              {/* Полное расписание */}
              <div className="mt-4 pt-4 border-t">
                <div className="text-sm font-medium mb-2">
                  Полное расписание:
                </div>
                <div className="grid grid-cols-2 gap-1 text-sm">
                  {Object.entries(pickupOrder.warehouse.working_hours).map(
                    ([day, hours]) => (
                      <div key={day} className="flex justify-between">
                        <span className="text-base-content/70">
                          {day === 'monday' && 'Пн'}
                          {day === 'tuesday' && 'Вт'}
                          {day === 'wednesday' && 'Ср'}
                          {day === 'thursday' && 'Чт'}
                          {day === 'friday' && 'Пт'}
                          {day === 'saturday' && 'Сб'}
                          {day === 'sunday' && 'Вс'}
                        </span>
                        <span className="font-medium">{hours}</span>
                      </div>
                    )
                  )}
                </div>
              </div>
            </div>
          </div>

          {/* Детали заказа */}
          <div className="card bg-base-100 shadow-lg">
            <div className="card-body p-6">
              <h4 className="font-semibold text-lg mb-4">Детали заказа</h4>

              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-base-content/70">Получатель:</span>
                  <span className="font-medium">
                    {pickupOrder.customer_name}
                  </span>
                </div>

                <div className="flex justify-between">
                  <span className="text-base-content/70">Телефон:</span>
                  <span className="font-medium">
                    {pickupOrder.customer_phone}
                  </span>
                </div>

                <div className="flex justify-between">
                  <span className="text-base-content/70">
                    Количество товаров:
                  </span>
                  <span className="font-medium">
                    {pickupOrder.items_count} шт.
                  </span>
                </div>

                <div className="flex justify-between">
                  <span className="text-base-content/70">Сумма заказа:</span>
                  <span className="font-medium text-primary">
                    {pickupOrder.total_amount} RSD
                  </span>
                </div>

                <div className="flex justify-between">
                  <span className="text-base-content/70">Дата создания:</span>
                  <span className="font-medium">
                    {formatDateTime(pickupOrder.created_at)}
                  </span>
                </div>

                {pickupOrder.notes && (
                  <div className="pt-3 border-t">
                    <div className="text-sm text-base-content/70 mb-1">
                      Примечания:
                    </div>
                    <div className="text-sm">{pickupOrder.notes}</div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Инструкции */}
      <div className="card bg-gradient-to-r from-info/5 to-info/10">
        <div className="card-body p-6">
          <h4 className="font-semibold text-lg mb-4 flex items-center gap-2">
            <InformationCircleIcon className="w-5 h-5 text-info" />
            Как получить заказ
          </h4>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-3">
              <div className="flex items-start gap-3">
                <div className="w-6 h-6 bg-primary text-primary-content rounded-full flex items-center justify-center text-sm font-bold flex-shrink-0">
                  1
                </div>
                <div>
                  <div className="font-medium">Приезжайте на склад</div>
                  <div className="text-sm text-base-content/70">
                    В рабочее время с документом удостоверяющим личность
                  </div>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <div className="w-6 h-6 bg-primary text-primary-content rounded-full flex items-center justify-center text-sm font-bold flex-shrink-0">
                  2
                </div>
                <div>
                  <div className="font-medium">
                    Назовите код или покажите QR
                  </div>
                  <div className="text-sm text-base-content/70">
                    Сотрудник склада найдет ваш заказ
                  </div>
                </div>
              </div>
            </div>

            <div className="space-y-3">
              <div className="flex items-start gap-3">
                <div className="w-6 h-6 bg-primary text-primary-content rounded-full flex items-center justify-center text-sm font-bold flex-shrink-0">
                  3
                </div>
                <div>
                  <div className="font-medium">Проверьте товары</div>
                  <div className="text-sm text-base-content/70">
                    Возможность примерки и проверки товаров
                  </div>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <div className="w-6 h-6 bg-primary text-primary-content rounded-full flex items-center justify-center text-sm font-bold flex-shrink-0">
                  4
                </div>
                <div>
                  <div className="font-medium">Получите заказ</div>
                  <div className="text-sm text-base-content/70">
                    Заберите свои товары и получите чек
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div className="mt-6 p-4 bg-warning/10 rounded-lg">
            <div className="text-sm">
              <strong>Важно:</strong> Код действует {getDaysLeft()} дней с
              момента готовности заказа. После истечения срока заказ может быть
              возвращен продавцу.
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
