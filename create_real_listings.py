#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import requests
import json
import time
import random
from datetime import datetime
import base64

# Configuration
API_BASE = "http://localhost:3000/api/v1"
JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3R1c2VyN0BleGFtcGxlLmNvbSIsImV4cCI6MTc1NjQyMTgxOSwiaWF0IjoxNzU2MzM1NDE5LCJpc19hZG1pbiI6ZmFsc2UsInVzZXJfaWQiOjd9.fFArmmTGPugYTVEkdAZMvlDrVqVqqgjRsxwHfRqzZmc"

headers = {
    "Authorization": f"Bearer {JWT_TOKEN}",
    "Content-Type": "application/json"
}

# Real estate listings in Novi Sad
real_estate_novi_sad = [
    {
        "title": {
            "ru": "Современная квартира в центре Нови Сада",
            "en": "Modern apartment in Novi Sad center",
            "sr": "Moderan stan u centru Novog Sada"
        },
        "description": {
            "ru": "Просторная двухкомнатная квартира с современной отделкой в самом центре города. Полностью меблирована, с видом на Дунай. Рядом парк, магазины, рестораны.",
            "en": "Spacious two-bedroom apartment with modern finishing in the city center. Fully furnished, with Danube river view. Near park, shops, restaurants.",
            "sr": "Prostran dvosoban stan sa modernim završnim radovima u samom centru grada. Potpuno namešten, sa pogledom na Dunav. Blizu parka, prodavnica, restorana."
        },
        "price": 850,
        "currency": "EUR",
        "category_id": 146,  # Real estate -> Apartments
        "location": {
            "latitude": 45.2551,
            "longitude": 19.8451,
            "address": "Булевар Михаила Пупина 10, Нови Сад",
            "city": "Novi Sad",
            "country": "Serbia"
        },
        "attributes": {
            "rooms": "2",
            "area": "65",
            "floor": "3",
            "total_floors": "5",
            "furnished": "true",
            "heating": "central"
        },
        "images": [
            "https://images.unsplash.com/photo-1502672260266-1c1ef2d93688",
            "https://images.unsplash.com/photo-1560448204-e02f11c3d0e2",
            "https://images.unsplash.com/photo-1558442086-8ea19a79cd4d"
        ]
    },
    {
        "title": {
            "ru": "Роскошный пентхаус с террасой",
            "en": "Luxury penthouse with terrace",
            "sr": "Luksuzni penthaus sa terasom"
        },
        "description": {
            "ru": "Эксклюзивный пентхаус площадью 120м² с панорамной террасой 40м². Дизайнерский ремонт, умный дом, 3 спальни, 2 ванные комнаты. Подземный гараж на 2 машины.",
            "en": "Exclusive 120m² penthouse with 40m² panoramic terrace. Designer renovation, smart home, 3 bedrooms, 2 bathrooms. Underground garage for 2 cars.",
            "sr": "Ekskluzivni penthaus od 120m² sa panoramskom terasom od 40m². Dizajnerska renovacija, pametan dom, 3 spavaće sobe, 2 kupatila. Podzemna garaža za 2 automobila."
        },
        "price": 2200,
        "currency": "EUR",
        "category_id": 146,
        "location": {
            "latitude": 45.2467,
            "longitude": 19.8515,
            "address": "Лимански парк, Нови Сад",
            "city": "Novi Sad",
            "country": "Serbia"
        },
        "attributes": {
            "rooms": "4",
            "area": "120",
            "floor": "8",
            "total_floors": "8",
            "furnished": "true",
            "heating": "floor",
            "parking": "true"
        },
        "images": [
            "https://images.unsplash.com/photo-1512917774080-9991f1c4c750",
            "https://images.unsplash.com/photo-1416331108676-a22ccb276e35",
            "https://images.unsplash.com/photo-1484154218962-a197022b5858"
        ]
    },
    {
        "title": {
            "ru": "Уютная студия возле университета",
            "en": "Cozy studio near university",
            "sr": "Udoban studio blizu univerziteta"
        },
        "description": {
            "ru": "Идеальная студия для студентов или молодых специалистов. Полностью оборудована, современная кухня, высокоскоростной интернет. В 5 минутах от университета.",
            "en": "Perfect studio for students or young professionals. Fully equipped, modern kitchen, high-speed internet. 5 minutes from university.",
            "sr": "Savršen studio za studente ili mlade profesionalce. Potpuno opremljen, moderna kuhinja, brzi internet. 5 minuta od univerziteta."
        },
        "price": 350,
        "currency": "EUR",
        "category_id": 146,
        "location": {
            "latitude": 45.2485,
            "longitude": 19.8335,
            "address": "Улица Данила Киша 15, Нови Сад",
            "city": "Novi Sad",
            "country": "Serbia"
        },
        "attributes": {
            "rooms": "1",
            "area": "32",
            "floor": "2",
            "total_floors": "4",
            "furnished": "true",
            "heating": "electric"
        },
        "images": [
            "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d",
            "https://images.unsplash.com/photo-1522708323590-d24dbb6b0267",
            "https://images.unsplash.com/photo-1502672023488-70e25813eb80"
        ]
    }
]

