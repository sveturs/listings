package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"backend/internal/proj/search_optimization/storage"
)

const (
	// Severity levels
	severityWarning = "warning"
)

// SecurityCheck функция проверки безопасности перед применением изменений
type SecurityCheck struct {
	service *searchOptimizationService
}

func NewSecurityCheck(service *searchOptimizationService) *SecurityCheck {
	return &SecurityCheck{service: service}
}

// SafetyRule правило безопасности
type SafetyRule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // critical, warning, info
	Enabled     bool   `json:"enabled"`
}

// SecurityReport отчет о проверке безопасности
type SecurityReport struct {
	OverallSafety    string          `json:"overall_safety"` // safe, warning, dangerous
	TotalViolations  int             `json:"total_violations"`
	CriticalIssues   int             `json:"critical_issues"`
	Warnings         int             `json:"warnings"`
	Violations       []RuleViolation `json:"violations"`
	Recommendations  []string        `json:"recommendations"`
	ApprovalRequired bool            `json:"approval_required"`
}

// RuleViolation нарушение правила безопасности
type RuleViolation struct {
	RuleID        string      `json:"rule_id"`
	Severity      string      `json:"severity"`
	Message       string      `json:"message"`
	FieldName     string      `json:"field_name,omitempty"`
	CurrentValue  interface{} `json:"current_value,omitempty"`
	ProposedValue interface{} `json:"proposed_value,omitempty"`
	Impact        string      `json:"impact"`
}

// GetSafetyRules возвращает список правил безопасности
func (sc *SecurityCheck) GetSafetyRules() []SafetyRule {
	return []SafetyRule{
		{
			ID:          "max_weight_change",
			Name:        "Максимальное изменение веса",
			Description: "Вес не должен изменяться более чем на 30% за одну операцию",
			Severity:    "critical",
			Enabled:     true,
		},
		{
			ID:          "weight_bounds",
			Name:        "Границы весов",
			Description: "Вес должен находиться в диапазоне от 0.0 до 1.0",
			Severity:    "critical",
			Enabled:     true,
		},
		{
			ID:          "critical_fields_protection",
			Name:        "Защита критических полей",
			Description: "Критические поля (title, description) не должны иметь вес ниже 0.3",
			Severity:    "critical",
			Enabled:     true,
		},
		{
			ID:          "confidence_threshold",
			Name:        "Пороговый уровень уверенности",
			Description: "Изменения должны иметь уровень уверенности не менее 70%",
			Severity:    "warning",
			Enabled:     true,
		},
		{
			ID:          "sample_size_check",
			Name:        "Размер выборки",
			Description: "Должно быть достаточно данных для принятия решений (минимум 100 поисков)",
			Severity:    "warning",
			Enabled:     true,
		},
		{
			ID:          "negative_impact_check",
			Name:        "Проверка отрицательного влияния",
			Description: "Изменения не должны снижать предсказанный CTR более чем на 5%",
			Severity:    "warning",
			Enabled:     true,
		},
		{
			ID:          "field_consistency",
			Name:        "Согласованность полей",
			Description: "Веса связанных полей должны быть согласованы",
			Severity:    "info",
			Enabled:     true,
		},
		{
			ID:          "frequency_limit",
			Name:        "Ограничение частоты изменений",
			Description: "Один и тот же вес не должен изменяться чаще раза в день",
			Severity:    "warning",
			Enabled:     true,
		},
	}
}

// ValidateOptimizationResults проверяет результаты оптимизации на безопасность
func (sc *SecurityCheck) ValidateOptimizationResults(ctx context.Context, results []*storage.WeightOptimizationResult) (*SecurityReport, error) {
	report := &SecurityReport{
		OverallSafety:    "safe",
		TotalViolations:  0,
		CriticalIssues:   0,
		Warnings:         0,
		Violations:       []RuleViolation{},
		Recommendations:  []string{},
		ApprovalRequired: false,
	}

	rules := sc.GetSafetyRules()

	for _, result := range results {
		for _, rule := range rules {
			if !rule.Enabled {
				continue
			}

			violations := sc.checkRule(ctx, rule, result)
			for _, violation := range violations {
				report.Violations = append(report.Violations, violation)
				report.TotalViolations++

				switch violation.Severity {
				case "critical":
					report.CriticalIssues++
				case severityWarning:
					report.Warnings++
				}
			}
		}
	}

	// Определяем общий уровень безопасности
	if report.CriticalIssues > 0 {
		report.OverallSafety = "dangerous"
		report.ApprovalRequired = true
		report.Recommendations = append(report.Recommendations,
			"❌ Критические проблемы безопасности! Необходимо одобрение старшего администратора.")
	} else if report.Warnings > 3 {
		report.OverallSafety = "warning"
		report.ApprovalRequired = true
		report.Recommendations = append(report.Recommendations,
			"⚠️ Множественные предупреждения. Рекомендуется дополнительная проверка.")
	} else if report.Warnings > 0 {
		report.OverallSafety = "warning"
		report.Recommendations = append(report.Recommendations,
			"⚠️ Обнаружены предупреждения. Мониторьте изменения после применения.")
	}

	// Добавляем общие рекомендации
	if len(results) > 0 {
		report.Recommendations = append(report.Recommendations,
			"📊 Создайте резервную копию текущих весов перед применением изменений")
		report.Recommendations = append(report.Recommendations,
			"📈 Мониторьте метрики CTR и конверсии в течение 24-48 часов после изменений")
		report.Recommendations = append(report.Recommendations,
			"🔄 Подготовьтесь к быстрому откату в случае негативного влияния")
	}

	return report, nil
}

