'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';

interface OrderFiltersProps {
  onApplyFilters: (filters: any) => void;
  onClearFilters: () => void;
}

export default function OrderFilters({
  onApplyFilters,
  onClearFilters,
}: OrderFiltersProps) {
  const t = useTranslations('orders');
  const [showFilters, setShowFilters] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo, setDateTo] = useState('');
  const [minPrice, setMinPrice] = useState('');
  const [maxPrice, setMaxPrice] = useState('');
  const [selectedStatuses, setSelectedStatuses] = useState<string[]>([]);

  const statusOptions = [
    { value: 'pending', label: t('status.pending') },
    { value: 'paid', label: t('status.paid') },
    { value: 'confirmed', label: t('status.confirmed') },
    { value: 'shipped', label: t('status.shipped') },
    { value: 'delivered', label: t('status.delivered') },
    { value: 'completed', label: t('status.completed') },
    { value: 'cancelled', label: t('status.cancelled') },
    { value: 'disputed', label: t('status.disputed') },
    { value: 'refunded', label: t('status.refunded') },
  ];

  const handleStatusToggle = (status: string) => {
    setSelectedStatuses((prev) =>
      prev.includes(status)
        ? prev.filter((s) => s !== status)
        : [...prev, status]
    );
  };

  const handleApplyFilters = () => {
    const filters = {
      search: searchQuery,
      dateFrom,
      dateTo,
      minPrice: minPrice ? parseFloat(minPrice) : undefined,
      maxPrice: maxPrice ? parseFloat(maxPrice) : undefined,
      statuses: selectedStatuses,
    };
    onApplyFilters(filters);
    setShowFilters(false);
  };

  const handleClearFilters = () => {
    setSearchQuery('');
    setDateFrom('');
    setDateTo('');
    setMinPrice('');
    setMaxPrice('');
    setSelectedStatuses([]);
    onClearFilters();
    setShowFilters(false);
  };

  const hasActiveFilters =
    searchQuery ||
    dateFrom ||
    dateTo ||
    minPrice ||
    maxPrice ||
    selectedStatuses.length > 0;

  return (
    <div className="mb-6">
      {/* Search Bar and Toggle */}
      <div className="flex flex-col sm:flex-row gap-4 mb-4">
        <div className="flex-1">
          <div className="form-control">
            <div className="input-group">
              <input
                type="text"
                placeholder={t('searchOrders')}
                className="input input-bordered w-full"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    handleApplyFilters();
                  }
                }}
              />
              <button className="btn btn-square" onClick={handleApplyFilters}>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-6 w-6"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                  />
                </svg>
              </button>
            </div>
          </div>
        </div>

        <button
          onClick={() => setShowFilters(!showFilters)}
          className={`btn ${showFilters ? 'btn-primary' : 'btn-outline'}`}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="h-5 w-5 mr-2"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"
            />
          </svg>
          {t('filter')}
          {hasActiveFilters && (
            <span className="badge badge-sm badge-secondary ml-2">
              {selectedStatuses.length +
                (dateFrom || dateTo ? 1 : 0) +
                (minPrice || maxPrice ? 1 : 0)}
            </span>
          )}
        </button>
      </div>

      {/* Expanded Filters */}
      {showFilters && (
        <div className="card bg-base-100 shadow-lg">
          <div className="card-body">
            <h3 className="card-title text-lg mb-4">{t('advancedFilters')}</h3>

            {/* Status Filters */}
            <div className="form-control mb-4">
              <label className="label">
                <span className="label-text font-semibold">
                  {t('status.label')}
                </span>
              </label>
              <div className="flex flex-wrap gap-2">
                {statusOptions.map((status) => (
                  <label key={status.value} className="cursor-pointer">
                    <input
                      type="checkbox"
                      className="checkbox checkbox-sm"
                      checked={selectedStatuses.includes(status.value)}
                      onChange={() => handleStatusToggle(status.value)}
                    />
                    <span className="label-text ml-2">{status.label}</span>
                  </label>
                ))}
              </div>
            </div>

            {/* Date Range */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-semibold">
                    {t('dateRange')}
                  </span>
                </label>
                <div className="flex gap-2">
                  <div className="flex-1">
                    <label className="label">
                      <span className="label-text text-xs">{t('from')}</span>
                    </label>
                    <input
                      type="date"
                      className="input input-bordered input-sm w-full"
                      value={dateFrom}
                      onChange={(e) => setDateFrom(e.target.value)}
                    />
                  </div>
                  <div className="flex-1">
                    <label className="label">
                      <span className="label-text text-xs">{t('to')}</span>
                    </label>
                    <input
                      type="date"
                      className="input input-bordered input-sm w-full"
                      value={dateTo}
                      onChange={(e) => setDateTo(e.target.value)}
                    />
                  </div>
                </div>
              </div>

              {/* Price Range */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-semibold">
                    {t('priceRange')}
                  </span>
                </label>
                <div className="flex gap-2">
                  <div className="flex-1">
                    <label className="label">
                      <span className="label-text text-xs">
                        {t('minPrice')}
                      </span>
                    </label>
                    <input
                      type="number"
                      className="input input-bordered input-sm w-full"
                      placeholder="0"
                      value={minPrice}
                      onChange={(e) => setMinPrice(e.target.value)}
                    />
                  </div>
                  <div className="flex-1">
                    <label className="label">
                      <span className="label-text text-xs">
                        {t('maxPrice')}
                      </span>
                    </label>
                    <input
                      type="number"
                      className="input input-bordered input-sm w-full"
                      placeholder="999999"
                      value={maxPrice}
                      onChange={(e) => setMaxPrice(e.target.value)}
                    />
                  </div>
                </div>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="card-actions justify-end mt-4">
              <button
                onClick={handleClearFilters}
                className="btn btn-ghost"
                disabled={!hasActiveFilters}
              >
                {t('clearFilters')}
              </button>
              <button onClick={handleApplyFilters} className="btn btn-primary">
                {t('applyFilters')}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Active Filters Display */}
      {hasActiveFilters && !showFilters && (
        <div className="alert alert-info">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            className="stroke-current shrink-0 w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            ></path>
          </svg>
          <div className="flex-1">
            <span>{t('activeFilters')}: </span>
            {searchQuery && (
              <span className="badge badge-sm mx-1">
                {t('search')}: {searchQuery}
              </span>
            )}
            {selectedStatuses.map((status) => (
              <span key={status} className="badge badge-sm mx-1">
                {statusOptions.find((s) => s.value === status)?.label}
              </span>
            ))}
            {(dateFrom || dateTo) && (
              <span className="badge badge-sm mx-1">
                {t('dateRange')}: {dateFrom || '...'} - {dateTo || '...'}
              </span>
            )}
            {(minPrice || maxPrice) && (
              <span className="badge badge-sm mx-1">
                {t('priceRange')}: {minPrice || '0'} - {maxPrice || '∞'} RSD
              </span>
            )}
          </div>
          <button onClick={handleClearFilters} className="btn btn-ghost btn-sm">
            {t('clearFilters')}
          </button>
        </div>
      )}
    </div>
  );
}
