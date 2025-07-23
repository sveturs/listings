'use client';

import { AnimatedSection } from '@/components/ui/AnimatedSection';

export default function BentoGridExamplePage() {
  return (
    <div className="container mx-auto p-4 max-w-7xl">
      <AnimatedSection animation="fadeIn">
        <h1 className="text-4xl font-bold mb-8">Bento Grid Layout</h1>
        <p className="text-lg text-base-content/70 mb-8">
          Modern grid layout inspired by Bento boxes for showcasing content
        </p>
      </AnimatedSection>

      {/* Main Bento Grid */}
      <AnimatedSection animation="slideUp" delay={0.1}>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-12">
          {/* Large Feature Card */}
          <div className="md:col-span-2 md:row-span-2">
            <div className="card bg-gradient-to-br from-primary to-secondary text-primary-content h-full">
              <div className="card-body">
                <h2 className="card-title text-2xl mb-4">Featured Property</h2>
                <div className="flex-1 flex items-center justify-center">
                  <div className="text-center">
                    <div className="text-6xl mb-4">🏡</div>
                    <p className="text-lg opacity-90">
                      Luxury Villa in Belgrade
                    </p>
                    <p className="text-3xl font-bold mt-2">€1,200/month</p>
                  </div>
                </div>
                <div className="card-actions justify-end">
                  <button className="btn btn-primary">View Details</button>
                </div>
              </div>
            </div>
          </div>

          {/* Stats Cards */}
          <div className="card bg-base-200">
            <div className="card-body">
              <div className="text-3xl mb-2">📊</div>
              <h3 className="font-semibold">Total Views</h3>
              <p className="text-2xl font-bold">12,543</p>
              <p className="text-sm text-success">+12% this week</p>
            </div>
          </div>

          <div className="card bg-base-200">
            <div className="card-body">
              <div className="text-3xl mb-2">💬</div>
              <h3 className="font-semibold">Messages</h3>
              <p className="text-2xl font-bold">48</p>
              <p className="text-sm text-warning">5 unread</p>
            </div>
          </div>

          {/* Wide Card */}
          <div className="md:col-span-2 card bg-accent text-accent-content">
            <div className="card-body">
              <h3 className="card-title">Quick Actions</h3>
              <div className="flex gap-2 flex-wrap">
                <button className="btn btn-sm">Post Listing</button>
                <button className="btn btn-sm">Search</button>
                <button className="btn btn-sm">Messages</button>
                <button className="btn btn-sm">Settings</button>
              </div>
            </div>
          </div>

          {/* Tall Card */}
          <div className="md:row-span-2 card bg-base-200">
            <div className="card-body">
              <h3 className="card-title mb-4">Recent Activity</h3>
              <div className="space-y-3 flex-1">
                {[1, 2, 3, 4].map((i) => (
                  <div key={i} className="flex items-center gap-3">
                    <div className="avatar placeholder">
                      <div className="bg-neutral text-neutral-content rounded-full w-8">
                        <span className="text-xs">U{i}</span>
                      </div>
                    </div>
                    <div className="flex-1">
                      <p className="text-sm font-medium">User {i}</p>
                      <p className="text-xs opacity-70">Viewed your listing</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Small Info Card */}
          <div className="card bg-info text-info-content">
            <div className="card-body">
              <div className="text-3xl mb-2">🔔</div>
              <h3 className="font-semibold">New Feature!</h3>
              <p className="text-sm">Try our new quick view</p>
            </div>
          </div>
        </div>
      </AnimatedSection>

      {/* Alternative Layouts */}
      <AnimatedSection animation="slideUp" delay={0.2}>
        <h2 className="text-2xl font-bold mb-6">Alternative Layouts</h2>

        {/* Layout 2: Dashboard Style */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
          <div className="md:col-span-2">
            <div className="card bg-base-100 shadow-xl h-full">
              <div className="card-body">
                <h3 className="card-title mb-4">Property Performance</h3>
                <div className="h-48 bg-base-200 rounded-lg flex items-center justify-center">
                  <p className="text-base-content/50">Chart Placeholder</p>
                </div>
              </div>
            </div>
          </div>
          <div className="space-y-4">
            <div className="card bg-success text-success-content">
              <div className="card-body p-4">
                <h4 className="font-semibold">Bookings</h4>
                <p className="text-2xl font-bold">23</p>
              </div>
            </div>
            <div className="card bg-warning text-warning-content">
              <div className="card-body p-4">
                <h4 className="font-semibold">Pending</h4>
                <p className="text-2xl font-bold">5</p>
              </div>
            </div>
          </div>
        </div>

        {/* Layout 3: Content Showcase */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="col-span-2 row-span-2 card bg-base-100 shadow-xl">
            <figure className="px-4 pt-4">
              <div className="bg-base-200 rounded-lg h-48 w-full flex items-center justify-center">
                <span className="text-6xl">🏠</span>
              </div>
            </figure>
            <div className="card-body">
              <h3 className="card-title">Main Feature</h3>
              <p>Your primary content goes here with more space</p>
            </div>
          </div>

          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="card bg-base-200">
              <div className="card-body p-4">
                <div className="text-2xl mb-2">🏢</div>
                <h4 className="font-medium">Item {i}</h4>
              </div>
            </div>
          ))}
        </div>
      </AnimatedSection>

      {/* Responsive Behavior */}
      <AnimatedSection animation="slideUp" delay={0.3}>
        <section className="card bg-base-100 shadow-xl mt-8">
          <div className="card-body">
            <h2 className="card-title mb-4">Responsive Behavior</h2>
            <p className="mb-4">
              The Bento Grid layout automatically adapts to different screen
              sizes:
            </p>
            <ul className="list-disc list-inside space-y-2">
              <li>Mobile: Single column layout</li>
              <li>Tablet: 2-3 columns with adjusted spans</li>
              <li>Desktop: Full 4-column grid with varied sizes</li>
            </ul>

            <div className="mt-4">
              <p className="text-sm text-base-content/70">
                Resize your browser window to see the responsive behavior in
                action!
              </p>
            </div>
          </div>
        </section>
      </AnimatedSection>

      {/* Code Example */}
      <AnimatedSection animation="slideUp" delay={0.4}>
        <section className="card bg-base-100 shadow-xl mt-8">
          <div className="card-body">
            <h2 className="card-title mb-4">Implementation</h2>
            <div className="mockup-code">
              <pre data-prefix="1">
                <code>{`<div className="grid grid-cols-1 md:grid-cols-4 gap-4">`}</code>
              </pre>
              <pre data-prefix="2">
                <code>{`  {/* Large card spanning 2 cols and 2 rows */}`}</code>
              </pre>
              <pre data-prefix="3">
                <code>{`  <div className="md:col-span-2 md:row-span-2">`}</code>
              </pre>
              <pre data-prefix="4">
                <code>{`    <div className="card h-full">...</div>`}</code>
              </pre>
              <pre data-prefix="5">
                <code>{`  </div>`}</code>
              </pre>
              <pre data-prefix="6">
                <code>{`  `}</code>
              </pre>
              <pre data-prefix="7">
                <code>{`  {/* Regular cards */}`}</code>
              </pre>
              <pre data-prefix="8">
                <code>{`  <div className="card">...</div>`}</code>
              </pre>
              <pre data-prefix="9">
                <code>{`  `}</code>
              </pre>
              <pre data-prefix="10">
                <code>{`  {/* Wide card spanning 2 columns */}`}</code>
              </pre>
              <pre data-prefix="11">
                <code>{`  <div className="md:col-span-2 card">...</div>`}</code>
              </pre>
              <pre data-prefix="12">
                <code>{`</div>`}</code>
              </pre>
            </div>
          </div>
        </section>
      </AnimatedSection>
    </div>
  );
}
