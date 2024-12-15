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
    Switch,
    Paper,
    FormGroup,
    Checkbox
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
        {
            bed_number: "1",
            price_per_night: 0,
            has_outlet: true,
            has_light: true,
            has_shelf: true,
            bed_type: 'single',
            images: [],
            imagePreviewUrls: []
        }
    ]);
    
    const [images, setImages] = useState([]);
    const [previewUrls, setPreviewUrls] = useState([]);
    const [errorMessage, setErrorMessage] = useState("");
    const [isSuccess, setIsSuccess] = useState(false);

    const handleBedImageChange = (index, e) => {
        const files = Array.from(e.target.files || []);
        if (files.length === 0) return;

        const validFiles = files.filter(file => {
            if (!file.type.startsWith('image/')) {
                setErrorMessage("Можно загружать только изображения");
                return false;
            }
            if (file.size > 15 * 1024 * 1024) {
                setErrorMessage("Размер файла не должен превышать 5MB");
                return false;
            }
            return true;
        });

        if (validFiles.length === 0) return;

        const newBeds = [...beds];
        newBeds[index] = {
            ...newBeds[index],
            images: [...(newBeds[index].images || []), ...validFiles]
        };

        validFiles.forEach(file => {
            const reader = new FileReader();
            reader.onloadend = () => {
                newBeds[index].imagePreviewUrls = [
                    ...(newBeds[index].imagePreviewUrls || []),
                    reader.result
                ];
                setBeds([...newBeds]);
            };
            reader.readAsDataURL(file);
        });
    };

    const handleLocationSelect = (location) => {
        setRoom(prev => ({
            ...prev,
            latitude: location.latitude,
            longitude: location.longitude,
            formatted_address: location.formatted_address,
            address_street: location.address_components?.street || '',
            address_city: location.address_components?.city || '',
            address_state: location.address_components?.state || '',
            address_country: location.address_components?.country || '',
            address_postal_code: location.address_components?.postal_code || ''
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
            if (file.size > 15 * 1024 * 1024) {
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
            const roomResponse = await axios.post("/api/v1/rooms", room);
            const roomId = roomResponse.data.data.id;
            console.log('Room response:', roomResponse);

            if (!roomId) {
                console.error('No room ID in response:', roomResponse);
                throw new Error('Failed to get room ID');
            }

            if (room.accommodation_type === 'bed' && beds.length > 0) {
                for (const bed of beds) {
                    console.log(`Creating bed for room ${roomId}:`, bed);
                    const bedResponse = await axios.post(`/api/v1/rooms/${roomId}/beds`, {
                        bed_number: bed.bed_number,
                        price_per_night: parseFloat(bed.price_per_night),
                        has_outlet: bed.has_outlet,
                        has_light: bed.has_light,
                        has_shelf: bed.has_shelf,
                        bed_type: bed.bed_type
                    });
                    console.log('Bed response:', bedResponse);

                    const bedId = bedResponse.data.data.id;

                    if (bed.images && bed.images.length > 0) {
                        const formData = new FormData();
                        bed.images.forEach(image => {
                            formData.append('images', image);
                        });

                        await axios.post(`/api/v1/rooms/${roomId}/beds/${bedId}/images`, formData, {
                            headers: {
                                'Content-Type': 'multipart/form-data'
                            }
                        });
                    }
                }
            }

            if (images.length > 0) {
                const formData = new FormData();
                images.forEach(image => {
                    formData.append('images', image);
                });

                await axios.post(`/api/v1/rooms/${roomId}/images`, formData, {
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
            setBeds([{
                bed_number: "1",
                price_per_night: 0,
                has_outlet: true,
                has_light: true,
                has_shelf: true,
                bed_type: 'single',
                images: [],
                imagePreviewUrls: []
            }]);
            setImages([]);
            setPreviewUrls([]);

        } catch (error) {
            console.error('Ошибка при добавлении:', error);
            setErrorMessage(error.response?.data?.error || "Ошибка при добавлении объекта размещения");
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
                    <Grid item xs={12}>
                        <TextField
                            label="Название"
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
                            <Grid item container xs={12} spacing={2}>
                                <Grid item xs={12} md={6}>
                                    <TextField
                                        label="Всего кроватей"
                                        type="number"
                                        fullWidth
                                        required
                                        value={room.total_beds || ''}
                                        onChange={(e) => {
                                            const totalBeds = parseInt(e.target.value) || 0;
                                            const defaultPrice = room.default_price || 0;

                                            setRoom({
                                                ...room,
                                                total_beds: totalBeds,
                                                available_beds: totalBeds
                                            });

                                            const newBeds = Array.from({ length: totalBeds }, (_, index) => ({
                                                bed_number: `${index + 1}`,
                                                price_per_night: defaultPrice,
                                                has_outlet: true,
                                                has_light: true,
                                                has_shelf: true,
                                                bed_type: 'single',
                                                images: [],
                                                imagePreviewUrls: []
                                            }));
                                            setBeds(newBeds);
                                        }}
                                    />
                                </Grid>
                                <Grid item xs={12} md={6}>
                                    <TextField
                                        label="Стандартная цена за ночь"
                                        type="number"
                                        fullWidth
                                        required
                                        value={room.default_price || ''}
                                        onChange={(e) => {
                                            const defaultPrice = parseFloat(e.target.value) || 0;
                                            setRoom({
                                                ...room,
                                                default_price: defaultPrice
                                            });
                                            setBeds(prevBeds =>
                                                prevBeds.map(bed => ({
                                                    ...bed,
                                                    price_per_night: defaultPrice
                                                }))
                                            );
                                        }}
                                    />
                                </Grid>
                            </Grid>

                            <Grid item xs={12}>
                                <Typography variant="h6" sx={{ mt: 2, mb: 2 }}>
                                    Информация о кроватях
                                </Typography>
                                {beds.map((bed, index) => (
                                    <Grid item xs={12} key={index}>
                                        <Paper sx={{ p: 2, mb: 2 }}>
                                            <Grid container spacing={2} alignItems="center">
                                                <Grid item xs={12} sm={6} md={3}>
                                                    <TextField
                                                        label={`Номер кровати ${index + 1}`}
                                                        value={bed.bed_number}
                                                        onChange={(e) => {
                                                            const newBeds = [...beds];
                                                            newBeds[index].bed_number = e.target.value;
                                                            setBeds(newBeds);
                                                        }}
                                                        fullWidth
                                                    />
                                                </Grid>
                                                
                                                <Grid item xs={12} sm={6} md={3}>
                                                    <TextField
                                                        label="Цена за ночь"
                                                        type="number"
                                                        value={bed.price_per_night}
                                                        onChange={(e) => {
                                                            const newBeds = [...beds];
                                                            newBeds[index].price_per_night = parseFloat(e.target.value);
                                                            setBeds(newBeds);
                                                        }}
                                                        fullWidth
                                                    />
                                                </Grid>

                                                <Grid item xs={12} sm={6} md={3}>
                                                    <FormControl fullWidth>
                                                        <InputLabel>Тип кровати</InputLabel>
                                                        <Select
                                                            value={bed.bed_type || 'single'}
                                                            onChange={(e) => {
                                                                const newBeds = [...beds];
                                                                newBeds[index].bed_type = e.target.value;
                                                                setBeds(newBeds);
                                                            }}
                                                        >
                                                            <MenuItem value="single">Отдельностоящая</MenuItem>
                                                            <MenuItem value="top">Верхний ярус</MenuItem>
                                                            <MenuItem value="bottom">Нижний ярус</MenuItem>
                                                        </Select>
                                                    </FormControl>
                                                </Grid>

                                                <Grid item xs={12}>
                                                    <FormGroup row>
                                                        <FormControlLabel
                                                            control={
                                                                <Checkbox
                                                                    checked={bed.has_outlet}
                                                                    onChange={(e) => {
                                                                        const newBeds = [...beds];
                                                                        newBeds[index].has_outlet = e.target.checked;
                                                                        setBeds(newBeds);
                                                                    }}
                                                                />
                                                            }
                                                            label="Розетка"
                                                        />
                                                        <FormControlLabel
                                                            control={
                                                                <Checkbox
                                                                    checked={bed.has_light}
                                                                    onChange={(e) => {
                                                                        const newBeds = [...beds];
                                                                        newBeds[index].has_light = e.target.checked;
                                                                        setBeds(newBeds);
                                                                    }}
                                                                />
                                                            }
                                                            label="Светильник"
                                                            />
                                                            <FormControlLabel
                                                                control={
                                                                    <Checkbox
                                                                        checked={bed.has_shelf}
                                                                        onChange={(e) => {
                                                                            const newBeds = [...beds];
                                                                            newBeds[index].has_shelf = e.target.checked;
                                                                            setBeds(newBeds);
                                                                        }}
                                                                    />
                                                                }
                                                                label="Полка"
                                                            />
                                                        </FormGroup>
                                                    </Grid>
    
                                                    <Grid item xs={12}>
                                                        <Box>
                                                            <Button
                                                                variant="contained"
                                                                component="label"
                                                                size="small"
                                                                startIcon={<CloudUploadIcon />}
                                                            >
                                                                Фото кровати
                                                                <input
                                                                    type="file"
                                                                    hidden
                                                                    multiple
                                                                    accept="image/*"
                                                                    onChange={(e) => handleBedImageChange(index, e)}
                                                                />
                                                            </Button>
                                                            {bed.imagePreviewUrls && bed.imagePreviewUrls.length > 0 && (
                                                                <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', mt: 2 }}>
                                                                    {bed.imagePreviewUrls.map((url, imgIndex) => (
                                                                        <Box
                                                                            key={imgIndex}
                                                                            sx={{
                                                                                position: 'relative',
                                                                                width: 60,
                                                                                height: 60
                                                                            }}
                                                                        >
                                                                            <img
                                                                                src={url}
                                                                                alt={`Preview ${imgIndex}`}
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
                                                                                    top: -8,
                                                                                    right: -8,
                                                                                    bgcolor: 'background.paper'
                                                                                }}
                                                                                onClick={() => {
                                                                                    const newBeds = [...beds];
                                                                                    newBeds[index].images.splice(imgIndex, 1);
                                                                                    newBeds[index].imagePreviewUrls.splice(imgIndex, 1);
                                                                                    setBeds(newBeds);
                                                                                }}
                                                                            >
                                                                                <DeleteIcon fontSize="small" />
                                                                            </IconButton>
                                                                        </Box>
                                                                    ))}
                                                                </Box>
                                                            )}
                                                        </Box>
                                                    </Grid>
                                                </Grid>
                                            </Paper>
                                        </Grid>
                                    ))}
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
    
                        <Grid item xs={12}>
                            <LocationPicker onLocationSelect={handleLocationSelect} />
                        </Grid>
    
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