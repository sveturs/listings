'use client';

import React, { useEffect, useRef, useState, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { toast } from 'react-hot-toast';
import configManager from '@/config';

interface QRBarcodeScannerProps {
  onScan?: (data: ScanResult) => void;
  onClose?: () => void;
  mode?: 'qr' | 'barcode' | 'both';
  autoSearch?: boolean;
  autoFill?: boolean;
}

interface ScanResult {
  type: 'qr' | 'barcode' | 'ean' | 'upc' | 'isbn' | 'unknown';
  format: string;
  rawValue: string;
  parsedData?: {
    productId?: string;
    url?: string;
    ean?: string;
    upc?: string;
    isbn?: string;
    attributes?: Record<string, any>;
  };
  confidence: number;
  timestamp: number;
}

interface ProductInfo {
  name: string;
  brand?: string;
  category?: string;
  attributes?: Record<string, any>;
  imageUrl?: string;
  price?: number;
}

export default function QRBarcodeScanner({
  onScan,
  onClose,
  mode = 'both',
  autoSearch = true,
  autoFill = true,
}: QRBarcodeScannerProps) {
  const t = useTranslations('scanner');
  const router = useRouter();
  const videoRef = useRef<HTMLVideoElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const streamRef = useRef<MediaStream | null>(null);

  const [isScanning, setIsScanning] = useState(false);
  const [lastScan, setLastScan] = useState<ScanResult | null>(null);
  const [productInfo, setProductInfo] = useState<ProductInfo | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [cameraError, setCameraError] = useState<string | null>(null);
  const [flashEnabled, setFlashEnabled] = useState(false);
  const [currentCamera, setCurrentCamera] = useState<'user' | 'environment'>(
    'environment'
  );
  const [scanHistory, setScanHistory] = useState<ScanResult[]>([]);

  const scanIntervalRef = useRef<NodeJS.Timeout | null>(null);
  const barcodeDetectorRef = useRef<any>(null);

  // Initialize camera and scanner
  useEffect(() => {
    initializeScanner();
    return () => {
      stopScanning();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentCamera]);

  const initializeScanner = async () => {
    try {
      // Check if BarcodeDetector is available
      if ('BarcodeDetector' in window) {
        const BarcodeDetector = (window as any).BarcodeDetector;
        const supportedFormats = await BarcodeDetector.getSupportedFormats();

        // Configure detector based on mode
        let formats = supportedFormats;
        if (mode === 'qr') {
          formats = formats.filter((f: string) => f.includes('qr'));
        } else if (mode === 'barcode') {
          formats = formats.filter((f: string) => !f.includes('qr'));
        }

        barcodeDetectorRef.current = new BarcodeDetector({ formats });
      } else {
        // Fallback to library-based scanning (e.g., ZXing)
        await loadFallbackScanner();
      }

      // Start camera
      await startCamera();
    } catch {
      console.error('Scanner initialization error');
      setCameraError(t('scannerInitError'));
    }
  };

  const loadFallbackScanner = async () => {
    // Dynamic import of fallback scanner library
    // TODO: Install @zxing/browser package for fallback scanner
    // const { BrowserMultiFormatReader } = await import('@zxing/browser');
    // barcodeDetectorRef.current = new BrowserMultiFormatReader();
    console.warn(
      'Fallback scanner not available - @zxing/browser not installed'
    );
  };

  const startCamera = async () => {
    try {
      const constraints: MediaStreamConstraints = {
        video: {
          facingMode: currentCamera,
          width: { ideal: 1920 },
          height: { ideal: 1080 },
        },
      };

      const stream = await navigator.mediaDevices.getUserMedia(constraints);
      streamRef.current = stream;

      if (videoRef.current) {
        videoRef.current.srcObject = stream;
        await videoRef.current.play();
        setIsScanning(true);
        startScanLoop();
      }

      // Check for torch support (flashlight)
      const track = stream.getVideoTracks()[0];
      const capabilities = track.getCapabilities ? track.getCapabilities() : {};
      if ('torch' in capabilities) {
        // Flash is available
      }
    } catch {
      console.error('Camera error');
      setCameraError(t('cameraAccessError'));
    }
  };

  const stopScanning = () => {
    if (scanIntervalRef.current) {
      clearInterval(scanIntervalRef.current);
      scanIntervalRef.current = null;
    }

    if (streamRef.current) {
      streamRef.current.getTracks().forEach((track) => track.stop());
      streamRef.current = null;
    }

    if (videoRef.current) {
      videoRef.current.srcObject = null;
    }

    setIsScanning(false);
  };

  const startScanLoop = () => {
    scanIntervalRef.current = setInterval(async () => {
      if (videoRef.current && videoRef.current.readyState === 4) {
        await scanFrame();
      }
    }, 100); // Scan 10 times per second
  };

  const scanFrame = async () => {
    if (!videoRef.current || !canvasRef.current) return;

    const canvas = canvasRef.current;
    const video = videoRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Set canvas size to match video
    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;

    // Draw current frame
    ctx.drawImage(video, 0, 0, canvas.width, canvas.height);

    // Apply image enhancements for better scanning
    enhanceImage(ctx, canvas.width, canvas.height);

    try {
      let detectedCodes: any[] = [];

      if (barcodeDetectorRef.current) {
        if ('detect' in barcodeDetectorRef.current) {
          // Native BarcodeDetector
          detectedCodes = await barcodeDetectorRef.current.detect(canvas);
        } else {
          // Fallback library - TODO: Implement when @zxing/browser is installed
          // const result =
          //   await barcodeDetectorRef.current.decodeFromCanvas(canvas);
          //     detectedCodes = [
          //       {
          //         rawValue: result.getText(),
          //         format: result.getBarcodeFormat(),
          //         boundingBox: result.getResultPoints(),
          //       },
          //     ];
          // }
        }
      }

      if (detectedCodes.length > 0) {
        const code = detectedCodes[0];
        await handleDetectedCode(code);
      }
    } catch {
      // Scanning failed for this frame, continue
    }
  };

  const enhanceImage = (
    ctx: CanvasRenderingContext2D,
    width: number,
    height: number
  ) => {
    // Enhance contrast and brightness for better scanning
    const imageData = ctx.getImageData(0, 0, width, height);
    const data = imageData.data;

    for (let i = 0; i < data.length; i += 4) {
      // Increase contrast
      data[i] = Math.min(255, data[i] * 1.2); // Red
      data[i + 1] = Math.min(255, data[i + 1] * 1.2); // Green
      data[i + 2] = Math.min(255, data[i + 2] * 1.2); // Blue
    }

    ctx.putImageData(imageData, 0, 0);
  };

  const handleDetectedCode = async (code: any) => {
    const scanResult: ScanResult = {
      type: detectCodeType(code.format),
      format: code.format,
      rawValue: code.rawValue,
      confidence: code.confidence || 1,
      timestamp: Date.now(),
      parsedData: parseCodeData(code.rawValue, code.format),
    };

    // Avoid duplicate scans
    if (
      lastScan &&
      lastScan.rawValue === scanResult.rawValue &&
      Date.now() - lastScan.timestamp < 2000
    ) {
      return;
    }

    setLastScan(scanResult);
    setScanHistory((prev) => [scanResult, ...prev.slice(0, 9)]);

    // Haptic feedback
    if ('vibrate' in navigator) {
      navigator.vibrate(200);
    }

    // Sound feedback
    playBeep();

    // Visual feedback
    showScanAnimation();

    if (onScan) {
      onScan(scanResult);
    }

    // Auto-process the scan
    if (autoSearch || autoFill) {
      await processScanResult(scanResult);
    }
  };

  const detectCodeType = (format: string): ScanResult['type'] => {
    const formatLower = format.toLowerCase();
    if (formatLower.includes('qr')) return 'qr';
    if (formatLower.includes('ean')) return 'ean';
    if (formatLower.includes('upc')) return 'upc';
    if (formatLower.includes('isbn')) return 'isbn';
    if (formatLower.includes('code')) return 'barcode';
    return 'unknown';
  };

  const parseCodeData = (
    rawValue: string,
    format: string
  ): ScanResult['parsedData'] => {
    const parsed: ScanResult['parsedData'] = {};

    // Check if it's a URL
    try {
      const url = new URL(rawValue);
      parsed.url = url.href;

      // Check if it's our marketplace URL
      if (url.hostname === 'svetu.rs' || url.hostname === 'localhost') {
        const pathParts = url.pathname.split('/');
        if (pathParts.includes('product')) {
          const productIndex = pathParts.indexOf('product');
          if (productIndex + 1 < pathParts.length) {
            parsed.productId = pathParts[productIndex + 1];
          }
        }
      }
    } catch {
      // Not a URL
    }

    // Check barcode format
    if (format.toLowerCase().includes('ean')) {
      parsed.ean = rawValue;
    } else if (format.toLowerCase().includes('upc')) {
      parsed.upc = rawValue;
    } else if (format.toLowerCase().includes('isbn')) {
      parsed.isbn = rawValue;
    }

    // Try to parse as JSON (for QR codes with data)
    try {
      const jsonData = JSON.parse(rawValue);
      if (jsonData.attributes) {
        parsed.attributes = jsonData.attributes;
      }
      if (jsonData.productId) {
        parsed.productId = jsonData.productId;
      }
    } catch {
      // Not JSON
    }

    return parsed;
  };

  const processScanResult = async (scanResult: ScanResult) => {
    setIsLoading(true);

    try {
      if (scanResult.parsedData?.productId) {
        // Navigate to product page
        if (autoSearch) {
          router.push(`/product/${scanResult.parsedData.productId}`);
        }
      } else if (scanResult.parsedData?.url) {
        // Handle URL
        if (autoSearch && scanResult.parsedData.url.includes('svetu.rs')) {
          window.location.href = scanResult.parsedData.url;
        }
      } else if (
        scanResult.parsedData?.ean ||
        scanResult.parsedData?.upc ||
        scanResult.parsedData?.isbn
      ) {
        // Look up product by barcode
        const productInfo = await lookupProductByBarcode(scanResult);
        setProductInfo(productInfo);

        if (autoFill && productInfo) {
          fillProductForm(productInfo);
        }
      }
    } catch {
      console.error('Error processing scan');
      toast.error(t('processingError'));
    } finally {
      setIsLoading(false);
    }
  };

  const lookupProductByBarcode = async (
    scanResult: ScanResult
  ): Promise<ProductInfo | null> => {
    try {
      const apiUrl = configManager.get('api.url');
      const response = await fetch(`${apiUrl}/api/v1/products/barcode-lookup`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          barcode: scanResult.rawValue,
          type: scanResult.type,
          format: scanResult.format,
        }),
      });

      if (!response.ok) {
        // Try external API
        return await lookupExternalBarcode(scanResult);
      }

      return await response.json();
    } catch {
      return await lookupExternalBarcode(scanResult);
    }
  };

  const lookupExternalBarcode = async (
    scanResult: ScanResult
  ): Promise<ProductInfo | null> => {
    // Mock implementation - would call external barcode API
    try {
      // Simulate API call to Open Food Facts, UPC Database, etc.
      const mockProducts: Record<string, ProductInfo> = {
        '8690504180098': {
          name: 'iPhone 13 Pro',
          brand: 'Apple',
          category: 'electronics',
          attributes: { color: 'Space Gray', storage: '256GB' },
        },
        '9780545010221': {
          name: 'Harry Potter and the Deathly Hallows',
          brand: 'Scholastic',
          category: 'books',
          attributes: { author: 'J.K. Rowling', pages: '784' },
        },
      };

      return mockProducts[scanResult.rawValue] || null;
    } catch {
      return null;
    }
  };

  const fillProductForm = (productInfo: ProductInfo) => {
    // Auto-fill create listing form
    if (window.location.pathname.includes('create')) {
      // Dispatch event to fill form
      window.dispatchEvent(
        new CustomEvent('fillProductData', {
          detail: productInfo,
        })
      );
      toast.success(t('autoFilled'));
    }
  };

  const playBeep = () => {
    const audioContext = new (window.AudioContext ||
      (window as any).webkitAudioContext)();
    const oscillator = audioContext.createOscillator();
    const gainNode = audioContext.createGain();

    oscillator.connect(gainNode);
    gainNode.connect(audioContext.destination);

    oscillator.frequency.value = 1000;
    oscillator.type = 'sine';

    gainNode.gain.setValueAtTime(0.3, audioContext.currentTime);
    gainNode.gain.exponentialRampToValueAtTime(
      0.01,
      audioContext.currentTime + 0.1
    );

    oscillator.start(audioContext.currentTime);
    oscillator.stop(audioContext.currentTime + 0.1);
  };

  const showScanAnimation = () => {
    // Visual feedback animation
    if (canvasRef.current) {
      const ctx = canvasRef.current.getContext('2d');
      if (!ctx) return;

      ctx.strokeStyle = '#00ff00';
      ctx.lineWidth = 5;
      ctx.strokeRect(0, 0, canvasRef.current.width, canvasRef.current.height);

      setTimeout(() => {
        if (canvasRef.current && videoRef.current) {
          ctx.drawImage(videoRef.current, 0, 0);
        }
      }, 200);
    }
  };

  const toggleFlash = useCallback(async () => {
    if (!streamRef.current) return;

    try {
      const track = streamRef.current.getVideoTracks()[0];
      const capabilities = track.getCapabilities ? track.getCapabilities() : {};

      if ('torch' in capabilities) {
        await (track as any).applyConstraints({
          advanced: [{ torch: !flashEnabled }],
        });
        setFlashEnabled(!flashEnabled);
      } else {
        toast.error(t('flashNotSupported'));
      }
    } catch {
      console.error('Flash toggle error:');
    }
  }, [flashEnabled, t]);

  const switchCamera = useCallback(() => {
    setCurrentCamera((prev) => (prev === 'user' ? 'environment' : 'user'));
  }, []);

  const manualInput = useCallback(() => {
    const input = prompt(t('enterBarcodeManually'));
    if (input) {
      handleDetectedCode({
        rawValue: input,
        format: 'manual',
        confidence: 1,
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [t]);

  return (
    <div className="fixed inset-0 z-50 bg-black">
      {/* Camera View */}
      <div className="relative h-full w-full">
        <video
          ref={videoRef}
          className="h-full w-full object-cover"
          playsInline
          muted
        />

        <canvas
          ref={canvasRef}
          className="absolute left-0 top-0 h-full w-full"
          style={{ display: 'none' }}
        />

        {/* Scanning Frame */}
        <div className="absolute inset-0 flex items-center justify-center">
          <div className="relative h-64 w-64">
            {/* Animated scanning line */}
            {isScanning && (
              <div className="absolute inset-x-0 top-0 h-0.5 animate-scan bg-gradient-to-r from-transparent via-green-400 to-transparent" />
            )}

            {/* Corner markers */}
            <div className="absolute left-0 top-0 h-8 w-8 border-l-4 border-t-4 border-white" />
            <div className="absolute right-0 top-0 h-8 w-8 border-r-4 border-t-4 border-white" />
            <div className="absolute bottom-0 left-0 h-8 w-8 border-b-4 border-l-4 border-white" />
            <div className="absolute bottom-0 right-0 h-8 w-8 border-b-4 border-r-4 border-white" />
          </div>
        </div>

        {/* Controls */}
        <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent p-4">
          <div className="flex items-center justify-between">
            <button
              className="btn btn-circle btn-ghost text-white"
              onClick={onClose}
            >
              ✕
            </button>

            <div className="flex gap-2">
              <button
                className="btn btn-circle btn-ghost text-white"
                onClick={toggleFlash}
              >
                {flashEnabled ? '🔦' : '💡'}
              </button>

              <button
                className="btn btn-circle btn-ghost text-white"
                onClick={switchCamera}
              >
                🔄
              </button>

              <button
                className="btn btn-circle btn-ghost text-white"
                onClick={manualInput}
              >
                ⌨️
              </button>
            </div>
          </div>

          {/* Scan instruction */}
          {!lastScan && (
            <p className="mt-2 text-center text-sm text-white">
              {t('pointCameraAtCode')}
            </p>
          )}
        </div>

        {/* Last Scan Result */}
        {lastScan && (
          <div className="absolute left-4 right-4 top-4">
            <div className="rounded-lg bg-white/90 p-4 backdrop-blur">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <p className="text-sm font-medium">
                    {t('scanned')} {lastScan.type}
                  </p>
                  <p className="mt-1 font-mono text-xs">{lastScan.rawValue}</p>

                  {productInfo && (
                    <div className="mt-2 border-t pt-2">
                      <p className="font-semibold">{productInfo.name}</p>
                      {productInfo.brand && (
                        <p className="text-sm text-gray-600">
                          {productInfo.brand}
                        </p>
                      )}
                    </div>
                  )}
                </div>

                {isLoading && (
                  <div className="loading loading-spinner loading-sm" />
                )}
              </div>

              {lastScan.parsedData?.productId && (
                <button
                  className="btn btn-primary btn-sm mt-2 w-full"
                  onClick={() =>
                    router.push(`/product/${lastScan.parsedData?.productId}`)
                  }
                >
                  {t('viewProduct')}
                </button>
              )}
            </div>
          </div>
        )}

        {/* Scan History */}
        {scanHistory.length > 1 && (
          <div className="absolute right-4 top-20">
            <button
              className="btn btn-circle btn-ghost btn-sm text-white"
              onClick={() => {
                // Show history modal
              }}
            >
              📜
            </button>
          </div>
        )}

        {/* Error Message */}
        {cameraError && (
          <div className="absolute inset-0 flex items-center justify-center bg-black/80">
            <div className="rounded-lg bg-white p-6 text-center">
              <p className="text-error">{cameraError}</p>
              <button
                className="btn btn-primary mt-4"
                onClick={initializeScanner}
              >
                {t('retry')}
              </button>
            </div>
          </div>
        )}
      </div>

      <style jsx>{`
        @keyframes scan {
          0% {
            transform: translateY(0);
          }
          100% {
            transform: translateY(256px);
          }
        }

        .animate-scan {
          animation: scan 2s linear infinite;
        }
      `}</style>
    </div>
  );
}
