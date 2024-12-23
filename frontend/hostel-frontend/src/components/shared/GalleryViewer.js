// frontend/hostel-frontend/src/components/shared/GalleryViewer.js
import React, { useState } from 'react';
import {
    Dialog,
    IconButton,
    Box,
    Grid,
    DialogContent,
    Stack,
} from '@mui/material';
import {
    Close as CloseIcon,
    ChevronLeft as ChevronLeftIcon,
    ChevronRight as ChevronRightIcon,
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
    };

    const handleClose = (e) => {
        e?.stopPropagation();
        setSelectedIndex(null);
        if (externalClose) {
            externalClose();
        }
    };

    const handlePrev = (e) => {
        e?.stopPropagation();
        setSelectedIndex(prev => (prev > 0 ? prev - 1 : images.length - 1));
    };

    const handleNext = (e) => {
        e?.stopPropagation();
        setSelectedIndex(prev => (prev < images.length - 1 ? prev + 1 : 0));
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
                        justifyContent: 'space-between'
                    }}
                >
                    {/* Основное изображение */}
                    <Box sx={{
                        flex: 1,
                        width: '100%',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        position: 'relative'
                    }}>
                        <IconButton
                            onClick={handleClose}
                            sx={{
                                position: 'absolute',
                                right: 8,
                                top: 8,
                                color: 'white',
                                zIndex: 1,
                                '&:hover': { bgcolor: 'rgba(255, 255, 255, 0.1)' }
                            }}
                        >
                            <CloseIcon />
                        </IconButton>

                        {images.length > 1 && (
                            <>
                                <IconButton
                                    onClick={handlePrev}
                                    sx={{
                                        position: 'absolute',
                                        left: 8,
                                        color: 'white',
                                        zIndex: 1,
                                        '&:hover': { bgcolor: 'rgba(255, 255, 255, 0.1)' }
                                    }}
                                >
                                    <ChevronLeftIcon />
                                </IconButton>
                                <IconButton
                                    onClick={handleNext}
                                    sx={{
                                        position: 'absolute',
                                        right: 8,
                                        color: 'white',
                                        zIndex: 1,
                                        '&:hover': { bgcolor: 'rgba(255, 255, 255, 0.1)' }
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
                                maxWidth: '100%',
                                maxHeight: 'calc(100vh - 120px)', // Оставляем место для превью
                                objectFit: 'contain',
                                cursor: 'pointer'
                            }}
                            onClick={(e) => {
                                e.stopPropagation();
                                if (images.length > 1) {
                                    handleNext(e);
                                }
                            }}
                        />
                    </Box>

                    {/* Полоса превью */}
                    {images.length > 1 && (
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