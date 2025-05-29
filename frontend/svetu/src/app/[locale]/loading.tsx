export default function Loading() {
  return (
    <div className="min-h-screen animate-pulse">
      <div className="py-20 text-center">
        <div className="h-10 w-64 bg-gray-300 rounded mx-auto mb-4"></div>
        <div className="h-6 w-48 bg-gray-300 rounded mx-auto mb-8"></div>
        <div className="h-12 w-40 bg-gray-300 rounded mx-auto"></div>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-12">
        {[1, 2, 3].map((id) => (
          <div key={id} className="border rounded-lg p-6">
            <div className="h-6 w-32 bg-gray-300 rounded mb-2"></div>
            <div className="h-4 w-48 bg-gray-300 rounded mb-4"></div>
            <div className="h-4 w-24 bg-gray-300 rounded"></div>
          </div>
        ))}
      </div>
    </div>
  );
}