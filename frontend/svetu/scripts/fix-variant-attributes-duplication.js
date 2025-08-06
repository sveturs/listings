const fs = require('fs');
const path = require('path');

// Функция для объединения дублированных секций variantAttributes
function fixVariantAttributesDuplication(locale) {
  const filePath = path.join(__dirname, `../src/messages/${locale}/admin.json`);

  console.log(`Обработка файла: ${filePath}`);

  const content = fs.readFileSync(filePath, 'utf8');
  const data = JSON.parse(content);

  // Проверяем наличие дублированных секций
  if (!data.admin || !data.admin.variantAttributes) {
    console.log(
      `Файл ${locale}/admin.json не содержит admin.variantAttributes`
    );
    return;
  }

  // Создаем объединенную секцию variantAttributes
  const mergedVariantAttributes = {
    // Основные ключи из второй секции (более полной)
    title:
      locale === 'ru'
        ? 'Вариативные атрибуты'
        : locale === 'en'
          ? 'Variant Attributes'
          : 'Varijantni atributi',
    description:
      locale === 'ru'
        ? 'Управление атрибутами для вариантов товаров'
        : locale === 'en'
          ? 'Manage attributes for product variants'
          : 'Upravljanje atributima za varijante proizvoda',
    addAttribute:
      locale === 'ru'
        ? 'Добавить вариативный атрибут'
        : locale === 'en'
          ? 'Add Variant Attribute'
          : 'Dodaj varijantni atribut',
    editAttribute:
      locale === 'ru'
        ? 'Редактировать вариативный атрибут'
        : locale === 'en'
          ? 'Edit Variant Attribute'
          : 'Izmeni varijantni atribut',
    systemName:
      locale === 'ru'
        ? 'Системное имя'
        : locale === 'en'
          ? 'System Name'
          : 'Sistemsko ime',
    displayName:
      locale === 'ru'
        ? 'Отображаемое название'
        : locale === 'en'
          ? 'Display Name'
          : 'Naziv za prikaz',
    type: locale === 'ru' ? 'Тип' : locale === 'en' ? 'Type' : 'Tip',
    isRequired:
      locale === 'ru'
        ? 'Обязательный'
        : locale === 'en'
          ? 'Required'
          : 'Obavezno',
    affectsStock:
      locale === 'ru'
        ? 'Влияет на остатки'
        : locale === 'en'
          ? 'Affects Stock'
          : 'Utiče na zalihe',
    settings:
      locale === 'ru'
        ? 'Настройки'
        : locale === 'en'
          ? 'Settings'
          : 'Podešavanja',
    sortOrder:
      locale === 'ru'
        ? 'Порядок сортировки'
        : locale === 'en'
          ? 'Sort Order'
          : 'Redosled sortiranja',
    systemNamePlaceholder:
      locale === 'ru'
        ? 'например: color, size, memory'
        : locale === 'en'
          ? 'e.g.: color, size, memory'
          : 'npr.: color, size, memory',
    displayNamePlaceholder:
      locale === 'ru'
        ? 'например: Цвет, Размер, Память'
        : locale === 'en'
          ? 'e.g.: Color, Size, Memory'
          : 'npr.: Boja, Veličina, Memorija',
    systemNameHint:
      locale === 'ru'
        ? 'Только латинские буквы, цифры и подчеркивания'
        : locale === 'en'
          ? 'Only latin letters, numbers and underscores'
          : 'Samo latinična slova, brojevi i donje crte',
    sortOrderHint:
      locale === 'ru'
        ? 'Порядок отображения в интерфейсе'
        : locale === 'en'
          ? 'Display order in interface'
          : 'Redosled prikaza u interfejsu',
    isRequiredHint:
      locale === 'ru'
        ? 'Обязательно ли указывать значение при создании варианта'
        : locale === 'en'
          ? 'Whether value is required when creating variant'
          : 'Da li je vrednost obavezna pri kreiranju varijante',
    affectsStockHint:
      locale === 'ru'
        ? 'Влияют ли разные значения на отдельный учет остатков'
        : locale === 'en'
          ? 'Whether different values affect separate stock tracking'
          : 'Da li različite vrednosti utiču na odvojeno praćenje zaliha',
    validationError:
      locale === 'ru'
        ? 'Заполните обязательные поля'
        : locale === 'en'
          ? 'Please fill in required fields'
          : 'Molimo popunite obavezna polja',
    allTypes:
      locale === 'ru'
        ? 'Все типы'
        : locale === 'en'
          ? 'All Types'
          : 'Svi tipovi',
    // Объединенная секция types со всеми типами
    types: {
      text: locale === 'ru' ? 'Текст' : locale === 'en' ? 'Text' : 'Tekst',
      number: locale === 'ru' ? 'Число' : locale === 'en' ? 'Number' : 'Broj',
      select: locale === 'ru' ? 'Выбор' : locale === 'en' ? 'Select' : 'Izbor',
      multiselect:
        locale === 'ru'
          ? 'Множественный выбор'
          : locale === 'en'
            ? 'Multiple select'
            : 'Višestruki izbor',
      boolean:
        locale === 'ru' ? 'Да/Нет' : locale === 'en' ? 'Boolean' : 'Da/Ne',
      date: locale === 'ru' ? 'Дата' : locale === 'en' ? 'Date' : 'Datum',
      range: locale === 'ru' ? 'Диапазон' : locale === 'en' ? 'Range' : 'Opseg',
      color: locale === 'ru' ? 'Цвет' : locale === 'en' ? 'Color' : 'Boja',
      size: locale === 'ru' ? 'Размер' : locale === 'en' ? 'Size' : 'Veličina',
      memory:
        locale === 'ru' ? 'Память' : locale === 'en' ? 'Memory' : 'Memorija',
      storage:
        locale === 'ru'
          ? 'Хранилище'
          : locale === 'en'
            ? 'Storage'
            : 'Skladištenje',
      material:
        locale === 'ru'
          ? 'Материал'
          : locale === 'en'
            ? 'Material'
            : 'Materijal',
      capacity:
        locale === 'ru'
          ? 'Емкость'
          : locale === 'en'
            ? 'Capacity'
            : 'Kapacitet',
      power: locale === 'ru' ? 'Мощность' : locale === 'en' ? 'Power' : 'Snaga',
      connectivity:
        locale === 'ru'
          ? 'Подключение'
          : locale === 'en'
            ? 'Connectivity'
            : 'Povezivanje',
      style: locale === 'ru' ? 'Стиль' : locale === 'en' ? 'Style' : 'Stil',
      pattern:
        locale === 'ru' ? 'Узор' : locale === 'en' ? 'Pattern' : 'Uzorak',
      weight: locale === 'ru' ? 'Вес' : locale === 'en' ? 'Weight' : 'Težina',
      bundle:
        locale === 'ru' ? 'Комплект' : locale === 'en' ? 'Bundle' : 'Paket',
    },
    // Остальные ключи из второй секции
    manageLinks:
      locale === 'ru'
        ? 'Управление связями'
        : locale === 'en'
          ? 'Manage Links'
          : 'Upravljanje vezama',
    manageMappings:
      locale === 'ru'
        ? 'Управление связями атрибутов'
        : locale === 'en'
          ? 'Manage Attribute Mappings'
          : 'Upravljanje mapiranjem atributa',
    mappingTitle:
      locale === 'ru'
        ? 'Связывание атрибутов для'
        : locale === 'en'
          ? 'Mapping attributes for'
          : 'Mapiranje atributa za',
    mappingDescription:
      locale === 'ru'
        ? 'Выберите атрибуты категорий, которые будут связаны с этим вариативным атрибутом'
        : locale === 'en'
          ? 'Select category attributes that will be linked to this variant attribute'
          : 'Izaberite atribute kategorija koji će biti povezani sa ovim varijantnim atributom',
    searchAttributes:
      locale === 'ru'
        ? 'Поиск атрибутов...'
        : locale === 'en'
          ? 'Search attributes...'
          : 'Pretraži atribute...',
    autoDetect:
      locale === 'ru'
        ? 'Автоопределение'
        : locale === 'en'
          ? 'Auto-detect'
          : 'Automatsko prepoznavanje',
    autoDetectSuccess:
      locale === 'ru'
        ? 'Автоматически найдено связей: {count}'
        : locale === 'en'
          ? 'Auto-detected links: {count}'
          : 'Automatski pronađeno veza: {count}',
    autoDetectNoResults:
      locale === 'ru'
        ? 'Автоматически не удалось найти подходящие атрибуты'
        : locale === 'en'
          ? 'Could not auto-detect matching attributes'
          : 'Nije moguće automatski pronaći odgovarajuće atribute',
    noAttributesFound:
      locale === 'ru'
        ? 'Атрибуты не найдены'
        : locale === 'en'
          ? 'No attributes found'
          : 'Atributi nisu pronađeni',
    previouslyLinked:
      locale === 'ru'
        ? 'Ранее связан'
        : locale === 'en'
          ? 'Previously linked'
          : 'Prethodno povezano',
    selectedCount:
      locale === 'ru'
        ? 'Выбрано атрибутов: {count}'
        : locale === 'en'
          ? 'Selected attributes: {count}'
          : 'Izabrano atributa: {count}',
    saveMappings:
      locale === 'ru'
        ? 'Сохранить связи'
        : locale === 'en'
          ? 'Save mappings'
          : 'Sačuvaj mapiranja',
    saving:
      locale === 'ru'
        ? 'Сохранение...'
        : locale === 'en'
          ? 'Saving...'
          : 'Čuvanje...',
    cancel: locale === 'ru' ? 'Отмена' : locale === 'en' ? 'Cancel' : 'Otkaži',
    loadMappingsError:
      locale === 'ru'
        ? 'Ошибка загрузки связей'
        : locale === 'en'
          ? 'Error loading mappings'
          : 'Greška pri učitavanju mapiranja',
    saveMappingsSuccess:
      locale === 'ru'
        ? 'Связи успешно сохранены'
        : locale === 'en'
          ? 'Mappings saved successfully'
          : 'Mapiranja uspešno sačuvana',
    saveMappingsError:
      locale === 'ru'
        ? 'Ошибка сохранения связей'
        : locale === 'en'
          ? 'Error saving mappings'
          : 'Greška pri čuvanju mapiranja',
    dragDropTitle:
      locale === 'ru'
        ? 'Перетаскивание атрибутов'
        : locale === 'en'
          ? 'Drag and Drop Attributes'
          : 'Prevuci i pusti atribute',
    dragDropHint:
      locale === 'ru'
        ? 'Перетащите атрибуты из левой колонки в правую для создания связей. Используйте автоопределение для быстрого поиска похожих атрибутов.'
        : locale === 'en'
          ? 'Drag attributes from the left column to the right to create links. Use auto-detect to quickly find similar attributes.'
          : 'Prevucite atribute iz leve kolone u desnu za kreiranje veza. Koristite automatsko prepoznavanje za brzo pronalaženje sličnih atributa.',
    availableAttributes:
      locale === 'ru'
        ? 'Доступные атрибуты'
        : locale === 'en'
          ? 'Available Attributes'
          : 'Dostupni atributi',
    linkedAttributes:
      locale === 'ru'
        ? 'Связанные атрибуты'
        : locale === 'en'
          ? 'Linked Attributes'
          : 'Povezani atributi',
    linkedWith:
      locale === 'ru'
        ? 'Связано с'
        : locale === 'en'
          ? 'Linked with'
          : 'Povezano sa',
    noMatchingAttributes:
      locale === 'ru'
        ? 'Нет подходящих атрибутов'
        : locale === 'en'
          ? 'No matching attributes'
          : 'Nema odgovarajućih atributa',
    noAvailableAttributes:
      locale === 'ru'
        ? 'Нет доступных атрибутов'
        : locale === 'en'
          ? 'No available attributes'
          : 'Nema dostupnih atributa',
    noLinkedAttributes:
      locale === 'ru'
        ? 'Нет связанных атрибутов. Перетащите атрибуты сюда.'
        : locale === 'en'
          ? 'No linked attributes. Drag attributes here.'
          : 'Nema povezanih atributa. Prevucite atribute ovde.',
    autoDetectError:
      locale === 'ru'
        ? 'Ошибка автоопределения'
        : locale === 'en'
          ? 'Auto-detect error'
          : 'Greška automatskog prepoznavanja',
    mappingsUpdated:
      locale === 'ru'
        ? 'Связи обновлены'
        : locale === 'en'
          ? 'Mappings updated'
          : 'Mapiranja ažurirana',
  };

  // Обновляем секцию variantAttributes
  data.admin.variantAttributes = mergedVariantAttributes;

  // Удаляем дублированные секции если они есть в корне объекта
  if (data.variantAttributes) {
    delete data.variantAttributes;
  }

  // Записываем обратно
  fs.writeFileSync(filePath, JSON.stringify(data, null, 2) + '\n', 'utf8');
  console.log(`Файл ${locale}/admin.json исправлен успешно`);
}

// Обрабатываем все языки
['ru', 'en', 'sr'].forEach((locale) => {
  try {
    fixVariantAttributesDuplication(locale);
  } catch (error) {
    console.error(`Ошибка при обработке ${locale}:`, error);
  }
});

console.log('\nИсправление дублирования variantAttributes завершено!');