# Real estate in Belgrade
real_estate_belgrade = [
    {
        "title": {
            "ru": "Элитная квартира на Врачаре",
            "en": "Elite apartment in Vračar",
            "sr": "Elitni stan na Vračaru"
        },
        "description": {
            "ru": "Роскошная трёхкомнатная квартира в престижном районе Врачар. Высокие потолки, паркет, две ванные комнаты, балкон с видом на храм Святого Саввы.",
            "en": "Luxurious three-bedroom apartment in prestigious Vračar area. High ceilings, parquet, two bathrooms, balcony with Saint Sava temple view.",
            "sr": "Luksuzni trosoban stan u prestižnom kraju Vračar. Visoki plafoni, parket, dva kupatila, balkon sa pogledom na hram Svetog Save."
        },
        "price": 1500,
        "currency": "EUR",
        "category_id": 146,
        "location": {
            "latitude": 44.7988,
            "longitude": 20.4685,
            "address": "Булевар Краља Александра 45, Београд",
            "city": "Belgrade",
            "country": "Serbia"
        },
        "attributes": {
            "rooms": "3",
            "area": "95",
            "floor": "4",
            "total_floors": "6",
            "furnished": "false",
            "heating": "central",
            "parking": "true"
        },
        "images": [
            "https://images.unsplash.com/photo-1567496898669-ee935f5f647a",
            "https://images.unsplash.com/photo-1565182999561-18d7dc61c393",
            "https://images.unsplash.com/photo-1556020685-ae41abfc9365"
        ]
    },
    {
        "title": {
            "ru": "Новая квартира в Београде на воде",
            "en": "New apartment in Belgrade Waterfront",
            "sr": "Nov stan u Beogradu na vodi"
        },
        "description": {
            "ru": "Современная квартира в новом комплексе Београд на воде. Панорамный вид на реку Саву, консьерж-сервис, фитнес-центр, подземный паркинг.",
            "en": "Modern apartment in new Belgrade Waterfront complex. Panoramic Sava river view, concierge service, fitness center, underground parking.",
            "sr": "Moderan stan u novom kompleksu Beograd na vodi. Panoramski pogled na reku Savu, konsjerž servis, fitnes centar, podzemni parking."
        },
        "price": 2800,
        "currency": "EUR",
        "category_id": 146,
        "location": {
            "latitude": 44.8078,
            "longitude": 20.4448,
            "address": "Савски трг 2, Београд",
            "city": "Belgrade",
            "country": "Serbia"
        },
        "attributes": {
            "rooms": "2",
            "area": "78",
            "floor": "12",
            "total_floors": "20",
            "furnished": "true",
            "heating": "central",
            "parking": "true"
        },
        "images": [
            "https://images.unsplash.com/photo-1545324418-cc1a3fa10c00",
            "https://images.unsplash.com/photo-1556912172-45b7abe8b7e1",
            "https://images.unsplash.com/photo-1560185127-6a86733ccc3f"
        ]
    }
]

