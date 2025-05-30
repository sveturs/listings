import type { NextConfig } from 'next';
import createNextIntlPlugin from 'next-intl/plugin';
import configManager from './src/config';

const withNextIntl = createNextIntlPlugin();

const nextConfig: NextConfig = {
  images: {
    remotePatterns: configManager.getImageHosts(),
  },
};

export default withNextIntl(nextConfig);
