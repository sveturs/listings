import React, { useState } from 'react';
import { 
  Container, 
  Typography, 
  Box, 
  Paper, 
  Divider, 
  ToggleButtonGroup, 
  ToggleButton,
  useTheme,
  useMediaQuery 
} from '@mui/material';

// Полные тексты политик на разных языках
const policyTexts = {
    en: {
      title: "Privacy Policy",
      lastUpdated: "Last updated: December 03, 2024",
      sections: [
        {
          title: "1. Introduction",
          content: `LandHub.rs ("we", "our", or "us") is committed to protecting your privacy. 
          This Privacy Policy explains our practices regarding the collection, use, disclosure, 
          and protection of your information when you use our service.`,
        },
        {
          title: "2. Information We Collect",
          subsections: [
            {
              title: "2.1 Information You Provide",
              items: [
                "Personal identification information (name, email address, phone number)",
                "Booking and reservation details",
                "Payment and transaction information",
                "Profile information and preferences",
                "Communication records and feedback",
                "Identity verification documents when required",
                "Any other information you choose to provide"
              ]
            },
            {
              title: "2.2 Information Automatically Collected",
              items: [
                "Device and browser information",
                "IP address and location data",
                "Cookie data and unique identifiers",
                "Usage statistics and interaction data",
                "Performance and error data",
                "Session duration and patterns",
                "Referral sources and navigation paths"
              ]
            },
            {
              title: "2.3 Information from Third Parties",
              items: [
                "Property information from Booking.com through their official API",
                "Payment processing services information",
                "Identity verification services data",
                "Analytics data",
                "Marketing and advertising partner data"
              ]
            }
          ]
        },
        {
          title: "3. Use of Information",
          content: "We use your information for the following purposes:",
          items: [
            "Processing and managing your bookings",
            "Providing customer support and responding to inquiries",
            "Improving and personalizing our services",
            "Sending important notifications about your bookings",
            "Marketing and promotional communications (with your consent)",
            "Preventing fraud and ensuring platform security",
            "Complying with legal obligations"
          ]
        },
        {
          title: "4. Data Protection and Security",
          content: "We implement robust security measures including:",
          items: [
            "Encryption of data in transit and at rest",
            "Regular security assessments and audits",
            "Access controls and authentication",
            "Employee training on data protection",
            "Incident response procedures",
            "Regular backups and disaster recovery plans"
          ]
        },
        {
          title: "5. Booking.com Integration",
          content: `Our platform integrates with Booking.com to provide additional accommodation options. 
          This integration involves:`,
          items: [
            "Secure API access following Booking.com's requirements",
            "Display of property information and availability",
            "Booking processing through Booking.com's systems",
            "Data sharing in accordance with Booking.com's terms",
            "Clear identification of Booking.com properties"
          ]
        },
        {
          title: "6. Your Rights",
          content: "Under applicable data protection laws, you have the following rights:",
          items: [
            "Right to access your personal data",
            "Right to correct inaccurate data",
            "Right to erasure ('right to be forgotten')",
            "Right to restrict processing",
            "Right to data portability",
            "Right to object to processing",
            "Right to withdraw consent"
          ]
        },
        {
          title: "7. Data Retention",
          content: "We retain your personal data for the following periods:",
          items: [
            "Account information: While your account is active plus 2 years",
            "Booking records: 7 years (tax requirements)",
            "Marketing data: Until consent withdrawal",
            "Technical logs: 26 months",
            "Payment information: As required by law"
          ]
        },
        {
          title: "8. International Data Transfers",
          content: `We may transfer your data internationally. We ensure appropriate safeguards through:`,
          items: [
            "Standard Contractual Clauses (SCCs)",
            "Adequacy decisions",
            "Binding Corporate Rules",
            "Data Processing Agreements"
          ]
        },
        {
          title: "9. Contact Us",
          content: `For any questions about this Privacy Policy or our data practices, contact us at:`,
          items: [
            "Email: privacy@landhub.rs",
            "Address: [Your Legal Address], Novi Sad, Serbia",
            "Phone: +381 XX XXX XXX",
            "Data Protection Officer: dpo@landhub.rs"
          ]
        }
      ]
    },
    sr: {
      title: "Politika Privatnosti",
      lastUpdated: "Poslednji put ažurirano: 03.12.2024",
      sections: [
        {
          title: "1. Uvod",
          content: `LandHub.rs („mi", „naš" ili „nas") je posvećen zaštiti vaše privatnosti. 
          U skladu sa Zakonom o zaštiti podataka o ličnosti Republike Srbije 
          („Sl. glasnik RS", br. 87/2018), ova Politika privatnosti objašnjava naše prakse 
          u vezi sa prikupljanjem, korišćenjem i zaštitom vaših podataka.`
        },
        {
          title: "2. Podaci koje prikupljamo",
          subsections: [
            {
              title: "2.1 Podaci koje vi dostavljate",
              items: [
                "Lični identifikacioni podaci (ime, email adresa, telefon)",
                "Podaci o rezervacijama",
                "Podaci o plaćanju",
                "Informacije o profilu",
                "Evidencija komunikacije",
                "Dokumenti za verifikaciju identiteta",
                "Ostali podaci koje odlučite da dostavite"
              ]
            },
            {
              title: "2.2 Automatski prikupljeni podaci",
              items: [
                "Podaci o uređaju i pretraživaču",
                "IP adresa i lokacijski podaci",
                "Kolačići i jedinstveni identifikatori",
                "Statistika korišćenja",
                "Podaci o performansama",
                "Trajanje sesije",
                "Izvori poseta"
              ]
            }
          ]
        },
        {
          title: "3. Pravni osnov obrade",
          content: "U skladu sa ZZPL, obrada se vrši na osnovu:",
          items: [
            "Vašeg pristanka",
            "Izvršenja ugovora",
            "Zakonskih obaveza",
            "Legitimnih interesa"
          ]
        },
        {
          title: "4. Mere zaštite podataka",
          content: "Primenjujemo sledeće mere zaštite:",
          items: [
            "Enkripcija podataka",
            "Redovne provere bezbednosti",
            "Kontrola pristupa",
            "Obuka zaposlenih",
            "Procedure za incidente",
            "Redovno pravljenje rezervnih kopija"
          ]
        },
        {
          title: "5. Integracija sa Booking.com",
          content: "Naša platforma je integrisana sa Booking.com sistemom:",
          items: [
            "Siguran API pristup",
            "Prikaz informacija o smeštaju",
            "Obrada rezervacija",
            "Deljenje podataka u skladu sa uslovima",
            "Jasno označavanje Booking.com ponuda"
          ]
        },
        {
          title: "6. Vaša prava",
          content: "Prema ZZPL, imate sledeća prava:",
          items: [
            "Pravo na pristup podacima",
            "Pravo na ispravku",
            "Pravo na brisanje",
            "Pravo na ograničenje obrade",
            "Pravo na prenosivost",
            "Pravo na prigovor",
            "Pravo na opoziv pristanka"
          ]
        },
        {
          title: "7. Rok čuvanja podataka",
          content: "Vaše podatke čuvamo u sledećim rokovima:",
          items: [
            "Podaci o nalogu: Dok je nalog aktivan plus 2 godine",
            "Podaci o rezervacijama: 7 godina",
            "Marketinški podaci: Do opoziva pristanka",
            "Tehnički zapisi: 26 meseci"
          ]
        },
        {
          title: "8. Kontakt",
          content: "Za sva pitanja o zaštiti podataka, kontaktirajte nas:",
          items: [
            "Email: privacy@landhub.rs",
            "Adresa: [Vaša adresa], Novi Sad, Srbija",
            "Telefon: +381 XX XXX XXX",
            "Lice za zaštitu podataka: dpo@landhub.rs"
          ]
        }
      ]
    },
    ru: {
      title: "Политика Конфиденциальности",
      lastUpdated: "Последнее обновление: 03.12.2024",
      sections: [
        {
          title: "1. Введение",
          content: `LandHub.rs ("мы", "наш" или "нас") обеспечивает защиту вашей конфиденциальности. 
          В соответствии с Федеральным законом №152-ФЗ "О персональных данных", 
          настоящая Политика конфиденциальности объясняет наши практики в отношении 
          сбора, использования и защиты ваших данных.`
        },
        {
          title: "2. Собираемые данные",
          subsections: [
            {
              title: "2.1 Предоставляемые вами данные",
              items: [
                "Личные идентификационные данные (имя, email, телефон)",
                "Данные бронирований",
                "Платежная информация",
                "Информация профиля",
                "Записи коммуникаций",
                "Документы для верификации личности",
                "Иные предоставляемые вами данные"
              ]
            },
            {
              title: "2.2 Автоматически собираемые данные",
              items: [
                "Данные об устройстве и браузере",
                "IP-адрес и данные о местоположении",
                "Cookies и уникальные идентификаторы",
                "Статистика использования",
                "Данные о производительности",
                "Длительность сессий",
                "Источники переходов"
              ]
            }
          ]
        },
        {
          title: "3. Локализация данных",
          content: "В соответствии с требованиями 152-ФЗ, мы обеспечиваем:",
          items: [
            "Хранение персональных данных граждан РФ на территории РФ",
            "Уведомление Роскомнадзора об обработке данных",
            "Соблюдение требований к защите данных",
            "Локальное резервное копирование"
          ]
        },
        {
          title: "4. Защита данных",
          content: "Мы применяем следующие меры защиты:",
          items: [
            "Шифрование данных",
            "Регулярные проверки безопасности",
            "Контроль доступа",
            "Обучение сотрудников",
            "Процедуры реагирования на инциденты",
            "Резервное копирование"
          ]
        },
        {
          title: "5. Интеграция с Booking.com",
          content: "Наша платформа интегрирована с системой Booking.com:",
          items: [
            "Безопасный API-доступ",
            "Отображение информации о размещении",
            "Обработка бронирований",
            "Обмен данными согласно условиям",
            "Четкая маркировка предложений Booking.com"
          ]
        },
        {
          title: "6. Ваши права",
          content: "В соответствии с 152-ФЗ, вы имеете право:",
          items: [
            "На доступ к своим данным",
            "На исправление данных",
            "На удаление данных",
            "На ограничение обработки",
            "На перенос данных",
            "На возражение против обработки",
            "На отзыв согласия"
          ]
        },
        {
          title: "7. Сроки хранения",
          content: "Мы храним ваши данные в следующие сроки:",
          items: [
            "Данные аккаунта: Пока аккаунт активен плюс 2 года",
            "Данные бронирований: 7 лет",
            "Маркетинговые данные: До отзыва согласия",
            "Технические логи: 26 месяцев"
          ]
        },
        {
          title: "8. Контакты",
          content: "По всем вопросам о защите данных обращайтесь:",
          items: [
            "Email: privacy@landhub.rs",
            "Адрес: [Ваш адрес], Нови-Сад, Сербия",
            "Телефон: +381 XX XXX XXX",
            "Ответственный за защиту данных: dpo@landhub.rs"
          ]
        }
      ]
    }
  };
