import React, { useCallback, useEffect, useRef } from 'react';
import { Box, Button, CircularProgress, Typography } from '@mui/material';
import { useTranslation } from 'react-i18next';
import { ChevronDown } from 'lucide-react';

const InfiniteScroll = ({ 
  hasMore, 
  loading, 
  onLoadMore, 
  children, 
  autoLoad = false,
  loadingMessage = 'Loading...',
  loadMoreButtonText = 'Show more',
  noMoreItemsText = 'No more items to show'
}) => {
  const { t } = useTranslation('marketplace');
  const observer = useRef(null);
  const loadMoreRef = useRef(null);

  const handleObserver = useCallback(
    (entries) => {
      const [target] = entries;
      if (target.isIntersecting && hasMore && !loading && autoLoad) {
        onLoadMore();
      }
    },
    [hasMore, loading, onLoadMore, autoLoad]
  );

  useEffect(() => {
    const currentRef = loadMoreRef.current;
    
    // Основное изменение: подключаем observer только если hasMore = true
    if (currentRef && hasMore) {
      observer.current = new IntersectionObserver(handleObserver, {
        root: null,
        rootMargin: '0px',
        threshold: 0.1
      });
      
      observer.current.observe(currentRef);
    }
    
    return () => {
      if (observer.current) {
        if (currentRef) {
          observer.current.unobserve(currentRef);
        }
        observer.current.disconnect(); // Полностью отключаем observer при размонтировании
      }
    };
  }, [handleObserver, hasMore]); // Важно добавить hasMore в массив зависимостей

  return (
    <>
      {children}
      
      <Box 
        ref={loadMoreRef} 
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          py: 4,
          minHeight: '100px'
        }}
      >
        {loading ? (
          <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 1 }}>
            <CircularProgress size={24} />
            <Typography variant="body2" color="text.secondary">
              {loadingMessage}
            </Typography>
          </Box>
        ) : hasMore ? (
          <Button
            variant="outlined"
            onClick={onLoadMore}
            startIcon={<ChevronDown />}
            disabled={loading}
          >
            {loadMoreButtonText}
          </Button>
        ) : (
          <Typography variant="body2" color="text.secondary">
            {noMoreItemsText}
          </Typography>
        )}
      </Box>
    </>
  );
};

export default InfiniteScroll;