// frontend/hostel-frontend/src/pages/marketplace/CreateListingPage.tsx
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
    useMediaQuery,
    Tooltip,
    CircularProgress,
    SelectChangeEvent
} from "@mui/material";
import { useNavigate } from 'react-router-dom';
import { Delete as DeleteIcon, CloudUpload as CloudUploadIcon } from '@mui/icons-material';
import LocationPicker from '../../components/global/LocationPicker';
import MiniMap from '../../components/maps/MiniMap';
import axios from "../../api/axios";
import { useLanguage } from '../../contexts/LanguageContext';
import ImageUploader from '../../components/marketplace/ImageUploader';
import CategorySelect from '../../components/marketplace/CategorySelect';
import { ChevronRight, ChevronLeft, Info } from 'lucide-react';
import { useLocation, LocationData } from '../../contexts/LocationContext';
import AttributeFields, { AttributeValue } from '../../components/marketplace/AttributeFields';

// Define interfaces for component
interface Category {
    id: string | number;
    name: string;
    description?: string;
    parent_id?: string | number;
    [key: string]: any;
}

interface ListingFormData {
    title: string;
    description: string;
    price: number;
    category_id: string;
    condition: 'new' | 'used';
    location: string;
    city: string;
    country: string;
    show_on_map: boolean;
    latitude: number | null;
    longitude: number | null;
    original_language: string;
    translations?: {
        [key: string]: {
            [key: string]: string;
        }
    };
    [key: string]: any;
}

interface GoogleTranslateInfo {
    used: number;
    limit: number;
}

interface ApiErrorResponse {
    response?: {
        data?: {
            error?: string;
        }
    }
}

type TranslationProvider = 'google' | 'openai';

