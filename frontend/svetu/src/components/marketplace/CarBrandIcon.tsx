import React from 'react';

interface CarBrandIconProps {
  brand: string;
  className?: string;
}

const CarBrandIcon: React.FC<CarBrandIconProps> = ({
  brand,
  className = 'w-16 h-16',
}) => {
  const normalizedBrand = brand.toLowerCase().replace(/[^a-z0-9]/g, '');

  const brandIcons: { [key: string]: React.ReactElement } = {
    volkswagen: (
      <svg viewBox="0 0 100 100" className={className} fill="currentColor">
        <circle cx="50" cy="50" r="45" fill="#1e3a5f" />
        <path d="M50 15 L35 65 H45 L50 30 L55 65 H65 Z" fill="white" />
        <path
          d="M30 40 L45 80 H55 L70 40 M40 50 H60"
          fill="white"
          stroke="white"
          strokeWidth="3"
        />
      </svg>
    ),
    mercedesbenz: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#333333" />
        <circle
          cx="50"
          cy="50"
          r="40"
          fill="none"
          stroke="white"
          strokeWidth="2"
        />
        <path
          d="M50 10 L50 50 L20 80 L50 50 L80 80 Z"
          fill="none"
          stroke="white"
          strokeWidth="3"
        />
      </svg>
    ),
    bmw: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#0066b1" />
        <circle cx="50" cy="50" r="40" fill="white" />
        <path d="M50 10 A40 40 0 0 1 90 50 L50 50 Z" fill="#0066b1" />
        <path d="M50 50 A40 40 0 0 1 10 50 L50 50 Z" fill="#0066b1" />
        <circle
          cx="50"
          cy="50"
          r="35"
          fill="none"
          stroke="#333"
          strokeWidth="2"
        />
      </svg>
    ),
    audi: (
      <svg viewBox="0 0 100 100" className={className}>
        <g fill="none" stroke="#333333" strokeWidth="4">
          <circle cx="20" cy="50" r="12" />
          <circle cx="35" cy="50" r="12" />
          <circle cx="50" cy="50" r="12" />
          <circle cx="65" cy="50" r="12" />
        </g>
      </svg>
    ),
    toyota: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse
          cx="50"
          cy="50"
          rx="45"
          ry="30"
          fill="none"
          stroke="#eb0a1e"
          strokeWidth="3"
        />
        <ellipse
          cx="50"
          cy="50"
          rx="30"
          ry="20"
          fill="none"
          stroke="#eb0a1e"
          strokeWidth="3"
        />
        <ellipse
          cx="50"
          cy="50"
          rx="15"
          ry="40"
          fill="none"
          stroke="#eb0a1e"
          strokeWidth="3"
        />
      </svg>
    ),
    ford: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="35" fill="#003478" />
        <text
          x="50"
          y="60"
          fontFamily="Arial Black"
          fontSize="24"
          fill="white"
          textAnchor="middle"
        >
          FORD
        </text>
      </svg>
    ),
    honda: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="30" width="80" height="40" fill="#cc0000" />
        <text
          x="50"
          y="55"
          fontFamily="Arial Black"
          fontSize="20"
          fill="white"
          textAnchor="middle"
        >
          HONDA
        </text>
      </svg>
    ),
    nissan: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#c3002f" />
        <rect x="20" y="45" width="60" height="10" fill="white" />
        <text
          x="50"
          y="40"
          fontFamily="Arial"
          fontSize="16"
          fill="white"
          textAnchor="middle"
        >
          NISSAN
        </text>
      </svg>
    ),
    mazda: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#0066cc" />
        <path
          d="M50 20 C30 30, 30 70, 50 80 C70 70, 70 30, 50 20"
          fill="none"
          stroke="white"
          strokeWidth="3"
        />
        <path d="M50 35 L35 50 L50 65 L65 50 Z" fill="white" />
      </svg>
    ),
    hyundai: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="30" fill="#002c5f" />
        <text
          x="50"
          y="55"
          fontFamily="Arial"
          fontSize="14"
          fontStyle="italic"
          fill="white"
          textAnchor="middle"
        >
          HYUNDAI
        </text>
      </svg>
    ),
    kia: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="35" fill="#bb162b" />
        <text
          x="50"
          y="60"
          fontFamily="Arial Black"
          fontSize="28"
          fill="white"
          textAnchor="middle"
        >
          KIA
        </text>
      </svg>
    ),
    peugeot: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="10" width="80" height="80" fill="#000000" />
        <path
          d="M50 25 C45 25, 40 30, 40 40 L40 60 L60 60 L60 40 C60 30, 55 25, 50 25 M45 35 L55 35 L55 50 L45 50 Z"
          fill="white"
        />
        <path d="M50 20 L45 25 L55 25 Z M35 70 L50 75 L65 70" fill="white" />
      </svg>
    ),
    renault: (
      <svg viewBox="0 0 100 100" className={className}>
        <path d="M50 10 L20 50 L50 90 L80 50 Z" fill="#ffcc00" />
        <path d="M50 20 L30 50 L50 80 L70 50 Z" fill="#000000" />
        <path d="M50 30 L40 50 L50 70 L60 50 Z" fill="#ffcc00" />
      </svg>
    ),
    fiat: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#96172e" />
        <text
          x="50"
          y="60"
          fontFamily="Arial Black"
          fontSize="24"
          fill="white"
          textAnchor="middle"
        >
          FIAT
        </text>
      </svg>
    ),
    opel: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#f7c901" />
        <path
          d="M25 50 L75 50 M50 50 L15 30 M50 50 L15 70 M50 50 L85 30 M50 50 L85 70"
          stroke="#000"
          strokeWidth="5"
        />
        <circle cx="50" cy="50" r="8" fill="#000" />
      </svg>
    ),
    skoda: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#4ba82e" />
        <path
          d="M50 20 L40 50 C40 50, 45 45, 50 55 C55 45, 60 50, 60 50 L50 20"
          fill="white"
        />
        <circle cx="50" cy="65" r="5" fill="white" />
      </svg>
    ),
    citroen: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M30 40 L50 25 L70 40 M30 55 L50 40 L70 55"
          fill="none"
          stroke="#ed1c24"
          strokeWidth="5"
        />
      </svg>
    ),
    volvo: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#003057" />
        <circle
          cx="50"
          cy="50"
          r="40"
          fill="none"
          stroke="#9ba8b4"
          strokeWidth="3"
        />
        <path
          d="M25 40 L75 40 L65 25"
          fill="none"
          stroke="#9ba8b4"
          strokeWidth="3"
        />
        <text
          x="50"
          y="65"
          fontFamily="Arial"
          fontSize="16"
          fill="#9ba8b4"
          textAnchor="middle"
        >
          VOLVO
        </text>
      </svg>
    ),
    seat: (
      <svg viewBox="0 0 100 100" className={className}>
        <path d="M50 20 L30 45 L30 55 L50 35 L70 55 L70 45 Z" fill="#ea1d2c" />
        <text
          x="50"
          y="75"
          fontFamily="Arial Black"
          fontSize="20"
          fill="#ea1d2c"
          textAnchor="middle"
        >
          SEAT
        </text>
      </svg>
    ),
    suzuki: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M20 40 C30 20, 70 20, 80 40 L80 50 C70 30, 30 30, 20 50 Z"
          fill="#e30613"
        />
        <text
          x="50"
          y="70"
          fontFamily="Arial"
          fontSize="18"
          fill="#e30613"
          textAnchor="middle"
        >
          SUZUKI
        </text>
      </svg>
    ),
    mitsubishi: (
      <svg viewBox="0 0 100 100" className={className}>
        <path d="M50 20 L35 50 L50 40 L65 50 Z" fill="#e60012" />
        <path d="M35 50 L20 80 L35 70 L50 80 Z" fill="#e60012" />
        <path d="M65 50 L50 80 L65 70 L80 80 Z" fill="#e60012" />
      </svg>
    ),
    chevrolet: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="40" width="80" height="20" fill="#ffc20e" />
        <path
          d="M35 40 L45 40 L50 50 L45 60 L35 60 M65 40 L55 40 L50 50 L55 60 L65 60"
          fill="#333"
        />
      </svg>
    ),
    jeep: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="35" width="80" height="30" fill="#374743" />
        <text
          x="50"
          y="55"
          fontFamily="Arial Black"
          fontSize="22"
          fill="white"
          textAnchor="middle"
        >
          JEEP
        </text>
      </svg>
    ),
    lexus: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="35" fill="#1a1a1a" />
        <text
          x="50"
          y="58"
          fontFamily="Arial"
          fontSize="20"
          fill="white"
          textAnchor="middle"
        >
          LEXUS
        </text>
      </svg>
    ),
    porsche: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M50 15 L30 35 L30 65 L50 85 L70 65 L70 35 Z"
          fill="#000"
          stroke="#b12b28"
          strokeWidth="3"
        />
        <path d="M50 30 L40 40 L40 60 L50 70 L60 60 L60 40 Z" fill="#b12b28" />
        <text
          x="50"
          y="52"
          fontFamily="Arial"
          fontSize="8"
          fill="white"
          textAnchor="middle"
        >
          PORSCHE
        </text>
      </svg>
    ),
    ferrari: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="30" width="80" height="40" fill="#fff200" />
        <path
          d="M20 30 L20 70 L30 70 L30 55 L45 55 L45 70 L55 70 L55 30 Z"
          fill="#008c45"
        />
        <path
          d="M60 30 L60 70 L70 70 L70 55 L80 55 L80 70 L90 70 L90 30 Z"
          fill="#cd212a"
        />
        <path
          d="M40 40 C42 35, 58 35, 60 40"
          fill="none"
          stroke="#000"
          strokeWidth="2"
        />
      </svg>
    ),
    lamborghini: (
      <svg viewBox="0 0 100 100" className={className}>
        <path d="M50 20 L30 40 L30 60 L50 80 L70 60 L70 40 Z" fill="#000" />
        <path d="M50 25 L35 40 L35 60 L50 75 L65 60 L65 40 Z" fill="#d4af37" />
        <path d="M50 35 L45 50 L50 65 L55 50 Z" fill="#000" />
      </svg>
    ),
    tesla: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M50 20 L40 25 L40 30 L45 30 L45 75 L55 75 L55 30 L60 30 L60 25 Z"
          fill="#cc0000"
        />
        <path d="M30 20 L70 20" stroke="#cc0000" strokeWidth="3" />
      </svg>
    ),
    subaru: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="35" fill="#013b7a" />
        <g fill="white">
          <circle cx="30" cy="40" r="5" />
          <circle cx="45" cy="35" r="5" />
          <circle cx="60" cy="38" r="5" />
          <circle cx="70" cy="45" r="5" />
          <circle cx="55" cy="50" r="5" />
          <circle cx="40" cy="50" r="5" />
        </g>
      </svg>
    ),
    alfaromeo: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#9d1a2e" />
        <rect x="48" y="10" width="4" height="80" fill="white" />
        <path d="M20 50 L80 50" stroke="white" strokeWidth="4" />
        <path
          d="M35 35 C40 30, 60 30, 65 35"
          fill="none"
          stroke="white"
          strokeWidth="3"
        />
      </svg>
    ),
    infiniti: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M30 50 Q50 30, 70 50 Q50 70, 30 50"
          fill="none"
          stroke="#2d2d2d"
          strokeWidth="4"
        />
        <path d="M50 20 L50 80" stroke="#2d2d2d" strokeWidth="3" />
      </svg>
    ),
    mini: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#000" />
        <circle cx="50" cy="50" r="40" fill="#fff" />
        <circle cx="50" cy="50" r="35" fill="#000" />
        <text
          x="50"
          y="58"
          fontFamily="Arial Black"
          fontSize="20"
          fill="white"
          textAnchor="middle"
        >
          MINI
        </text>
      </svg>
    ),
    jaguar: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="30" fill="#005a2b" />
        <path
          d="M20 50 Q35 40, 50 50 Q65 40, 80 50"
          fill="none"
          stroke="white"
          strokeWidth="3"
        />
        <text
          x="50"
          y="65"
          fontFamily="Arial"
          fontSize="14"
          fill="white"
          textAnchor="middle"
        >
          JAGUAR
        </text>
      </svg>
    ),
    landrover: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="35" fill="#003a2d" />
        <text
          x="50"
          y="45"
          fontFamily="Arial"
          fontSize="12"
          fill="white"
          textAnchor="middle"
        >
          LAND
        </text>
        <text
          x="50"
          y="60"
          fontFamily="Arial"
          fontSize="12"
          fill="white"
          textAnchor="middle"
        >
          ROVER
        </text>
      </svg>
    ),
    rollsroyce: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="30" width="40" height="40" fill="#680021" />
        <rect x="50" y="30" width="40" height="40" fill="#000" />
        <text
          x="50"
          y="52"
          fontFamily="serif"
          fontSize="16"
          fill="white"
          textAnchor="middle"
        >
          RR
        </text>
      </svg>
    ),
    bentley: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#333" />
        <path
          d="M20 50 L30 30 L40 50 L50 30 L60 50 L70 30 L80 50"
          fill="none"
          stroke="white"
          strokeWidth="3"
        />
        <text
          x="50"
          y="70"
          fontFamily="serif"
          fontSize="16"
          fill="white"
          textAnchor="middle"
        >
          B
        </text>
      </svg>
    ),
    astonmartin: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M20 50 L50 30 L80 50 M30 50 L50 40 L70 50"
          fill="none"
          stroke="#003d2b"
          strokeWidth="4"
        />
        <text
          x="50"
          y="70"
          fontFamily="Arial"
          fontSize="10"
          fill="#003d2b"
          textAnchor="middle"
        >
          ASTON MARTIN
        </text>
      </svg>
    ),
    mclaren: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M20 50 Q35 30, 50 50 Q65 30, 80 50"
          fill="none"
          stroke="#ff8700"
          strokeWidth="5"
        />
        <text
          x="50"
          y="70"
          fontFamily="Arial"
          fontSize="14"
          fill="#ff8700"
          textAnchor="middle"
        >
          McLaren
        </text>
      </svg>
    ),
    maserati: (
      <svg viewBox="0 0 100 100" className={className}>
        <path d="M50 20 L45 80 L50 70 L55 80 L50 20" fill="#0c2340" />
        <path d="M40 40 L35 50 L40 45 M60 40 L65 50 L60 45" fill="#0c2340" />
      </svg>
    ),
    lotus: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#ffed00" />
        <circle cx="50" cy="50" r="40" fill="#046a38" />
        <path
          d="M50 25 Q40 35, 40 50 Q50 40, 50 25 Q50 40, 60 50 Q60 35, 50 25"
          fill="#ffed00"
        />
        <text
          x="50"
          y="70"
          fontFamily="Arial"
          fontSize="12"
          fill="white"
          textAnchor="middle"
        >
          LOTUS
        </text>
      </svg>
    ),
    bugatti: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="35" fill="#ce1126" />
        <rect x="20" y="45" width="60" height="10" fill="white" />
        <text
          x="50"
          y="40"
          fontFamily="serif"
          fontSize="14"
          fill="white"
          textAnchor="middle"
        >
          BUGATTI
        </text>
      </svg>
    ),
    genesis: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M20 50 L50 30 L80 50 M30 55 L50 40 L70 55"
          fill="none"
          stroke="#8b6f47"
          strokeWidth="4"
        />
        <text
          x="50"
          y="75"
          fontFamily="Arial"
          fontSize="12"
          fill="#8b6f47"
          textAnchor="middle"
        >
          GENESIS
        </text>
      </svg>
    ),
    acura: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M50 20 L30 70 L40 70 L50 35 L60 70 L70 70 Z"
          fill="none"
          stroke="#000"
          strokeWidth="4"
        />
      </svg>
    ),
    cadillac: (
      <svg viewBox="0 0 100 100" className={className}>
        <path d="M50 20 L35 35 L35 65 L50 80 L65 65 L65 35 Z" fill="#a4343a" />
        <path d="M50 30 L40 40 L40 60 L50 70 L60 60 L60 40 Z" fill="white" />
        <path d="M50 40 L45 45 L45 55 L50 60 L55 55 L55 45 Z" fill="#a4343a" />
      </svg>
    ),
    lincoln: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="35" width="80" height="30" fill="#000" />
        <path
          d="M30 50 L50 35 L70 50 L50 65 Z"
          fill="none"
          stroke="white"
          strokeWidth="3"
        />
      </svg>
    ),
    buick: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="30" cy="50" r="15" fill="#c1272d" />
        <circle cx="50" cy="50" r="15" fill="#fff" />
        <circle cx="70" cy="50" r="15" fill="#0039a6" />
      </svg>
    ),
    gmc: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="35" width="80" height="30" fill="#cc0000" />
        <text
          x="50"
          y="55"
          fontFamily="Arial Black"
          fontSize="24"
          fill="white"
          textAnchor="middle"
        >
          GMC
        </text>
      </svg>
    ),
    chrysler: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M50 20 L35 35 L20 50 L35 65 L50 80 L65 65 L80 50 L65 35 Z"
          fill="none"
          stroke="#263c80"
          strokeWidth="4"
        />
        <circle cx="50" cy="50" r="10" fill="#263c80" />
      </svg>
    ),
    dodge: (
      <svg viewBox="0 0 100 100" className={className}>
        <path d="M20 30 L20 70 L50 70 L80 50 L80 30 Z" fill="#cc0000" />
        <text
          x="50"
          y="55"
          fontFamily="Arial Black"
          fontSize="18"
          fill="white"
          textAnchor="middle"
        >
          DODGE
        </text>
      </svg>
    ),
    ram: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M30 40 Q30 20, 50 20 Q70 20, 70 40 L70 60 L60 60 L60 40 Q60 30, 50 30 Q40 30, 40 40 L40 60 L30 60 Z"
          fill="#000"
        />
        <text
          x="50"
          y="80"
          fontFamily="Arial Black"
          fontSize="20"
          fill="#000"
          textAnchor="middle"
        >
          RAM
        </text>
      </svg>
    ),
    saab: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="40" r="30" fill="#0066cc" />
        <path
          d="M50 20 Q60 25, 65 35 Q60 45, 50 50 Q40 45, 35 35 Q40 25, 50 20"
          fill="red"
        />
        <text
          x="50"
          y="80"
          fontFamily="Arial"
          fontSize="18"
          fill="#0066cc"
          textAnchor="middle"
        >
          SAAB
        </text>
      </svg>
    ),
    hummer: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="35" width="80" height="30" fill="#464646" />
        <text
          x="50"
          y="55"
          fontFamily="Arial Black"
          fontSize="18"
          fill="#f7c319"
          textAnchor="middle"
        >
          HUMMER
        </text>
      </svg>
    ),
    smart: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="35" cy="50" r="25" fill="#000" />
        <circle cx="65" cy="50" r="25" fill="#7fba00" />
        <text
          x="35"
          y="55"
          fontFamily="Arial"
          fontSize="14"
          fill="white"
          textAnchor="middle"
        >
          s
        </text>
        <text
          x="65"
          y="55"
          fontFamily="Arial"
          fontSize="14"
          fill="white"
          textAnchor="middle"
        >
          m
        </text>
      </svg>
    ),
    dacia: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="35" width="80" height="30" fill="#0055a4" />
        <text
          x="50"
          y="55"
          fontFamily="Arial"
          fontSize="18"
          fill="white"
          textAnchor="middle"
        >
          DACIA
        </text>
      </svg>
    ),
    lada: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="40" ry="30" fill="#cc0000" />
        <rect x="35" y="40" width="30" height="20" fill="white" />
        <path
          d="M40 45 L50 55 L60 45"
          fill="none"
          stroke="#cc0000"
          strokeWidth="3"
        />
      </svg>
    ),
    geely: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#003d7a" />
        <circle cx="50" cy="50" r="35" fill="#c8102e" />
        <circle cx="50" cy="50" r="25" fill="#003d7a" />
        <text
          x="50"
          y="55"
          fontFamily="Arial"
          fontSize="14"
          fill="white"
          textAnchor="middle"
        >
          GEELY
        </text>
      </svg>
    ),
    byd: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="35" fill="#ee1c25" />
        <text
          x="50"
          y="58"
          fontFamily="Arial Black"
          fontSize="24"
          fill="white"
          textAnchor="middle"
        >
          BYD
        </text>
      </svg>
    ),
    greatwall: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#00a651" />
        <rect x="30" y="45" width="40" height="10" fill="white" />
        <text
          x="50"
          y="40"
          fontFamily="Arial"
          fontSize="12"
          fill="white"
          textAnchor="middle"
        >
          GREAT WALL
        </text>
      </svg>
    ),
    haval: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="30" width="80" height="40" fill="#cc0000" />
        <text
          x="50"
          y="55"
          fontFamily="Arial Black"
          fontSize="20"
          fill="white"
          textAnchor="middle"
        >
          HAVAL
        </text>
      </svg>
    ),
    chery: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="35" fill="#000" />
        <path
          d="M30 50 Q50 30, 70 50"
          fill="none"
          stroke="white"
          strokeWidth="4"
        />
        <text
          x="50"
          y="70"
          fontFamily="Arial"
          fontSize="16"
          fill="white"
          textAnchor="middle"
        >
          CHERY
        </text>
      </svg>
    ),
    mg: (
      <svg viewBox="0 0 100 100" className={className}>
        <polygon points="50,20 30,70 70,70" fill="#c4161c" />
        <text
          x="50"
          y="55"
          fontFamily="Arial Black"
          fontSize="20"
          fill="white"
          textAnchor="middle"
        >
          MG
        </text>
      </svg>
    ),
    zastava: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#003893" />
        <path
          d="M50 20 L60 45 L85 45 L65 60 L75 85 L50 65 L25 85 L35 60 L15 45 L40 45 Z"
          fill="#ffcd00"
        />
      </svg>
    ),
    yugo: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="35" width="80" height="30" fill="#cc0000" />
        <rect x="10" y="35" width="27" height="30" fill="#003893" />
        <rect x="37" y="35" width="26" height="30" fill="white" />
        <text
          x="50"
          y="55"
          fontFamily="Arial Black"
          fontSize="18"
          fill="#003893"
          textAnchor="middle"
        >
          YUGO
        </text>
      </svg>
    ),
    fap: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#0055a4" />
        <text
          x="50"
          y="60"
          fontFamily="Arial Black"
          fontSize="24"
          fill="white"
          textAnchor="middle"
        >
          FAP
        </text>
      </svg>
    ),
    imt: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="35" width="80" height="30" fill="#ff6b35" />
        <text
          x="50"
          y="55"
          fontFamily="Arial Black"
          fontSize="20"
          fill="white"
          textAnchor="middle"
        >
          IMT
        </text>
      </svg>
    ),
    tata: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="35" fill="#003f87" />
        <text
          x="50"
          y="58"
          fontFamily="Arial"
          fontSize="20"
          fill="white"
          textAnchor="middle"
        >
          TATA
        </text>
      </svg>
    ),
    mahindra: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M20 50 L30 30 L40 50 L50 30 L60 50 L70 30 L80 50"
          fill="none"
          stroke="#ed1c24"
          strokeWidth="5"
        />
        <text
          x="50"
          y="75"
          fontFamily="Arial"
          fontSize="14"
          fill="#ed1c24"
          textAnchor="middle"
        >
          MAHINDRA
        </text>
      </svg>
    ),
    ds: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="35" width="35" height="30" fill="#000" />
        <rect x="45" y="35" width="45" height="30" fill="#8b7355" />
        <text
          x="27"
          y="55"
          fontFamily="Arial"
          fontSize="20"
          fill="white"
          textAnchor="middle"
        >
          D
        </text>
        <text
          x="67"
          y="55"
          fontFamily="Arial"
          fontSize="20"
          fill="white"
          textAnchor="middle"
        >
          S
        </text>
      </svg>
    ),
    cupra: (
      <svg viewBox="0 0 100 100" className={className}>
        <path d="M50 20 L30 80 L70 80 Z" fill="#8b7355" />
        <text
          x="50"
          y="65"
          fontFamily="Arial Black"
          fontSize="14"
          fill="white"
          textAnchor="middle"
        >
          CUPRA
        </text>
      </svg>
    ),
    polestar: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M50 20 L55 35 L70 35 L58 45 L63 60 L50 50 L37 60 L42 45 L30 35 L45 35 Z"
          fill="none"
          stroke="#4c9fd7"
          strokeWidth="3"
        />
        <path
          d="M50 20 L55 35 L70 35 L58 45 L63 60 L50 50 L37 60 L42 45 L30 35 L45 35 Z"
          fill="#4c9fd7"
          opacity="0.3"
        />
      </svg>
    ),
    rivian: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="35" cy="50" r="15" fill="#2d5016" />
        <circle cx="50" cy="50" r="15" fill="#2d5016" />
        <circle cx="65" cy="50" r="15" fill="#2d5016" />
        <circle cx="42.5" cy="50" r="15" fill="#5a9a30" />
        <circle cx="57.5" cy="50" r="15" fill="#5a9a30" />
      </svg>
    ),
    lucid: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M20 50 L40 30 L50 40 L60 30 L80 50"
          fill="none"
          stroke="#000"
          strokeWidth="4"
        />
        <text
          x="50"
          y="75"
          fontFamily="Arial"
          fontSize="16"
          fill="#000"
          textAnchor="middle"
        >
          LUCID
        </text>
      </svg>
    ),
    vinfast: (
      <svg viewBox="0 0 100 100" className={className}>
        <path d="M30 40 L50 20 L70 40 L50 60 Z" fill="#0066cc" />
        <path d="M50 30 L60 40 L50 50 L40 40 Z" fill="white" />
        <text
          x="50"
          y="80"
          fontFamily="Arial"
          fontSize="12"
          fill="#0066cc"
          textAnchor="middle"
        >
          VINFAST
        </text>
      </svg>
    ),
    nio: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#000" />
        <path
          d="M50 25 C60 25, 70 35, 70 45 C70 55, 60 65, 50 65 C40 65, 30 55, 30 45 C30 35, 40 25, 50 25"
          fill="none"
          stroke="white"
          strokeWidth="3"
        />
        <text
          x="50"
          y="55"
          fontFamily="Arial"
          fontSize="16"
          fill="white"
          textAnchor="middle"
        >
          NIO
        </text>
      </svg>
    ),
    xpeng: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M20 30 L40 70 M40 30 L20 70 M60 30 L60 70 M60 30 L80 30 M60 50 L75 50"
          fill="none"
          stroke="#00a0e9"
          strokeWidth="5"
        />
      </svg>
    ),
    lixiang: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="20" y="35" width="60" height="30" fill="#000" />
        <text
          x="50"
          y="55"
          fontFamily="Arial"
          fontSize="16"
          fill="white"
          textAnchor="middle"
        >
          Li Auto
        </text>
      </svg>
    ),
    hongqi: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="35" width="80" height="30" fill="#cc0000" />
        <path
          d="M50 20 L55 35 L70 35 L58 45 L63 60 L50 50 L37 60 L42 45 L30 35 L45 35 Z"
          fill="#ffcd00"
        />
      </svg>
    ),
    daewoo: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="35" fill="#003f7c" />
        <text
          x="50"
          y="58"
          fontFamily="Arial"
          fontSize="16"
          fill="white"
          textAnchor="middle"
        >
          DAEWOO
        </text>
      </svg>
    ),
    scion: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="35" width="80" height="30" fill="#000" />
        <text
          x="50"
          y="55"
          fontFamily="Arial Black"
          fontSize="20"
          fill="#7c878e"
          textAnchor="middle"
        >
          SCION
        </text>
      </svg>
    ),
    pontiac: (
      <svg viewBox="0 0 100 100" className={className}>
        <path d="M50 20 L35 50 L50 80 L65 50 Z" fill="#cc0000" />
        <path d="M50 30 L40 50 L50 70 L60 50 Z" fill="white" />
        <path d="M50 40 L45 50 L50 60 L55 50 Z" fill="#cc0000" />
      </svg>
    ),
    saturn: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="30" fill="#cc0000" />
        <ellipse
          cx="50"
          cy="50"
          rx="45"
          ry="10"
          fill="none"
          stroke="#cc0000"
          strokeWidth="5"
        />
        <text
          x="50"
          y="55"
          fontFamily="Arial"
          fontSize="12"
          fill="white"
          textAnchor="middle"
        >
          SATURN
        </text>
      </svg>
    ),
    mercury: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="35" r="20" fill="#000" />
        <path d="M50 55 L50 75 M40 65 L60 65" stroke="#000" strokeWidth="4" />
        <text
          x="50"
          y="90"
          fontFamily="Arial"
          fontSize="12"
          fill="#000"
          textAnchor="middle"
        >
          MERCURY
        </text>
      </svg>
    ),
    oldsmobile: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="25" fill="#cc0000" />
        <text
          x="50"
          y="55"
          fontFamily="serif"
          fontSize="14"
          fill="white"
          textAnchor="middle"
        >
          Oldsmobile
        </text>
      </svg>
    ),
    plymouth: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#003f7c" />
        <path
          d="M50 20 L40 40 L30 50 L40 60 L50 80 L60 60 L70 50 L60 40 Z"
          fill="white"
        />
        <text
          x="50"
          y="55"
          fontFamily="Arial"
          fontSize="10"
          fill="#003f7c"
          textAnchor="middle"
        >
          PLYMOUTH
        </text>
      </svg>
    ),
    geo: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#0066cc" />
        <text
          x="50"
          y="60"
          fontFamily="Arial Black"
          fontSize="28"
          fill="white"
          textAnchor="middle"
        >
          GEO
        </text>
      </svg>
    ),
    eagle: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M50 20 L30 40 L20 50 L30 45 L40 50 L50 40 L60 50 L70 45 L80 50 L70 40 Z"
          fill="#8b4513"
        />
        <text
          x="50"
          y="75"
          fontFamily="Arial"
          fontSize="16"
          fill="#8b4513"
          textAnchor="middle"
        >
          EAGLE
        </text>
      </svg>
    ),
    fisker: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="45" r="30" fill="#ff8c00" />
        <circle cx="50" cy="45" r="25" fill="white" />
        <circle cx="50" cy="45" r="20" fill="#ff8c00" />
        <text
          x="50"
          y="80"
          fontFamily="Arial"
          fontSize="14"
          fill="#ff8c00"
          textAnchor="middle"
        >
          FISKER
        </text>
      </svg>
    ),
    karma: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M50 20 Q30 35, 30 50 Q30 65, 50 80 Q70 65, 70 50 Q70 35, 50 20"
          fill="none"
          stroke="#b8860b"
          strokeWidth="4"
        />
        <text
          x="50"
          y="55"
          fontFamily="Arial"
          fontSize="14"
          fill="#b8860b"
          textAnchor="middle"
        >
          KARMA
        </text>
      </svg>
    ),
    ineos: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="35" width="80" height="30" fill="#000" />
        <text
          x="50"
          y="55"
          fontFamily="Arial Black"
          fontSize="18"
          fill="white"
          textAnchor="middle"
        >
          INEOS
        </text>
      </svg>
    ),
    isuzu: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="35" width="80" height="30" fill="#cc0000" />
        <text
          x="50"
          y="55"
          fontFamily="Arial Black"
          fontSize="20"
          fill="white"
          textAnchor="middle"
        >
          ISUZU
        </text>
      </svg>
    ),
    daihatsu: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="30" width="80" height="40" fill="#cc0000" />
        <path d="M50 35 L40 55 L50 65 L60 55 Z" fill="white" />
        <text
          x="50"
          y="50"
          fontFamily="Arial"
          fontSize="8"
          fill="#cc0000"
          textAnchor="middle"
        >
          D
        </text>
      </svg>
    ),
    lancia: (
      <svg viewBox="0 0 100 100" className={className}>
        <path d="M50 20 L30 50 L50 80 L70 50 Z" fill="#001489" />
        <text
          x="50"
          y="55"
          fontFamily="serif"
          fontSize="14"
          fill="white"
          textAnchor="middle"
        >
          LANCIA
        </text>
      </svg>
    ),
    rover: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M50 20 L40 30 L30 50 L40 70 L50 80 L60 70 L70 50 L60 30 Z"
          fill="none"
          stroke="#003d2b"
          strokeWidth="4"
        />
        <text
          x="50"
          y="55"
          fontFamily="Arial"
          fontSize="12"
          fill="#003d2b"
          textAnchor="middle"
        >
          ROVER
        </text>
      </svg>
    ),
    proton: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="40" r="25" fill="#ffd700" />
        <path
          d="M50 20 L45 40 L30 35 L40 45 L35 60 L50 50 L65 60 L60 45 L70 35 L55 40 Z"
          fill="#003f7c"
        />
        <text
          x="50"
          y="80"
          fontFamily="Arial"
          fontSize="14"
          fill="#003f7c"
          textAnchor="middle"
        >
          PROTON
        </text>
      </svg>
    ),
    perodua: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="35" fill="#cc0000" />
        <text
          x="50"
          y="58"
          fontFamily="Arial"
          fontSize="14"
          fill="white"
          textAnchor="middle"
        >
          PERODUA
        </text>
      </svg>
    ),
    alpine: (
      <svg viewBox="0 0 100 100" className={className}>
        <path d="M50 20 L30 60 L45 60 L50 45 L55 60 L70 60 Z" fill="#0055a4" />
        <text
          x="50"
          y="80"
          fontFamily="Arial"
          fontSize="14"
          fill="#0055a4"
          textAnchor="middle"
        >
          ALPINE
        </text>
      </svg>
    ),
    abarth: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M50 25 C45 20, 35 20, 30 30 C25 40, 25 60, 35 70 L50 75 L65 70 C75 60, 75 40, 70 30 C65 20, 55 20, 50 25"
          fill="#cc0000"
        />
        <path
          d="M40 35 C40 35, 45 45, 50 40 C55 45, 60 35, 60 35"
          fill="#ffd700"
        />
        <path d="M45 50 L50 60 L55 50" fill="#000" />
      </svg>
    ),
    datsun: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#cc0000" />
        <rect x="20" y="45" width="60" height="10" fill="white" />
        <text
          x="50"
          y="40"
          fontFamily="Arial"
          fontSize="14"
          fill="white"
          textAnchor="middle"
        >
          DATSUN
        </text>
      </svg>
    ),
    maruti: (
      <svg viewBox="0 0 100 100" className={className}>
        <ellipse cx="50" cy="50" rx="45" ry="35" fill="#003f7c" />
        <text
          x="50"
          y="58"
          fontFamily="Arial"
          fontSize="16"
          fill="white"
          textAnchor="middle"
        >
          MARUTI
        </text>
      </svg>
    ),
    amgeneral: (
      <svg viewBox="0 0 100 100" className={className}>
        <rect x="10" y="35" width="80" height="30" fill="#4b5320" />
        <text
          x="50"
          y="55"
          fontFamily="Arial Black"
          fontSize="14"
          fill="white"
          textAnchor="middle"
        >
          AM GENERAL
        </text>
      </svg>
    ),
    maybach: (
      <svg viewBox="0 0 100 100" className={className}>
        <path
          d="M50 20 L40 50 L45 50 L50 30 L55 50 L60 50 Z M40 50 L35 70 L40 70 L45 55 L50 70 L55 70 L60 55 L65 70 L70 70 L65 50 L60 50 L55 65 L50 50 L45 65 L40 50 Z"
          fill="#8b7355"
        />
      </svg>
    ),
    panoz: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#000" />
        <path
          d="M50 20 L35 50 L50 45 L65 50 Z M35 50 L50 80 L50 55 M65 50 L50 80 L50 55"
          fill="none"
          stroke="#cc0000"
          strokeWidth="3"
        />
      </svg>
    ),
    spyker: (
      <svg viewBox="0 0 100 100" className={className}>
        <circle cx="50" cy="50" r="45" fill="#ff8c00" />
        <path
          d="M50 20 L40 40 L30 50 L40 60 L50 80 L60 60 L70 50 L60 40 Z"
          fill="none"
          stroke="white"
          strokeWidth="3"
        />
        <circle cx="50" cy="50" r="10" fill="white" />
      </svg>
    ),
  };

  const icon = brandIcons[normalizedBrand];

  if (icon) {
    return icon;
  }

  // Fallback с первой буквой для неизвестных марок
  return (
    <div
      className={`${className} bg-neutral text-neutral-content rounded-full flex items-center justify-center font-bold text-2xl`}
    >
      {brand.charAt(0).toUpperCase()}
    </div>
  );
};

export default CarBrandIcon;