# Cars
cars = [
    {
        "title": {
            "ru": "BMW X5 2021 - идеальное состояние",
            "en": "BMW X5 2021 - perfect condition",
            "sr": "BMW X5 2021 - savršeno stanje"
        },
        "description": {
            "ru": "BMW X5 xDrive40i в отличном состоянии. Полная сервисная история, один владелец, гаражное хранение. M-пакет, панорамная крыша, адаптивная подвеска.",
            "en": "BMW X5 xDrive40i in excellent condition. Full service history, one owner, garage kept. M-package, panoramic roof, adaptive suspension.",
            "sr": "BMW X5 xDrive40i u odličnom stanju. Kompletna servisna istorija, jedan vlasnik, garažirano. M-paket, panoramski krov, adaptivno vešanje."
        },
        "price": 65000,
        "currency": "EUR",
        "category_id": 129,  # Cars
        "location": {
            "latitude": 45.2671,
            "longitude": 19.8335,
            "address": "Футошки пут 12, Нови Сад",
            "city": "Novi Sad",
            "country": "Serbia"
        },
        "attributes": {
            "make": "BMW",
            "model": "X5",
            "year": "2021",
            "mileage": "28000",
            "fuel_type": "petrol",
            "transmission": "automatic",
            "power_kw": "250",
            "color": "black"
        },
        "images": [
            "https://images.unsplash.com/photo-1555215858-9dc80e68c2c8",
            "https://images.unsplash.com/photo-1617531653332-bd46c24f2068",
            "https://images.unsplash.com/photo-1616455579100-2ceaa4eb2d37"
        ]
    },
    {
        "title": {
            "ru": "Mercedes-Benz E-Class 2020",
            "en": "Mercedes-Benz E-Class 2020",
            "sr": "Mercedes-Benz E-Class 2020"
        },
        "description": {
            "ru": "Mercedes E220d AMG Line. Безаварийный, полный AMG пакет, матричные фары, массаж сидений, проекционный дисплей.",
            "en": "Mercedes E220d AMG Line. Accident-free, full AMG package, matrix lights, seat massage, head-up display.",
            "sr": "Mercedes E220d AMG Line. Bez udesa, pun AMG paket, matrix svetla, masaža sedišta, head-up display."
        },
        "price": 48000,
        "currency": "EUR",
        "category_id": 129,
        "location": {
            "latitude": 44.8125,
            "longitude": 20.4612,
            "address": "Булевар Михаила Пупина 165, Нови Београд",
            "city": "Belgrade",
            "country": "Serbia"
        },
        "attributes": {
            "make": "Mercedes-Benz",
            "model": "E-Class",
            "year": "2020",
            "mileage": "45000",
            "fuel_type": "diesel",
            "transmission": "automatic",
            "power_kw": "143",
            "color": "silver"
        },
        "images": [
            "https://images.unsplash.com/photo-1618843479313-40f8afb4b4d8",
            "https://images.unsplash.com/photo-1606664515524-ed2f786a0bd6",
            "https://images.unsplash.com/photo-1614162692292-7ac56d7f7f1e"
        ]
    },
    {
        "title": {
            "ru": "Audi Q3 2022 - экономичный кроссовер",
            "en": "Audi Q3 2022 - economical crossover",
            "sr": "Audi Q3 2022 - ekonomičan krosover"
        },
        "description": {
            "ru": "Audi Q3 35 TFSI S-line. Виртуальная приборная панель, матричные фары, спортивные сиденья, система навигации MMI.",
            "en": "Audi Q3 35 TFSI S-line. Virtual cockpit, matrix headlights, sport seats, MMI navigation system.",
            "sr": "Audi Q3 35 TFSI S-line. Virtual cockpit, matrix farovi, sportska sedišta, MMI navigacija."
        },
        "price": 42000,
        "currency": "EUR",
        "category_id": 129,
        "location": {
            "latitude": 45.2551,
            "longitude": 19.8451,
            "address": "Темерински пут 25, Нови Сад",
            "city": "Novi Sad",
            "country": "Serbia"
        },
        "attributes": {
            "make": "Audi",
            "model": "Q3",
            "year": "2022",
            "mileage": "15000",
            "fuel_type": "petrol",
            "transmission": "automatic",
            "power_kw": "110",
            "color": "white"
        },
        "images": [
            "https://images.unsplash.com/photo-1606611013016-969c19c0f0f0",
            "https://images.unsplash.com/photo-1614026480218-547ae7abd05f",
            "https://images.unsplash.com/photo-1609521263047-f8f205293f24"
        ]
    }
]

