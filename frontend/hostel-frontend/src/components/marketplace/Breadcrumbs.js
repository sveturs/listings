import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { 
    Breadcrumbs as MuiBreadcrumbs, 
    Typography, 
    Box,
    useTheme,
    useMediaQuery 
} from '@mui/material';
import { ChevronRight } from 'lucide-react';

const Breadcrumbs = ({ paths }) => {
    const theme = useTheme();
    const navigate = useNavigate();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    console.log('Breadcrumbs received paths:', paths);

    if (!paths || paths.length === 0) {
        return null;
    }

    const handleCategoryClick = (categoryId, event) => {
        event.preventDefault();
        navigate(`/marketplace?category_id=${categoryId}`);
    };

    return (
        <Box sx={{ 
            mb: 2,
            py: 1,
            overflowX: 'auto',
            whiteSpace: 'nowrap',
            '&::-webkit-scrollbar': {
                display: 'none'
            },
            scrollbarWidth: 'none'
        }}>
            <MuiBreadcrumbs 
                separator={<ChevronRight size={16} />}
                aria-label="breadcrumb"
                sx={{
                    '& .MuiBreadcrumbs-ol': {
                        flexWrap: 'nowrap'
                    },
                    '& .MuiBreadcrumbs-li': {
                        display: 'flex',
                        alignItems: 'center'
                    }
                }}
            >
                <Link 
                    to="/marketplace"
                    style={{ 
                        color: theme.palette.text.secondary,
                        textDecoration: 'none',
                        padding: '4px 8px',
                        borderRadius: '4px',
                        transition: 'all 0.2s'
                    }}
                >
                    Главная
                </Link>

                {paths.map((path, index) => {
                    const isLast = index === paths.length - 1;
                    
                    if (isLast) {
                        return (
                            <Typography 
                                key={path.id} 
                                color="text.primary"
                                sx={{ 
                                    fontSize: '0.875rem',
                                    padding: '4px 8px',
                                    fontWeight: 500
                                }}
                            >
                                {path.name}
                            </Typography>
                        );
                    }

                    return (
                        <Link
                            key={path.id}
                            to={`/marketplace?category_id=${path.id}`}
                            onClick={(e) => handleCategoryClick(path.id, e)}
                            style={{ 
                                color: theme.palette.text.secondary,
                                textDecoration: 'none',
                                padding: '4px 8px',
                                borderRadius: '4px',
                                fontSize: '0.875rem',
                                transition: 'all 0.2s',
                                '&:hover': {
                                    backgroundColor: theme.palette.action.hover,
                                    color: theme.palette.primary.main
                                }
                            }}
                        >
                            {path.name}
                        </Link>
                    );
                })}
            </MuiBreadcrumbs>
        </Box>
    );
};

export default Breadcrumbs;