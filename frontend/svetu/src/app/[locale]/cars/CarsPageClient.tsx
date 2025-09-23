"use client";

import { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Image from "next/image";
import { apiClient } from "@/utils/api";
import type { components } from "@/types/generated/api";
import { Car, Search, TrendingUp, Calendar } from "lucide-react";

type CarMake = components["schemas"]["models.CarMake"];
type CarModel = components["schemas"]["models.CarModel"];
type MarketplaceListing = components["schemas"]["handler.MarketplaceListingWithLocation"];

interface CarsPageClientProps {
  locale: string;
}

export default function CarsPageClient({ locale }: CarsPageClientProps) {
  const t = useTranslations("cars");
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [popularMakes, setPopularMakes] = useState<CarMake[]>([]);
  const [latestListings, setLatestListings] = useState<MarketplaceListing[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [stats, setStats] = useState({
    totalListings: 0,
    totalMakes: 0,
    totalModels: 0,
  });

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);

      // Загружаем популярные марки
      const makesResponse = await apiClient.get("/api/v1/cars/makes", {
        params: { limit: 12 },
      });
      if (makesResponse.data?.data) {
        setPopularMakes(makesResponse.data.data.slice(0, 12));
      }

      // Загружаем последние автомобильные объявления
      const searchParams = {
        limit: 8,
        offset: 0,
        category_ids: "10101,10102,10103,10104", // Автомобильные категории
        sort: "created_at_desc",
      };

      const listingsResponse = await apiClient.post(
        "/api/v1/marketplace/search",
        searchParams
      );

      if (listingsResponse.data?.data?.items) {
        setLatestListings(listingsResponse.data.data.items);
        setStats({
          totalListings: listingsResponse.data.data.total || 0,
          totalMakes: makesResponse.data.data?.length || 0,
          totalModels: 3788, // Из БД
        });
      }
    } catch (error) {
      console.error("Error loading cars data:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = () => {
    if (searchQuery) {
      router.push(`/${locale}?q=${encodeURIComponent(searchQuery)}&category=10101`);
    }
  };

  const handleMakeClick = (makeSlug: string) => {
    router.push(`/${locale}?category=10101&car_make=${makeSlug}`);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-100">
      {/* Hero Section */}
      <div className="hero min-h-[400px] bg-gradient-to-br from-primary to-primary-focus text-primary-content">
        <div className="hero-content text-center">
          <div className="max-w-md">
            <h1 className="text-5xl font-bold mb-5 flex items-center justify-center gap-3">
              <Car className="w-12 h-12" />
              {t("heroTitle")}
            </h1>
            <p className="mb-8">{t("heroDescription")}</p>

            {/* Search Bar */}
            <div className="join w-full">
              <input
                type="text"
                placeholder={t("searchPlaceholder")}
                className="input input-bordered join-item flex-1 text-base-content"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              />
              <button
                className="btn btn-secondary join-item"
                onClick={handleSearch}
              >
                <Search className="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Statistics */}
      <div className="bg-base-200 py-8">
        <div className="container mx-auto px-4">
          <div className="stats stats-horizontal shadow w-full">
            <div className="stat">
              <div className="stat-figure text-primary">
                <Car className="w-8 h-8" />
              </div>
              <div className="stat-title">{t("stats.listings")}</div>
              <div className="stat-value text-primary">
                {stats.totalListings.toLocaleString()}
              </div>
            </div>
            <div className="stat">
              <div className="stat-figure text-secondary">
                <TrendingUp className="w-8 h-8" />
              </div>
              <div className="stat-title">{t("stats.makes")}</div>
              <div className="stat-value text-secondary">
                {stats.totalMakes}
              </div>
            </div>
            <div className="stat">
              <div className="stat-figure text-accent">
                <Calendar className="w-8 h-8" />
              </div>
              <div className="stat-title">{t("stats.models")}</div>
              <div className="stat-value text-accent">
                {stats.totalModels.toLocaleString()}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Popular Makes */}
      <div className="container mx-auto px-4 py-12">
        <h2 className="text-3xl font-bold mb-8">{t("popularMakes")}</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
          {popularMakes.map((make) => (
            <div
              key={make.id}
              className="card bg-base-100 shadow-xl cursor-pointer hover:shadow-2xl transition-shadow"
              onClick={() => handleMakeClick(make.slug)}
            >
              <div className="card-body items-center text-center p-4">
                {make.logo_url ? (
                  <div className="avatar">
                    <div className="w-16 rounded-full">
                      <Image
                        src={make.logo_url}
                        alt={make.name}
                        width={64}
                        height={64}
                      />
                    </div>
                  </div>
                ) : (
                  <div className="avatar placeholder">
                    <div className="bg-neutral text-neutral-content rounded-full w-16">
                      <span className="text-2xl">
                        {make.name.charAt(0).toUpperCase()}
                      </span>
                    </div>
                  </div>
                )}
                <h3 className="card-title text-sm">{make.name}</h3>
              </div>
            </div>
          ))}
        </div>

        <div className="text-center mt-8">
          <Link href={`/${locale}?category=10101`} className="btn btn-primary">
            {t("viewAllMakes")}
          </Link>
        </div>
      </div>

      {/* Categories */}
      <div className="bg-base-200 py-12">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold mb-8">{t("carCategories")}</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <Link
              href={`/${locale}?category=10101`}
              className="btn btn-outline btn-lg"
            >
              {t("categories.passenger")}
            </Link>
            <Link
              href={`/${locale}?category=10102`}
              className="btn btn-outline btn-lg"
            >
              {t("categories.suv")}
            </Link>
            <Link
              href={`/${locale}?category=10103`}
              className="btn btn-outline btn-lg"
            >
              {t("categories.commercial")}
            </Link>
            <Link
              href={`/${locale}?category=10104`}
              className="btn btn-outline btn-lg"
            >
              {t("categories.motorcycle")}
            </Link>
          </div>
        </div>
      </div>

      {/* Latest Listings */}
      {latestListings.length > 0 && (
        <div className="container mx-auto px-4 py-12">
          <h2 className="text-3xl font-bold mb-8">{t("latestListings")}</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {latestListings.map((listing) => (
              <Link
                key={listing.id}
                href={`/${locale}/listing/${listing.id}`}
                className="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow"
              >
                <figure className="aspect-[4/3] relative">
                  {listing.images && listing.images[0] ? (
                    <Image
                      src={listing.images[0].thumbnail_url || listing.images[0].url}
                      alt={listing.title}
                      fill
                      className="object-cover"
                    />
                  ) : (
                    <div className="w-full h-full bg-base-200 flex items-center justify-center">
                      <Car className="w-12 h-12 text-base-content/30" />
                    </div>
                  )}
                </figure>
                <div className="card-body p-4">
                  <h3 className="card-title text-base line-clamp-1">
                    {listing.title}
                  </h3>
                  {listing.price && (
                    <p className="text-lg font-bold text-primary">
                      €{listing.price.toLocaleString()}
                    </p>
                  )}
                  <p className="text-sm text-base-content/60">
                    {listing.location?.city || listing.location?.country}
                  </p>
                </div>
              </Link>
            ))}
          </div>

          <div className="text-center mt-8">
            <Link
              href={`/${locale}?category=10101,10102,10103,10104`}
              className="btn btn-primary btn-lg"
            >
              {t("viewAllListings")}
            </Link>
          </div>
        </div>
      )}

      {/* Quick Links */}
      <div className="bg-base-200 py-12">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold mb-8">{t("quickLinks")}</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title">{t("sellYourCar")}</h3>
                <p>{t("sellYourCarDescription")}</p>
                <div className="card-actions justify-end">
                  <Link
                    href={`/${locale}/create-listing-choice`}
                    className="btn btn-primary"
                  >
                    {t("createListing")}
                  </Link>
                </div>
              </div>
            </div>

            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title">{t("priceAnalysis")}</h3>
                <p>{t("priceAnalysisDescription")}</p>
                <div className="card-actions justify-end">
                  <button className="btn btn-secondary" disabled>
                    {t("comingSoon")}
                  </button>
                </div>
              </div>
            </div>

            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title">{t("vinDecoder")}</h3>
                <p>{t("vinDecoderDescription")}</p>
                <div className="card-actions justify-end">
                  <button className="btn btn-accent" disabled>
                    {t("comingSoon")}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}