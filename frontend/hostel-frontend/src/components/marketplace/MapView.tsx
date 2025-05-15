// frontend/hostel-frontend/src/components/marketplace/MapView.tsx
import React, { useState, useEffect, useRef, useMemo, RefObject } from 'react';
import { useTranslation } from 'react-i18next';
import { Navigation, X, List, Maximize2, Store } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { Listing } from '../../types/listing';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import {
  Box,
  Paper,
  Typography,
  Chip,
  Button,
  Card,
  CardContent,
  CardMedia,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  useTheme,
  useMediaQuery,
  Collapse,
  IconButton,
  Modal,
  SelectChangeEvent
} from '@mui/material';
import { TILE_LAYER_URL, TILE_LAYER_ATTRIBUTION } from '../maps/map-constants';
import '../maps/leaflet-icons'; // For fixing Leaflet icons
import FullscreenMap from '../maps/FullscreenMap';
import { useLocation } from '../../contexts/LocationContext';
import CentralAttributeFilters from './CentralAttributeFilters';

// TypeScript interfaces

interface ImageObject {
  id?: number | string;
  file_path?: string;
  is_main?: boolean;
  storage_type?: string;
  public_url?: string;
  [key: string]: any;
}

interface MapFilters {
  category_id?: string | number;
  latitude?: number;
  longitude?: number;
  distance?: string;
  attributeFilters?: Record<string, any>;
  [key: string]: any;
}

interface ListingPreviewProps {
  listing: Listing;
  onClose: () => void;
  onNavigate: (listingId: string | number) => void;
}

interface MapViewProps {
  listings: Listing[];
  filters: MapFilters;
  userLocation?: {
    latitude: number | null;
    longitude: number | null;
    city?: string;
    country?: string;
  };
  onFilterChange: (filters: MapFilters) => void;
  onMapClose: () => void;
}

interface MapCenter {
  latitude: number;
  longitude: number;
  title: string;
}

interface MapMarker {
  latitude: number;
  longitude: number;
  title: string;
  tooltip: string;
  id: number | string;
  listing: Listing;
}

interface StorefrontData {
  listings: Listing[];
  name: string;
  location: [number, number];
  address?: string;
}

interface StorefrontLocations {
  [key: string]: {
    listings: Listing[];
    storefront_id: number | string;
    name: string;
    location: [number, number];
    address?: string;
  };
}

