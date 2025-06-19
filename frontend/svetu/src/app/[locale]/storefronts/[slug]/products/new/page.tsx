'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { toast } from '@/utils/toast';
import { apiClient } from '@/services/api-client';
import type { components } from '@/types/generated/api';

type CreateProductRequest =
  components['schemas']['backend_internal_domain_models.CreateProductRequest'];
type MarketplaceCategory =
  components['schemas']['backend_internal_domain_models.MarketplaceCategory'];

interface PageProps {
  params: Promise<{
    locale: string;
    slug: string;
  }>;
}

export default function NewProductPage({ params }: PageProps) {
  const [slug, setSlug] = useState<string>('');
  
  useEffect(() => {
    params.then((p) => setSlug(p.slug));
  }, [params]);
  const t = useTranslations();
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [categories, setCategories] = useState<MarketplaceCategory[]>([]);
  const [loadingCategories, setLoadingCategories] = useState(true);

  const [formData, setFormData] = useState<CreateProductRequest>({
    name: '',
    description: '',
    price: 0,
    currency: 'RSD',
    category_id: 0,
    stock_quantity: 0,
    is_active: true,
    attributes: {},
  });

  // Load categories
  useEffect(() => {
    const loadCategories = async () => {
      try {
        const response = await apiClient.get('/api/v1/marketplace/categories/tree');
        if (response.data) {
          const flattenCategories = (cats: any[]): MarketplaceCategory[] => {
            let result: MarketplaceCategory[] = [];
            cats.forEach((cat) => {
              result.push(cat);
              if (cat.children && cat.children.length > 0) {
                result = result.concat(flattenCategories(cat.children));
              }
            });
            return result;
          };
          setCategories(flattenCategories(response.data));
        }
      } catch (error) {
        console.error('Failed to load categories:', error);
        toast.error(t('storefronts.products.errorLoadingCategories'));
      } finally {
        setLoadingCategories(false);
      }
    };
    loadCategories();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!slug) {
      return; // Wait for slug to be loaded
    }

    if (!formData.name || !formData.description || !formData.category_id) {
      toast.error(t('storefronts.products.fillRequiredFields'));
      return;
    }

    setLoading(true);
    try {
      const response = await apiClient.post(
        `/api/v1/storefronts/${slug}/products`,
        formData
      );
      if (response.data) {
        toast.success(t('storefronts.products.productCreated'));
        router.push(`/storefronts/${slug}/products`);
      }
    } catch (error: any) {
      console.error('Failed to create product:', error);
      toast.error(
        error.response?.data?.error ||
          t('storefronts.products.errorCreatingProduct')
      );
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value, type } = e.target;

    if (type === 'checkbox') {
      const checked = (e.target as HTMLInputElement).checked;
      setFormData((prev) => ({ ...prev, [name]: checked }));
    } else if (
      name === 'price' ||
      name === 'stock_quantity' ||
      name === 'category_id'
    ) {
      setFormData((prev) => ({ ...prev, [name]: Number(value) }));
    } else {
      setFormData((prev) => ({ ...prev, [name]: value }));
    }
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      {/* Header */}
      <div className="mb-8">
        <Link
          href={`/storefronts/${slug}/products`}
          className="inline-flex items-center text-primary hover:underline mb-4"
        >
          <ArrowLeftIcon className="w-4 h-4 mr-2" />
          {t('common.back')}
        </Link>
        <h1 className="text-3xl font-bold">
          {t('storefronts.products.addNewProduct')}
        </h1>
      </div>

      {/* Form */}
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Basic Information */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">
              {t('storefronts.products.basicInformation')}
            </h2>

            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('storefronts.products.productName')} *
                </span>
              </label>
              <input
                type="text"
                name="name"
                value={formData.name}
                onChange={handleChange}
                className="input input-bordered"
                required
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">
                  {t('storefronts.products.description')} *
                </span>
              </label>
              <textarea
                name="description"
                value={formData.description}
                onChange={handleChange}
                className="textarea textarea-bordered h-24"
                required
              />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text">
                    {t('storefronts.products.category')} *
                  </span>
                </label>
                <select
                  name="category_id"
                  value={formData.category_id}
                  onChange={handleChange}
                  className="select select-bordered"
                  required
                  disabled={loadingCategories}
                >
                  <option value={0}>
                    {t('storefronts.products.selectCategory')}
                  </option>
                  {categories.map((cat) => (
                    <option key={cat.id} value={cat.id}>
                      {cat.parent_id ? '— ' : ''}
                      {cat.name}
                    </option>
                  ))}
                </select>
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text">
                    {t('storefronts.products.sku')}
                  </span>
                </label>
                <input
                  type="text"
                  name="sku"
                  value={formData.sku || ''}
                  onChange={handleChange}
                  className="input input-bordered"
                  placeholder={t('storefronts.products.skuPlaceholder')}
                />
              </div>
            </div>
          </div>
        </div>

        {/* Pricing */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">
              {t('storefronts.products.pricing')}
            </h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text">
                    {t('storefronts.products.price')} *
                  </span>
                </label>
                <input
                  type="number"
                  name="price"
                  value={formData.price}
                  onChange={handleChange}
                  className="input input-bordered"
                  min="0"
                  step="0.01"
                  required
                />
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text">
                    {t('storefronts.products.currency')}
                  </span>
                </label>
                <select
                  name="currency"
                  value={formData.currency}
                  onChange={handleChange}
                  className="select select-bordered"
                >
                  <option value="RSD">RSD</option>
                  <option value="EUR">EUR</option>
                  <option value="USD">USD</option>
                </select>
              </div>
            </div>
          </div>
        </div>

        {/* Inventory */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">
              {t('storefronts.products.inventory')}
            </h2>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="form-control">
                <label className="label">
                  <span className="label-text">
                    {t('storefronts.products.stockQuantity')}
                  </span>
                </label>
                <input
                  type="number"
                  name="stock_quantity"
                  value={formData.stock_quantity}
                  onChange={handleChange}
                  className="input input-bordered"
                  min="0"
                />
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text">
                    {t('storefronts.products.barcode')}
                  </span>
                </label>
                <input
                  type="text"
                  name="barcode"
                  value={formData.barcode || ''}
                  onChange={handleChange}
                  className="input input-bordered"
                  placeholder={t('storefronts.products.barcodePlaceholder')}
                />
              </div>
            </div>

            <div className="form-control">
              <label className="label cursor-pointer">
                <span className="label-text">
                  {t('storefronts.products.activeProduct')}
                </span>
                <input
                  type="checkbox"
                  name="is_active"
                  checked={formData.is_active}
                  onChange={handleChange}
                  className="checkbox checkbox-primary"
                />
              </label>
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex justify-end gap-4">
          <Link
            href={`/storefronts/${slug}/products`}
            className="btn btn-ghost"
          >
            {t('common.cancel')}
          </Link>
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading && <span className="loading loading-spinner"></span>}
            {t('storefronts.products.createProduct')}
          </button>
        </div>
      </form>
    </div>
  );
}