# Electronics
electronics = [
    {
        "title": {
            "ru": "iPhone 14 Pro Max 256GB",
            "en": "iPhone 14 Pro Max 256GB",
            "sr": "iPhone 14 Pro Max 256GB"
        },
        "description": {
            "ru": "Новый iPhone 14 Pro Max в заводской пленке. Цвет Deep Purple, 256GB памяти. Полный комплект, гарантия 2 года. Возможна рассрочка.",
            "en": "Brand new iPhone 14 Pro Max in factory seal. Deep Purple color, 256GB storage. Complete package, 2 year warranty. Installment available.",
            "sr": "Novi iPhone 14 Pro Max u fabrickoj foliji. Deep Purple boja, 256GB memorije. Kompletan paket, 2 godine garancije. Moguća rata."
        },
        "price": 1200,
        "currency": "EUR",
        "category_id": 104,  # Mobile phones
        "location": {
            "latitude": 45.2551,
            "longitude": 19.8451,
            "address": "Дунавска 15, Нови Сад",
            "city": "Novi Sad",
            "country": "Serbia"
        },
        "attributes": {
            "brand": "Apple",
            "model": "iPhone 14 Pro Max",
            "storage": "256GB",
            "color": "Deep Purple",
            "condition": "new"
        },
        "images": [
            "https://images.unsplash.com/photo-1678652197831-2d180705cd2c",
            "https://images.unsplash.com/photo-1678685888221-cda773a3dcdb",
            "https://images.unsplash.com/photo-1695048133142-1a20484d2569"
        ]
    },
    {
        "title": {
            "ru": "MacBook Pro 14\" M3 Pro",
            "en": "MacBook Pro 14\" M3 Pro",
            "sr": "MacBook Pro 14\" M3 Pro"
        },
        "description": {
            "ru": "MacBook Pro 14 дюймов с процессором M3 Pro. 18GB RAM, 512GB SSD. Идеален для профессиональной работы. AppleCare+ до 2025.",
            "en": "MacBook Pro 14-inch with M3 Pro processor. 18GB RAM, 512GB SSD. Perfect for professional work. AppleCare+ until 2025.",
            "sr": "MacBook Pro 14 inča sa M3 Pro procesorom. 18GB RAM, 512GB SSD. Savršen za profesionalni rad. AppleCare+ do 2025."
        },
        "price": 2300,
        "currency": "EUR",
        "category_id": 105,  # Computers
        "location": {
            "latitude": 44.8125,
            "longitude": 20.4612,
            "address": "Кнез Михаилова 30, Београд",
            "city": "Belgrade",
            "country": "Serbia"
        },
        "attributes": {
            "brand": "Apple",
            "model": "MacBook Pro 14",
            "processor": "M3 Pro",
            "ram": "18GB",
            "storage": "512GB",
            "condition": "like new"
        },
        "images": [
            "https://images.unsplash.com/photo-1517336714731-489689fd1ca8",
            "https://images.unsplash.com/photo-1611186871348-b1ce696e52c9",
            "https://images.unsplash.com/photo-1541807084-5c52b6b3adef"
        ]
    },
    {
        "title": {
            "ru": "Samsung QLED TV 65\" 4K",
            "en": "Samsung QLED TV 65\" 4K",
            "sr": "Samsung QLED TV 65\" 4K"
        },
        "description": {
            "ru": "Samsung QLED телевизор 65 дюймов с разрешением 4K. Quantum Dot технология, 120Hz, поддержка HDR10+. Smart TV с Tizen OS.",
            "en": "Samsung QLED TV 65 inches with 4K resolution. Quantum Dot technology, 120Hz, HDR10+ support. Smart TV with Tizen OS.",
            "sr": "Samsung QLED TV 65 inča sa 4K rezolucijom. Quantum Dot tehnologija, 120Hz, HDR10+ podrška. Smart TV sa Tizen OS."
        },
        "price": 1100,
        "currency": "EUR",
        "category_id": 106,  # TVs
        "location": {
            "latitude": 45.2467,
            "longitude": 19.8515,
            "address": "Булевар Ослобођења 88, Нови Сад",
            "city": "Novi Sad",
            "country": "Serbia"
        },
        "attributes": {
            "brand": "Samsung",
            "screen_size": "65",
            "resolution": "4K",
            "technology": "QLED",
            "smart_tv": "true"
        },
        "images": [
            "https://images.unsplash.com/photo-1593359677879-a4bb92f829d1",
            "https://images.unsplash.com/photo-1567690187548-f07b1d7bf5a9",
            "https://images.unsplash.com/photo-1558888401-3cc1de77652d"
        ]
    }
]

