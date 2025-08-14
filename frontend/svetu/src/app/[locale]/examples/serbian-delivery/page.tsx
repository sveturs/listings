'use client';

import { useState } from 'react';
import {
  TruckIcon,
  MapPinIcon,
  CreditCardIcon,
  ShieldCheckIcon,
  ClockIcon,
  CubeIcon,
  ChartBarIcon,
  UserGroupIcon,
  BuildingStorefrontIcon,
  CheckCircleIcon,
  ArrowRightIcon,
  BanknotesIcon,
  DocumentCheckIcon,
  ArrowPathIcon,
  GlobeAltIcon,
  DocumentTextIcon,
  ServerIcon,
  BoltIcon,
} from '@heroicons/react/24/outline';
import { StarIcon } from '@heroicons/react/24/solid';

// Import delivery components
import SerbianDeliveryMethodSelector from './components/SerbianDeliveryMethodSelector';
import SerbianTrackingWidget from './components/SerbianTrackingWidget';
import SerbianSellerShipmentInterface from './components/SerbianSellerShipmentInterface';
import SerbianParcelShopMap from './components/SerbianParcelShopMap';
import SerbianDeliveryCalculator from './components/SerbianDeliveryCalculator';
import SerbianBulkShipmentManager from './components/SerbianBulkShipmentManager';

// New Bex API components
import BexApiIntegration from './components/BexApiIntegration';
import BexShipmentCreator from './components/BexShipmentCreator';
import BexCustomsManager from './components/BexCustomsManager';
import BexLabelGenerator from './components/BexLabelGenerator';

