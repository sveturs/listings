'use client';

import { useState } from 'react';
import {
  DocumentArrowUpIcon,
  TrashIcon,
  ArrowPathIcon,
} from '@heroicons/react/24/outline';

export default function SerbianBulkShipmentManager() {
  const [orders] = useState([
    {
      id: '001',
      customer: 'Марија Стојановић',
      city: 'Београд',
      amount: 3500,
      status: 'pending',
    },
    {
      id: '002',
      customer: 'Петар Михаиловић',
      city: 'Нови Сад',
      amount: 2800,
      status: 'pending',
    },
    {
      id: '003',
      customer: 'Ана Јовановић',
      city: 'Ниш',
      amount: 4200,
      status: 'processing',
    },
  ]);

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h3 className="card-title mb-4">Масовна обрада пошиљки</h3>

        {/* Upload Section */}
        <div className="border-2 border-dashed border-base-300 rounded-lg p-6 text-center mb-4">
          <DocumentArrowUpIcon className="w-12 h-12 mx-auto text-base-content/50 mb-2" />
          <p className="text-base-content/70">
            Превуци CSV фајл или кликни да избереш
          </p>
          <p className="text-xs text-base-content/50 mt-1">
            Подржани формати: AKS, Post Express, City Express
          </p>
          <button className="btn btn-primary btn-sm mt-2">Избери фајл</button>
        </div>

        {/* Orders List */}
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <h4 className="font-semibold">Наруџбине (3)</h4>
            <div className="flex gap-2">
              <button className="btn btn-ghost btn-sm">Селектуј све</button>
              <button className="btn btn-primary btn-sm gap-2">
                <ArrowPathIcon className="w-4 h-4" />
                Обради
              </button>
            </div>
          </div>

          {orders.map((order) => (
            <div
              key={order.id}
              className="flex items-center gap-3 p-3 bg-base-200 rounded"
            >
              <input type="checkbox" className="checkbox checkbox-sm" />
              <div className="flex-1">
                <div className="flex items-center justify-between">
                  <span className="font-medium">{order.customer}</span>
                  <div className="flex items-center gap-2">
                    <span className="text-sm">{order.city}</span>
                    <span className="font-semibold">{order.amount} РСД</span>
                    <div
                      className={`badge badge-sm ${
                        order.status === 'pending'
                          ? 'badge-warning'
                          : 'badge-info'
                      }`}
                    >
                      {order.status === 'pending' ? 'На чекању' : 'У обради'}
                    </div>
                  </div>
                </div>
              </div>
              <button className="btn btn-ghost btn-sm text-error">
                <TrashIcon className="w-4 h-4" />
              </button>
            </div>
          ))}
        </div>

        {/* Serbian Courier Options */}
        <div className="mt-4">
          <h5 className="font-medium mb-2">Избор курирске службе:</h5>
          <div className="grid grid-cols-2 gap-2">
            <label className="cursor-pointer">
              <input
                type="radio"
                name="courier"
                className="radio radio-sm"
                defaultChecked
              />
              <span className="ml-2">AKS Express</span>
            </label>
            <label className="cursor-pointer">
              <input type="radio" name="courier" className="radio radio-sm" />
              <span className="ml-2">Post Express</span>
            </label>
            <label className="cursor-pointer">
              <input type="radio" name="courier" className="radio radio-sm" />
              <span className="ml-2">City Express</span>
            </label>
            <label className="cursor-pointer">
              <input type="radio" name="courier" className="radio radio-sm" />
              <span className="ml-2">Yettel Post</span>
            </label>
          </div>
        </div>
      </div>
    </div>
  );
}
