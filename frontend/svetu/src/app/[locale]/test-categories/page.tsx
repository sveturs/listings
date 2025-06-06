'use client';

import { useState, useEffect } from 'react';
import { adminApi, Category } from '@/services/admin';

export default function TestCategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await adminApi.categories.getAll();
      setCategories(data);
    } catch (error) {
      console.error('Failed to load categories:', error);
      setError(
        error instanceof Error ? error.message : 'Failed to load categories'
      );
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-4">
        Test Categories Page (Outside Admin)
      </h1>

      {error && (
        <div className="alert alert-error mb-4">
          <span>Error: {error}</span>
        </div>
      )}

      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title">Categories from API/Mock</h2>
          <p>Total categories: {categories.length}</p>

          {categories.length === 0 ? (
            <p>No categories found</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Slug</th>
                    <th>Parent ID</th>
                  </tr>
                </thead>
                <tbody>
                  {categories.map((cat) => (
                    <tr key={cat.id}>
                      <td>{cat.id}</td>
                      <td>{cat.name}</td>
                      <td>{cat.slug}</td>
                      <td>{cat.parent_id || '-'}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          <div className="mt-4">
            <button className="btn btn-primary" onClick={loadCategories}>
              Reload Categories
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
