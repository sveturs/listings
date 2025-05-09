// frontend/hostel-frontend/src/components/store/ImportSourceList.tsx
import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
    Box,
    Typography,
    Card,
    CardHeader,
    CardContent,
    Divider,
    IconButton,
    Button,
    Chip,
    Grid,
    Paper,
    Badge,
    Tooltip,
    Stack
} from '@mui/material';
import {
    Trash2,
    Edit,
    Database,
    Upload,
    BarChart,
    Clock,
    RefreshCw,
    Plus,
    ExternalLink,
    Tag
} from 'lucide-react';
import AddImportSourceModal from './AddImportSourceModal';
import ImportModal from './ImportModal';
import CategoryMappingModal from './CategoryMappingModal';

export interface ImportSource {
    id: string;
    type: string;
    url?: string;
    schedule?: string;
    created_at: string;
    last_import_status?: string;
    last_import_at?: string;
}

interface ImportSourceListProps {
    sources: ImportSource[];
    storefrontId: string;
    onUpdate: () => void;
    onDelete: (sourceId: string) => Promise<void>;
    onFetchHistory: (sourceId: string) => void;
    onRunDirectSync: (sourceId: string) => void;
}

const ImportSourceList: React.FC<ImportSourceListProps> = ({ 
    sources, 
    storefrontId, 
    onUpdate, 
    onDelete, 
    onFetchHistory, 
    onRunDirectSync 
}) => {

    const { t } = useTranslation(['common', 'marketplace']);

    const [openAddModal, setOpenAddModal] = useState<boolean>(false);
    const [openImportModal, setOpenImportModal] = useState<boolean>(false);
    const [openMappingModal, setOpenMappingModal] = useState<boolean>(false);
    const [selectedSource, setSelectedSource] = useState<ImportSource | null>(null);
    const [editingSource, setEditingSource] = useState<ImportSource | null>(null);

    const handleEdit = (source: ImportSource): void => {
        setEditingSource(source);
        setOpenAddModal(true);
    };

    const handleAdd = (): void => {
        setEditingSource(null);
        setOpenAddModal(true);
    };

    const handleImport = (source: ImportSource): void => {
        setSelectedSource(source);
        setOpenImportModal(true);
    };

    const handleCategoryMapping = (source: ImportSource): void => {
        setSelectedSource(source);
        setOpenMappingModal(true);
    };

    const handleMappingSave = (): void => {
        onUpdate();
    };

    const handleDelete = async (sourceId: string): Promise<void> => {
        if (window.confirm(t('marketplace:store.import.deleteConfirm'))) {
            await onDelete(sourceId);
        }
    };

    const formatScheduleLabel = (schedule: string | undefined): string => {
        if (!schedule) return t('marketplace:store.import.scheduleNone', { defaultValue: 'Только вручную' });

        switch (schedule.toLowerCase()) {
            case 'hourly':
                return t('marketplace:store.import.scheduleHourly', { defaultValue: 'Каждый час' });
            case 'daily':
                return t('marketplace:store.import.scheduleDaily', { defaultValue: 'Ежедневно' });
            case 'weekly':
                return t('marketplace:store.import.scheduleWeekly', { defaultValue: 'Еженедельно' });
            case 'monthly':
                return t('marketplace:store.import.scheduleMonthly', { defaultValue: 'Ежемесячно' });
            default:
                return schedule;
        }
    };

    const getStatusColor = (status: string | undefined): 'success' | 'warning' | 'error' | 'default' => {
        if (!status) return 'default';

        switch (status.toLowerCase()) {
            case 'success':
                return 'success';
            case 'partial':
                return 'warning';
            case 'failed':
                return 'error';
            default:
                return 'default';
        }
    };

    return (
        <Box mb={4}>
            <Box mb={3} display="flex" justifyContent="space-between" alignItems="center">
                <Typography variant="h6" gutterBottom>
                    {t('marketplace:store.import.sources')}
                </Typography>

                <Button
                    variant="contained"
                    startIcon={<Plus />}
                    onClick={handleAdd}
                >
                    {t('marketplace:store.import.addSource')}
                </Button>
            </Box>

            {sources.length === 0 ? (
                <Paper sx={{ p: 4, textAlign: 'center' }}>
                    <Database size={64} stroke={1} style={{ margin: '20px auto', opacity: 0.5 }} />
                    <Typography variant="h6" gutterBottom>
                        {t('marketplace:store.import.noSources')}
                    </Typography>
                    <Typography variant="body1" color="text.secondary" paragraph>
                        {t('marketplace:store.import.createFirstSource')}
                    </Typography>

                </Paper>
            ) : (
                <Grid container spacing={3}>
                    {sources.map((source) => (
                        <Grid item xs={12} md={6} key={source.id}>
                            <Card sx={{ height: '100%' }}>
                                <CardHeader
                                    title={
                                        <Box display="flex" alignItems="center" gap={1}>
                                            <Typography variant="subtitle1">
                                                {t('marketplace:store.import.sourceType', { type: source.type.toUpperCase() })}
                                            </Typography>
                                            {source.schedule && (
                                                <Tooltip title={t('marketplace:store.import.scheduledImport')}>
                                                    <Badge color="primary" variant="dot">
                                                        <Clock size={16} />
                                                    </Badge>
                                                </Tooltip>
                                            )}
                                        </Box>
                                    }
                                    subheader={
                                        <>
                                            <Typography variant="caption" color="text.secondary">
                                                {t('marketplace:store.import.created', { date: new Date(source.created_at).toLocaleDateString() })}
                                            </Typography>
                                            {source.schedule && (
                                                <Chip
                                                    size="small"
                                                    label={formatScheduleLabel(source.schedule)}
                                                    color="primary"
                                                    variant="outlined"
                                                    sx={{ ml: 1 }}
                                                />
                                            )}
                                        </>
                                    }
                                    action={
                                        <Box>
                                            <IconButton onClick={() => handleEdit(source)} title={t('common:buttons.edit')}>
                                                <Edit size={18} />
                                            </IconButton>
                                            <IconButton onClick={() => handleDelete(source.id)} title={t('common:buttons.delete')}>
                                                <Trash2 size={18} />
                                            </IconButton>
                                        </Box>
                                    }
                                />
                                <Divider />
                                <CardContent>
                                    <Stack spacing={2}>
                                        <Box>
                                            <Typography variant="subtitle2">
                                                {t('marketplace:store.import.sourceDetails')}
                                            </Typography>

                                            {source.url ? (
                                                <Box display="flex" alignItems="center" gap={1}>
                                                    <ExternalLink size={14} />
                                                    <Typography variant="body2" component="a" href={source.url} target="_blank" rel="noopener noreferrer" sx={{ wordBreak: 'break-all' }}>
                                                        {source.url}
                                                    </Typography>
                                                </Box>
                                            ) : (
                                                <Typography variant="body2" color="text.secondary">
                                                    {t('marketplace:store.import.manualUpload')}
                                                </Typography>
                                            )}
                                        </Box>

                                        <Box>
                                            <Typography variant="subtitle2">
                                                {t('marketplace:store.import.importStatus')}
                                            </Typography>

                                            <Box display="flex" alignItems="center" gap={2}>
                                                <Chip
                                                    label={source.last_import_status || t('marketplace:store.import.statusNone')}
                                                    size="small"
                                                    color={getStatusColor(source.last_import_status)}
                                                />

                                                <Typography variant="body2" color="text.secondary">
                                                    {source.last_import_at
                                                        ? t('marketplace:store.import.lastRunAt', {
                                                            date: new Date(source.last_import_at).toLocaleString()
                                                        })
                                                        : t('marketplace:store.import.never')
                                                    }
                                                </Typography>
                                            </Box>
                                        </Box>

                                        <Divider />
                                        <Box display="flex" justifyContent="space-between" flexWrap="wrap" gap={1}>
                                            <Button
                                                variant="outlined"
                                                size="small"
                                                startIcon={<BarChart />}
                                                onClick={() => onFetchHistory(source.id)}
                                            >
                                                {t('marketplace:store.import.history')}
                                            </Button>

                                            <Button
                                                variant="outlined"
                                                size="small"
                                                startIcon={<Tag />}
                                                onClick={() => handleCategoryMapping(source)}
                                                title={t('marketplace:store.categoryMapping.buttonTitle')}
                                            >
                                                {t('marketplace:store.categoryMapping.button')}
                                            </Button>

                                            {source.url && (
                                                <Button
                                                    variant="outlined"
                                                    size="small"
                                                    startIcon={<RefreshCw />}
                                                    onClick={() => onRunDirectSync(source.id)}
                                                    title={t('marketplace:store.import.runSync')}
                                                >
                                                    {t('marketplace:store.import.sync')}
                                                </Button>
                                            )}

                                            <Button
                                                variant="contained"
                                                size="small"
                                                startIcon={<Upload />}
                                                onClick={() => handleImport(source)}
                                            >
                                                {t('marketplace:store.import.importData')}
                                            </Button>
                                        </Box>

                                    </Stack>
                                </CardContent>
                            </Card>
                        </Grid>
                    ))}
                </Grid>
            )}

            {/* Модальное окно добавления/редактирования источника */}
            <AddImportSourceModal
                open={openAddModal}
                onClose={() => setOpenAddModal(false)}
                onSuccess={onUpdate}
                storefrontId={storefrontId}
                initialData={editingSource}
            />

            {/* Модальное окно импорта */}
            <ImportModal
                open={openImportModal && selectedSource !== null}
                onClose={() => {
                    setOpenImportModal(false);
                    setSelectedSource(null);
                }}
                sourceId={selectedSource?.id}
                onSuccess={(result) => {
                    onUpdate();
                    if (selectedSource) {
                        onFetchHistory(selectedSource.id);
                    }
                }}
            />

            {/* Модальное окно сопоставления категорий */}
            <CategoryMappingModal
                open={openMappingModal && selectedSource !== null}
                onClose={() => {
                    setOpenMappingModal(false);
                    setSelectedSource(null);
                }}
                sourceId={selectedSource?.id}
                onSave={handleMappingSave}
            />
        </Box>
    );
};

export default ImportSourceList;