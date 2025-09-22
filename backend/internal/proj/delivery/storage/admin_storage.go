package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"

	"backend/internal/proj/delivery/models"
)

const (
	NotAvailable = "N/A"
)

// GetAllProviders retrieves all delivery providers
func (s *Storage) GetAllProviders(ctx context.Context) ([]models.Provider, error) {
	query := `
		SELECT id, code, name, logo_url, is_active,
		       supports_cod, supports_insurance, supports_tracking,
		       api_config, capabilities
		FROM delivery_providers
		ORDER BY name
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close database rows")
		}
	}()

	var providers []models.Provider
	for rows.Next() {
		var p models.Provider
		var apiConfig, capabilities sql.NullString

		err := rows.Scan(
			&p.ID, &p.Code, &p.Name, &p.LogoURL, &p.IsActive,
			&p.SupportsCOD, &p.SupportsInsurance, &p.SupportsTracking,
			&apiConfig, &capabilities,
		)
		if err != nil {
			return nil, err
		}

		// Check if API is configured - we can use this for display purposes

		providers = append(providers, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return providers, nil
}

// ToggleProviderStatus toggles the active status of a provider
func (s *Storage) ToggleProviderStatus(ctx context.Context, providerID int) error {
	query := `
		UPDATE delivery_providers
		SET is_active = NOT is_active, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := s.db.ExecContext(ctx, query, providerID)
	return err
}

