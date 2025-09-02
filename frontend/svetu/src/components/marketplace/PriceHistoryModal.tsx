'use client';

import React, { useEffect, useState, useCallback } from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  ChartOptions,
} from 'chart.js';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { X, AlertTriangle, TrendingDown } from 'lucide-react';
import { useTranslations } from 'next-intl';

// Регистрируем компоненты Chart.js
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

interface PricePoint {
  price: number;
  created_at: string;
}

interface PriceHistoryModalProps {
  listingId: number;
  isOpen: boolean;
  onClose: () => void;
}

export const PriceHistoryModal: React.FC<PriceHistoryModalProps> = ({
  listingId,
  isOpen,
  onClose,
}) => {
  const t = useTranslations('marketplace');
  const [priceHistory, setPriceHistory] = useState<PricePoint[]>([]);
  const [loading, setLoading] = useState(true);
  const [hasManipulation, setHasManipulation] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchPriceHistory = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch(
        `/api/v1/marketplace/listings/${listingId}/price-history`
      );

      if (!response.ok) {
        throw new Error(t('priceHistory.failedToLoad'));
      }

      const data = await response.json();

      // Сортируем данные от старых к новым для правильного отображения на графике
      const sortedHistory = (data.data || []).sort(
        (a: PricePoint, b: PricePoint) =>
          new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
      );

      setPriceHistory(sortedHistory);
      checkForManipulation(sortedHistory);
      setLoading(false);
    } catch (error) {
      console.error('Error fetching price history:', error);
      setError(t('priceHistory.error'));
      setLoading(false);
    }
  }, [listingId, t]);

  useEffect(() => {
    if (isOpen && listingId) {
      fetchPriceHistory();
    }
  }, [isOpen, listingId, fetchPriceHistory]);

  const checkForManipulation = (history: PricePoint[]) => {
    if (history.length < 3) return;

    // Проверка на манипуляции (резкий рост > 30% с последующим снижением)
    for (let i = 1; i < history.length; i++) {
      const prevPrice = history[i - 1].price;
      const currPrice = history[i].price;
      const changePercent = ((currPrice - prevPrice) / prevPrice) * 100;

      // Если цена выросла более чем на 30%
      if (changePercent > 30) {
        // Проверяем последующее снижение
        for (let j = i + 1; j < history.length; j++) {
          const futurePrice = history[j].price;
          if (futurePrice < prevPrice * 1.1) {
            setHasManipulation(true);
            return;
          }
        }
      }
    }
  };

  // Подготовка данных для графика
  const chartData = {
    labels: priceHistory.map((point) =>
      format(new Date(point.created_at), 'dd MMM yyyy', { locale: ru })
    ),
    datasets: [
      {
        label: t('priceHistory.currentPrice'),
        data: priceHistory.map((point) => point.price),
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
        tension: 0.1,
        pointRadius: 5,
        pointHoverRadius: 7,
        pointBackgroundColor: priceHistory.map((_, index) => {
          // Выделяем последнюю точку (текущую цену) красным цветом
          if (index === priceHistory.length - 1 && priceHistory.length > 1) {
            const lastPrice = priceHistory[index].price;
            const prevPrice = priceHistory[index - 1].price;
            if (lastPrice < prevPrice) {
              return 'rgb(34, 197, 94)'; // Зеленый для снижения цены
            }
          }
          return 'rgb(59, 130, 246)'; // Стандартный синий
        }),
        pointBorderColor: 'white',
        pointBorderWidth: 2,
      },
    ],
  };

  const chartOptions: ChartOptions<'line'> = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
      },
      title: {
        display: true,
        text: t('priceHistory.chartTitle'),
      },
      tooltip: {
        callbacks: {
          label: function (context) {
            return t('priceHistory.priceLabel', {
              price: context.parsed.y.toLocaleString(),
            });
          },
        },
      },
    },
    scales: {
      y: {
        beginAtZero: false,
        ticks: {
          callback: function (value) {
            return `${value} ${t('priceHistory.currency')}`;
          },
        },
      },
    },
  };

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box max-w-4xl">
        <button
          className="btn btn-sm btn-circle absolute right-2 top-2"
          onClick={onClose}
          type="button"
        >
          <X />
        </button>

        <h3 className="font-bold text-lg mb-4 flex items-center gap-2">
          <TrendingDown className="w-5 h-5" />
          {t('priceHistory.title')}
        </h3>

        {hasManipulation && (
          <div className="alert alert-warning mb-4">
            <AlertTriangle className="w-5 h-5" />
            <span>{t('priceHistory.manipulationWarning')}</span>
          </div>
        )}

        {loading ? (
          <div className="flex justify-center py-8">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        ) : error ? (
          <div className="alert alert-error">
            <span>{error}</span>
          </div>
        ) : priceHistory.length === 0 ? (
          <div className="text-center py-8 text-base-content/60">
            <p>{t('priceHistory.noData')}</p>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="h-64 w-full">
              <Line data={chartData} options={chartOptions} />
            </div>

            <div className="stats stats-horizontal w-full">
              <div className="stat">
                <div className="stat-title">
                  {t('priceHistory.currentPrice')}
                </div>
                <div className="stat-value text-lg">
                  {priceHistory[
                    priceHistory.length - 1
                  ]?.price.toLocaleString()}{' '}
                  {t('priceHistory.currency')}
                </div>
              </div>

              <div className="stat">
                <div className="stat-title">{t('priceHistory.maxPrice')}</div>
                <div className="stat-value text-lg">
                  {Math.max(
                    ...priceHistory.map((p) => p.price)
                  ).toLocaleString()}{' '}
                  {t('priceHistory.currency')}
                </div>
              </div>

              <div className="stat">
                <div className="stat-title">{t('priceHistory.minPrice')}</div>
                <div className="stat-value text-lg">
                  {Math.min(
                    ...priceHistory.map((p) => p.price)
                  ).toLocaleString()}{' '}
                  {t('priceHistory.currency')}
                </div>
              </div>
            </div>

            {/* Таблица с историей изменений */}
            <div className="divider">История изменений</div>
            <div className="overflow-x-auto max-h-48">
              <table className="table table-zebra table-sm">
                <thead>
                  <tr>
                    <th>Дата</th>
                    <th>Цена</th>
                    <th>Изменение</th>
                  </tr>
                </thead>
                <tbody>
                  {priceHistory.map((point, index) => {
                    const prevPrice =
                      index > 0 ? priceHistory[index - 1].price : null;
                    const changePercent = prevPrice
                      ? (((point.price - prevPrice) / prevPrice) * 100).toFixed(
                          1
                        )
                      : null;

                    return (
                      <tr key={index}>
                        <td>
                          {format(
                            new Date(point.created_at),
                            'dd.MM.yyyy HH:mm',
                            { locale: ru }
                          )}
                        </td>
                        <td className="font-semibold">
                          {point.price.toLocaleString()} РСД
                        </td>
                        <td>
                          {changePercent && (
                            <span
                              className={`badge badge-sm ${
                                parseFloat(changePercent) < 0
                                  ? 'badge-success'
                                  : parseFloat(changePercent) > 0
                                    ? 'badge-error'
                                    : 'badge-ghost'
                              }`}
                            >
                              {parseFloat(changePercent) > 0 && '+'}
                              {changePercent}%
                            </span>
                          )}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>
        )}

        <div className="modal-action">
          <button className="btn" onClick={onClose}>
            {t('priceHistory.close')}
          </button>
        </div>
      </div>
    </div>
  );
};
