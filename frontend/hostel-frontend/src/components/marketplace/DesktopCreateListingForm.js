// /data/proj/hostel-booking-system/frontend/hostel-frontend/src/components/marketplace/DesktopCreateListingForm.js

import React from 'react';
import { useTranslation } from 'react-i18next'; // Добавляем импорт
import {
    TextField,
    Button,
    Grid,
    FormControlLabel,
    Switch,
    IconButton,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Box,
} from "@mui/material";
import { Delete as DeleteIcon } from '@mui/icons-material';
import LocationPicker from '../global/LocationPicker';
import MiniMap from '../maps/MiniMap';
import ImageUploader from './ImageUploader';
import HierarchicalSelect from './HierarchicalSelect';

const DesktopCreateListingForm = ({ 
    listing, 
    setListing, 
    categories, 
    images, 
    setImages, 
    previewUrls, 
    setPreviewUrls,
    onSubmit
}) => {
    const { t, i18n } = useTranslation('marketplace'); // Добавляем i18n

    const handleSubmitForm = (e) => {
        e.preventDefault();
        // Добавляем язык интерфейса как оригинальный язык объявления
        onSubmit({
            ...listing,
            original_language: i18n.language
        });
    };

    return (
        <form onSubmit={handleSubmitForm}>
            <Grid container spacing={3}>
                <Grid item xs={12}>
                    <FormControl fullWidth required>
                        <InputLabel>{t('listings.create.category')}</InputLabel>
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
                        label={t('listings.create.name')}
                        fullWidth
                        required
                        value={listing.title}
                        onChange={(e) => setListing({ ...listing, title: e.target.value })}
                    />
                </Grid>


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
                        gridTemplateColumns: 'repeat(auto-fill, minmax(100px, 1fr))',
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

                <Grid item xs={6}>
                    <TextField
                        label="Цена"
                        type="number"
                        fullWidth
                        required
                        value={listing.price}
                        onChange={(e) => setListing({ ...listing, price: e.target.value })}
                    />
                </Grid>

                <Grid item xs={6}>
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
                    <Box sx={{ mb: 1 }}>
                        <LocationPicker onLocationSelect={(location) => {
                            setListing(prev => ({
                                ...prev,
                                latitude: location.latitude,
                                longitude: location.longitude,
                                location: location.formatted_address,
                                city: location.address_components?.city || '',
                                country: location.address_components?.country || ''
                            }));
                        }} />
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
                        label="Показывать местоположение на карте"
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
                        size="large"
                        disabled={!listing.title || !listing.description || !listing.category_id || listing.price <= 0}
                    >
                        {t('listings.create.submit')}
                    </Button>
                </Grid>
            </Grid>
        </form>
    );
};

export default DesktopCreateListingForm;