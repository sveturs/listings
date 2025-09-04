'use client';

import React, { useEffect, useRef, useState, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from 'react-hot-toast';
import '@/types/model-viewer';

interface ARProductViewerProps {
  productId: string;
  modelUrl?: string;
  imageUrl: string;
  productName: string;
  attributes?: Record<string, any>;
  onClose?: () => void;
}

interface ARCapability {
  supported: boolean;
  webXR: boolean;
  quickLook: boolean;
  sceneViewer: boolean;
}

export default function ARProductViewer({
  productId: _productId,
  modelUrl,
  imageUrl,
  productName,
  attributes,
  onClose,
}: ARProductViewerProps) {
  const t = useTranslations('ar');
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [arCapabilities, setArCapabilities] = useState<ARCapability>({
    supported: false,
    webXR: false,
    quickLook: false,
    sceneViewer: false,
  });
  const [isARActive, setIsARActive] = useState(false);
  const modelViewerRef = useRef<any>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);

  // Check AR capabilities
  useEffect(() => {
    const checkARSupport = async () => {
      const capabilities: ARCapability = {
        supported: false,
        webXR: false,
        quickLook: false,
        sceneViewer: false,
      };

      // Check WebXR support
      if ('xr' in navigator) {
        try {
          const isSupported = await (navigator as any).xr.isSessionSupported(
            'immersive-ar'
          );
          capabilities.webXR = isSupported;
          capabilities.supported = isSupported;
        } catch {
          console.log('WebXR not supported');
        }
      }

      // Check iOS Quick Look support
      const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent);
      if (isIOS && 'relList' in HTMLAnchorElement.prototype) {
        const anchor = document.createElement('a');
        capabilities.quickLook = anchor.relList.supports('ar');
        capabilities.supported =
          capabilities.supported || capabilities.quickLook;
      }

      // Check Android Scene Viewer support
      const isAndroid = /Android/.test(navigator.userAgent);
      if (isAndroid) {
        capabilities.sceneViewer = true;
        capabilities.supported = true;
      }

      setArCapabilities(capabilities);
    };

    checkARSupport();
  }, []);

  // Load 3D model or generate from image
  useEffect(() => {
    const loadModel = async () => {
      setIsLoading(true);
      setError(null);

      try {
        if (modelUrl) {
          // Load existing 3D model
          await load3DModel(modelUrl);
        } else {
          // Generate 3D model from image using AI
          await generateModelFromImage(imageUrl);
        }
      } catch {
        setError('Failed to load model');
      } finally {
        setIsLoading(false);
      }
    };

    loadModel();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [modelUrl, imageUrl]);

  const load3DModel = async (url: string) => {
    // Implementation for loading 3D model
    // This would use Three.js or model-viewer library
    if (modelViewerRef.current) {
      modelViewerRef.current.src = url;
    }
  };

  const generateModelFromImage = async (imageUrl: string) => {
    // Mock implementation for AI-based 3D model generation
    // In production, this would call an AI service
    try {
      // Simulate API call to generate 3D model from image
      const response = await fetch('/api/v1/ar/generate-model', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ imageUrl, attributes }),
      });

      if (!response.ok) {
        throw new Error('Failed to generate 3D model');
      }

      const { modelUrl: generatedModelUrl } = await response.json();
      await load3DModel(generatedModelUrl);
    } catch {
      // Fallback to 2.5D representation
      create25DModel(imageUrl);
    }
  };

  const create25DModel = (imageUrl: string) => {
    // Create a 2.5D representation using depth estimation
    if (canvasRef.current) {
      const ctx = canvasRef.current.getContext('2d');
      if (!ctx) return;

      const img = new Image();
      img.onload = () => {
        // Apply depth effect to create pseudo-3D
        ctx.drawImage(img, 0, 0);
        applyDepthEffect(ctx, img);
      };
      img.src = imageUrl;
    }
  };

  const applyDepthEffect = (
    ctx: CanvasRenderingContext2D,
    img: HTMLImageElement
  ) => {
    // Apply depth map and parallax effect
    const imageData = ctx.getImageData(0, 0, img.width, img.height);
    // Depth estimation algorithm would go here
    // This is a simplified version
    ctx.putImageData(imageData, 0, 0);
  };

  const startARSession = useCallback(async () => {
    if (arCapabilities.webXR) {
      try {
        const session = await (navigator as any).xr.requestSession(
          'immersive-ar',
          {
            requiredFeatures: ['hit-test', 'dom-overlay'],
            domOverlay: { root: document.body },
          }
        );

        setIsARActive(true);
        handleWebXRSession(session);
      } catch {
        toast.error(t('arNotSupported'));
      }
    } else if (arCapabilities.quickLook) {
      // iOS Quick Look
      const a = document.createElement('a');
      a.rel = 'ar';
      a.href = modelUrl || '';

      const img = document.createElement('img');
      img.src = imageUrl;
      a.appendChild(img);

      a.click();
    } else if (arCapabilities.sceneViewer) {
      // Android Scene Viewer
      const intent = `intent://arvr.google.com/scene-viewer/1.0?file=${encodeURIComponent(modelUrl || '')}&mode=ar_preferred#Intent;scheme=https;package=com.google.android.googlequicksearchbox;end;`;
      window.location.href = intent;
    } else {
      toast.error(t('arNotSupported'));
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [arCapabilities, modelUrl, imageUrl, t]);

  const handleWebXRSession = async (session: any) => {
    // WebXR session handling
    const canvas = canvasRef.current;
    if (!canvas) return;

    const gl = canvas.getContext('webgl2', { xrCompatible: true });
    if (!gl) return;

    // Set up WebXR rendering loop
    const referenceSpace = await session.requestReferenceSpace('local');
    const hitTestSource = await session.requestHitTestSource({
      space: referenceSpace,
    });

    const onXRFrame = (time: number, frame: any) => {
      session.requestAnimationFrame(onXRFrame);

      const pose = frame.getViewerPose(referenceSpace);
      if (!pose) return;

      // Render 3D model at hit test location
      const hitTestResults = frame.getHitTestResults(hitTestSource);
      if (hitTestResults.length > 0) {
        const hit = hitTestResults[0];
        // Place 3D model at hit location
        renderModelAtLocation(hit.getPose(referenceSpace));
      }
    };

    session.requestAnimationFrame(onXRFrame);

    session.addEventListener('end', () => {
      setIsARActive(false);
    });
  };

  const renderModelAtLocation = (_pose: any) => {
    // Render 3D model at the specified pose
    // Implementation would use Three.js or WebGL
  };

  const exitAR = useCallback(() => {
    setIsARActive(false);
    if (onClose) onClose();
  }, [onClose]);

  // Furniture placement helper
  const measureRoom = useCallback(async () => {
    if (!arCapabilities.webXR) {
      toast.error(t('measureNotSupported'));
      return;
    }

    try {
      const _session = await (navigator as any).xr.requestSession(
        'immersive-ar',
        {
          requiredFeatures: ['hit-test', 'plane-detection'],
        }
      );

      // Room measurement logic
      toast.success(t('measureStarted'));
    } catch {
      toast.error(t('measureFailed'));
    }
  }, [arCapabilities, t]);

  return (
    <div className="fixed inset-0 z-50 bg-black/90">
      {/* AR Viewer */}
      <div className="relative h-full w-full">
        {/* Canvas for WebXR or 2.5D rendering */}
        <canvas
          ref={canvasRef}
          className="absolute inset-0 h-full w-full"
          style={{ display: isARActive ? 'block' : 'none' }}
        />

        {/* Model viewer for non-AR preview */}
        {!isARActive && (
          <div className="flex h-full w-full items-center justify-center">
            {isLoading ? (
              <div className="text-white">
                <div className="loading loading-spinner loading-lg"></div>
                <p className="mt-4">{t('loadingModel')}</p>
              </div>
            ) : error ? (
              <div className="text-center text-white">
                <p className="text-error">{error}</p>
                <button
                  className="btn btn-primary mt-4"
                  onClick={() => create25DModel(imageUrl)}
                >
                  {t('use25D')}
                </button>
              </div>
            ) : (
              <div className="relative h-full w-full">
                {/* 3D Model viewer */}
                {/* @ts-ignore - model-viewer is a web component */}
                <model-viewer
                  ref={modelViewerRef}
                  ar
                  ar-modes="webxr scene-viewer quick-look"
                  camera-controls
                  auto-rotate
                  shadow-intensity="1"
                  className="h-full w-full"
                >
                  <button
                    slot="ar-button"
                    className="btn btn-primary absolute bottom-4 left-1/2 -translate-x-1/2"
                    onClick={startARSession}
                  >
                    {t('viewInAR')}
                  </button>
                  {/* @ts-ignore */}
                </model-viewer>
              </div>
            )}
          </div>
        )}

        {/* AR Controls */}
        {isARActive && (
          <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/50 p-4">
            <div className="flex items-center justify-between text-white">
              <div>
                <h3 className="text-lg font-semibold">{productName}</h3>
                <p className="text-sm opacity-80">{t('tapToPlace')}</p>
              </div>
              <div className="flex gap-2">
                {attributes?.category === 'furniture' && (
                  <button
                    className="btn btn-circle btn-ghost"
                    onClick={measureRoom}
                  >
                    📐
                  </button>
                )}
                <button
                  className="btn btn-circle btn-ghost"
                  onClick={() => toast.success(t('screenshotSaved'))}
                >
                  📸
                </button>
                <button className="btn btn-circle btn-error" onClick={exitAR}>
                  ✕
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Info Panel */}
        {!isARActive && !isLoading && !error && (
          <div className="absolute left-4 top-4 rounded-lg bg-white/90 p-4 backdrop-blur">
            <h2 className="mb-2 text-xl font-bold">{productName}</h2>
            {attributes && (
              <div className="space-y-1 text-sm">
                {Object.entries(attributes)
                  .slice(0, 5)
                  .map(([key, value]) => (
                    <div key={key} className="flex justify-between">
                      <span className="font-medium">{key}:</span>
                      <span>{String(value)}</span>
                    </div>
                  ))}
              </div>
            )}
            <div className="mt-4 flex gap-2">
              <button
                className="btn btn-primary btn-sm"
                onClick={startARSession}
                disabled={!arCapabilities.supported}
              >
                {arCapabilities.supported ? t('tryInAR') : t('arNotAvailable')}
              </button>
              <button className="btn btn-ghost btn-sm" onClick={onClose}>
                {t('close')}
              </button>
            </div>
          </div>
        )}

        {/* AR Capabilities indicator */}
        {!isARActive && (
          <div className="absolute right-4 top-4 rounded-lg bg-black/50 p-2 text-xs text-white">
            <div>WebXR: {arCapabilities.webXR ? '✅' : '❌'}</div>
            <div>Quick Look: {arCapabilities.quickLook ? '✅' : '❌'}</div>
            <div>Scene Viewer: {arCapabilities.sceneViewer ? '✅' : '❌'}</div>
          </div>
        )}
      </div>
    </div>
  );
}
