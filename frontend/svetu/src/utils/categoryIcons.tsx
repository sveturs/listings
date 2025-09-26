import { FaCar, FaTshirt } from 'react-icons/fa';
import {
  BsHouseDoor,
  BsLaptop,
  BsBriefcase,
  BsTools,
  BsPalette,
  BsHandbag,
  BsPhone,
  BsGem,
} from 'react-icons/bs';

// Маппинг имен иконок на компоненты
const iconComponents: { [key: string]: any } = {
  BsHouseDoor,
  FaCar,
  BsLaptop,
  FaTshirt,
  BsBriefcase,
  BsTools,
  BsPalette,
  BsHandbag,
  BsPhone,
  BsGem,
};

// Функция для получения компонента иконки по имени
export const getCategoryIcon = (iconName?: string) => {
  if (!iconName) return BsHandbag;
  return iconComponents[iconName] || BsHandbag;
};
