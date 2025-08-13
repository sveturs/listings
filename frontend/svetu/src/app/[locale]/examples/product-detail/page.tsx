'use client';

import React, { useState, useEffect } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { motion, AnimatePresence } from 'framer-motion';
import { 
  FaArrowLeft, FaHeart, FaShare, FaShieldAlt, FaTruck, FaStore,
  FaStar, FaMapMarkerAlt, FaEye, FaCheck,
  FaShoppingCart, FaThumbsUp, FaUserCheck,
  FaBox, FaCreditCard, FaUndo, FaPhoneAlt, FaFlag,
  FaChevronRight, FaExpand, FaFire,
  FaBolt, FaPercent, FaChartLine
} from 'react-icons/fa';
import { BsChat, BsLightningFill } from 'react-icons/bs';

export default function ProductDetailPage() {
  const [selectedImage, setSelectedImage] = useState(0);
  const [activeTab, setActiveTab] = useState('description');
  const [quantity, setQuantity] = useState(1);
  const [showPhone, setShowPhone] = useState(false);
  const [isFavorite, setIsFavorite] = useState(false);
  const [showImageModal, setShowImageModal] = useState(false);
  const [selectedSize, setSelectedSize] = useState('256GB'); // Значение по умолчанию
  const [selectedColor, setSelectedColor] = useState('Natural Titanium'); // Значение по умолчанию
  const [imageError, setImageError] = useState<{[key: number]: boolean}>({});
  const [showPriceHistory, setShowPriceHistory] = useState(false);

  // Данные товара
  const product = {
    id: 1,
    title: 'iPhone 15 Pro 256GB Natural Titanium',
    price: 1299,
    oldPrice: 1499,
    discount: 13,
    rating: 4.8,
    reviews: 234,
    sold: 1856,
    views: 15234,
    favorites: 432,
    seller: {
      name: 'TechStore Official',
      rating: 4.9,
      reviews: 3421,
      registered: '2019',
      responseTime: '~15 мин',
      deliveryRate: 98.5,
      phone: '+381 69 123 4567',
      verified: true,
      official: true
    },
    location: 'Белград, Нови Београд',
    distance: '2.3 км от вас',
    images: [
      'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=800',
      'https://images.unsplash.com/photo-1695048132647-d8a0e5a2c4f0?w=800',
      'https://images.unsplash.com/photo-1695048133098-1b91c7e5d4c9?w=800',
      'https://images.unsplash.com/photo-1695048133092-50ec87e2d9e7?w=800',
      'https://images.unsplash.com/photo-1695048133086-9d79c7e1e7f5?w=800'
    ],
    characteristics: {
      'Бренд': 'Apple',
      'Модель': 'iPhone 15 Pro',
      'Память': '256 GB',
      'Цвет': 'Natural Titanium',
      'Дисплей': '6.1" Super Retina XDR',
      'Процессор': 'A17 Pro',
      'Камера': '48 MP основная',
      'Состояние': 'Новый',
      'Гарантия': '24 месяца'
    },
    sizes: ['128GB', '256GB', '512GB', '1TB'],
    colors: ['Natural Titanium', 'Blue Titanium', 'White Titanium', 'Black Titanium'],
    inStock: true,
    fastDelivery: true,
    freeDelivery: true,
    blackFriday: true,
    priceHistory: [
      { date: '01.10', price: 1499 },
      { date: '15.10', price: 1449 },
      { date: '01.11', price: 1399 },
      { date: '15.11', price: 1349 },
      { date: '25.11', price: 1299 }
    ]
  };

  // Похожие товары
  const similarProducts = [
    {
      id: 2,
      title: 'iPhone 15 128GB Blue',
      price: 999,
      oldPrice: 1099,
      image: 'https://images.unsplash.com/photo-1678685888221-cda773a3dcdb?w=400',
      rating: 4.7,
      reviews: 156
    },
    {
      id: 3,
      title: 'iPhone 14 Pro 256GB Space Black',
      price: 1099,
      oldPrice: 1299,
      image: 'https://images.unsplash.com/photo-1678911820864-e2c567c655d7?w=400',
      rating: 4.8,
      reviews: 423
    },
    {
      id: 4,
      title: 'Samsung Galaxy S24 Ultra',
      price: 1399,
      image: 'https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?w=400',
      rating: 4.9,
      reviews: 89
    },
    {
      id: 5,
      title: 'Google Pixel 8 Pro',
      price: 899,
      oldPrice: 999,
      image: 'https://images.unsplash.com/photo-1598327105666-5b89351aff97?w=400',
      rating: 4.6,
      reviews: 234
    }
  ];

  // Отзывы
  const reviews = [
    {
      id: 1,
      user: 'Александр М.',
      rating: 5,
      date: '2 дня назад',
      text: 'Отличный телефон! Доставка быстрая, все как в описании. Продавец отвечает моментально.',
      images: ['https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=200'],
      helpful: 23,
      verified: true
    },
    {
      id: 2,
      user: 'Мария К.',
      rating: 5,
      date: '1 неделю назад',
      text: 'Купила в подарок мужу, он в восторге! Телефон оригинальный, проверили по IMEI.',
      helpful: 15,
      verified: true
    },
    {
      id: 3,
      user: 'Дмитрий П.',
      rating: 4,
      date: '2 недели назад',
      text: 'Хороший телефон, но доставка задержалась на день. В остальном все отлично.',
      helpful: 8,
      verified: true
    }
  ];

  return (
    <div className="min-h-screen bg-base-200">
      {/* Хлебные крошки */}
      <div className="bg-base-100 border-b">
        <div className="container mx-auto px-4 py-3">
          <div className="flex items-center gap-2 text-sm">
            <Link href="/ru/examples/ideal-marketplace" className="hover:text-primary">
              <FaArrowLeft className="inline mr-2" />
              Главная
            </Link>
            <FaChevronRight className="text-xs opacity-50" />
            <Link href="#" className="hover:text-primary">Электроника</Link>
            <FaChevronRight className="text-xs opacity-50" />
            <Link href="#" className="hover:text-primary">Смартфоны</Link>
            <FaChevronRight className="text-xs opacity-50" />
            <span className="opacity-70">iPhone 15 Pro</span>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-6">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Левая колонка - Изображения */}
          <div className="lg:col-span-2">
            <div className="bg-base-100 rounded-xl p-6">
              {/* Основное изображение */}
              <div className="relative aspect-square mb-4 bg-base-200 rounded-lg overflow-hidden">
                {product.blackFriday && (
                  <div className="absolute top-4 left-4 z-10">
                    <div className="badge badge-error gap-1 px-3 py-4">
                      <FaPercent />
                      <span className="font-bold">Black Friday</span>
                    </div>
                  </div>
                )}
                <div className="absolute top-4 right-4 z-10 flex flex-col gap-2">
                  <button 
                    onClick={() => setIsFavorite(!isFavorite)}
                    className={`btn btn-circle btn-sm ${isFavorite ? 'btn-error' : 'btn-ghost bg-base-100'}`}
                  >
                    <FaHeart className={isFavorite ? 'text-white' : ''} />
                  </button>
                  <button className="btn btn-circle btn-sm btn-ghost bg-base-100">
                    <FaShare />
                  </button>
                  <button 
                    onClick={() => setShowImageModal(true)}
                    className="btn btn-circle btn-sm btn-ghost bg-base-100"
                  >
                    <FaExpand />
                  </button>
                </div>
                {!imageError[selectedImage] ? (
                  <Image
                    src={product.images[selectedImage]}
                    alt={product.title}
                    width={800}
                    height={800}
                    className="object-contain w-full h-full cursor-zoom-in"
                    onClick={() => setShowImageModal(true)}
                    onError={() => setImageError(prev => ({...prev, [selectedImage]: true}))}
                  />
                ) : (
                  <div className="flex items-center justify-center w-full h-full bg-base-300">
                    <div className="text-center">
                      <FaBox className="text-4xl opacity-30 mb-2" />
                      <p className="text-sm opacity-50">Изображение недоступно</p>
                    </div>
                  </div>
                )}
              </div>

              {/* Миниатюры */}
              <div className="flex gap-2 overflow-x-auto pb-2">
                {product.images.map((img, idx) => (
                  <button
                    key={idx}
                    onClick={() => setSelectedImage(idx)}
                    className={`flex-shrink-0 w-20 h-20 rounded-lg overflow-hidden border-2 transition-all ${
                      selectedImage === idx ? 'border-primary ring-2 ring-primary/20' : 'border-base-300'
                    }`}
                  >
                    {!imageError[idx] ? (
                      <Image
                        src={img}
                        alt={`${product.title} ${idx + 1}`}
                        width={80}
                        height={80}
                        className="object-cover w-full h-full"
                        onError={() => setImageError(prev => ({...prev, [idx]: true}))}
                      />
                    ) : (
                      <div className="flex items-center justify-center w-full h-full bg-base-300">
                        <FaBox className="text-xs opacity-30" />
                      </div>
                    )}
                  </button>
                ))}
              </div>

              {/* Табы с информацией */}
              <div className="mt-8">
                <div className="tabs tabs-boxed bg-base-200">
                  <button 
                    className={`tab ${activeTab === 'description' ? 'tab-active' : ''}`}
                    onClick={() => setActiveTab('description')}
                  >
                    Описание
                  </button>
                  <button 
                    className={`tab ${activeTab === 'characteristics' ? 'tab-active' : ''}`}
                    onClick={() => setActiveTab('characteristics')}
                  >
                    Характеристики
                  </button>
                  <button 
                    className={`tab ${activeTab === 'reviews' ? 'tab-active' : ''}`}
                    onClick={() => setActiveTab('reviews')}
                  >
                    Отзывы ({product.reviews})
                  </button>
                  <button 
                    className={`tab ${activeTab === 'delivery' ? 'tab-active' : ''}`}
                    onClick={() => setActiveTab('delivery')}
                  >
                    Доставка
                  </button>
                </div>

                <div className="mt-6">
                  {activeTab === 'description' && (
                    <div className="prose max-w-none">
                      <h3>О товаре</h3>
                      <p>
                        iPhone 15 Pro - это вершина инженерной мысли Apple. Корпус из титана аэрокосмического класса 
                        делает его невероятно прочным и легким одновременно. Новый чип A17 Pro обеспечивает 
                        беспрецедентную производительность.
                      </p>
                      <h4>Основные преимущества:</h4>
                      <ul>
                        <li>Титановый корпус - прочность и легкость</li>
                        <li>Чип A17 Pro с 3-нм техпроцессом</li>
                        <li>Камера 48 МП с оптическим зумом</li>
                        <li>Action Button для быстрого доступа</li>
                        <li>USB-C с поддержкой USB 3</li>
                      </ul>
                      <p>
                        Телефон новый, запечатанный, с официальной гарантией Apple 24 месяца. 
                        В комплекте все оригинальные аксессуары.
                      </p>
                    </div>
                  )}

                  {activeTab === 'characteristics' && (
                    <div className="space-y-3">
                      {Object.entries(product.characteristics).map(([key, value]) => (
                        <div key={key} className="flex justify-between py-2 border-b border-base-300">
                          <span className="text-base-content/70">{key}</span>
                          <span className="font-medium">{value}</span>
                        </div>
                      ))}
                    </div>
                  )}

                  {activeTab === 'reviews' && (
                    <div className="space-y-4">
                      {/* Общая статистика */}
                      <div className="bg-base-200 rounded-lg p-4">
                        <div className="flex items-center gap-4">
                          <div className="text-center">
                            <div className="text-3xl font-bold">{product.rating}</div>
                            <div className="rating rating-sm">
                              {[1,2,3,4,5].map(i => (
                                <input 
                                  key={i} 
                                  type="radio" 
                                  name="product-rating"
                                  className="mask mask-star-2 bg-warning" 
                                  checked={i === Math.round(product.rating)}
                                  readOnly
                                />
                              ))}
                            </div>
                            <div className="text-sm opacity-70">{product.reviews} отзывов</div>
                          </div>
                          <div className="flex-1">
                            <div className="space-y-1">
                              {[5,4,3,2,1].map(stars => (
                                <div key={stars} className="flex items-center gap-2">
                                  <span className="text-xs w-3">{stars}</span>
                                  <FaStar className="text-warning text-xs" />
                                  <progress 
                                    className="progress progress-warning w-full" 
                                    value={stars === 5 ? 75 : stars === 4 ? 20 : 5} 
                                    max="100"
                                  />
                                </div>
                              ))}
                            </div>
                          </div>
                        </div>
                      </div>

                      {/* Список отзывов */}
                      {reviews.map(review => (
                        <div key={review.id} className="border-b border-base-300 pb-4">
                          <div className="flex items-start gap-3">
                            <div className="avatar placeholder">
                              <div className="bg-primary text-primary-content rounded-full w-10">
                                <span>{review.user[0]}</span>
                              </div>
                            </div>
                            <div className="flex-1">
                              <div className="flex items-center gap-2 mb-1">
                                <span className="font-medium">{review.user}</span>
                                {review.verified && (
                                  <span className="badge badge-success badge-sm gap-1">
                                    <FaCheck className="text-xs" />
                                    Проверенный покупатель
                                  </span>
                                )}
                              </div>
                              <div className="rating rating-sm mb-2">
                                {[1,2,3,4,5].map(i => (
                                  <input 
                                    key={i} 
                                    type="radio" 
                                    name={`review-rating-${review.id}`}
                                    className="mask mask-star-2 bg-warning" 
                                    checked={i <= review.rating}
                                    readOnly
                                  />
                                ))}
                                <span className="text-sm ml-2 opacity-70">{review.date}</span>
                              </div>
                              <p className="text-sm">{review.text}</p>
                              {review.images && (
                                <div className="flex gap-2 mt-2">
                                  {review.images.map((img, idx) => (
                                    <Image
                                      key={idx}
                                      src={img}
                                      alt="Review"
                                      width={60}
                                      height={60}
                                      className="rounded cursor-pointer hover:opacity-80"
                                    />
                                  ))}
                                </div>
                              )}
                              <div className="flex items-center gap-4 mt-3">
                                <button className="btn btn-ghost btn-xs gap-1">
                                  <FaThumbsUp />
                                  Полезно ({review.helpful})
                                </button>
                                <button className="btn btn-ghost btn-xs">
                                  Ответить
                                </button>
                              </div>
                            </div>
                          </div>
                        </div>
                      ))}
                      <button className="btn btn-outline btn-block">
                        Показать все отзывы
                      </button>
                    </div>
                  )}

                  {activeTab === 'delivery' && (
                    <div className="space-y-4">
                      <div className="alert alert-success">
                        <FaTruck />
                        <span>Бесплатная доставка при заказе от €50</span>
                      </div>
                      <div className="space-y-3">
                        <div className="flex items-start gap-3">
                          <FaBolt className="text-warning mt-1" />
                          <div>
                            <div className="font-medium">Экспресс-доставка</div>
                            <div className="text-sm opacity-70">Сегодня до 18:00 - €5</div>
                          </div>
                        </div>
                        <div className="flex items-start gap-3">
                          <FaTruck className="text-primary mt-1" />
                          <div>
                            <div className="font-medium">Стандартная доставка</div>
                            <div className="text-sm opacity-70">1-2 рабочих дня - €3</div>
                          </div>
                        </div>
                        <div className="flex items-start gap-3">
                          <FaStore className="text-info mt-1" />
                          <div>
                            <div className="font-medium">Самовывоз</div>
                            <div className="text-sm opacity-70">Белград, ТЦ "Delta City" - бесплатно</div>
                          </div>
                        </div>
                      </div>
                      <div className="divider"></div>
                      <div className="space-y-2">
                        <div className="flex items-center gap-2">
                          <FaUndo className="text-success" />
                          <span>Возврат товара в течение 14 дней</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <FaShieldAlt className="text-info" />
                          <span>Защита покупателя SveTu</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <FaCreditCard className="text-primary" />
                          <span>Оплата при получении</span>
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>

          {/* Правая колонка - Информация о покупке */}
          <div className="space-y-4">
            {/* Цена и действия */}
            <div className="bg-base-100 rounded-xl p-6 sticky top-4">
              {/* История цены */}
              {product.blackFriday && (
                <div className="mb-4">
                  <button 
                    onClick={() => setShowPriceHistory(!showPriceHistory)}
                    className="flex items-center gap-2 text-sm text-primary hover:underline"
                  >
                    <FaChartLine />
                    История цены
                    <span className="badge badge-success badge-sm">-{product.discount}%</span>
                  </button>
                  {showPriceHistory && (
                    <div className="mt-2 p-3 bg-base-200 rounded-lg">
                      <div className="text-xs space-y-1">
                        {product.priceHistory.map((item, idx) => (
                          <div key={idx} className="flex justify-between">
                            <span className="opacity-70">{item.date}</span>
                            <span className={idx === product.priceHistory.length - 1 ? 'text-success font-bold' : ''}>
                              €{item.price}
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              )}

              {/* Цена */}
              <div className="mb-4">
                <div className="flex items-center gap-2 mb-1">
                  {product.oldPrice && (
                    <>
                      <span className="text-2xl line-through opacity-50">€{product.oldPrice}</span>
                      <span className="badge badge-error">-{product.discount}%</span>
                    </>
                  )}
                </div>
                <div className="text-4xl font-bold text-primary">€{product.price}</div>
                {product.blackFriday && (
                  <div className="flex items-center gap-1 mt-2 text-sm text-error">
                    <FaFire />
                    <span>Цена по Black Friday!</span>
                  </div>
                )}
              </div>

              {/* Выбор параметров */}
              {product.colors && (
                <div className="mb-4">
                  <label className="text-sm font-medium mb-2 block">Цвет:</label>
                  <div className="flex flex-wrap gap-2">
                    {product.colors.map(color => (
                      <button
                        key={color}
                        onClick={() => setSelectedColor(color)}
                        className={`btn btn-sm ${selectedColor === color ? 'btn-primary' : 'btn-outline'}`}
                      >
                        {color}
                      </button>
                    ))}
                  </div>
                </div>
              )}

              {product.sizes && (
                <div className="mb-4">
                  <label className="text-sm font-medium mb-2 block">Память:</label>
                  <div className="flex flex-wrap gap-2">
                    {product.sizes.map(size => (
                      <button
                        key={size}
                        onClick={() => setSelectedSize(size)}
                        className={`btn btn-sm ${selectedSize === size ? 'btn-primary' : 'btn-outline'}`}
                      >
                        {size}
                      </button>
                    ))}
                  </div>
                </div>
              )}

              {/* Количество */}
              <div className="mb-4">
                <label className="text-sm font-medium mb-2 block">Количество:</label>
                <div className="join">
                  <button 
                    className="btn join-item"
                    onClick={() => setQuantity(Math.max(1, quantity - 1))}
                  >
                    -
                  </button>
                  <input 
                    type="number" 
                    className="input input-bordered join-item w-20 text-center" 
                    value={quantity}
                    onChange={(e) => setQuantity(parseInt(e.target.value) || 1)}
                  />
                  <button 
                    className="btn join-item"
                    onClick={() => setQuantity(quantity + 1)}
                  >
                    +
                  </button>
                </div>
                {product.inStock && (
                  <div className="text-sm text-success mt-1">✓ В наличии</div>
                )}
              </div>

              {/* Кнопки действий */}
              <div className="space-y-2">
                <button 
                  className="btn btn-primary btn-block"
                  onClick={() => {
                    if (product.sizes && !selectedSize) {
                      alert('Пожалуйста, выберите размер памяти');
                      return;
                    }
                    if (product.colors && !selectedColor) {
                      alert('Пожалуйста, выберите цвет');
                      return;
                    }
                    alert(`Добавлено в корзину: ${product.title}${selectedSize ? ' ' + selectedSize : ''}${selectedColor ? ' ' + selectedColor : ''} x${quantity}`);
                  }}
                >
                  <FaShoppingCart />
                  Добавить в корзину
                </button>
                <button className="btn btn-secondary btn-block">
                  Купить сейчас
                </button>
                <button 
                  className="btn btn-outline btn-block"
                  onClick={() => setShowPhone(!showPhone)}
                >
                  <FaPhoneAlt />
                  {showPhone ? product.seller.phone : 'Показать телефон'}
                </button>
                <button className="btn btn-outline btn-block">
                  <BsChat />
                  Написать продавцу
                </button>
              </div>

              {/* Преимущества */}
              <div className="mt-6 space-y-2">
                {product.fastDelivery && (
                  <div className="flex items-center gap-2 text-sm">
                    <BsLightningFill className="text-warning" />
                    <span>Быстрая доставка</span>
                  </div>
                )}
                {product.freeDelivery && (
                  <div className="flex items-center gap-2 text-sm">
                    <FaTruck className="text-success" />
                    <span>Бесплатная доставка от €50</span>
                  </div>
                )}
                <div className="flex items-center gap-2 text-sm">
                  <FaShieldAlt className="text-info" />
                  <span>Защита покупателя</span>
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <FaUndo className="text-primary" />
                  <span>Возврат 14 дней</span>
                </div>
              </div>
            </div>

            {/* Информация о продавце */}
            <div className="bg-base-100 rounded-xl p-6">
              <div className="flex items-start gap-3 mb-4">
                <div className="avatar placeholder">
                  <div className="bg-primary text-primary-content rounded-full w-12">
                    <span>TS</span>
                  </div>
                </div>
                <div className="flex-1">
                  <div className="flex items-center gap-2">
                    <h3 className="font-bold">{product.seller.name}</h3>
                    {product.seller.official && (
                      <span className="badge badge-primary badge-sm">Official</span>
                    )}
                    {product.seller.verified && (
                      <FaUserCheck className="text-success" />
                    )}
                  </div>
                  <div className="flex items-center gap-1 text-sm">
                    <div className="rating rating-xs">
                      {[1,2,3,4,5].map(i => (
                        <input 
                          key={i} 
                          type="radio" 
                          name="seller-rating"
                          className="mask mask-star-2 bg-warning" 
                          checked={i === Math.round(product.seller.rating)}
                          readOnly
                        />
                      ))}
                    </div>
                    <span>{product.seller.rating}</span>
                    <span className="opacity-70">({product.seller.reviews})</span>
                  </div>
                </div>
              </div>

              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="opacity-70">На площадке с</span>
                  <span>{product.seller.registered}</span>
                </div>
                <div className="flex justify-between">
                  <span className="opacity-70">Время ответа</span>
                  <span>{product.seller.responseTime}</span>
                </div>
                <div className="flex justify-between">
                  <span className="opacity-70">Доставка вовремя</span>
                  <span className="text-success">{product.seller.deliveryRate}%</span>
                </div>
              </div>

              <div className="mt-4 pt-4 border-t border-base-300">
                <Link href="#" className="btn btn-sm btn-outline btn-block">
                  Все товары продавца
                </Link>
              </div>
            </div>

            {/* Местоположение */}
            <div className="bg-base-100 rounded-xl p-6">
              <h3 className="font-bold mb-3 flex items-center gap-2">
                <FaMapMarkerAlt className="text-primary" />
                Местоположение
              </h3>
              <p className="text-sm mb-2">{product.location}</p>
              <p className="text-sm text-primary">{product.distance}</p>
              <div className="mt-3 h-32 bg-base-200 rounded-lg flex items-center justify-center">
                <span className="text-sm opacity-50">Карта</span>
              </div>
            </div>

            {/* Статистика */}
            <div className="bg-base-100 rounded-xl p-6">
              <div className="grid grid-cols-2 gap-4 text-center">
                <div>
                  <div className="text-2xl font-bold text-primary">{product.views.toLocaleString()}</div>
                  <div className="text-xs opacity-70">
                    <FaEye className="inline mr-1" />
                    Просмотров
                  </div>
                </div>
                <div>
                  <div className="text-2xl font-bold text-error">{product.favorites}</div>
                  <div className="text-xs opacity-70">
                    <FaHeart className="inline mr-1" />
                    В избранном
                  </div>
                </div>
                <div>
                  <div className="text-2xl font-bold text-success">{product.sold.toLocaleString()}</div>
                  <div className="text-xs opacity-70">
                    <FaBox className="inline mr-1" />
                    Продано
                  </div>
                </div>
                <div>
                  <div className="text-2xl font-bold text-warning">{product.rating}</div>
                  <div className="text-xs opacity-70">
                    <FaStar className="inline mr-1" />
                    Рейтинг
                  </div>
                </div>
              </div>
            </div>

            {/* Пожаловаться */}
            <button className="btn btn-ghost btn-sm btn-block">
              <FaFlag />
              Пожаловаться
            </button>
          </div>
        </div>

        {/* Похожие товары */}
        <div className="mt-12">
          <h2 className="text-2xl font-bold mb-6">Похожие товары</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-4">
            {similarProducts.map(item => (
              <div key={item.id} className="card bg-base-100 hover:shadow-xl transition-shadow">
                <figure className="relative h-48">
                  {item.oldPrice && (
                    <div className="badge badge-error absolute top-2 right-2 z-10">
                      -{Math.round((1 - item.price/item.oldPrice) * 100)}%
                    </div>
                  )}
                  <Image
                    src={item.image}
                    alt={item.title}
                    width={200}
                    height={200}
                    className="object-cover w-full h-full"
                  />
                </figure>
                <div className="card-body p-4">
                  <h3 className="text-sm line-clamp-2">{item.title}</h3>
                  <div className="flex items-center gap-1 text-xs">
                    <div className="rating rating-xs">
                      {[1,2,3,4,5].map(i => (
                        <input 
                          key={i} 
                          type="radio" 
                          name={`similar-rating-${item.id}`}
                          className="mask mask-star-2 bg-warning" 
                          checked={i === Math.round(item.rating)}
                          readOnly
                        />
                      ))}
                    </div>
                    <span>({item.reviews})</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-lg font-bold">€{item.price}</span>
                    {item.oldPrice && (
                      <span className="text-sm line-through opacity-50">€{item.oldPrice}</span>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Недавно просмотренные */}
        <div className="mt-12">
          <h2 className="text-2xl font-bold mb-6">Вы недавно смотрели</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-4">
            {similarProducts.slice(0, 3).map(item => (
              <div key={item.id} className="card bg-base-100 hover:shadow-xl transition-shadow">
                <figure className="relative h-48">
                  <Image
                    src={item.image}
                    alt={item.title}
                    width={200}
                    height={200}
                    className="object-cover w-full h-full"
                  />
                </figure>
                <div className="card-body p-4">
                  <h3 className="text-sm line-clamp-2">{item.title}</h3>
                  <div className="text-lg font-bold">€{item.price}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Модальное окно для просмотра изображений */}
      <AnimatePresence>
        {showImageModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-50 bg-black/90 flex items-center justify-center p-4"
            onClick={() => setShowImageModal(false)}
          >
            <button 
              className="absolute top-4 right-4 btn btn-circle btn-ghost text-white"
              onClick={() => setShowImageModal(false)}
            >
              ✕
            </button>
            <div className="relative max-w-5xl max-h-[90vh]">
              <Image
                src={product.images[selectedImage]}
                alt={product.title}
                width={1200}
                height={1200}
                className="object-contain max-h-[90vh] w-auto"
              />
              <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2">
                {product.images.map((_, idx) => (
                  <button
                    key={idx}
                    onClick={(e) => {
                      e.stopPropagation();
                      setSelectedImage(idx);
                    }}
                    className={`w-2 h-2 rounded-full transition-all ${
                      selectedImage === idx ? 'bg-white w-8' : 'bg-white/50'
                    }`}
                  />
                ))}
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}