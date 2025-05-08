// frontend/hostel-frontend/src/components/maps/leaflet-icons.js
import L from 'leaflet';



const iconUrl = '/bed-marker.svg';
const hotelIconUrl = '/hotel-marker.svg';
const shadowUrl = 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-shadow.png';

L.Icon.Default.mergeOptions({
  iconRetinaUrl: iconUrl,
  iconUrl: iconUrl,
  shadowUrl: shadowUrl,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  tooltipAnchor: [16, -28],
  shadowSize: [41, 41]
});

export const hotelIcon = L.icon({
  iconUrl: hotelIconUrl,
  shadowUrl: shadowUrl,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  tooltipAnchor: [16, -28],
  shadowSize: [41, 41]
});

export const customIcon = L.icon({
  iconUrl: iconUrl,
  shadowUrl: shadowUrl,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  tooltipAnchor: [16, -28],
  shadowSize: [41, 41]
});

export const createColorIcon = (color) => {
  return L.divIcon({
    className: `custom-div-icon`,
    html: `<div style="background-color:${color};width:20px;height:20px;border-radius:50%;"></div>`,
    iconSize: [20, 20],
    iconAnchor: [10, 10]
  });
};