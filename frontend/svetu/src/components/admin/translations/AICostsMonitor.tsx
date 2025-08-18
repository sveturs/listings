'use client';

import { useState, useEffect, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import {
  CurrencyDollarIcon,
  ExclamationTriangleIcon,
  ArrowTrendingUpIcon,
  ArrowPathIcon,
  CalendarIcon,
  ChartBarIcon,
  InformationCircleIcon,
} from '@heroicons/react/24/outline';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts';

interface ProviderCosts {
  provider: string;
  total_cost: number;
  total_tokens: number;
  total_requests: number;
  last_updated: string;
  daily_costs: Record<string, number>;
  hourly_costs: Record<string, number>;
}

interface CostsSummary {
  total_cost: number;
  total_tokens: number;
  total_requests: number;
  today_cost: number;
  month_cost: number;
  by_provider: Record<string, ProviderCosts>;
  today_by_provider: Record<string, number>;
  month_by_provider: Record<string, number>;
}

interface CostAlerts {
  alerts: string[];
  daily_limit: number;
  monthly_limit: number;
  has_alerts: boolean;
}

export default function AICostsMonitor() {
  const _t = useTranslations('admin');
  const [costsSummary, setCostsSummary] = useState<CostsSummary | null>(null);
  const [alerts, setAlerts] = useState<CostAlerts | null>(null);
  const [_loading, _setLoading] = useState(false); // Changed to false to show demo data immediately
  const [refreshing, setRefreshing] = useState(false);
  const [dailyLimit, setDailyLimit] = useState(100);
  const [monthlyLimit, setMonthlyLimit] = useState(2000);
  const [isDemo, setIsDemo] = useState(true); // Flag to show we're in demo mode

  // Initialize with demo data
  const initializeDemoData = () => {
    setCostsSummary({
      total_cost: 0,
      total_tokens: 0,
      total_requests: 0,
      today_cost: 0,
      month_cost: 0,
      by_provider: {
        openai: {
          provider: 'openai',
          total_cost: 0,
          total_tokens: 0,
          total_requests: 0,
          last_updated: new Date().toISOString(),
          daily_costs: {},
          hourly_costs: {},
        },
        google: {
          provider: 'google',
          total_cost: 0,
          total_tokens: 0,
          total_requests: 0,
          last_updated: new Date().toISOString(),
          daily_costs: {},
          hourly_costs: {},
        },
        deepl: {
          provider: 'deepl',
          total_cost: 0,
          total_tokens: 0,
          total_requests: 0,
          last_updated: new Date().toISOString(),
          daily_costs: {},
          hourly_costs: {},
        },
        claude: {
          provider: 'claude',
          total_cost: 0,
          total_tokens: 0,
          total_requests: 0,
          last_updated: new Date().toISOString(),
          daily_costs: {},
          hourly_costs: {},
        },
      },
      today_by_provider: {},
      month_by_provider: {},
    });

    setAlerts({
      alerts: [] as string[],
      daily_limit: dailyLimit,
      monthly_limit: monthlyLimit,
      has_alerts: false,
    });
  };

  const fetchCostsData = async () => {
    setRefreshing(true);

    try {
      // Use local API route that handles authentication
      const baseUrl =
        typeof window !== 'undefined'
          ? window.location.origin
          : 'http://localhost:3001';
      const response = await fetch(`${baseUrl}/api/admin/translations/costs`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });

      console.log('Costs API response:', response.status);

      if (!response.ok) {
        console.error('Failed to fetch costs:', response.status);
        initializeDemoData();
        setIsDemo(true);
        setRefreshing(false);
        return;
      }

      const data = await response.json();
      console.log('Costs data received:', data);

      if (data.success && data.data) {
        console.log('Setting real costs data:', data.data);
        setCostsSummary(data.data);
        setIsDemo(false);

        // Check alerts
        const alertsData = {
          alerts: [] as string[],
          daily_limit: dailyLimit,
          monthly_limit: monthlyLimit,
          has_alerts: false,
        };

        if (data.data.today_cost > dailyLimit * 0.8) {
          alertsData.alerts.push(
            `Дневной расход приближается к лимиту: ${formatCurrency(data.data.today_cost)} из ${formatCurrency(dailyLimit)}`
          );
          alertsData.has_alerts = true;
        }

        if (data.data.month_cost > monthlyLimit * 0.8) {
          alertsData.alerts.push(
            `Месячный расход приближается к лимиту: ${formatCurrency(data.data.month_cost)} из ${formatCurrency(monthlyLimit)}`
          );
          alertsData.has_alerts = true;
        }

        setAlerts(alertsData);
      } else {
        initializeDemoData();
        setIsDemo(true);
      }
    } catch (error) {
      console.error('Error fetching costs:', error);
      initializeDemoData();
      setIsDemo(true);
    } finally {
      setRefreshing(false);
    }
  };

  useEffect(() => {
    fetchCostsData();

    // Auto-refresh every 30 seconds
    const interval = setInterval(() => {
      fetchCostsData();
    }, 30000);

    return () => clearInterval(interval);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [dailyLimit, monthlyLimit]);

  const resetProviderCosts = async (provider: string) => {
    if (!confirm(`Вы уверены, что хотите сбросить счетчики для ${provider}?`)) {
      return;
    }

    // For now, just refresh the data since reset endpoint might not exist
    // TODO: Implement reset endpoint when needed
    console.log('Reset requested for provider:', provider);
    alert('Функция сброса будет добавлена в следующей версии');
    fetchCostsData();
  };

  const formatCurrency = (amount: number) => {
    // Для очень маленьких сумм показываем больше знаков после запятой
    const fractionDigits = amount < 0.01 && amount > 0 ? 6 : 2;
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: fractionDigits,
      maximumFractionDigits: fractionDigits,
    }).format(amount);
  };

  const formatNumber = (num: number) => {
    return new Intl.NumberFormat('ru-RU').format(num);
  };

  // Подготовка данных для графика расходов по дням
  const chartData = useMemo(() => {
    if (!costsSummary) return [];

    // Собираем все уникальные даты из всех провайдеров
    const allDates = new Set<string>();
    Object.values(costsSummary.by_provider).forEach((provider) => {
      Object.keys(provider.daily_costs).forEach((date) => allDates.add(date));
    });

    // Сортируем даты и берем последние 7 дней
    const sortedDates = Array.from(allDates).sort();
    const last7Days = sortedDates.slice(-7);

    // Формируем данные для графика
    return last7Days.map((date) => {
      const dataPoint: any = {
        date: new Date(date).toLocaleDateString('ru-RU', {
          day: 'numeric',
          month: 'short',
        }),
        total: 0,
      };

      // Добавляем данные по каждому провайдеру
      Object.entries(costsSummary.by_provider).forEach(([provider, data]) => {
        const cost = data.daily_costs[date] || 0;
        dataPoint[provider] = cost;
        dataPoint.total += cost;
      });

      return dataPoint;
    });
  }, [costsSummary]);

  // Подготовка данных для круговой диаграммы
  const pieData = useMemo(() => {
    if (!costsSummary) return [];

    return Object.entries(costsSummary.by_provider)
      .filter(([_, data]) => data.total_cost > 0)
      .map(([provider, data]) => ({
        name: provider.charAt(0).toUpperCase() + provider.slice(1),
        value: data.total_cost,
      }));
  }, [costsSummary]);

  // Цвета для провайдеров
  const COLORS = {
    openai: '#10a37f',
    google: '#4285f4',
    deepl: '#0f2b46',
    claude: '#d97757',
  };

  const PIE_COLORS = ['#10a37f', '#4285f4', '#0f2b46', '#d97757'];

  if (_loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Demo Mode Notice */}
      {isDemo && (
        <div className="alert alert-info">
          <InformationCircleIcon className="h-6 w-6" />
          <div>
            <h3 className="font-bold">Демо режим</h3>
            <p className="text-sm">
              Система мониторинга расходов готова к работе. Данные будут
              отображаться после начала использования AI переводов.
            </p>
          </div>
        </div>
      )}

      {/* Alerts Section */}
      {alerts?.has_alerts && (
        <div className="alert alert-warning">
          <ExclamationTriangleIcon className="h-6 w-6" />
          <div>
            <h3 className="font-bold">Предупреждения о расходах</h3>
            <ul className="mt-2 space-y-1">
              {alerts.alerts.map((alert, index) => (
                <li key={index} className="text-sm">
                  {alert}
                </li>
              ))}
            </ul>
          </div>
        </div>
      )}

      {/* Summary Stats */}
      <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-primary">
            <CurrencyDollarIcon className="h-8 w-8" />
          </div>
          <div className="stat-title">Всего расходов</div>
          <div className="stat-value text-2xl">
            {formatCurrency(costsSummary?.total_cost || 0)}
          </div>
          <div className="stat-desc">За все время</div>
        </div>

        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-info">
            <CalendarIcon className="h-8 w-8" />
          </div>
          <div className="stat-title">Сегодня</div>
          <div className="stat-value text-2xl">
            {formatCurrency(costsSummary?.today_cost || 0)}
          </div>
          <div className="stat-desc">Лимит: {formatCurrency(dailyLimit)}</div>
        </div>

        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-success">
            <ArrowTrendingUpIcon className="h-8 w-8" />
          </div>
          <div className="stat-title">Этот месяц</div>
          <div className="stat-value text-2xl">
            {formatCurrency(costsSummary?.month_cost || 0)}
          </div>
          <div className="stat-desc">Лимит: {formatCurrency(monthlyLimit)}</div>
        </div>

        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-secondary">
            <ChartBarIcon className="h-8 w-8" />
          </div>
          <div className="stat-title">Токены</div>
          <div className="stat-value text-2xl">
            {formatNumber(costsSummary?.total_tokens || 0)}
          </div>
          <div className="stat-desc">Использовано токенов</div>
        </div>

        <div className="stat bg-base-100 rounded-lg border border-base-300">
          <div className="stat-figure text-accent">
            <ArrowPathIcon className="h-8 w-8" />
          </div>
          <div className="stat-title">Запросы</div>
          <div className="stat-value text-2xl">
            {formatNumber(costsSummary?.total_requests || 0)}
          </div>
          <div className="stat-desc">Всего запросов</div>
        </div>
      </div>

      {/* Provider Details */}
      <div className="card bg-base-100">
        <div className="card-body">
          <div className="flex justify-between items-center mb-4">
            <h2 className="card-title">Расходы по провайдерам</h2>
            <button
              className="btn btn-sm btn-ghost"
              onClick={() => fetchCostsData()}
              disabled={refreshing}
            >
              <ArrowPathIcon
                className={`h-4 w-4 ${refreshing ? 'animate-spin' : ''}`}
              />
              Обновить
            </button>
          </div>

          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th>Провайдер</th>
                  <th>Всего расходов</th>
                  <th>Сегодня</th>
                  <th>Месяц</th>
                  <th>Токены</th>
                  <th>Запросы</th>
                  <th>Действия</th>
                </tr>
              </thead>
              <tbody>
                {costsSummary?.by_provider &&
                Object.keys(costsSummary.by_provider).length > 0
                  ? Object.entries(costsSummary.by_provider).map(
                      ([provider, data]) => (
                        <tr key={provider}>
                          <td className="font-medium capitalize">{provider}</td>
                          <td>{formatCurrency(data.total_cost)}</td>
                          <td>
                            {formatCurrency(
                              costsSummary.today_by_provider?.[provider] || 0
                            )}
                          </td>
                          <td>
                            {formatCurrency(
                              costsSummary.month_by_provider?.[provider] || 0
                            )}
                          </td>
                          <td>{formatNumber(data.total_tokens)}</td>
                          <td>{formatNumber(data.total_requests)}</td>
                          <td>
                            <button
                              className="btn btn-xs btn-ghost"
                              onClick={() => resetProviderCosts(provider)}
                            >
                              Сбросить
                            </button>
                          </td>
                        </tr>
                      )
                    )
                  : // Show default providers when no data
                    ['openai', 'google', 'deepl', 'claude'].map((provider) => (
                      <tr key={provider}>
                        <td className="font-medium capitalize">{provider}</td>
                        <td>{formatCurrency(0)}</td>
                        <td>{formatCurrency(0)}</td>
                        <td>{formatCurrency(0)}</td>
                        <td>{formatNumber(0)}</td>
                        <td>{formatNumber(0)}</td>
                        <td>
                          <button className="btn btn-xs btn-ghost" disabled>
                            Сбросить
                          </button>
                        </td>
                      </tr>
                    ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* Settings */}
      <div className="card bg-base-100">
        <div className="card-body">
          <h2 className="card-title mb-4">Настройки лимитов</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">Дневной лимит (USD)</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={dailyLimit}
                onChange={(e) => setDailyLimit(Number(e.target.value))}
                min={1}
                step={10}
              />
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">Месячный лимит (USD)</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={monthlyLimit}
                onChange={(e) => setMonthlyLimit(Number(e.target.value))}
                min={1}
                step={100}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Графики расходов */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* График расходов по дням */}
        <div className="card bg-base-100">
          <div className="card-body">
            <h2 className="card-title mb-4">Расходы по дням</h2>
            {chartData.length > 0 ? (
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" />
                  <YAxis />
                  <Tooltip
                    formatter={(value: number) => formatCurrency(value)}
                    labelStyle={{ color: '#000' }}
                  />
                  <Legend />
                  <Line
                    type="monotone"
                    dataKey="total"
                    stroke="#8b5cf6"
                    strokeWidth={2}
                    name="Всего"
                  />
                  {Object.keys(COLORS).map((provider) => (
                    <Line
                      key={provider}
                      type="monotone"
                      dataKey={provider}
                      stroke={COLORS[provider as keyof typeof COLORS]}
                      strokeWidth={1}
                      name={
                        provider.charAt(0).toUpperCase() + provider.slice(1)
                      }
                    />
                  ))}
                </LineChart>
              </ResponsiveContainer>
            ) : (
              <div className="h-64 flex items-center justify-center border-2 border-dashed border-base-300 rounded-lg">
                <div className="text-center text-base-content/60">
                  <ChartBarIcon className="h-12 w-12 mx-auto mb-2" />
                  <p>Нет данных для отображения</p>
                  <p className="text-sm mt-1">
                    Данные появятся после использования AI переводов
                  </p>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Круговая диаграмма распределения по провайдерам */}
        <div className="card bg-base-100">
          <div className="card-body">
            <h2 className="card-title mb-4">Распределение по провайдерам</h2>
            {pieData.length > 0 ? (
              <ResponsiveContainer width="100%" height={300}>
                <PieChart>
                  <Pie
                    data={pieData}
                    cx="50%"
                    cy="50%"
                    labelLine={false}
                    label={({ name, percent }) =>
                      `${name} ${((percent || 0) * 100).toFixed(0)}%`
                    }
                    outerRadius={80}
                    fill="#8884d8"
                    dataKey="value"
                  >
                    {pieData.map((entry, index) => (
                      <Cell
                        key={`cell-${index}`}
                        fill={PIE_COLORS[index % PIE_COLORS.length]}
                      />
                    ))}
                  </Pie>
                  <Tooltip
                    formatter={(value: number) => formatCurrency(value)}
                  />
                </PieChart>
              </ResponsiveContainer>
            ) : (
              <div className="h-64 flex items-center justify-center border-2 border-dashed border-base-300 rounded-lg">
                <div className="text-center text-base-content/60">
                  <ChartBarIcon className="h-12 w-12 mx-auto mb-2" />
                  <p>Нет данных для отображения</p>
                  <p className="text-sm mt-1">
                    Данные появятся после использования AI переводов
                  </p>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Почасовая статистика за сегодня */}
      <div className="card bg-base-100">
        <div className="card-body">
          <h2 className="card-title mb-4">Почасовая статистика за сегодня</h2>
          {(() => {
            // Подготовка данных для почасового графика
            const hourlyData: any[] = [];
            const today = new Date().toISOString().split('T')[0];

            if (costsSummary) {
              for (let hour = 0; hour < 24; hour++) {
                const hourKey = `${today}T${hour.toString().padStart(2, '0')}`;
                let totalCost = 0;

                Object.values(costsSummary.by_provider).forEach((provider) => {
                  if (provider.hourly_costs[hourKey]) {
                    totalCost += provider.hourly_costs[hourKey];
                  }
                });

                if (totalCost > 0) {
                  hourlyData.push({
                    hour: `${hour}:00`,
                    cost: totalCost,
                  });
                }
              }
            }

            return hourlyData.length > 0 ? (
              <ResponsiveContainer width="100%" height={200}>
                <BarChart data={hourlyData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="hour" />
                  <YAxis />
                  <Tooltip
                    formatter={(value: number) => formatCurrency(value)}
                  />
                  <Bar dataKey="cost" fill="#8b5cf6" />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="h-32 flex items-center justify-center border-2 border-dashed border-base-300 rounded-lg">
                <div className="text-center text-base-content/60">
                  <p>Нет активности за сегодня</p>
                </div>
              </div>
            );
          })()}
        </div>
      </div>
    </div>
  );
}
