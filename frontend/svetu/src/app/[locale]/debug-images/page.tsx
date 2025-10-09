'use client';

import { useState, useEffect } from 'react';
import { MarketplaceService } from '@/services/c2c';
import { C2CItem } from '@/types/c2c';
import configManager from '@/config';
import SafeImage from '@/components/SafeImage';

export default function DebugImagesPage() {
  const [items, setItems] = useState<C2CItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const response = await MarketplaceService.search({
        page: 0,
        size: 5,
      });
      setItems(response.data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  const getImageUrl = (publicUrl?: string) => {
    if (!publicUrl) return null;
    return configManager.buildImageUrl(publicUrl);
  };

  if (loading) return <div className="p-8">Loading...</div>;
  if (error) return <div className="p-8 text-red-500">Error: {error}</div>;

  return (
    <div className="container mx-auto p-8">
      <h1 className="text-2xl font-bold mb-6">Debug: Images</h1>

      <div className="mb-8">
        <h2 className="text-lg font-semibold mb-4">Config Info:</h2>
        <div className="bg-gray-100 p-4 rounded">
          <p>
            <strong>API URL:</strong> {configManager.getApiUrl()}
          </p>
          <p>
            <strong>Minio URL:</strong> {configManager.getMinioUrl()}
          </p>
          <p>
            <strong>Image Base URL:</strong> {configManager.getImageBaseUrl()}
          </p>
          <p>
            <strong>Is Production:</strong>{' '}
            {configManager.isProduction().toString()}
          </p>
        </div>
      </div>

      <div className="space-y-8">
        {items.map((item) => (
          <div key={item.id} className="border p-4 rounded">
            <h3 className="font-semibold mb-2">{item.title}</h3>

            <div className="mb-4">
              <h4 className="font-medium">Raw Images Data:</h4>
              <pre className="bg-gray-100 p-2 text-xs overflow-x-auto">
                {JSON.stringify(item.images, null, 2)}
              </pre>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {item.images?.map((image, index) => {
                const fullUrl = getImageUrl(image.public_url);
                return (
                  <div key={image.id} className="border p-2">
                    <p className="text-xs mb-2">
                      <strong>Image {index + 1}:</strong>{' '}
                      {image.is_main ? '(MAIN)' : ''}
                    </p>
                    <p className="text-xs mb-2">
                      <strong>Raw URL:</strong> {image.public_url}
                    </p>
                    <p className="text-xs mb-2">
                      <strong>Built URL:</strong> {fullUrl}
                    </p>

                    <div className="relative w-full h-48 bg-gray-200 border">
                      <p className="text-xs mb-1">
                        <strong>Safe URL Check:</strong>{' '}
                        {fullUrl ? 'PASSED' : 'BLOCKED'}
                      </p>
                      <SafeImage
                        src={fullUrl}
                        alt={`${item.title} - Image ${index + 1}`}
                        fill
                        className="object-cover"
                        fallback={
                          <div className="w-full h-full flex items-center justify-center text-gray-500 border border-red-200">
                            <span>No Image (Fallback)</span>
                          </div>
                        }
                        onError={(e) => {
                          console.error('Image load error:', fullUrl, e);
                        }}
                      />
                    </div>
                  </div>
                );
              })}
            </div>

            {(!item.images || item.images.length === 0) && (
              <p className="text-gray-500">No images for this item</p>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
