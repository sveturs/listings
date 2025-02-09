import React from 'react';
import { useTranslation } from 'react-i18next';
import {
    Select,
    MenuItem,
    FormControl,
    Box,
    Typography,
    useTheme,
    useMediaQuery
} from '@mui/material';

const LANGUAGES = [
    { code: 'en', name: 'English', flag: '🇬🇧' },
    { code: 'ru', name: 'Русский', flag: '🇷🇺' },
    { code: 'sr', name: 'Српски', flag: '🇷🇸' }
];

const LanguageSwitcher = () => {
    const { i18n } = useTranslation();
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const handleLanguageChange = (event) => {
        const newLang = event.target.value;
        i18n.changeLanguage(newLang);
        localStorage.setItem('preferredLanguage', newLang);
        document.documentElement.lang = newLang;
    };

    const currentLanguage = LANGUAGES.find(lang => lang.code === i18n.language);

    return (
        <FormControl size="small">
            <Select
                value={i18n.language}
                onChange={handleLanguageChange}
                sx={{
                    minWidth: isMobile ? 'auto' : 120,
                    '& .MuiSelect-select': {
                        display: 'flex',
                        alignItems: 'center',
                        gap: 1,
                        py: isMobile ? 1 : 1.2,
                        pl: isMobile ? 1.5 : 2,
                        pr: isMobile ? 3 : 4
                    }
                }}
            >
                {LANGUAGES.map((lang) => (
                    <MenuItem 
                        key={lang.code} 
                        value={lang.code}
                        sx={{
                            minWidth: isMobile ? 'auto' : 120,
                            py: isMobile ? 1 : 1.2,
                            px: isMobile ? 1.5 : 2,
                        }}
                    >
                        <Box sx={{ 
                            display: 'flex', 
                            alignItems: 'center', 
                            gap: 1,
                            width: '100%',
                            justifyContent: isMobile ? 'center' : 'flex-start'
                        }}>
                            <Typography 
                                variant="body2" 
                                component="span"
                                sx={{ fontSize: isMobile ? '1.2rem' : '1rem' }}
                            >
                                {lang.flag}
                            </Typography>
                            {!isMobile && (
                                <Typography variant="body2">
                                    {lang.name}
                                </Typography>
                            )}
                        </Box>
                    </MenuItem>
                ))}
            </Select>
        </FormControl>
    );
};

export default LanguageSwitcher;