// src/components/marketplace/Breadcrumbs.js
import React from 'react';
import { Link } from 'react-router-dom';
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
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));


    if (!paths || paths.length === 0) return null;

    return (
        <Box sx={{ 
            mb: 2,
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
                    }
                }}
            >
                <Link 
                    to="/marketplace"
                    style={{ 
                        color: 'inherit', 
                        textDecoration: 'none',
                        '&:hover': {
                            textDecoration: 'underline'
                        }
                    }}
                >
                    Главная
                </Link>
                {paths.map((path, index) => {
                    const isLast = index === paths.length - 1;
                    const to = `/marketplace?category=${path.id}`;
                    
                    return isLast ? (
                        <Typography 
                            key={path.id} 
                            color="text.primary"
                            sx={{ fontSize: '0.875rem' }}
                        >
                            {path.name}
                        </Typography>
                    ) : (
                        <Link
                            key={path.id}
                            to={to}
                            style={{ 
                                color: 'inherit',
                                textDecoration: 'none',
                                '&:hover': {
                                    textDecoration: 'underline'
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