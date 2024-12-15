// src/components/ImageGallery.js
import React, { useState } from 'react';
import {
    Dialog,
    IconButton,
    Box,
    Grid,
    DialogContent,
} from '@mui/material';
import {
    Close as CloseIcon,
    ChevronLeft as ChevronLeftIcon,
    ChevronRight as ChevronRightIcon,
} from '@mui/icons-material';

const BACKEND_URL = process.env.REACT_APP_BACKEND_URL;

const ImageGallery = ({ images }) => {
    const [selectedIndex, setSelectedIndex] = useState(null);

    const handleOpen = (index) => setSelectedIndex(index);
    const handleClose = () => setSelectedIndex(null);
    const handlePrev = () => setSelectedIndex(prev => (prev > 0 ? prev - 1 : images.length - 1));
    const handleNext = () => setSelectedIndex(prev => (prev < images.length - 1 ? prev + 1 : 0));

    return (
        <>
            {/* Превью изображений */}
            <Grid container spacing={1}>
                {images?.map((image, index) => (
                    <Grid item xs={4} sm={3} md={2} key={index}>
                        <Box
                            component="img"
                            src={`${BACKEND_URL}/uploads/${image.file_path}`}
                            alt={`Preview ${index + 1}`}
                            sx={{
                                width: '100%',
                                height: '100px',
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

            {/* Диалог для полноразмерного просмотра */}
            <Dialog
                open={selectedIndex !== null}
                onClose={handleClose}
                maxWidth="xl"
                fullWidth
            >
                <DialogContent sx={{ position: 'relative', p: 0, bgcolor: 'black' }}>
                    <IconButton
                        onClick={handleClose}
                        sx={{ position: 'absolute', right: 8, top: 8, color: 'white' }}
                    >
                        <CloseIcon />
                    </IconButton>
                    
                    {selectedIndex !== null && (
                        <Box
                            component="img"
                            src={`${BACKEND_URL}/uploads/${images[selectedIndex].file_path}`}
                            alt={`Image ${selectedIndex + 1}`}
                            sx={{
                                width: '100%',
                                height: '90vh',
                                objectFit: 'contain'
                            }}
                        />
                    )}

                    <IconButton
                        onClick={handlePrev}
                        sx={{ position: 'absolute', left: 8, top: '50%', color: 'white' }}
                    >
                        <ChevronLeftIcon />
                    </IconButton>

                    <IconButton
                        onClick={handleNext}
                        sx={{ position: 'absolute', right: 8, top: '50%', color: 'white' }}
                    >
                        <ChevronRightIcon />
                    </IconButton>
                </DialogContent>
            </Dialog>
        </>
    );
};

export default ImageGallery;