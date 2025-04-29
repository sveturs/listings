//frontend/hostel-frontend/src/components/reviews/PhotoViewer.js
import React from 'react';
import {
    Dialog,
    IconButton,
    Box,
} from '@mui/material';
import {
    ChevronLeft as ChevronLeftIcon,
    ChevronRight as ChevronRightIcon,
    Close as CloseIcon,
} from '@mui/icons-material';

const PhotoViewer = ({ open, onClose, photos, currentIndex = 0 }) => {
    const [index, setIndex] = React.useState(currentIndex);
    const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

    const handlePrevious = (e) => {
        e.stopPropagation();
        setIndex((prev) => (prev > 0 ? prev - 1 : photos.length - 1));
    };

    const handleNext = (e) => {
        e.stopPropagation();
        setIndex((prev) => (prev < photos.length - 1 ? prev + 1 : 0));
    };

    return (
        <Dialog
            open={open}
            onClose={onClose}
            maxWidth="xl"
            fullWidth
            onClick={onClose}
        >
            <Box
                sx={{
                    position: 'relative',
                    height: 'calc(100vh - 64px)',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    bgcolor: 'black'
                }}
            >
                {/* Кнопка закрытия */}
                <IconButton
                    onClick={onClose}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        color: 'white'
                    }}
                >
                    <CloseIcon />
                </IconButton>

                {/* Навигационные кнопки */}
                {photos.length > 1 && (
                    <>
                        <IconButton
                            onClick={handlePrevious}
                            sx={{
                                position: 'absolute',
                                left: 8,
                                color: 'white',
                                '&:hover': { bgcolor: 'rgba(255, 255, 255, 0.1)' }
                            }}
                        >
                            <ChevronLeftIcon fontSize="large" />
                        </IconButton>
                        <IconButton
                            onClick={handleNext}
                            sx={{
                                position: 'absolute',
                                right: 8,
                                color: 'white',
                                '&:hover': { bgcolor: 'rgba(255, 255, 255, 0.1)' }
                            }}
                        >
                            <ChevronRightIcon fontSize="large" />
                        </IconButton>
                    </>
                )}

                {/* Изображение */}
                <Box
                    component="img"
                    src={photos[index].includes('listings/')
                        ? `${BACKEND_URL}/listings/${photos[index].split('/').pop()}`
                        : `${BACKEND_URL}/uploads/${photos[index]}`}
                    alt={`Photo ${index + 1}`}
                    sx={{
                        maxHeight: '100%',
                        maxWidth: '100%',
                        objectFit: 'contain',
                        cursor: 'zoom-out'
                    }}
                    onClick={(e) => {
                        e.stopPropagation();
                        onClose();
                    }}
                />
            </Box>
        </Dialog>
    );
};

export default PhotoViewer;