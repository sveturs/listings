// frontend/hostel-frontend/src/components/marketplace/InfiniteScroll.js
import React, { useCallback, useEffect, useRef, useState } from 'react';
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
  const [requestSent, setRequestSent] = useState(false);

  // Сбрасываем флаг requestSent когда изменяется hasMore или loading
  useEffect(() => {
    if (!loading && hasMore) {
      setRequestSent(false);
    }
  }, [hasMore, loading]);

  const handleObserver = useCallback(
    (entries) => {
      const [target] = entries;
      // Добавляем проверку requestSent, чтобы предотвратить повторные запросы
      if (target.isIntersecting && hasMore && !loading && autoLoad && !requestSent) {
        console.log('InfiniteScroll: Observer triggered, loading more items');
        setRequestSent(true);
        onLoadMore();
      }
    },
    [hasMore, loading, onLoadMore, autoLoad, requestSent]
  );

  // Обработчик кнопки "Показать еще"
  const handleLoadMoreClick = () => {
    if (!loading && hasMore && !requestSent) {
      console.log('InfiniteScroll: Load more button clicked');
      setRequestSent(true);
      onLoadMore();
    }
  };

  useEffect(() => {
    const currentRef = loadMoreRef.current;
    // Подключаем observer только если hasMore = true и !requestSent и !loading
    if (currentRef && hasMore && !requestSent && !loading) {
      console.log('InfiniteScroll: Setting up intersection observer, hasMore =', hasMore);
      // Очищаем предыдущий observer если он существует
      if (observer.current) {
        observer.current.disconnect();
      }
      observer.current = new IntersectionObserver(handleObserver, {
        root: null,
        rootMargin: '200px', // Увеличиваем margin для более раннего срабатывания
        threshold: 0.1
      });
      observer.current.observe(currentRef);
    } else if ((!hasMore || loading || requestSent) && observer.current) {
      console.log('InfiniteScroll: Disconnecting observer, hasMore =', hasMore, 'loading =', loading, 'requestSent =', requestSent);
      observer.current.disconnect();
    }
    return () => {
      if (observer.current) {
        observer.current.disconnect();
      }
    };
  }, [handleObserver, hasMore, requestSent, loading]);

  // Отладочное логирование
  useEffect(() => {
    console.log(`InfiniteScroll state: hasMore=${hasMore}, loading=${loading}, requestSent=${requestSent}`);
  }, [hasMore, loading, requestSent]);

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
            onClick={handleLoadMoreClick}
            startIcon={<ChevronDown />}
            disabled={loading || requestSent}
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