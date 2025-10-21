import { Metadata } from 'next';
import TestDataCleanupClient from './TestDataCleanupClient';

export const metadata: Metadata = {
  title: 'Test Data Cleanup | Admin',
  description: 'Manage and cleanup test data from database',
};

export default function TestDataCleanupPage() {
  return <TestDataCleanupClient />;
}
