import React, { useEffect, useRef } from 'react';

const ProtectedPhoneNumber = ({ phone, visible }) => {
  const canvasRef = useRef(null);

  useEffect(() => {
    if (visible && phone && canvasRef.current) {
      const canvas = canvasRef.current;
      const ctx = canvas.getContext('2d');
      ctx.font = '16px Arial';
      ctx.fillStyle = '#1976d2';
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      // Убираем лишние символы из номера, оставляем только цифры
      const cleanedPhone = phone.replace(/[^\d]/g, '');

      // Добавляем '+' в начало номера, если его нет
      const normalizedPhone = cleanedPhone.startsWith('381') ? `+${cleanedPhone}` : `+${cleanedPhone}`;

      // Форматируем номер телефона
      let formattedPhone = '';
      if (normalizedPhone.startsWith('+381')) {
        // Форматируем для Сербии (+381)
        formattedPhone = normalizedPhone.replace(
          /\+381(\d{2,3})(\d{3})(\d{2})(\d{2})/,
          '+381 $1 $2 $3 $4'
        );
      } else {
        // Форматируем для остальных стран
        formattedPhone = normalizedPhone.replace(
          /\+(\d{1})(\d{3})(\d{3})(\d{2})(\d{2})/,
          '+$1 $2 $3 $4 $5'
        );
      }

      ctx.fillText(formattedPhone, 10, 20);
    }
  }, [phone, visible]);

  if (!visible || !phone) return null;

  return (
    <canvas 
      ref={canvasRef}
      width={250}
      height={30}
      style={{ marginLeft: 8 }}
    />
  );
};

export default ProtectedPhoneNumber;
