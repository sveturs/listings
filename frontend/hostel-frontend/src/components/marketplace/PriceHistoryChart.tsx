// src/components/marketplace/PriceHistoryChart.tsx
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { 
    LineChart, 
    Line, 
    XAxis, 
    YAxis, 
    CartesianGrid, 
    Tooltip as RechartsTooltip, 
    Legend,
    ResponsiveContainer, 
    TooltipProps
} from 'recharts';
import { 
    Box, 
    Typography, 
    CircularProgress, 
    Alert 
} from '@mui/material';
import axios from '../../api/axios';

interface PriceHistoryItem {
    effective_from: string;
    effective_to?: string | null;
    price: number;
}

interface ChartDataItem {
    date: string;
    price: number;
    label?: string;
}

interface PriceHistoryChartProps {
    listingId: number | string;
}

interface CustomTooltipProps extends TooltipProps<number, string> {
    active?: boolean;
    payload?: Array<{
        value: number;
        payload: ChartDataItem;
    }>;
    label?: string;
}

const PriceHistoryChart: React.FC<PriceHistoryChartProps> = ({ listingId }) => {
    const { t, i18n } = useTranslation('marketplace');
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [data, setData] = useState<ChartDataItem[]>([]);

    useEffect(() => {
        const fetchPriceHistory = async (): Promise<void> => {
            try {
                setLoading(true);
                const response = await axios.get(`/api/v1/marketplace/listings/${listingId}/price-history`);
                
                if (response.data && response.data.data) {
                    // Преобразуем данные для графика и сортируем по дате от старых к новым
                    const historyData = response.data.data as PriceHistoryItem[];
                    const chartData = historyData
                        .sort((a, b) => new Date(a.effective_from).getTime() - new Date(b.effective_from).getTime())
                        .map(item => ({
                            date: new Date(item.effective_from).toLocaleDateString(i18n.language),
                            price: item.price,
                            // Добавляем метку "текущая цена" для последней записи
                            label: item.effective_to ? undefined : t('listings.details.priceHistory.currentPrice')
                        }));
                    
                    setData(chartData);
                } else {
                    setData([]);
                }
                
                setLoading(false);
            } catch (err) {
                console.error('Error fetching price history:', err);
                setError(t('listings.details.priceHistory.loadError'));
                setLoading(false);
            }
        };

        fetchPriceHistory();
    }, [listingId, t, i18n.language]);

    // Функция форматирования числа с учетом локали
    const formatPrice = (value?: number): string => {
        if (value === undefined) return '';
        return new Intl.NumberFormat(i18n.language, {
            style: 'currency',
            currency: 'RSD',
            maximumFractionDigits: 0,
            notation: 'standard' // Используем стандартную нотацию
        }).format(value);
    };

    // Кастомный тултип для графика
    const CustomTooltip: React.FC<CustomTooltipProps> = ({ active, payload, label }) => {
        if (active && payload && payload.length > 0) {
            return (
                <Box
                    sx={{
                        backgroundColor: 'background.paper',
                        p: 1.5,
                        border: '1px solid',
                        borderColor: 'divider',
                        borderRadius: 1,
                        boxShadow: 1
                    }}
                >
                    <Typography variant="subtitle2">{label}</Typography>
                    <Typography variant="body2" color="primary.main">
                        {formatPrice(payload[0].value)}
                    </Typography>
                    {payload[0].payload.label && (
                        <Typography variant="caption" color="success.main">
                            {payload[0].payload.label}
                        </Typography>
                    )}
                </Box>
            );
        }
        return null;
    };

    if (loading) {
        return (
            <Box display="flex" justifyContent="center" alignItems="center" height={300}>
                <CircularProgress />
            </Box>
        );
    }

    if (error) {
        return (
            <Alert severity="error">{error}</Alert>
        );
    }

    if (data.length === 0) {
        return (
            <Alert severity="info">
                {t('listings.details.priceHistory.noData')}
            </Alert>
        );
    }

    return (
        <Box sx={{ width: '100%', height: 300 }}>
            <ResponsiveContainer>
                <LineChart
                    data={data}
                    margin={{ top: 10, right: 30, left: 50, bottom: 0 }} // Увеличиваем левый отступ
                >
                    <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
                    <XAxis 
                        dataKey="date" 
                        padding={{ left: 10, right: 10 }}
                        tick={{ fontSize: 12 }}
                    />
                    <YAxis 
                        tickFormatter={formatPrice}
                        tick={{ fontSize: 12 }}
                        domain={['dataMin', 'dataMax']} // Автоматическое масштабирование
                    />
                    <RechartsTooltip content={<CustomTooltip />} />
                    <Legend />
                    <Line 
                        type="monotone" 
                        dataKey="price" 
                        name={t('listings.details.price')}
                        stroke="#3f51b5" 
                        activeDot={{ r: 8 }}
                        strokeWidth={2}
                    />
                </LineChart>
            </ResponsiveContainer>
        </Box>
    );
};

export default PriceHistoryChart;