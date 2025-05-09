// frontend/hostel-frontend/src/components/marketplace/MapView.tsx
import React, { useState, useEffect, useRef, useMemo, RefObject } from 'react';
import { useTranslation } from 'react-i18next';
import { Navigation, X, List, Maximize2, Store } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
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
import '../maps/leaflet-icons'; // ;O 8A?@02;5=8O 8:>=>: Leaflet
import FullscreenMap from '../maps/FullscreenMap';
import { useLocation } from '../../contexts/LocationContext';
import CentralAttributeFilters from './CentralAttributeFilters';

// ?@545;5=85 B8?>2
export interface ImageObject {
  id?: number | string;
  file_path?: string;
  public_url?: string;
  is_main?: boolean;
  storage_type?: string;
  [key: string]: any;
}

export interface Listing {
  id: number | string;
  title: string;
  price: number;
  latitude?: number;
  longitude?: number;
  location?: string;
  address?: string;
  condition?: string;
  images?: (string | ImageObject)[];
  storefront_id?: number | string;
  storefront_name?: string;
  isPartOfStorefront?: boolean;
  isUniqueLocation?: boolean;
  storefrontName?: string;
  storefrontItemCount?: number;
  show_on_map?: boolean;
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

// ><?>=5=B 4;O ?@54?@>A<>B@0 >1JO2;5=8O ?@8 :;8:5 ?> <0@:5@C
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

  // ?B8<878@>20==0O =>@<0;870F8O ?CB8 : 87>1@065=8N
  const getImageUrl = (images?: (string | ImageObject)[]): string => {
    if (!images || !Array.isArray(images) || images.length === 0) {
      return '/placeholder-listing.jpg';
    }

    // 0E>48< 3;02=>5 87>1@065=85 8;8 8A?>;L7C5< ?5@2>5 2 A?8A:5
    let mainImage = images.find(img => img && (img as ImageObject).is_main === true) || images[0];

    // A?>;L7C5< ?5@5<5==CN >:@C65=8O 87 window.ENV 2<5AB> process.env
    const baseUrl = (window as any).ENV?.REACT_APP_MINIO_URL || (window as any).ENV?.REACT_APP_BACKEND_URL || '';
    console.log('MapView: Using baseUrl from env:', baseUrl);

    // 1. !B@>:>2K5 ?CB8 (4;O >1@0B=>9 A>2<5AB8<>AB8)
    if (typeof mainImage === 'string') {
      console.log('MapView: Processing string image path:', mainImage);

      // B=>A8B5;L=K9 ?CBL MinIO
      if (mainImage.startsWith('/listings/')) {
        const url = `${baseUrl}${mainImage}`;
        console.log('MapView: Using MinIO relative path:', url);
        return url;
      }

      // ID/filename.jpg (?@O<>9 ?CBL MinIO)
      if (mainImage.match(/^\d+\/[^\/]+$/)) {
        const url = `${baseUrl}/listings/${mainImage}`;
        console.log('MapView: Using direct MinIO path pattern:', url);
        return url;
      }

      // >:0;L=>5 E@0=8;8I5 (>1@0B=0O A>2<5AB8<>ABL)
      const url = `${baseUrl}/uploads/${mainImage}`;
      console.log('MapView: Using local storage path:', url);
      return url;
    }

    // 2. 1J5:BK A 8=D>@<0F859 > D09;5
    if (typeof mainImage === 'object' && mainImage !== null) {
      const imageObj = mainImage as ImageObject;
      console.log('MapView: Processing image object:', imageObj);

      // @8>@8B5B 1: A?>;L7C5< PublicURL 5A;8 >= 4>ABC?5=
      if (imageObj.public_url && typeof imageObj.public_url === 'string' && imageObj.public_url.trim() !== '') {
        const publicUrl = imageObj.public_url;
        console.log('MapView: Found public_url string:', publicUrl);

        // 1A>;NB=K9 URL
        if (publicUrl.startsWith('http')) {
          console.log('MapView: Using absolute URL:', publicUrl);
          return publicUrl;
        }
        // B=>A8B5;L=K9 URL A /listings/
        else if (publicUrl.startsWith('/listings/')) {
          const url = `${baseUrl}${publicUrl}`;
          console.log('MapView: Using public_url with listings path:', url);
          return url;
        }
        // @C3>9 >B=>A8B5;L=K9 URL
        else {
          const url = `${baseUrl}${publicUrl}`;
          console.log('MapView: Using general relative public_url:', url);
          return url;
        }
      }

      // @8>@8B5B 2: $>@<8@C5< URL =0 >A=>25 B8?0 E@0=8;8I0 8 ?CB8 : D09;C
      if (imageObj.file_path) {
        if (imageObj.storage_type === 'minio' || imageObj.file_path.includes('listings/')) {
          // #G8BK205< 2>7<>6=>ABL =0;8G8O ?@5D8:A0 listings/ 2 ?CB8
          const filePath = imageObj.file_path.includes('listings/')
            ? imageObj.file_path.replace('listings/', '') 
            : imageObj.file_path;

          const url = `${baseUrl}/listings/${filePath}`;
          console.log('MapView: Constructed MinIO URL from path:', url);
          return url;
        }

        // >:0;L=>5 E@0=8;8I5
        const url = `${baseUrl}/uploads/${imageObj.file_path}`;
        console.log('MapView: Using local storage path from object:', url);
        return url;
      }
    }

    console.log('MapView: Could not determine image URL, using placeholder');
    return '/placeholder-listing.jpg';
  };

