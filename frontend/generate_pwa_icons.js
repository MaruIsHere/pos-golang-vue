import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const publicDir = path.join(__dirname, 'public');
if (!fs.existsSync(publicDir)) {
  fs.mkdirSync(publicDir, { recursive: true });
}

// Generate high quality SVG App Icon
const svgContent = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="512" height="512">
  <defs>
    <linearGradient id="bgGrad" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#4f46e5" />
      <stop offset="50%" stop-color="#7c3aed" />
      <stop offset="100%" stop-color="#0f172a" />
    </linearGradient>
    <filter id="shadow" x="-20%" y="-20%" width="140%" height="140%">
      <feDropShadow dx="0" dy="8" stdDeviation="12" flood-color="#000000" flood-opacity="0.4" />
    </filter>
  </defs>

  <!-- Background Card -->
  <rect width="512" height="512" rx="128" fill="url(#bgGrad)" />

  <!-- Shopping Cart & POS Icon -->
  <g filter="url(#shadow)" transform="translate(40, 30)">
    <!-- POS Terminal Body -->
    <rect x="120" y="80" width="192" height="260" rx="24" fill="#1e293b" stroke="#6366f1" stroke-width="8" />
    <!-- Screen -->
    <rect x="144" y="104" width="144" height="120" rx="12" fill="#020617" stroke="#38bdf8" stroke-width="4" />
    <text x="216" y="160" font-family="Arial, sans-serif" font-size="28" font-weight="bold" fill="#38bdf8" text-anchor="middle">POS PRO</text>
    <text x="216" y="195" font-family="Arial, sans-serif" font-size="18" font-weight="bold" fill="#34d399" text-anchor="middle">Rp 150.000</text>
    
    <!-- Keypad Buttons -->
    <circle cx="168" cy="256" r="10" fill="#6366f1" />
    <circle cx="216" cy="256" r="10" fill="#6366f1" />
    <circle cx="264" cy="256" r="10" fill="#6366f1" />

    <circle cx="168" cy="288" r="10" fill="#6366f1" />
    <circle cx="216" cy="288" r="10" fill="#6366f1" />
    <circle cx="264" cy="288" r="10" fill="#6366f1" />

    <!-- Receipt paper coming out -->
    <path d="M 160 80 L 160 30 L 272 30 L 272 80 Z" fill="#ffffff" />
    <line x1="176" y1="45" x2="256" y2="45" stroke="#94a3b8" stroke-width="4" stroke-linecap="round" />
    <line x1="176" y1="60" x2="236" y2="60" stroke="#94a3b8" stroke-width="4" stroke-linecap="round" />
  </g>
</svg>`;

fs.writeFileSync(path.join(publicDir, 'icon.svg'), svgContent);
fs.writeFileSync(path.join(publicDir, 'pwa-192x192.png'), svgContent);
fs.writeFileSync(path.join(publicDir, 'pwa-512x512.png'), svgContent);
fs.writeFileSync(path.join(publicDir, 'apple-touch-icon.png'), svgContent);

console.log('PWA Icons generated successfully in public/');
