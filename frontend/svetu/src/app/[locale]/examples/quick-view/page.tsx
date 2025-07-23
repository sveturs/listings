'use client';

import { useState } from 'react';
import { AnimatedSection } from '@/components/ui/AnimatedSection';
import {
  XMarkIcon,
  HeartIcon,
  ShareIcon,
  MapPinIcon,
  CurrencyEuroIcon,
  CalendarIcon,
  UserIcon,
} from '@heroicons/react/24/outline';
import { HeartIcon as HeartIconSolid } from '@heroicons/react/24/solid';

interface Listing {
  id: number;
  title: string;
  price: number;
  location: string;
  type: string;
  rating: number;
  images: number;
}

export default function QuickViewExamplePage() {
  const [selectedListing, setSelectedListing] = useState<Listing | null>(null);
  const [isFavorite, setIsFavorite] = useState(false);

  const listings: Listing[] = [
    {
      id: 1,
      title: 'Cozy Studio in City Center',
      price: 450,
      location: 'Belgrade, Serbia',
      type: 'Studio',
      rating: 4.5,
      images: 5,
    },
    {
      id: 2,
      title: 'Modern 2BR Apartment',
      price: 750,
      location: 'Novi Sad, Serbia',
      type: '2 Bedroom',
      rating: 4.8,
      images: 8,
    },
    {
      id: 3,
      title: 'Luxury Penthouse',
      price: 1200,
      location: 'Belgrade, Serbia',
      type: 'Penthouse',
      rating: 5.0,
      images: 12,
    },
    {
      id: 4,
      title: 'Student Room Near University',
      price: 250,
      location: 'Belgrade, Serbia',
      type: 'Room',
      rating: 4.2,
      images: 3,
    },
  ];

  const openQuickView = (listing: Listing) => {
    setSelectedListing(listing);
    setIsFavorite(false);
  };

  const closeQuickView = () => {
    setSelectedListing(null);
    setIsFavorite(false);
  };

  return (
    <div className="container mx-auto p-4 max-w-6xl">
      <AnimatedSection animation="fadeIn">
        <h1 className="text-4xl font-bold mb-8">Quick View Modal</h1>
        <p className="text-lg text-base-content/70 mb-8">
          Preview listings quickly without leaving the current page
        </p>
      </AnimatedSection>

      {/* Listing Grid */}
      <AnimatedSection animation="slideUp" delay={0.1}>
        <section className="mb-8">
          <h2 className="text-2xl font-bold mb-6">
            Click any listing for quick preview
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {listings.map((listing, index) => (
              <AnimatedSection
                key={listing.id}
                animation="slideUp"
                delay={0.1 * index}
              >
                <div
                  className="card bg-base-100 shadow-xl cursor-pointer hover:shadow-2xl transition-shadow"
                  onClick={() => openQuickView(listing)}
                >
                  <figure className="px-4 pt-4">
                    <div className="bg-base-200 rounded-lg h-48 w-full flex items-center justify-center">
                      <span className="text-6xl">🏠</span>
                    </div>
                  </figure>
                  <div className="card-body">
                    <h3 className="card-title">{listing.title}</h3>
                    <p className="text-base-content/70">{listing.location}</p>
                    <div className="flex justify-between items-center mt-4">
                      <span className="text-2xl font-bold">
                        €{listing.price}/mo
                      </span>
                      <div className="badge badge-ghost">{listing.type}</div>
                    </div>
                    <div className="flex items-center gap-2 mt-2">
                      <div className="rating rating-sm">
                        {[1, 2, 3, 4, 5].map((star) => (
                          <input
                            key={star}
                            type="radio"
                            className="mask mask-star-2 bg-orange-400"
                            checked={star <= Math.floor(listing.rating)}
                            readOnly
                          />
                        ))}
                      </div>
                      <span className="text-sm text-base-content/70">
                        ({listing.rating})
                      </span>
                    </div>
                  </div>
                </div>
              </AnimatedSection>
            ))}
          </div>
        </section>
      </AnimatedSection>

      {/* Quick View Modal */}
      {selectedListing && (
        <div className="modal modal-open">
          <div className="modal-box max-w-4xl">
            {/* Header */}
            <div className="flex justify-between items-start mb-4">
              <h3 className="font-bold text-xl">{selectedListing.title}</h3>
              <button
                className="btn btn-sm btn-circle btn-ghost"
                onClick={closeQuickView}
              >
                <XMarkIcon className="h-5 w-5" />
              </button>
            </div>

            {/* Image Gallery */}
            <div className="grid grid-cols-2 gap-2 mb-4">
              <div className="col-span-2 md:col-span-1 bg-base-200 rounded-lg h-64 flex items-center justify-center">
                <span className="text-8xl">🏠</span>
              </div>
              <div className="grid grid-cols-2 gap-2">
                {[1, 2, 3, 4].map((i) => (
                  <div
                    key={i}
                    className="bg-base-200 rounded-lg h-30 flex items-center justify-center"
                  >
                    <span className="text-4xl">📷</span>
                  </div>
                ))}
              </div>
            </div>

            {/* Quick Info */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <div className="mb-4">
                  <h4 className="font-semibold mb-2">Details</h4>
                  <div className="space-y-2">
                    <div className="flex items-center gap-2">
                      <MapPinIcon className="h-5 w-5 text-base-content/50" />
                      <span>{selectedListing.location}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <CurrencyEuroIcon className="h-5 w-5 text-base-content/50" />
                      <span className="font-bold">
                        €{selectedListing.price}/month
                      </span>
                    </div>
                    <div className="flex items-center gap-2">
                      <CalendarIcon className="h-5 w-5 text-base-content/50" />
                      <span>Available immediately</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <UserIcon className="h-5 w-5 text-base-content/50" />
                      <span>Posted by John Doe</span>
                    </div>
                  </div>
                </div>

                <div className="mb-4">
                  <h4 className="font-semibold mb-2">Features</h4>
                  <div className="flex flex-wrap gap-2">
                    <div className="badge badge-outline">WiFi</div>
                    <div className="badge badge-outline">Parking</div>
                    <div className="badge badge-outline">Pet Friendly</div>
                    <div className="badge badge-outline">Furnished</div>
                  </div>
                </div>
              </div>

              <div>
                <h4 className="font-semibold mb-2">Description</h4>
                <p className="text-base-content/80 mb-4">
                  Beautiful and modern accommodation in the heart of the city.
                  Perfect for students and young professionals. Close to public
                  transport, shops, and restaurants. All utilities included in
                  the price.
                </p>

                <div className="flex gap-2">
                  <button
                    className="btn btn-primary flex-1"
                    onClick={() => alert('Contact feature would open here')}
                  >
                    Contact Owner
                  </button>
                  <button
                    className={`btn btn-square ${isFavorite ? 'btn-error' : 'btn-ghost'}`}
                    onClick={() => setIsFavorite(!isFavorite)}
                  >
                    {isFavorite ? (
                      <HeartIconSolid className="h-5 w-5" />
                    ) : (
                      <HeartIcon className="h-5 w-5" />
                    )}
                  </button>
                  <button className="btn btn-square btn-ghost">
                    <ShareIcon className="h-5 w-5" />
                  </button>
                </div>
              </div>
            </div>

            <div className="modal-action">
              <button className="btn" onClick={closeQuickView}>
                Close
              </button>
              <button className="btn btn-primary">View Full Details</button>
            </div>
          </div>
          <div className="modal-backdrop" onClick={closeQuickView}></div>
        </div>
      )}

      {/* Features */}
      <AnimatedSection animation="slideUp" delay={0.3}>
        <section className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h2 className="card-title mb-4">Quick View Features</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <h3 className="font-semibold mb-2">⚡ Instant Preview</h3>
                <p className="text-base-content/70">
                  View essential listing information without page navigation
                </p>
              </div>
              <div>
                <h3 className="font-semibold mb-2">📸 Image Gallery</h3>
                <p className="text-base-content/70">
                  Browse through property images in a compact layout
                </p>
              </div>
              <div>
                <h3 className="font-semibold mb-2">💬 Quick Actions</h3>
                <p className="text-base-content/70">
                  Contact owner, save to favorites, or share directly from
                  preview
                </p>
              </div>
              <div>
                <h3 className="font-semibold mb-2">📱 Mobile Optimized</h3>
                <p className="text-base-content/70">
                  Responsive design works perfectly on all devices
                </p>
              </div>
            </div>
          </div>
        </section>
      </AnimatedSection>

      {/* Implementation */}
      <AnimatedSection animation="slideUp" delay={0.4}>
        <section className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">Implementation</h2>
            <div className="mockup-code">
              <pre data-prefix="1">
                <code>{`// Open modal with listing data`}</code>
              </pre>
              <pre data-prefix="2">
                <code>{`const openQuickView = (listing) => {`}</code>
              </pre>
              <pre data-prefix="3">
                <code>{`  setSelectedListing(listing);`}</code>
              </pre>
              <pre data-prefix="4">
                <code>{`};`}</code>
              </pre>
              <pre data-prefix="5">
                <code>{``}</code>
              </pre>
              <pre data-prefix="6">
                <code>{`// Modal component`}</code>
              </pre>
              <pre data-prefix="7">
                <code>{`{selectedListing && (`}</code>
              </pre>
              <pre data-prefix="8">
                <code>{`  <div className="modal modal-open">`}</code>
              </pre>
              <pre data-prefix="9">
                <code>{`    <div className="modal-box">`}</code>
              </pre>
              <pre data-prefix="10">
                <code>{`      {/* Quick view content */}`}</code>
              </pre>
              <pre data-prefix="11">
                <code>{`    </div>`}</code>
              </pre>
              <pre data-prefix="12">
                <code>{`  </div>`}</code>
              </pre>
              <pre data-prefix="13">
                <code>{`)}`}</code>
              </pre>
            </div>
          </div>
        </section>
      </AnimatedSection>
    </div>
  );
}