const PrivacyPolicy = () => {
  const [language, setLanguage] = useState('en');
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const handleLanguageChange = (event, newLanguage) => {
    if (newLanguage !== null) {
      setLanguage(newLanguage);
    }
  };

  const renderContent = (sections) => {
    return sections.map((section, index) => (
      <Box key={index} mb={4}>
        <Typography variant={isMobile ? "h6" : "h5"} component="h2" gutterBottom>
          {section.title}
        </Typography>
        
        {section.content && (
          <Typography paragraph>
            {section.content}
          </Typography>
        )}
        
        {section.items && (
          <Box sx={{ pl: 2 }}>
            <ul>
              {section.items.map((item, itemIndex) => (
                <Typography component="li" key={itemIndex} paragraph>
                  {item}
                </Typography>
              ))}
            </ul>
          </Box>
        )}
        
        {section.subsections && section.subsections.map((subsection, subIndex) => (
          <Box key={subIndex} ml={2} mb={2}>
            <Typography variant="h6" component="h3" gutterBottom>
              {subsection.title}
            </Typography>
            {subsection.items && (
              <ul>
                {subsection.items.map((item, itemIndex) => (
                  <Typography component="li" key={itemIndex} paragraph>
                    {item}
                  </Typography>
                ))}
              </ul>
            )}
          </Box>
        ))}
      </Box>
    ));
  };

  return (
    <Container maxWidth="md">
      <Box py={4}>
        <Box mb={3} display="flex" justifyContent="flex-end">
          <ToggleButtonGroup
            value={language}
            exclusive
            onChange={handleLanguageChange}
            size="small"
          >
            <ToggleButton value="en">English</ToggleButton>
            <ToggleButton value="sr">Srpski</ToggleButton>
            <ToggleButton value="ru">Русский</ToggleButton>
          </ToggleButtonGroup>
        </Box>

        <Paper elevation={1} sx={{ p: { xs: 2, md: 4 } }}>
          <Typography 
            variant={isMobile ? "h4" : "h3"} 
            component="h1" 
            gutterBottom
          >
            {policyTexts[language].title}
          </Typography>
          
          <Typography 
            variant="subtitle1" 
            color="text.secondary" 
            paragraph
          >
            {policyTexts[language].lastUpdated}
          </Typography>

          <Divider sx={{ my: 3 }} />

          {renderContent(policyTexts[language].sections)}
        </Paper>
      </Box>
    </Container>
  );
};

export default PrivacyPolicy;