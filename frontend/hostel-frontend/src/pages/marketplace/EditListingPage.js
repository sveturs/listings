import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
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
import { Delete as DeleteIcon, CloudUpload as CloudUploadIcon } from '@mui/icons-material';
import LocationPicker from '../../components/global/LocationPicker';
import MiniMap from '../../components/maps/MiniMap';
import { GoogleMap, Marker } from '@react-google-maps/api';
import axios from '../../api/axios';
import { useAuth } from '../../contexts/AuthContext';
import { useLanguage } from '../../contexts/LanguageContext';
import ImageUploader from '../../components/marketplace/ImageUploader';


const EditListingPage = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const { user } = useAuth();
    const { language } = useLanguage();

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
        longitude: null
    });

    const [images, setImages] = useState([]);
    const [previewUrls, setPreviewUrls] = useState([]);
    const [categories, setCategories] = useState([]);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState(false);
    const [showExpandedMap, setShowExpandedMap] = useState(false);
    const [loading, setLoading] = useState(true);

    // Загрузка данных объявления и категорий
    useEffect(() => {
        const fetchData = async () => {
            try {
                const [listingResponse, categoriesResponse] = await Promise.all([
                    axios.get(`/api/v1/marketplace/listings/${id}`),
                    axios.get("/api/v1/marketplace/categories")
                ]);

                const listingData = listingResponse.data.data;

                // Проверяем, является ли текущий пользователь владельцем
                if (listingData.user_id !== user?.id) {
                    navigate('/marketplace');
                    return;
                }

                setListing({
                    title: listingData.title,
                    description: listingData.description,
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

                // Загружаем существующие изображения
                if (listingData.images) {
                    setPreviewUrls(listingData.images.map(img =>
                        `${process.env.REACT_APP_BACKEND_URL}/uploads/${img.file_path}`
                    ));
                }

                setCategories(categoriesResponse.data.data || []);
                setLoading(false);
            } catch (err) {
                setError("Ошибка при загрузке данных");
                setLoading(false);
            }
        };

        if (user?.id) {
            fetchData();
        }
    }, [id, user, navigate]);

    const handleLocationSelect = (location) => {
        setListing(prev => ({
            ...prev,
            latitude: location.latitude,
            longitude: location.longitude,
            location: location.formatted_address,
            city: location.address_components?.city || '',
            country: location.address_components?.country || ''
        }));
    };

    const handleImageChange = (e) => {
        const files = Array.from(e.target.files || []);
        if (files.length === 0) return;

        const validFiles = files.filter(file => {
            if (!file.type.startsWith('image/')) {
                setError("Можно загружать только изображения");
                return false;
            }
            if (file.size > 15 * 1024 * 1024) {
                setError("Размер файла не должен превышать 15MB");
                return false;
            }
            return true;
        });

        if (validFiles.length === 0) return;

        setImages(prev => [...prev, ...validFiles]);

        validFiles.forEach(file => {
            const reader = new FileReader();
            reader.onloadend = () => {
                setPreviewUrls(prev => [...prev, reader.result]);
            };
            reader.readAsDataURL(file);
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError("");
        setSuccess(false);

        try {
            // Обновляем основные данные объявления
            await axios.put(`/api/v1/marketplace/listings/${id}`, {
                ...listing,
                price: parseFloat(listing.price),
                original_language: language
            });

            // Если есть новые изображения, загружаем их
            if (images.length > 0) {
                const formData = new FormData();
                images.forEach(image => {
                    formData.append('images', image);
                });

                await axios.post(`/api/v1/marketplace/listings/${id}/images`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
            }

            setSuccess(true);
            setTimeout(() => {
                navigate(`/marketplace/listings/${id}`);
            }, 1500);

        } catch (error) {
            console.error('Ошибка при обновлении объявления:', error);
            setError(error.response?.data?.error || "Ошибка при обновлении объявления");
        }
    };

    if (loading) {
        return (
            <Container maxWidth="md">
                <Box sx={{ mt: 4, textAlign: 'center' }}>
                    <Typography>Загрузка...</Typography>
                </Box>
            </Container>
        );
    }

    return (
        <Container maxWidth="md">
            <Box sx={{ mt: 4, mb: 4 }}>
                <Typography variant="h4" gutterBottom>
                    Редактирование объявления
                </Typography>

                {error && (
                    <Alert severity="error" sx={{ mb: 2 }}>
                        {error}
                    </Alert>
                )}

                {success && (
                    <Alert severity="success" sx={{ mb: 2 }}>
                        Объявление успешно обновлено!
                    </Alert>
                )}

                <Paper sx={{ p: 3 }}>
                    <form onSubmit={handleSubmit}>
                        <Grid container spacing={3}>
                            <Grid item xs={12}>
                                <TextField
                                    label="Заголовок"
                                    fullWidth
                                    required
                                    value={listing.title}
                                    onChange={(e) => setListing({ ...listing, title: e.target.value })}
                                />
                            </Grid>

                            <Grid item xs={12}>
                                <TextField
                                    label="Описание"
                                    fullWidth
                                    required
                                    multiline
                                    rows={4}
                                    value={listing.description}
                                    onChange={(e) => setListing({ ...listing, description: e.target.value })}
                                />
                            </Grid>

                            <Grid item xs={12} sm={6}>
                                <TextField
                                    label="Цена"
                                    type="number"
                                    fullWidth
                                    required
                                    value={listing.price}
                                    onChange={(e) => setListing({ ...listing, price: e.target.value })}
                                />
                            </Grid>

                            <Grid item xs={12} sm={6}>
                                <FormControl fullWidth required>
                                    <InputLabel>Категория</InputLabel>
                                    <Select
                                        value={listing.category_id}
                                        onChange={(e) => setListing({ ...listing, category_id: e.target.value })}
                                    >
                                        {categories.map((category) => (
                                            <MenuItem key={category.id} value={category.id}>
                                                {category.name}
                                            </MenuItem>
                                        ))}
                                    </Select>
                                </FormControl>
                            </Grid>

                            <Grid item xs={12} sm={6}>
                                <FormControl fullWidth required>
                                    <InputLabel>Состояние</InputLabel>
                                    <Select
                                        value={listing.condition}
                                        onChange={(e) => setListing({ ...listing, condition: e.target.value })}
                                    >
                                        <MenuItem value="new">Новое</MenuItem>
                                        <MenuItem value="used">Б/у</MenuItem>
                                    </Select>
                                </FormControl>
                            </Grid>

                            <Grid item xs={12}>
                                <Typography variant="h6" gutterBottom>
                                    Местоположение
                                </Typography>
                                <LocationPicker
                                    onLocationSelect={handleLocationSelect}
                                    initialLocation={{
                                        latitude: listing.latitude,
                                        longitude: listing.longitude,
                                        formatted_address: listing.location
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
                                    label="Показывать местоположение на карте"
                                    sx={{ mt: 1 }}
                                />
                            </Grid>

                            <Grid item xs={12}>
                                <Typography variant="h6" gutterBottom>
                                    Фотографии
                                </Typography>
                                <ImageUploader
                                    onImagesSelected={(processedImages) => {
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
                                                alt={`Preview ${index}`}
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
                                                    URL.revokeObjectURL(url); // Освобождаем память
                                                }}
                                            >
                                                <DeleteIcon />
                                            </IconButton>
                                        </Box>
                                    ))}
                                </Box>
                            </Grid>

                            <Grid item xs={12}>
                                <Box sx={{ display: 'flex', gap: 2 }}>
                                    <Button
                                        type="submit"
                                        variant="contained"
                                        color="primary"
                                        fullWidth
                                        size="large"
                                        disabled={!listing.title || !listing.description || !listing.category_id || listing.price <= 0}
                                    >
                                        Сохранить изменения
                                    </Button>
                                    <Button
                                        variant="outlined"
                                        fullWidth
                                        size="large"
                                        onClick={() => navigate(`/marketplace/listings/${id}`)}
                                    >
                                        Отмена
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
                    sx={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        p: 2
                    }}
                >
                    <Paper
                        sx={{
                            position: 'relative',
                            width: '100%',
                            maxWidth: 1200,
                            maxHeight: '90vh',
                            overflow: 'hidden'
                        }}
                    >
                        <GoogleMap
                            mapContainerStyle={{
                                width: '100%',
                                height: '80vh'
                            }}
                            center={{
                                lat: listing.latitude,
                                lng: listing.longitude
                            }}
                            zoom={15}
                            options={{
                                zoomControl: true,
                                mapTypeControl: true,
                                streetViewControl: true,
                                gestureHandling: "greedy"
                            }}
                        >
                            <Marker
                                position={{
                                    lat: listing.latitude,
                                    lng: listing.longitude
                                }}
                                title={listing.title}
                            />
                        </GoogleMap>
                    </Paper>
                </Modal>
            )}
        </Container>
    );
};

export default EditListingPage;