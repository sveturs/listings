import React, { useState, useEffect, useRef } from 'react';
import { Box, Paper } from '@mui/material';
import { keyframes } from '@mui/system';

const fadeOut = keyframes`
  0% { opacity: 1; }
  90% { opacity: 1; }
  100% { opacity: 0; }
`;

const PhonePopup = ({ phone, visible, onClose }) => {
    const canvasRef = useRef(null);
  
    useEffect(() => {
        if (visible && phone && canvasRef.current) {
          const canvas = canvasRef.current;
          const ctx = canvas.getContext('2d');
          ctx.font = '16px Arial';
          ctx.fillStyle = '#000000';
          
          const cleanedPhone = phone.replace(/[^\d]/g, '');
          const formattedPhone = cleanedPhone.startsWith('381') 
            ? cleanedPhone.replace(/(\d{3})(\d{2})(\d{3})(\d{2})(\d{2})/, '+$1 $2 $3 $4 $5')
            : cleanedPhone.replace(/(\d{1})(\d{3})(\d{3})(\d{2})(\d{2})/, '+$1 $2 $3 $4 $5');
          
          // Измеряем ширину текста
          const textWidth = ctx.measureText(formattedPhone).width;
          // Устанавливаем ширину canvas с небольшим отступом
          canvas.width = textWidth + 20;
          
          // Перерисовываем после изменения размера
          ctx.font = '16px Arial';
          ctx.fillStyle = '#000000';
          ctx.fillText(formattedPhone, 10, 20);
        }
      }, [phone, visible]);
  
    useEffect(() => {
      if (visible) {
        const timer = setTimeout(onClose, 10000);
        return () => clearTimeout(timer);
      }
    }, [visible, onClose]);
  
    if (!visible || !phone) return null;
  
    return (
      <Box
        sx={{
          position: 'absolute',
          bottom: '100%',
          left: '50%',
          transform: 'translateX(-50%)',
          mb: 1,
          animation: `${fadeOut} 10s ease-in forwards`
        }}
      >
        <Paper
          elevation={3}
          sx={{
            p: 1.5,
            bgcolor: 'background.paper',
            borderRadius: 1,
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
        </Paper>
      </Box>
    );
  };

export default PhonePopup;