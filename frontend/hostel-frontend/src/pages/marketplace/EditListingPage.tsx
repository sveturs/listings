// frontend/hostel-frontend/src/pages/marketplace/EditListingPage.tsx
import React, { useState, useEffect, ChangeEvent, FormEvent } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { useLanguage } from '../../contexts/LanguageContext';
import LocationPicker from '../../components/global/LocationPicker';
import MiniMap from '../../components/maps/MiniMap';
import ImageUploader from '../../components/marketplace/ImageUploader';
import FullscreenMap from '../../components/maps/FullscreenMap';
import { Delete as DeleteIcon } from '@mui/icons-material';
import AttributeFields from '../../components/marketplace/AttributeFields';
import CategorySelect from '../../components/marketplace/CategorySelect';
import {
    Container,
    TextField,
    Button,
    Typography,
    Box,
    Alert,
    Grid,
    FormControlLabel,
    Switch,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Paper,
    Modal,
    IconButton
} from '@mui/material';
import axios from '../../api/axios';

interface ListingImage {
    id: number;
    file_path: string;
    listing_id: number;
}

interface Translation {
    title: string;
    description: string;
}

interface Translations {
    [language: string]: Translation;
}

interface LocationData {
    latitude: number | null;
    longitude: number | null;
    formatted_address: string;
    address_components?: {
        city: string;
        country: string;
    };
}

interface Category {
    id: number;
    name: string;
    parent_id: number | null;
    level: number;
    [key: string]: any;
}

interface Listing {
    id?: number;
    title: string;
    description: string;
    price: number;
    category_id: string | number;
    condition: string;
    location: string;
    city: string;
    country: string;
    show_on_map: boolean;
    latitude: number | null;
    longitude: number | null;
    user_id?: number;
    original_language?: string;
    translations?: Translations;
    images?: ListingImage[];
    attributes?: AttributeValue[];
    [key: string]: any;
}

interface AttributeValue {
    attribute_id: number;
    attribute_name: string;
    attribute_type: string;
    display_name: string;
    display_value?: string;
    text_value?: string;
    numeric_value?: number;
    boolean_value?: boolean;
    value: any;
    listing_id?: number;
}

interface RouteParams {
    id: string;
}

