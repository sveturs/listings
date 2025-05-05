import React, { useState, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Box,
  TextField,
  InputAdornment,
  IconButton,
  Paper,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Typography,
  Chip,
  CircularProgress,
  Divider,
  Button,
  useTheme,
  useMediaQuery
} from '@mui/material';
import {
  Search as SearchIcon,
  Clear as ClearIcon,
  History as HistoryIcon,
  LocationOn as LocationIcon,
  Category as CategoryIcon,
  KeyboardReturnOutlined as ReturnIcon,
  ListAlt,
  ListAlt as ListingIcon
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';
import axios from '../../api/axios';

const SearchContainer = styled(Box)(({ theme, drawerOpen, drawerWidth }) => ({
  position: 'absolute',
  zIndex: 1000,
  top: theme.spacing(3), // Увеличение отступа сверху
  left: '50%',
  transform: 'translateX(-50%)',
  width: '90%',
  maxWidth: 500,
  pointerEvents: 'auto',
  [theme.breakpoints.down('sm')]: {
    top: theme.spacing(1),
    width: '95%',
  }
}));

const SearchField = styled(TextField)(({ theme }) => ({
  '& .MuiOutlinedInput-root': {
    borderRadius: 24,
    backgroundColor: theme.palette.background.paper,
    '&.Mui-focused': {
      boxShadow: '0 2px 12px rgba(0, 0, 0, 0.1)'
    }
  }
}));

const SearchResults = styled(Paper)(({ theme }) => ({
  marginTop: theme.spacing(1),
  maxHeight: 400,
  overflow: 'auto',
  boxShadow: '0 4px 20px rgba(0, 0, 0, 0.15)',
  position: 'absolute',
  width: '100%',
  zIndex: 1010, // Higher than the search box
  backgroundColor: theme.palette.background.paper, // White background 
  [theme.breakpoints.down('sm')]: {
    maxHeight: 300
  }
}));

const RecentSearches = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexWrap: 'wrap',
  gap: theme.spacing(1),
  padding: theme.spacing(2)
}));

