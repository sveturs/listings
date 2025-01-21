// frontend/hostel-frontend/src/pages/marketplace/CreateListingPage.js
import React, { useState, useEffect } from "react";
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
    useMediaQuery,
    Modal
} from "@mui/material";
import { Delete as DeleteIcon, CloudUpload as CloudUploadIcon } from '@mui/icons-material';
import LocationPicker from '../../components/global/LocationPicker';
import MiniMap from '../../components/maps/MiniMap';
import { GoogleMap, Marker } from '@react-google-maps/api';
import axios from "../../api/axios";
import { useLanguage } from '../../contexts/LanguageContext';
import ImageUploader from '../../components/marketplace/ImageUploader';




const CreateListing = () => {
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
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
    useEffect(() => {
        const fetchCategories = async () => {
            try {
                const response = await axios.get("/api/v1/marketplace/categories");
                setCategories(response.data.data || []);
            } catch (err) {
                setError("Ошибка при загрузке категорий");
            }
        };
        fetchCategories();
    }, []);

    const handleImageChange = (e) => {
        console.log('handleImageChange triggered');
        const files = Array.from(e.target.files || []);
        console.log('Selected files:', files);

        if (files.length === 0) {
            console.log('No files selected');
            return;
        }

        const validFiles = files.filter(file => {
            console.log('Checking file:', file.name, 'type:', file.type, 'size:', file.size);

            if (!file.type.startsWith('image/')) {
                console.log('Invalid file type:', file.type);
                setError("Можно загружать только изображения");
                return false;
            }
            if (file.size > 15 * 1024 * 1024) {
                console.log('File too large:', file.size);
                setError("Размер файла не должен превышать 15MB");
                return false;
            }
            return true;
        });

        console.log('Valid files:', validFiles);

        if (validFiles.length === 0) return;

        validFiles.forEach(file => {
            const reader = new FileReader();
            reader.onloadend = () => {
                console.log('File read successfully:', file.name);
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
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        console.log('Form submission started');
        setError("");
        setSuccess(false);

        try {
            // Создаем объект с данными листинга
            const listingData = {
                ...listing,
                price: parseFloat(listing.price),
                original_language: language
            };

            console.log('Submitting form data:', listingData);
            const response = await axios.post("/api/v1/marketplace/listings", listingData);
            const listingId = response.data.data.id;
            console.log('Listing created:', listingId);

            if (images.length > 0) {
                console.log('Preparing to upload images:', images.length);
                const formData = new FormData();
                images.forEach((image, index) => {
                    console.log('Adding image to formData:', image.name);
                    formData.append('images', image);
                    if (index === 0) {
                        formData.append('main_image_index', '0');
                    }
                });

                console.log('Uploading images...');
                await axios.post(`/api/v1/marketplace/listings/${listingId}/images`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
                console.log('Images uploaded successfully');
            }

            setSuccess(true);

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
                longitude: null
            });
            setImages([]);
            setPreviewUrls([]);

        } catch (error) {
            console.error('Error details:', error);
            setError(error.response?.data?.error || "Ошибка при создании объявления");
        }
    };

    return (
<Container 
    maxWidth="md" 
    disableGutters={isMobile}
    sx={{ 
        mx: isMobile ? 0 : 'auto',  // Убираем автоматические отступы на мобильном
        width: isMobile ? '100%' : 'auto'  // Принудительно растягиваем на всю ширину
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
                        Объявление успешно создано!
                    </Alert>
                )}
    
    <Paper sx={{ 
    p: isMobile ? '8px 0' : 3,  // Оставляем только вертикальный отступ 8px
    boxShadow: isMobile ? 'none' : 1,
    bgcolor: isMobile ? 'transparent' : 'background.paper',
    width: '100%'
}}>
                    <form onSubmit={handleSubmit}>
                        <Grid container spacing={isMobile ? 2 : 3}>
                            {/* Базовая информация */}
                            <Grid item xs={12}>
                            <FormControl fullWidth required size={isMobile ? "small" : "medium"}>
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
                            <Grid item xs={12}>
                                <TextField
                                    label="Наименование товара"
                                    fullWidth
                                    required
                                    value={listing.title}
                                    onChange={(e) => setListing({ ...listing, title: e.target.value })}
                                    size={isMobile ? "small" : "medium"}
                                />
                            </Grid>
    
                            <Grid item xs={12}>
                                <TextField
                                    label="Описание"
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
                                label="Цена"
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



                        {/* Местоположение */}
                        <Grid item xs={12}>
                            <Box sx={{ mb: 1 }}>
                                <LocationPicker onLocationSelect={handleLocationSelect} />
                            </Box>

                            {listing.latitude && listing.longitude && (
                                <MiniMap
                                    latitude={listing.latitude}
                                    longitude={listing.longitude}
                                    address={listing.location}
                                    onExpand={() => setShowExpandedMap(true)}
                                />
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

                        {/* Фотографии */}
                        <Grid item xs={12}>
                            <ImageUploader
                                onImagesSelected={(processedImages) => {
                                    setImages(processedImages.map(img => img.file));
                                    setPreviewUrls(processedImages.map(img => img.preview));
                                }}
                            />

                            <Box sx={{ 
                                mt: 2, 
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
                                            alt={`Preview ${index}`}
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
                            <Button
                                type="submit"
                                variant="contained"
                                color="primary"
                                fullWidth
                                size={isMobile ? "large" : "large"}
                                disabled={!listing.title || !listing.description || !listing.category_id || listing.price <= 0}
                            >
                                Создать объявление
                            </Button>
                        </Grid>
                    </Grid>
                </form>
            </Paper>
        </Box>

        {/* Модальное окно с картой */}
        {showExpandedMap && (
            <Modal
                open={showExpandedMap}
                onClose={() => setShowExpandedMap(false)}
                sx={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    p: 1
                }}
            >
                {/* ... оставляем существующий код модального окна ... */}
            </Modal>
        )}
    </Container>
);

};

export default CreateListing;