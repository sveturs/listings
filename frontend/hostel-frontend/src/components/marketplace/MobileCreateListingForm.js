// /data/proj/hostel-booking-system/frontend/hostel-frontend/src/components/marketplace/MobileCreateListingForm.js

import React from 'react';
import { useTranslation } from 'react-i18next'; // Добавляем импорт
import { Box, TextField, Button, FormControl, InputLabel, Select, MenuItem, FormControlLabel, Switch, Alert } from '@mui/material';
import LocationPicker from '../global/LocationPicker';
import MiniMap from '../maps/MiniMap';
import ImageUploader from './ImageUploader';

const MobileCreateListingForm = ({
    listing,
    setListing,
    categories,
    images,
    setImages,
    previewUrls,
    setPreviewUrls,
    onSubmit,
    error,
    success
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
        <Box sx={{ px: 2, pb: 2 }}>
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

            <form onSubmit={handleSubmitForm} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                <TextField
                    label={t('listings.create.name')}
                    fullWidth
                    required
                    value={listing.title}
                    onChange={(e) => setListing({ ...listing, title: e.target.value })}
                />



                <TextField
                    label="Описание"
                    fullWidth
                    required
                    multiline
                    rows={4}
                    value={listing.description}
                    onChange={(e) => setListing({ ...listing, description: e.target.value })}
                />

                <Box sx={{ display: 'flex', gap: 2 }}>
                    <TextField
                        label="Цена"
                        type="number"
                        fullWidth
                        required
                        value={listing.price}
                        onChange={(e) => setListing({ ...listing, price: e.target.value })}
                    />

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
                </Box>

                <FormControl fullWidth required>
                    <InputLabel>{t('listings.create.category')}</InputLabel>
                    <Box sx={{ mt: 2 }}>
                        <VirtualizedCategoryTree
                            selectedId={listing.category_id}
                            onSelectCategory={(id) => {
                                setListing(prev => ({
                                    ...prev,
                                    category_id: id
                                }));
                            }}
                        />
                    </Box>
                </FormControl>

                <Box sx={{ mt: 1 }}>
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
                    />
                </Box>

                {listing.latitude && listing.longitude && (
                    <MiniMap
                        latitude={listing.latitude}
                        longitude={listing.longitude}
                        address={listing.location}
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
                />

                <ImageUploader
                    onImagesSelected={(processedImages) => {
                        setImages(processedImages.map(img => img.file));
                        setPreviewUrls(processedImages.map(img => img.preview));
                    }}
                />

                <Box sx={{
                    display: 'grid',
                    gridTemplateColumns: 'repeat(auto-fill, minmax(80px, 1fr))',
                    gap: 1,
                    mt: 1
                }}>
                    {previewUrls.map((url, index) => (
                        <Box
                            key={index}
                            sx={{
                                position: 'relative',
                                paddingTop: '100%',
                                borderRadius: 1,
                                overflow: 'hidden'
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
                                    objectFit: 'cover'
                                }}
                            />
                        </Box>
                    ))}
                </Box>

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
            </form>
        </Box>
    );
};

export default MobileCreateListingForm;