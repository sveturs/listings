import React from 'react';
import {
  // Transport icons
  Car,
  Truck,
  Bike,
  Ship,
  Anchor,
  Sailboat,

  // Industry icons
  Factory,
  Tractor,
  Wheat,

  // Tech/Tools icons
  Cog,
  Wrench,
  Settings,

  // General icons
  Globe,
  Flag,
  CreditCard,
  Star,
  Clock,
  Calendar,
  Battery,
  Zap,
  Leaf,
  Crown,
  Shield,
  Music,
  Map,
  Mountain,
  Building2,

  // Default icon
  Package,
} from 'lucide-react';

// Map of string icon names to React components
const iconMap: Record<string, React.ComponentType<any>> = {
  // Cars & Transport
  car: Car,
  truck: Truck,
  bus: Truck, // Using Truck as Bus is not available
  van: Truck, // Using Truck as Van is not available
  trailer: Truck,
  motorcycle: Bike,
  scooter: Bike,
  bicycle: Bike,

  // Water transport
  ship: Ship,
  anchor: Anchor,
  sailboat: Sailboat,
  water: Ship,

  // Industry & Agriculture
  factory: Factory,
  tractor: Tractor,
  wheat: Wheat,

  // Tech & Tools
  cog: Cog,
  gear: Settings,
  tools: Wrench, // Using Wrench as Tools icon
  wrench: Wrench,
  engine: Cog,

  // General
  globe: Globe,
  flag: Flag,
  'id-card': CreditCard,
  star: Star,
  clock: Clock,
  calendar: Calendar,
  battery: Battery,
  bolt: Zap,
  leaf: Leaf,
  crown: Crown,
  shield: Shield,
  music: Music,
  map: Map,
  mountain: Mountain,
  city: Building2,

  // Special vehicles
  racing: Car,
  'car-side': Car,
  caravan: Truck,
  speed: Zap,
  vintage: Clock,
  gem: Star,
  snowflake: Star,
  golf: Flag,
  'quad-bike': Bike,
  triangle: Shield,
};

export function getCategoryIcon(
  iconName?: string
): React.ComponentType<any> | null {
  if (!iconName) return null;

  // Return the mapped icon or default Package icon
  return iconMap[iconName.toLowerCase()] || Package;
}

export function renderCategoryIcon(
  iconName?: string,
  className?: string
): React.ReactNode {
  if (!iconName) return null;

  // Проверяем, является ли это эмодзи (не латинские символы)
  const isEmoji = /[^\u0000-\u007F]/.test(iconName);
  if (isEmoji) {
    return <span className={className}>{iconName}</span>;
  }

  const IconComponent = getCategoryIcon(iconName);
  if (!IconComponent) return null;

  return <IconComponent className={className} />;
}
