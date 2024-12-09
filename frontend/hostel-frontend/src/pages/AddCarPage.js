import React, { useState } from "react";
import {
    Container,
    TextField,
    Button,
    Box,
    Alert,
    Grid,
    Typography,
    Paper,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Chip,
    Input,
    Divider,
} from "@mui/material";
import {
    DirectionsCar as CarIcon,
    CloudUpload as UploadIcon,
    AddLocation as LocationIcon,
} from '@mui/icons-material';
import axios from '../api/axios';
import LocationPicker from '../components/LocationPicker';

const FEATURES = [
    'Кондиционер',
    'Климат-контроль',
    'Круиз-контроль',
    'Парктроники',
    'Камера заднего вида',
    'Навигация',
    'Bluetooth',
    'USB',
    'AUX',
    'MP3',
    'CD',
    'Кожаный салон',
    'Люк',
    'Панорамная крыша',
    'Подогрев сидений',
    'Электропривод сидений',
    'Электропривод зеркал',
    'Электропривод окон',
];

const AddCarPage = () => {
    const [formData, setFormData] = useState({
        make: "",
        model: "",
        year: "",
        price_per_day: "",
        transmission: "",
        fuel_type: "",
        seats: "",
        location: "",
        latitude: null,
        longitude: null,
        formatted_address: "",
        features: [],
        description: "",
    });

    const [images, setImages] = useState([]);
    const [previewUrls, setPreviewUrls] = useState([]);
    const [successMessage, setSuccessMessage] = useState("");
    const [errorMessage, setErrorMessage] = useState("");

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleImageChange = (e) => {
        const files = Array.from(e.target.files);
        
        // Проверяем только размер
        const validFiles = files.filter(file => {
            if (file.size > 10 * 1024 * 1024) { // Увеличим до 10MB
                setErrorMessage("Размер файла не должен превышать 10MB");
                return false;
            }
            return true;
        });
    
        if (validFiles.length) {
            setImages(prev => [...prev, ...validFiles]);
            // Создаем превью
            validFiles.forEach(file => {
                const reader = new FileReader();
                reader.onloadend = () => {
                    setPreviewUrls(prev => [...prev, reader.result]);
                };
                reader.readAsDataURL(file);
            });
        }
    };

    
    const handleLocationSelect = (location) => {
        setFormData(prev => ({
            ...prev,
            location: location.formatted_address,
            latitude: location.latitude,
            longitude: location.longitude,
            formatted_address: location.formatted_address
        }));
    };

    const handleFeatureToggle = (feature) => {
        setFormData(prev => ({
            ...prev,
            features: prev.features.includes(feature)
                ? prev.features.filter(f => f !== feature)
                : [...prev.features, feature]
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setErrorMessage("");
        setSuccessMessage("");

        try {
            // Сначала создаем автомобиль
            const carResponse = await axios.post("/api/v1/cars", {
                ...formData,
                year: parseInt(formData.year),
                price_per_day: parseFloat(formData.price_per_day),
                seats: parseInt(formData.seats),
                availability: true
            });

            const carId = carResponse.data.data.id;

            // Если есть изображения, загружаем их
            if (images.length > 0) {
                const formData = new FormData();
                images.forEach(image => {
                    formData.append('images', image);
                });

                await axios.post(`/api/v1/cars/${carId}/images`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
            }

            setSuccessMessage("Автомобиль успешно добавлен!");
            setFormData({
                make: "",
                model: "",
                year: "",
                price_per_day: "",
                transmission: "",
                fuel_type: "",
                seats: "",
                location: "",
                latitude: null,
                longitude: null,
                formatted_address: "",
                features: [],
                description: "",
            });
            setImages([]);
            setPreviewUrls([]);

        } catch (error) {
            setErrorMessage(error.response?.data?.error || "Ошибка при добавлении автомобиля");
        }
    };

    return (
        <Container>
            <Box sx={{ mt: 4, mb: 4 }}>
                <Typography variant="h4" component="h1" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
                    <CarIcon sx={{ mr: 2 }} />
                    Добавить автомобиль
                </Typography>

                {successMessage && (
                    <Alert severity="success" sx={{ mb: 2 }}>{successMessage}</Alert>
                )}
                {errorMessage && (
                    <Alert severity="error" sx={{ mb: 2 }}>{errorMessage}</Alert>
                )}

                <form onSubmit={handleSubmit}>
                    <Paper sx={{ p: 3, mb: 3 }}>
                        <Typography variant="h6" gutterBottom>
                            Основная информация
                        </Typography>
                        <Grid container spacing={2}>
                            <Grid item xs={12} sm={6} md={4}>
                                <TextField
                                    label="Марка автомобиля"
                                    name="make"
                                    value={formData.make}
                                    onChange={handleChange}
                                    fullWidth
                                    required
                                />
                            </Grid>
                            <Grid item xs={12} sm={6} md={4}>
                                <TextField
                                    label="Модель"
                                    name="model"
                                    value={formData.model}
                                    onChange={handleChange}
                                    fullWidth
                                    required
                                />
                            </Grid>
                            <Grid item xs={12} sm={6} md={4}>
                                <TextField
                                    label="Год выпуска"
                                    name="year"
                                    type="number"
                                    value={formData.year}
                                    onChange={handleChange}
                                    fullWidth
                                    required
                                    InputProps={{ inputProps: { min: 1900, max: new Date().getFullYear() } }}
                                />
                            </Grid>
                            <Grid item xs={12} sm={6} md={4}>
                                <TextField
                                    label="Цена за день"
                                    name="price_per_day"
                                    type="number"
                                    value={formData.price_per_day}
                                    onChange={handleChange}
                                    fullWidth
                                    required
                                    InputProps={{
                                        inputProps: { min: 0 },
                                        startAdornment: <span>₽&nbsp;</span>
                                    }}
                                />
                            </Grid>
                            <Grid item xs={12} sm={6} md={4}>
                                <FormControl fullWidth required>
                                    <InputLabel>Коробка передач</InputLabel>
                                    <Select
                                        name="transmission"
                                        value={formData.transmission}
                                        onChange={handleChange}
                                        label="Коробка передач"
                                    >
                                        <MenuItem value="automatic">Автомат</MenuItem>
                                        <MenuItem value="manual">Механика</MenuItem>
                                    </Select>
                                </FormControl>
                            </Grid>
                            <Grid item xs={12} sm={6} md={4}>
                                <FormControl fullWidth required>
                                    <InputLabel>Тип топлива</InputLabel>
                                    <Select
                                        name="fuel_type"
                                        value={formData.fuel_type}
                                        onChange={handleChange}
                                        label="Тип топлива"
                                    >
                                        <MenuItem value="petrol">Бензин</MenuItem>
                                        <MenuItem value="diesel">Дизель</MenuItem>
                                        <MenuItem value="electric">Электро</MenuItem>
                                        <MenuItem value="hybrid">Гибрид</MenuItem>
                                    </Select>
                                </FormControl>
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    label="Количество мест"
                                    name="seats"
                                    type="number"
                                    value={formData.seats}
                                    onChange={handleChange}
                                    fullWidth
                                    required
                                    InputProps={{ inputProps: { min: 1, max: 9 } }}
                                />
                            </Grid>
                        </Grid>
                    </Paper>

                    <Paper sx={{ p: 3, mb: 3 }}>
                        <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
                            <LocationIcon sx={{ mr: 1 }} />
                            Местоположение
                        </Typography>
                        <LocationPicker onLocationSelect={handleLocationSelect} />
                    </Paper>

                    <Paper sx={{ p: 3, mb: 3 }}>
                        <Typography variant="h6" gutterBottom>
                            Особенности и комплектация
                        </Typography>
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 2 }}>
                            {FEATURES.map(feature => (
                                <Chip
                                    key={feature}
                                    label={feature}
                                    onClick={() => handleFeatureToggle(feature)}
                                    color={formData.features.includes(feature) ? "primary" : "default"}
                                    sx={{ '&:hover': { opacity: 0.8 } }}
                                />
                            ))}
                        </Box>
                        <TextField
                            label="Описание автомобиля"
                            name="description"
                            value={formData.description}
                            onChange={handleChange}
                            fullWidth
                            multiline
                            rows={4}
                            sx={{ mt: 2 }}
                        />
                    </Paper>

                    <Paper sx={{ p: 3, mb: 3 }}>
                        <Typography variant="h6" gutterBottom>
                            Фотографии
                        </Typography>
                        <Box sx={{ mb: 2 }}>
                            <Button
                                variant="outlined"
                                component="label"
                                startIcon={<UploadIcon />}
                            >
                                Загрузить фото
                                <Input
                                    type="file"
                                    hidden
                                    multiple
                                    accept="image/*"
                                    onChange={handleImageChange}
                                />
                            </Button>
                        </Box>
                        <Grid container spacing={2}>
                            {previewUrls.map((url, index) => (
                                <Grid item xs={6} sm={4} md={3} key={index}>
                                    <Box
                                        sx={{
                                            width: '100%',
                                            paddingTop: '75%',
                                            position: 'relative',
                                            borderRadius: 1,
                                            overflow: 'hidden'
                                        }}
                                    >
                                        <img
                                            src={url}
                                            alt={`Preview ${index + 1}`}
                                            style={{
                                                position: 'absolute',
                                                top: 0,
                                                left: 0,
                                                width: '100%',
                                                height: '100%',
                                                objectFit: 'cover'
                                            }}
                                        />
                                    </Box>
                                </Grid>
                            ))}
                        </Grid>
                    </Paper>

                    <Divider sx={{ my: 3 }} />

                    <Button
                        type="submit"
                        variant="contained"
                        color="primary"
                        size="large"
                        fullWidth
                        disabled={!formData.make || !formData.model || !formData.price_per_day || !formData.location}
                        startIcon={<CarIcon />}
                    >
                        Добавить автомобиль
                    </Button>
                </form>
            </Box>
        </Container>
    );
};

export default AddCarPage;