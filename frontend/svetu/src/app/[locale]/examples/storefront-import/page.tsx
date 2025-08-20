'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { AnimatedSection } from '@/components/ui/AnimatedSection';
import { SveTuLogoStatic } from '@/components/logos/SveTuLogoStatic';

const StorefrontImportDemo = () => {
  const _t = useTranslations('storefronts');
  const [selectedMethod, setSelectedMethod] = useState<string>('');
  const [currentStep, setCurrentStep] = useState(1);
  const [uploadedFile, setUploadedFile] = useState<File | null>(null);
  const [csvData, setCsvData] = useState<string[][]>([]);
  const [fieldMapping, setFieldMapping] = useState<Record<string, string>>({});
  const [previewData, setPreviewData] = useState<any[]>([]);
  const [xmlData, setXmlData] = useState<any>(null);
  const [excelSheets, setExcelSheets] = useState<string[]>([]);
  const [selectedSheet, setSelectedSheet] = useState<string>('');
  const [apiSettings, setApiSettings] = useState<Record<string, string>>({});
  const [ftpSettings, setFtpSettings] = useState<Record<string, string>>({});
  const [webhookSettings, setWebhookSettings] = useState<
    Record<string, string>
  >({});
  const [importProgress, setImportProgress] = useState(0);
  const [isImporting, setIsImporting] = useState(false);

  // Демо данные для сербского рынка
  const sampleCsvData = [
    [
      'Naziv',
      'Cena (RSD)',
      'Količina',
      'Kategorija',
      'Opis',
      'Brend',
      'Šifra',
      'PDV (%)',
      'Zemlja porekla',
    ],
    [
      'Samsung Galaxy S24 Ultra 256GB',
      '149990',
      '12',
      'Elektronika > Mobilni telefoni',
      'Najnoviji Samsung flagship sa S Pen olovkom',
      'Samsung',
      'SM-S928B',
      '20',
      'Južna Koreja',
    ],
    [
      'Laptop HP Pavilion 15"',
      '89990',
      '5',
      'Elektronika > Računari',
      'Intel Core i5, 8GB RAM, 512GB SSD',
      'HP',
      'HP-PAV15',
      '20',
      'Kina',
    ],
    [
      'Wireless slušalice JBL',
      '12990',
      '30',
      'Elektronika > Audio',
      'Bluetooth slušalice sa noise cancellation',
      'JBL',
      'JBL-WH900',
      '20',
      'Kina',
    ],
    [
      'Kafemat DeLonghi',
      '34990',
      '8',
      'Kućni aparati > Kuhinja',
      'Automatski espresso aparat sa mlečnom penom',
      'DeLonghi',
      'DL-EC235',
      '20',
      'Italija',
    ],
  ];

  const sampleXmlData = {
    shop: {
      name: 'TehnoSvet Srbija',
      company: 'DOO TehnoSvet Beograd',
      url: 'https://tehnosvet.rs',
      pib: '12345678',
      categories: [
        { id: '1', name: 'Elektronika' },
        { id: '2', name: 'Mobilni telefoni', parentId: '1' },
        { id: '3', name: 'Računari', parentId: '1' },
        { id: '4', name: 'Kućni aparati' },
      ],
      offers: [
        {
          id: '12345',
          name: 'Xiaomi Redmi Note 12 Pro',
          categoryId: '2',
          price: '32990',
          currency: 'RSD',
          vendor: 'Xiaomi',
          model: 'Redmi Note 12 Pro',
          description: 'Xiaomi telefon sa 108MP kamerom i 120W brzim punjenjem',
          picture: 'https://example.com/xiaomi-note12.jpg',
          available: 'true',
          pdv: '20',
          zemlja_porekla: 'Kina',
          params: [
            { name: 'Ekran', value: '6.67"' },
            { name: 'Memorija', value: '128GB' },
            { name: 'Boja', value: 'Plava' },
            { name: 'Garantija', value: '24 meseca' },
          ],
        },
        {
          id: '12346',
          name: 'Lenovo IdeaPad Gaming 3',
          categoryId: '3',
          price: '89990',
          currency: 'RSD',
          vendor: 'Lenovo',
          model: 'IdeaPad Gaming 3',
          description:
            'Gaming laptop sa AMD Ryzen 5 procesorom i GTX 1650 grafikom',
          picture: 'https://example.com/lenovo-gaming.jpg',
          available: 'true',
          pdv: '20',
          zemlja_porekla: 'Kina',
          params: [
            { name: 'Procesor', value: 'AMD Ryzen 5 5600H' },
            { name: 'RAM', value: '8GB DDR4' },
            { name: 'Grafika', value: 'GTX 1650 4GB' },
            { name: 'Garantija', value: '24 meseca' },
          ],
        },
      ],
    },
  };

  const sampleExcelSheets = ['Proizvodi', 'Zalihe', 'Cene', 'Karakteristike'];

  const sampleApiIntegrations = [
    {
      name: 'PANTHEON ERP',
      icon: '🏢',
      description:
        'Najpoznatiji regionalni ERP za Srbiju i JI Evropu (Datalab)',
      fields: ['server_url', 'database', 'username', 'password'],
      status: 'available',
      sampleData: {
        server_url: 'http://pantheon-server.mojafirma.rs:8080',
        database: 'TragovinaPreduzeće',
        username: 'admin',
        password: '••••••••',
      },
      features: [
        'PDV saglasnost',
        'E-pošta fiskalizacija',
        'Srpsko zakonodavstvo',
        'Multi-valute',
      ],
      lastSync: 'Pre 3 minuta',
      productsCount: 3247,
    },
    {
      name: 'KupujemProdajem',
      icon: '🛒',
      description: 'Uvoz oglasa sa najveće srpske trading platforme',
      fields: ['username', 'password', 'category_filter'],
      status: 'available',
      sampleData: {
        username: 'moj_nalog@email.rs',
        password: '••••••••',
        category_filter: 'Elektronika',
      },
      features: ['Web scraping', 'Auto kategorije', 'Slika preuzimanje'],
      lastSync: 'Pre 1 sat',
      productsCount: 1235,
    },
    {
      name: 'Limundo',
      icon: '🏪',
      description: 'Aukcijska platforma - import završenih aukcija',
      fields: ['api_key', 'seller_id', 'min_price'],
      status: 'available',
      sampleData: {
        api_key: 'lim_12345...',
        seller_id: 'prodavac123',
        min_price: '1000',
      },
      features: ['Auction API', 'Sold items', 'Price history'],
      lastSync: 'Pre 30 minuta',
      productsCount: 892,
    },
    {
      name: 'Microsoft Dynamics 365 BC',
      icon: '🛍️',
      description: 'Microsoft ERP sa srpskom lokalizacijom (bivši Navision)',
      fields: ['tenant_url', 'client_id', 'client_secret', 'company_id'],
      status: 'available',
      sampleData: {
        tenant_url: 'https://businesscentral.dynamics.com/',
        client_id: 'app_12345...',
        client_secret: '••••••••',
        company_id: 'CRONUS-Srbija',
      },
      features: [
        'Cloud i On-premise',
        'Srpski jezik',
        'PDV podrška',
        'Power BI integracija',
      ],
      lastSync: 'Pre 1 sat',
      productsCount: 4532,
    },
    {
      name: 'Halo Oglasi',
      icon: '🔷',
      description: 'Import oglasa sa popularne platforme za oglase',
      fields: ['username', 'password', 'location'],
      status: 'coming_soon',
      sampleData: {
        username: 'moj_nalog',
        password: '••••••••',
        location: 'Beograd',
      },
      features: ['Web scraping', 'Location filter', 'Auto refresh'],
      lastSync: 'Uskoro',
      productsCount: 0,
    },
    {
      name: 'WooCommerce',
      icon: '🏬',
      description: 'Import iz postojećeg WordPress/WooCommerce sajta',
      fields: ['shop_url', 'consumer_key', 'consumer_secret'],
      status: 'available',
      sampleData: {
        shop_url: 'https://moj-shop.rs',
        consumer_key: 'ck_12345...',
        consumer_secret: 'cs_67890...',
      },
      features: ['REST API v3', 'Webhook support', 'Varijante proizvoda'],
      lastSync: 'Pre 45 minuta',
      productsCount: 567,
    },
  ];

  const importMethods = [
    {
      id: 'csv',
      title: '📊 CSV fajlovi',
      description: 'Uvoz proizvoda iz CSV sa podešavanjem polja',
      features: [
        'Auto prepoznavanje separatora',
        'Pregled podataka',
        'Mapiranje polja',
      ],
      icon: '📊',
      color: 'bg-gradient-to-r from-green-500 to-emerald-500',
    },
    {
      id: 'xml',
      title: '📋 XML/Digital Vision katalozi',
      description: 'Uvoz iz srpskih digitalnih kataloga i XML formata',
      features: [
        'Validacija šeme',
        'Ugneždene kategorije',
        'Auto mapiranje polja',
      ],
      icon: '📋',
      color: 'bg-gradient-to-r from-blue-500 to-cyan-500',
    },
    {
      id: 'excel',
      title: '📈 Excel fajlovi',
      description: 'Podrška za .xlsx, .xls sa više listova',
      features: ['Izbor lista', 'Formule i stilovi', 'Batch obrada'],
      icon: '📈',
      color: 'bg-gradient-to-r from-orange-500 to-red-500',
    },
    {
      id: 'api',
      title: '🔗 API integracije',
      description: 'Povezivanje sa popularnim srpskim platformama',
      features: ['PANTHEON ERP', 'KupujemProdajem', 'Microsoft Dynamics 365'],
      icon: '🔗',
      color: 'bg-gradient-to-r from-purple-500 to-pink-500',
    },
    {
      id: 'ftp',
      title: '📁 FTP/SFTP sinhronizacija',
      description: 'Automatsko ažuriranje kataloga proizvoda',
      features: ['Raspored', 'Inkrementalna ažuriranja', 'Obaveštenja'],
      icon: '📁',
      color: 'bg-gradient-to-r from-indigo-500 to-purple-500',
    },
    {
      id: 'webhook',
      title: '⚡ Webhook obaveštenja',
      description: 'Trenutna sinhronizacija zaliha i cena',
      features: ['Real-time ažuriranja', 'Batch operacije', 'Retry mehanizam'],
      icon: '⚡',
      color: 'bg-gradient-to-r from-yellow-500 to-orange-500',
    },
  ];

  const csvFields = [
    { id: 'name', label: 'Naziv proizvoda', required: true },
    { id: 'price', label: 'Cena (RSD)', required: true },
    { id: 'quantity', label: 'Količina', required: false },
    { id: 'category', label: 'Kategorija', required: true },
    { id: 'description', label: 'Opis', required: false },
    { id: 'brand', label: 'Brend', required: false },
    { id: 'sku', label: 'Šifra proizvoda', required: false },
    { id: 'pdv', label: 'PDV (%)', required: false },
    { id: 'origin_country', label: 'Zemlja porekla', required: false },
  ];

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setUploadedFile(file);
      if (selectedMethod === 'csv') {
        setCsvData(sampleCsvData);
      } else if (selectedMethod === 'xml') {
        setXmlData(sampleXmlData);
      } else if (selectedMethod === 'excel') {
        setExcelSheets(sampleExcelSheets);
      }
      setCurrentStep(2);
    }
  };

  const handleMethodSelect = (methodId: string) => {
    setSelectedMethod(methodId);
    if (methodId === 'xml') {
      setXmlData(sampleXmlData);
      setCurrentStep(2);
    } else if (methodId === 'api') {
      setCurrentStep(2);
    } else if (methodId === 'ftp') {
      setCurrentStep(2);
    } else if (methodId === 'webhook') {
      setCurrentStep(2);
    }
  };

  const handleFieldMapping = (csvColumn: string, systemField: string) => {
    setFieldMapping((prev) => ({
      ...prev,
      [csvColumn]: systemField,
    }));
  };

  const generatePreview = () => {
    let mapped: any[] = [];

    // Генерируем реалистичные сербские товары в зависимости от метода импорта
    if (selectedMethod === 'csv') {
      mapped = [
        {
          name: 'Samsung Galaxy S24 Ultra',
          price: '189999',
          stock: 25,
          category: 'Mobilni telefoni',
          sku: 'SAM-S24U-256',
          brand: 'Samsung',
          pdv: 38000,
        },
        {
          name: 'iPhone 15 Pro Max',
          price: '219999',
          stock: 18,
          category: 'Mobilni telefoni',
          sku: 'APL-15PM-512',
          brand: 'Apple',
          pdv: 44000,
        },
        {
          name: 'Bosch WAX32EH0BY Mašina za pranje veša',
          price: '89999',
          stock: 12,
          category: 'Bela tehnika',
          sku: 'BSH-WAX32',
          brand: 'Bosch',
          pdv: 18000,
        },
        {
          name: 'LG OLED55C3PUA Smart TV',
          price: '169999',
          stock: 8,
          category: 'TV i Audio',
          sku: 'LG-OLED55C3',
          brand: 'LG',
          pdv: 34000,
        },
        {
          name: 'Sony PlayStation 5 Slim',
          price: '69999',
          stock: 35,
          category: 'Gaming',
          sku: 'SNY-PS5S-1TB',
          brand: 'Sony',
          pdv: 14000,
        },
        {
          name: 'Gorenje NRK6193TX Kombinovani frižider',
          price: '79999',
          stock: 6,
          category: 'Bela tehnika',
          sku: 'GOR-NRK6193',
          brand: 'Gorenje',
          pdv: 16000,
        },
        {
          name: 'Dell XPS 15 Laptop',
          price: '249999',
          stock: 10,
          category: 'Laptopovi',
          sku: 'DELL-XPS15-I7',
          brand: 'Dell',
          pdv: 50000,
        },
        {
          name: 'Xiaomi Mi Band 8',
          price: '5999',
          stock: 120,
          category: 'Pametni satovi',
          sku: 'XMI-MB8',
          brand: 'Xiaomi',
          pdv: 1200,
        },
      ];
    } else if (selectedMethod === 'xml') {
      mapped = [
        {
          name: 'Electrolux EES47320L Mašina za sudove',
          price: '94999',
          stock: 15,
          category: 'Bela tehnika',
          sku: 'ELX-EES47',
          brand: 'Electrolux',
          pdv: 19000,
        },
        {
          name: 'Lenovo ThinkPad X1 Carbon',
          price: '289999',
          stock: 7,
          category: 'Laptopovi',
          sku: 'LEN-X1C-G11',
          brand: 'Lenovo',
          pdv: 58000,
        },
        {
          name: 'Asus ROG Strix Gaming stolica',
          price: '49999',
          stock: 22,
          category: 'Gaming oprema',
          sku: 'ASUS-ROG-CH',
          brand: 'Asus',
          pdv: 10000,
        },
        {
          name: 'JBL Charge 5 Bluetooth zvučnik',
          price: '18999',
          stock: 45,
          category: 'Audio',
          sku: 'JBL-CH5-BLK',
          brand: 'JBL',
          pdv: 3800,
        },
        {
          name: 'Philips Airfryer XXL',
          price: '32999',
          stock: 28,
          category: 'Mali kućni aparati',
          sku: 'PHL-AF-XXL',
          brand: 'Philips',
          pdv: 6600,
        },
        {
          name: 'Canon EOS R6 Mark II',
          price: '399999',
          stock: 4,
          category: 'Foto oprema',
          sku: 'CAN-R6M2',
          brand: 'Canon',
          pdv: 80000,
        },
      ];
    } else if (selectedMethod === 'excel') {
      mapped = [
        {
          name: 'Midea MDV-V28G Klima uređaj',
          price: '54999',
          stock: 18,
          category: 'Klimatizacija',
          sku: 'MID-V28G',
          brand: 'Midea',
          pdv: 11000,
        },
        {
          name: 'HP Pavilion 15 Laptop',
          price: '84999',
          stock: 20,
          category: 'Laptopovi',
          sku: 'HP-PAV15-R5',
          brand: 'HP',
          pdv: 17000,
        },
        {
          name: 'Garmin Fenix 7X Solar',
          price: '119999',
          stock: 9,
          category: 'Pametni satovi',
          sku: 'GAR-F7XS',
          brand: 'Garmin',
          pdv: 24000,
        },
        {
          name: 'DeLonghi Magnifica S',
          price: '62999',
          stock: 14,
          category: 'Kafe aparati',
          sku: 'DLG-MAGS',
          brand: 'DeLonghi',
          pdv: 12600,
        },
        {
          name: 'Beko B5RCNA406HXB Frižider',
          price: '114999',
          stock: 5,
          category: 'Bela tehnika',
          sku: 'BEK-B5RC',
          brand: 'Beko',
          pdv: 23000,
        },
      ];
    } else if (selectedMethod === 'api') {
      mapped = [
        {
          name: 'Apple MacBook Air M3',
          price: '189999',
          stock: 12,
          category: 'Laptopovi',
          sku: 'APL-MBA-M3',
          brand: 'Apple',
          pdv: 38000,
        },
        {
          name: 'Samsung Neo QLED 65"',
          price: '229999',
          stock: 6,
          category: 'TV i Audio',
          sku: 'SAM-QN65',
          brand: 'Samsung',
          pdv: 46000,
        },
        {
          name: 'Dyson V15 Detect Usisivač',
          price: '89999',
          stock: 16,
          category: 'Mali kućni aparati',
          sku: 'DYS-V15D',
          brand: 'Dyson',
          pdv: 18000,
        },
        {
          name: 'Microsoft Xbox Series X',
          price: '64999',
          stock: 25,
          category: 'Gaming',
          sku: 'MS-XSX-1TB',
          brand: 'Microsoft',
          pdv: 13000,
        },
        {
          name: 'Huawei Watch GT 4',
          price: '32999',
          stock: 38,
          category: 'Pametni satovi',
          sku: 'HUA-GT4',
          brand: 'Huawei',
          pdv: 6600,
        },
        {
          name: 'Nikon Z9 Mirrorless',
          price: '649999',
          stock: 2,
          category: 'Foto oprema',
          sku: 'NIK-Z9',
          brand: 'Nikon',
          pdv: 130000,
        },
      ];
    } else if (selectedMethod === 'ftp') {
      mapped = [
        {
          name: 'Vox WMI1495T14A Mašina za pranje',
          price: '44999',
          stock: 20,
          category: 'Bela tehnika',
          sku: 'VOX-WMI1495',
          brand: 'Vox',
          pdv: 9000,
        },
        {
          name: 'Tesla TV 55" 4K Smart',
          price: '59999',
          stock: 15,
          category: 'TV i Audio',
          sku: 'TSL-55K620',
          brand: 'Tesla',
          pdv: 12000,
        },
        {
          name: 'Razer DeathAdder V3 Pro',
          price: '22999',
          stock: 42,
          category: 'Gaming oprema',
          sku: 'RZR-DAV3P',
          brand: 'Razer',
          pdv: 4600,
        },
        {
          name: 'Whirlpool W7 821O OX Frižider',
          price: '99999',
          stock: 8,
          category: 'Bela tehnika',
          sku: 'WHP-W7821',
          brand: 'Whirlpool',
          pdv: 20000,
        },
        {
          name: 'iPad Pro 12.9" M2',
          price: '169999',
          stock: 11,
          category: 'Tableti',
          sku: 'APL-IPADP13',
          brand: 'Apple',
          pdv: 34000,
        },
      ];
    } else if (selectedMethod === 'webhook') {
      mapped = [
        {
          name: 'Siemens iQ700 Ugradna rerna',
          price: '129999',
          stock: 7,
          category: 'Ugradna tehnika',
          sku: 'SIE-IQ700',
          brand: 'Siemens',
          pdv: 26000,
        },
        {
          name: 'OnePlus 12 Pro',
          price: '134999',
          stock: 19,
          category: 'Mobilni telefoni',
          sku: 'OP-12PRO',
          brand: 'OnePlus',
          pdv: 27000,
        },
        {
          name: 'Bose QuietComfort Ultra',
          price: '54999',
          stock: 24,
          category: 'Audio',
          sku: 'BOSE-QCU',
          brand: 'Bose',
          pdv: 11000,
        },
        {
          name: 'Indesit IWC 71252 Mašina za veš',
          price: '34999',
          stock: 30,
          category: 'Bela tehnika',
          sku: 'IND-IWC71',
          brand: 'Indesit',
          pdv: 7000,
        },
        {
          name: 'GoPro HERO12 Black',
          price: '64999',
          stock: 18,
          category: 'Akcione kamere',
          sku: 'GP-H12B',
          brand: 'GoPro',
          pdv: 13000,
        },
        {
          name: 'AEG L8FEC68S ProSteam',
          price: '149999',
          stock: 4,
          category: 'Bela tehnika',
          sku: 'AEG-L8FEC',
          brand: 'AEG',
          pdv: 30000,
        },
        {
          name: 'Marshall Stanmore III',
          price: '44999',
          stock: 16,
          category: 'Audio',
          sku: 'MAR-STN3',
          brand: 'Marshall',
          pdv: 9000,
        },
      ];
    }

    // Додајемо додатне податке за сваки производ
    mapped = mapped.map((item, index) => ({
      ...item,
      id: index + 1,
      description: `Висококвалитетан производ од реномираног произвођача ${item.brand}`,
      importDate: new Date().toLocaleDateString('sr-RS'),
      status: 'active',
    }));

    setPreviewData(mapped);
    setCurrentStep(3);
  };

  const simulateImport = () => {
    setIsImporting(true);
    setImportProgress(0);

    const interval = setInterval(() => {
      setImportProgress((prev) => {
        if (prev >= 100) {
          clearInterval(interval);
          setIsImporting(false);
          setCurrentStep(4);
          return 100;
        }
        return prev + 10;
      });
    }, 200);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-100 to-base-200">
      {/* Header */}
      <div className="navbar bg-base-100 shadow-lg">
        <div className="navbar-start">
          <SveTuLogoStatic variant="gradient" width={120} height={40} />
        </div>
        <div className="navbar-center">
          <h1 className="text-xl font-bold">
            🏪 Demo: Uvoz proizvoda u prodavnice
          </h1>
        </div>
        <div className="navbar-end">
          <div className="badge badge-primary badge-lg">SRBIJA DEMO</div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-6 max-w-7xl">
        {/* Progress Steps */}
        <AnimatedSection animation="fadeIn">
          <div className="steps steps-horizontal mb-8 w-full">
            <div className={`step ${currentStep >= 1 ? 'step-primary' : ''}`}>
              Izbor metoda
            </div>
            <div className={`step ${currentStep >= 2 ? 'step-primary' : ''}`}>
              Podešavanje
            </div>
            <div className={`step ${currentStep >= 3 ? 'step-primary' : ''}`}>
              Pregled
            </div>
            <div className={`step ${currentStep >= 4 ? 'step-primary' : ''}`}>
              Uvoz
            </div>
          </div>
        </AnimatedSection>

        {/* Step 1: Method Selection */}
        {currentStep === 1 && (
          <AnimatedSection animation="fadeIn">
            <div className="text-center mb-8">
              <h2 className="text-3xl font-bold mb-4">
                Izaberite način uvoza proizvoda
              </h2>
              <p className="text-lg text-base-content/70">
                Podržavamo više načina uvoza podataka prilagodenih srpskom
                tržištu i lokalnim potrebama
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {importMethods.map((method, index) => (
                <AnimatedSection
                  key={method.id}
                  animation="slideUp"
                  delay={index * 0.1}
                >
                  <div
                    className={`card bg-base-100 shadow-xl hover:shadow-2xl transition-all duration-300 hover:-translate-y-2 cursor-pointer ${
                      selectedMethod === method.id ? 'ring-2 ring-primary' : ''
                    }`}
                    onClick={() => handleMethodSelect(method.id)}
                  >
                    <div className="card-body">
                      <div
                        className={`w-16 h-16 rounded-xl ${method.color} flex items-center justify-center text-3xl mb-4`}
                      >
                        {method.icon}
                      </div>
                      <h3 className="card-title text-lg">{method.title}</h3>
                      <p className="text-base-content/70 mb-4">
                        {method.description}
                      </p>

                      <div className="space-y-2">
                        {method.features.map((feature, idx) => (
                          <div key={idx} className="flex items-center gap-2">
                            <svg
                              className="w-4 h-4 text-success"
                              fill="currentColor"
                              viewBox="0 0 20 20"
                            >
                              <path
                                fillRule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                                clipRule="evenodd"
                              />
                            </svg>
                            <span className="text-sm">{feature}</span>
                          </div>
                        ))}
                      </div>

                      {selectedMethod === method.id && (
                        <div className="card-actions justify-end mt-4">
                          {(method.id === 'csv' ||
                            method.id === 'xml' ||
                            method.id === 'excel') && (
                            <div className="form-control w-full">
                              <label className="label">
                                <span className="label-text">
                                  {method.id === 'csv' && 'Učitajte CSV fajl'}
                                  {method.id === 'xml' &&
                                    'Загрузите XML/YML файл'}
                                  {method.id === 'excel' &&
                                    'Загрузите Excel файл'}
                                </span>
                              </label>
                              <input
                                type="file"
                                accept={
                                  method.id === 'csv'
                                    ? '.csv'
                                    : method.id === 'xml'
                                      ? '.xml,.yml'
                                      : '.xlsx,.xls'
                                }
                                onChange={handleFileUpload}
                                className="file-input file-input-bordered file-input-primary w-full"
                              />
                            </div>
                          )}
                          {(method.id === 'api' ||
                            method.id === 'ftp' ||
                            method.id === 'webhook') && (
                            <button
                              className="btn btn-primary"
                              onClick={() => setCurrentStep(2)}
                            >
                              Настроить подключение
                            </button>
                          )}
                        </div>
                      )}
                    </div>
                  </div>
                </AnimatedSection>
              ))}
            </div>
          </AnimatedSection>
        )}

        {/* Step 2: CSV Field Mapping */}
        {currentStep === 2 && selectedMethod === 'csv' && (
          <AnimatedSection animation="fadeIn">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* CSV Preview */}
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="card-title">📊 Pregled CSV fajla</h3>
                  <p className="text-base-content/70 mb-4">
                    Fajl: {uploadedFile?.name} ({csvData.length - 1} redova)
                  </p>

                  <div className="overflow-x-auto">
                    <table className="table table-zebra">
                      <thead>
                        <tr>
                          {csvData[0]?.map((header, index) => (
                            <th key={index} className="text-xs">
                              {header}
                            </th>
                          ))}
                        </tr>
                      </thead>
                      <tbody>
                        {csvData.slice(1, 4).map((row, rowIndex) => (
                          <tr key={rowIndex}>
                            {row.map((cell, cellIndex) => (
                              <td
                                key={cellIndex}
                                className="text-xs max-w-32 truncate"
                              >
                                {cell}
                              </td>
                            ))}
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>

              {/* Field Mapping */}
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="card-title">🎯 Mapiranje polja</h3>
                  <p className="text-base-content/70 mb-4">
                    Navedite koje CSV kolone odgovaraju poljima proizvoda
                  </p>

                  <div className="space-y-4">
                    {csvData[0]?.map((csvColumn, index) => (
                      <div key={index} className="form-control">
                        <label className="label">
                          <span className="label-text font-medium">
                            {csvColumn}
                          </span>
                        </label>
                        <select
                          className="select select-bordered"
                          value={fieldMapping[csvColumn] || ''}
                          onChange={(e) =>
                            handleFieldMapping(csvColumn, e.target.value)
                          }
                        >
                          <option value="">-- Ne koristiti --</option>
                          {csvFields.map((field) => (
                            <option key={field.id} value={field.id}>
                              {field.label} {field.required && '(obavezno)'}
                            </option>
                          ))}
                        </select>
                      </div>
                    ))}
                  </div>

                  <div className="card-actions justify-end mt-6">
                    <button
                      className="btn btn-ghost"
                      onClick={() => setCurrentStep(1)}
                    >
                      Nazad
                    </button>
                    <button
                      className="btn btn-primary"
                      onClick={generatePreview}
                      disabled={!fieldMapping.name || !fieldMapping.price}
                    >
                      Kreiraj pregled
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </AnimatedSection>
        )}

        {/* Step 2: XML Configuration */}
        {currentStep === 2 && selectedMethod === 'xml' && (
          <AnimatedSection animation="fadeIn">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* XML Preview */}
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="card-title">📋 XML/YML struktura</h3>
                  <p className="text-base-content/70 mb-4">
                    Pronađeno: {xmlData?.shop?.offers?.length || 0} proizvoda,{' '}
                    {xmlData?.shop?.categories?.length || 0} kategorija
                  </p>

                  <div className="space-y-4">
                    <div className="collapse collapse-arrow bg-base-200">
                      <input type="checkbox" defaultChecked />
                      <div className="collapse-title font-medium">
                        📁 Informacije o prodavnici
                      </div>
                      <div className="collapse-content">
                        <div className="space-y-2">
                          <p>
                            <strong>Naziv:</strong> {xmlData?.shop?.name}
                          </p>
                          <p>
                            <strong>Kompanija:</strong> {xmlData?.shop?.company}
                          </p>
                          <p>
                            <strong>URL:</strong> {xmlData?.shop?.url}
                          </p>
                        </div>
                      </div>
                    </div>

                    <div className="collapse collapse-arrow bg-base-200">
                      <input type="checkbox" />
                      <div className="collapse-title font-medium">
                        📂 Kategorije ({xmlData?.shop?.categories?.length})
                      </div>
                      <div className="collapse-content">
                        <div className="max-h-32 overflow-y-auto space-y-1">
                          {xmlData?.shop?.categories?.map((cat: any) => (
                            <div key={cat.id} className="text-sm">
                              • {cat.name} (ID: {cat.id})
                            </div>
                          ))}
                        </div>
                      </div>
                    </div>

                    <div className="collapse collapse-arrow bg-base-200">
                      <input type="checkbox" />
                      <div className="collapse-title font-medium">
                        📦 Proizvodi ({xmlData?.shop?.offers?.length})
                      </div>
                      <div className="collapse-content">
                        <div className="max-h-32 overflow-y-auto space-y-1">
                          {xmlData?.shop?.offers
                            ?.slice(0, 3)
                            .map((offer: any) => (
                              <div key={offer.id} className="text-sm">
                                • {offer.name} - €{offer.price}
                              </div>
                            ))}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* XML Settings */}
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="card-title">⚙️ Podešavanja XML uvoza</h3>

                  <div className="space-y-4">
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Ažuriranje proizvoda</span>
                      </label>
                      <select className="select select-bordered">
                        <option>Kreiraj nove + ažuriraj postojeće</option>
                        <option>Samo kreiraj nove</option>
                        <option>Samo ažuriraj postojeće</option>
                      </select>
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Kreiranje kategorija</span>
                      </label>
                      <label className="cursor-pointer label">
                        <span className="label-text">
                          Automatski kreiraj kategorije koje nedostaju
                        </span>
                        <input
                          type="checkbox"
                          defaultChecked
                          className="checkbox checkbox-primary"
                        />
                      </label>
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Obrada slika</span>
                      </label>
                      <label className="cursor-pointer label">
                        <span className="label-text">
                          Učitaj slike preko linkova
                        </span>
                        <input
                          type="checkbox"
                          defaultChecked
                          className="checkbox checkbox-primary"
                        />
                      </label>
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Validacija podataka</span>
                      </label>
                      <label className="cursor-pointer label">
                        <span className="label-text">
                          Striktna provera cena i šifri
                        </span>
                        <input
                          type="checkbox"
                          defaultChecked
                          className="checkbox checkbox-primary"
                        />
                      </label>
                    </div>
                  </div>

                  <div className="card-actions justify-end mt-6">
                    <button
                      className="btn btn-ghost"
                      onClick={() => setCurrentStep(1)}
                    >
                      Nazad
                    </button>
                    <button
                      className="btn btn-primary"
                      onClick={generatePreview}
                    >
                      Kreiraj pregled
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </AnimatedSection>
        )}

        {/* Step 2: Excel Configuration */}
        {currentStep === 2 && selectedMethod === 'excel' && (
          <AnimatedSection animation="fadeIn">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* Excel Sheets */}
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="card-title">📈 Listovi Excel fajla</h3>

                  <div className="space-y-3">
                    {excelSheets.map((sheet, index) => (
                      <div key={index} className="form-control">
                        <label className="cursor-pointer">
                          <input
                            type="radio"
                            name="excel-sheet"
                            className="radio radio-primary mr-3"
                            value={sheet}
                            checked={selectedSheet === sheet}
                            onChange={(e) => setSelectedSheet(e.target.value)}
                          />
                          <span className="label-text font-medium">
                            {sheet}
                          </span>
                        </label>
                        <div className="text-xs text-base-content/60 ml-6">
                          {sheet === 'Proizvodi' &&
                            'Osnovne informacije o proizvodima'}
                          {sheet === 'Zalihe' && 'Količina na lagerima'}
                          {sheet === 'Cene' && 'Trenutne cene i popusti'}
                          {sheet === 'Karakteristike' &&
                            'Detaljne karakteristike proizvoda'}
                        </div>
                      </div>
                    ))}
                  </div>

                  <div className="mt-4 p-3 bg-info/10 rounded-lg">
                    <div className="flex items-center gap-2">
                      <svg
                        className="w-4 h-4 text-info"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                          clipRule="evenodd"
                        />
                      </svg>
                      <span className="text-sm font-medium">Savet:</span>
                    </div>
                    <p className="text-sm mt-1">
                      Izaberite list sa osnovnim podacima o proizvodima. Ostale
                      listove možete povezati kasnije.
                    </p>
                  </div>
                </div>
              </div>

              {/* Excel Preview */}
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="card-title">👁️ Pregled podataka</h3>
                  {selectedSheet ? (
                    <div className="space-y-4">
                      <p className="text-base-content/70">
                        List: <strong>{selectedSheet}</strong>
                      </p>

                      <div className="overflow-x-auto">
                        <table className="table table-zebra table-xs">
                          <thead>
                            <tr>
                              <th>Колонка A</th>
                              <th>Колонка B</th>
                              <th>Колонка C</th>
                              <th>Колонка D</th>
                            </tr>
                          </thead>
                          <tbody>
                            <tr>
                              <td>Название</td>
                              <td>Cena</td>
                              <td>Количество</td>
                              <td>Kategorija</td>
                            </tr>
                            <tr>
                              <td>Proizvod 1</td>
                              <td>599</td>
                              <td>12</td>
                              <td>Electronics</td>
                            </tr>
                            <tr>
                              <td>Proizvod 2</td>
                              <td>799</td>
                              <td>8</td>
                              <td>Computers</td>
                            </tr>
                          </tbody>
                        </table>
                      </div>

                      <div className="stats stats-vertical lg:stats-horizontal shadow">
                        <div className="stat">
                          <div className="stat-title">Строк</div>
                          <div className="stat-value text-sm">1,247</div>
                        </div>
                        <div className="stat">
                          <div className="stat-title">Колонок</div>
                          <div className="stat-value text-sm">12</div>
                        </div>
                        <div className="stat">
                          <div className="stat-title">Размер</div>
                          <div className="stat-value text-sm">2.3MB</div>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div className="text-center py-8">
                      <div className="text-6xl mb-4">📊</div>
                      <p className="text-base-content/60">
                        Izaberite list za pregled
                      </p>
                    </div>
                  )}

                  <div className="card-actions justify-end mt-6">
                    <button
                      className="btn btn-ghost"
                      onClick={() => setCurrentStep(1)}
                    >
                      Nazad
                    </button>
                    <button
                      className="btn btn-primary"
                      onClick={generatePreview}
                      disabled={!selectedSheet}
                    >
                      Kreiraj pregled
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </AnimatedSection>
        )}

        {/* Step 2: API Configuration */}
        {currentStep === 2 && selectedMethod === 'api' && (
          <AnimatedSection animation="fadeIn">
            <div className="space-y-6">
              <div className="text-center">
                <h2 className="text-2xl font-bold mb-2">
                  🔗 Podešavanje API integracija
                </h2>
                <p className="text-base-content/70">
                  Izaberite platformu i podešavajte konekciju sa srpskim ERP
                  sistemima
                </p>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                {sampleApiIntegrations.map((integration, index) => (
                  <AnimatedSection
                    key={integration.name}
                    animation="slideUp"
                    delay={index * 0.1}
                  >
                    <div
                      className={`card bg-base-100 shadow-xl h-full ${integration.status === 'coming_soon' ? 'opacity-60' : 'hover:shadow-2xl'} transition-all duration-300`}
                    >
                      <div className="card-body">
                        <div className="flex items-center gap-4 mb-4">
                          <div className="text-4xl">{integration.icon}</div>
                          <div className="flex-1">
                            <div className="flex items-center gap-2">
                              <h3 className="card-title text-lg">
                                {integration.name}
                              </h3>
                              {integration.status === 'coming_soon' ? (
                                <div className="badge badge-warning badge-sm">
                                  Soon
                                </div>
                              ) : (
                                <div className="badge badge-success badge-sm">
                                  Active
                                </div>
                              )}
                            </div>
                            <p className="text-sm text-base-content/70">
                              {integration.description}
                            </p>
                          </div>
                        </div>

                        {integration.status === 'available' && (
                          <>
                            {/* Statistics */}
                            <div className="stats stats-vertical bg-base-200 shadow mb-4">
                              <div className="stat py-2">
                                <div className="stat-title text-xs">
                                  Proizvoda
                                </div>
                                <div className="stat-value text-lg">
                                  {integration.productsCount.toLocaleString()}
                                </div>
                              </div>
                              <div className="stat py-2">
                                <div className="stat-title text-xs">
                                  Poslednja sinhronizacija
                                </div>
                                <div className="stat-value text-sm">
                                  {integration.lastSync}
                                </div>
                              </div>
                            </div>

                            {/* Features */}
                            <div className="mb-4">
                              <div className="text-sm font-medium mb-2">
                                Mogućnosti:
                              </div>
                              <div className="flex flex-wrap gap-1">
                                {integration.features.map((feature, idx) => (
                                  <span
                                    key={idx}
                                    className="badge badge-outline badge-xs"
                                  >
                                    {feature}
                                  </span>
                                ))}
                              </div>
                            </div>

                            {/* Configuration Fields */}
                            <div className="collapse collapse-arrow bg-base-200">
                              <input type="checkbox" />
                              <div className="collapse-title text-sm font-medium">
                                Podešavanja konekcije
                              </div>
                              <div className="collapse-content">
                                <div className="space-y-3">
                                  {integration.fields.map((field) => (
                                    <div key={field} className="form-control">
                                      <label className="label py-1">
                                        <span className="label-text text-xs capitalize">
                                          {field.replace('_', ' ')}
                                        </span>
                                      </label>
                                      <input
                                        type={
                                          field.includes('password') ||
                                          field.includes('secret') ||
                                          field.includes('token')
                                            ? 'password'
                                            : 'text'
                                        }
                                        placeholder={
                                          integration.sampleData[
                                            field as keyof typeof integration.sampleData
                                          ] ||
                                          `Unesite ${field.replace('_', ' ')}`
                                        }
                                        className="input input-bordered input-sm"
                                        value={apiSettings[field] || ''}
                                        onChange={(e) =>
                                          setApiSettings((prev) => ({
                                            ...prev,
                                            [field]: e.target.value,
                                          }))
                                        }
                                      />
                                    </div>
                                  ))}
                                </div>
                              </div>
                            </div>

                            <div className="card-actions justify-between mt-4">
                              <button className="btn btn-xs btn-outline">
                                Test
                              </button>
                              <button
                                className="btn btn-xs btn-primary"
                                onClick={generatePreview}
                              >
                                Sinhronizuj
                              </button>
                            </div>
                          </>
                        )}

                        {integration.status === 'coming_soon' && (
                          <div className="text-center py-4">
                            <div className="text-2xl mb-2">🚧</div>
                            <p className="text-sm text-base-content/60">
                              Integracija u razvoju
                            </p>
                            <button className="btn btn-ghost btn-xs mt-2">
                              Obavesti o pokretanju
                            </button>
                          </div>
                        )}
                      </div>
                    </div>
                  </AnimatedSection>
                ))}
              </div>

              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="card-title">⚙️ Dodatna podešavanja</h3>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-3">
                      <h4 className="font-semibold">Sinhronizacija</h4>
                      <label className="cursor-pointer label">
                        <span className="label-text">
                          Automatsko ažuriranje cena
                        </span>
                        <input
                          type="checkbox"
                          className="checkbox checkbox-primary"
                        />
                      </label>
                      <label className="cursor-pointer label">
                        <span className="label-text">
                          Sinhronizacija zaliha
                        </span>
                        <input
                          type="checkbox"
                          defaultChecked
                          className="checkbox checkbox-primary"
                        />
                      </label>
                      <label className="cursor-pointer label">
                        <span className="label-text">Uvoz novih proizvoda</span>
                        <input
                          type="checkbox"
                          defaultChecked
                          className="checkbox checkbox-primary"
                        />
                      </label>
                    </div>

                    <div className="space-y-3">
                      <h4 className="font-semibold">Raspored</h4>
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">
                            Učestalost sinhronizacije
                          </span>
                        </label>
                        <select className="select select-bordered">
                          <option>Svaki sat</option>
                          <option>Svakih 6 sati</option>
                          <option>Dnevno</option>
                          <option>Nedeljno</option>
                          <option>Вручную</option>
                        </select>
                      </div>
                    </div>
                  </div>

                  <div className="card-actions justify-center mt-6">
                    <button
                      className="btn btn-ghost"
                      onClick={() => setCurrentStep(1)}
                    >
                      Nazad na izbor metoda
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </AnimatedSection>
        )}

        {/* Step 2: FTP Configuration */}
        {currentStep === 2 && selectedMethod === 'ftp' && (
          <AnimatedSection animation="fadeIn">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* FTP Settings */}
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="card-title">📁 Настройки FTP/SFTP</h3>

                  <div className="space-y-4">
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Тип подключения</span>
                      </label>
                      <div className="flex gap-4">
                        <label className="cursor-pointer">
                          <input
                            type="radio"
                            name="ftp-type"
                            value="ftp"
                            className="radio radio-primary mr-2"
                            defaultChecked
                          />
                          FTP
                        </label>
                        <label className="cursor-pointer">
                          <input
                            type="radio"
                            name="ftp-type"
                            value="sftp"
                            className="radio radio-primary mr-2"
                          />
                          SFTP
                        </label>
                        <label className="cursor-pointer">
                          <input
                            type="radio"
                            name="ftp-type"
                            value="ftps"
                            className="radio radio-primary mr-2"
                          />
                          FTPS
                        </label>
                      </div>
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Сервер</span>
                      </label>
                      <input
                        type="text"
                        placeholder="ftp.example.com"
                        className="input input-bordered"
                        value={ftpSettings.server || ''}
                        onChange={(e) =>
                          setFtpSettings((prev) => ({
                            ...prev,
                            server: e.target.value,
                          }))
                        }
                      />
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Port</span>
                        </label>
                        <input
                          type="number"
                          placeholder="21"
                          className="input input-bordered"
                          value={ftpSettings.port || ''}
                          onChange={(e) =>
                            setFtpSettings((prev) => ({
                              ...prev,
                              port: e.target.value,
                            }))
                          }
                        />
                      </div>
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Папка</span>
                        </label>
                        <input
                          type="text"
                          placeholder="/uploads"
                          className="input input-bordered"
                          value={ftpSettings.folder || ''}
                          onChange={(e) =>
                            setFtpSettings((prev) => ({
                              ...prev,
                              folder: e.target.value,
                            }))
                          }
                        />
                      </div>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Логин</span>
                        </label>
                        <input
                          type="text"
                          placeholder="username"
                          className="input input-bordered"
                          value={ftpSettings.username || ''}
                          onChange={(e) =>
                            setFtpSettings((prev) => ({
                              ...prev,
                              username: e.target.value,
                            }))
                          }
                        />
                      </div>
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Lozinka</span>
                        </label>
                        <input
                          type="password"
                          placeholder="••••••••"
                          className="input input-bordered"
                          value={ftpSettings.password || ''}
                          onChange={(e) =>
                            setFtpSettings((prev) => ({
                              ...prev,
                              password: e.target.value,
                            }))
                          }
                        />
                      </div>
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Maska fajlova</span>
                      </label>
                      <input
                        type="text"
                        placeholder="products_*.csv"
                        className="input input-bordered"
                        value={ftpSettings.fileMask || ''}
                        onChange={(e) =>
                          setFtpSettings((prev) => ({
                            ...prev,
                            fileMask: e.target.value,
                          }))
                        }
                      />
                      <div className="label">
                        <span className="label-text-alt">
                          Primer: products_*.csv, catalog_*.xml
                        </span>
                      </div>
                    </div>

                    <div className="flex gap-2">
                      <button className="btn btn-outline btn-sm">
                        Test konekcije
                      </button>
                      <button className="btn btn-outline btn-sm">
                        Pregled fajlova
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              {/* Schedule & Automation */}
              <div className="space-y-6">
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <h3 className="card-title">⏰ Raspored sinhronizacije</h3>

                    <div className="space-y-4">
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Učestalost provere</span>
                        </label>
                        <select className="select select-bordered">
                          <option>Svakih 15 minuta</option>
                          <option>Svaki sat</option>
                          <option>Svakih 6 sati</option>
                          <option>Dnevno u 03:00</option>
                          <option>По будням в 08:00</option>
                          <option>Вручную</option>
                        </select>
                      </div>

                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Obrada fajlova</span>
                        </label>
                        <label className="cursor-pointer label">
                          <span className="label-text">
                            Удалять обработанные файлы
                          </span>
                          <input
                            type="checkbox"
                            className="checkbox checkbox-primary"
                          />
                        </label>
                        <label className="cursor-pointer label">
                          <span className="label-text">Создавать бэкапы</span>
                          <input
                            type="checkbox"
                            defaultChecked
                            className="checkbox checkbox-primary"
                          />
                        </label>
                        <label className="cursor-pointer label">
                          <span className="label-text">
                            Отправлять уведомления
                          </span>
                          <input
                            type="checkbox"
                            defaultChecked
                            className="checkbox checkbox-primary"
                          />
                        </label>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <h3 className="card-title">📊 Statistika i logovi</h3>

                    <div className="stats stats-vertical shadow mb-4">
                      <div className="stat">
                        <div className="stat-title">
                          Poslednja sinhronizacija
                        </div>
                        <div className="stat-value text-sm">Pre 2 sata</div>
                        <div className="stat-desc">Uspešno</div>
                      </div>
                      <div className="stat">
                        <div className="stat-title">Файлов обработано</div>
                        <div className="stat-value text-sm">1,247</div>
                        <div className="stat-desc">За последний месяц</div>
                      </div>
                      <div className="stat">
                        <div className="stat-title">Proizvoda uvezeno</div>
                        <div className="stat-value text-sm">15,892</div>
                        <div className="stat-desc">Всего</div>
                      </div>
                    </div>

                    {/* Activity Log */}
                    <div className="collapse collapse-arrow bg-base-200">
                      <input type="checkbox" />
                      <div className="collapse-title text-sm font-medium">
                        📝 Журнал активности (последние 10 записей)
                      </div>
                      <div className="collapse-content">
                        <div className="max-h-48 overflow-y-auto space-y-2">
                          {[
                            {
                              time: '14:32',
                              status: 'success',
                              message:
                                'Obrađen fajl products_2024_01_15.csv (247 proizvoda)',
                            },
                            {
                              time: '14:30',
                              status: 'info',
                              message:
                                'Konekcija sa FTP serverom uspostavljena',
                            },
                            {
                              time: '12:15',
                              status: 'success',
                              message:
                                'Sinhronizacija završena: +12 proizvoda, ~3 ažuriranja',
                            },
                            {
                              time: '12:10',
                              status: 'warning',
                              message:
                                'Файл prices_old.csv пропущен (устаревший формат)',
                            },
                            {
                              time: '10:45',
                              status: 'success',
                              message:
                                'Uvoz categories_update.xml završen uspešno',
                            },
                            {
                              time: '09:22',
                              status: 'error',
                              message: 'Greška učitanja slike: timeout',
                            },
                            {
                              time: '08:30',
                              status: 'success',
                              message: 'Плановая синхронизация запущена',
                            },
                          ].map((log, idx) => (
                            <div
                              key={idx}
                              className="flex items-center gap-3 p-2 bg-base-100 rounded text-xs"
                            >
                              <div className="text-base-content/50">
                                {log.time}
                              </div>
                              <div
                                className={`w-2 h-2 rounded-full ${
                                  log.status === 'success'
                                    ? 'bg-success'
                                    : log.status === 'warning'
                                      ? 'bg-warning'
                                      : log.status === 'error'
                                        ? 'bg-error'
                                        : 'bg-info'
                                }`}
                              ></div>
                              <div className="flex-1">{log.message}</div>
                            </div>
                          ))}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div className="lg:col-span-2">
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <div className="card-actions justify-center">
                      <button
                        className="btn btn-ghost"
                        onClick={() => setCurrentStep(1)}
                      >
                        Nazad na izbor metoda
                      </button>
                      <button
                        className="btn btn-primary"
                        onClick={generatePreview}
                      >
                        Запустить синхронизацию
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </AnimatedSection>
        )}

        {/* Step 2: Webhook Configuration */}
        {currentStep === 2 && selectedMethod === 'webhook' && (
          <AnimatedSection animation="fadeIn">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* Webhook Settings */}
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="card-title">⚡ Podešavanja Webhook-a</h3>

                  <div className="space-y-4">
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">URL endpoint-a</span>
                      </label>
                      <input
                        type="url"
                        placeholder="https://your-store.com/webhook"
                        className="input input-bordered"
                        value={webhookSettings.url || ''}
                        onChange={(e) =>
                          setWebhookSettings((prev) => ({
                            ...prev,
                            url: e.target.value,
                          }))
                        }
                      />
                      <div className="label">
                        <span className="label-text-alt">
                          URL za prijem obaveštenja o promenama
                        </span>
                      </div>
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Tajni ključ</span>
                      </label>
                      <div className="input-group">
                        <input
                          type="text"
                          placeholder="webhook_secret_key"
                          className="input input-bordered flex-1"
                          value={webhookSettings.secret || ''}
                          onChange={(e) =>
                            setWebhookSettings((prev) => ({
                              ...prev,
                              secret: e.target.value,
                            }))
                          }
                        />
                        <button className="btn btn-square btn-outline">
                          <svg
                            className="w-4 h-4"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                            />
                          </svg>
                        </button>
                      </div>
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Tipovi događaja</span>
                      </label>
                      <div className="space-y-2">
                        <label className="cursor-pointer label">
                          <span className="label-text">
                            Kreiranje proizvoda
                          </span>
                          <input
                            type="checkbox"
                            defaultChecked
                            className="checkbox checkbox-primary"
                          />
                        </label>
                        <label className="cursor-pointer label">
                          <span className="label-text">Ažuriranje cene</span>
                          <input
                            type="checkbox"
                            defaultChecked
                            className="checkbox checkbox-primary"
                          />
                        </label>
                        <label className="cursor-pointer label">
                          <span className="label-text">Promena zaliha</span>
                          <input
                            type="checkbox"
                            defaultChecked
                            className="checkbox checkbox-primary"
                          />
                        </label>
                        <label className="cursor-pointer label">
                          <span className="label-text">Brisanje proizvoda</span>
                          <input
                            type="checkbox"
                            className="checkbox checkbox-primary"
                          />
                        </label>
                      </div>
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Format podataka</span>
                      </label>
                      <select className="select select-bordered">
                        <option>JSON</option>
                        <option>XML</option>
                        <option>Form Data</option>
                      </select>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Tajmaut (sek)</span>
                        </label>
                        <input
                          type="number"
                          placeholder="30"
                          className="input input-bordered"
                          value={webhookSettings.timeout || ''}
                          onChange={(e) =>
                            setWebhookSettings((prev) => ({
                              ...prev,
                              timeout: e.target.value,
                            }))
                          }
                        />
                      </div>
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Ponavljanja</span>
                        </label>
                        <input
                          type="number"
                          placeholder="3"
                          className="input input-bordered"
                          value={webhookSettings.retries || ''}
                          onChange={(e) =>
                            setWebhookSettings((prev) => ({
                              ...prev,
                              retries: e.target.value,
                            }))
                          }
                        />
                      </div>
                    </div>

                    <div className="flex gap-2">
                      <button className="btn btn-outline btn-sm">
                        Test webhook
                      </button>
                      <button className="btn btn-outline btn-sm">
                        Pregled logova
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              {/* Webhook Examples & Testing */}
              <div className="space-y-6">
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <h3 className="card-title">📝 Primer payload</h3>

                    <div className="mockup-code text-xs">
                      <pre>
                        <code>{`{
  "event": "product.updated",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "id": 12345,
    "name": "iPhone 15 Pro",
    "price": 999.00,
    "stock": 15,
    "category": "Electronics",
    "updated_fields": ["price", "stock"]
  },
  "signature": "sha256=abc123..."
}`}</code>
                      </pre>
                    </div>
                  </div>
                </div>

                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <h3 className="card-title">🔄 Status isporuke</h3>

                    <div className="space-y-3">
                      <div className="flex justify-between items-center p-3 bg-success/10 rounded-lg">
                        <div className="flex items-center gap-2">
                          <div className="w-2 h-2 bg-success rounded-full"></div>
                          <span className="text-sm">Kreiranje proizvoda</span>
                        </div>
                        <span className="text-xs text-success">200 OK</span>
                      </div>

                      <div className="flex justify-between items-center p-3 bg-success/10 rounded-lg">
                        <div className="flex items-center gap-2">
                          <div className="w-2 h-2 bg-success rounded-full"></div>
                          <span className="text-sm">Ažuriranje cene</span>
                        </div>
                        <span className="text-xs text-success">200 OK</span>
                      </div>

                      <div className="flex justify-between items-center p-3 bg-error/10 rounded-lg">
                        <div className="flex items-center gap-2">
                          <div className="w-2 h-2 bg-error rounded-full"></div>
                          <span className="text-sm">Promena zaliha</span>
                        </div>
                        <span className="text-xs text-error">500 Error</span>
                      </div>
                    </div>

                    <div className="stats stats-vertical shadow mt-4">
                      <div className="stat">
                        <div className="stat-title">Uspešno</div>
                        <div className="stat-value text-sm text-success">
                          98.5%
                        </div>
                      </div>
                      <div className="stat">
                        <div className="stat-title">Prosečno vreme</div>
                        <div className="stat-value text-sm">147ms</div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div className="lg:col-span-2">
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <div className="card-actions justify-center">
                      <button
                        className="btn btn-ghost"
                        onClick={() => setCurrentStep(1)}
                      >
                        Nazad na izbor metoda
                      </button>
                      <button
                        className="btn btn-primary"
                        onClick={generatePreview}
                      >
                        Aktiviraj webhook
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </AnimatedSection>
        )}

        {/* Step 3: Preview */}
        {currentStep === 3 && (
          <AnimatedSection animation="fadeIn">
            <div className="space-y-6">
              <div className="text-center">
                <h2 className="text-2xl font-bold mb-2">
                  Pregled proizvoda za uvoz
                </h2>
                <p className="text-base-content/70">
                  Proverite podatke pre uvoza. Pronađeno {previewData.length}{' '}
                  proizvoda.
                </p>
              </div>

              {/* Statistics */}
              <div className="stats shadow w-full">
                <div className="stat">
                  <div className="stat-figure text-primary">
                    <svg
                      className="w-8 h-8"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z" />
                    </svg>
                  </div>
                  <div className="stat-title">Ukupno proizvoda</div>
                  <div className="stat-value text-primary">
                    {previewData.length}
                  </div>
                  <div className="stat-desc">Spremno za uvoz</div>
                </div>

                <div className="stat">
                  <div className="stat-figure text-secondary">
                    <svg
                      className="w-8 h-8"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z"
                        clipRule="evenodd"
                      />
                    </svg>
                  </div>
                  <div className="stat-title">Ažuriranja</div>
                  <div className="stat-value text-secondary">0</div>
                  <div className="stat-desc">Postojećih proizvoda</div>
                </div>

                <div className="stat">
                  <div className="stat-figure text-accent">
                    <svg
                      className="w-8 h-8"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                        clipRule="evenodd"
                      />
                    </svg>
                  </div>
                  <div className="stat-title">Validnih</div>
                  <div className="stat-value text-accent">
                    {previewData.length}
                  </div>
                  <div className="stat-desc">Bez grešaka</div>
                </div>
              </div>

              {/* Products Grid */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {previewData.map((product, index) => (
                  <AnimatedSection
                    key={index}
                    animation="slideUp"
                    delay={index * 0.1}
                  >
                    <div className="card bg-base-100 shadow-xl">
                      <div className="card-body">
                        <h3 className="card-title text-lg">{product.name}</h3>
                        <div className="space-y-2">
                          <div className="flex justify-between">
                            <span className="text-base-content/70">Cena:</span>
                            <span className="font-bold">€{product.price}</span>
                          </div>
                          {product.quantity && (
                            <div className="flex justify-between">
                              <span className="text-base-content/70">
                                Količina:
                              </span>
                              <span>{product.quantity}</span>
                            </div>
                          )}
                          {product.category && (
                            <div className="flex justify-between">
                              <span className="text-base-content/70">
                                Kategorija:
                              </span>
                              <span className="text-sm">
                                {product.category}
                              </span>
                            </div>
                          )}
                          {product.brand && (
                            <div className="flex justify-between">
                              <span className="text-base-content/70">
                                Бренд:
                              </span>
                              <span>{product.brand}</span>
                            </div>
                          )}
                        </div>
                        {product.description && (
                          <p className="text-sm text-base-content/70 mt-2">
                            {product.description.substring(0, 100)}...
                          </p>
                        )}
                      </div>
                    </div>
                  </AnimatedSection>
                ))}
              </div>

              {/* Actions */}
              {/* Actions */}
              <div className="flex justify-center gap-4">
                <button
                  className="btn btn-ghost"
                  onClick={() => setCurrentStep(2)}
                  disabled={isImporting}
                >
                  Nazad na podešavanja
                </button>
                <button
                  className={`btn btn-primary btn-lg ${isImporting ? 'loading' : ''}`}
                  onClick={simulateImport}
                  disabled={isImporting}
                >
                  {isImporting
                    ? 'Uvozimo...'
                    : `Počni uvoz (${previewData.length} proizvoda)`}
                </button>
              </div>

              {/* Progress indicator */}
              {isImporting && (
                <AnimatedSection animation="fadeIn">
                  <div className="card bg-base-100 shadow-xl mt-8">
                    <div className="card-body">
                      <h3 className="card-title text-center">
                        🚀 Izvršava se uvoz proizvoda
                      </h3>

                      <div className="space-y-4">
                        <div className="text-center">
                          <div
                            className="radial-progress text-primary text-2xl"
                            style={
                              {
                                '--value': importProgress,
                                '--size': '6rem',
                              } as React.CSSProperties
                            }
                          >
                            {importProgress}%
                          </div>
                        </div>

                        <progress
                          className="progress progress-primary w-full"
                          value={importProgress}
                          max="100"
                        ></progress>

                        <div className="text-center text-sm text-base-content/70">
                          {importProgress < 20 && 'Provera podataka...'}
                          {importProgress >= 20 &&
                            importProgress < 40 &&
                            'Kreiranje kategorija...'}
                          {importProgress >= 40 &&
                            importProgress < 60 &&
                            'Učitavanje slika...'}
                          {importProgress >= 60 &&
                            importProgress < 80 &&
                            'Dodavanje proizvoda...'}
                          {importProgress >= 80 &&
                            importProgress < 100 &&
                            'Finalizacija...'}
                          {importProgress === 100 && 'Završeno!'}
                        </div>

                        <div className="stats stats-horizontal shadow w-full">
                          <div className="stat">
                            <div className="stat-title">Обработано</div>
                            <div className="stat-value text-sm">
                              {Math.floor(
                                (previewData.length * importProgress) / 100
                              )}
                            </div>
                            <div className="stat-desc">
                              из {previewData.length}
                            </div>
                          </div>
                          <div className="stat">
                            <div className="stat-title">Скорость</div>
                            <div className="stat-value text-sm">12.5</div>
                            <div className="stat-desc">proizvoda/sek</div>
                          </div>
                          <div className="stat">
                            <div className="stat-title">Осталось</div>
                            <div className="stat-value text-sm">
                              {Math.max(
                                0,
                                Math.ceil((100 - importProgress) * 0.2)
                              )}
                              с
                            </div>
                            <div className="stat-desc">približno</div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </AnimatedSection>
              )}
            </div>
          </AnimatedSection>
        )}

        {/* Step 4: Success */}
        {currentStep === 4 && (
          <AnimatedSection animation="fadeIn">
            <div>
              {/* Success Header */}
              <div className="text-center mb-8">
                <div className="inline-flex items-center justify-center w-24 h-24 bg-success/20 rounded-full mb-4">
                  <svg
                    className="w-12 h-12 text-success"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                      clipRule="evenodd"
                    />
                  </svg>
                </div>
                <h2 className="text-3xl font-bold text-success mb-2">
                  Uvoz je uspešno završen!
                </h2>
                <p className="text-lg text-base-content/70">
                  Svih {previewData.length} proizvoda je uspešno dodato u vašu
                  prodavnicu
                </p>
              </div>

              {/* Import Statistics */}
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
                <div className="stat bg-base-100 shadow-xl">
                  <div className="stat-figure text-success">
                    <svg
                      className="w-8 h-8"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                  </div>
                  <div className="stat-title">Uspešno uvezeno</div>
                  <div className="stat-value text-success">
                    {previewData.length}
                  </div>
                  <div className="stat-desc">novih proizvoda</div>
                </div>

                <div className="stat bg-base-100 shadow-xl">
                  <div className="stat-figure text-primary">
                    <svg
                      className="w-8 h-8"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                  </div>
                  <div className="stat-title">Ukupna vrednost</div>
                  <div className="stat-value text-primary">
                    {(
                      previewData.reduce((sum, p) => sum + Number(p.price), 0) /
                      1000000
                    ).toFixed(1)}
                    M
                  </div>
                  <div className="stat-desc">RSD bez PDV</div>
                </div>

                <div className="stat bg-base-100 shadow-xl">
                  <div className="stat-figure text-info">
                    <svg
                      className="w-8 h-8"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
                      />
                    </svg>
                  </div>
                  <div className="stat-title">Ukupne zalihe</div>
                  <div className="stat-value text-info">
                    {previewData.reduce((sum, p) => sum + p.stock, 0)}
                  </div>
                  <div className="stat-desc">komada na lageru</div>
                </div>

                <div className="stat bg-base-100 shadow-xl">
                  <div className="stat-figure text-accent">
                    <svg
                      className="w-8 h-8"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M13 10V3L4 14h7v7l9-11h-7z"
                      />
                    </svg>
                  </div>
                  <div className="stat-title">Brzina uvoza</div>
                  <div className="stat-value text-accent">2.3s</div>
                  <div className="stat-desc">
                    {Math.round(previewData.length / 2.3)} proizvoda/sek
                  </div>
                </div>
              </div>

              {/* Category Breakdown */}
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
                {/* Categories */}
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <h3 className="card-title mb-4">
                      📂 Raspodela po kategorijama
                    </h3>
                    <div className="space-y-3">
                      {(() => {
                        // Подсчитываем категории
                        const categoryCount: Record<string, number> = {};
                        previewData.forEach((p) => {
                          categoryCount[p.category] =
                            (categoryCount[p.category] || 0) + 1;
                        });

                        // Сортируем по количеству и берем топ-4
                        const topCategories = Object.entries(categoryCount)
                          .sort(([, a], [, b]) => b - a)
                          .slice(0, 4);

                        const colors = [
                          'bg-blue-500',
                          'bg-green-500',
                          'bg-purple-500',
                          'bg-orange-500',
                        ];

                        return topCategories.map(([category, count], idx) => (
                          <div
                            key={category}
                            className="flex items-center justify-between"
                          >
                            <div className="flex items-center gap-3">
                              <div
                                className={`w-3 h-3 rounded-full ${colors[idx]}`}
                              ></div>
                              <span>{category}</span>
                            </div>
                            <div className="flex items-center gap-2">
                              <span className="font-bold">{count}</span>
                              <span className="text-sm text-base-content/60">
                                (
                                {Math.round((count / previewData.length) * 100)}
                                %)
                              </span>
                            </div>
                          </div>
                        ));
                      })()}
                    </div>
                  </div>
                </div>

                {/* Price Range */}
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <h3 className="card-title mb-4">💰 Raspon cena</h3>
                    <div className="space-y-3">
                      {(() => {
                        const prices = previewData.map((p) => Number(p.price));
                        const minPrice = Math.min(...prices);
                        const maxPrice = Math.max(...prices);
                        const avgPrice =
                          prices.reduce((sum, p) => sum + p, 0) / prices.length;
                        const totalValue = prices.reduce(
                          (sum, p) => sum + p,
                          0
                        );
                        const totalPDV = previewData.reduce(
                          (sum, p) => sum + (p.pdv || 0),
                          0
                        );

                        return (
                          <>
                            <div className="flex items-center justify-between">
                              <span className="text-sm text-base-content/60">
                                Najniža cena:
                              </span>
                              <span className="font-bold">
                                {minPrice.toLocaleString('sr-RS')} RSD
                              </span>
                            </div>
                            <div className="flex items-center justify-between">
                              <span className="text-sm text-base-content/60">
                                Najviša cena:
                              </span>
                              <span className="font-bold">
                                {maxPrice.toLocaleString('sr-RS')} RSD
                              </span>
                            </div>
                            <div className="flex items-center justify-between">
                              <span className="text-sm text-base-content/60">
                                Prosečna cena:
                              </span>
                              <span className="font-bold text-lg text-primary">
                                {Math.round(avgPrice).toLocaleString('sr-RS')}{' '}
                                RSD
                              </span>
                            </div>
                            <div className="divider"></div>
                            <div className="flex items-center justify-between">
                              <span className="text-sm text-base-content/60">
                                Ukupno bez PDV:
                              </span>
                              <span className="font-bold">
                                {totalValue.toLocaleString('sr-RS')} RSD
                              </span>
                            </div>
                            <div className="flex items-center justify-between">
                              <span className="text-sm text-base-content/60">
                                Ukupan PDV (20%):
                              </span>
                              <span className="font-bold text-warning">
                                {totalPDV.toLocaleString('sr-RS')} RSD
                              </span>
                            </div>
                            <div className="flex items-center justify-between">
                              <span className="text-sm text-base-content/60">
                                Ukupno sa PDV:
                              </span>
                              <span className="font-bold text-success text-lg">
                                {(totalValue + totalPDV).toLocaleString(
                                  'sr-RS'
                                )}{' '}
                                RSD
                              </span>
                            </div>
                          </>
                        );
                      })()}
                    </div>
                  </div>
                </div>
              </div>

              {/* Imported Products Preview */}
              <div className="card bg-base-100 shadow-xl mb-8">
                <div className="card-body">
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="card-title">
                      📋 Poslednji uvezeni proizvodi
                    </h3>
                    <button className="btn btn-sm btn-ghost">
                      Pogledaj sve
                    </button>
                  </div>

                  <div className="overflow-x-auto">
                    <table className="table table-zebra">
                      <thead>
                        <tr>
                          <th>Slika</th>
                          <th>Naziv</th>
                          <th>Kategorija</th>
                          <th>Cena</th>
                          <th>Zalihe</th>
                          <th>Status</th>
                        </tr>
                      </thead>
                      <tbody>
                        {previewData.slice(0, 5).map((product, idx) => (
                          <tr key={idx}>
                            <td>
                              <div className="w-12 h-12 bg-base-200 rounded-lg flex items-center justify-center">
                                📦
                              </div>
                            </td>
                            <td>
                              <div className="font-bold">{product.name}</div>
                              <div className="text-sm opacity-50">
                                SKU: {product.sku}
                              </div>
                            </td>
                            <td>{product.category}</td>
                            <td className="font-semibold">
                              {Number(product.price).toLocaleString('sr-RS')}{' '}
                              RSD
                            </td>
                            <td>
                              <span className="badge badge-ghost badge-sm">
                                {product.stock} kom
                              </span>
                            </td>
                            <td>
                              <span className="badge badge-success badge-sm">
                                Aktivan
                              </span>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>

              {/* Actions After Import */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
                <div className="card bg-gradient-to-r from-blue-500 to-blue-600 text-white">
                  <div className="card-body">
                    <h3 className="card-title text-white">🎯 Optimizuj SEO</h3>
                    <p className="text-white/90">
                      Dodaj meta opise i ključne reči za bolje rangiranje
                    </p>
                    <div className="card-actions justify-end">
                      <button className="btn btn-sm btn-white/20 hover:btn-white/30">
                        Počni
                      </button>
                    </div>
                  </div>
                </div>

                <div className="card bg-gradient-to-r from-green-500 to-green-600 text-white">
                  <div className="card-body">
                    <h3 className="card-title text-white">📸 Dodaj slike</h3>
                    <p className="text-white/90">
                      Učitaj visokokvalitetne fotografije proizvoda
                    </p>
                    <div className="card-actions justify-end">
                      <button className="btn btn-sm btn-white/20 hover:btn-white/30">
                        Upload
                      </button>
                    </div>
                  </div>
                </div>

                <div className="card bg-gradient-to-r from-purple-500 to-purple-600 text-white">
                  <div className="card-body">
                    <h3 className="card-title text-white">📢 Promoviši</h3>
                    <p className="text-white/90">
                      Podeli nove proizvode na društvenim mrežama
                    </p>
                    <div className="card-actions justify-end">
                      <button className="btn btn-sm btn-white/20 hover:btn-white/30">
                        Podeli
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              {/* Footer Actions */}
              <div className="flex justify-center gap-4">
                <button
                  className="btn btn-ghost"
                  onClick={() => {
                    setCurrentStep(1);
                    setSelectedMethod('');
                    setUploadedFile(null);
                    setCsvData([]);
                    setFieldMapping({});
                    setPreviewData([]);
                  }}
                >
                  🔄 Uvezi još proizvoda
                </button>
                <button className="btn btn-primary">
                  🛍️ Idi na prodavnicu
                </button>
                <button className="btn btn-secondary">
                  📊 Prikaži analitiku
                </button>
              </div>
            </div>
          </AnimatedSection>
        )}

        {/* Feature Highlights */}
        {currentStep === 1 && (
          <AnimatedSection animation="fadeIn" delay={0.6}>
            <div className="mt-16 p-8 bg-base-200 rounded-xl">
              <div className="text-center mb-8">
                <h2 className="text-2xl font-bold mb-4">
                  🚀 Mogućnosti sistema uvoza
                </h2>
                <p className="text-base-content/70">
                  Moćni alati za vlasnike velikih lančanih prodavnica
                </p>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                <div className="text-center">
                  <div className="w-16 h-16 bg-gradient-to-r from-blue-500 to-cyan-500 rounded-xl flex items-center justify-center text-3xl mx-auto mb-4">
                    🔄
                  </div>
                  <h3 className="font-bold mb-2">Autosinhronizacija</h3>
                  <p className="text-sm text-base-content/70">
                    Automatsko ažuriranje zaliha i cena u real-time
                  </p>
                </div>

                <div className="text-center">
                  <div className="w-16 h-16 bg-gradient-to-r from-green-500 to-emerald-500 rounded-xl flex items-center justify-center text-3xl mx-auto mb-4">
                    🎯
                  </div>
                  <h3 className="font-bold mb-2">Pametna validacija</h3>
                  <p className="text-sm text-base-content/70">
                    AI provera kvaliteta podataka sa PDV kalkulacijama
                  </p>
                </div>

                <div className="text-center">
                  <div className="w-16 h-16 bg-gradient-to-r from-purple-500 to-pink-500 rounded-xl flex items-center justify-center text-3xl mx-auto mb-4">
                    📊
                  </div>
                  <h3 className="font-bold mb-2">Analitika uvoza</h3>
                  <p className="text-sm text-base-content/70">
                    Detaljna statistika i izveštaji o uvozu proizvoda
                  </p>
                </div>

                <div className="text-center">
                  <div className="w-16 h-16 bg-gradient-to-r from-orange-500 to-red-500 rounded-xl flex items-center justify-center text-3xl mx-auto mb-4">
                    🛡️
                  </div>
                  <h3 className="font-bold mb-2">Rollback izmena</h3>
                  <p className="text-sm text-base-content/70">
                    Mogućnost poništavanja uvoza jednim klikom
                  </p>
                </div>
              </div>
            </div>
          </AnimatedSection>
        )}
      </div>
    </div>
  );
};

export default StorefrontImportDemo;
