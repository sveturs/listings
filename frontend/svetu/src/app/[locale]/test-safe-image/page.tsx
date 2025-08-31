'use client';

import SafeImage from '@/components/SafeImage';
import { getSafeImageUrl } from '@/utils/imageUtils';
import configManager from '@/config';

export default function TestSafeImagePage() {
  const testUrl = '/listings/177/1750065511512010029.jpg';
  const safeUrl = getSafeImageUrl(testUrl);
  const builtUrl = configManager.buildImageUrl(testUrl);

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Test SafeImage</h1>

      <div className="space-y-4">
        <div>
          <h2 className="font-bold">Test URL:</h2>
          <p>{testUrl}</p>
        </div>

        <div>
          <h2 className="font-bold">Safe URL result:</h2>
          <p>{safeUrl || 'null'}</p>
        </div>

        <div>
          <h2 className="font-bold">Built URL (configManager):</h2>
          <p>{builtUrl}</p>
        </div>

        <div>
          <h2 className="font-bold">Direct img tag with relative URL:</h2>
          <img
            src={testUrl}
            alt="Direct test relative"
            className="w-32 h-32 object-cover border"
            onError={(e) => console.error('Direct img error (relative):', e)}
            onLoad={() => console.log('Direct img loaded (relative)')}
          />
        </div>

        <div>
          <h2 className="font-bold">Direct img tag with built URL:</h2>
          <img
            src={builtUrl}
            alt="Direct test built"
            className="w-32 h-32 object-cover border"
            onError={(e) => console.error('Direct img error (built):', e)}
            onLoad={() => console.log('Direct img loaded (built)')}
          />
        </div>

        <div>
          <h2 className="font-bold">SafeImage component with relative URL:</h2>
          <div className="w-32 h-32 relative border">
            <SafeImage
              src={testUrl}
              alt="SafeImage test relative"
              fill
              className="object-cover"
            />
          </div>
        </div>

        <div>
          <h2 className="font-bold">SafeImage component with built URL:</h2>
          <div className="w-32 h-32 relative border">
            <SafeImage
              src={builtUrl}
              alt="SafeImage test built"
              fill
              className="object-cover"
            />
          </div>
        </div>
      </div>
    </div>
  );
}
