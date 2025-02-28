// frontend/hostel-frontend/src/components/shared/GalleryViewer.js
import React, { useState } from 'react';
import {
    Dialog,
    IconButton,
    Box,
    Grid,
    DialogContent,
    Stack,
    Tooltip,
    Zoom,
} from '@mui/material';
import {
    Close as CloseIcon,
    ChevronLeft as ChevronLeftIcon,
    ChevronRight as ChevronRightIcon,
    ZoomIn as ZoomInIcon,
    ZoomOut as ZoomOutIcon,
} from '@mui/icons-material';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

const GalleryViewer = ({
    images,
    open: externalOpen,
    onClose: externalClose,
    initialIndex = 0,
    galleryMode = 'thumbnails',
    thumbnailSize = { width: '100%', height: '100px' },
    gridColumns = { xs: 4, sm: 3, md: 2 }
}) => {
    const [selectedIndex, setSelectedIndex] = useState(galleryMode === 'fullscreen' ? initialIndex : null);
    const [isZoomed, setIsZoomed] = useState(false);
    const isOpen = externalOpen !== undefined ? externalOpen : selectedIndex !== null;

    if (!images || images.length === 0) return null;

    const getImageUrl = (image) => {
        if (!image) return '';
        if (typeof image === 'string') return `${BACKEND_URL}/uploads/${image}`;
        if (image.file_path) return `${BACKEND_URL}/uploads/${image.file_path}`;
        return '';
    };

    const handleOpen = (index) => {
        setSelectedIndex(index);
        setIsZoomed(false);
    };

    const handleClose = (e) => {
        e?.stopPropagation();
        setSelectedIndex(null);
        setIsZoomed(false);
        if (externalClose) {
            externalClose();
        }
    };

    const handlePrev = (e) => {
        e?.stopPropagation();
        setIsZoomed(false);
        setSelectedIndex(prev => (prev > 0 ? prev - 1 : images.length - 1));
    };

    const handleNext = (e) => {
        e?.stopPropagation();
        setIsZoomed(false);
        setSelectedIndex(prev => (prev < images.length - 1 ? prev + 1 : 0));
    };

    const toggleZoom = (e) => {
        e.stopPropagation();
        setIsZoomed(!isZoomed);
    };

    return (
        <>
            {/* Превью изображений */}
            {galleryMode === 'thumbnails' && (
                <Grid container spacing={1}>
                    {images.map((image, index) => (
                        <Grid item {...gridColumns} key={index}>
                            <Box
                                component="img"
                                src={getImageUrl(image)}
                                alt={`Preview ${index + 1}`}
                                sx={{
                                    width: thumbnailSize.width,
                                    height: thumbnailSize.height,
                                    objectFit: 'cover',
                                    borderRadius: 1,
                                    cursor: 'pointer',
                                    '&:hover': {
                                        opacity: 0.8,
                                        transform: 'scale(1.05)',
                                        transition: 'all 0.2s'
                                    }
                                }}
                                onClick={() => handleOpen(index)}
                            />
                        </Grid>
                    ))}
                </Grid>
            )}

            {/* Полноэкранный просмотр */}
            <Dialog
                open={isOpen}
                onClose={handleClose}
                maxWidth="xl"
                fullWidth
                onClick={handleClose}
                sx={{
                    '.MuiDialog-paper': {
                        m: 0,
                        maxHeight: '100vh',
                        maxWidth: '100vw',
                        width: '100%',
                        height: '100%',
                        bgcolor: 'black'
                    }
                }}
            >
                <DialogContent
                    sx={{
                        position: 'relative',
                        p: 0,
                        height: '100vh',
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                        overflow: isZoomed ? 'auto' : 'hidden'
                    }}
                >
                    {/* Основное изображение */}
                    <Box sx={{
                        flex: 1,
                        width: '100%',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        position: 'relative',
                        overflow: isZoomed ? 'auto' : 'hidden'
                    }}>
                        <IconButton
                            onClick={handleClose}
                            sx={{
                                position: 'absolute',
                                right: 16,
                                top: 16,
                                color: 'white',
                                zIndex: 10,
                                bgcolor: 'rgba(0, 0, 0, 0.3)',
                                '&:hover': { bgcolor: 'rgba(0, 0, 0, 0.5)' }
                            }}
                        >
                            <CloseIcon />
                        </IconButton>

                        <Tooltip title={isZoomed ? "Уменьшить" : "Увеличить до оригинального размера"}>
                            <IconButton
                                onClick={toggleZoom}
                                sx={{
                                    position: 'absolute',
                                    right: 16,
                                    top: 70,
                                    color: 'white',
                                    zIndex: 10,
                                    bgcolor: 'rgba(0, 0, 0, 0.3)',
                                    '&:hover': { bgcolor: 'rgba(0, 0, 0, 0.5)' }
                                }}
                            >
                                {isZoomed ? <ZoomOutIcon /> : <ZoomInIcon />}
                            </IconButton>
                        </Tooltip>

                        {images.length > 1 && !isZoomed && (
                            <>
                                <IconButton
                                    onClick={handlePrev}
                                    sx={{
                                        position: 'absolute',
                                        left: 16,
                                        backgroundColor: 'rgba(0, 0, 0, 0.3)',
                                        color: 'white',
                                        zIndex: 10,
                                        '&:hover': { bgcolor: 'rgba(0, 0, 0, 0.5)' }
                                    }}
                                >
                                    <ChevronLeftIcon />
                                </IconButton>
                                <IconButton
                                    onClick={handleNext}
                                    sx={{
                                        position: 'absolute',
                                        right: 16,
                                        backgroundColor: 'rgba(0, 0, 0, 0.3)',
                                        color: 'white',
                                        zIndex: 10,
                                        '&:hover': { bgcolor: 'rgba(0, 0, 0, 0.5)' }
                                    }}
                                >
                                    <ChevronRightIcon />
                                </IconButton>
                            </>
                        )}

                        <Box
                            component="img"
                            src={getImageUrl(images[selectedIndex || 0])}
                            alt={`Image ${(selectedIndex || 0) + 1}`}
                            sx={{
                                maxWidth: isZoomed ? 'none' : '100%',
                                maxHeight: isZoomed ? 'none' : 'calc(100vh - 120px)', // Оставляем место для превью
                                width: isZoomed ? 'auto' : 'auto',
                                height: isZoomed ? 'auto' : 'auto',
                                objectFit: 'contain',
                                cursor: isZoomed ? 'zoom-out' : 'zoom-in',
                                transition: 'transform 0.3s ease'
                            }}
                            onClick={(e) => {
                                e.stopPropagation();
                                toggleZoom(e);
                            }}
                        />
                    </Box>

                    {/* Полоса превью */}
                  11  {images.length > 1 && !isZoomed && (
                        <Stack
                            direction="row"
                            spacing={1}
                            sx={{
                                p: 1,
                                width: '100%',
                                overflowX: 'auto',
                                bgcolor: 'rgba(0, 0, 0, 0.5)',
                                height: 100,
                                alignItems: 'center'
                            }}
                            onClick={(e) => e.stopPropagation()}
                        >
                            {images.map((image, index) => (
                                <Box
                                    key={index}
                                    component="img"
                                    src={getImageUrl(image)}
                                    alt={`Thumbnail ${index + 1}`}
                                    onClick={() => setSelectedIndex(index)}
                                    sx={{
                                        height: 80,
                                        width: 'auto',
                                        cursor: 'pointer',
                                        borderRadius: 1,
                                        opacity: selectedIndex === index ? 1 : 0.6,
                                        transition: 'all 0.2s',
                                        border: selectedIndex === index ? '2px solid white' : 'none',
                                        '&:hover': {
                                            opacity: 1
                                        }
                                    }}
                                />
                            ))}
                        </Stack>
                    )}
                </DialogContent>
            </Dialog>
        </>
    );
};

export default GalleryViewer;