const CreateListing: React.FC = () => {
    const { t, i18n } = useTranslation('marketplace');
    const theme = useTheme();
    const { language } = useLanguage();
    const { userLocation } = useLocation();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
    const navigate = useNavigate();

    const [listing, setListing] = useState<ListingFormData>({
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
        longitude: null,
        original_language: i18n.language // Set current interface language as original language
    });

    const [images, setImages] = useState<File[]>([]);
    const [previewUrls, setPreviewUrls] = useState<string[]>([]);
    const [categories, setCategories] = useState<Category[]>([]);
    const [error, setError] = useState<string>("");
    const [success, setSuccess] = useState<boolean>(false);
    const [showExpandedMap, setShowExpandedMap] = useState<boolean>(false);
    const [locationWarning, setLocationWarning] = useState<boolean>(false);
    const [attributeValues, setAttributeValues] = useState<AttributeValue[]>([]);
    const [googleTranslateInfo, setGoogleTranslateInfo] = useState<GoogleTranslateInfo>({ used: 0, limit: 100 });
    const [loadingLimits, setLoadingLimits] = useState<boolean>(false);
    const [translationProvider, setTranslationProvider] = useState<TranslationProvider>(
        (localStorage.getItem('preferredTranslationProvider') as TranslationProvider) || 'google'
    );

    const getTranslatedText = (field: string): string => {
        if (!listing) return '';

        // If current language matches original language
        if (language === listing.original_language) {
            return listing[field];
        }

        // Try to get translation
        if (listing.translations &&
            listing.translations[language] &&
            listing.translations[language][field]) {
            return listing.translations[language][field];
        }

        // If no translation, return original
        return listing[field];
    };

    useEffect(() => {
        // If we have user location data, use it
        if (userLocation && userLocation.lat && userLocation.lon) {
            // Create location object for LocationPicker
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

            // Set location in component state
            setListing(prev => ({
                ...prev,
                latitude: userLocation.lat,
                longitude: userLocation.lon,
                location: initialLocation.formatted_address,
                city: userLocation.city || '',
                country: userLocation.country || 'Serbia'
            }));

            console.log("Set initial location from context:", initialLocation);
        }
    }, [userLocation]);

    useEffect(() => {
        const fetchCategories = async (): Promise<void> => {
            try {
                const response = await axios.get("/api/v1/marketplace/categories");
                setCategories(response.data.data || []);
            } catch (err) {
                setError(t('listings.create.error'));
            }
        };
        fetchCategories();
        
        // Get Google Translate limits information
        setLoadingLimits(true);
        axios.get('/api/v1/translation/limits')
            .then(response => {
                const { google } = response.data.data;
                setGoogleTranslateInfo({
                    used: google.used,
                    limit: google.limit
                });
            })
            .catch(error => {
                console.error('Error getting translation limits:', error);
                // Use default values in case of error
                setGoogleTranslateInfo({ used: 0, limit: 100 });
            })
            .finally(() => {
                setLoadingLimits(false);
            });
    }, [t]);

    const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
        const fileList = e.target.files || [];
        if (fileList.length === 0) return;

        const files = Array.from(fileList) as File[];
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
                setPreviewUrls(prev => [...prev, reader.result as string]);
            };
            reader.onerror = (error) => {
                console.error('Error reading file:', error);
            };
            reader.readAsDataURL(file);
        });

        setImages(prev => [...prev, ...validFiles]);
    };

    interface LocationSelectData {
        latitude: number;
        longitude: number;
        formatted_address: string;
        address_components?: {
            city?: string;
            country?: string;
            [key: string]: any;
        };
    }

    const handleLocationSelect = (location: LocationSelectData): void => {
        setListing(prev => ({
            ...prev,
            latitude: location.latitude,
            longitude: location.longitude,
            location: location.formatted_address,
            city: location.address_components?.city || '',
            country: location.address_components?.country || ''
        }));

        // Check for coordinates and show warning if they're missing
        setLocationWarning(!location.latitude || !location.longitude);
    };

    // Convert attributes array for server submission
    const prepareAttributesForSubmission = (attributes: AttributeValue[]): AttributeValue[] => {
        return attributes.map(attr => {
            // Create a copy of the attribute for modification
            const newAttr = { ...attr };

            // Ensure correct format for numeric values
            if (attr.attribute_type === 'number' && attr.value !== undefined) {
                // Convert to number and ensure it's stored in the correct field
                let numValue: number;

                // Handle fractional values for engine capacity
                if (attr.attribute_name === 'engine_capacity') {
                    // For engine capacity, use parseFloat to preserve decimal places
                    numValue = parseFloat(attr.value as string);
                    // Round to one decimal place
                    if (!isNaN(numValue)) {
                        numValue = Math.round(numValue * 10) / 10;
                    }
                } else if (attr.attribute_name === 'year') {
                    // For manufacturing year, handle separately to avoid resetting
                    numValue = parseInt(attr.value as string);
                    // Log for debugging
                    console.log(`Preparing year for submission: ${attr.value} -> ${numValue}`);

                    // Check for valid year
                    if (isNaN(numValue) || numValue < 1900 || numValue > new Date().getFullYear() + 1) {
                        // If year is invalid, set current year
                        const currentYear = new Date().getFullYear();
                        console.warn(`Invalid year ${numValue}, setting current year ${currentYear}`);
                        numValue = currentYear;
                    }
                } else {
                    // For other numeric fields, still use parseFloat
                    numValue = parseFloat(attr.value as string);
                }

                if (!isNaN(numValue)) {
                    // Ensure numeric_value is always set for numeric attributes
                    newAttr.numeric_value = numValue;
                    // Update main value for consistency
                    newAttr.value = numValue;

                    // Also update display value
                    if (attr.attribute_name === 'year') {
                        newAttr.display_value = String(numValue);
                    }
                } else {
                    console.error(`Failed to convert "${attr.value}" to number for attribute ${attr.attribute_name}`);
                }
            }

            return newAttr;
        });
    };
 
    const handleSubmit = async (e: React.FormEvent): Promise<void> => {
        e.preventDefault();
        setError("");
        setSuccess(false);

        try {
            // Get user profile for city, if not specified
            let userCity = listing.city;
            let userCountry = listing.country;

            // If coordinates or city not specified, try to get from profile
            if ((!listing.latitude || !listing.longitude || !listing.city) && !locationWarning) {
                try {
                    const profileResponse = await axios.get('/api/v1/users/profile');
                    const userProfile = profileResponse.data.data;

                    // Use city from profile if available and not already specified
                    if (!listing.city && userProfile.city) {
                        userCity = userProfile.city;
                    }

                    // Use country from profile if available and not already specified
                    if (!listing.country && userProfile.country) {
                        userCountry = userProfile.country;
                    }

                    setLocationWarning(true);
                } catch (profileError) {
                    console.error("Failed to fetch user profile:", profileError);
                }
            }

            // Process attributes before submission
            const processedAttributes = prepareAttributesForSubmission(attributeValues);

            // Use translation provider from component state
            console.log(`Using translation provider: ${translationProvider}`);
            
            // Log data for debugging
            console.log("Submitting listing:", {
                ...listing,
                price: parseFloat(listing.price.toString()),
                original_language: i18n.language,
                attributes: processedAttributes,
                city: userCity || listing.city,
                country: userCountry || listing.country
            });

            const response = await axios.post(`/api/v1/marketplace/listings?translation_provider=${translationProvider}`, {
                ...listing,
                price: parseFloat(listing.price.toString()),
                original_language: i18n.language,
                attributes: processedAttributes,
                city: userCity || listing.city,
                country: userCountry || listing.country
            });

            const listingId = response.data.data.id;

            // Check attribute saving
            console.log("Listing created with ID:", listingId);
            console.log("Checking attribute submission:", processedAttributes);


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
                original_language: 'sr'
            });
            setImages([]);
            setPreviewUrls([]);

            setTimeout(() => {
                navigate(`/marketplace/listings/${listingId}`);
            }, 1500);

        } catch (error) {
            const apiError = error as ApiErrorResponse;
            setError(apiError.response?.data?.error || t('listings.create.error'));
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
                                    {/* @ts-ignore - Working around MUI typing issues */}
                                    <InputLabel shrink>{t('listings.create.category')}</InputLabel>
                                    <CategorySelect
                                        categories={categories}
                                        value={listing.category_id}
                                        onChange={(value: string) => setListing({ ...listing, category_id: value })}
                                        error={!listing.category_id}
                                    />
                                </FormControl>
                            </Grid>

                            <Grid item xs={12}>
                                <Alert severity="info" sx={{ mb: 2 }}>
                                  {t('listings.create.languageInfo', {
                                    language: i18n.language === 'ru' ? 'русском' : 
                                              i18n.language === 'en' ? 'английском' : 'сербском',
                                    defaultValue: 'Вы создаете объявление на {{language}} языке. Система автоматически переведет его на другие поддерживаемые языки.'
                                  })}
                                </Alert>
                                
                                {/* Translation provider selection */}
                                <FormControl
                                    fullWidth
                                    sx={{ mb: 2 }}
                                    size={isMobile ? "small" : "medium"}
                                >
                                    {/* @ts-ignore - Working around MUI typing issues */}
                                    <InputLabel>{t('store.listings.selectTranslationProvider', 'Сервис перевода')}</InputLabel>
                                    <Select
                                        value={translationProvider}
                                        onChange={(e: SelectChangeEvent<TranslationProvider>) => {
                                            setTranslationProvider(e.target.value as TranslationProvider);
                                            localStorage.setItem('preferredTranslationProvider', e.target.value as string);
                                        }}
                                    >
                                        <MenuItem value="google">
                                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                                <Typography>{t('store.listings.googleTranslate', 'Google Translate')}</Typography>
                                                <Typography variant="caption" color="success.main">
                                                    ({t('store.listings.freeWithinLimits', 'Бесплатно в пределах лимита')})
                                                </Typography>
                                            </Box>
                                        </MenuItem>
                                        <MenuItem value="openai">
                                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                                <Typography>{t('store.listings.openaiTranslate', 'OpenAI (GPT)')}</Typography>
                                                <Typography variant="caption" color="text.secondary">
                                                    ({t('store.listings.higherQuality', 'Высокое качество')})
                                                </Typography>
                                            </Box>
                                        </MenuItem>
                                    </Select>
                                </FormControl>
                                
                                {/* Google Translate limits information */}
                                {translationProvider === 'google' && (
                                    <Box sx={{ mb: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
                                        <Info size={16} color={theme.palette.info.main} />
                                        <Typography variant="caption" color="text.secondary">
                                            {loadingLimits ? (
                                                <CircularProgress size={12} />
                                            ) : (
                                                t('store.listings.googleTranslateLimits', {
                                                    used: googleTranslateInfo.used,
                                                    limit: googleTranslateInfo.limit,
                                                    defaultValue: 'Использовано {{used}} из {{limit}} в этом месяце.'
                                                })
                                            )}
                                        </Typography>
                                    </Box>
                                )}
                                
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
                                    onChange={(e) => setListing({ ...listing, price: parseFloat(e.target.value) })}
                                    size={isMobile ? "small" : "medium"}
                                />
                            </Grid>

                            <Grid item xs={6}>
                                <FormControl fullWidth required size={isMobile ? "small" : "medium"}>
                                    {/* @ts-ignore - Working around MUI typing issues */}
                                    <InputLabel>{t('listings.create.condition.label')}</InputLabel>
                                    <Select<'new' | 'used'>
                                        value={listing.condition}
                                        onChange={(e) => setListing({ ...listing, condition: e.target.value as 'new' | 'used' })}
                                    >
                                        <MenuItem value="new">{t('listings.create.condition.new')}</MenuItem>
                                        <MenuItem value="used">{t('listings.create.condition.used')}</MenuItem>
                                    </Select>
                                </FormControl>
                            </Grid>

                            <Grid item xs={12}>
                                <Box sx={{ mb: 1 }}>
                                    <LocationPicker
                                        onLocationSelect={handleLocationSelect}
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
                                <AttributeFields
                                    categoryId={listing.category_id}
                                    value={attributeValues}
                                    onChange={setAttributeValues}
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
                                    disabled={!listing.title || !listing.description || !listing.category_id || listing.price <= 0}
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