// frontend/hostel-frontend/src/pages/marketplace/CreateListingPage.js
import React, { useState, useEffect } from "react";
import { useTranslation } from 'react-i18next';

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
    IconButton,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Paper,
    useTheme,
    useMediaQuery
} from "@mui/material";
import { useNavigate } from 'react-router-dom';
import { Delete as DeleteIcon, CloudUpload as CloudUploadIcon } from '@mui/icons-material';
import LocationPicker from '../../components/global/LocationPicker';
import MiniMap from '../../components/maps/MiniMap';
import { GoogleMap, Marker } from '@react-google-maps/api';
import axios from "../../api/axios";
import { useLanguage } from '../../contexts/LanguageContext';
import ImageUploader from '../../components/marketplace/ImageUploader';
import CategorySelect from '../../components/marketplace/CategorySelect';
import { ChevronRight, ChevronLeft } from 'lucide-react';
import { useLocation } from '../../contexts/LocationContext';
import AutoPropertiesForm from '../../components/marketplace/AutoPropertiesForm';

const CreateListing = () => {
    const { t, i18n } = useTranslation('marketplace');
    const theme = useTheme();
    const { language } = useLanguage();
    const { userLocation } = useLocation();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
    const navigate = useNavigate();

    const [listing, setListing] = useState({
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
        original_language: 'sr',
        longitude: null
    });
    const [autoProperties, setAutoProperties] = useState({
        brand: "",
        model: "",
        year: "",
        mileage: "",
        fuel_type: "",
        transmission: "",
        engine_capacity: "",
        power: "",
        color: "",
        body_type: "",
        drive_type: "",
        number_of_doors: "",
        number_of_seats: "",
        additional_features: ""
    });
    const [isAutoCategory, setIsAutoCategory] = useState(false);
    const [images, setImages] = useState([]);
    const [previewUrls, setPreviewUrls] = useState([]);
    const [categories, setCategories] = useState([]);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState(false);
    const [showExpandedMap, setShowExpandedMap] = useState(false);
    const [locationWarning, setLocationWarning] = useState(false);
    const checkAutoCategory = async (categoryId) => {
        console.log(`Проверка категории ${categoryId} на принадлежность к авто`);
    
        if (!categoryId) {
            setIsAutoCategory(false);
            return;
        }
    
        // Жесткая проверка на известные автомобильные категории
        const autoRootCategory = 2000;
        const autoCategoriesIds = [2000, 2100, 2200, 2210, 2220, 2230, 2240, 2300, 2310, 2315, 2320, 2325, 2330, 2335, 2340, 2345, 2350, 2355, 2360, 2365];
        
        if (autoCategoriesIds.includes(Number(categoryId))) {
            console.log(`Категория ${categoryId} является автомобильной (локальная проверка)`);
            setIsAutoCategory(true);
            return;
        }
        
        // Проверяем, является ли категория потомком автомобильной категории
        // Для этого нужно узнать родительскую категорию текущей категории
        const categoryPath = categories.find(cat => cat.id === Number(categoryId))?.path;
        if (categoryPath && categoryPath.includes(`/${autoRootCategory}/`)) {
            console.log(`Категория ${categoryId} является потомком автомобильной категории`);
            setIsAutoCategory(true);
            return;
        }
        
        // Если не можем определить локально, делаем запрос
        try {
            console.log(`Отправка запроса на сервер для проверки категории ${categoryId}`);
            const response = await axios.get(`/api/v1/auto/category?category_id=${categoryId}`);
            console.log(`Ответ сервера:`, response.data);
    
            // Дополнительная проверка ответа
            if (response.data && response.data.data && 
                typeof response.data.data.is_auto === 'boolean') {
                setIsAutoCategory(response.data.data.is_auto);
                console.log(`Категория ${categoryId} ${response.data.data.is_auto ? 'является' : 'не является'} автомобильной (серверная проверка)`);
            } else {
                console.log(`Получен неожиданный формат ответа, устанавливаем isAutoCategory = false`);
                setIsAutoCategory(false);
            }
        } catch (err) {
            console.error('Ошибка проверки категории:', err);
            setIsAutoCategory(false);
        }
    };

    const getTranslatedText = (field) => {
        if (!listing) return '';

        // Если текущий язык совпадает с языком оригинала
        if (language === listing.original_language) {
            return listing[field];
        }

        // Пытаемся получить перевод
        if (listing.translations &&
            listing.translations[language] &&
            listing.translations[language][field]) {
            return listing.translations[language][field];
        }

        // Если перевода нет, возвращаем оригинал
        return listing[field];
    };
    const handleCategoryChange = (categoryId) => {
        setListing(prev => ({
            ...prev,
            category_id: categoryId
        }));

        // Проверяем, является ли категория автомобильной
        checkAutoCategory(categoryId);
    };

    useEffect(() => {
        // Если у нас есть данные о местоположении пользователя, используем их
        if (userLocation && userLocation.lat && userLocation.lon) {
            // Создаем объект местоположения для LocationPicker
            const initialLocation = {
                latitude: userLocation.lat,
                longitude: userLocation.lon,
                formatted_address: userLocation.city ?
                    `${userLocation.city}, ${userLocation.country || 'Serbia'}` :
                    'Your location',
                address_components: {
                    city: userLocation.city || '',
                    country: userLocation.country || 'Serbia'
                }
            };

            // Устанавливаем местоположение в состояние компонента
            setListing(prev => ({
                ...prev,
                latitude: userLocation.lat,
                longitude: userLocation.lon,
                location: userLocation.city ?
                    `${userLocation.city}, ${userLocation.country || 'Serbia'}` :
                    'Your location',
                city: userLocation.city || '',
                country: userLocation.country || 'Serbia'
            }));

        }
    }, [userLocation]);

    useEffect(() => {
        const fetchCategories = async () => {
            try {
                const response = await axios.get("/api/v1/marketplace/categories");
                setCategories(response.data.data || []);
            } catch (err) {
                setError(t('listings.create.error'));
            }
        };
        fetchCategories();
    }, [t]);

    const handleImageChange = (e) => {
        const files = Array.from(e.target.files || []);
        if (files.length === 0) return;

        const validFiles = files.filter(file => {
            if (!file.type.startsWith('image/')) {
                setError(t('listings.create.photos.onlyImages'));
                return false;
            }
            if (file.size > 15 * 1024 * 1024) {
                setError(t('listings.create.photos.maxSize', { size: '15MB' }));
                return false;
            }
            return true;
        });

        if (validFiles.length === 0) return;

        validFiles.forEach(file => {
            const reader = new FileReader();
            reader.onloadend = () => {
                setPreviewUrls(prev => [...prev, reader.result]);
            };
            reader.onerror = (error) => {
                console.error('Error reading file:', error);
            };
            reader.readAsDataURL(file);
        });

        setImages(prev => [...prev, ...validFiles]);
    };

    const handleLocationSelect = (location) => {
        setListing(prev => ({
            ...prev,
            latitude: location.latitude,
            longitude: location.longitude,
            location: location.formatted_address,
            city: location.address_components?.city || '',
            country: location.address_components?.country || ''
        }));

        // Проверяем наличие координат и показываем предупреждение, если их нет
        setLocationWarning(!location.latitude || !location.longitude);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError("");
        setSuccess(false);

        try {
            const listingData = {
                ...listing,
                price: parseFloat(listing.price),
                original_language: i18n.language // Устанавливаем текущий язык интерфейса как язык оригинала
            };

            let listingId;

            // Если это автомобильная категория, используем специальный API
            if (isAutoCategory) {
                const response = await axios.post("/api/v1/auto/listings", {
                    listing: listingData,
                    auto_properties: autoProperties
                });
                listingId = response.data.data.id;
            } else {
                // Используем обычный API для обычных объявлений
                const response = await axios.post("/api/v1/marketplace/listings", listingData);
                listingId = response.data.data.id;
            }

            // Загрузка изображений
            if (images.length > 0) {
                const formData = new FormData();
                images.forEach((image, index) => {
                    formData.append('images', image);
                    if (index === 0) {
                        formData.append('main_image_index', '0');
                    }
                });

                await axios.post(`/api/v1/marketplace/listings/${listingId}/images`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
            }

            setSuccess(true);

            // Сброс состояний
            setListing({
                title: "",
                description: "",
                price: 0,
                category_id: "",
                condition: "new",
                location: "",
                city: "",
                country: "",
                latitude: null,
                longitude: null,
                show_on_map: true,
                original_language: i18n.language
            });

            setAutoProperties({
                brand: "",
                model: "",
                year: "",
                mileage: "",
                fuel_type: "",
                transmission: "",
                engine_capacity: "",
                power: "",
                color: "",
                body_type: "",
                drive_type: "",
                number_of_doors: "",
                number_of_seats: "",
                additional_features: ""
            });

            setImages([]);
            setPreviewUrls([]);

            setTimeout(() => {
                navigate(`/marketplace/listings/${listingId}`);
            }, 1500);

        } catch (error) {
            setError(error.response?.data?.error || t('listings.create.error'));
        }
    };

    return (
        <Container
            maxWidth="md"
            disableGutters={isMobile}
            sx={{
                mx: isMobile ? 0 : 'auto',
                width: isMobile ? '100%' : 'auto'
            }}
        >
            <Box sx={{
                mt: isMobile ? 0 : 4,
                mb: isMobile ? 0 : 4,
                px: 0,
                width: '100%'
            }}>
                {error && (
                    <Alert severity="error" sx={{ mb: 2 }}>
                        {error}
                    </Alert>
                )}

                {success && (
                    <Alert severity="success" sx={{ mb: 2 }}>
                        {t('listings.create.success')}
                    </Alert>
                )}

                <Paper sx={{
                    p: isMobile ? '8px 0' : 3,
                    boxShadow: isMobile ? 'none' : 1,
                    bgcolor: isMobile ? 'transparent' : 'background.paper',
                    width: '100%'
                }}>
                    <form onSubmit={handleSubmit}>
                        <Grid container spacing={isMobile ? 2 : 3}>
                            <Grid item xs={12}>
                                <FormControl fullWidth required error={!listing.category_id}>
                                    <InputLabel shrink>{t('listings.create.category')}</InputLabel>
                                    <CategorySelect
                                        categories={categories}
                                        value={listing.category_id}
                                        onChange={handleCategoryChange}
                                        error={!listing.category_id}
                                    />
                                </FormControl>
                            </Grid>

                            <Grid item xs={12}>
                                <TextField
                                    label={t('listings.create.name')}
                                    fullWidth
                                    required
                                    value={listing.title}
                                    onChange={(e) => setListing({ ...listing, title: e.target.value })}
                                    size={isMobile ? "small" : "medium"}
                                />
                            </Grid>

                            <Grid item xs={12} sx={{ mb: 0.1 }}>
                                <ImageUploader
                                    onImagesSelected={(processedImages) => {
                                        setImages(processedImages.map(img => img.file));
                                        setPreviewUrls(processedImages.map(img => img.preview));
                                    }}
                                    maxImages={10}
                                    maxSizeMB={15}
                                />

                                <Box sx={{
                                    mt: 0.1,
                                    display: 'grid',
                                    gridTemplateColumns: `repeat(auto-fill, minmax(${isMobile ? '80px' : '100px'}, 1fr))`,
                                    gap: 1
                                }}>
                                    {previewUrls.map((url, index) => (
                                        <Box
                                            key={index}
                                            sx={{
                                                position: 'relative',
                                                paddingTop: '100%'
                                            }}
                                        >
                                            <img
                                                src={url}
                                                alt={`${t('listings.create.photos.preview')} ${index + 1}`}
                                                style={{
                                                    position: 'absolute',
                                                    top: 0,
                                                    left: 0,
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
                                                    top: 4,
                                                    right: 4,
                                                    backgroundColor: 'rgba(255,255,255,0.8)',
                                                    '&:hover': {
                                                        backgroundColor: 'rgba(255,255,255,0.9)'
                                                    }
                                                }}
                                                onClick={() => {
                                                    setImages(prev => prev.filter((_, i) => i !== index));
                                                    setPreviewUrls(prev => prev.filter((_, i) => i !== index));
                                                    URL.revokeObjectURL(url);
                                                }}
                                            >
                                                <DeleteIcon fontSize="small" />
                                            </IconButton>
                                        </Box>
                                    ))}
                                </Box>
                            </Grid>

                            <Grid item xs={12}>
                                <TextField
                                    label={t('listings.create.description')}
                                    fullWidth
                                    required
                                    multiline
                                    rows={isMobile ? 3 : 4}
                                    value={listing.description}
                                    onChange={(e) => setListing({ ...listing, description: e.target.value })}
                                    size={isMobile ? "small" : "medium"}
                                />
                            </Grid>

                            <Grid item xs={6}>
                                <TextField
                                    label={t('listings.create.price')}
                                    type="number"
                                    fullWidth
                                    required
                                    value={listing.price}
                                    onChange={(e) => setListing({ ...listing, price: e.target.value })}
                                    size={isMobile ? "small" : "medium"}
                                />
                            </Grid>

                            <Grid item xs={6}>
                                <FormControl fullWidth required size={isMobile ? "small" : "medium"}>
                                    <InputLabel>{t('listings.create.condition.label')}</InputLabel>
                                    <Select
                                        value={listing.condition}
                                        onChange={(e) => setListing({ ...listing, condition: e.target.value })}
                                    >
                                        <MenuItem value="new">{t('listings.create.condition.new')}</MenuItem>
                                        <MenuItem value="used">{t('listings.create.condition.used')}</MenuItem>
                                    </Select>
                                </FormControl>
                            </Grid>

                            {/* Добавляем форму автомобильных свойств, если выбрана автомобильная категория */}
                            {isAutoCategory && (
                                <Grid item xs={12}>
                                    <AutoPropertiesForm
                                        values={autoProperties}
                                        onChange={setAutoProperties}
                                    />
                                </Grid>
                            )}

                            <Grid item xs={12}>
                                <Box sx={{ mb: 1 }}>
                                    <LocationPicker
                                        onLocationSelect={(location) => {
                                            setListing(prev => ({
                                                ...prev,
                                                latitude: location.latitude,
                                                longitude: location.longitude,
                                                location: location.formatted_address,
                                                city: location.address_components?.city || '',
                                                country: location.address_components?.country || ''
                                            }));
                                        }}
                                        initialLocation={{
                                            latitude: listing.latitude,
                                            longitude: listing.longitude,
                                            formatted_address: listing.location,
                                            address_components: {
                                                city: listing.city,
                                                country: listing.country
                                            }
                                        }}
                                    />
                                </Box>

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
                                    label={t('listings.create.location.showOnMap')}
                                    sx={{ mt: 1 }}
                                />
                            </Grid>

                            <Grid item xs={12}>
                                <Button
                                    id="createAnnouncementButton"
                                    type="submit"
                                    variant="contained"
                                    color="primary"
                                    fullWidth
                                    size={isMobile ? "large" : "large"}
                                    disabled={
                                        !listing.title ||
                                        !listing.description ||
                                        !listing.category_id ||
                                        listing.price <= 0 ||
                                        // Добавляем проверку обязательных полей для автомобилей
                                        (isAutoCategory && (!autoProperties.brand || !autoProperties.model || !autoProperties.year))
                                    }
                                >
                                    {t('listings.create.submit')}
                                </Button>
                            </Grid>
                        </Grid>
                    </form>
                </Paper>
            </Box>
        </Container>
    );
};

export default CreateListing;