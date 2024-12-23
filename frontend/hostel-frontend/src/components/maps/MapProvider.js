//frontend/hostel-frontend/src/components/maps/MapProvider.js
import React, { createContext, useContext, useState, useCallback } from 'react';
import { LoadScript } from '@react-google-maps/api';

const MapContext = createContext(null);
const libraries = ["places", "geometry"];

export const MapProvider = ({ children }) => {
    const [isLoaded, setIsLoaded] = useState(false);

    const handleLoad = useCallback(() => {
        setIsLoaded(true);
    }, []);

    return (
        <MapContext.Provider value={{ isLoaded }}>
            <LoadScript
                googleMapsApiKey={process.env.REACT_APP_GOOGLE_MAPS_API_KEY}
                libraries={libraries}
                onLoad={handleLoad}
            >
                {children}
            </LoadScript>
        </MapContext.Provider>
    );
};

export const useMap = () => {
    const context = useContext(MapContext);
    if (!context) {
        throw new Error('useMap must be used within a MapProvider');
    }
    return context;
};

export default MapProvider;