// Listing preview component when clicking on a marker
const ListingPreview: React.FC<ListingPreviewProps> = ({ listing, onClose, onNavigate }) => {
  const { t } = useTranslation(['marketplace', 'common']);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const navigate = useNavigate();

  if (!listing) return null;

  const formatPrice = (price: number): string => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RSD',
      maximumFractionDigits: 0
    }).format(price);
  };

  // Optimized image path normalization
  const getImageUrl = (images?: (string | ImageObject)[]): string => {
    if (!images || !Array.isArray(images) || images.length === 0) {
      return '/placeholder-listing.jpg';
    }

    // Find main image or use the first one in the list
    let mainImage = images.find(img => img && (img as ImageObject).is_main === true) || images[0];

    // Use environment variable from window.ENV instead of process.env
    const baseUrl = (window as any).ENV?.REACT_APP_MINIO_URL || (window as any).ENV?.REACT_APP_BACKEND_URL || '';
    console.log('MapView: Using baseUrl from env:', baseUrl);

    // 1. String paths (for backward compatibility)
    if (typeof mainImage === 'string') {
      console.log('MapView: Processing string image path:', mainImage);

      // MinIO relative path
      if (mainImage.startsWith('/listings/')) {
        const url = `${baseUrl}${mainImage}`;
        console.log('MapView: Using MinIO relative path:', url);
        return url;
      }

      // ID/filename.jpg (direct MinIO path)
      if (mainImage.match(/^\d+\/[^\/]+$/)) {
        const url = `${baseUrl}/listings/${mainImage}`;
        console.log('MapView: Using direct MinIO path pattern:', url);
        return url;
      }

      // Local storage (backward compatibility)
      const url = `${baseUrl}/uploads/${mainImage}`;
      console.log('MapView: Using local storage path:', url);
      return url;
    }

    // 2. File info objects
    if (typeof mainImage === 'object' && mainImage !== null) {
      const imageObj = mainImage as ImageObject;
      console.log('MapView: Processing image object:', imageObj);

      // Priority 1: Use PublicURL if available
      if (imageObj.public_url && typeof imageObj.public_url === 'string' && imageObj.public_url.trim() !== '') {
        const publicUrl = imageObj.public_url;
        console.log('MapView: Found public_url string:', publicUrl);

        // Absolute URL
        if (publicUrl.startsWith('http')) {
          console.log('MapView: Using absolute URL:', publicUrl);
          return publicUrl;
        }
        // Relative URL with /listings/
        else if (publicUrl.startsWith('/listings/')) {
          const url = `${baseUrl}${publicUrl}`;
          console.log('MapView: Using public_url with listings path:', url);
          return url;
        }
        // Other relative URL
        else {
          const url = `${baseUrl}${publicUrl}`;
          console.log('MapView: Using general relative public_url:', url);
          return url;
        }
      }

      // Priority 2: Form URL based on storage type and file path
      if (imageObj.file_path) {
        if (imageObj.storage_type === 'minio' || imageObj.file_path.includes('listings/')) {
          // Consider the possibility of listings/ prefix in the path
          const filePath = imageObj.file_path.includes('listings/')
            ? imageObj.file_path.replace('listings/', '') 
            : imageObj.file_path;

          const url = `${baseUrl}/listings/${filePath}`;
          console.log('MapView: Constructed MinIO URL from path:', url);
          return url;
        }

        // Local storage
        const url = `${baseUrl}/uploads/${imageObj.file_path}`;
        console.log('MapView: Using local storage path from object:', url);
        return url;
      }
    }

    console.log('MapView: Could not determine image URL, using placeholder');
    return '/placeholder-listing.jpg';
  };

  const imageUrl = getImageUrl(listing.images);

  // Determine if we're showing a storefront card or listing
  const isStorefrontCard = Boolean(listing.isPartOfStorefront ||
    (listing.storefront_id && listing.isUniqueLocation));

  console.log('ListingPreview: Determining card type:', {
    isStorefrontCard: isStorefrontCard,
    isPartOfStorefront: listing.isPartOfStorefront,
    hasStorefrontId: Boolean(listing.storefront_id),
    isUnique: listing.isUniqueLocation
  });

  return (
    <Card
      sx={{
        position: 'absolute',
        bottom: isMobile ? 0 : 16,
        left: isMobile ? 0 : 16,
        maxWidth: isMobile ? '100%' : 400,
        width: isMobile ? '100%' : 'auto',
        zIndex: 1000,
        borderRadius: isMobile ? '8px 8px 0 0' : 1,
      }}
    >
      <Box sx={{ position: 'relative' }}>
        <IconButton
          onClick={onClose}
          sx={{
            position: 'absolute',
            top: 8,
            right: 8,
            bgcolor: 'background.paper',
            opacity: 0.8,
            '&:hover': { bgcolor: 'background.paper', opacity: 1 },
            zIndex: 10
          }}
        >
          <X size={16} />
        </IconButton>

        {isStorefrontCard && (
          <Box
            sx={{
              position: 'absolute',
              top: 0,
              left: 0,
              right: 0,
              backgroundColor: theme.palette.primary?.main,
              color: 'white',
              padding: '4px 12px',
              display: 'flex',
              alignItems: 'center',
              zIndex: 5
            }}
          >
            <Store size={16} style={{ marginRight: 6 }} />
            <Typography variant="subtitle2" fontWeight="bold">
              {listing.storefrontName || t('common:map.storefront')}
              {listing.isUniqueLocation && (
                <Chip
                  size="small"
                  label={t('common:map.uniqueLocation', { defaultValue: 'Unique location' })}
                  sx={{ ml: 1, height: 18, fontSize: '0.65rem' }}
                  color="secondary"
                />
              )}
            </Typography>
          </Box>
        )}

        {imageUrl && (
          <CardMedia
            component="img"
            height={isStorefrontCard ? 160 : 140}
            image={imageUrl}
            alt={listing.title}
            sx={{
              pt: isStorefrontCard ? '24px' : 0
            }}
          />
        )}

        <CardContent>
          {isStorefrontCard && (
            <Box sx={{ mb: 2 }}>
              <Typography variant="body2" color="text.secondary">
                {t('common:map.items', { count: listing.storefrontItemCount })}
              </Typography>
              <Typography variant="h6" fontWeight="bold" sx={{ mt: 0.5 }}>
                {listing.title}
              </Typography>
            </Box>
          )}

          {!isStorefrontCard && (
            <Typography variant="subtitle1" noWrap gutterBottom>
              {listing.title}
            </Typography>
          )}

          <Typography variant="h6" color="primary" gutterBottom>
            {formatPrice(listing.price)}
          </Typography>

          <Box display="flex" justifyContent="space-between" alignItems="center" mt={1}>
            <Chip
              label={listing.condition === 'new' ? t('listings.conditions.new') : t('listings.conditions.used')}
              size="small"
              color={listing.condition === 'new' ? 'success' : 'default'}
            />

            {isStorefrontCard ? (
              <Button
                variant="contained"
                color="primary"
                size="small"
                fullWidth
                sx={{ ml: 1 }}
                onClick={() => {
                  console.log(`Navigating to storefront page: /shop/${listing.storefront_id}?highlightedListingId=${listing.id}`, {
                    storefront_id: listing.storefront_id,
                    listing_id: listing.id
                  });
                  navigate(`/shop/${listing.storefront_id}?highlightedListingId=${listing.id}`);
                }}
              >
                {listing.isUniqueLocation
                  ? t('common:map.viewItemInStorefront', { defaultValue: 'View in storefront' })
                  : t('common:map.viewStorefront')}
              </Button>
            ) : (
              <Button
                variant="contained"
                size="small"
                onClick={() => {
                  console.log(`Navigating to listing page: ${listing.id}`);
                  onNavigate(listing.id);
                }}
              >
                {t('listings.details.viewDetails')}
              </Button>
            )}
          </Box>
        </CardContent>
      </Box>
    </Card>
  );
};