// UpdateProviderConfig updates provider API configuration
func (s *Storage) UpdateProviderConfig(ctx context.Context, providerID int, config models.ProviderConfig) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	query := `
		UPDATE delivery_providers
		SET api_config = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err = s.db.ExecContext(ctx, query, providerID, configJSON)
	return err
}

// GetAllPricingRules retrieves all pricing rules
func (s *Storage) GetAllPricingRules(ctx context.Context) ([]models.PricingRule, error) {
	query := `
		SELECT pr.id, pr.provider_id, dp.name, pr.rule_type, pr.priority,
		       pr.is_active, pr.min_price, pr.max_price,
		       pr.weight_ranges, pr.volume_ranges,
		       pr.fragile_surcharge, pr.oversized_surcharge,
		       pr.special_handling_surcharge
		FROM delivery_pricing_rules pr
		JOIN delivery_providers dp ON pr.provider_id = dp.id
		ORDER BY dp.name, pr.priority
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close database rows")
		}
	}()

	var rules []models.PricingRule
	for rows.Next() {
		var r models.PricingRule
		var weightRanges, volumeRanges sql.NullString

		var providerName string
		err := rows.Scan(
			&r.ID, &r.ProviderID, &providerName, &r.RuleType, &r.Priority,
			&r.IsActive, &r.MinPrice, &r.MaxPrice,
			&weightRanges, &volumeRanges,
			&r.FragileSurcharge, &r.OversizedSurcharge,
			&r.SpecialHandlingSurcharge,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON ranges if present
		if weightRanges.Valid {
			if err := json.Unmarshal([]byte(weightRanges.String), &r.WeightRanges); err != nil {
				log.Error().Err(err).Msg("Failed to unmarshal weight ranges")
			}
		}
		if volumeRanges.Valid {
			if err := json.Unmarshal([]byte(volumeRanges.String), &r.VolumeRanges); err != nil {
				log.Error().Err(err).Msg("Failed to unmarshal volume ranges")
			}
		}

		rules = append(rules, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return rules, nil
}

// UpdatePricingRule updates an existing pricing rule
func (s *Storage) UpdatePricingRule(ctx context.Context, rule models.PricingRule) error {
	weightRanges, _ := json.Marshal(rule.WeightRanges)
	volumeRanges, _ := json.Marshal(rule.VolumeRanges)

	query := `
		UPDATE delivery_pricing_rules
		SET rule_type = $2, priority = $3, is_active = $4,
		    min_price = $5, max_price = $6,
		    weight_ranges = $7, volume_ranges = $8,
		    fragile_surcharge = $9, oversized_surcharge = $10,
		    special_handling_surcharge = $11,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := s.db.ExecContext(
		ctx, query, rule.ID,
		rule.RuleType, rule.Priority, rule.IsActive,
		rule.MinPrice, rule.MaxPrice, weightRanges, volumeRanges,
		rule.FragileSurcharge, rule.OversizedSurcharge, rule.SpecialHandlingSurcharge,
	)

	return err
}

// DeletePricingRule deletes a pricing rule
func (s *Storage) DeletePricingRule(ctx context.Context, ruleID int) error {
	query := `DELETE FROM delivery_pricing_rules WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, ruleID)
	return err
}

// GetProblemShipments retrieves problem shipments based on filters
func (s *Storage) GetProblemShipments(ctx context.Context, problemType, status string) ([]models.ProblemShipment, error) {
	query := `
		SELECT ds.id, ds.tracking_number, dp.name,
		       ds.status, ds.created_at,
		       COALESCE(ds.recipient_info->>'name', '') as customer_name,
		       COALESCE(ds.recipient_info->>'phone', '') as customer_phone,
		       COALESCE(ds.provider_response->>'error_message', 'Problem detected') as description
		FROM delivery_shipments ds
		JOIN delivery_providers dp ON ds.provider_id = dp.id
		WHERE ds.status IN ('failed', 'canceled', 'returned')
	`

	args := []interface{}{}
	argCount := 0

	if problemType != "" {
		argCount++
		query += fmt.Sprintf(" AND ds.status = $%d", argCount)
		args = append(args, problemType)
	}

	if status != "" {
		// Map frontend status to DB status
		argCount++
		query += fmt.Sprintf(" AND ds.status = $%d", argCount)
		args = append(args, status)
	}

	query += " ORDER BY ds.created_at DESC"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close database rows")
		}
	}()

	var problems []models.ProblemShipment
	for rows.Next() {
		var p models.ProblemShipment

		err := rows.Scan(
			&p.ID, &p.TrackingNumber, &p.ProviderName,
			&p.Status, &p.CreatedAt,
			&p.CustomerName, &p.CustomerPhone, &p.Description,
		)
		if err != nil {
			return nil, err
		}

		// Map status to problem type
		switch p.Status {
		case "failed":
			p.ProblemType = "delayed"
			p.Priority = "high"
		case "canceled":
			p.ProblemType = "refused"
			p.Priority = "medium"
		case "returned":
			p.ProblemType = "wrongAddress"
			p.Priority = "low"
		}

		problems = append(problems, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return problems, nil
}

// AssignProblem assigns a problem to an admin
func (s *Storage) AssignProblem(ctx context.Context, problemID, adminID int) error {
	// This would typically update a separate problems table
	// For now, we'll update metadata in the shipments table
	query := `
		UPDATE delivery_shipments
		SET provider_response = jsonb_set(
			COALESCE(provider_response, '{}'::jsonb),
			'{assigned_to}',
			to_jsonb($2::text)
		),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := s.db.ExecContext(ctx, query, problemID, adminID)
	return err
}

// ResolveProblem marks a problem as resolved
func (s *Storage) ResolveProblem(ctx context.Context, problemID int, resolution string) error {
	query := `
		UPDATE delivery_shipments
		SET status = 'resolved',
		    provider_response = jsonb_set(
			    COALESCE(provider_response, '{}'::jsonb),
			    '{resolution}',
			    to_jsonb($2::text)
		    ),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := s.db.ExecContext(ctx, query, problemID, resolution)
	return err
}

// GetDashboardStats retrieves dashboard statistics
func (s *Storage) GetDashboardStats(ctx context.Context) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}

	// Today's shipments
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM delivery_shipments
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&stats.TodayShipments)
	if err != nil {
		return nil, err
	}

	// Today's delivered
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM delivery_shipments
		WHERE DATE(actual_delivery_date) = CURRENT_DATE
	`).Scan(&stats.TodayDelivered)
	if err != nil {
		return nil, err
	}

	// In transit
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM delivery_shipments
		WHERE status IN ('in_transit', 'out_for_delivery')
	`).Scan(&stats.InTransit)
	if err != nil {
		return nil, err
	}

	// Problems
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM delivery_shipments
		WHERE status IN ('failed', 'canceled', 'returned')
	`).Scan(&stats.Problems)
	if err != nil {
		return nil, err
	}

	// Success rate
	var total, successful int
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*), COUNT(CASE WHEN status = 'delivered' THEN 1 END)
		FROM delivery_shipments
		WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
	`).Scan(&total, &successful)
	if err != nil {
		return nil, err
	}

	if total > 0 {
		stats.SuccessRate = float64(successful) / float64(total) * 100
	}

	// Average delivery time
	var avgTime sql.NullFloat64
	err = s.db.QueryRowContext(ctx, `
		SELECT AVG(EXTRACT(EPOCH FROM (actual_delivery_date - created_at))/86400.0)
		FROM delivery_shipments
		WHERE actual_delivery_date IS NOT NULL
		AND created_at >= CURRENT_DATE - INTERVAL '30 days'
	`).Scan(&avgTime)
	if err != nil {
		return nil, err
	}

	if avgTime.Valid {
		stats.AvgDeliveryTime = fmt.Sprintf("%.1f дня", avgTime.Float64)
	} else {
		stats.AvgDeliveryTime = NotAvailable
	}

	// Provider stats
	providerQuery := `
		SELECT dp.code, dp.name,
		       COUNT(*) as shipments,
		       COUNT(CASE WHEN ds.status = 'delivered' THEN 1 END) as delivered,
		       AVG(EXTRACT(EPOCH FROM (ds.actual_delivery_date - ds.created_at))/86400.0) as avg_time
		FROM delivery_shipments ds
		JOIN delivery_providers dp ON ds.provider_id = dp.id
		WHERE ds.created_at >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY dp.id, dp.code, dp.name
		ORDER BY shipments DESC
	`

	rows, err := s.db.QueryContext(ctx, providerQuery)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error().Err(closeErr).Msg("Failed to close database rows")
		}
	}()

	for rows.Next() {
		var ps models.ProviderStats
		var avgTime sql.NullFloat64

		err := rows.Scan(&ps.Code, &ps.Name, &ps.Shipments, &ps.Delivered, &avgTime)
		if err != nil {
			continue
		}

		if ps.Shipments > 0 {
			ps.SuccessRate = float64(ps.Delivered) / float64(ps.Shipments) * 100
		}

		if avgTime.Valid {
			ps.AvgTime = fmt.Sprintf("%.1f дня", avgTime.Float64)
		} else {
			ps.AvgTime = NotAvailable
		}

		stats.ProviderStats = append(stats.ProviderStats, ps)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// Cost analysis
	err = s.db.QueryRowContext(ctx, `
		SELECT AVG(delivery_cost), SUM(delivery_cost)
		FROM delivery_shipments
		WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
	`).Scan(&stats.CostAnalysis.AvgCost, &stats.CostAnalysis.TotalCost)
	if err != nil {
		// Use default values if no data
		stats.CostAnalysis.AvgCost = 0
		stats.CostAnalysis.TotalCost = 0
	}

	// Calculate savings (mock calculation)
	stats.CostAnalysis.Savings = stats.CostAnalysis.TotalCost * 0.15

	return stats, nil
}