# Furniture
furniture = [
    {
        "title": {
            "ru": "Итальянский кожаный диван",
            "en": "Italian leather sofa",
            "sr": "Italijanska kožna garnitura"
        },
        "description": {
            "ru": "Роскошный трёхместный диван из натуральной итальянской кожи. Цвет коньяк, ручная работа, эргономичный дизайн. Размеры: 230x95x85 см.",
            "en": "Luxurious three-seater sofa made of genuine Italian leather. Cognac color, handmade, ergonomic design. Dimensions: 230x95x85 cm.",
            "sr": "Luksuzna trosed od prave italijanske kože. Konjak boja, ručni rad, ergonomski dizajn. Dimenzije: 230x95x85 cm."
        },
        "price": 1800,
        "currency": "EUR",
        "category_id": 164,  # Furniture
        "location": {
            "latitude": 44.7866,
            "longitude": 20.4489,
            "address": "Теразије 25, Београд",
            "city": "Belgrade",
            "country": "Serbia"
        },
        "attributes": {
            "material": "leather",
            "seats": "3",
            "color": "cognac",
            "condition": "new"
        },
        "images": [
            "https://images.unsplash.com/photo-1555041469-a586c61ea9bc",
            "https://images.unsplash.com/photo-1493663284031-b7e3aefcae8e",
            "https://images.unsplash.com/photo-1549187774-b4e9b0445b41"
        ]
    },
    {
        "title": {
            "ru": "Обеденный стол из массива дуба",
            "en": "Solid oak dining table",
            "sr": "Trpezarijski sto od punog hrasta"
        },
        "description": {
            "ru": "Массивный обеденный стол из натурального дуба. Вместимость 8 человек, раздвижная конструкция. Размеры: 200-250x100x75 см.",
            "en": "Massive dining table made of solid oak. Seats 8 people, extendable design. Dimensions: 200-250x100x75 cm.",
            "sr": "Masivan trpezarijski sto od prirodnog hrasta. Za 8 osoba, rasklopiv dizajn. Dimenzije: 200-250x100x75 cm."
        },
        "price": 950,
        "currency": "EUR",
        "category_id": 164,
        "location": {
            "latitude": 45.2551,
            "longitude": 19.8451,
            "address": "Железничка 4, Нови Сад",
            "city": "Novi Sad",
            "country": "Serbia"
        },
        "attributes": {
            "material": "oak",
            "seats": "8",
            "extendable": "true",
            "condition": "new"
        },
        "images": [
            "https://images.unsplash.com/photo-1549497538-303791108f95",
            "https://images.unsplash.com/photo-1571089086084-e8b3f2f3e72f",
            "https://images.unsplash.com/photo-1581539250439-c96689b516dd"
        ]
    }
]

