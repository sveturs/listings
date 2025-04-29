// frontend/hostel-frontend/src/components/shared/GalleryViewer.js
import React, { useState, useRef, useEffect } from 'react';
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
    const [isTransitioning, setIsTransitioning] = useState(false);
    const [transitionDirection, setTransitionDirection] = useState(null);
    const [dragPosition, setDragPosition] = useState(0); // Новое состояние для перетаскивания
    const [isDragging, setIsDragging] = useState(false); // Флаг, что идет перетаскивание
    const isOpen = externalOpen !== undefined ? externalOpen : selectedIndex !== null;

    // Ссылки для обработки свайпов
    const containerRef = useRef(null);
    const touchStartX = useRef(null);
    const touchStartY = useRef(null);
    const currentTouchX = useRef(null); // Текущая позиция касания

    // Добавляем эффект для обработки клавиш клавиатуры
    useEffect(() => {
        const handleKeyDown = (e) => {
            if (!isOpen) return;

            if (e.key === 'ArrowLeft') {
                handlePrev(e);
            } else if (e.key === 'ArrowRight') {
                handleNext(e);
            } else if (e.key === 'Escape') {
                handleClose(e);
            }
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => {
            window.removeEventListener('keydown', handleKeyDown);
        };
    }, [isOpen, selectedIndex]); // eslint-disable-line react-hooks/exhaustive-deps

    if (!images || images.length === 0) return null;

    const getImageUrl = (image) => {
        if (!image) return '/placeholder.jpg';
        
        const baseUrl = process.env.REACT_APP_BACKEND_URL || '';
        
        // Для строк (обратная совместимость)
        if (typeof image === 'string') {
            // Если это путь из MinIO (начинается с /listings/ или содержит listings/)
            if (image.startsWith('/listings/') || image.includes('listings/')) {
                return `${baseUrl}${image}`;
            }
            return `${baseUrl}/uploads/${image}`;
        }
        
        // Для объектов с информацией о файле
        if (image.file_path) {
            // Для MinIO объектов
            if (image.storage_type === 'minio' || image.file_path.includes('listings/')) {
                if (image.public_url && image.public_url.startsWith('/listings/')) {
                    return `${baseUrl}${image.public_url}`;
                }
                return `${baseUrl}/listings/${image.file_path}`;
            }
            
            // Для локального хранилища
            return `${baseUrl}/uploads/${image.file_path}`;
        }
        
        return '/placeholder.jpg';
    }
    
    

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
        if (isDragging) return; // Не меняем при активном перетаскивании
        setIsZoomed(false);

        // Моментально меняем изображение без анимации
        setSelectedIndex(prev => (prev > 0 ? prev - 1 : images.length - 1));
    };

    const handleNext = (e) => {
        e?.stopPropagation();
        if (isDragging) return; // Не меняем при активном перетаскивании
        setIsZoomed(false);

        // Моментально меняем изображение без анимации
        setSelectedIndex(prev => (prev < images.length - 1 ? prev + 1 : 0));
    };

    const toggleZoom = (e) => {
        e.stopPropagation();
        setIsZoomed(!isZoomed);
    };

    // Обработчик прокрутки колесика мыши
    const handleWheel = (e) => {
        if (isZoomed) return; // Не обрабатываем в режиме зума

        if (e.deltaY < 0) {
            // Прокрутка вверх - следующая фотография
            handleNext(e);
        } else if (e.deltaY > 0) {
            // Прокрутка вниз - предыдущая фотография
            handlePrev(e);
        }
        e.preventDefault(); // Предотвращаем стандартную прокрутку страницы
    };

    // Проверка, является ли устройство мобильным
    const isMobile = () => {
        return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
    };

    // Обработчик начала касания (для свайпов)
    const handleTouchStart = (e) => {
        if (isZoomed) return; // Не обрабатываем в режиме зума
        if (!isMobile()) return; // Работаем только на мобильных устройствах

        touchStartX.current = e.touches[0].clientX;
        touchStartY.current = e.touches[0].clientY;
        currentTouchX.current = e.touches[0].clientX;
        setIsDragging(true);
        setDragPosition(0);
    };

    // Обработчик перемещения пальца
    const handleTouchMove = (e) => {
        if (isZoomed || !isDragging || !touchStartX.current || !isMobile()) return;

        const touchX = e.touches[0].clientX;
        const deltaX = touchX - touchStartX.current;

        // Проверяем, что свайп в основном горизонтальный
        const touchY = e.touches[0].clientY;
        const deltaY = Math.abs(touchY - touchStartY.current);

        // Если движение больше вертикальное, то прекращаем обработку
        if (deltaY > Math.abs(deltaX) * 0.8) {
            return;
        }

        // Обновляем позицию для эффекта перетаскивания
        currentTouchX.current = touchX;
        setDragPosition(deltaX);

        // Предотвращаем прокрутку страницы при свайпе
        e.preventDefault();
    };

    // Обработчик окончания касания (для свайпов)
    const handleTouchEnd = (e) => {
        if (isZoomed || !isDragging || !touchStartX.current || !isMobile()) {
            setIsDragging(false);
            setDragPosition(0);
            return;
        }

        const touchEndX = e.changedTouches[0].clientX;
        const deltaX = touchEndX - touchStartX.current;

        // Определяем направление и порог перелистывания
        const threshold = window.innerWidth * 0.15; // 15% ширины экрана

        if (Math.abs(deltaX) > threshold) {
            if (deltaX > 0) {
                // Свайп вправо - предыдущая фотография
                setSelectedIndex(prev => (prev > 0 ? prev - 1 : images.length - 1));
            } else {
                // Свайп влево - следующая фотография
                setSelectedIndex(prev => (prev < images.length - 1 ? prev + 1 : 0));
            }
        } else {
            // Если свайп не достиг порога, возвращаем изображение в исходное положение с анимацией
            setDragPosition(0);
            // Добавляем небольшую задержку перед сбросом isDragging для плавного возврата
            setTimeout(() => {
                setIsDragging(false);
            }, 150);
            return;
        }

        // Сбрасываем все значения перетаскивания
        touchStartX.current = null;
        touchStartY.current = null;
        currentTouchX.current = null;
        setIsDragging(false);
        setDragPosition(0);
    };

    // Обработчик отмены касания
    const handleTouchCancel = () => {
        if (!isMobile()) return;

        touchStartX.current = null;
        touchStartY.current = null;
        currentTouchX.current = null;
        setIsDragging(false);
        setDragPosition(0);
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
                    ref={containerRef}
                    onWheel={handleWheel}
                    onTouchStart={handleTouchStart}
                    onTouchMove={handleTouchMove}
                    onTouchEnd={handleTouchEnd}
                    onTouchCancel={handleTouchCancel}
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
                                transition: isDragging ? 'none' : 'transform 0.15s ease-out',
                                transform: isDragging
                                    ? `translateX(${dragPosition}px)`
                                    : 'translateX(0)',
                                opacity: 1,
                                willChange: 'transform' // Оптимизация производительности анимации
                            }}
                            onClick={(e) => {
                                e.stopPropagation();
                                if (!isDragging) { // Только если не перетаскиваем
                                    toggleZoom(e);
                                }
                            }}
                        />
                    </Box>

                    {/* Полоса превью */}
                    {images.length > 1 && !isZoomed && (
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