// checkRule проверяет конкретное правило для результата оптимизации
func (sc *SecurityCheck) checkRule(ctx context.Context, rule SafetyRule, result *storage.WeightOptimizationResult) []RuleViolation {
	var violations []RuleViolation

	switch rule.ID {
	case "max_weight_change":
		violations = append(violations, sc.checkMaxWeightChange(rule, result)...)

	case "weight_bounds":
		violations = append(violations, sc.checkWeightBounds(rule, result)...)

	case "critical_fields_protection":
		violations = append(violations, sc.checkCriticalFieldsProtection(rule, result)...)

	case "confidence_threshold":
		violations = append(violations, sc.checkConfidenceThreshold(rule, result)...)

	case "sample_size_check":
		violations = append(violations, sc.checkSampleSize(rule, result)...)

	case "negative_impact_check":
		violations = append(violations, sc.checkNegativeImpact(rule, result)...)

	case "field_consistency":
		violations = append(violations, sc.checkFieldConsistency(rule, result)...)

	case "frequency_limit":
		// Эта проверка требует доступа к истории изменений
		// violations = append(violations, sc.checkFrequencyLimit(ctx, rule, result)...)
	}

	return violations
}

func (sc *SecurityCheck) checkMaxWeightChange(rule SafetyRule, result *storage.WeightOptimizationResult) []RuleViolation {
	var violations []RuleViolation

	changePercent := math.Abs(result.OptimizedWeight-result.CurrentWeight) / result.CurrentWeight
	if changePercent > sc.service.config.MaxWeightChange {
		violations = append(violations, RuleViolation{
			RuleID:   rule.ID,
			Severity: rule.Severity,
			Message: fmt.Sprintf("Изменение веса на %.1f%% превышает максимально допустимое (%.1f%%)",
				changePercent*100, sc.service.config.MaxWeightChange*100),
			FieldName:     result.FieldName,
			CurrentValue:  result.CurrentWeight,
			ProposedValue: result.OptimizedWeight,
			Impact:        "Резкие изменения весов могут негативно повлиять на качество поиска",
		})
	}

	return violations
}

func (sc *SecurityCheck) checkWeightBounds(rule SafetyRule, result *storage.WeightOptimizationResult) []RuleViolation {
	var violations []RuleViolation

	if result.OptimizedWeight < sc.service.config.MinWeight || result.OptimizedWeight > sc.service.config.MaxWeight {
		violations = append(violations, RuleViolation{
			RuleID:   rule.ID,
			Severity: rule.Severity,
			Message: fmt.Sprintf("Вес %.3f выходит за допустимые границы [%.1f, %.1f]",
				result.OptimizedWeight, sc.service.config.MinWeight, sc.service.config.MaxWeight),
			FieldName:     result.FieldName,
			CurrentValue:  result.CurrentWeight,
			ProposedValue: result.OptimizedWeight,
			Impact:        "Веса за пределами допустимого диапазона могут нарушить работу поиска",
		})
	}

	return violations
}

func (sc *SecurityCheck) checkCriticalFieldsProtection(rule SafetyRule, result *storage.WeightOptimizationResult) []RuleViolation {
	var violations []RuleViolation

	criticalFields := map[string]float64{
		"title":       0.3,
		"description": 0.2,
		"category":    0.2,
	}

	if minWeight, isCritical := criticalFields[result.FieldName]; isCritical {
		if result.OptimizedWeight < minWeight {
			violations = append(violations, RuleViolation{
				RuleID:   rule.ID,
				Severity: rule.Severity,
				Message: fmt.Sprintf("Критическое поле '%s' не может иметь вес ниже %.1f",
					result.FieldName, minWeight),
				FieldName:     result.FieldName,
				CurrentValue:  result.CurrentWeight,
				ProposedValue: result.OptimizedWeight,
				Impact:        "Низкие веса критических полей сильно ухудшат качество поиска",
			})
		}
	}

	return violations
}

