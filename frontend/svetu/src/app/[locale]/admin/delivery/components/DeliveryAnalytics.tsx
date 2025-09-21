'use client';

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';

interface AnalyticsData {
  totalShipments: number;
  deliveryRate: number;
  avgDeliveryTime: string;
  customerSatisfaction: number;
  costPerShipment: number;
  problemRate: number;
}

export default function DeliveryAnalytics() {
  const t = useTranslations('admin.delivery.analytics');
  const [analyticsData, setAnalyticsData] = useState<AnalyticsData>({
    totalShipments: 0,
    deliveryRate: 0,
    avgDeliveryTime: '0 дней',
    customerSatisfaction: 0,
    costPerShipment: 0,
    problemRate: 0,
  });
  const [timeRange, setTimeRange] = useState('30d');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchAnalytics();
  }, [timeRange]);

  const fetchAnalytics = async () => {
    try {
      // Mock data
      setTimeout(() => {
        setAnalyticsData({
          totalShipments: 1247,
          deliveryRate: 94.3,
          avgDeliveryTime: '2.1 дня',
          customerSatisfaction: 4.5,
          costPerShipment: 12.5,
          problemRate: 2.3,
        });
        setLoading(false);
      }, 500);
    } catch (error) {
      console.error('Failed to fetch analytics:', error);
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header with filters */}
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold">{t('title')}</h2>
        <div className="flex gap-4">
          <select
            className="select select-bordered select-sm"
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value)}
          >
            <option value="7d">Последние 7 дней</option>
            <option value="30d">Последние 30 дней</option>
            <option value="90d">Последние 90 дней</option>
            <option value="365d">Последний год</option>
          </select>
          <button className="btn btn-primary btn-sm">{t('export')}</button>
        </div>
      </div>

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-4">
        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-title text-sm">
            {t('metrics.totalShipments')}
          </div>
          <div className="stat-value text-2xl">
            {analyticsData.totalShipments}
          </div>
          <div className="stat-desc">↗︎ +15% от предыдущего периода</div>
        </div>

        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-title text-sm">{t('metrics.deliveryRate')}</div>
          <div className="stat-value text-2xl">
            {analyticsData.deliveryRate}%
          </div>
          <div className="stat-desc">↗︎ +2.3% от предыдущего</div>
        </div>

        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-title text-sm">
            {t('metrics.avgDeliveryTime')}
          </div>
          <div className="stat-value text-2xl">
            {analyticsData.avgDeliveryTime}
          </div>
          <div className="stat-desc">↘︎ -4ч от предыдущего</div>
        </div>

        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-title text-sm">
            {t('metrics.customerSatisfaction')}
          </div>
          <div className="stat-value text-2xl">
            {analyticsData.customerSatisfaction}
            <span className="text-lg">/5</span>
          </div>
          <div className="stat-desc">↗︎ +0.2 от предыдущего</div>
        </div>

        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-title text-sm">
            {t('metrics.costPerShipment')}
          </div>
          <div className="stat-value text-2xl">
            €{analyticsData.costPerShipment}
          </div>
          <div className="stat-desc">↘︎ -€0.50 от предыдущего</div>
        </div>

        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-title text-sm">{t('metrics.problemRate')}</div>
          <div className="stat-value text-2xl">
            {analyticsData.problemRate}%
          </div>
          <div className="stat-desc">↘︎ -0.5% от предыдущего</div>
        </div>
      </div>

      {/* Charts Section */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Delivery Trends Chart */}
        <div className="card bg-base-100 shadow">
          <div className="card-body">
            <h3 className="card-title text-base">
              {t('charts.deliveryTrends')}
            </h3>
            <div className="h-64 flex items-center justify-center bg-base-200 rounded">
              <span className="text-base-content/50">
                График трендов доставок
              </span>
            </div>
          </div>
        </div>

        {/* Provider Comparison */}
        <div className="card bg-base-100 shadow">
          <div className="card-body">
            <h3 className="card-title text-base">
              {t('charts.providerComparison')}
            </h3>
            <div className="space-y-3">
              {[
                { name: 'Post Express', value: 94, color: 'primary' },
                { name: 'BEX Express', value: 95, color: 'secondary' },
                { name: 'AKS', value: 92, color: 'accent' },
                { name: 'D Express', value: 97, color: 'success' },
                { name: 'City Express', value: 91, color: 'warning' },
              ].map((provider) => (
                <div key={provider.name}>
                  <div className="flex justify-between mb-1">
                    <span className="text-sm">{provider.name}</span>
                    <span className="text-sm font-semibold">
                      {provider.value}%
                    </span>
                  </div>
                  <progress
                    className={`progress progress-${provider.color} h-2`}
                    value={provider.value}
                    max="100"
                  />
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Cost Analysis */}
        <div className="card bg-base-100 shadow">
          <div className="card-body">
            <h3 className="card-title text-base">{t('charts.costAnalysis')}</h3>
            <div className="h-64 flex items-center justify-center bg-base-200 rounded">
              <span className="text-base-content/50">
                График анализа затрат
              </span>
            </div>
          </div>
        </div>

        {/* Geographic Distribution */}
        <div className="card bg-base-100 shadow">
          <div className="card-body">
            <h3 className="card-title text-base">
              {t('charts.geographicDistribution')}
            </h3>
            <div className="space-y-2">
              <div className="flex justify-between">
                <span>Белград</span>
                <span className="font-semibold">42%</span>
              </div>
              <div className="flex justify-between">
                <span>Нови-Сад</span>
                <span className="font-semibold">18%</span>
              </div>
              <div className="flex justify-between">
                <span>Ниш</span>
                <span className="font-semibold">12%</span>
              </div>
              <div className="flex justify-between">
                <span>Крагуевац</span>
                <span className="font-semibold">8%</span>
              </div>
              <div className="flex justify-between">
                <span>Суботица</span>
                <span className="font-semibold">6%</span>
              </div>
              <div className="flex justify-between">
                <span>Другие города</span>
                <span className="font-semibold">14%</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Detailed Report Section */}
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <h3 className="card-title text-base mb-4">Детальный отчет</h3>
          <div className="overflow-x-auto">
            <table className="table table-sm">
              <thead>
                <tr>
                  <th>Метрика</th>
                  <th>Текущий период</th>
                  <th>Предыдущий период</th>
                  <th>Изменение</th>
                  <th>Тренд</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Общий объем доставок</td>
                  <td>1,247</td>
                  <td>1,084</td>
                  <td>+163</td>
                  <td className="text-success">↗︎ +15.0%</td>
                </tr>
                <tr>
                  <td>Успешные доставки</td>
                  <td>1,176</td>
                  <td>1,005</td>
                  <td>+171</td>
                  <td className="text-success">↗︎ +17.0%</td>
                </tr>
                <tr>
                  <td>Возвраты</td>
                  <td>42</td>
                  <td>48</td>
                  <td>-6</td>
                  <td className="text-success">↘︎ -12.5%</td>
                </tr>
                <tr>
                  <td>Проблемные доставки</td>
                  <td>29</td>
                  <td>31</td>
                  <td>-2</td>
                  <td className="text-success">↘︎ -6.5%</td>
                </tr>
                <tr>
                  <td>Общая стоимость</td>
                  <td>€15,587.50</td>
                  <td>€14,092.00</td>
                  <td>+€1,495.50</td>
                  <td className="text-warning">↗︎ +10.6%</td>
                </tr>
                <tr>
                  <td>Средняя маржа</td>
                  <td>18.5%</td>
                  <td>17.2%</td>
                  <td>+1.3%</td>
                  <td className="text-success">↗︎ +7.6%</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}
