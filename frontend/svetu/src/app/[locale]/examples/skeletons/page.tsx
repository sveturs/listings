'use client';

import { useState } from 'react';
import { AnimatedSection } from '@/components/ui/AnimatedSection';

export default function SkeletonsExamplePage() {
  const [loading, setLoading] = useState(true);

  const toggleLoading = () => {
    setLoading(true);
    setTimeout(() => setLoading(false), 2000);
  };

  return (
    <div className="container mx-auto p-4 max-w-6xl">
      <AnimatedSection animation="fadeIn">
        <h1 className="text-4xl font-bold mb-8">Skeleton Loaders</h1>
        <p className="text-lg text-base-content/70 mb-8">
          Skeleton screens provide visual feedback while content is loading
        </p>
      </AnimatedSection>

      <AnimatedSection animation="slideUp" delay={0.1}>
        <div className="mb-8">
          <button className="btn btn-primary" onClick={toggleLoading}>
            Reload All Examples (2s delay)
          </button>
        </div>
      </AnimatedSection>

      {/* Card Skeleton */}
      <AnimatedSection animation="slideUp" delay={0.2}>
        <section className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">Card Skeleton</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  {loading ? (
                    <>
                      <div className="skeleton h-32 w-full mb-4"></div>
                      <div className="skeleton h-4 w-28 mb-2"></div>
                      <div className="skeleton h-4 w-full"></div>
                      <div className="skeleton h-4 w-full"></div>
                      <div className="skeleton h-4 w-3/4"></div>
                    </>
                  ) : (
                    <>
                      <div className="bg-primary/20 h-32 rounded-lg mb-4 flex items-center justify-center">
                        <span className="text-4xl">🏠</span>
                      </div>
                      <h3 className="font-semibold">Listing Title {i}</h3>
                      <p className="text-base-content/70">
                        This is a sample listing description that appears after
                        the skeleton loader finishes.
                      </p>
                    </>
                  )}
                </div>
              </div>
            ))}
          </div>
        </section>
      </AnimatedSection>

      {/* Text Content Skeleton */}
      <AnimatedSection animation="slideUp" delay={0.3}>
        <section className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h2 className="card-title mb-4">Text Content Skeleton</h2>
            {loading ? (
              <div className="space-y-2">
                <div className="skeleton h-6 w-48 mb-4"></div>
                <div className="skeleton h-4 w-full"></div>
                <div className="skeleton h-4 w-full"></div>
                <div className="skeleton h-4 w-3/4"></div>
                <div className="skeleton h-4 w-full"></div>
                <div className="skeleton h-4 w-5/6"></div>
              </div>
            ) : (
              <div>
                <h3 className="text-xl font-semibold mb-4">Article Title</h3>
                <p className="text-base-content/80">
                  Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed
                  do eiusmod tempor incididunt ut labore et dolore magna aliqua.
                  Ut enim ad minim veniam, quis nostrud exercitation ullamco
                  laboris.
                </p>
              </div>
            )}
          </div>
        </section>
      </AnimatedSection>

      {/* Avatar & User Info Skeleton */}
      <AnimatedSection animation="slideUp" delay={0.4}>
        <section className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h2 className="card-title mb-4">User Profile Skeleton</h2>
            <div className="flex items-center gap-4">
              {loading ? (
                <>
                  <div className="skeleton w-16 h-16 rounded-full shrink-0"></div>
                  <div className="flex-1">
                    <div className="skeleton h-4 w-32 mb-2"></div>
                    <div className="skeleton h-3 w-48"></div>
                  </div>
                </>
              ) : (
                <>
                  <div className="avatar">
                    <div className="w-16 rounded-full bg-primary text-primary-content">
                      <span className="text-2xl flex items-center justify-center h-full">
                        JD
                      </span>
                    </div>
                  </div>
                  <div>
                    <h3 className="font-semibold">John Doe</h3>
                    <p className="text-sm text-base-content/70">
                      john.doe@example.com
                    </p>
                  </div>
                </>
              )}
            </div>
          </div>
        </section>
      </AnimatedSection>

      {/* Table Skeleton */}
      <AnimatedSection animation="slideUp" delay={0.5}>
        <section className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h2 className="card-title mb-4">Table Skeleton</h2>
            <div className="overflow-x-auto">
              <table className="table">
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>Email</th>
                    <th>Status</th>
                    <th>Action</th>
                  </tr>
                </thead>
                <tbody>
                  {[1, 2, 3].map((i) => (
                    <tr key={i}>
                      {loading ? (
                        <>
                          <td>
                            <div className="skeleton h-4 w-24"></div>
                          </td>
                          <td>
                            <div className="skeleton h-4 w-32"></div>
                          </td>
                          <td>
                            <div className="skeleton h-4 w-16"></div>
                          </td>
                          <td>
                            <div className="skeleton h-4 w-20"></div>
                          </td>
                        </>
                      ) : (
                        <>
                          <td>User {i}</td>
                          <td>user{i}@example.com</td>
                          <td>
                            <span className="badge badge-success">Active</span>
                          </td>
                          <td>
                            <button className="btn btn-sm">Edit</button>
                          </td>
                        </>
                      )}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </section>
      </AnimatedSection>

      {/* Custom Pulse Animation */}
      <AnimatedSection animation="slideUp" delay={0.6}>
        <section className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">Custom Pulse Effects</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <h3 className="font-semibold mb-2">Wave Effect</h3>
                <div className="space-y-2">
                  <div className="skeleton h-4 w-full animate-pulse"></div>
                  <div
                    className="skeleton h-4 w-full animate-pulse"
                    style={{ animationDelay: '0.1s' }}
                  ></div>
                  <div
                    className="skeleton h-4 w-full animate-pulse"
                    style={{ animationDelay: '0.2s' }}
                  ></div>
                </div>
              </div>
              <div>
                <h3 className="font-semibold mb-2">Shimmer Effect</h3>
                <div className="relative overflow-hidden">
                  <div className="skeleton h-32 w-full"></div>
                  <div className="absolute inset-0 -translate-x-full animate-[shimmer_2s_infinite]">
                    <div className="h-full w-1/2 bg-gradient-to-r from-transparent via-white/20 to-transparent"></div>
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
