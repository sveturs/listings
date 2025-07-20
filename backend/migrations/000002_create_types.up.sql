-- Custom types migration

CREATE TYPE public.geo_source_type AS ENUM (
    'marketplace_listing',
    'storefront_product',
    'storefront'
);

CREATE TYPE public.location_privacy_level AS ENUM (
    'exact',
    'approximate',
    'city_only',
    'hidden'
);

CREATE TYPE public.storefront_geo_strategy AS ENUM (
    'storefront_location',
    'individual_location'
);