const GISSearchPanel = ({ onSearch, onLocationSelect, drawerOpen = false, drawerWidth = 400 }) => {
  const { t } = useTranslation('gis');
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const [query, setQuery] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [recentSearches, setRecentSearches] = useState([]);
  const [showResults, setShowResults] = useState(false);
  const [loading, setLoading] = useState(false);
  const [searchType, setSearchType] = useState('');
  const inputRef = useRef(null);
  const debounceRef = useRef(null);

  useEffect(() => {
    // Load recent searches from localStorage
    const storedSearches = localStorage.getItem('gisRecentSearches');
    if (storedSearches) {
      try {
        setRecentSearches(JSON.parse(storedSearches).slice(0, 5));
      } catch (e) {
        console.error('Error parsing recent searches:', e);
      }
    }
  }, []);

  useEffect(() => {
    if (!query) {
      setSearchResults([]);
      return;
    }

    // Show loading state immediately for better UX
    if (query.trim().length >= 2) {
      setLoading(true);
    }

    // Improved debouncing with reference
    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }

    debounceRef.current = setTimeout(() => {
      if (query.trim().length >= 2) {
        performSearch();
      } else {
        setLoading(false);
      }
    }, 300); // Reduced debounce time for more responsive results

    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, [query]);
  
  // Focus the search field on component mount
  useEffect(() => {
    // Focus the search input on mount for better UX
    const searchInput = document.querySelector('input[type="text"]');
    if (searchInput && !isMobile) {
      setTimeout(() => {
        searchInput.focus();
      }, 500);
    }
  }, []);
  
  // Add click-outside handler to hide search results
  useEffect(() => {
    const handleClickOutside = (e) => {
      if (inputRef.current && !inputRef.current.contains(e.target)) {
        setShowResults(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  const performSearch = async () => {
    if (!query.trim() || query.trim().length < 2) return;
    
    setLoading(true);
    try {
      // Try enhanced suggestions API first
      try {
        // First attempt - enhanced suggestions API with rich results
        const enhancedResponse = await axios.get('/api/v1/marketplace/enhanced-suggestions', {
          params: { q: query.trim(), size: 8 }
        });
        
        if (enhancedResponse.data && enhancedResponse.data.data) {
          console.log('Using enhanced suggestions API, query:', query.trim(), 'results:', enhancedResponse.data.data);
          
          // Map the enhanced suggestions to our format with appropriate icons
          const enhancedResults = enhancedResponse.data.data.map(item => {
            // Default values
            let icon = <SearchIcon color="action" />;
            let secondaryText = '';
            let borderColor = '#808080'; // Default gray
            
            // Type-specific processing
            if (item.type === 'product' || item.type === 'listing') {
              icon = <ListingIcon color="action" />;
              borderColor = '#4682B4'; // Blue for listings
            } else if (item.type === 'category') {
              icon = <CategoryIcon color="secondary" />;
              secondaryText = item.path ? item.path.map(cat => cat.name).join(' > ') : '';
              borderColor = '#6BAB33'; // Green for categories
            } else if (item.type === 'attribute') {
              icon = <ListAlt color="action" />;
              secondaryText = item.attribute_name ? `${item.attribute_name}: ${item.attribute_value || item.title}` : '';
              borderColor = '#FFA500'; // Orange for attributes
            } else if (item.type === 'address') {
              icon = <LocationIcon color="primary" />;
              borderColor = '#FF5722'; // Orange-red for addresses
            }
            
            return {
              ...item,
              icon,
              secondaryText: item.secondaryText || secondaryText,
              borderColor,
              display: item.display || item.title
            };
          });
          
          setSearchResults(enhancedResults);
          setSearchType(enhancedResults.length > 0 ? enhancedResults[0].type : '');
          setLoading(false);
          return;
        }
      } catch (err) {
        console.log('Enhanced suggestions API not available, falling back to standard search');
        // Continue with fallback searches
      }

      // Fallback searches if enhanced API fails
      // 1. Search addresses
      const addressResponse = await axios.get('/api/v1/geocode/get-city-suggestions', {
        params: { q: query.trim(), limit: 5 }
      });

      // 2. Search listings
      const listingsResponse = await axios.get('/api/v1/marketplace/search', {
        params: { query: query.trim(), size: 5 }
      });

      // 3. Search categories
      const categoriesResponse = await axios.get('/api/v1/marketplace/category-suggestions', {
        params: { q: query.trim() }
      });

      // Process addresses with priority
      const addressResults = (addressResponse.data?.data || []).map(item => ({
        ...item,
        type: 'address',
        title: item.city || item.address,
        icon: <LocationIcon color="primary" />,
        display: item.city || item.address,
        priority: 1, // Higher priority for addresses
        secondaryText: item.country || '',
        borderColor: '#FF5722' // Orange-red for addresses
      }));

      // Process listings with medium priority
      const listingResults = (listingsResponse.data?.data || []).map(item => ({
        ...item,
        type: 'listing',
        title: item.title || item.name,
        icon: <ListingIcon color="action" />,
        display: item.title || item.name,
        priority: 2, // Medium priority for listings
        secondaryText: item.price ? `${t('card.price')}: ${item.price} RSD` : '',
        borderColor: '#4682B4' // Blue for listings
      }));

      // Process categories with lower priority
      const categoryResults = (categoriesResponse.data?.data || []).map(item => ({
        ...item,
        type: 'category',
        title: item.name,
        icon: <CategoryIcon color="secondary" />,
        display: item.name,
        priority: 3, // Lower priority for categories
        borderColor: '#6BAB33', // Green for categories
        secondaryText: item.parent_name || ''
      }));

      // Combine all results and sort by priority
      const combinedResults = [
        ...addressResults,
        ...listingResults.slice(0, 3), // Limit to 3 listings
        ...categoryResults.slice(0, 3)  // Limit to 3 categories
      ].sort((a, b) => a.priority - b.priority);

      setSearchResults(combinedResults);
      setSearchType(combinedResults.length > 0 ? combinedResults[0].type : '');
    } catch (error) {
      console.error('Search error:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = () => {
    if (!query.trim()) return;
    
    // Add to recent searches
    const updatedSearches = [
      { query: query.trim(), timestamp: Date.now() },
      ...recentSearches.filter(item => item.query !== query.trim())
    ].slice(0, 5);
    
    setRecentSearches(updatedSearches);
    localStorage.setItem('gisRecentSearches', JSON.stringify(updatedSearches));
    
    // Perform search
    if (onSearch) {
      onSearch(query.trim());
    }
    
    // Hide results
    setShowResults(false);
  };

  const handleResultClick = (result) => {
    // Add to recent searches - use appropriate field based on result type
    let searchText = '';
    switch (result.type) {
      case 'address':
        searchText = result.address || result.city || result.title;
        break;
      case 'category':
        searchText = result.name || result.title;
        break;
      case 'attribute':
        searchText = result.attribute_value || result.title;
        break;
      case 'product':
      case 'listing':
      default:
        searchText = result.title || result.name || result.display;
        break;
    }
    
    const updatedSearches = [
      { query: searchText, timestamp: Date.now() },
      ...recentSearches.filter(item => item.query !== searchText)
    ].slice(0, 5);
    
    setRecentSearches(updatedSearches);
    localStorage.setItem('gisRecentSearches', JSON.stringify(updatedSearches));
    
    // Handle different result types
    if (result.type === 'address' && onLocationSelect) {
      // For address - center map on location coordinates
      onLocationSelect({
        latitude: result.latitude || result.lat,
        longitude: result.longitude || result.lon,
        address: result.address || result.city || result.title
      });
      setQuery(searchText);
    } 
    else if ((result.type === 'product' || result.type === 'listing') && result.id) {
      // For products/listings with ID - direct to listing details or search by title
      if (onSearch) {
        onSearch(searchText, { listingId: result.id });
        setQuery(searchText);
      }
      
      // If we have direct listing navigation from GIS interface in the future:
      // window.location.href = `/marketplace/listings/${result.id}`;
    } 
    else if (result.type === 'category' && result.id && onSearch) {
      // For categories - search with category filter
      onSearch('', { categoryId: result.id });
      setQuery(searchText);
    }
    else if (result.type === 'attribute' && onSearch) {
      // For attributes - search with attribute filter
      setQuery(searchText);
      
      // Set up attribute filters if we have the attribute name
      if (result.attribute_name) {
        onSearch(searchText, { 
          attributeFilters: { [result.attribute_name]: result.attribute_value || result.title }
        });
      } else {
        onSearch(searchText);
      }
    }
    else if (onSearch) {
      // Generic fallback for other result types
      onSearch(searchText);
      setQuery(searchText);
    }
    
    // Hide results
    setShowResults(false);
  };

  const handleRecentSearchClick = (searchItem) => {
    setQuery(searchItem.query);
    if (onSearch) {
      onSearch(searchItem.query);
    }
    setShowResults(false);
  };

  const handleClearSearch = () => {
    setQuery('');
    setSearchResults([]);
    if (onSearch) {
      onSearch('');
    }
  };

  const handleInputFocus = () => {
    setShowResults(true);
  };

  const handleKeyDown = (e) => {
    if (e.key === 'Enter') {
      handleSearch();
    }
  };

  return (
    <SearchContainer ref={inputRef} drawerOpen={drawerOpen} drawerWidth={drawerWidth}>
      <SearchField
        fullWidth
        placeholder={t('search.placeholder')}
        variant="outlined"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onFocus={handleInputFocus}
        onKeyDown={handleKeyDown}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <SearchIcon color="action" />
            </InputAdornment>
          ),
          endAdornment: (
            <InputAdornment position="end">
              {loading ? (
                <CircularProgress size={20} />
              ) : query ? (
                <IconButton
                  size="small"
                  onClick={handleClearSearch}
                  edge="end"
                  aria-label="clear search"
                >
                  <ClearIcon fontSize="small" />
                </IconButton>
              ) : null}
            </InputAdornment>
          ),
        }}
      />

      {showResults && (
        <SearchResults elevation={4} sx={{ boxShadow: '0 4px 12px rgba(0, 0, 0, 0.08)', borderRadius: '8px' }}>
          {query.trim() === '' && recentSearches.length > 0 && (
            <>
              <Box px={2} py={1}>
                <Typography variant="subtitle2" color="textSecondary">
                  {t('search.recent')}
                </Typography>
              </Box>
              <RecentSearches>
                {recentSearches.map((item, index) => (
                  <Chip
                    key={index}
                    icon={<HistoryIcon fontSize="small" />}
                    label={item.query}
                    onClick={() => handleRecentSearchClick(item)}
                    variant="outlined"
                    size="small"
                  />
                ))}
              </RecentSearches>
              <Divider />
            </>
          )}

          {query.trim() !== '' && (
            <>
              <Box px={2} py={1} display="flex" justifyContent="space-between" alignItems="center">
                <Typography variant="subtitle2" color="textSecondary">
                  {searchResults.length > 0 
                    ? t('search.suggestions') 
                    : t('search.noResults')}
                </Typography>
                <Button
                  size="small"
                  endIcon={<ReturnIcon />}
                  onClick={handleSearch}
                >
                  {t('search.keywordSearch')}
                </Button>
              </Box>

              {searchResults.length > 0 && (
                <List disablePadding>
                  {searchType !== '' && (
                    <ListItem sx={{ bgcolor: 'action.hover', py: 0.5 }}>
                      <ListItemText 
                        primary={
                          <Typography variant="caption" color="textSecondary">
                            {searchType === 'address' 
                              ? t('search.addressSearch') 
                              : searchType === 'category' 
                                ? t('search.categorySearch') 
                                : t('search.keywordSearch')}
                          </Typography>
                        } 
                      />
                    </ListItem>
                  )}
                  
                  {searchResults.map((result, index) => (
                    <React.Fragment key={index}>
                      {index > 0 && result.type !== searchResults[index - 1].type && (
                        <ListItem sx={{ bgcolor: 'action.hover', py: 0.5 }}>
                          <ListItemText 
                            primary={
                              <Typography variant="caption" color="textSecondary">
                                {result.type === 'address' 
                                  ? t('search.addressSearch') 
                                  : result.type === 'category' 
                                    ? t('search.categorySearch') 
                                    : t('search.keywordSearch')}
                              </Typography>
                            } 
                          />
                        </ListItem>
                      )}
                      
                      <ListItem
                        button
                        onClick={() => handleResultClick(result)}
                        sx={{ 
                          padding: '10px 16px',
                          borderBottom: index < searchResults.length - 1 ? '1px solid #f0f0f0' : 'none',
                          position: 'relative',
                          '&::before': {
                            content: '""',
                            position: 'absolute',
                            left: 0,
                            top: '20%',
                            height: '60%',
                            width: '4px',
                            backgroundColor: result.borderColor || 
                              (result.type === 'address' ? '#FF5722' : // Orange for addresses
                               result.type === 'category' ? '#6BAB33' : // Green for categories
                               result.type === 'attribute' ? '#FFA500' : // Orange for attributes
                               '#4682B4'), // Blue for listings/products
                            borderRadius: '0 4px 4px 0'
                          },
                          '&:hover': { 
                            backgroundColor: 'rgba(0, 0, 0, 0.03)'
                          },
                          transition: 'background-color 0.2s ease'
                        }}
                      >
                        <ListItemIcon sx={{ minWidth: 36 }}>
                          {result.icon}
                        </ListItemIcon>
                        <ListItemText 
                          primary={
                            <Typography
                              variant="body2"
                              sx={{
                                fontWeight: result.type === 'listing' || result.type === 'product' ? 500 : 400,
                                color: 'text.primary',
                                fontSize: '0.95rem',
                                lineHeight: '1.25',
                                overflow: 'hidden',
                                textOverflow: 'ellipsis',
                                whiteSpace: 'nowrap'
                              }}
                            >
                              {result.title || 
                               (result.type === 'address' ? result.address : 
                                result.type === 'category' ? result.name : 
                                result.display || '')}
                            </Typography>
                          }
                          secondary={
                            result.secondaryText ? (
                              <Typography
                                variant="caption"
                                sx={{
                                  color: 'text.secondary',
                                  fontSize: '0.8rem',
                                  display: 'block',
                                  maxWidth: '100%',
                                  overflow: 'hidden',
                                  textOverflow: 'ellipsis',
                                  whiteSpace: 'nowrap'
                                }}
                              >
                                {result.secondaryText}
                              </Typography>
                            ) : null
                          }
                          primaryTypographyProps={{ noWrap: true }}
                          secondaryTypographyProps={{ noWrap: true }}
                        />
                      </ListItem>
                    </React.Fragment>
                  ))}
                </List>
              )}
            </>
          )}
        </SearchResults>
      )}
    </SearchContainer>
  );
};

export default GISSearchPanel;