const MapView: React.FC<MapViewProps> = ({ listings, filters, onFilterChange, onMapClose }) => {
  const { t } = useTranslation('marketplace');
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const navigate = useNavigate();
  const { userLocation, detectUserLocation } = useLocation();

  // Create proper coordinates object from userLocation
  const locationCoordinates = userLocation ? {
    latitude: userLocation.lat,
    longitude: userLocation.lon
  } : null;

  // Log for debugging
  console.log("MapView userLocation:", userLocation);
  console.log("MapView locationCoordinates:", locationCoordinates);

  const mapRef = useRef<L.Map | null>(null);
  const markersRef = useRef<L.Marker[]>([]);
  const mapContainerRef = useRef<HTMLDivElement | null>(null);
  const [selectedListing, setSelectedListing] = useState<Listing | null>(null);
  const [mapReady, setMapReady] = useState<boolean>(false);
  
  useEffect(() => {
    if (selectedListing) {
      console.log("Selected listing state changed:", {
        id: selectedListing.id,
        isPartOfStorefront: selectedListing.isPartOfStorefront,
        storefront_id: selectedListing.storefront_id,
        isUniqueLocation: selectedListing.isUniqueLocation
      });
    }
  }, [selectedListing]);
  
  // State for fullscreen map modal
  const [expandedMapOpen, setExpandedMapOpen] = useState<boolean>(false);
  // Center for fullscreen map
  const [expandedMapCenter, setExpandedMapCenter] = useState<MapCenter | null>(null);
  // Markers for fullscreen map
  const [expandedMapMarkers, setExpandedMapMarkers] = useState<MapMarker[]>([]);

  // Ensure we always have userLocationState
  const [userLocationState, setUserLocationState] = useState<{latitude: number, longitude: number} | null>(locationCoordinates);
  
  const handleAttributeFilterChange = (newAttrFilters: Record<string, any>): void => {
    console.log("MapView: Received new attribute filters:", newAttrFilters);
    if (onFilterChange) {
      onFilterChange({
        ...filters,
        attributeFilters: newAttrFilters
      });
    }
  };

  // Handler for resetting attribute filters
  const resetAttributeFilters = (): void => {
    if (onFilterChange) {
      onFilterChange({
        ...filters,
        attributeFilters: {}
      });
    }
  };

  // Update userLocationState when userLocation changes
  useEffect(() => {
    if (userLocation) {
      setUserLocationState({
        latitude: userLocation.lat,
        longitude: userLocation.lon
      });
      console.log("Updated userLocationState from userLocation:", {
        latitude: userLocation.lat,
        longitude: userLocation.lon
      });
    }
  }, [userLocation]);

  // Function to convert distance string (e.g., "5km") to meters
  const getRadiusInMeters = (distanceString: string | undefined): number => {
    if (!distanceString) return 5000; // Default 5 km

    const match = distanceString.match(/^(\d+)km$/);
    if (match) {
      return parseInt(match[1]) * 1000;
    }
    return 5000;
  };

  // Initialize map
  useEffect(() => {
    if (!mapContainerRef.current || mapRef.current) return;

    // Add delay to ensure DOM is ready
    const timer = setTimeout(() => {
      try {
        // Set map center based on context data
        const initialPosition = userLocation
          ? [userLocation.lat, userLocation.lon]  // Use data from context
          : [45.2671, 19.8335]; // Novi Sad coordinates by default

        // Check that container exists and has dimensions
        if (!mapContainerRef.current ||
          mapContainerRef.current.clientWidth === 0 ||
          mapContainerRef.current.clientHeight === 0) {
          console.log("Map container not ready yet, retrying...");
          return;
        }

        // Create map with preferCanvas option for better performance
        mapRef.current = L.map(mapContainerRef.current, {
          preferCanvas: true,
          attributionControl: false,
          zoomControl: true,
          inertia: true,
          fadeAnimation: true,
          zoomAnimation: true,
          renderer: L.canvas()
        }).setView(initialPosition as [number, number], 13);

        // Add tile layer
        L.tileLayer(TILE_LAYER_URL, {
          attribution: TILE_LAYER_ATTRIBUTION,
          maxZoom: 19
        }).addTo(mapRef.current);

        // If user location is available, add marker and circle
        if (userLocation) {
          try {
            L.circle(initialPosition as [number, number], {
              color: theme.palette.primary?.main,
              fillColor: theme.palette.primary?.light,
              fillOpacity: 0.2,
              radius: getRadiusInMeters(filters.distance)
            }).addTo(mapRef.current);

            L.marker(initialPosition as [number, number], {
              icon: L.divIcon({
                html: `<div style="
                  background-color: ${theme.palette.primary?.main};
                  width: 16px;
                  height: 16px;
                  border-radius: 50%;
                  border: 2px solid white;
                  box-shadow: 0 0 4px rgba(0,0,0,0.3);
                "></div>`,
                className: 'my-location-marker',
                iconSize: [20, 20],
                iconAnchor: [10, 10]
              })
            }).addTo(mapRef.current)
              .bindTooltip(t('listings.map.yourLocation'), { permanent: false });
          } catch (innerError) {
            console.error("Error adding user marker:", innerError);
          }
        }

        // Handler for fixing _leaflet_pos issue
        mapRef.current.on('zoomanim', (e) => {
          // Do nothing, but this helps prevent errors
        });

        setMapReady(true);
      } catch (error) {
        console.error("Error initializing map:", error);
      }
    }, 300); // Increase delay to ensure DOM readiness

    return () => {
      clearTimeout(timer);
      if (mapRef.current) {
        try {
          // Explicitly remove all layers before removing map
          mapRef.current.eachLayer((layer) => {
            mapRef.current?.removeLayer(layer);
          });
          mapRef.current.remove();
        } catch (error) {
          console.error("Error removing map:", error);
        }
        mapRef.current = null;
      }
      setMapReady(false); // Reset map ready state
    };
  }, [userLocation, theme, t, filters.distance]);

  // Update radius circle when distance filter changes
  useEffect(() => {
    if (!mapRef.current || !userLocation || !filters.distance || !mapReady) return;

    try {
      // Remove old circles
      mapRef.current.eachLayer(layer => {
        if (layer instanceof L.Circle) {
          mapRef.current?.removeLayer(layer);
        }
      });

      // Add new circle with current radius
      const radiusInMeters = getRadiusInMeters(filters.distance);
      L.circle([userLocation.lat, userLocation.lon], {
        color: theme.palette.primary?.main,
        fillColor: theme.palette.primary?.light,
        fillOpacity: 0.2,
        radius: radiusInMeters
      }).addTo(mapRef.current);
    } catch (error) {
      console.error("Error updating radius circle:", error);
    }
  }, [filters.distance, userLocation, theme, mapReady]);

  // Update listing markers
  useEffect(() => {
    if (!mapRef.current || !mapReady) return;

    try {
      // Remove old markers
      markersRef.current.forEach(marker => {
        try {
          mapRef.current?.removeLayer(marker);
        } catch (error) {
          console.error("Error removing marker:", error);
        }
      });
      markersRef.current = [];

      // Check for listings with coordinates
      const validListings = listings.filter(listing =>
        listing.latitude && listing.longitude &&
        listing.show_on_map !== false
      );

      if (validListings.length === 0) return;

      // Create marker group for auto-scaling
      const markerGroup = L.featureGroup();

      const storefrontListings = new Map<string | number, StorefrontData>();
      const storefrontUniqueLocations = new Map<string, {
        listings: Listing[];
        storefront_id: string | number;
        name: string;
        location: [number, number];
        address?: string;
      }>();
      const individualListings: Listing[] = [];

      const isUniqueLocation = (storefrontId: string | number, lat: number, lng: number): boolean => {
        if (!storefrontListings.has(storefrontId)) return true;

        const mainLocation = storefrontListings.get(storefrontId)!.location;


        const threshold = 0.0001;
        const isDifferent = Math.abs(mainLocation[0] - lat) > threshold ||
          Math.abs(mainLocation[1] - lng) > threshold;


        if (isDifferent) {
          console.log(`Found unique coordinates for storefront listing ${storefrontId}: [${lat}, ${lng}] different from main [${mainLocation[0]}, ${mainLocation[1]}]`);
        }

        return isDifferent;
      };

      validListings.forEach(listing => {
        if (listing.storefront_id) {
          if (!storefrontListings.has(listing.storefront_id)) {
            storefrontListings.set(listing.storefront_id, {
              listings: [],
              name: listing.storefront_name || t('listings.map.storefront'),
              location: [listing.latitude!, listing.longitude!],
              address: listing.location || listing.address
            });
          }
          storefrontListings.get(listing.storefront_id)!.listings.push(listing);
        } else {

          individualListings.push(listing);
        }
      });


      const MAX_INDIVIDUAL_LISTINGS = 1000;
      const limitedIndividualListings = individualListings.slice(0, MAX_INDIVIDUAL_LISTINGS);


      const MAX_STOREFRONTS = 200;

      const limitedStorefrontIds = Array.from(storefrontListings.keys()).slice(0, MAX_STOREFRONTS);


      const MAX_ITEMS_PER_STOREFRONT = 10000000;

      console.log(`Processing ${limitedStorefrontIds.length} storefronts out of ${storefrontListings.size}`);
      console.log(`Processing ${limitedIndividualListings.length} individual listings (limit: ${MAX_INDIVIDUAL_LISTINGS})`);


      console.log("Starting search for listings with unique locations...");
      validListings.forEach(listing => {
        if (listing.storefront_id && listing.latitude && listing.longitude) {

          console.log(`Checking listing: ID=${listing.id}, storefront_id=${listing.storefront_id}, coordinates=[${listing.latitude}, ${listing.longitude}]`);


          if (isUniqueLocation(listing.storefront_id, listing.latitude, listing.longitude)) {

            const locationKey = `${listing.storefront_id}_${listing.latitude}_${listing.longitude}`;
            console.log(`Unique location: key=${locationKey}`);

            if (!storefrontUniqueLocations.has(locationKey)) {

              const storefrontName = listing.storefront_name ||
                (storefrontListings.has(listing.storefront_id) ?
                  storefrontListings.get(listing.storefront_id)!.name :
                  t('listings.map.storefront'));

              storefrontUniqueLocations.set(locationKey, {
                listings: [],
                storefront_id: listing.storefront_id,
                name: storefrontName,
                location: [listing.latitude, listing.longitude],
                address: listing.location || listing.address
              });

              console.log(`Created new unique location entry: ${storefrontName}, coordinates=[${listing.latitude}, ${listing.longitude}]`);
            }


            storefrontUniqueLocations.get(locationKey)!.listings.push(listing);
            console.log(`Listing ID=${listing.id} added to unique locations list, total: ${storefrontUniqueLocations.get(locationKey)!.listings.length}`);


            if (storefrontListings.has(listing.storefront_id)) {
              const mainStorefrontListings = storefrontListings.get(listing.storefront_id)!.listings;
              const index = mainStorefrontListings.findIndex(item => item.id === listing.id);
              if (index >= 0) {
                mainStorefrontListings.splice(index, 1);
                console.log(`Listing ID=${listing.id} removed from main storefront list, remaining: ${mainStorefrontListings.length}`);
              } else {
                console.log(`Listing ID=${listing.id} not found in main storefront list`);
              }
            } else {
              console.log(`Storefront ID=${listing.storefront_id} not found in main storefronts list`);
            }
          }
        }
      });


      console.log(`Found ${storefrontListings.size} storefronts on map`);
      storefrontListings.forEach((storefront, id) => {
        console.log(`Storefront ID ${id}: ${storefront.name}, listings: ${storefront.listings.length}`);
      });

      console.log(`Found ${storefrontUniqueLocations.size} unique locations for storefront items`);
      storefrontUniqueLocations.forEach((location, key) => {
        console.log(`Location: ${key}, Storefront ID: ${location.storefront_id}, listings: ${location.listings.length}`);
      });


      limitedStorefrontIds.forEach(storefrontId => {
        try {
          const storefront = storefrontListings.get(storefrontId);
          if (!storefront) return;
          

          const limitedStorefrontListings = storefront.listings;

          const storeMarker = L.marker(storefront.location, {
            icon: L.divIcon({
              html: `
                <div style="
                  width: 42px;
                  height: 42px;
                  display: flex;
                  flex-direction: column;
                  align-items: center;
                  justify-content: center;
                ">
                  <div style="
                    background-color: ${theme.palette.primary?.main};
                    color: white;
                    width: 40px;
                    height: 32px;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    border-radius: 4px;
                    position: relative;
                    border: 2px solid white;
                    box-shadow: 0 2px 5px rgba(0,0,0,0.3);
                  ">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z"></path>
                      <polyline points="9 22 9 12 15 12 15 22"></polyline>
                    </svg>
                    <div style="
                      position: absolute;
                      top: -10px;
                      right: -10px;
                      background-color: ${theme.palette.error?.main};
                      color: white;
                      border-radius: 50%;
                      width: 22px;
                      height: 22px;
                      font-size: 12px;
                      font-weight: bold;
                      display: flex;
                      align-items: center;
                      justify-content: center;
                      border: 1px solid white;
                    ">${storefront.listings.length > 999 ? '999+' : storefront.listings.length}</div>
                  </div>
                  <div style="
                    width: 0;
                    height: 0;
                    border-left: 8px solid transparent;
                    border-right: 8px solid transparent;
                    border-top: 10px solid ${theme.palette.primary?.main};
                    margin-top: -2px;
                  "></div>
                </div>
              `,
              className: 'storefront-marker',
              iconSize: [42, 42],
              iconAnchor: [21, 42]
            })
          })
          .bindTooltip(`${storefront.name} (${storefront.listings.length} ${t('common:map.items')})`)
            .on('click', () => {
              // Show first listing from storefront with storefront tag
              if (limitedStorefrontListings.length > 0) {
                const firstListing = limitedStorefrontListings[0];
                // Add information that it's part of a storefront
                firstListing.isPartOfStorefront = true;
                firstListing.storefrontName = storefront.name;
                firstListing.storefrontItemCount = storefront.listings.length;
                firstListing.storefront_id = storefrontId; // Add storefront ID for navigation
                // Add listing ID for highlighting when navigating to storefront page
                firstListing.id = firstListing.id || limitedStorefrontListings[0].id;
                setSelectedListing(firstListing);
              }
            });

          storeMarker.addTo(mapRef.current);
          markerGroup.addLayer(storeMarker);
          markersRef.current.push(storeMarker);
        } catch (error) {
          console.error("Error adding storefront marker:", error);
        }
      });

      // Add markers for unique location storefront listings
      storefrontUniqueLocations.forEach((location, locationKey) => {
        try {
          // Create marker with indication that it's part of storefront, but with unique location
          const uniqueStoreMarker = L.marker(location.location, {
            icon: L.divIcon({
              html: `
                <div style="
                  width: 36px;
                  height: 36px;
                  display: flex;
                  flex-direction: column;
                  align-items: center;
                  justify-content: center;
                ">
                  <div style="
                    background-color: ${theme.palette.secondary?.main};
                    color: white;
                    width: 34px;
                    height: 28px;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    border-radius: 4px;
                    position: relative;
                    border: 2px solid white;
                    box-shadow: 0 2px 5px rgba(0,0,0,0.3);
                  ">
                      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
                      <circle cx="12" cy="10" r="3"></circle>
                    </svg>
                    <div style="
                      position: absolute;
                      top: -8px;
                      right: -8px;
                      background-color: ${theme.palette.error?.main};
                      color: white;
                      border-radius: 50%;
                      width: 20px;
                      height: 20px;
                      font-size: 10px;
                      font-weight: bold;
                      display: flex;
                      align-items: center;
                      justify-content: center;
                      border: 1px solid white;
                    ">${location.listings.length > 99 ? '99+' : location.listings.length}</div>
                  </div>
                  <div style="
                    width: 0;
                    height: 0;
                    border-left: 7px solid transparent;
                    border-right: 7px solid transparent;
                    border-top: 8px solid ${theme.palette.secondary?.main};
                    margin-top: -2px;
                  "></div>
                </div>
              `,
              className: 'storefront-unique-marker',
              iconSize: [36, 36],
              iconAnchor: [18, 36]
            })
          })
            .bindTooltip(`${location.name} (${t('common:map.uniqueLocation')}) - ${location.listings.length} ${t('common:map.items')}`)

            .on('click', () => {
              // Show first listing from this unique location
              if (location.listings.length > 0) {
                // Create completely new object (important to avoid React issues)
                const uniqueListing = JSON.parse(JSON.stringify(location.listings[0]));

                // Set necessary attributes
                uniqueListing.isPartOfStorefront = true; // Key property!
                uniqueListing.storefrontName = location.name;
                uniqueListing.storefrontItemCount = location.listings.length;
                uniqueListing.storefront_id = location.storefront_id;
                uniqueListing.isUniqueLocation = true;

                // Ensure listing has ID
                if (!uniqueListing.id && location.listings[0].id) {
                  uniqueListing.id = location.listings[0].id;
                }

                console.log("Prepared object for unique location:", {
                  id: uniqueListing.id,
                  title: uniqueListing.title,
                  storefront_id: uniqueListing.storefront_id,
                  isPartOfStorefront: uniqueListing.isPartOfStorefront,
                  isUniqueLocation: uniqueListing.isUniqueLocation
                });

                // Check for id and storefront_id for navigation
                if (uniqueListing.id && uniqueListing.storefront_id) {
                  console.log("Setting selected listing for unique location");
                  setSelectedListing(uniqueListing);
                } else {
                  console.error("Could not get listing ID or storefront ID for unique location");
                }
              }
            })

          uniqueStoreMarker.addTo(mapRef.current);
          markerGroup.addLayer(uniqueStoreMarker);
          markersRef.current.push(uniqueStoreMarker);

          console.log(`Added marker for unique location items: ${location.name}, coordinates=[${location.location[0]}, ${location.location[1]}], listings: ${location.listings.length}`);
        } catch (error) {
          console.error("Error adding unique location storefront marker:", error);
        }
      });

      // Add markers for individual listings
      limitedIndividualListings.forEach(listing => {
        try {
          if (!listing.latitude || !listing.longitude) return;
          
          const marker = L.marker([listing.latitude, listing.longitude])
            .bindTooltip(`${listing.price.toLocaleString()} RSD`)
            .on('click', () => {
              setSelectedListing(listing);
            });

          marker.addTo(mapRef.current);
          markerGroup.addLayer(marker);
          markersRef.current.push(marker);
        } catch (error) {
          console.error("Error adding marker:", error);
        }
      });

      // Set map bounds to include all markers
      // if no user location
      if (!userLocation && markerGroup.getLayers().length > 0) {
        try {
          mapRef.current.fitBounds(markerGroup.getBounds(), {
            padding: [50, 50],
            maxZoom: 15
          });
        } catch (error) {
          console.error("Error fitting bounds:", error);
        }
      }
    } catch (error) {
      console.error("Error updating markers:", error);
    }
  }, [listings, mapReady, userLocation, t, theme]);

  // Handle search radius change
  const handleRadiusChange = (event: SelectChangeEvent<string>): void => {
    onFilterChange({ ...filters, distance: event.target.value });
  };

  // Navigate to listing details
  const handleNavigateToListing = (listingId: number | string): void => {
    navigate(`/marketplace/listings/${listingId}`);
  };

  // Handler for opening fullscreen map
  const handleExpandMap = (): void => {
    // Get list of all listings with coordinates
    const validListings = listings.filter(listing =>
      listing.latitude && listing.longitude && listing.show_on_map !== false
    );

    // Form markers for fullscreen map
    const markersForFullscreen = validListings.map(listing => ({
      latitude: listing.latitude!,
      longitude: listing.longitude!,
      title: listing.title,
      tooltip: `${listing.price.toLocaleString()} RSD`,
      id: listing.id,
      listing: listing // Pass full listing data
    }));

    // Determine center for fullscreen map with guaranteed coordinates
    let center: MapCenter | null = null;

    // Case 1: Selected listing
    if (selectedListing && selectedListing.latitude && selectedListing.longitude) {
      center = {
        latitude: selectedListing.latitude,
        longitude: selectedListing.longitude,
        title: selectedListing.title
      };
    }
    // Case 2: User location
    else if (userLocation && userLocation.lat && userLocation.lon) {
      center = {
        latitude: userLocation.lat,
        longitude: userLocation.lon,
        title: t('listings.map.yourLocation')
      };
    }
    // Case 3: Current map center
    else if (mapRef.current) {
      try {
        const mapCenter = mapRef.current.getCenter();
        center = {
          latitude: mapCenter.lat,
          longitude: mapCenter.lng,
          title: t('listings.map.mapCenter')
        };
      } catch (error) {
        console.error("Error getting map center:", error);
      }
    }
    // Case 4: First listing from list
    else if (validListings.length > 0) {
      const firstListing = validListings[0];
      center = {
        latitude: firstListing.latitude!,
        longitude: firstListing.longitude!,
        title: firstListing.title
      };
    }
    // Case 5: Default coordinates (Novi Sad)
    else {
      center = {
        latitude: 45.2671,
        longitude: 19.8335,
        title: "Novi Sad"
      };
    }

    // Extra check before setting state
    if (!center || !center.latitude || !center.longitude) {
      center = {
        latitude: 45.2671,
        longitude: 19.8335,
        title: "Novi Sad"
      };
    }

    // Ensure we have numeric values for coordinates
    center.latitude = Number(center.latitude);
    center.longitude = Number(center.longitude);

    // Set state and open modal
    setExpandedMapCenter(center);
    setExpandedMapMarkers(markersForFullscreen);
    setExpandedMapOpen(true);
  };

  // Function to detect user location
  const handleDetectLocation = async (): Promise<void> => {
    try {
      // Use function from location context
      const locationData = await detectUserLocation();

      // If successfully got location, update filters
      onFilterChange({
        ...filters,
        latitude: locationData.lat,
        longitude: locationData.lon,
        distance: filters.distance || '5km'
      });

      // Center map on new coordinates
      if (mapRef.current) {
        mapRef.current.setView([locationData.lat, locationData.lon], 13);
      }
    } catch (error) {
      console.error("Error getting location:", error);
      alert(t('listings.map.locationError'));
    }
  };

  // Determine if map is available (not using distance without coordinates)
  const isMapAvailable = useMemo(() => {
    return !filters.distance || (filters.latitude && filters.longitude);
  }, [filters.distance, filters.latitude, filters.longitude]);

  const isDistanceWithoutCoordinates = filters.distance && (!filters.latitude || !filters.longitude);

  return (
    <Box
      sx={{
        position: 'relative',
        width: '100%',
        display: 'flex',
        flexDirection: 'column',
        minHeight: isMobile ? 'calc(100vh - 120px)' : 'auto'
      }}
    >
      {/* Map toolbar */}
      <Paper
        elevation={3}
        sx={{
          p: 2,
          mb: 2,
          zIndex: 1000,
          display: 'flex',
          flexWrap: 'wrap',
          alignItems: 'center',
          justifyContent: 'space-between',
          gap: 2
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          {/* Show information about number of listings and storefronts on map */}
          {(() => {
            // Filter listings with coordinates and show_on_map flag
            const validListings = listings.filter(l => l.latitude && l.longitude && l.show_on_map !== false);

            // Count unique storefronts
            const storefrontsSet = new Set<string | number>();
            validListings.forEach(l => {
              if (l.storefront_id) {
                storefrontsSet.add(l.storefront_id);
              }
            });

            // Count individual listings (without storefront)
            const individualListingsCount = validListings.filter(l => !l.storefront_id).length;

            return (
              <Box sx={{ display: 'flex', gap: 1 }}>
                <Chip
                  label={`${validListings.length} ${t('listings.map.itemsOnMap')}`}
                  color="primary"
                  variant="outlined"
                />
                <Chip
                  label={`${storefrontsSet.size} ${t('listings.map.storefrontsOnMap', { defaultValue: 'storefrontsOnMap' })}`}
                  color="secondary"
                  variant="outlined"
                />
                <Chip
                  label={`${individualListingsCount} ${t('listings.map.individualItemsOnMap', { defaultValue: 'individualItemsOnMap' })}`}
                  color="info"
                  variant="outlined"
                />
              </Box>
            );
          })()}
        </Box>

        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Button
            variant="outlined"
            startIcon={<List />}
            onClick={onMapClose}
          >
            {isMobile ? t('listings.map.list') : t('listings.map.backToList')}
          </Button>
        </Box>
      </Paper>

      {/* Map container - important! Fixed height and other properties */}
      <Box
        sx={{
          position: 'relative',
          width: '100%',
          height: isMobile ? '50vh' : '60vh', // Fixed height
          borderRadius: 1,
          overflow: 'hidden',
          marginBottom: 3
        }}
      >
        {/* "Expand" button in MiniMap style */}
        <IconButton
          onClick={handleExpandMap}
          sx={{
            position: 'absolute',
            top: 8,
            right: 8,
            bgcolor: 'background.paper',
            '&:hover': {
              bgcolor: 'background.paper',
            },
            zIndex: 1000,
            boxShadow: '0 2px 6px rgba(0,0,0,0.1)'
          }}
          size="small"
        >
          <Maximize2 size={20} />
        </IconButton>

        {/* Map container must have explicit dimensions */}
        <div
          ref={mapContainerRef}
          style={{
            width: '100%',
            height: '100%',
            position: 'absolute',
            top: 0,
            left: 0
          }}
        />

        {!mapReady && (
          <Box
            sx={{
              position: 'absolute',
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              bgcolor: 'background.paper',
              zIndex: 999
            }}
          >
            <Typography variant="h6">
              {t('listings.map.loadingMap', { defaultValue: 'Loading map...' })}
            </Typography>
          </Box>
        )}

        {/* Location detection button */}
        {!userLocation && (
          <Button
            variant="contained"
            color="primary"
            startIcon={<Navigation />}
            sx={{
              position: 'absolute',
              bottom: 16,
              right: 16,
              zIndex: 1000
            }}
            onClick={handleDetectLocation}
          >
            {t('listings.map.useMyLocation')}
          </Button>
        )}
      </Box>

      {/* Attribute filters block with map - important! */}
      {filters.category_id && (
        <Box sx={{ width: '100%', marginBottom: 3 }}>
          <CentralAttributeFilters
            categoryId={filters.category_id}
            onFilterChange={handleAttributeFilterChange}
            filters={filters.attributeFilters || {}}
            resetAttributeFilters={resetAttributeFilters}
          />
        </Box>
      )}

      {/* Selected listing info */}
      {selectedListing && (
        <ListingPreview
          listing={selectedListing}
          onClose={() => setSelectedListing(null)}
          onNavigate={handleNavigateToListing}
        />
      )}

      {/* Fullscreen map modal */}
      {/* @ts-ignore */}
      <Modal
        open={expandedMapOpen}
        onClose={() => setExpandedMapOpen(false)}
        aria-labelledby="expanded-map-modal"
        aria-describedby="expanded-map-view"
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          p: 3
        }}
        children={
          <Paper
            sx={{
              width: '90%',
              maxWidth: 1200,
              maxHeight: '90vh',
              bgcolor: 'background.paper',
              borderRadius: 2,
              boxShadow: 24,
              position: 'relative',
              overflow: 'hidden'
            }}
          >
            <Box sx={{ position: 'absolute', top: 8, right: 8, zIndex: 1050 }}>
              <IconButton
                onClick={() => setExpandedMapOpen(false)}
                sx={{
                  bgcolor: 'background.paper',
                  '&:hover': {
                    bgcolor: 'background.paper',
                  },
                  boxShadow: '0 2px 6px rgba(0,0,0,0.2)'
                }}
              >
                <X size={20} />
              </IconButton>
            </Box>

            {expandedMapCenter && (
              <FullscreenMap
                latitude={expandedMapCenter.latitude}
                longitude={expandedMapCenter.longitude}
                title={expandedMapCenter.title}
                markers={expandedMapMarkers.map(marker => {
                  // Создаем новый объект маркера с точными типами
                  const newMarker: any = {
                    latitude: marker.latitude,
                    longitude: marker.longitude,
                    title: marker.title,
                    tooltip: marker.tooltip,
                    id: marker.id
                  };

                  // Если есть данные объявления, преобразуем их
                  if (marker.listing) {
                    newMarker.listing = {
                      id: marker.listing.id,
                      title: marker.listing.title,
                      price: marker.listing.price,
                      // Принудительно преобразуем тип condition
                      condition: (marker.listing.condition === 'new' ? 'new' : 'used') as 'new' | 'used',
                      images: marker.listing.images
                    };
                  }

                  return newMarker;
                })}
              />
            )}
          </Paper>
        }
      />
    </Box>
  );
};

export default MapView;