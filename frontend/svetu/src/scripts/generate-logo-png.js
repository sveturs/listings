const fs = require('fs');
const { createCanvas } = require('canvas');

// Create a 48x48 canvas
const width = 48;
const height = 48;
const canvas = createCanvas(width, height);
const ctx = canvas.getContext('2d');

// Clear canvas
ctx.fillStyle = 'transparent';
ctx.fillRect(0, 0, width, height);

// Calculate sizes
const scale = width / 100; // Original is 100x100
const tileSize = 25 * scale;
const gap = 2.5 * scale;
const offset = (width - (tileSize * 3 + gap * 2)) / 2;

// Tile positions
const positions = [
  { x: 0, y: 0 },
  { x: 1, y: 0 },
  { x: 2, y: 0 },
  { x: 0, y: 1 },
  { x: 1, y: 1 },
  { x: 2, y: 1 },
  { x: 0, y: 2 },
  { x: 1, y: 2 },
  { x: 2, y: 2 },
];

// Colors and icons
const colors = [
  '#2196F3',
  '#4CAF50',
  '#F44336',
  '#FF9800',
  '#673AB7',
  '#00BCD4',
  '#FFEB3B',
  '#607D8B',
  '#9C27B0',
];
const icons = ['🛒', '🏪', '🛍️', '📦', '🏠', '🤝', '📱', '💳', '💰'];

// Draw tiles
positions.forEach((pos, index) => {
  const x = offset + pos.x * (tileSize + gap);
  const y = offset + pos.y * (tileSize + gap);

  // Create gradient
  const gradient = ctx.createLinearGradient(x, y, x + tileSize, y + tileSize);
  gradient.addColorStop(0, colors[index]);
  gradient.addColorStop(1, colors[index] + 'B3'); // 70% opacity

  // Shadow
  ctx.shadowColor = 'rgba(0, 0, 0, 0.2)';
  ctx.shadowBlur = 3;
  ctx.shadowOffsetX = 2;
  ctx.shadowOffsetY = 2;

  // Draw rounded rectangle
  const radius = tileSize * 0.15;
  ctx.beginPath();
  ctx.moveTo(x + radius, y);
  ctx.lineTo(x + tileSize - radius, y);
  ctx.quadraticCurveTo(x + tileSize, y, x + tileSize, y + radius);
  ctx.lineTo(x + tileSize, y + tileSize - radius);
  ctx.quadraticCurveTo(
    x + tileSize,
    y + tileSize,
    x + tileSize - radius,
    y + tileSize
  );
  ctx.lineTo(x + radius, y + tileSize);
  ctx.quadraticCurveTo(x, y + tileSize, x, y + tileSize - radius);
  ctx.lineTo(x, y + radius);
  ctx.quadraticCurveTo(x, y, x + radius, y);
  ctx.closePath();

  ctx.fillStyle = gradient;
  ctx.fill();

  // Reset shadow for text
  ctx.shadowColor = 'transparent';
  ctx.shadowBlur = 0;
  ctx.shadowOffsetX = 0;
  ctx.shadowOffsetY = 0;

  // Draw emoji icon
  ctx.font = `${tileSize * 0.7}px Arial`;
  ctx.fillStyle = 'white';
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  ctx.fillText(icons[index], x + tileSize / 2, y + tileSize / 2);
});

// Save as PNG
const buffer = canvas.toBuffer('image/png');
fs.writeFileSync(
  '/data/hostel-booking-system/frontend/svetu/public/logos/svetu-gradient-48x48.png',
  buffer
);

console.log('Logo saved to public/logos/svetu-gradient-48x48.png');