const EditListingPage: React.FC = () => {
    const { t, i18n } = useTranslation('marketplace', 'common');
    const [currentLanguage, setCurrentLanguage] = useState<string>(i18n.language);
    const { id } = useParams<keyof RouteParams>();
    const navigate = useNavigate();
    const { user } = useAuth();
    const { language } = useLanguage();

    const [listing, setListing] = useState<Listing>({
        title: "",
        description: "",
        price: 0,
        category_id: "",
        condition: "new",
        location: "",
        city: "",
        country: "",
        show_on_map: true,
        latitude: null,
        longitude: null
    });

    const [images, setImages] = useState<File[]>([]);
    const [previewUrls, setPreviewUrls] = useState<string[]>([]);
    const [categories, setCategories] = useState<Category[]>([]);
    const [error, setError] = useState<string>("");
    const [success, setSuccess] = useState<boolean>(false);
    const [showExpandedMap, setShowExpandedMap] = useState<boolean>(false);
    const [loading, setLoading] = useState<boolean>(true);
    const [attributeValues, setAttributeValues] = useState<AttributeValue[]>([]);
    const [attributesLoaded, setAttributesLoaded] = useState<boolean>(false);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [listingResponse, categoriesResponse] = await Promise.all([
                    axios.get(`/api/v1/marketplace/listings/${id}`),
                    axios.get("/api/v1/marketplace/categories")
                ]);

                const listingData = listingResponse.data.data;

                if (listingData.user_id !== user?.id) {
                    navigate('/marketplace');
                    return;
                }

                // Получаем текст на нужном языке
                const title = i18n.language === listingData.original_language
                    ? listingData.title
                    : listingData.translations?.[i18n.language]?.title || listingData.title;

                const description = i18n.language === listingData.original_language
                    ? listingData.description
                    : listingData.translations?.[i18n.language]?.description || listingData.description;

                setListing({
                    ...listingData,
                    title,
                    description,
                    price: listingData.price,
                    category_id: listingData.category_id,
                    condition: listingData.condition,
                    location: listingData.location,
                    city: listingData.city,
                    country: listingData.country,
                    show_on_map: listingData.show_on_map,
                    latitude: listingData.latitude,
                    longitude: listingData.longitude
                });

                if (listingData.images) {
                    setPreviewUrls(listingData.images.map((img: ListingImage) =>
                        `${process.env.REACT_APP_BACKEND_URL}/uploads/${img.file_path}`
                    ));
                }

                // Просто сохраняем категории, CategorySelect компонент сам обработает переводы
                setCategories(categoriesResponse.data.data || []);

                // Обработка атрибутов объявления, если они есть
                if (listingData.attributes && listingData.attributes.length > 0) {
                    console.log("Загружаем атрибуты объявления:", listingData.attributes);

                    const formattedAttributes = listingData.attributes.map((attr: any) => {
                        // Для отладки выводим каждый атрибут
                        console.log(`Атрибут ${attr.attribute_name} (${attr.attribute_type}): ${attr.display_value}`);

                        // Создаем основную структуру атрибута
                        const formattedAttr: AttributeValue = {
                            attribute_id: attr.attribute_id,
                            attribute_name: attr.attribute_name,
                            attribute_type: attr.attribute_type,
                            display_name: attr.display_name,
                            display_value: attr.display_value,
                            listing_id: listingData.id,
                            value: '' // Значение будет установлено ниже в зависимости от типа
                        };

                        // Устанавливаем значение в зависимости от типа атрибута
                        switch (attr.attribute_type) {
                            case 'text':
                            case 'select':
                                formattedAttr.text_value = attr.text_value || '';
                                formattedAttr.value = attr.text_value || attr.display_value || '';
                                break;
                            case 'number':
                                const numValue = attr.numeric_value !== null ? attr.numeric_value :
                                    (attr.display_value ? parseFloat(attr.display_value) : 0);

                                formattedAttr.numeric_value = numValue;
                                formattedAttr.value = numValue;

                                // Особый случай для года - проверяем на валидность
                                if (attr.attribute_name === 'year' && (numValue < 1900 || numValue > new Date().getFullYear() + 1)) {
                                    const currentYear = new Date().getFullYear();
                                    formattedAttr.numeric_value = currentYear;
                                    formattedAttr.value = currentYear;
                                    console.log(`Корректировка невалидного года ${numValue} -> ${currentYear}`);
                                }
                                break;
                            case 'boolean':
                                // Преобразуем разные форматы в булево значение
                                let boolValue = false;
                                if (attr.boolean_value !== null && attr.boolean_value !== undefined) {
                                    boolValue = !!attr.boolean_value;
                                } else if (attr.display_value) {
                                    boolValue = attr.display_value.toLowerCase() === 'да' ||
                                        attr.display_value.toLowerCase() === 'true' ||
                                        attr.display_value === '1';
                                }
                                formattedAttr.boolean_value = boolValue;
                                formattedAttr.value = boolValue;
                                break;
                            default:
                                // Если тип не определен, используем отображаемое значение
                                formattedAttr.value = attr.display_value || '';
                        }

                        return formattedAttr;
                    });

                    console.log("Устанавливаем атрибуты:", formattedAttributes);
                    setAttributeValues(formattedAttributes);
                } else {
                    console.log("Атрибуты не найдены в объявлении");
                }

                setLoading(false);
                setAttributesLoaded(true);
            } catch (err) {
                console.error("Ошибка при загрузке данных:", err);
                setError(t('listings.edit.errors.loadFailed'));
                setLoading(false);
                setAttributesLoaded(true);
            }
        };

        if (user?.id) {
            fetchData();
        }
    }, [id, user, navigate, t, i18n.language]);

    // Добавляем эффект для отслеживания изменения языка
    useEffect(() => {
        const updateContent = async () => {
            try {
                const response = await axios.get(`/api/v1/marketplace/listings/${id}`);
                const listingData = response.data.data;

                // Получаем текст на выбранном языке
                const title = i18n.language === listingData.original_language
                    ? listingData.title
                    : listingData.translations?.[i18n.language]?.title || listingData.title;

                const description = i18n.language === listingData.original_language
                    ? listingData.description
                    : listingData.translations?.[i18n.language]?.description || listingData.description;

                setListing(prev => ({
                    ...prev,
                    title,
                    description
                }));
            } catch (error) {
                console.error('Error updating content for new language:', error);
            }
        };

        if (currentLanguage !== i18n.language) {
            updateContent();
            setCurrentLanguage(i18n.language);
        }
    }, [i18n.language, id, currentLanguage]);

    const handleLocationSelect = (location: LocationData) => {
        setListing(prev => ({
            ...prev,
            latitude: location.latitude,
            longitude: location.longitude,
            location: location.formatted_address,
            city: location.address_components?.city || '',
            country: location.address_components?.country || ''
        }));
    };

    // Отслеживаем изменения атрибутов и выводим их в консоль для отладки
    useEffect(() => {
        if (attributeValues.length > 0) {
            console.log("Текущие значения атрибутов:", attributeValues);
        }
    }, [attributeValues]);

    // Преобразует массив атрибутов для отправки на сервер
    const prepareAttributesForSubmission = (attributes: AttributeValue[]): AttributeValue[] => {
        return attributes.map(attr => {
            // Создаем копию атрибута для изменения
            const newAttr = { ...attr };

            // Обеспечиваем правильный формат числовых значений
            if (attr.attribute_type === 'number' && attr.value !== undefined) {
                // Преобразуем в число и убеждаемся, что оно сохранено в правильном поле
                let numValue: number;

                // Обрабатываем дробные значения для объема двигателя
                if (attr.attribute_name === 'engine_capacity') {
                    // Для объема двигателя используем parseFloat для сохранения дробной части
                    numValue = parseFloat(attr.value);
                    // Округляем до одного знака после запятой
                    if (!isNaN(numValue)) {
                        numValue = Math.round(numValue * 10) / 10;
                    }
                } else if (attr.attribute_name === 'year') {
                    // Для года выпуска обрабатываем отдельно, чтобы избежать сброса
                    numValue = parseInt(attr.value);
                    // Логируем для отладки
                    console.log(`Подготовка года выпуска для отправки: ${attr.value} -> ${numValue}`);

                    // Проверяем на валидный год
                    if (isNaN(numValue) || numValue < 1900 || numValue > new Date().getFullYear() + 1) {
                        // Если год невалидный, устанавливаем текущий год
                        const currentYear = new Date().getFullYear();
                        console.warn(`Невалидный год ${numValue}, устанавливаем текущий год ${currentYear}`);
                        numValue = currentYear;
                    }
                } else {
                    // Для других числовых полей по-прежнему используем parseFloat
                    numValue = parseFloat(attr.value);
                }

                if (!isNaN(numValue)) {
                    // Гарантируем, что numeric_value всегда устанавливается для числовых атрибутов
                    newAttr.numeric_value = numValue;
                    // Обновляем основное значение для согласованности
                    newAttr.value = numValue;

                    // Также обновляем отображаемое значение
                    if (attr.attribute_name === 'year') {
                        newAttr.display_value = String(numValue);
                    } else {
                        newAttr.display_value = String(numValue);
                    }
                } else {
                    console.error(`Не удалось преобразовать "${attr.value}" в число для атрибута ${attr.attribute_name}`);
                }
            } else if (attr.attribute_type === 'boolean') {
                // Для булевых значений
                newAttr.boolean_value = Boolean(attr.value);
                newAttr.value = newAttr.boolean_value;
                newAttr.display_value = newAttr.boolean_value ? 'Да' : 'Нет';
            } else if (attr.attribute_type === 'text' || attr.attribute_type === 'select') {
                // Для текстовых значений
                newAttr.text_value = String(attr.value || '');
                newAttr.value = newAttr.text_value;
                newAttr.display_value = newAttr.text_value;
            }

            return newAttr;
        });
    };

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError("");
        setSuccess(false);

        try {
            if (attributeValues.length === 0 && attributesLoaded) {
                console.log("Атрибуты отсутствуют, но страница загружена. Пробуем повторно получить атрибуты для категории", listing.category_id);

                try {
                    const attributesResponse = await axios.get(`/api/v1/marketplace/categories/${listing.category_id}/attributes`);
                    if (attributesResponse.data?.data && attributesResponse.data.data.length > 0) {
                        console.log("Получены атрибуты категории:", attributesResponse.data.data);

                        // Инициализируем атрибуты со значениями по умолчанию, так как исходные значения потеряны
                        const defaultAttributes = attributesResponse.data.data.map((attr: any) => {
                            const attrValue: AttributeValue = {
                                attribute_id: attr.id,
                                attribute_name: attr.name,
                                attribute_type: attr.attribute_type,
                                display_name: attr.display_name,
                                value: ''
                            };

                            switch (attr.attribute_type) {
                                case 'text':
                                case 'select':
                                    attrValue.text_value = '';
                                    break;
                                case 'number':
                                    attrValue.numeric_value = 0;
                                    break;
                                case 'boolean':
                                    attrValue.boolean_value = false;
                                    break;
                            }

                            return attrValue;
                        });

                        console.log("Инициализированы атрибуты по умолчанию:", defaultAttributes);
                        // Обновляем атрибуты для отправки
                        setAttributeValues(defaultAttributes);
                    }
                } catch (error) {
                    console.error("Ошибка при получении атрибутов категории:", error);
                }
            }

            const processedAttributes = prepareAttributesForSubmission(attributeValues);
            console.log("Отправляем на сервер атрибуты:", processedAttributes);

            // Подготавливаем данные для обновления
            const listingData: Listing = {
                ...listing,
                price: parseFloat(listing.price.toString()),
                attributes: processedAttributes
            };

            // Всегда обновляем основные данные объявления
            await axios.put(`/api/v1/marketplace/listings/${id}`, listingData);

            // Если язык не совпадает с оригиналом, также обновляем переводы
            if (i18n.language !== listing.original_language) {
                // Получаем провайдер перевода из локального хранилища или используем значение по умолчанию
                const translationProvider = localStorage.getItem('preferredTranslationProvider') || 'google';
                
                await axios.put(`/api/v1/marketplace/translations/${id}?translation_provider=${translationProvider}`, {
                    language: i18n.language,
                    translations: {
                        title: listing.title,
                        description: listing.description
                    },
                    is_verified: true,
                    provider: translationProvider // Добавляем информацию о провайдере перевода
                });
            }

            // Отправляем новые изображения, если они есть
            if (images.length > 0) {
                const formData = new FormData();
                images.forEach((file) => {
                    formData.append('images', file);
                });

                await axios.post(
                    `/api/v1/marketplace/listings/${id}/images`,
                    formData,
                    {
                        headers: {
                            'Content-Type': 'multipart/form-data'
                        }
                    }
                );
            }

            setSuccess(true);
            setTimeout(() => {
                navigate(`/marketplace/listings/${id}`);
            }, 1500);
        } catch (error) {
            console.error("Ошибка при обновлении:", error);
            setError(t('listings.edit.errors.updateFailed'));
        }
    };
  
    if (loading) {
        return (
            <Container maxWidth="md">
                <Box sx={{ mt: 4, textAlign: 'center' }}>
                    <Typography>{t('listings.edit.loading')}</Typography>
                </Box>
            </Container>
        );
    }

    return (
        <Container maxWidth="md">
            <Box sx={{ mt: 4, mb: 4 }}>
                <Typography variant="h4" gutterBottom>
                    {t('listings.edit.title')}
                </Typography>

                {/* Новый информационный блок о языке редактирования */}
                <Alert severity="info" sx={{ mb: 2 }}>
                  {t('listings.edit.languageInfo', {
                    currentLanguage: i18n.language === 'ru' ? 'русском' : 
                                     i18n.language === 'en' ? 'английском' : 'сербском',
                    originalLanguage: listing?.original_language === 'ru' ? 'русском' :
                                      listing?.original_language === 'en' ? 'английском' : 'сербском',
                    defaultValue: 'Вы редактируете объявление на {{currentLanguage}} языке. Оригинал объявления на {{originalLanguage}} языке.'
                  })}
                  <br />
                  {i18n.language !== listing?.original_language && (
                    t('listings.edit.translationNote', {
                      defaultValue: 'Изменения будут сохранены как перевод, оригинал останется без изменений.'
                    })
                  )}
                </Alert>
                {error && (
                    <Alert severity="error" sx={{ mb: 2 }}>
                        {error}
                    </Alert>
                )}

                {success && (
                    <Alert severity="success" sx={{ mb: 2 }}>
                        {t('listings.edit.success')}
                    </Alert>
                )}

                <Paper sx={{ p: 3 }}>
                    <form onSubmit={handleSubmit}>
                        <Grid container spacing={3}>
                            <Grid item xs={12}>
                                <TextField
                                    label={t('listings.create.name')}
                                    fullWidth
                                    required
                                    value={listing.title || ''}
                                    onChange={(e: ChangeEvent<HTMLInputElement>) => setListing({ ...listing, title: e.target.value })}
                                />
                            </Grid>

                            <Grid item xs={12}>
                                <TextField
                                    label={t('listings.create.description')}
                                    fullWidth
                                    required
                                    multiline
                                    rows={4}
                                    value={listing.description || ''}
                                    onChange={(e: ChangeEvent<HTMLInputElement>) => setListing({ ...listing, description: e.target.value })}
                                />
                            </Grid>

                            <Grid item xs={12} sm={6}>
                                <TextField
                                    label={t('listings.create.price')}
                                    type="number"
                                    fullWidth
                                    required
                                    value={listing.price || ''}
                                    onChange={(e: ChangeEvent<HTMLInputElement>) => setListing({ ...listing, price: Number(e.target.value) })}
                                />
                            </Grid>

                            <Grid item xs={12} sm={6}>
                                <FormControl fullWidth required error={!listing.category_id}>
                                    <InputLabel
                                        shrink
                                        sx={{
                                            backgroundColor: 'background.paper',
                                            px: 1,
                                            '&.MuiInputLabel-shrink': {
                                                transform: 'translate(14px, -6px) scale(0.75)'
                                            }
                                        }}
                                        children={t('listings.create.category')}
                                    />
                                    {/* Используем CategorySelect с параметром displaySelectedValueOnly */}
                                    <CategorySelect
                                        categories={categories}
                                        value={listing.category_id}
                                        onChange={(value: string | number) => setListing({ ...listing, category_id: value })}
                                        error={!listing.category_id}
                                    />
                                </FormControl>
                            </Grid>

                            <Grid item xs={12} sm={6}>
                                <FormControl fullWidth required>
                                    <InputLabel children={t('listings.create.condition.label')} />
                                    <Select
                                        value={listing.condition || 'new'}
                                        onChange={(e) => setListing({ ...listing, condition: e.target.value as string })}
                                    >
                                        <MenuItem value="new">{t('listings.create.condition.new')}</MenuItem>
                                        <MenuItem value="used">{t('listings.create.condition.used')}</MenuItem>
                                    </Select>
                                </FormControl>
                            </Grid>

                            <Grid item xs={12}>
                                <Typography variant="h6" gutterBottom>
                                    {t('listings.edit.location.title')}
                                </Typography>
                                <LocationPicker
                                    onLocationSelect={handleLocationSelect}
                                    initialLocation={{
                                        latitude: listing.latitude,
                                        longitude: listing.longitude,
                                        formatted_address: listing.location,
                                        address_components: {
                                            city: listing.city || '',
                                            country: listing.country || ''
                                        }
                                    }}
                                />

                                {listing.latitude && listing.longitude && (
                                    <Box sx={{ mt: 2 }}>
                                        <MiniMap
                                            latitude={listing.latitude}
                                            longitude={listing.longitude}
                                            address={listing.location}
                                            onExpand={() => setShowExpandedMap(true)}
                                        />
                                    </Box>
                                )}

                                <FormControlLabel
                                    control={
                                        <Switch
                                            checked={listing.show_on_map}
                                            onChange={(e) => setListing(prev => ({
                                                ...prev,
                                                show_on_map: e.target.checked
                                            }))}
                                        />
                                    }
                                    label={t('listings.edit.location.showOnMap')}
                                    sx={{ mt: 1 }}
                                />
                            </Grid>

                            <Grid item xs={12}>
                                <Typography variant="h6" gutterBottom>
                                    {t('listings.edit.photos.title')}
                                </Typography>
                                <ImageUploader
                                    onImagesSelected={(processedImages: any[]) => {
                                        setImages(processedImages.map(img => img.file));
                                        setPreviewUrls(processedImages.map(img => img.preview));
                                    }}
                                    maxImages={10}
                                    maxSizeMB={1}
                                />

                                <Box sx={{ mt: 2, display: 'flex', gap: 2, flexWrap: 'wrap' }}>
                                    {previewUrls.map((url, index) => (
                                        <Box
                                            key={index}
                                            sx={{ position: 'relative', width: 100, height: 100 }}
                                        >
                                            <img
                                                src={url}
                                                alt={t('listings.edit.photos.preview', { index: index + 1 })}
                                                style={{
                                                    width: '100%',
                                                    height: '100%',
                                                    objectFit: 'cover',
                                                    borderRadius: '4px'
                                                }}
                                            />
                                            <IconButton
                                                size="small"
                                                sx={{
                                                    position: 'absolute',
                                                    top: -10,
                                                    right: -10,
                                                    bgcolor: 'background.paper'
                                                }}
                                                onClick={() => {
                                                    setImages(prev => prev.filter((_, i) => i !== index));
                                                    setPreviewUrls(prev => prev.filter((_, i) => i !== index));
                                                    URL.revokeObjectURL(url);
                                                }}
                                            >
                                                <DeleteIcon />
                                            </IconButton>
                                        </Box>
                                    ))}
                                </Box>
                            </Grid>
                            <Grid item xs={12}>
                                {listing.category_id && (
                                    <AttributeFields
                                        key={`attr-fields-${attributesLoaded ? 'loaded' : 'loading'}-${listing.category_id}`}
                                        categoryId={listing.category_id}
                                        value={attributeValues}
                                        onChange={(newValues: AttributeValue[]) => {
                                            console.log("Новые значения атрибутов:", newValues);
                                            setAttributeValues(newValues);
                                        }}
                                    />
                                )}
                            </Grid>
                            <Grid item xs={12}>
                                <Box sx={{ display: 'flex', gap: 2 }}>
                                    <Button
                                        type="submit"
                                        variant="contained"
                                        color="primary"
                                        fullWidth
                                        size="large"
                                        disabled={!listing.title || !listing.description || !listing.category_id || !listing.price || Number(listing.price) <= 0}
                                    >
                                        {t('listings.edit.saveChanges')}
                                    </Button>
                                    <Button
                                        variant="outlined"
                                        fullWidth
                                        size="large"
                                        onClick={() => navigate(`/marketplace/listings/${id}`)}
                                    >
                                        {t('buttons.cancel', { ns: 'common' })}
                                    </Button>
                                </Box>
                            </Grid>
                        </Grid>
                    </form>
                </Paper>
            </Box>

            {showExpandedMap && listing.latitude && listing.longitude && (
                <Modal
                    open={showExpandedMap}
                    onClose={() => setShowExpandedMap(false)}
                    keepMounted
                    children={
                        <Paper
                            sx={{
                                position: 'relative',
                                width: '100%',
                                maxWidth: 1200,
                                maxHeight: '90vh',
                                overflow: 'hidden'
                            }}
                        >
                            <FullscreenMap
                                latitude={listing.latitude}
                                longitude={listing.longitude}
                                title={listing.title}
                            />
                        </Paper>
                    }
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        p: 2
                    }}
                />
            )}
        </Container>
    );
};

export default EditListingPage;