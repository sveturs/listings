'use client';

import { useState, useEffect } from 'react';
import { adminApi, Category } from '@/services/admin';

export default function TestCategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = async () => {
    try {
      setLoading(true);
      const data = await adminApi.categories.getAll();
      setCategories(data);
    } catch (error) {
      console.error('Failed to load categories:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div>Loading categories...</div>;
  }

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-4">Test Categories Page</h1>
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title">Categories from API/Mock</h2>
          <p>Total categories: {categories.length}</p>
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
        </div>
      </div>
    </div>
  );
}
