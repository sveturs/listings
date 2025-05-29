import { Listing, ListingStatus, PaginatedResponse } from '@/types/listing';

export const mockListings: Listing[] = [
  {
    id: '1',
    title: 'Modern Apartment in City Center',
    description: 'Beautiful 2-bedroom apartment with stunning city views',
    price: 250000,
    currency: 'USD',
    location: 'Belgrade, Serbia',
    category: { id: '1', name: 'Real Estate', slug: 'real-estate' },
    attributes: [],
    images: [
      { id: '1', url: 'https://via.placeholder.com/400x300', order: 0 }
    ],
    status: ListingStatus.ACTIVE,
    userId: '1',
    user: { id: '1', name: 'John Doe', rating: 4.5 },
    viewCount: 156,
    favoriteCount: 12,
    createdAt: '2024-01-15T10:00:00Z',
    updatedAt: '2024-01-15T10:00:00Z',
  },
  {
    id: '2',
    title: 'Cozy Studio Near University',
    description: 'Perfect for students, fully furnished',
    price: 450,
    currency: 'EUR',
    location: 'Novi Sad, Serbia',
    category: { id: '1', name: 'Real Estate', slug: 'real-estate' },
    attributes: [],
    images: [
      { id: '2', url: 'https://via.placeholder.com/400x300', order: 0 }
    ],
    status: ListingStatus.ACTIVE,
    userId: '2',
    user: { id: '2', name: 'Jane Smith', rating: 5.0 },
    viewCount: 89,
    favoriteCount: 7,
    createdAt: '2024-01-14T10:00:00Z',
    updatedAt: '2024-01-14T10:00:00Z',
  },
  {
    id: '3',
    title: 'Luxury Villa with Pool',
    description: 'Spacious 5-bedroom villa with private pool and garden',
    price: 750000,
    currency: 'EUR',
    location: 'Belgrade, Serbia',
    category: { id: '1', name: 'Real Estate', slug: 'real-estate' },
    attributes: [],
    images: [
      { id: '3', url: 'https://via.placeholder.com/400x300', order: 0 }
    ],
    status: ListingStatus.ACTIVE,
    userId: '3',
    user: { id: '3', name: 'Mike Johnson', rating: 4.8 },
    viewCount: 234,
    favoriteCount: 28,
    createdAt: '2024-01-13T10:00:00Z',
    updatedAt: '2024-01-13T10:00:00Z',
  },
];

export function getMockListings(page = 1, pageSize = 12): PaginatedResponse<Listing> {
  const start = (page - 1) * pageSize;
  const end = start + pageSize;
  const items = mockListings.slice(start, end);
  
  return {
    items,
    total: mockListings.length,
    page,
    pageSize,
    totalPages: Math.ceil(mockListings.length / pageSize),
  };
}

export function getMockListingById(id: string): Listing | undefined {
  return mockListings.find(listing => listing.id === id);
}