export default function SerbianDeliveryExamplesPage() {
  const [activeTab, setActiveTab] = useState('overview');
  const [_selectedDeliveryMethod, _setSelectedDeliveryMethod] = useState('bex');

  const features = [
    {
      icon: ServerIcon,
      title: 'BexExpress RESTful API',
      description: 'Пуна интеграција са API за аутоматизацију достава',
      badge: 'NEW API',
      color: 'bg-purple-100',
      iconColor: 'text-purple-600',
    },
    {
      icon: GlobeAltIcon,
      title: 'Међународна достава',
      description: 'Царинска документа и међународне пошиљке преко BEX',
      badge: 'International',
      color: 'bg-green-100',
      iconColor: 'text-green-600',
    },
    {
      icon: MapPinIcon,
      title: 'Parcel Shop мрежа',
      description: '500+ локација за преузимање широм Србије',
      badge: 'Доступно',
      color: 'bg-blue-100',
      iconColor: 'text-blue-600',
    },
    {
      icon: BoltIcon,
      title: 'Real-time праћење',
      description: 'Праћење статуса пошиљке у реалном времену',
      badge: 'Live',
      color: 'bg-orange-100',
      iconColor: 'text-orange-600',
    },
    {
      icon: DocumentTextIcon,
      title: 'Аутоматске адреснице',
      description: 'Генерисање и штампање адресница преко API',
      badge: 'Automation',
      color: 'bg-indigo-100',
      iconColor: 'text-indigo-600',
    },
    {
      icon: ShieldCheckIcon,
      title: 'Осигурање до 100,000 РСД',
      description: 'Заштита вредне робе са потпуним осигурањем',
      badge: 'Заштићено',
      color: 'bg-red-100',
      iconColor: 'text-red-600',
    },
  ];

  const stats = [
    { label: 'Градова покривености', value: '180+', icon: MapPinIcon },
    {
      label: 'Пунктова за преузимање',
      value: '500+',
      icon: BuildingStorefrontIcon,
    },
    { label: 'Просечно време доставе', value: '1-2 дана', icon: ClockIcon },
    { label: 'Достављених пакета', value: '5M+', icon: CubeIcon },
  ];

  const testimonials = [
    {
      name: 'Милош Јовановић',
      role: 'Продавац електронике',
      rating: 5,
      text: 'AKS је најбржи за Београд. Пакети стигну за дан, купци су задовољни!',
      avatar: '👨‍💼',
    },
    {
      name: 'Јелена Петровић',
      role: 'Купац',
      rating: 5,
      text: 'Post Express пункт код куће је прави погодак. Преузимам када стигнем с посла.',
      avatar: '👩',
    },
    {
      name: 'Марко Николић',
      role: 'Власник радње',
      rating: 5,
      text: 'Ситићарго ради одлично за скупље ствари. Никад проблема са осигурањем.',
      avatar: '👨‍💻',
    },
  ];

  const tabs = [
    { id: 'overview', label: 'Преглед', icon: ChartBarIcon },
    { id: 'bex-api', label: 'BEX API', icon: ServerIcon, badge: 'NEW' },
    { id: 'shipment', label: 'Пошиљке', icon: CubeIcon },
    { id: 'customs', label: 'Царина', icon: GlobeAltIcon },
    { id: 'tracking', label: 'Праћење', icon: MapPinIcon },
    { id: 'labels', label: 'Адреснице', icon: DocumentTextIcon },
    { id: 'parcel-shops', label: 'Parcel Shops', icon: BuildingStorefrontIcon },
    { id: 'calculator', label: 'Калкулатор', icon: CreditCardIcon },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-b from-base-100 to-base-200">
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-blue-600 to-red-600 text-white">
        <div className="container mx-auto px-4 py-6 md:py-12">
          <div className="flex flex-col sm:flex-row items-center gap-3 mb-4">
            <div className="p-3 bg-white/20 rounded-xl backdrop-blur-sm">
              <TruckIcon className="w-6 h-6 sm:w-8 sm:h-8" />
            </div>
            <div className="text-center sm:text-left">
              <h1 className="text-2xl sm:text-3xl md:text-4xl font-bold">
                Интеграција српских курирских служби
              </h1>
              <p className="text-sm sm:text-base text-white/80 mt-2">
                AKS, Post Express, City Express и Ситићарго за ваш маркетплејс
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Tabs Navigation */}
      <div className="bg-base-100 border-b sticky top-0 z-40 backdrop-blur-lg bg-opacity-90">
        <div className="container mx-auto px-2 sm:px-4">
          <div className="flex gap-1 sm:gap-2 overflow-x-auto py-2 sm:py-4 scrollbar-hide">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`
                  flex items-center gap-1 sm:gap-2 px-2 sm:px-4 py-2 rounded-lg transition-all
                  whitespace-nowrap text-xs sm:text-sm font-medium min-w-fit relative
                  ${
                    activeTab === tab.id
                      ? 'bg-blue-600 text-white shadow-lg'
                      : 'bg-base-200 hover:bg-base-300'
                  }
                `}
              >
                <tab.icon className="w-4 h-4 sm:w-5 sm:h-5 flex-shrink-0" />
                <span className="hidden xs:inline">{tab.label}</span>
                {tab.badge && (
                  <span className="badge badge-xs badge-warning absolute -top-1 -right-1">
                    {tab.badge}
                  </span>
                )}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="container mx-auto px-4 py-4 sm:py-8">
        {activeTab === 'overview' && (
          <div className="space-y-8">
            {/* Features Grid */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
              {features.map((feature, index) => (
                <div
                  key={index}
                  className="card bg-base-100 shadow-xl hover:shadow-2xl transition-all hover:-translate-y-1"
                >
                  <div className="card-body">
                    <div className="flex items-start justify-between">
                      <div
                        className={`p-3 ${feature.color || 'bg-blue-100'} rounded-lg`}
                      >
                        <feature.icon
                          className={`w-6 h-6 ${feature.iconColor || 'text-blue-600'}`}
                        />
                      </div>
                      {feature.badge && (
                        <div className="badge badge-primary badge-sm">
                          {feature.badge}
                        </div>
                      )}
                    </div>
                    <h3 className="card-title text-lg mt-4">{feature.title}</h3>
                    <p className="text-base-content/70">
                      {feature.description}
                    </p>
                  </div>
                </div>
              ))}
            </div>

            {/* Stats */}
            <div className="bg-base-100 rounded-2xl shadow-xl p-4 sm:p-8">
              <h2 className="text-xl sm:text-2xl font-bold mb-4 sm:mb-6 text-center">
                Кључни показатељи српског тржишта
              </h2>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-6">
                {stats.map((stat, index) => (
                  <div key={index} className="text-center">
                    <div className="flex justify-center mb-3">
                      <div className="p-3 bg-red-100 rounded-full">
                        <stat.icon className="w-8 h-8 text-red-600" />
                      </div>
                    </div>
                    <div className="text-xl sm:text-3xl font-bold text-blue-600">
                      {stat.value}
                    </div>
                    <div className="text-xs sm:text-sm text-base-content/60 mt-1">
                      {stat.label}
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Process Steps */}
            <div className="bg-gradient-to-r from-blue-50 to-red-50 rounded-2xl p-4 sm:p-8">
              <h2 className="text-xl sm:text-2xl font-bold mb-4 sm:mb-6">
                Како функционише српска достава
              </h2>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-6">
                {[
                  {
                    step: '1',
                    title: 'Наруџбина',
                    desc: 'Купац бира начин доставе и плаћања',
                  },
                  {
                    step: '2',
                    title: 'Припрема',
                    desc: 'Продавац пакује и предаје куриру',
                  },
                  {
                    step: '3',
                    title: 'Достава',
                    desc: 'AKS/Post Express доставља робу',
                  },
                  {
                    step: '4',
                    title: 'Преузимање',
                    desc: 'Купац преузима и плаћа пошљом',
                  },
                ].map((item, index) => (
                  <div key={index} className="relative">
                    {index < 3 && (
                      <ArrowRightIcon className="absolute top-8 -right-3 w-6 h-6 text-blue-300 hidden md:block" />
                    )}
                    <div className="text-center">
                      <div className="w-12 h-12 sm:w-16 sm:h-16 bg-blue-600 text-white rounded-full flex items-center justify-center text-lg sm:text-2xl font-bold mx-auto mb-2 sm:mb-4">
                        {item.step}
                      </div>
                      <h3 className="text-sm sm:text-base font-semibold mb-1 sm:mb-2">
                        {item.title}
                      </h3>
                      <p className="text-xs sm:text-sm text-base-content/60">
                        {item.desc}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Testimonials */}
            <div>
              <h2 className="text-xl sm:text-2xl font-bold mb-4 sm:mb-6">
                Мишљења корисника
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 sm:gap-6">
                {testimonials.map((testimonial, index) => (
                  <div key={index} className="card bg-base-100 shadow-xl">
                    <div className="card-body">
                      <div className="flex items-center gap-3 mb-4">
                        <div className="text-4xl">{testimonial.avatar}</div>
                        <div>
                          <div className="font-semibold">
                            {testimonial.name}
                          </div>
                          <div className="text-sm text-base-content/60">
                            {testimonial.role}
                          </div>
                        </div>
                      </div>
                      <div className="flex gap-1 mb-3">
                        {[...Array(testimonial.rating)].map((_, i) => (
                          <StarIcon key={i} className="w-5 h-5 text-warning" />
                        ))}
                      </div>
                      <p className="text-base-content/80 italic">
                        &ldquo;{testimonial.text}&rdquo;
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {activeTab === 'bex-api' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>🔗 BexExpress RESTful API интеграција</h2>
              <p>
                Комплетна интеграција са BexExpress API за аутоматизацију
                достава, праћење пошиљки и управљање међународним пошиљкама.
              </p>
            </div>
            <BexApiIntegration />
          </div>
        )}

        {activeTab === 'shipment' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>📦 Креирање и управљање пошиљкама</h2>
              <p>
                Користите postShipments API за креирање домаћих пошиљки са свим
                опцијама: осигурање, откупнина, повратни документи.
              </p>
            </div>
            <BexShipmentCreator />
          </div>
        )}

        {activeTab === 'customs' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>🌍 Међународна достава и царина</h2>
              <p>
                postShipmentsCustoms API за ИНО пошиљке са царинском
                документацијом, HS кодовима и DDP опцијама.
              </p>
            </div>
            <BexCustomsManager />
          </div>
        )}

        {activeTab === 'labels' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>🏷️ Генерисање адресница и налепница</h2>
              <p>
                Аутоматско генерисање адресница у A4/A6 формату са баркодовима и
                позиционирањем за масовну штампу.
              </p>
            </div>
            <BexLabelGenerator />
          </div>
        )}

        {activeTab === 'parcel-shops' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>🏪 Мрежа Parcel Shop локација</h2>
              <p>
                Преглед свих пунктова за преузимање широм Србије са радним
                временом и GPS координатама.
              </p>
            </div>
            <SerbianParcelShopMap />
          </div>
        )}

        {activeTab === 'tracking' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>Праћење пошиљке</h2>
              <p>
                Купци и продавци могу да прате статус доставе у реалном времену
                кроз српске курирске службе.
              </p>
            </div>
            <SerbianTrackingWidget />
          </div>
        )}

        {activeTab === 'calculator' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>Калкулатор трошкова доставе</h2>
              <p>
                Израчунајте цену доставе у зависности од параметара пошиљке и
                руте кроз српске курирске службе.
              </p>
            </div>
            <SerbianDeliveryCalculator />
          </div>
        )}
      </div>

      {/* CTA Section */}
      <div className="bg-gradient-to-r from-blue-600 to-red-600 text-white mt-8 sm:mt-16">
        <div className="container mx-auto px-4 py-6 sm:py-12">
          <div className="text-center">
            <h2 className="text-xl sm:text-3xl font-bold mb-2 sm:mb-4">
              Спремни сте за српски тржиште?
            </h2>
            <p className="text-sm sm:text-xl mb-4 sm:mb-8 opacity-90">
              Интегришите српске курирске службе и проширите продају широм
              Србије
            </p>
            <div className="flex flex-col sm:flex-row gap-2 sm:gap-4 justify-center">
              <button className="btn btn-sm sm:btn-lg bg-white text-blue-600 hover:bg-white/90">
                <DocumentCheckIcon className="w-4 h-4 sm:w-5 sm:h-5" />
                API документација
              </button>
              <button className="btn btn-sm sm:btn-lg btn-outline border-white text-white hover:bg-white/20">
                <ArrowPathIcon className="w-4 h-4 sm:w-5 sm:h-5" />
                Почни интеграцију
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
