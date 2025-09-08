-- Удаление добавленных ключевых слов для DJ оборудования

DELETE FROM category_keywords 
WHERE source = 'manual' 
  AND category_id IN (10191, 1103)
  AND keyword IN (
    -- English
    'dj', 'dj equipment', 'dj controller', 'dj mixer', 'turntable', 'cdj',
    'dj case', 'flight case', 'equipment case', 'audio', 'audio equipment',
    'sound system', 'speakers', 'amplifier', 'headphones', 'microphone',
    'video', 'music',
    -- Russian
    'диджей', 'dj оборудование', 'диджейское оборудование', 'dj контроллер',
    'dj микшер', 'dj пульт', 'кейс', 'кейс для оборудования', 
    'транспортировочный кейс', 'аудио', 'аудио оборудование',
    'звуковое оборудование', 'колонки', 'акустика', 'усилитель',
    'наушники', 'микрофон', 'видео', 'музыка', 'музыкальный',
    -- Serbian
    'dj oprema', 'dj kontroler', 'dj mikser', 'gramofon', 'kofer za opremu',
    'audio oprema', 'zvučnici', 'pojačalo', 'slušalice', 'mikrofon',
    'muzika', 'muzički'
  );