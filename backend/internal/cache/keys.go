package cache

import "fmt"

// Константы для префиксов ключей кеша
const (
	PrefixCategories      = "categories:"
	PrefixCategoryTree    = "category_tree:"
	PrefixCategoryAttrs   = "category_attrs:"
	PrefixAttributeGroups = "attribute_groups:"
	PrefixAttribute       = "attribute:"
	PrefixTranslations    = "translations:"
)

// BuildCategoriesKey формирует ключ для списка всех категорий
func BuildCategoriesKey(locale string) string {
	return fmt.Sprintf("%sall:%s", PrefixCategories, locale)
}

// BuildCategoryTreeKey формирует ключ для дерева категорий
func BuildCategoryTreeKey(locale string, activeOnly bool) string {
	if activeOnly {
		return fmt.Sprintf("%s%s:active", PrefixCategoryTree, locale)
	}
	return fmt.Sprintf("%s%s:all", PrefixCategoryTree, locale)
}

// BuildCategoryKey формирует ключ для конкретной категории
func BuildCategoryKey(categoryID int64) string {
	return fmt.Sprintf("%sid:%d", PrefixCategories, categoryID)
}

// BuildCategoryAttributesKey формирует ключ для атрибутов категории
func BuildCategoryAttributesKey(categoryID int64, locale string) string {
	return fmt.Sprintf("%s%d:%s", PrefixCategoryAttrs, categoryID, locale)
}

// BuildAttributeGroupsKey формирует ключ для групп атрибутов категории
func BuildAttributeGroupsKey(categoryID int64, locale string) string {
	return fmt.Sprintf("%scategory:%d:%s", PrefixAttributeGroups, categoryID, locale)
}

// BuildAttributeKey формирует ключ для конкретного атрибута
func BuildAttributeKey(attributeID int64) string {
	return fmt.Sprintf("%sid:%d", PrefixAttribute, attributeID)
}

// BuildTranslationStatusKey формирует ключ для статуса переводов
func BuildTranslationStatusKey(entityType string) string {
	return fmt.Sprintf("%sstatus:%s", PrefixTranslations, entityType)
}

// BuildCategoryInvalidationPattern формирует паттерн для инвалидации всех ключей категории
func BuildCategoryInvalidationPattern(categoryID int64) string {
	return fmt.Sprintf("*:%d:*", categoryID)
}

// BuildAllCategoriesInvalidationPattern формирует паттерн для инвалидации всех категорий
func BuildAllCategoriesInvalidationPattern() string {
	return fmt.Sprintf("%s*", PrefixCategories)
}

// BuildAttributeInvalidationPattern формирует паттерн для инвалидации атрибута
func BuildAttributeInvalidationPattern(attributeID int64) string {
	return fmt.Sprintf("*attribute*:%d:*", attributeID)
}
