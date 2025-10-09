'use client';

import { useState, useEffect } from 'react';
import { MarketplaceService } from '@/services/c2c';
import configManager from '@/config';
import { EnhancedListingCard } from '@/components/c2c/EnhancedListingCard';

export default function DebugMarketplacePage() {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchData() {
      try {
        const result = await MarketplaceService.search({
          sort_by: 'date_desc',
          page: 0,
          size: 20,
        });
        setData(result);
      } catch (error) {
        console.error('Error:', error);
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, []);

  if (loading) return <div>Loading...</div>;

  const listing177 = data?.data?.find((item: any) => item.id === 177);

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">
        Debug Marketplace - Listing 177
      </h1>

      {listing177 && (
        <div className="space-y-4">
          <div>
            <h2 className="font-bold">Raw Data:</h2>
            <pre className="bg-gray-100 p-2 rounded overflow-auto">
              {JSON.stringify(listing177, null, 2)}
            </pre>
          </div>

          <div>
            <h2 className="font-bold">Images Test:</h2>
            {listing177.images?.map((img: any, idx: number) => {
              const url = img.public_url;
              const builtUrl = configManager.buildImageUrl(url);
              return (
                <div key={idx} className="border p-2 mb-2">
                  <p>Original URL: {url}</p>
                  <p>Built URL: {builtUrl}</p>
                  <p>Direct test:</p>
                  {}
                  <img
                    src={builtUrl}
                    alt="Test"
                    className="w-32 h-32 object-cover"
                    onError={(e) => console.error('Image load error:', e)}
                    onLoad={() => console.log('Image loaded successfully')}
                  />
                </div>
              );
            })}
          </div>

          <div>
            <h2 className="font-bold">EnhancedListingCard Component:</h2>
            <div className="w-64">
              <EnhancedListingCard
                item={listing177}
                locale="en"
                viewMode="grid"
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
