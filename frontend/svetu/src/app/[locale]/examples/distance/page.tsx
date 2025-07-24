'use client';

import { useState } from 'react';
import { AnimatedSection } from '@/components/ui/AnimatedSection';
import { MapPinIcon } from '@heroicons/react/24/outline';

export default function DistanceExamplePage() {
  const [unit, setUnit] = useState<'km' | 'mi'>('km');

  const distances = [
    { meters: 250, walkTime: 3 },
    { meters: 500, walkTime: 6 },
    { meters: 1000, walkTime: 12 },
    { meters: 2500, walkTime: 30 },
    { meters: 5000, walkTime: 60 },
    { meters: 10000, walkTime: 120 },
  ];

  const formatDistance = (meters: number) => {
    if (unit === 'km') {
      return meters < 1000 ? `${meters}m` : `${(meters / 1000).toFixed(1)}km`;
    } else {
      const miles = meters * 0.000621371;
      return miles < 0.5
        ? `${Math.round(meters * 3.28084)}ft`
        : `${miles.toFixed(1)}mi`;
    }
  };

  const getDistanceColor = (meters: number) => {
    if (meters <= 500) return 'text-success';
    if (meters <= 1000) return 'text-info';
    if (meters <= 2500) return 'text-warning';
    return 'text-error';
  };

  const getDistanceWidth = (meters: number) => {
    const maxDistance = 10000;
    return `${(meters / maxDistance) * 100}%`;
  };

  return (
    <div className="container mx-auto p-4 max-w-6xl">
      <AnimatedSection animation="fadeIn">
        <h1 className="text-4xl font-bold mb-8">Distance Visualization</h1>
        <p className="text-lg text-base-content/70 mb-8">
          Visual representation of distances with walking time estimates
        </p>
      </AnimatedSection>

      {/* Unit Toggle */}
      <AnimatedSection animation="slideUp" delay={0.1}>
        <div className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h2 className="card-title mb-4">Unit Selection</h2>
            <div className="flex gap-2">
              <button
                className={`btn ${unit === 'km' ? 'btn-primary' : 'btn-ghost'}`}
                onClick={() => setUnit('km')}
              >
                Metric (km/m)
              </button>
              <button
                className={`btn ${unit === 'mi' ? 'btn-primary' : 'btn-ghost'}`}
                onClick={() => setUnit('mi')}
              >
                Imperial (mi/ft)
              </button>
            </div>
          </div>
        </div>
      </AnimatedSection>

      {/* Distance Cards */}
      <AnimatedSection animation="slideUp" delay={0.2}>
        <section className="mb-8">
          <h2 className="text-2xl font-bold mb-6">Distance Cards</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {distances.map((item, index) => (
              <AnimatedSection
                key={item.meters}
                animation="slideUp"
                delay={0.1 * index}
              >
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <div className="flex items-start justify-between">
                      <div>
                        <h3
                          className={`text-3xl font-bold ${getDistanceColor(item.meters)}`}
                        >
                          {formatDistance(item.meters)}
                        </h3>
                        <p className="text-base-content/70">
                          ~{item.walkTime} min walk
                        </p>
                      </div>
                      <div
                        className={`text-4xl ${getDistanceColor(item.meters)}`}
                      >
                        {item.meters <= 500 && '🚶‍♂️'}
                        {item.meters > 500 && item.meters <= 1000 && '🚴‍♂️'}
                        {item.meters > 1000 && item.meters <= 2500 && '🚗'}
                        {item.meters > 2500 && '🚌'}
                      </div>
                    </div>
                    <progress
                      className={`progress ${
                        item.meters <= 500
                          ? 'progress-success'
                          : item.meters <= 1000
                            ? 'progress-info'
                            : item.meters <= 2500
                              ? 'progress-warning'
                              : 'progress-error'
                      } w-full mt-4`}
                      value={item.meters}
                      max="10000"
                    ></progress>
                  </div>
                </div>
              </AnimatedSection>
            ))}
          </div>
        </section>
      </AnimatedSection>

      {/* Visual Distance Bar */}
      <AnimatedSection animation="slideUp" delay={0.3}>
        <section className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h2 className="card-title mb-6">Distance Comparison</h2>
            <div className="space-y-4">
              {distances.map((item) => (
                <div key={item.meters}>
                  <div className="flex justify-between mb-2">
                    <span className="font-medium">
                      {formatDistance(item.meters)}
                    </span>
                    <span className="text-sm text-base-content/70">
                      ~{item.walkTime} min
                    </span>
                  </div>
                  <div className="relative bg-base-200 rounded-full h-8 overflow-hidden">
                    <div
                      className={`absolute inset-y-0 left-0 ${
                        item.meters <= 500
                          ? 'bg-success'
                          : item.meters <= 1000
                            ? 'bg-info'
                            : item.meters <= 2500
                              ? 'bg-warning'
                              : 'bg-error'
                      } rounded-full transition-all duration-1000 ease-out flex items-center justify-end pr-2`}
                      style={{ width: getDistanceWidth(item.meters) }}
                    >
                      <MapPinIcon className="h-5 w-5 text-white" />
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </section>
      </AnimatedSection>

      {/* Circular Distance Indicator */}
      <AnimatedSection animation="slideUp" delay={0.4}>
        <section className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h2 className="card-title mb-6">Circular Indicators</h2>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
              {distances.slice(0, 4).map((item) => (
                <div key={item.meters} className="text-center">
                  <div className="relative inline-flex">
                    <div
                      className={`radial-progress ${getDistanceColor(item.meters)}`}
                      style={{ '--value': (item.meters / 10000) * 100 } as any}
                      role="progressbar"
                    >
                      <span className="text-xl font-bold">
                        {formatDistance(item.meters)}
                      </span>
                    </div>
                  </div>
                  <p className="text-sm text-base-content/70 mt-2">
                    {item.walkTime} min
                  </p>
                </div>
              ))}
            </div>
          </div>
        </section>
      </AnimatedSection>

      {/* Map-style Distance */}
      <AnimatedSection animation="slideUp" delay={0.5}>
        <section className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h2 className="card-title mb-6">Map-style Visualization</h2>
            <div className="relative bg-base-200 rounded-lg p-8 h-64 flex items-center justify-center">
              <div className="absolute inset-0 flex items-center justify-center">
                {[100, 200, 300, 400].map((radius) => (
                  <div
                    key={radius}
                    className="absolute border border-base-content/20 rounded-full"
                    style={{
                      width: `${radius}px`,
                      height: `${radius}px`,
                    }}
                  />
                ))}
              </div>
              <div className="relative z-10">
                <div className="bg-primary text-primary-content rounded-full w-4 h-4"></div>
                <p className="absolute top-6 left-1/2 -translate-x-1/2 whitespace-nowrap text-sm font-medium">
                  Your Location
                </p>
              </div>
              {/* Sample points */}
              <div className="absolute top-1/4 left-1/3">
                <div className="bg-success rounded-full w-3 h-3"></div>
                <p className="text-xs mt-1">250m</p>
              </div>
              <div className="absolute bottom-1/3 right-1/4">
                <div className="bg-warning rounded-full w-3 h-3"></div>
                <p className="text-xs mt-1">1.2km</p>
              </div>
            </div>
          </div>
        </section>
      </AnimatedSection>

      {/* Usage Examples */}
      <AnimatedSection animation="slideUp" delay={0.6}>
        <section className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">Usage in Listings</h2>
            <div className="space-y-4">
              {/* Example listing card */}
              <div className="card bg-base-200">
                <div className="card-body">
                  <h3 className="font-semibold mb-2">
                    Modern Apartment in City Center
                  </h3>
                  <div className="flex items-center gap-4 text-sm">
                    <div className="flex items-center gap-1">
                      <MapPinIcon className="h-4 w-4" />
                      <span className="text-success font-medium">
                        450m from center
                      </span>
                    </div>
                    <span className="text-base-content/70">• 5 min walk</span>
                  </div>
                </div>
              </div>

              <div className="card bg-base-200">
                <div className="card-body">
                  <h3 className="font-semibold mb-2">
                    Cozy Studio Near University
                  </h3>
                  <div className="flex items-center gap-4 text-sm">
                    <div className="flex items-center gap-1">
                      <MapPinIcon className="h-4 w-4" />
                      <span className="text-warning font-medium">
                        2.3km from campus
                      </span>
                    </div>
                    <span className="text-base-content/70">• 28 min walk</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
      </AnimatedSection>
    </div>
  );
}