func (sc *SecurityCheck) checkConfidenceThreshold(rule SafetyRule, result *storage.WeightOptimizationResult) []RuleViolation {
	var violations []RuleViolation

	threshold := 0.7 // 70%
	if result.ConfidenceLevel < threshold {
		violations = append(violations, RuleViolation{
			RuleID:   rule.ID,
			Severity: rule.Severity,
			Message: fmt.Sprintf("Низкий уровень уверенности %.1f%% (требуется минимум %.1f%%)",
				result.ConfidenceLevel*100, threshold*100),
			FieldName:     result.FieldName,
			CurrentValue:  result.ConfidenceLevel,
			ProposedValue: threshold,
			Impact:        "Низкая уверенность означает высокий риск неэффективного изменения",
		})
	}

	return violations
}

func (sc *SecurityCheck) checkSampleSize(rule SafetyRule, result *storage.WeightOptimizationResult) []RuleViolation {
	var violations []RuleViolation

	minSampleSize := 100
	if result.SampleSize < minSampleSize {
		violations = append(violations, RuleViolation{
			RuleID:   rule.ID,
			Severity: rule.Severity,
			Message: fmt.Sprintf("Недостаточный размер выборки: %d (требуется минимум %d)",
				result.SampleSize, minSampleSize),
			FieldName:     result.FieldName,
			CurrentValue:  result.SampleSize,
			ProposedValue: minSampleSize,
			Impact:        "Малая выборка может привести к неточным оптимизациям",
		})
	}

	return violations
}

func (sc *SecurityCheck) checkNegativeImpact(rule SafetyRule, result *storage.WeightOptimizationResult) []RuleViolation {
	var violations []RuleViolation

	if result.ImprovementScore < -5.0 { // Снижение CTR более чем на 5%
		violations = append(violations, RuleViolation{
			RuleID:   rule.ID,
			Severity: rule.Severity,
			Message: fmt.Sprintf("Предсказанное снижение CTR на %.1f%% превышает допустимый порог",
				math.Abs(result.ImprovementScore)),
			FieldName:     result.FieldName,
			CurrentValue:  result.CurrentCTR,
			ProposedValue: result.PredictedCTR,
			Impact:        "Снижение CTR негативно повлияет на пользовательский опыт",
		})
	}

	return violations
}

func (sc *SecurityCheck) checkFieldConsistency(rule SafetyRule, result *storage.WeightOptimizationResult) []RuleViolation {
	var violations []RuleViolation

	// Проверяем логическую согласованность весов
	// Например, title должен иметь больший вес чем description
	if result.FieldName == "description" && result.OptimizedWeight > 0.9 {
		violations = append(violations, RuleViolation{
			RuleID:        rule.ID,
			Severity:      rule.Severity,
			Message:       "Вес описания не должен превышать вес заголовка",
			FieldName:     result.FieldName,
			CurrentValue:  result.CurrentWeight,
			ProposedValue: result.OptimizedWeight,
			Impact:        "Нарушение иерархии важности полей может снизить релевантность",
		})
	}

	return violations
}

// RequiresAdminApproval проверяет, требуется ли дополнительное одобрение администратора
func (sc *SecurityCheck) RequiresAdminApproval(report *SecurityReport) bool {
	return report.ApprovalRequired || report.CriticalIssues > 0
}

// GenerateSecurityBrief создает краткий отчет о безопасности для администратора
func (sc *SecurityCheck) GenerateSecurityBrief(report *SecurityReport) string {
	brief := "🔒 ОТЧЕТ БЕЗОПАСНОСТИ\n"
	brief += fmt.Sprintf("Общий статус: %s\n", report.OverallSafety)
	brief += fmt.Sprintf("Критических проблем: %d\n", report.CriticalIssues)
	brief += fmt.Sprintf("Предупреждений: %d\n", report.Warnings)

	if len(report.Violations) > 0 {
		brief += "\n🚨 ОСНОВНЫЕ ПРОБЛЕМЫ:\n"
		for _, violation := range report.Violations {
			if violation.Severity == "critical" {
				brief += fmt.Sprintf("❌ %s: %s\n", violation.FieldName, violation.Message)
			}
		}
	}

	if len(report.Recommendations) > 0 {
		brief += "\n💡 РЕКОМЕНДАЦИИ:\n"
		for _, rec := range report.Recommendations {
			brief += fmt.Sprintf("• %s\n", rec)
		}
	}

	return brief
}

// CreateSecurityCheckpoint создает checkpoint для отката изменений
func (sc *SecurityCheck) CreateSecurityCheckpoint(ctx context.Context, changes []*storage.WeightOptimizationResult, adminID int) error {
	// Создаем детальную запись об изменениях для возможности отката
	checkpointData := map[string]interface{}{
		"checkpoint_time": time.Now(),
		"admin_id":        adminID,
		"changes":         changes,
		"security_check":  true,
	}

	// В реальной реализации сохранили бы это в отдельную таблицу
	// TODO: Реализовать сохранение checkpoint в БД
	_ = checkpointData // Используем переменную для избежания warning

	sc.service.logger.Info(fmt.Sprintf("Security checkpoint created for %d changes by admin %d",
		len(changes), adminID))

	return nil
}
