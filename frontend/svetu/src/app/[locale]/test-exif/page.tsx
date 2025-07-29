'use client';

import { useState } from 'react';
import { extractExifData } from '@/utils/exifUtils';

export default function TestExifPage() {
  const [exifData, setExifData] = useState<any>(null);
  const [error, setError] = useState<string>('');

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      console.log('Testing EXIF extraction for file:', file.name);
      const data = await extractExifData(file);
      setExifData(data);
      setError('');
    } catch (err) {
      console.error('EXIF extraction error:', err);
      setError(String(err));
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">EXIF Data Test</h1>

      <div className="mb-4">
        <input
          type="file"
          accept="image/*"
          onChange={handleFileChange}
          className="file-input file-input-bordered w-full max-w-xs"
        />
      </div>

      {error && (
        <div className="alert alert-error mb-4">
          <span>{error}</span>
        </div>
      )}

      {exifData && (
        <div className="bg-base-200 p-4 rounded">
          <h2 className="text-xl font-semibold mb-2">Extracted EXIF Data:</h2>
          <pre className="overflow-auto">
            {JSON.stringify(exifData, null, 2)}
          </pre>
        </div>
      )}
    </div>
  );
}
