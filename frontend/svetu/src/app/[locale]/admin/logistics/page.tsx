import { Metadata } from 'next';
import LogisticsPageClient from './LogisticsPageClient';

export const metadata: Metadata = {
  title: 'Logistics Monitoring | Admin',
  description: 'Monitor and manage logistics operations',
};

export default function LogisticsPage() {
  return <LogisticsPageClient />;
}
