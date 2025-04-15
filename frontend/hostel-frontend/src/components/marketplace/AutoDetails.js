// frontend/hostel-frontend/src/components/marketplace/AutoDetails.js
import React from 'react';
import { useTranslation } from 'react-i18next';
import {
    Box,
    Grid,
    Typography,
    Paper,
    Divider,
    Stack,
    Chip
} from '@mui/material';
import {
    Car,
    Calendar,
    Gauge,
    Droplet,
    RotateCw,
    Activity,
    PaintBucket,
    CreditCard,
    Truck,
    Ruler,
    Users
} from 'lucide-react';

const AutoDetails = ({ autoProperties }) => {
    const { t } = useTranslation(['marketplace', 'common']);

    if (!autoProperties) {
        return null;
    }

    const formatValue = (value, suffix = '') => {
        if (value === null || value === undefined || value === '') {
            return t('auto.details.not_specified', { defaultValue: 'Не указано' });
        }
        return `${value}${suffix}`;
    };

    const translateFuelType = (fuelType) => {
        if (!fuelType) return t('auto.details.not_specified', { defaultValue: 'Не указано' });
        return t(`auto.properties.fuel_types.${fuelType.toLowerCase()}`, { defaultValue: fuelType });
    };

    const translateTransmission = (transmission) => {
        if (!transmission) return t('auto.details.not_specified', { defaultValue: 'Не указано' });
        return t(`auto.properties.transmissions.${transmission.toLowerCase()}`, { defaultValue: transmission });
    };

    const translateBodyType = (bodyType) => {
        if (!bodyType) return t('auto.details.not_specified', { defaultValue: 'Не указано' });
        return t(`auto.properties.body_types.${bodyType.toLowerCase()}`, { defaultValue: bodyType });
    };

    const translateDriveType = (driveType) => {
        if (!driveType) return t('auto.details.not_specified', { defaultValue: 'Не указано' });
        return t(`auto.properties.drive_types.${driveType.toLowerCase()}`, { defaultValue: driveType });
    };

    const renderDetailItem = (icon, label, value, important = false) => (
        <Grid item xs={6} sm={4} md={3}>
            <Box sx={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                textAlign: 'center',
                p: 1.5,
                borderRadius: 1,
                height: '100%',
                backgroundColor: important ? 'rgba(0, 0, 0, 0.03)' : 'transparent'
            }}>
                <Box sx={{ color: 'primary.main', mb: 1 }}>
                    {icon}
                </Box>
                <Typography variant="caption" color="text.secondary" gutterBottom>
                    {label}
                </Typography>
                <Typography variant="body2" fontWeight={important ? 'medium' : 'normal'}>
                    {value}
                </Typography>
            </Box>
        </Grid>
    );

    return (
        <Paper sx={{ p: 2, mb: 3 }}>
            <Box sx={{ mb: 2 }}>
                <Typography variant="h6" gutterBottom>
                    {t('auto.details.title', { defaultValue: 'Характеристики автомобиля' })}
                </Typography>
                <Stack direction="row" spacing={1} sx={{ mb: 2 }}>
                    <Chip
                        icon={<Car size={16} />}
                        label={`${autoProperties.brand} ${autoProperties.model}`}
                        color="primary"
                        variant="outlined"
                    />
                    <Chip
                        icon={<Calendar size={16} />}
                        label={formatValue(autoProperties.year)}
                        variant="outlined"
                    />
                    {autoProperties.body_type && (
                        <Chip
                            icon={<Truck size={16} />}
                            label={translateBodyType(autoProperties.body_type)}
                            variant="outlined"
                        />
                    )}
                </Stack>
            </Box>

            <Divider sx={{ mb: 2 }} />

            <Grid container spacing={2}>
                {/* Основная информация */}
                {renderDetailItem(
                    <Car size={24} />,
                    t('auto.properties.brand_model', { defaultValue: 'Марка и модель' }),
                    `${autoProperties.brand} ${autoProperties.model}`,
                    true
                )}

                {renderDetailItem(
                    <Calendar size={24} />,
                    t('auto.properties.year', { defaultValue: 'Год выпуска' }),
                    formatValue(autoProperties.year),
                    true
                )}

                {renderDetailItem(
                    <Gauge size={24} />,
                    t('auto.properties.mileage', { defaultValue: 'Пробег' }),
                    formatValue(autoProperties.mileage, ' км'),
                    true
                )}

                {renderDetailItem(
                    <Droplet size={24} />,
                    t('auto.properties.fuel_type', { defaultValue: 'Тип топлива^' }),
                    translateFuelType(autoProperties.fuel_type)
                )}

                {renderDetailItem(
                    <RotateCw size={24} />,
                    t('auto.properties.transmission', { defaultValue: 'Трансмиссия' }),
                    translateTransmission(autoProperties.transmission)
                )}

                {renderDetailItem(
                    <Activity size={24} />,
                    t('auto.properties.engine', { defaultValue: 'Двигатель' }),
                    (autoProperties.engine_capacity !== null || autoProperties.power !== null) ?
                        (autoProperties.engine_capacity !== null ? formatValue(autoProperties.engine_capacity, ' л') : '') +
                        (autoProperties.engine_capacity !== null && autoProperties.power !== null ? ' / ' : '') +
                        (autoProperties.power !== null ? formatValue(autoProperties.power, ' л.с.') : '')
                        : t('auto.details.not_specified', { defaultValue: 'Не указано' })
                )}

                {renderDetailItem(
                    <PaintBucket size={24} />,
                    t('auto.properties.color', { defaultValue: 'Цвет' }),
                    formatValue(autoProperties.color)
                )}

                {renderDetailItem(
                    <Truck size={24} />,
                    t('auto.properties.body_type', { defaultValue: 'Тип кузова' }),
                    translateBodyType(autoProperties.body_type)
                )}

                {renderDetailItem(
                    <CreditCard size={24} />,
                    t('auto.properties.drive_type', { defaultValue: 'Привод' }),
                    translateDriveType(autoProperties.drive_type)
                )}

                {renderDetailItem(
                    <Ruler size={24} />,
                    t('auto.properties.doors', { defaultValue: 'Двери' }),
                    formatValue(autoProperties.number_of_doors)
                )}

                {renderDetailItem(
                    <Users size={24} />,
                    t('auto.properties.seats', { defaultValue: 'Места' }),
                    formatValue(autoProperties.number_of_seats)
                )}
            </Grid>

            {autoProperties.additional_features && (
                <>
                    <Divider sx={{ my: 2 }} />
                    <Box>
                        <Typography variant="subtitle2" gutterBottom>
                            {t('auto.properties.features', { defaultValue: 'Дополнительные особенности' })}
                        </Typography>
                        <Typography variant="body2" sx={{ whiteSpace: 'pre-line' }}>
                            {autoProperties.additional_features}
                        </Typography>
                    </Box>
                </>
            )}
        </Paper>
    );
};

export default AutoDetails;