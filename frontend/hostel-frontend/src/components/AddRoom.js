import React, { useState } from "react";
import LocationPicker from './LocationPicker';
import {
    TextField,
    Button,
    Container,
    Typography,
    Box,
    Alert,
    Grid,
    IconButton,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    FormControlLabel,
    Switch
} from "@mui/material";
import { Delete as DeleteIcon, CloudUpload as CloudUploadIcon } from '@mui/icons-material';
import axios from "../api/axios";

const AddRoom = () => {
    const [room, setRoom] = useState({
        name: "",
        accommodation_type: "room",
        capacity: 0,
        price_per_night: 0,
        address_street: "",
        address_city: "",
        address_state: "",
        address_country: "",
        address_postal_code: "",
        is_shared: false,
        total_beds: null,
        available_beds: null,
        has_private_bathroom: true,
        latitude: null,
        longitude: null,
        formatted_address: ''
    });

    const [beds, setBeds] = useState([
        { bed_number: "1", price_per_night: 0 }
    ]);
    const [images, setImages] = useState([]);
    const [previewUrls, setPreviewUrls] = useState([]);
    const [errorMessage, setErrorMessage] = useState("");
    const [isSuccess, setIsSuccess] = useState(false);

    const handleLocationSelect = (location) => {
        setRoom(prev => ({
            ...prev,
            latitude: location.latitude,
            longitude: location.longitude,
            formatted_address: location.formatted_address,
            address_street: location.address_components.find(
                c => c.types.includes('route')
            )?.long_name || '',
            address_city: location.address_components.find(
                c => c.types.includes('locality')
            )?.long_name || '',
            address_state: location.address_components.find(
                c => c.types.includes('administrative_area_level_1')
            )?.long_name || '',
            address_country: location.address_components.find(
                c => c.types.includes('country')
            )?.long_name || '',
            address_postal_code: location.address_components.find(
                c => c.types.includes('postal_code')
            )?.long_name || ''
        }));
    };

    const handleImageChange = (e) => {
        const files = Array.from(e.target.files || []);
        if (files.length === 0) return;

        const validFiles = files.filter(file => {
            if (!file.type.startsWith('image/')) {
                setErrorMessage("Можно загружать только изображения");
                return false;
            }
            if (file.size > 5 * 1024 * 1024) {
                setErrorMessage("Размер файла не должен превышать 5MB");
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
            reader.onerror = () => {
                setErrorMessage("Ошибка при чтении файла: " + file.name);
            };
            reader.readAsDataURL(file);
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setErrorMessage("");
        setIsSuccess(false);

        try {
            if (room.accommodation_type === 'bed' && (!room.total_beds || beds.length === 0)) {
                setErrorMessage("Добавьте информацию о кроватях");
                return;
            }

            const roomData = {
                ...room,
                latitude: parseFloat(room.latitude),
                longitude: parseFloat(room.longitude)
            };

            // Создаем комнату
            const roomResponse = await axios.post("/rooms", roomData);
            const roomId = roomResponse.data.id;

            // Если тип размещения - койко-места, создаем кровати
            if (room.accommodation_type === 'bed' && beds.length > 0) {
                await Promise.all(
                    beds.map(async (bed) => {
                        try {
                            await axios.post(`/rooms/${roomId}/beds`, {
                                bed_number: bed.bed_number,
                                price_per_night: parseFloat(bed.price_per_night)
                            });
                        } catch (error) {
                            console.error('Ошибка добавления кровати:', error);
                            throw error;
                        }
                    })
                );
            }

            // Загружаем изображения
            if (images.length > 0) {
                const formData = new FormData();
                images.forEach(image => {
                    formData.append('images', image);
                });

                await axios.post(`/rooms/${roomId}/images`, formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
            }

            setIsSuccess(true);
            setRoom({
                name: "",
                accommodation_type: "room",
                capacity: 0,
                price_per_night: 0,
                address_street: "",
                address_city: "",
                address_state: "",
                address_country: "",
                address_postal_code: "",
                is_shared: false,
                total_beds: null,
                available_beds: null,
                has_private_bathroom: true,
                latitude: null,
                longitude: null,
                formatted_address: ''
            });
            setBeds([{ bed_number: "1", price_per_night: 0 }]);
            setImages([]);
            setPreviewUrls([]);
        } catch (error) {
            console.error('Ошибка при добавлении:', error);
            setErrorMessage(error.response?.data || "Ошибка при добавлении объекта размещения");
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>
                Добавить объект размещения
            </Typography>
            {errorMessage && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {errorMessage}
                </Alert>
            )}
            {isSuccess && (
                <Alert severity="success" sx={{ mb: 2 }}>
                    Объект размещения успешно добавлен
                </Alert>
            )}
            <form onSubmit={handleSubmit}>
                <Grid container spacing={2}>
                    {/* Базовая информация */}
                    <Grid item xs={12}>
                        <TextField
                            label="Название"
                            fullWidth
                            required
                            value={room.name}
                            onChange={(e) => setRoom({ ...room, name: e.target.value })}
                        />
                    </Grid>

                    {/* Тип размещения */}
                    <Grid item xs={12}>
                        <FormControl fullWidth>
                            <InputLabel>Тип размещения</InputLabel>
                            <Select
                                value={room.accommodation_type}
                                onChange={(e) => {
                                    const newType = e.target.value;
                                    setRoom(prev => ({
                                        ...prev,
                                        accommodation_type: newType,
                                        total_beds: newType === 'bed' ? prev.total_beds : null,
                                        available_beds: newType === 'bed' ? prev.available_beds : null,
                                        is_shared: newType === 'bed' ? true : false,
                                    }));
                                }}
                            >
                                <MenuItem value="bed">Койко-место</MenuItem>
                                <MenuItem value="room">Комната</MenuItem>
                                <MenuItem value="apartment">Квартира</MenuItem>
                            </Select>
                        </FormControl>
                    </Grid>

                    {room.accommodation_type === 'bed' ? (
                        <>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    label="Всего кроватей"
                                    type="number"
                                    fullWidth
                                    required
                                    value={room.total_beds || ''}
                                    onChange={(e) => setRoom({
                                        ...room,
                                        total_beds: parseInt(e.target.value) || 0,
                                        available_beds: parseInt(e.target.value) || 0
                                    })}
                                />
                            </Grid>
                            <Grid item xs={12}>
                                {beds.map((bed, index) => (
                                    <Box key={index} sx={{ display: 'flex', gap: 2, mb: 2 }}>
                                        <TextField
                                            label={`Номер кровати ${index + 1}`}
                                            value={bed.bed_number}
                                            onChange={(e) => {
                                                const newBeds = [...beds];
                                                newBeds[index].bed_number = e.target.value;
                                                setBeds(newBeds);
                                            }}
                                        />
                                        <TextField
                                            label="Цена за ночь"
                                            type="number"
                                            value={bed.price_per_night}
                                            onChange={(e) => {
                                                const newBeds = [...beds];
                                                newBeds[index].price_per_night = parseFloat(e.target.value);
                                                setBeds(newBeds);
                                            }}
                                        />
                                        <IconButton onClick={() => {
                                            const newBeds = [...beds];
                                            newBeds.splice(index, 1);
                                            setBeds(newBeds);
                                        }}>
                                            <DeleteIcon />
                                        </IconButton>
                                    </Box>
                                ))}
                                <Button
                                    variant="outlined"
                                    onClick={() => setBeds([...beds, { bed_number: `${beds.length + 1}`, price_per_night: 0 }])}
                                >
                                    Добавить кровать
                                </Button>
                            </Grid>
                        </>
                    ) : (
                        <>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    label="Вместимость"
                                    type="number"
                                    fullWidth
                                    required
                                    value={room.capacity}
                                    onChange={(e) => setRoom({ ...room, capacity: parseInt(e.target.value) || 0 })}
                                />
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <TextField
                                    label="Цена за ночь"
                                    type="number"
                                    fullWidth
                                    required
                                    value={room.price_per_night}
                                    onChange={(e) => setRoom({ ...room, price_per_night: parseFloat(e.target.value) || 0 })}
                                />
                            </Grid>
                        </>
                    )}

                    <Grid item xs={12}>
                        <FormControlLabel
                            control={
                                <Switch
                                    checked={room.is_shared}
                                    onChange={(e) => setRoom({ ...room, is_shared: e.target.checked })}
                                />
                            }
                            label="Общее помещение"
                        />
                    </Grid>

                    <Grid item xs={12}>
                        <FormControlLabel
                            control={
                                <Switch
                                    checked={room.has_private_bathroom}
                                    onChange={(e) => setRoom({ ...room, has_private_bathroom: e.target.checked })}
                                />
                            }
                            label="Отдельная ванная комната"
                        />
                    </Grid>

                    {/* Выбор местоположения */}
                    <Grid item xs={12}>
                        <LocationPicker onLocationSelect={handleLocationSelect} />
                    </Grid>

                    {/* Адрес */}
                    <Grid item xs={12}>
                        <Typography variant="h6">Адрес</Typography>
                    </Grid>
                    <Grid item xs={12}>
                        <TextField
                            label="Улица"
                            fullWidth
                            required
                            value={room.address_street}
                            onChange={(e) => setRoom({ ...room, address_street: e.target.value })}
                        />
                    </Grid>
                    <Grid item xs={12} md={6}>
                        <TextField
                            label="Город"
                            fullWidth
                            required
                            value={room.address_city}
                            onChange={(e) => setRoom({ ...room, address_city: e.target.value })}
                        />
                    </Grid>
                    <Grid item xs={12} md={6}>
                        <TextField
                            label="Область/Регион"
                            fullWidth
                            value={room.address_state}
                            onChange={(e) => setRoom({ ...room, address_state: e.target.value })}
                        />
                    </Grid>
                    <Grid item xs={12} md={6}>
                        <TextField
                            label="Страна"
                            fullWidth
                            required
                            value={room.address_country}
                            onChange={(e) => setRoom({ ...room, address_country: e.target.value })}
                        />
                    </Grid>
                    <Grid item xs={12} md={6}>
                        <TextField
                            label="Почтовый индекс"
                            fullWidth
                            value={room.address_postal_code}
                            onChange={(e) => setRoom({ ...room, address_postal_code: e.target.value })}
                        />
                    </Grid>

                    {/* Изображения */}
                    <Grid item xs={12}>
                        <Typography variant="h6">Фотографии</Typography>
                        <Box sx={{ mt: 1, mb: 2 }}>
                            <Button
                                variant="contained"
                                component="label"
                                startIcon={<CloudUploadIcon />}
                            >
                                Загрузить изображения
                                <input
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
                                <Grid item xs={12} sm={4} key={index}>
                                    <Box sx={{ position: 'relative' }}>
                                        <img
                                            src={url}
                                            alt={`Preview ${index}`}
                                            style={{
                                                width: '100%',
                                                height: '200px',
                                                objectFit: 'cover',
                                                borderRadius: '4px'
                                            }}
                                        />
                                        <IconButton
                                            sx={{
                                                position: 'absolute',
                                                top: 8,
                                                right: 8,
                                                bgcolor: 'rgba(255, 255, 255, 0.8)'
                                            }}
                                            onClick={() => {
                                                setImages(prev => prev.filter((_, i) => i !== index));
                                                setPreviewUrls(prev => prev.filter((_, i) => i !== index));
                                            }}
                                        >
                                            <DeleteIcon />
                                        </IconButton>
                                    </Box>
                                </Grid>
                            ))}
                        </Grid>
                    </Grid>

                    <Grid item xs={12}>
                        <Button
                            type="submit"
                            variant="contained"
                            color="primary"
                            fullWidth
                            size="large"
                            disabled={!room.name ||
                                !room.latitude ||
                                !room.longitude ||
                                !room.address_city ||
                                !room.address_country ||
                                (room.accommodation_type === 'bed' && (!room.total_beds || beds.length === 0))}
                        >
                            Добавить
                        </Button>
                    </Grid>
                </Grid>
            </form>
        </Container>
    );
};

export default AddRoom;