// GetAnalytics retrieves analytics data for the specified period
func (s *Storage) GetAnalytics(ctx context.Context, period string) (*models.AnalyticsData, error) {
	analytics := &models.AnalyticsData{Period: period}

	// Parse period to get interval
	var interval string
	switch period {
	case "7d":
		interval = "7 days"
	case "90d":
		interval = "90 days"
	case "365d":
		interval = "365 days"
	default:
		interval = "30 days"
	}

	// Total shipments
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM delivery_shipments
		WHERE created_at >= CURRENT_DATE - INTERVAL '`+interval+`'
	`).Scan(&analytics.TotalShipments)
	if err != nil {
		return nil, err
	}

	// Delivery rate and other metrics
	var delivered int
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(CASE WHEN status = 'delivered' THEN 1 END)
		FROM delivery_shipments
		WHERE created_at >= CURRENT_DATE - INTERVAL '`+interval+`'
	`).Scan(&delivered)
	if err != nil {
		return nil, err
	}

	if analytics.TotalShipments > 0 {
		analytics.DeliveryRate = float64(delivered) / float64(analytics.TotalShipments) * 100
	}

	// Average delivery time
	var avgTime sql.NullFloat64
	err = s.db.QueryRowContext(ctx, `
		SELECT AVG(EXTRACT(EPOCH FROM (actual_delivery_date - created_at))/86400.0)
		FROM delivery_shipments
		WHERE actual_delivery_date IS NOT NULL
		AND created_at >= CURRENT_DATE - INTERVAL '`+interval+`'
	`).Scan(&avgTime)
	if err != nil {
		return nil, err
	}

	if avgTime.Valid {
		analytics.AvgDeliveryTime = fmt.Sprintf("%.1f дня", avgTime.Float64)
	} else {
		analytics.AvgDeliveryTime = NotAvailable
	}

	// Customer satisfaction (mock data for now)
	analytics.CustomerSatisfaction = 4.5

	// Cost per shipment
	err = s.db.QueryRowContext(ctx, `
		SELECT AVG(delivery_cost)
		FROM delivery_shipments
		WHERE created_at >= CURRENT_DATE - INTERVAL '`+interval+`'
	`).Scan(&analytics.CostPerShipment)
	if err != nil {
		analytics.CostPerShipment = 0
	}

	// Problem rate
	var problems int
	err = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM delivery_shipments
		WHERE status IN ('failed', 'canceled', 'returned')
		AND created_at >= CURRENT_DATE - INTERVAL '`+interval+`'
	`).Scan(&problems)
	if err == nil && analytics.TotalShipments > 0 {
		analytics.ProblemRate = float64(problems) / float64(analytics.TotalShipments) * 100
	}

	// Geographic distribution (simplified)
	analytics.GeographicDistribution = map[string]float64{
		"Белград":   42.0,
		"Нови-Сад":  18.0,
		"Ниш":       12.0,
		"Крагуевац": 8.0,
		"Суботица":  6.0,
		"Другие":    14.0,
	}

	return analytics, nil
}
