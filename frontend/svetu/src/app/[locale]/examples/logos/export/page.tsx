'use client';

import React, { useRef, useEffect } from 'react';
import { SveTuLogoStatic } from '@/components/logos/SveTuLogoStatic';

export default function ExportLogoPage() {
  const svgRef = useRef<SVGSVGElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    // Convert SVG to PNG when component mounts
    if (svgRef.current && canvasRef.current) {
      convertSvgToPng();
    }
  }, []);

  const convertSvgToPng = () => {
    if (!svgRef.current || !canvasRef.current) return;

    const svgElement = svgRef.current;
    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Get SVG string
    const svgString = new XMLSerializer().serializeToString(svgElement);
    const svgBlob = new Blob([svgString], { type: 'image/svg+xml;charset=utf-8' });
    const url = URL.createObjectURL(svgBlob);

    // Create image from SVG
    const img = new Image();
    img.onload = () => {
      ctx.clearRect(0, 0, 48, 48);
      ctx.drawImage(img, 0, 0, 48, 48);
      URL.revokeObjectURL(url);
    };
    img.src = url;
  };

  const downloadPng = () => {
    if (!canvasRef.current) return;

    canvasRef.current.toBlob((blob) => {
      if (!blob) return;
      
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'svetu-gradient-48x48.png';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    });
  };

  const downloadSvg = () => {
    if (!svgRef.current) return;

    const svgString = new XMLSerializer().serializeToString(svgRef.current);
    const blob = new Blob([svgString], { type: 'image/svg+xml;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'svetu-gradient-48x48.svg';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  return (
    <div className="min-h-screen bg-gray-100 py-8">
      <div className="container mx-auto px-4">
        <h1 className="text-3xl font-bold text-center mb-8">Экспорт логотипа</h1>
        
        <div className="max-w-2xl mx-auto bg-white rounded-lg shadow-lg p-8">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            {/* SVG Preview */}
            <div>
              <h2 className="text-xl font-semibold mb-4">SVG (48×48)</h2>
              <div className="border-2 border-gray-300 rounded p-4 bg-gray-50 flex justify-center items-center" style={{ width: '200px', height: '200px' }}>
                <div ref={(el) => el && (svgRef.current = el.querySelector('svg'))}>
                  <SveTuLogoStatic variant="gradient" width={48} height={48} />
                </div>
              </div>
              <button
                onClick={downloadSvg}
                className="mt-4 btn btn-primary btn-sm"
              >
                Скачать SVG
              </button>
            </div>

            {/* Canvas Preview */}
            <div>
              <h2 className="text-xl font-semibold mb-4">PNG (48×48)</h2>
              <div className="border-2 border-gray-300 rounded p-4 bg-gray-50 flex justify-center items-center" style={{ width: '200px', height: '200px' }}>
                <canvas
                  ref={canvasRef}
                  width={48}
                  height={48}
                  style={{ imageRendering: 'pixelated', width: '48px', height: '48px' }}
                />
              </div>
              <button
                onClick={downloadPng}
                className="mt-4 btn btn-primary btn-sm"
              >
                Скачать PNG
              </button>
            </div>
          </div>

          {/* All variants */}
          <div className="mt-8 pt-8 border-t">
            <h2 className="text-xl font-semibold mb-4">Все варианты для экспорта</h2>
            <div className="grid grid-cols-5 gap-4">
              {['gradient', 'minimal', 'retro', 'neon', 'glassmorphic'].map((variant) => (
                <div key={variant} className="text-center">
                  <div className="border rounded p-2 mb-2">
                    <SveTuLogoStatic variant={variant as any} width={48} height={48} />
                  </div>
                  <p className="text-sm text-gray-600 capitalize">{variant}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}