  const imageUrl = getImageUrl(listing.images);

  // ?@545;O5<, ?>:07K205< ;8 <K :0@B>G:C 28B@8=K 8;8 >1JO2;5=8O
  const isStorefrontCard = Boolean(listing.isPartOfStorefront ||
    (listing.storefront_id && listing.isUniqueLocation));

  console.log('ListingPreview: >?@545;5=85 B8?0 :0@B>G:8:', {
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
              backgroundColor: theme.palette.primary.main,
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
                  label={t('common:map.uniqueLocation', { defaultValue: '#=8:0;L=K9 04@5A' })}
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
                  console.log(`5@5E>4 =0 AB@0=8FC 28B@8=K: /shop/${listing.storefront_id}?highlightedListingId=${listing.id}`, {
                    storefront_id: listing.storefront_id,
                    listing_id: listing.id
                  });
                  navigate(`/shop/${listing.storefront_id}?highlightedListingId=${listing.id}`);
                }}
              >
                {listing.isUniqueLocation
                  ? t('common:map.viewItemInStorefront', { defaultValue: 'B:@KBL 2 28B@8=5' })
                  : t('common:map.viewStorefront')}
              </Button>
            ) : (
              <Button
                variant="contained"
                size="small"
                onClick={() => {
                  console.log(`5@5E>4 =0 AB@0=8FC >1JO2;5=8O: ${listing.id}`);
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

  // !>7405< ?@028;L=K9 >1J5:B A :>>@48=0B0<8 87 userLocation
  const locationCoordinates = userLocation ? {
    latitude: userLocation.lat,
    longitude: userLocation.lon
  } : null;

  // >38@C5< 4;O >B;04:8
  console.log("MapView userLocation:", userLocation);
  console.log("MapView locationCoordinates:", locationCoordinates);

  const mapRef = useRef<L.Map | null>(null);
  const markersRef = useRef<L.Marker[]>([]);
  const mapContainerRef = useRef<HTMLDivElement | null>(null);
  const [selectedListing, setSelectedListing] = useState<Listing | null>(null);
  const [mapReady, setMapReady] = useState<boolean>(false);
  
  useEffect(() => {
    if (selectedListing) {
      console.log("!>AB>O=85 selectedListing 87<5=8;>AL:", {
        id: selectedListing.id,
        isPartOfStorefront: selectedListing.isPartOfStorefront,
        storefront_id: selectedListing.storefront_id,
        isUniqueLocation: selectedListing.isUniqueLocation
      });
    }
  }, [selectedListing]);
  
  // !>AB>O=85 4;O <>40;L=>3> >:=0 A ?>;=>M:@0==>9 :0@B>9
  const [expandedMapOpen, setExpandedMapOpen] = useState<boolean>(false);
  // &5=B@ 4;O ?>;=>M:@0==>9 :0@BK
  const [expandedMapCenter, setExpandedMapCenter] = useState<MapCenter | null>(null);
  // 0@:5@K 4;O ?>;=>M:@0==>9 :0@BK
  const [expandedMapMarkers, setExpandedMapMarkers] = useState<MapMarker[]>([]);

  // #1548<AO, GB> C =0A 2A5340 5ABL userLocationState
  const [userLocationState, setUserLocationState] = useState<{latitude: number, longitude: number} | null>(locationCoordinates);
  
  const handleAttributeFilterChange = (newAttrFilters: Record<string, any>): void => {
    console.log("MapView: ?>;CG5=K =>2K5 0B@81CB=K5 D8;LB@K:", newAttrFilters);
    if (onFilterChange) {
      onFilterChange({
        ...filters,
        attributeFilters: newAttrFilters
      });
    }
  };

  // 1@01>BG8: A1@>A0 0B@81CB=KE D8;LB@>2
  const resetAttributeFilters = (): void => {
    if (onFilterChange) {
      onFilterChange({
        ...filters,
        attributeFilters: {}
      });
    }
  };

  // 1=>2;O5< userLocationState ?@8 87<5=5=88 userLocation
  useEffect(() => {
    if (userLocation) {
      setUserLocationState({
        latitude: userLocation.lat,
        longitude: userLocation.lon
      });
      console.log("1=>2;5= userLocationState 87 userLocation:", {
        latitude: userLocation.lat,
        longitude: userLocation.lon
      });
    }
  }, [userLocation]);

  // $C=:F8O 4;O ?@5>1@07>20=8O @0AAB>O=8O (=0?@8<5@, "5km") 2 <5B@K
  const getRadiusInMeters = (distanceString: string | undefined): number => {
    if (!distanceString) return 5000; // > C<>;G0=8N 5 :<

    const match = distanceString.match(/^(\d+)km$/);
    if (match) {
      return parseInt(match[1]) * 1000;
    }
    return 5000;
  };

  // =8F80;870F8O :0@BK
  useEffect(() => {
    if (!mapContainerRef.current || mapRef.current) return;

    // >1028< 7045@6:C, GB>1K C1548BLAO, GB> DOM 3>B>2
    const timer = setTimeout(() => {
      try {
        // #AB0=02;8205< F5=B@ :0@BK =0 >A=>25 40==KE 87 :>=B5:AB0
        const initialPosition = userLocation
          ? [userLocation.lat, userLocation.lon]  // A?>;L7C5< 40==K5 87 :>=B5:AB0
          : [45.2671, 19.8335]; // >>@48=0BK >28-!040 ?> C<>;G0=8N

        // @>25@:0, GB> :>=B59=5@ ACI5AB2C5B 8 8<55B @07<5@K
        if (!mapContainerRef.current ||
          mapContainerRef.current.clientWidth === 0 ||
          mapContainerRef.current.clientHeight === 0) {
          console.log("Map container not ready yet, retrying...");
          return;
        }

        // !>7405< :0@BC A >?F859 preferCanvas 4;O ;CGH59 ?@>872>48B5;L=>AB8
        mapRef.current = L.map(mapContainerRef.current, {
          preferCanvas: true,
          attributionControl: false,
          zoomControl: true,
          inertia: true,
          fadeAnimation: true,
          zoomAnimation: true,
          renderer: L.canvas()
        }).setView(initialPosition as [number, number], 13);

        // >102;O5< A;>9 B09;>2
        L.tileLayer(TILE_LAYER_URL, {
          attribution: TILE_LAYER_ATTRIBUTION,
          maxZoom: 19
        }).addTo(mapRef.current);

        // A;8 5ABL <5AB>?>;>65=85 ?>;L7>20B5;O, 4>102;O5< <0@:5@ 8 :@C3
        if (userLocation) {
          try {
            L.circle(initialPosition as [number, number], {
              color: theme.palette.primary.main,
              fillColor: theme.palette.primary.light,
              fillOpacity: 0.2,
              radius: getRadiusInMeters(filters.distance)
            }).addTo(mapRef.current);

            L.marker(initialPosition as [number, number], {
              icon: L.divIcon({
                html: `<div style="
                  background-color: ${theme.palette.primary.main};
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

        // 1@01>BG8: 4;O 8A?@02;5=8O ?@>1;5<K _leaflet_pos
        mapRef.current.on('zoomanim', (e) => {
          // 8G53> =5 45;05<, => MB> ?><>305B ?@54>B2@0B8BL >H81:C
        });

        setMapReady(true);
      } catch (error) {
        console.error("Error initializing map:", error);
      }
    }, 300); // #25;8G8205< 7045@6:C 4;O 30@0=B88 3>B>2=>AB8 DOM

    return () => {
      clearTimeout(timer);
      if (mapRef.current) {
        try {
          // /2=> C40;O5< 2A5 A;>8 ?5@54 C40;5=85< :0@BK
          mapRef.current.eachLayer((layer) => {
            mapRef.current?.removeLayer(layer);
          });
          mapRef.current.remove();
        } catch (error) {
          console.error("Error removing map:", error);
        }
        mapRef.current = null;
      }
      setMapReady(false); // !1@0AK205< A>AB>O=85 3>B>2=>AB8 :0@BK
    };
  }, [userLocation, theme, t, filters.distance]);

  // 1=>2;O5< :@C3 @048CA0 ?@8 87<5=5=88 D8;LB@0 @0AAB>O=8O
  useEffect(() => {
    if (!mapRef.current || !userLocation || !filters.distance || !mapReady) return;

    try {
      // #40;O5< AB0@K5 :@C38
      mapRef.current.eachLayer(layer => {
        if (layer instanceof L.Circle) {
          mapRef.current?.removeLayer(layer);
        }
      });

      // >102;O5< =>2K9 :@C3 A 0:BC0;L=K< @048CA><
      const radiusInMeters = getRadiusInMeters(filters.distance);
      L.circle([userLocation.lat, userLocation.lon], {
        color: theme.palette.primary.main,
        fillColor: theme.palette.primary.light,
        fillOpacity: 0.2,
        radius: radiusInMeters
      }).addTo(mapRef.current);
    } catch (error) {
      console.error("Error updating radius circle:", error);
    }
  }, [filters.distance, userLocation, theme, mapReady]);

  // 1=>2;O5< <0@:5@K >1JO2;5=89
  useEffect(() => {
    if (!mapRef.current || !mapReady) return;

    try {
      // #40;O5< AB0@K5 <0@:5@K
      markersRef.current.forEach(marker => {
        try {
          mapRef.current?.removeLayer(marker);
        } catch (error) {
          console.error("Error removing marker:", error);
        }
      });
      markersRef.current = [];

      // @>25@O5< =0;8G85 >1JO2;5=89 A :>>@48=0B0<8
      const validListings = listings.filter(listing =>
        listing.latitude && listing.longitude &&
        listing.show_on_map !== false
      );

      if (validListings.length === 0) return;

      // !>7405< 3@C??C <0@:5@>2 4;O 02B><0AHB018@>20=8O
      const markerGroup = L.featureGroup();

      // @C??8@>2:0 >1JO2;5=89 ?> 28B@8=0< 8 >?@545;5=85 C=8:0;L=KE <5AB>?>;>65=89
      const storefrontListings = new Map<string | number, StorefrontData>(); // Map 4;O E@0=5=8O >1JO2;5=89 ?> 28B@8=0< ?> >A=>2=><C 04@5AC
      const storefrontUniqueLocations = new Map<string, {
        listings: Listing[];
        storefront_id: string | number;
        name: string;
        location: [number, number];
        address?: string;
      }>(); // Map 4;O E@0=5=8O >1JO2;5=89 28B@8= A C=8:0;L=K<8 04@5A0<8
      const individualListings: Listing[] = []; // 1JO2;5=8O 157 28B@8=K

      // $C=:F8O 4;O ?@>25@:8, O2;ONBAO ;8 :>>@48=0BK C=8:0;L=K<8
      const isUniqueLocation = (storefrontId: string | number, lat: number, lng: number): boolean => {
        // A;8 C =0A =5B 70?8A8 > 3;02=>< <5AB>?>;>65=88 28B@8=K, B> AG8B05< C=8:0;L=K<
        if (!storefrontListings.has(storefrontId)) return true;

        const mainLocation = storefrontListings.get(storefrontId)!.location;

        // @>25@O5<, >B;8G0NBAO ;8 :>>@48=0BK >B >A=>2=>3> <5AB>?>;>65=8O 28B@8=K
        // A?>;L7C5< =51>;LHCN ?>3@5H=>ABL, GB>1K 871560BL ?@>1;5< A B>G=>ABLN
        const threshold = 0.0001; // ?@8<5@=> 10 <5B@>2
        const isDifferent = Math.abs(mainLocation[0] - lat) > threshold ||
          Math.abs(mainLocation[1] - lng) > threshold;

        // K2>48< >B;04>G=CN 8=D>@<0F8N ?@8 >1=0@C65=88 C=8:0;L=KE :>>@48=0B
        if (isDifferent) {
          console.log(`1=0@C65=K C=8:0;L=K5 :>>@48=0BK 4;O >1JO2;5=8O 28B@8=K ${storefrontId}: [${lat}, ${lng}] >B;8G0NBAO >B >A=>2=KE [${mainLocation[0]}, ${mainLocation[1]}]`);
        }

        return isDifferent;
      };

      // !=0G0;0 ?@>E>48< ?> 2A5< >1JO2;5=8O< 4;O A1>@0 >A=>2=>9 8=D>@<0F88 > 28B@8=0E
      validListings.forEach(listing => {
        if (listing.storefront_id) {
          if (!storefrontListings.has(listing.storefront_id)) {
            // !>7405< =>2CN 70?8AL > 28B@8=5
            storefrontListings.set(listing.storefront_id, {
              listings: [],
              name: listing.storefront_name || t('listings.map.storefront'),
              location: [listing.latitude!, listing.longitude!],
              address: listing.location || listing.address
            });
          }
          // >102;O5< >1JO2;5=85 2 >1I89 A?8A>: 28B@8=K
          storefrontListings.get(listing.storefront_id)!.listings.push(listing);
        } else {
          // 1JO2;5=8O 157 28B@8=K A@07C 4>102;O5< 2 A?8A>: 8=48284C0;L=KE
          individualListings.push(listing);
        }
      });

      // 3@0=8G8205< :>;8G5AB2> 8=48284C0;L=KE >1JO2;5=89
      const MAX_INDIVIDUAL_LISTINGS = 1000;
      const limitedIndividualListings = individualListings.slice(0, MAX_INDIVIDUAL_LISTINGS);

      // 3@0=8G8205< :>;8G5AB2> 28B@8=
      const MAX_STOREFRONTS = 200;
      // 5@5< B>;L:> ?5@2K5 MAX_STOREFRONTS 28B@8= 5A;8 8E 1>;LH5
      const limitedStorefrontIds = Array.from(storefrontListings.keys()).slice(0, MAX_STOREFRONTS);

      // 3@0=8G8205< :>;8G5AB2> >1JO2;5=89 2 >4=>9 28B@8=5
      const MAX_ITEMS_PER_STOREFRONT = 10000000;

      console.log(`1@010BK205< ${limitedStorefrontIds.length} 28B@8= 87 ${storefrontListings.size}`);
      console.log(`1@010BK205< ${limitedIndividualListings.length} 8=48284C0;L=KE >1JO2;5=89 (;8<8B: ${MAX_INDIVIDUAL_LISTINGS})`);

      // "5?5@L ?@>E>48< ?> 2A5< >1JO2;5=8O< A=>20 4;O >?@545;5=8O C=8:0;L=KE <5AB>?>;>65=89
      console.log("0G8=05< ?>8A: >1JO2;5=89 A C=8:0;L=K<8 <5AB>?>;>65=8O<8...");
      validListings.forEach(listing => {
        if (listing.storefront_id && listing.latitude && listing.longitude) {
          // >38@C5< 8=D>@<0F8N > :064>< >1JO2;5=88 28B@8=K
          console.log(`@>25@:0 >1JO2;5=8O: ID=${listing.id}, storefront_id=${listing.storefront_id}, :>>@48=0BK=[${listing.latitude}, ${listing.longitude}]`);

          // @>25@O5<, C=8:0;L=> ;8 <5AB>?>;>65=85 4;O 40==>9 28B@8=K
          if (isUniqueLocation(listing.storefront_id, listing.latitude, listing.longitude)) {
            // !>7405< C=8:0;L=K9 :;NG 4;O MB>3> <5AB>?>;>65=8O
            const locationKey = `${listing.storefront_id}_${listing.latitude}_${listing.longitude}`;
            console.log(` #, !": key=${locationKey}`);

            if (!storefrontUniqueLocations.has(locationKey)) {
              // A;8 MB> ?5@2>5 >1JO2;5=85 A B0:8<8 :>>@48=0B0<8, A>7405< =>2CN 70?8AL
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

              console.log(`!>740=0 =>20O 70?8AL C=8:0;L=>3> <5AB>?>;>65=8O: ${storefrontName}, :>>@48=0BK=[${listing.latitude}, ${listing.longitude}]`);
            }

            // >102;O5< >1JO2;5=85 2 A?8A>: C=8:0;L=KE <5AB>?>;>65=89
            storefrontUniqueLocations.get(locationKey)!.listings.push(listing);
            console.log(`1JO2;5=85 ID=${listing.id} 4>102;5=> 2 A?8A>: C=8:0;L=KE <5AB>?>;>65=89, 2A53>: ${storefrontUniqueLocations.get(locationKey)!.listings.length}`);

            //  C40;O5< 53> 87 >A=>2=>3> A?8A:0 28B@8=K, GB>1K 871560BL 4C1;8@>20=8O
            if (storefrontListings.has(listing.storefront_id)) {
              const mainStorefrontListings = storefrontListings.get(listing.storefront_id)!.listings;
              const index = mainStorefrontListings.findIndex(item => item.id === listing.id);
              if (index >= 0) {
                mainStorefrontListings.splice(index, 1);
                console.log(`1JO2;5=85 ID=${listing.id} C40;5=> 87 >A=>2=>3> A?8A:0 28B@8=K, >AB02H8EAO: ${mainStorefrontListings.length}`);
              } else {
                console.log(`1JO2;5=85 ID=${listing.id} =5 =0945=> 2 >A=>2=>< A?8A:5 28B@8=K`);
              }
            } else {
              console.log(`8B@8=0 ID=${listing.storefront_id} =5 =0945=0 2 A?8A:5 >A=>2=KE 28B@8=`);
            }
          }
        }
      });

      // B;04>G=0O 8=D>@<0F8O
      console.log(`0945=> ${storefrontListings.size} 28B@8= =0 :0@B5`);
      storefrontListings.forEach((storefront, id) => {
        console.log(`8B@8=0 ID ${id}: ${storefront.name}, >1JO2;5=89: ${storefront.listings.length}`);
      });

      console.log(`0945=> ${storefrontUniqueLocations.size} C=8:0;L=KE <5AB>?>;>65=89 B>20@>2 28B@8=`);
      storefrontUniqueLocations.forEach((location, key) => {
        console.log(`675: ${key}, 9856 ID: ${location.storefront_id}, >1JO2;5=89: ${location.listings.length}`);
      });

      // >102;O5< <0@:5@K 4;O 28B@8=
      limitedStorefrontIds.forEach(storefrontId => {
        try {
          const storefront = storefrontListings.get(storefrontId);
          if (!storefront) return;
          
          // 3@0=8G8205< :>;8G5AB2> B>20@>2 2 >4=>9 28B@8=5 4> @07C<=>3> ?@545;0
          const limitedStorefrontListings = storefront.listings;

          // !>7405< <0@:5@ 4;O 28B@8=K 2 2845 <03078=0
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
                    background-color: ${theme.palette.primary.main};
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
                      background-color: ${theme.palette.error.main};
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
                    border-top: 10px solid ${theme.palette.primary.main};
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
              // @8 :;8:5 ?>:07K205< ?5@2>5 >1JO2;5=85 87 28B@8=K A <5B:>9 28B@8=K
              if (limitedStorefrontListings.length > 0) {
                const firstListing = limitedStorefrontListings[0];
                // >102;O5< 8=D>@<0F8N > B><, GB> MB> G0ABL 28B@8=K
                firstListing.isPartOfStorefront = true;
                firstListing.storefrontName = storefront.name;
                firstListing.storefrontItemCount = storefront.listings.length;
                firstListing.storefront_id = storefrontId; // >102;O5< ID 28B@8=K 4;O ?5@5E>40
                // >102;O5< ID >1JO2;5=8O 4;O 2K45;5=8O ?@8 ?5@5E>45 =0 AB@0=8FC 28B@8=K
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

      // >102;O5< <0@:5@K 4;O C=8:0;L=KE <5AB>?>;>65=89 >1JO2;5=89 28B@8=K
      storefrontUniqueLocations.forEach((location, locationKey) => {
        try {
          // !>7405< <0@:5@ A 8=48:0F859, GB> MB> G0ABL 28B@8=K, => A C=8:0;L=K< <5AB>?>;>65=85<
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
                    background-color: ${theme.palette.secondary.main};
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
                      background-color: ${theme.palette.error.main};
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
                    border-top: 8px solid ${theme.palette.secondary.main};
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
              // @8 :;8:5 ?>:07K205< ?5@2>5 >1JO2;5=85 A MB>3> C=8:0;L=>3> <5AB>?>;>65=8O
              if (location.listings.length > 0) {
                // !>7405< ?>;=>ABLN =>2K9 >1J5:B (206=> 4;O 871560=8O ?@>1;5< A React)
                const uniqueListing = JSON.parse(JSON.stringify(location.listings[0]));

                // #AB0=02;8205< =5>1E>48<K5 0B@81CBK
                uniqueListing.isPartOfStorefront = true; // ;NG52>5 A2>9AB2>!
                uniqueListing.storefrontName = location.name;
                uniqueListing.storefrontItemCount = location.listings.length;
                uniqueListing.storefront_id = location.storefront_id;
                uniqueListing.isUniqueLocation = true;

                // #156405<AO, GB> C ;8AB8=30 5ABL ID
                if (!uniqueListing.id && location.listings[0].id) {
                  uniqueListing.id = location.listings[0].id;
                }

                console.log(">43>B>2;5= >1J5:B 4;O C=8:0;L=>3> <5AB>?>;>65=8O:", {
                  id: uniqueListing.id,
                  title: uniqueListing.title,
                  storefront_id: uniqueListing.storefront_id,
                  isPartOfStorefront: uniqueListing.isPartOfStorefront,
                  isUniqueLocation: uniqueListing.isUniqueLocation
                });

                // @>25@O5< =0;8G85 id 8 storefront_id 4;O ?5@5E>40
                if (uniqueListing.id && uniqueListing.storefront_id) {
                  console.log("#AB0=02;8205< 2K45;5==K9 ;8AB8=3 4;O C=8:0;L=>3> <5AB>?>;>65=8O");
                  setSelectedListing(uniqueListing);
                } else {
                  console.error("5 C40;>AL ?>;CG8BL ID >1JO2;5=8O 8;8 28B@8=K 4;O C=8:0;L=>3> <5AB>?>;>65=8O");
                }
              }
            })

          uniqueStoreMarker.addTo(mapRef.current);
          markerGroup.addLayer(uniqueStoreMarker);
          markersRef.current.push(uniqueStoreMarker);

          console.log(`>102;5= <0@:5@ 4;O B>20@>2 A C=8:0;L=K< <5AB>?>;>65=85<: ${location.name}, :>>@48=0BK=[${location.location[0]}, ${location.location[1]}], >1JO2;5=89: ${location.listings.length}`);
        } catch (error) {
          console.error("Error adding unique location storefront marker:", error);
        }
      });

      // >102;O5< <0@:5@K 4;O 8=48284C0;L=KE >1JO2;5=89
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

      // #AB0=02;8205< 3@0=8FK :0@BK, GB>1K 1K;8 284=K 2A5 <0@:5@K
      // 5A;8 =5B ?>;L7>20B5;LA:>3> <5AB>?>;>65=8O
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

  // 1@01>BG8: 87<5=5=8O @048CA0 ?>8A:0
  const handleRadiusChange = (event: SelectChangeEvent<string>): void => {
    onFilterChange({ ...filters, distance: event.target.value });
  };

  // 02830F8O : ?>4@>1=>ABO< >1JO2;5=8O
  const handleNavigateToListing = (listingId: number | string): void => {
    navigate(`/marketplace/listings/${listingId}`);
  };

  // 1@01>BG8: 4;O >B:@KB8O ?>;=>M:@0==>9 :0@BK
  const handleExpandMap = (): void => {
    // >;CG05< A?8A>: 2A5E >1JO2;5=89 A :>>@48=0B0<8
    const validListings = listings.filter(listing =>
      listing.latitude && listing.longitude && listing.show_on_map !== false
    );

    // $>@<8@C5< <0@:5@K 4;O ?>;=>M:@0==>9 :0@BK
    const markersForFullscreen = validListings.map(listing => ({
      latitude: listing.latitude!,
      longitude: listing.longitude!,
      title: listing.title,
      tooltip: `${listing.price.toLocaleString()} RSD`,
      id: listing.id,
      listing: listing // 5@5405< ?>;=K5 40==K5 > ;8AB8=35
    }));

    // ?@545;O5< F5=B@ 4;O ?>;=>M:@0==>9 :0@BK A 30@0=B8@>20==K<8 :>>@48=0B0<8
    let center: MapCenter | null = null;

    // 5@2K9 A;CG09: 2K1@0==>5 >1JO2;5=85
    if (selectedListing && selectedListing.latitude && selectedListing.longitude) {
      center = {
        latitude: selectedListing.latitude,
        longitude: selectedListing.longitude,
        title: selectedListing.title
      };
    }
    // B>@>9 A;CG09: <5AB>?>;>65=85 ?>;L7>20B5;O
    else if (userLocation && userLocation.lat && userLocation.lon) {
      center = {
        latitude: userLocation.lat,
        longitude: userLocation.lon,
        title: t('listings.map.yourLocation')
      };
    }
    // "@5B89 A;CG09: B5:CI89 F5=B@ :0@BK
    else if (mapRef.current) {
      try {
        const mapCenter = mapRef.current.getCenter();
        center = {
          latitude: mapCenter.lat,
          longitude: mapCenter.lng,
          title: t('listings.map.mapCenter')
        };
      } catch (error) {
        console.error("H81:0 ?@8 ?>;CG5=88 F5=B@0 :0@BK:", error);
      }
    }
    // '5B25@BK9 A;CG09: ?5@2>5 >1JO2;5=85 87 A?8A:0
    else if (validListings.length > 0) {
      const firstListing = validListings[0];
      center = {
        latitude: firstListing.latitude!,
        longitude: firstListing.longitude!,
        title: firstListing.title
      };
    }
    // OBK9 A;CG09: D8:A8@>20==K5 :>>@48=0BK ?> C<>;G0=8N (>28-!04)
    else {
      center = {
        latitude: 45.2671,
        longitude: 19.8335,
        title: ">28-!04"
      };
    }

    // >?>;=8B5;L=0O ?@>25@:0 ?5@54 CAB0=>2:>9 A>AB>O=8O
    if (!center || !center.latitude || !center.longitude) {
      center = {
        latitude: 45.2671,
        longitude: 19.8335,
        title: ">28-!04"
      };
    }

    // @>25@O5<, GB> C =0A 5ABL G8A;>2K5 7=0G5=8O 4;O :>>@48=0B
    center.latitude = Number(center.latitude);
    center.longitude = Number(center.longitude);

    // #AB0=02;8205< A>AB>O=85 8 >B:@K205< <>40;L=>5 >:=>
    setExpandedMapCenter(center);
    setExpandedMapMarkers(markersForFullscreen);
    setExpandedMapOpen(true);
  };

  // $C=:F8O 4;O >?@545;5=8O <5AB>?>;>65=8O ?>;L7>20B5;O
  const handleDetectLocation = async (): Promise<void> => {
    try {
      // A?>;L7C5< DC=:F8N 87 :>=B5:AB0 <5AB>?>;>65=8O
      const locationData = await detectUserLocation();

      // A;8 CA?5H=> ?>;CG8;8 <5AB>?>;>65=85, >1=>2;O5< D8;LB@K
      onFilterChange({
        ...filters,
        latitude: locationData.lat,
        longitude: locationData.lon,
        distance: filters.distance || '5km'
      });

      // &5=B@8@C5< :0@BC =0 =>2KE :>>@48=0B0E
      if (mapRef.current) {
        mapRef.current.setView([locationData.lat, locationData.lon], 13);
      }
    } catch (error) {
      console.error("Error getting location:", error);
      alert(t('listings.map.locationError'));
    }
  };

  // ?@545;O5<, 4>ABC?=0 ;8 :0@B0 (=5 8A?>;L7C5BAO 2 70?@>A5 distance 157 :>>@48=0B)
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
      {/* 0=5;L 8=AB@C<5=B>2 :0@BK */}
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
          {/* >:07K205< 8=D>@<0F8N > :>;8G5AB25 >1JO2;5=89 8 28B@8= =0 :0@B5 */}
          {(() => {
            // $8;LB@C5< >1JO2;5=8O A :>>@48=0B0<8 8 D;03>< show_on_map
            const validListings = listings.filter(l => l.latitude && l.longitude && l.show_on_map !== false);

            // >4AG8BK205< :>;8G5AB2> C=8:0;L=KE 28B@8=
            const storefrontsSet = new Set<string | number>();
            validListings.forEach(l => {
              if (l.storefront_id) {
                storefrontsSet.add(l.storefront_id);
              }
            });

            // >4AG8BK205< :>;8G5AB2> G0AB=KE >1JO2;5=89 (157 28B@8=K)
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

      {/* >=B59=5@ :0@BK - 206=>! $8:A8@>20==0O 2KA>B0 8 >AB0;L=K5 A2>9AB20 */}
      <Box
        sx={{
          position: 'relative',
          width: '100%',
          height: isMobile ? '50vh' : '60vh', // $8:A8@>20==0O 2KA>B0
          borderRadius: 1,
          overflow: 'hidden',
          marginBottom: 3
        }}
      >
        {/* =>?:0 " 0725@=CBL" 2 AB8;5 MiniMap */}
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

        {/* >=B59=5@ 4;O :0@BK 4>;65= 8<5BL O2=K5 @07<5@K */}
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
              {t('listings.map.loadingMap', { defaultValue: '03@C7:0 :0@BK...' })}
            </Typography>
          </Box>
        )}

        {/* =>?:0 >?@545;5=8O <5AB>?>;>65=8O */}
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

      {/* ;>: A 0B@81CB=K<8 D8;LB@0<8  :0@B>9 - 206=>! */}
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

      {/* =D>@<0F8O > 2K1@0==>< >1JO2;5=88 */}
      {selectedListing && (
        <ListingPreview
          listing={selectedListing}
          onClose={() => setSelectedListing(null)}
          onNavigate={handleNavigateToListing}
        />
      )}

      {/* >40;L=>5 >:=> A ?>;=>M:@0==>9 :0@B>9 */}
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
      >
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
              markers={expandedMapMarkers}
            />
          )}
        </Paper>
      </Modal>
    </Box>
  );
};

export default MapView;