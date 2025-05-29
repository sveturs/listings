export default function Loading() {
  return (
    <div className="max-w-4xl mx-auto animate-pulse">
      <div className="h-6 w-32 bg-gray-300 rounded mb-6"></div>
      <div className="h-10 w-64 bg-gray-300 rounded mb-6"></div>
      <div className="bg-gray-200 h-64 rounded-lg mb-6"></div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <div className="border rounded-lg p-4">
          <div className="h-4 w-16 bg-gray-300 rounded mb-2"></div>
          <div className="h-8 w-32 bg-gray-300 rounded"></div>
        </div>
        <div className="border rounded-lg p-4">
          <div className="h-4 w-20 bg-gray-300 rounded mb-2"></div>
          <div className="h-6 w-40 bg-gray-300 rounded"></div>
        </div>
      </div>
      <div className="border rounded-lg p-6">
        <div className="h-4 w-24 bg-gray-300 rounded mb-3"></div>
        <div className="space-y-2">
          <div className="h-4 w-full bg-gray-300 rounded"></div>
          <div className="h-4 w-full bg-gray-300 rounded"></div>
          <div className="h-4 w-3/4 bg-gray-300 rounded"></div>
        </div>
      </div>
    </div>
  );
}