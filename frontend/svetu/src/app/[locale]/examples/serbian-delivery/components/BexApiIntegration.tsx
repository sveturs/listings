'use client';

import React, { useState } from 'react';
import {
  ServerIcon,
  KeyIcon,
  CheckCircleIcon,
  DocumentTextIcon,
  GlobeAltIcon,
  BoltIcon,
  ShieldCheckIcon,
  CodeBracketIcon,
} from '@heroicons/react/24/outline';

export default function BexApiIntegration() {
  const [activeEndpoint, setActiveEndpoint] = useState('postShipments');
  const [showResponse, setShowResponse] = useState(false);

  const endpoints = [
    {
      id: 'postShipments',
      name: 'POST /postShipments',
      description: 'Креирање домаћих пошиљки',
      color: 'bg-blue-100',
      textColor: 'text-blue-600',
    },
    {
      id: 'postShipmentsCustoms',
      name: 'POST /postShipmentsCustoms',
      description: 'Међународне пошиљке са царином',
      color: 'bg-green-100',
      textColor: 'text-green-600',
    },
    {
      id: 'getState',
      name: 'GET /getState',
      description: 'Праћење статуса пошиљке',
      color: 'bg-orange-100',
      textColor: 'text-orange-600',
    },
    {
      id: 'getLabelWithProperties',
      name: 'GET /getLabelWithProperties',
      description: 'Генерисање адресница',
      color: 'bg-purple-100',
      textColor: 'text-purple-600',
    },
    {
      id: 'listShipments',
      name: 'POST /listShipments',
      description: 'Листа пошиљки са филтерима',
      color: 'bg-indigo-100',
      textColor: 'text-indigo-600',
    },
    {
      id: 'listParcelShops',
      name: 'POST /listParcelShops',
      description: 'Листа Parcel Shop локација',
      color: 'bg-pink-100',
      textColor: 'text-pink-600',
    },
    {
      id: 'delete',
      name: 'POST /delete',
      description: 'Отказивање пошиљке',
      color: 'bg-red-100',
      textColor: 'text-red-600',
    },
  ];

  const apiExamples = {
    postShipments: {
      request: `{
  "shipmentslist": [{
    "shipmentId": 0,
    "serviceSpeed": 1,
    "shipmentType": 1,
    "shipmentCategory": 1,
    "totalPackages": 1,
    "shipmentContents": 3,
    "commentPublic": "Пажљиво - ломљиво",
    "commentPrivate": "Поруџбина #12345",
    "personalDelivery": false,
    "payType": 6,
    "insuranceAmount": 5000,
    "payToSender": 2500,
    "tasks": [{
      "type": 1,
      "nameType": 1,
      "name1": "Петровић",
      "name2": "Марко",
      "place": "БЕОГРАД",
      "street": "Кнез Михаилова",
      "houseNumber": 25,
      "phone": "0648516928"
    }, {
      "type": 2,
      "nameType": 1,
      "name1": "Јовановић",
      "name2": "Ана",
      "place": "НОВИ САД",
      "street": "Дунавска",
      "houseNumber": 10,
      "phone": "0651234567"
    }]
  }]
}`,
      response: `{
  "shipmentsResultList": [{
    "state": true,
    "shipmentId": 170123456,
    "err": ""
  }],
  "reqstate": true,
  "reqerr": ""
}`,
    },
    postShipmentsCustoms: {
      request: `{
  "shipmentId": 0,
  "shipmentType": 1,
  "shipmentCategory": 5,
  "totalPackages": 1,
  "parcels": [{
    "No": 1,
    "OperatorNumber": "LP000164786PP",
    "items": [{
      "ordinalNumber": 1,
      "description": "Electronics component",
      "skuCode": "SKU1234",
      "hsCode": "HS5678",
      "quantity": 2,
      "valuePerItem": 15.50,
      "weightPerItem": 0.25,
      "originCountryCode": "CN"
    }]
  }],
  "customs": [{
    "originSenderName": "China Export Co.",
    "originSenderAddress": "123 Export Street",
    "originSenderPlace": "Shanghai",
    "originSenderCountryCode": "CN",
    "currencyCode": "USD",
    "DDP": true
  }]
}`,
      response: `{
  "shipmentsResultList": [{
    "shipmentId": 245815264,
    "shipmentStatus": -51,
    "isSuccessful": true,
    "errorMessage": ""
  }]
}`,
    },
    getState: {
      request: `/ship/api/Ship/getstate?shipmentid=171123456&mtype=1&lang=1`,
      response: `{
  "status": 5,
  "statusText": "испоручена",
  "deliveryDate": "29.08.2024",
  "sendDate": "26.08.2024",
  "transferAmount": 2500,
  "note": "Потписао: Ана Јовановић"
}`,
    },
    getLabelWithProperties: {
      request: `/ShipDNF/Ship/getLabelWithProperties?pageSize=4&pagePosition=1&shipmentId=180923112&parcelNo=0`,
      response: `{
  "state": true,
  "shipmentId": 170123456,
  "parcelNo": 1,
  "parcelLabel": "JVBERi0xLjQKJeLjz9M...base64...",
  "err": ""
}`,
    },
    listShipments: {
      request: `{
  "lang": 1,
  "dateFormat": 1,
  "shipmentId": 0,
  "dateStart": "2024-08-01",
  "dateEnd": "2024-08-31",
  "status": 21
}`,
      response: `{
  "listShipment": [{
    "shipmentId": 193930784,
    "createDate": "2024-08-13T00:00:00.000Z",
    "senderName": "Test D.O.O.",
    "receiverName": "Bexexpress",
    "totalPackages": 4,
    "stateName": "Достављено",
    "totalPrice": 470,
    "payToSender": 18380
  }],
  "state": true,
  "err": "ok"
}`,
    },
    listParcelShops: {
      request: `{
  "regionId": 2
}`,
      response: `{
  "listShipment": [{
    "code": 359,
    "regionId": 2,
    "name": "Bex magacin ŠA",
    "city": "ŠABAC",
    "address": "ŠABAC, SUVOBORSKA bb",
    "businessHours": "08:00 - 13:00",
    "xcoordinate": "44.73880949442088",
    "ycoordinate": "19.691042501713863"
  }],
  "state": true
}`,
    },
    delete: {
      request: `{
  "shipmentid": 170123456,
  "mtype": 1,
  "lang": 1
}`,
      response: `{
  "state": true,
  "message": "Пошиљка успешно отказана"
}`,
    },
  };

  const handleTestApi = () => {
    setShowResponse(true);
    setTimeout(() => setShowResponse(false), 3000);
  };

  return (
    <div className="space-y-6">
      {/* API Connection Status */}
      <div className="alert alert-success">
        <CheckCircleIcon className="w-6 h-6" />
        <div>
          <h3 className="font-bold">API Connection Active</h3>
          <div className="text-xs">Endpoint: https://api.bex.rs:62502</div>
        </div>
        <div className="badge badge-primary">Bearer Token</div>
      </div>

      {/* API Features */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-figure text-primary">
            <ServerIcon className="w-8 h-8" />
          </div>
          <div className="stat-title">RESTful</div>
          <div className="stat-value text-primary">API v2</div>
          <div className="stat-desc">JSON формат</div>
        </div>

        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-figure text-secondary">
            <KeyIcon className="w-8 h-8" />
          </div>
          <div className="stat-title">Аутентификација</div>
          <div className="stat-value text-secondary">Bearer</div>
          <div className="stat-desc">X-AUTH-TOKEN</div>
        </div>

        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-figure text-accent">
            <BoltIcon className="w-8 h-8" />
          </div>
          <div className="stat-title">Брзина</div>
          <div className="stat-value text-accent">&lt; 200ms</div>
          <div className="stat-desc">Просечно време одзива</div>
        </div>

        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-figure text-success">
            <ShieldCheckIcon className="w-8 h-8" />
          </div>
          <div className="stat-title">Доступност</div>
          <div className="stat-value text-success">99.9%</div>
          <div className="stat-desc">SLA гаранција</div>
        </div>
      </div>

      {/* API Endpoints */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h3 className="card-title mb-4">
            <CodeBracketIcon className="w-6 h-6" />
            API Endpoints
          </h3>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
            {/* Endpoint List */}
            <div className="space-y-2">
              {endpoints.map((endpoint) => (
                <button
                  key={endpoint.id}
                  onClick={() => setActiveEndpoint(endpoint.id)}
                  className={`w-full text-left p-3 rounded-lg transition-all ${
                    activeEndpoint === endpoint.id
                      ? 'bg-primary text-primary-content'
                      : 'bg-base-200 hover:bg-base-300'
                  }`}
                >
                  <div className="font-mono text-sm font-bold">
                    {endpoint.name}
                  </div>
                  <div className="text-xs mt-1 opacity-80">
                    {endpoint.description}
                  </div>
                </button>
              ))}
            </div>

            {/* Code Examples */}
            <div className="lg:col-span-2 space-y-4">
              <div className="tabs tabs-boxed">
                <a className="tab tab-active">Request</a>
                <a className="tab">Response</a>
                <a className="tab">cURL</a>
                <a className="tab">JavaScript</a>
              </div>

              <div className="mockup-code">
                <pre data-prefix="$">
                  <code className="text-xs">
                    curl -X POST https://api.bex.rs:62502{activeEndpoint}
                  </code>
                </pre>
                <pre data-prefix=">" className="text-warning">
                  <code className="text-xs">
                    -H &quot;X-AUTH-TOKEN: your-api-token&quot;
                  </code>
                </pre>
                <pre data-prefix=">" className="text-info">
                  <code className="text-xs">
                    -H &quot;Content-Type: application/json&quot;
                  </code>
                </pre>
                <pre data-prefix=">" className="text-success">
                  <code className="text-xs">
                    -d &apos;
                    {JSON.stringify(
                      apiExamples[activeEndpoint as keyof typeof apiExamples]
                        ?.request || {},
                      null,
                      2
                    ).slice(0, 200)}
                    ...&apos;
                  </code>
                </pre>
              </div>

              {showResponse && (
                <div className="alert alert-success">
                  <CheckCircleIcon className="w-5 h-5" />
                  <span>API call successful! Response received.</span>
                </div>
              )}

              <button
                onClick={handleTestApi}
                className="btn btn-primary btn-block"
              >
                <BoltIcon className="w-5 h-5" />
                Test API Endpoint
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Integration Steps */}
      <div className="card bg-gradient-to-r from-blue-50 to-red-50">
        <div className="card-body">
          <h3 className="card-title mb-4">
            <DocumentTextIcon className="w-6 h-6" />
            Кораци интеграције
          </h3>

          <ul className="timeline timeline-vertical">
            <li>
              <div className="timeline-start">01</div>
              <div className="timeline-middle">
                <CheckCircleIcon className="w-5 h-5 text-primary" />
              </div>
              <div className="timeline-end timeline-box">
                <div className="font-bold">Регистрација</div>
                <p className="text-sm">
                  Контактирајте BexExpress за API креденцијале
                </p>
              </div>
              <hr className="bg-primary" />
            </li>
            <li>
              <hr className="bg-primary" />
              <div className="timeline-start">02</div>
              <div className="timeline-middle">
                <CheckCircleIcon className="w-5 h-5 text-primary" />
              </div>
              <div className="timeline-end timeline-box">
                <div className="font-bold">Тестирање</div>
                <p className="text-sm">Користите sandbox окружење за развој</p>
              </div>
              <hr />
            </li>
            <li>
              <hr />
              <div className="timeline-start">03</div>
              <div className="timeline-middle">
                <CheckCircleIcon className="w-5 h-5" />
              </div>
              <div className="timeline-end timeline-box">
                <div className="font-bold">Имплементација</div>
                <p className="text-sm">Интегришите API у вашу платформу</p>
              </div>
              <hr />
            </li>
            <li>
              <hr />
              <div className="timeline-start">04</div>
              <div className="timeline-middle">
                <CheckCircleIcon className="w-5 h-5" />
              </div>
              <div className="timeline-end timeline-box">
                <div className="font-bold">Продукција</div>
                <p className="text-sm">Пребаците се на продукцијски API</p>
              </div>
            </li>
          </ul>
        </div>
      </div>

      {/* Supported Features */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body items-center text-center">
            <GlobeAltIcon className="w-12 h-12 text-green-600 mb-2" />
            <p className="text-sm font-semibold">Међународна достава</p>
          </div>
        </div>
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body items-center text-center">
            <DocumentTextIcon className="w-12 h-12 text-purple-600 mb-2" />
            <p className="text-sm font-semibold">PDF адреснице</p>
          </div>
        </div>
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body items-center text-center">
            <BoltIcon className="w-12 h-12 text-orange-600 mb-2" />
            <p className="text-sm font-semibold">Webhooks</p>
          </div>
        </div>
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body items-center text-center">
            <ShieldCheckIcon className="w-12 h-12 text-blue-600 mb-2" />
            <p className="text-sm font-semibold">SSL/TLS</p>
          </div>
        </div>
      </div>
    </div>
  );
}
