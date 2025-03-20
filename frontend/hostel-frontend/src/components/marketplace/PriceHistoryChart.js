// src/components/marketplace/PriceHistoryChart.js
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { 
    LineChart, 
    Line, 
    XAxis, 
    YAxis, 
    CartesianGrid, 
    Tooltip, 
    Legend,
    ResponsiveContainer 
} from 'recharts';
import { 
    Box, 
    Typography, 
    CircularProgress, 
    Alert 
} from '@mui/material';
import axios from '../../api/axios';

const PriceHistoryChart = ({ listingId }) => {
    const { t, i18n } = useTranslation('marketplace');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [data, setData] = useState([]);

    useEffect(() => {
        const fetchPriceHistory = async () => {
            try {
                setLoading(true);
                const response = await axios.get(`/api/v1/marketplace/listings/${listingId}/price-history`);
                
                if (response.data && response.data.data) {
                    // Преобразуем данные для графика и сортируем по дате от старых к новым
                    const chartData = response.data.data
                        .sort((a, b) => new Date(a.effective_from) - new Date(b.effective_from))
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
    const formatPrice = (value) => {
        return new Intl.NumberFormat(i18n.language, {
            style: 'currency',
            currency: 'RSD',
            maximumFractionDigits: 0,
            notation: 'standard' // Используем стандартную нотацию
        }).format(value);
    };

    // Кастомный тултип для графика
    const CustomTooltip = ({ active, payload, label }) => {
        if (active && payload && payload.length) {
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
                    <Tooltip content={<CustomTooltip />} />
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