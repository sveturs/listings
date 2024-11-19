import React, { useState } from "react";
import { Delete as DeleteIcon, CloudUpload as CloudUploadIcon } from '@mui/icons-material';
import axios from "../api/axios";
import { 
    TextField, Button, Container, Typography, 
    Box, Alert, Grid, IconButton,
    FormControl, InputLabel, Select, MenuItem,
    FormControlLabel, Switch 
} from "@mui/material";

const AddRoom = () => {
    const [room, setRoom] = useState({
        name: "",
        accommodation_type: "room", // Добавляем тип размещения
        capacity: 0,
        price_per_night: 0,
        address_street: "",
        address_city: "",
        address_state: "",
        address_country: "",
        address_postal_code: "",
        is_shared: false,         // Общее помещение
        total_beds: null,         // Всего кроватей (для типа bed)
        available_beds: null,     // Доступно кроватей (для типа bed)
        has_private_bathroom: true // Наличие отдельной ванной
    });
    const [beds, setBeds] = useState([
        { bed_number: "1", price_per_night: 0 }
    ]);
    const [images, setImages] = useState([]);
    const [previewUrls, setPreviewUrls] = useState([]);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState(false);

    const handleImageChange = (e) => {
        console.log('handleImageChange called', e.target.files); // Отладочный вывод
        const files = Array.from(e.target.files || []);
        if (files.length === 0) {
            console.log('No files selected'); // Отладочный вывод
            return;
        }

        const validFiles = files.filter(file => {
            if (!file.type.startsWith('image/')) {
                setError("Можно загружать только изображения");
                console.log('Invalid file type:', file.type); // Отладочный вывод
                return false;
            }
            if (file.size > 5 * 1024 * 1024) {
                setError("Размер файла не должен превышать 5MB");
                console.log('File too large:', file.size); // Отладочный вывод
                return false;
            }
            return true;
        });
        if (validFiles.length === 0) {
            console.log('No valid files after filtering'); // Отладочный вывод
            return;
        }
        setImages(prev => [...prev, ...validFiles]);
        
        // Создаем превью
        validFiles.forEach(file => {
            const reader = new FileReader();
            reader.onloadend = () => {
                setPreviewUrls(prev => [...prev, reader.result]);
            };
            reader.onerror = () => {
                console.error('Error reading file:', file.name); // Отладочный вывод
                setError("Ошибка при чтении файла: " + file.name);
            };
            reader.readAsDataURL(file);
        });
    };

    const removeImage = (index) => {
        setImages(prev => prev.filter((_, i) => i !== index));
        setPreviewUrls(prev => prev.filter((_, i) => i !== index));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError("");
        setSuccess(false);
    
        try {
            // Валидация в зависимости от типа размещения
            if (room.accommodation_type === 'bed') {
                if (!room.total_beds || beds.length === 0) {
                    setError("Добавьте информацию о кроватях");
                    return;
                }
            }
    
            // Подготовка данных комнаты
            const roomData = {
                ...room,
                // Если это не койко-место, очищаем связанные поля
                total_beds: room.accommodation_type === 'bed' ? room.total_beds : null,
                available_beds: room.accommodation_type === 'bed' ? room.total_beds : null,
                // Цена за ночь будет общей для комнаты/квартиры, для койко-мест - отдельно
                price_per_night: room.accommodation_type === 'bed' ? 0 : room.price_per_night
            };
    
            // Создаем комнату
            console.log('Отправляемые данные комнаты:', roomData);
            const roomResponse = await axios.post("/rooms", roomData);
            const roomId = roomResponse.data.id;
    
            // Если тип размещения - койко-места, создаем кровати
            if (room.accommodation_type === 'bed' && beds.length > 0) {
                console.log('Добавление кроватей для комнаты:', roomId);
                await Promise.all(
                    beds.map(async (bed) => {
                        try {
                            const bedData = {
                                bed_number: bed.bed_number,
                                price_per_night: parseFloat(bed.price_per_night)
                            };
                            console.log('Отправка данных кровати:', bedData);
                            await axios.post(`/rooms/${roomId}/beds`, bedData);
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
    
            setSuccess(true);
            // Очистка формы
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
                has_private_bathroom: true
            });
            setBeds([{ bed_number: "1", price_per_night: 0 }]);
            setImages([]);
            setPreviewUrls([]);
        } catch (error) {
            console.error('Ошибка при добавлении:', error);
            setError(error.response?.data || "Ошибка при добавлении объекта размещения");
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>
                Добавить комнату
            </Typography>
            {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                </Alert>
            )}
            {success && (
                <Alert severity="success" sx={{ mb: 2 }}>
                    Комната успешно добавлена!
                </Alert>
            )}
            <form onSubmit={handleSubmit}>
                <Grid container spacing={2}>
                    <Grid item xs={12}>
                        <TextField
                            label="Название комнаты"
                            fullWidth
                            required
                            value={room.name}
                            onChange={(e) => setRoom({ ...room, name: e.target.value })}
                        />
                    </Grid>

                    <Grid item xs={12}>
    <FormControl fullWidth>
        <InputLabel>Тип размещения</InputLabel>
        <Select
            value={room.accommodation_type}
            onChange={(e) => setRoom({ ...room, accommodation_type: e.target.value })}
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
                inputProps={{ min: 1 }}
            />
        </Grid>
        <Grid item xs={12}>
            <Typography variant="h6" sx={{ mt: 2, mb: 1 }}>
                Кровати
            </Typography>
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
    <Grid item xs={12} md={6}>
        <TextField
            label="Вместимость"
            type="number"
            fullWidth
            required
            value={room.capacity}
            onChange={(e) => setRoom({ ...room, capacity: parseInt(e.target.value) || 0 })}
            inputProps={{ min: 1 }}
        />
    </Grid>
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

                    <Grid item xs={12} md={6}>
                        <TextField
                            label="Вместимость"
                            type="number"
                            fullWidth
                            required
                            value={room.capacity}
                            onChange={(e) => setRoom({ ...room, capacity: parseInt(e.target.value) || 0 })}
                            inputProps={{ min: 1 }}
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
                            inputProps={{ min: 0, step: "0.01" }}
                        />
                    </Grid>

                    {/* Адрес */}
                    <Grid item xs={12}>
                        <Typography variant="h6" sx={{ mt: 2, mb: 1 }}>
                            Адрес
                        </Typography>
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
                        <Typography variant="h6" sx={{ mt: 2, mb: 1 }}>
                            Фотографии комнаты
                        </Typography>
                    </Grid>
            <Grid item xs={12}>
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
                            onClick={(e) => e.target.value = null} // Сброс значения input
                        />
                    </Button>
                </Box>
                {/* Добавим информацию о загруженных файлах */}
                {images.length > 0 && (
                    <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                        Выбрано файлов: {images.length}
                    </Typography>
                )}
            </Grid>
                </Grid>

                <Grid container spacing={2} sx={{ mb: 2 }}>
                    {previewUrls.map((url, index) => (
                        <Grid item xs={12} sm={6} md={4} key={index}>
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
                                        bgcolor: 'rgba(255, 255, 255, 0.8)',
                                        '&:hover': {
                                            bgcolor: 'rgba(255, 255, 255, 0.9)',
                                        }
                                    }}
                                    onClick={() => removeImage(index)}
                                >
                                    <DeleteIcon />
                                </IconButton>
                            </Box>
                        </Grid>
                    ))}
                </Grid>

                <Box sx={{ mt: 3 }}>
                    <Button 
                        type="submit" 
                        variant="contained" 
                        color="primary"
                        fullWidth
                        size="large"
                        disabled={!room.name || !room.capacity || !room.price_per_night || 
                                !room.address_street || !room.address_city || !room.address_country}
                    >
                        Добавить комнату
                    </Button>
                </Box>
            </form>
        </Container>
    );
};

export default AddRoom;