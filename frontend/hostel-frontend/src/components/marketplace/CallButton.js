import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Button, Box, CircularProgress } from '@mui/material';
import { Phone } from 'lucide-react';

const PhonePopup = ({ phone, visible, onClose }) => {
 const canvasRef = React.useRef(null);

 React.useEffect(() => {
   if (visible && phone && canvasRef.current) {
     const canvas = canvasRef.current;
     const ctx = canvas.getContext('2d');
     ctx.font = '16px Arial';
     ctx.fillStyle = '#000000';
     
     const cleanedPhone = phone.replace(/[^\d]/g, '');
     const formattedPhone = cleanedPhone.startsWith('381')
       ? cleanedPhone.replace(/(\d{3})(\d{2})(\d{3})(\d{2})(\d{2})/, '+$1 $2 $3 $4 $5')
       : cleanedPhone.replace(/(\d{1})(\d{3})(\d{3})(\d{2})(\d{2})/, '+$1 $2 $3 $4 $5');
     
     const textWidth = ctx.measureText(formattedPhone).width;
     canvas.width = textWidth + 20;
     
     ctx.font = '16px Arial';
     ctx.fillStyle = '#000000';
     ctx.fillText(formattedPhone, 10, 20);
   }
 }, [phone, visible]);

 if (!visible || !phone) return null;

 return (
   <Box
     sx={{
       position: 'absolute',
       bottom: '100%',
       left: '50%',
       transform: 'translateX(-50%)',
       mb: 1,
       bgcolor: 'background.paper',
       borderRadius: 1,
       p: 1.5,
       boxShadow: 3,
       '&:after': {
         content: '""',
         position: 'absolute',
         top: '100%',
         left: '50%',
         transform: 'translateX(-50%)',
         border: '8px solid transparent',
         borderTopColor: 'background.paper'
       }
     }}
   >
     <canvas
       ref={canvasRef}
       width={200}
       height={30}
     />
   </Box>
 );
};

const CallButton = ({ phone, isMobile }) => {
 const { t } = useTranslation('marketplace');
 const [showPhone, setShowPhone] = useState(false);
 const [loading, setLoading] = useState(false);

 const formatPhone = (phone) => {
   const cleaned = phone.replace(/[^\d]/g, '');
   return cleaned.startsWith('381') 
     ? cleaned.replace(/(\d{3})(\d{2})(\d{3})(\d{2})(\d{2})/, '+$1 $2 $3 $4 $5')
     : cleaned.replace(/(\d{1})(\d{3})(\d{3})(\d{2})(\d{2})/, '+$1 $2 $3 $4 $5');
 };

 const handleMouseEnter = () => !isMobile && setShowPhone(true);
 const handleMouseLeave = () => !isMobile && setShowPhone(false);

 const handleCallClick = () => {
   if (!phone) {
     alert(t('listings.details.contact.errors.noPhone'));
     return;
   }

   if (isMobile) {
     setShowPhone(true);
     setLoading(true);
     setTimeout(() => {
       window.location.href = `tel:${phone}`;
       setShowPhone(false);
       setLoading(false);
     }, 1500);
   } else {
     window.location.href = `tel:${phone}`;
   }
 };

 return (
   <Box sx={{ position: 'relative', width: '100%' }}>
     <Button
       variant="contained"
       fullWidth
       startIcon={<Phone size={20} />}
       onClick={handleCallClick}
       onMouseEnter={handleMouseEnter}
       onMouseLeave={handleMouseLeave}
       sx={{ 
         height: 48,
         backgroundColor: 'success.main',
         '&:hover': { backgroundColor: 'success.dark' }
       }}
     >
       {isMobile && showPhone ? (
         <Box sx={{ 
           display: 'flex', 
           alignItems: 'center', 
           gap: 1,
           fontSize: '0.875rem',
           whiteSpace: 'nowrap'
         }}>
           {formatPhone(phone)}
           {loading && (
             <CircularProgress
               size={16}
               sx={{ color: 'common.white' }}
             />
           )}
         </Box>
       ) : (
         t('listings.details.contact.call')
       )}
     </Button>

     {!isMobile && <PhonePopup phone={phone} visible={showPhone} onClose={() => setShowPhone(false)} />}
   </Box>
 );
};

export default CallButton;