def download_and_upload_image(image_url):
    """Download image from URL and upload to MinIO via API"""
    try:
        # Download image
        response = requests.get(image_url + "?w=800&q=80", timeout=10)
        if response.status_code != 200:
            print(f"Failed to download image: {image_url}")
            return None
        
        # Prepare multipart upload
        files = {
            'images': ('image.jpg', response.content, 'image/jpeg')
        }
        
        # Upload to API
        upload_response = requests.post(
            f"{API_BASE}/images/upload",
            headers={"Authorization": f"Bearer {JWT_TOKEN}"},
            files=files
        )
        
        if upload_response.status_code == 200:
            result = upload_response.json()
            if result.get('data') and len(result['data']) > 0:
                return result['data'][0]
        
        print(f"Failed to upload image: {upload_response.text}")
        return None
    except Exception as e:
        print(f"Error processing image {image_url}: {e}")
        return None

def create_listing(listing_data):
    """Create a listing via API"""
    try:
        # First upload images
        uploaded_images = []
        for img_url in listing_data.get('images', []):
            uploaded_img = download_and_upload_image(img_url)
            if uploaded_img:
                uploaded_images.append(uploaded_img)
            time.sleep(1)  # Rate limiting
        
        # Prepare listing payload
        payload = {
            "title": listing_data['title']['ru'],
            "description": listing_data['description']['ru'],
            "price": listing_data['price'],
            "currency": listing_data['currency'],
            "category_id": listing_data['category_id'],
            "latitude": listing_data['location']['latitude'],
            "longitude": listing_data['location']['longitude'],
            "address": listing_data['location']['address'],
            "city": listing_data['location']['city'],
            "country": listing_data['location']['country'],
            "images": uploaded_images,
            "attributes": listing_data.get('attributes', {}),
            "translations": {
                "title": listing_data['title'],
                "description": listing_data['description']
            }
        }
        
        # Create listing
        response = requests.post(
            f"{API_BASE}/marketplace/listings",
            headers=headers,
            json=payload
        )
        
        if response.status_code == 200 or response.status_code == 201:
            result = response.json()
            print(f"✅ Created listing: {listing_data['title']['ru']}")
            return result.get('data')
        else:
            print(f"❌ Failed to create listing: {response.text}")
            return None
            
    except Exception as e:
        print(f"Error creating listing: {e}")
        return None

def main():
    """Main function to create all listings"""
    print("🚀 Starting to create real listings...")
    
    all_listings = []
    
    # Create real estate listings in Novi Sad
    print("\n📍 Creating real estate listings in Novi Sad...")
    for listing in real_estate_novi_sad:
        result = create_listing(listing)
        if result:
            all_listings.append(result)
        time.sleep(2)
    
    # Create real estate listings in Belgrade
    print("\n📍 Creating real estate listings in Belgrade...")
    for listing in real_estate_belgrade:
        result = create_listing(listing)
        if result:
            all_listings.append(result)
        time.sleep(2)
    
    # Create car listings
    print("\n🚗 Creating car listings...")
    for listing in cars:
        result = create_listing(listing)
        if result:
            all_listings.append(result)
        time.sleep(2)
    
    # Create electronics listings
    print("\n📱 Creating electronics listings...")
    for listing in electronics:
        result = create_listing(listing)
        if result:
            all_listings.append(result)
        time.sleep(2)
    
    # Create furniture listings
    print("\n🪑 Creating furniture listings...")
    for listing in furniture:
        result = create_listing(listing)
        if result:
            all_listings.append(result)
        time.sleep(2)
    
    print(f"\n✨ Successfully created {len(all_listings)} listings!")
    
    # Save listing IDs for future reference
    with open('/data/hostel-booking-system/created_listings.json', 'w') as f:
        json.dump(all_listings, f, indent=2)
    
    print("📝 Listing IDs saved to created_listings.json")

if __name__ == "__main__":
    main()