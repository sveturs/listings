'use client';

import { useTranslations } from 'next-intl';
import Link from 'next/link';
import {
  MagnifyingGlassIcon,
  MapIcon,
  CheckCircleIcon,
  TruckIcon,
  CurrencyDollarIcon,
  BuildingStorefrontIcon,
  ExclamationTriangleIcon,
  XCircleIcon,
} from '@heroicons/react/24/outline';

type TestStatus = 'working' | 'issues' | 'not_working';

interface TestCardProps {
  title: string;
  description: string;
  href: string;
  icon: React.ReactNode;
  status: TestStatus;
  responseTime?: number;
  issueDescription?: string;
}

function TestCard({
  title,
  description,
  href,
  icon,
  status,
  responseTime,
  issueDescription,
}: TestCardProps) {
  const t = useTranslations('postexpressTest.index.testCard');

  const getBadgeClass = () => {
    switch (status) {
      case 'working':
        return 'badge-success';
      case 'issues':
        return 'badge-warning';
      case 'not_working':
        return 'badge-error';
    }
  };

  const getBadgeIcon = () => {
    switch (status) {
      case 'working':
        return <CheckCircleIcon className="w-4 h-4" />;
      case 'issues':
        return <ExclamationTriangleIcon className="w-4 h-4" />;
      case 'not_working':
        return <XCircleIcon className="w-4 h-4" />;
    }
  };

  const getBadgeText = () => {
    switch (status) {
      case 'working':
        return t('working');
      case 'issues':
        return t('issues');
      case 'not_working':
        return t('notWorking');
    }
  };

  const getBorderClass = () => {
    switch (status) {
      case 'working':
        return 'border-success hover:border-success';
      case 'issues':
        return 'border-warning hover:border-warning';
      case 'not_working':
        return 'border-error hover:border-error';
    }
  };

  return (
    <Link href={href}>
      <div
        className={`card bg-base-100 shadow-xl hover:shadow-2xl transition-all cursor-pointer h-full border-2 ${getBorderClass()}`}
      >
        <div className="card-body">
          <div className="flex items-start justify-between">
            <div className="text-primary">{icon}</div>
            <div className={`badge ${getBadgeClass()} gap-2`}>
              {getBadgeIcon()}
              {getBadgeText()}
            </div>
          </div>

          <h2 className="card-title mt-4">{title}</h2>
          <p className="text-base-content/70 text-sm">{description}</p>

          {responseTime && (
            <div className="text-xs text-base-content/50 mt-2">
              {t('responseTime', { time: responseTime })}
            </div>
          )}

          {issueDescription && (
            <div className="alert alert-warning py-2 mt-2">
              <div className="text-xs">{issueDescription}</div>
            </div>
          )}

          <div className="card-actions justify-end mt-4">
            <button className="btn btn-primary btn-sm">{t('button')} →</button>
          </div>
        </div>
      </div>
    </Link>
  );
}

