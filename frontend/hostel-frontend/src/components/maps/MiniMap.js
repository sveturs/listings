//frontend/hostel-frontend/src/components/maps/MiniMap.js
import React from 'react';
import { GoogleMap, Marker } from '@react-google-maps/api';
import { Box, IconButton, Typography } from '@mui/material';
import { MapPin, Maximize2 } from 'lucide-react';

const MiniMap = ({ latitude, longitude, title, address, onClick, onExpand }) => {
    const mapContainerStyle = {
        width: '100%',
        height: '200px',
        borderRadius: '4px'
    };

    const center = {
        lat: latitude,
        lng: longitude
    };

    const options = {
        disableDefaultUI: true,
        zoomControl: true,
        clickableIcons: false,
        scrollwheel: false,
        gestureHandling: "greedy"
    };

    return (
        <Box sx={{ position: 'relative' }}>
            <GoogleMap
                mapContainerStyle={mapContainerStyle}
                center={center}
                zoom={14}
                options={options}
                onClick={onClick}
            >
                <Marker position={center} />
            </GoogleMap>

            {onExpand && (
                <IconButton
                    onClick={onExpand}
                    sx={{
                        position: 'absolute',
                        top: 8,
                        right: 8,
                        bgcolor: 'background.paper',
                        '&:hover': {
                            bgcolor: 'background.paper',
                        }
                    }}
                >
                    <Maximize2 size={20} />
                </IconButton>
            )}

            <Box
                sx={{
                    position: 'absolute',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    bgcolor: 'rgba(255, 255, 255, 0.9)',
                    p: 1,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1
                }}
            >
                <MapPin size={16} />
                <Typography variant="body2" noWrap>
                    {address}
                </Typography>
            </Box>
        </Box>
    );
};

export default MiniMap;
