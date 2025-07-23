'use client';

import { useState } from 'react';
import { AnimatedSection } from '@/components/ui/AnimatedSection';
import {
  HomeIcon,
  MagnifyingGlassIcon,
  PlusCircleIcon,
  ChatBubbleLeftIcon,
  UserIcon,
} from '@heroicons/react/24/outline';
import {
  HomeIcon as HomeIconSolid,
  MagnifyingGlassIcon as MagnifyingGlassIconSolid,
  PlusCircleIcon as PlusCircleIconSolid,
  ChatBubbleLeftIcon as ChatBubbleLeftIconSolid,
  UserIcon as UserIconSolid,
} from '@heroicons/react/24/solid';

export default function NavigationExamplePage() {
  const [activeTab, setActiveTab] = useState('home');
  const [theme, setTheme] = useState<'light' | 'dark'>('light');

  const navItems = [
    { id: 'home', label: 'Home', Icon: HomeIcon, IconActive: HomeIconSolid },
    {
      id: 'search',
      label: 'Search',
      Icon: MagnifyingGlassIcon,
      IconActive: MagnifyingGlassIconSolid,
    },
    {
      id: 'create',
      label: 'Create',
      Icon: PlusCircleIcon,
      IconActive: PlusCircleIconSolid,
    },
    {
      id: 'messages',
      label: 'Messages',
      Icon: ChatBubbleLeftIcon,
      IconActive: ChatBubbleLeftIconSolid,
    },
    {
      id: 'profile',
      label: 'Profile',
      Icon: UserIcon,
      IconActive: UserIconSolid,
    },
  ];

  return (
    <div className="container mx-auto p-4 max-w-6xl min-h-screen">
      <AnimatedSection animation="fadeIn">
        <h1 className="text-4xl font-bold mb-8">Mobile Navigation Examples</h1>
        <p className="text-lg text-base-content/70 mb-8">
          Responsive bottom navigation optimized for mobile devices
        </p>
      </AnimatedSection>

      {/* Phone Mockup */}
      <AnimatedSection animation="slideUp" delay={0.1}>
        <div className="flex justify-center mb-8">
          <div className="mockup-phone">
            <div className="camera"></div>
            <div className="display">
              <div
                className={`artboard artboard-demo phone-1 ${theme === 'dark' ? 'dark' : ''}`}
              >
                {/* Phone Screen Content */}
                <div className="h-full flex flex-col bg-base-100">
                  {/* Header */}
                  <div className="navbar bg-base-200 px-4">
                    <div className="flex-1">
                      <span className="text-xl font-bold">Sve Tu</span>
                    </div>
                    <button
                      className="btn btn-ghost btn-circle"
                      onClick={() =>
                        setTheme(theme === 'light' ? 'dark' : 'light')
                      }
                    >
                      {theme === 'light' ? '🌙' : '☀️'}
                    </button>
                  </div>

                  {/* Content */}
                  <div className="flex-1 p-4 overflow-y-auto">
                    <div className="space-y-4">
                      {activeTab === 'home' && (
                        <div>
                          <h2 className="text-2xl font-bold mb-4">Home</h2>
                          <div className="space-y-3">
                            {[1, 2, 3].map((i) => (
                              <div key={i} className="card bg-base-200">
                                <div className="card-body p-4">
                                  <h3 className="font-semibold">Listing {i}</h3>
                                  <p className="text-sm">
                                    Sample content for home feed
                                  </p>
                                </div>
                              </div>
                            ))}
                          </div>
                        </div>
                      )}
                      {activeTab === 'search' && (
                        <div>
                          <h2 className="text-2xl font-bold mb-4">Search</h2>
                          <input
                            type="text"
                            placeholder="Search listings..."
                            className="input input-bordered w-full"
                          />
                        </div>
                      )}
                      {activeTab === 'create' && (
                        <div>
                          <h2 className="text-2xl font-bold mb-4">
                            Create Listing
                          </h2>
                          <button className="btn btn-primary w-full">
                            + New Listing
                          </button>
                        </div>
                      )}
                      {activeTab === 'messages' && (
                        <div>
                          <h2 className="text-2xl font-bold mb-4">Messages</h2>
                          <p className="text-base-content/70">
                            No new messages
                          </p>
                        </div>
                      )}
                      {activeTab === 'profile' && (
                        <div>
                          <h2 className="text-2xl font-bold mb-4">Profile</h2>
                          <div className="avatar mb-4">
                            <div className="w-24 rounded-full bg-primary text-primary-content">
                              <span className="text-3xl flex items-center justify-center h-full">
                                JD
                              </span>
                            </div>
                          </div>
                          <p className="font-semibold">John Doe</p>
                        </div>
                      )}
                    </div>
                  </div>

                  {/* Bottom Navigation */}
                  <div className="btm-nav bg-base-200">
                    {navItems.map(({ id, label, Icon, IconActive }) => {
                      const isActive = activeTab === id;
                      const IconComponent = isActive ? IconActive : Icon;

                      return (
                        <button
                          key={id}
                          className={isActive ? 'active' : ''}
                          onClick={() => setActiveTab(id)}
                        >
                          <IconComponent className="h-5 w-5" />
                          <span className="btm-nav-label text-xs">{label}</span>
                        </button>
                      );
                    })}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </AnimatedSection>

      {/* Navigation Variations */}
      <AnimatedSection animation="slideUp" delay={0.2}>
        <section className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h2 className="card-title mb-4">Navigation Styles</h2>

            {/* Style 1: Default */}
            <div className="mb-6">
              <h3 className="font-semibold mb-2">Default Style</h3>
              <div className="btm-nav btm-nav-lg relative h-16 bg-base-200 rounded-lg">
                {navItems.slice(0, 4).map(({ id, label, Icon }) => (
                  <button key={id} className={id === 'home' ? 'active' : ''}>
                    <Icon className="h-5 w-5" />
                    <span className="btm-nav-label">{label}</span>
                  </button>
                ))}
              </div>
            </div>

            {/* Style 2: With Badge */}
            <div className="mb-6">
              <h3 className="font-semibold mb-2">With Notifications</h3>
              <div className="btm-nav btm-nav-lg relative h-16 bg-base-200 rounded-lg">
                {navItems.slice(0, 4).map(({ id, label, Icon }) => (
                  <button
                    key={id}
                    className={id === 'messages' ? 'active' : ''}
                  >
                    <div className="relative">
                      <Icon className="h-5 w-5" />
                      {id === 'messages' && (
                        <span className="badge badge-sm badge-error absolute -top-1 -right-1">
                          3
                        </span>
                      )}
                    </div>
                    <span className="btm-nav-label">{label}</span>
                  </button>
                ))}
              </div>
            </div>

            {/* Style 3: Icons Only */}
            <div>
              <h3 className="font-semibold mb-2">Icons Only (Compact)</h3>
              <div className="btm-nav btm-nav-lg relative h-12 bg-base-200 rounded-lg">
                {navItems.map(({ id, Icon }) => (
                  <button key={id} className={id === 'search' ? 'active' : ''}>
                    <Icon className="h-6 w-6" />
                  </button>
                ))}
              </div>
            </div>
          </div>
        </section>
      </AnimatedSection>

      {/* Implementation Code */}
      <AnimatedSection animation="slideUp" delay={0.3}>
        <section className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">Implementation</h2>
            <div className="mockup-code">
              <pre data-prefix="1">
                <code>{`<div className="btm-nav">`}</code>
              </pre>
              <pre data-prefix="2">
                <code>{`  <button className="active">`}</code>
              </pre>
              <pre data-prefix="3">
                <code>{`    <HomeIcon className="h-5 w-5" />`}</code>
              </pre>
              <pre data-prefix="4">
                <code>{`    <span className="btm-nav-label">Home</span>`}</code>
              </pre>
              <pre data-prefix="5">
                <code>{`  </button>`}</code>
              </pre>
              <pre data-prefix="6">
                <code>{`  <button>`}</code>
              </pre>
              <pre data-prefix="7">
                <code>{`    <SearchIcon className="h-5 w-5" />`}</code>
              </pre>
              <pre data-prefix="8">
                <code>{`    <span className="btm-nav-label">Search</span>`}</code>
              </pre>
              <pre data-prefix="9">
                <code>{`  </button>`}</code>
              </pre>
              <pre data-prefix="10">
                <code>{`</div>`}</code>
              </pre>
            </div>
          </div>
        </section>
      </AnimatedSection>
    </div>
  );
}