export default function PostExpressAPIIndexPage() {
  const t = useTranslations('postexpressTest.index');

  // Полностью рабочие тесты
  const workingTests: TestCardProps[] = [
    {
      title: t('tests.tx3.title'),
      description: t('tests.tx3.description'),
      href: '/examples/postexpress-api/tx3-settlements',
      icon: <MagnifyingGlassIcon className="w-8 h-8" />,
      status: 'working',
      responseTime: 147,
    },
    {
      title: t('tests.tx4.title'),
      description: t('tests.tx4.description'),
      href: '/examples/postexpress-api/tx4-streets',
      icon: <MapIcon className="w-8 h-8" />,
      status: 'working',
      responseTime: 196,
    },
    {
      title: t('tests.tx6.title'),
      description: t('tests.tx6.description'),
      href: '/examples/postexpress-api/tx6-validate',
      icon: <CheckCircleIcon className="w-8 h-8" />,
      status: 'working',
      responseTime: 157,
    },
    {
      title: t('tests.tx73Standard.title'),
      description: t('tests.tx73Standard.description'),
      href: '/examples/postexpress-api/tx73-standard',
      icon: <TruckIcon className="w-8 h-8" />,
      status: 'working',
      responseTime: 71,
    },
  ];

  // Тесты с недостатками
  const testsWithIssues: TestCardProps[] = [
    {
      title: t('tests.tx11.title'),
      description: t('tests.tx11.description'),
      href: '/examples/postexpress-api/tx11-postage',
      icon: <CurrencyDollarIcon className="w-8 h-8" />,
      status: 'issues',
      responseTime: 81,
      issueDescription: t('issueDescriptions.tx11'),
    },
  ];

  // Нерабочие тесты
  const notWorkingTests: TestCardProps[] = [
    {
      title: t('tests.tx9.title'),
      description: t('tests.tx9.description'),
      href: '/examples/postexpress-api/tx9-availability',
      icon: <TruckIcon className="w-8 h-8" />,
      status: 'not_working',
      issueDescription: t('issueDescriptions.tx9'),
    },
    {
      title: t('tests.tx73Cod.title'),
      description: t('tests.tx73Cod.description'),
      href: '/examples/postexpress-api/tx73-cod',
      icon: <CurrencyDollarIcon className="w-8 h-8" />,
      status: 'not_working',
      issueDescription: t('issueDescriptions.tx73Cod'),
    },
    {
      title: t('tests.tx73ParcelLocker.title'),
      description: t('tests.tx73ParcelLocker.description'),
      href: '/examples/postexpress-api/tx73-parcel-locker',
      icon: <BuildingStorefrontIcon className="w-8 h-8" />,
      status: 'not_working',
      issueDescription: t('issueDescriptions.tx73ParcelLocker'),
    },
  ];

  const totalTests =
    workingTests.length + testsWithIssues.length + notWorkingTests.length;

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="bg-gradient-to-r from-blue-600 to-indigo-600 text-white py-12">
        <div className="container mx-auto px-4">
          <h1 className="text-5xl font-bold mb-4">
            {t('title')}
          </h1>
          <p className="text-xl opacity-90 max-w-3xl">
            {t('subtitle')}
          </p>
          <div className="mt-6 flex gap-2 flex-wrap">
            <div className="badge badge-success badge-lg gap-2">
              <CheckCircleIcon className="w-4 h-4" />
              API Ready
            </div>
            <div className="badge badge-warning badge-lg gap-2">
              <CheckCircleIcon className="w-4 h-4" />
              Test Mode
            </div>
          </div>
        </div>
      </div>

      {/* Stats */}
      <div className="container mx-auto px-4 py-8">
        <div className="stats shadow w-full">
          <div className="stat">
            <div className="stat-figure text-primary">
              <CheckCircleIcon className="w-8 h-8" />
            </div>
            <div className="stat-title">{t('stats.totalTests')}</div>
            <div className="stat-value text-primary">{totalTests}</div>
            <div className="stat-desc">{t('stats.totalTestsDesc')}</div>
          </div>

          <div className="stat">
            <div className="stat-figure text-success">
              <CheckCircleIcon className="w-8 h-8" />
            </div>
            <div className="stat-title">{t('stats.fullyWorking')}</div>
            <div className="stat-value text-success">{workingTests.length}</div>
            <div className="stat-desc">{t('stats.fullyWorkingDesc')}</div>
          </div>

          <div className="stat">
            <div className="stat-figure text-warning">
              <ExclamationTriangleIcon className="w-8 h-8" />
            </div>
            <div className="stat-title">{t('stats.withIssues')}</div>
            <div className="stat-value text-warning">
              {testsWithIssues.length}
            </div>
            <div className="stat-desc">{t('stats.withIssuesDesc')}</div>
          </div>

          <div className="stat">
            <div className="stat-figure text-error">
              <XCircleIcon className="w-8 h-8" />
            </div>
            <div className="stat-title">{t('stats.notWorking')}</div>
            <div className="stat-value text-error">
              {notWorkingTests.length}
            </div>
            <div className="stat-desc">{t('stats.notWorkingDesc')}</div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 pb-12">
        {/* Полностью рабочие тесты */}
        <div className="mb-12">
          <h2 className="text-3xl font-bold mb-6 flex items-center gap-3">
            <CheckCircleIcon className="w-8 h-8 text-success" />
            {t('sections.fullyWorking')}
            <span className="badge badge-success badge-lg">
              {workingTests.length}
            </span>
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {workingTests.map((test) => (
              <TestCard key={test.href} {...test} />
            ))}
          </div>
        </div>

        {/* Тесты с недостатками */}
        {testsWithIssues.length > 0 && (
          <div className="mb-12">
            <h2 className="text-3xl font-bold mb-6 flex items-center gap-3">
              <ExclamationTriangleIcon className="w-8 h-8 text-warning" />
              {t('sections.withIssues')}
              <span className="badge badge-warning badge-lg">
                {testsWithIssues.length}
              </span>
            </h2>
            <div className="alert alert-warning mb-6">
              <ExclamationTriangleIcon className="w-6 h-6" />
              <div>
                <h3 className="font-bold">{t('issuesAlert.title')}</h3>
                <div className="text-sm">
                  {t('issuesAlert.description')}
                </div>
              </div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {testsWithIssues.map((test) => (
                <TestCard key={test.href} {...test} />
              ))}
            </div>
          </div>
        )}

        {/* Нерабочие тесты */}
        {notWorkingTests.length > 0 && (
          <div className="mb-12">
            <h2 className="text-3xl font-bold mb-6 flex items-center gap-3">
              <XCircleIcon className="w-8 h-8 text-error" />
              {t('sections.notWorking')}
              <span className="badge badge-error badge-lg">
                {notWorkingTests.length}
              </span>
            </h2>
            <div className="alert alert-error mb-6">
              <XCircleIcon className="w-6 h-6" />
              <div>
                <h3 className="font-bold">{t('issuesAlert.title')}</h3>
                <div className="text-sm">
                  {t('issuesAlert.description')}
                </div>
              </div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {notWorkingTests.map((test) => (
                <TestCard key={test.href} {...test} />
              ))}
            </div>
          </div>
        )}

        {/* Info Section */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h3 className="card-title">{t('about.title')}</h3>
            <p className="text-base-content/70">
              {t('about.description')}
            </p>
            <div className="divider"></div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
              <div>
                <h4 className="font-semibold mb-2">{t('about.features.realtime.title')}</h4>
                <p className="text-base-content/60">
                  {t('about.features.realtime.description')}
                </p>
              </div>
              <div>
                <h4 className="font-semibold mb-2">{t('about.features.requestResponse.title')}</h4>
                <p className="text-base-content/60">
                  {t('about.features.requestResponse.description')}
                </p>
              </div>
              <div>
                <h4 className="font-semibold mb-2">{t('about.features.prefilledForms.title')}</h4>
                <p className="text-base-content/60">
                  {t('about.features.prefilledForms.description')}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
