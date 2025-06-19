'use client';

import { useRouter } from 'next/navigation';

export default function TestDashboardPage() {
  const router = useRouter();

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-6">Тестовая страница для дашбордов витрин</h1>
      
      <div className="space-y-4">
        <div className="card bg-base-200">
          <div className="card-body">
            <h2 className="card-title">Витрина test@user.rs</h2>
            <p>Владелец: test@user.rs (ID: 15)</p>
            <div className="flex gap-2 mt-4">
              <button 
                className="btn btn-primary"
                onClick={() => router.push('/en/storefronts/test-store-737754917/dashboard')}
              >
                Открыть дашборд
              </button>
              <button 
                className="btn btn-secondary"
                onClick={() => router.push('/en/storefronts/test-store-737754917/analytics')}
              >
                Открыть аналитику
              </button>
              <button 
                className="btn btn-accent"
                onClick={() => router.push('/en/storefronts/test-store-737754917')}
              >
                Публичная страница
              </button>
            </div>
          </div>
        </div>

        <div className="card bg-base-200">
          <div className="card-body">
            <h2 className="card-title">Витрина testuser-electronics-store</h2>
            <p>Владелец: test@user.rs (ID: 15)</p>
            <div className="flex gap-2 mt-4">
              <button 
                className="btn btn-primary"
                onClick={() => router.push('/en/storefronts/testuser-electronics-store/dashboard')}
              >
                Открыть дашборд
              </button>
              <button 
                className="btn btn-secondary"
                onClick={() => router.push('/en/storefronts/testuser-electronics-store/analytics')}
              >
                Открыть аналитику
              </button>
              <button 
                className="btn btn-accent"
                onClick={() => router.push('/en/storefronts/testuser-electronics-store')}
              >
                Публичная страница
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}