package translation_admin

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
)

// GetVersionHistory retrieves version history for a translation
func (s *Service) GetVersionHistory(ctx context.Context, entityType string, entityID int) (*models.VersionHistoryResponse, error) {
	s.logger.Info().
		Str("entity_type", entityType).
		Int("entity_id", entityID).
		Msg("Getting version history")

	// Get all versions for the entity
	versions, err := s.translationRepo.GetVersionsByEntity(ctx, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	if len(versions) == 0 {
		return &models.VersionHistoryResponse{
			TranslationID:  0,
			CurrentVersion: 0,
			Versions:       []models.TranslationVersion{},
			TotalVersions:  0,
		}, nil
	}

	// Get the maximum version number across all translations
	var maxVersion int
	for _, version := range versions {
		if version.Version > maxVersion {
			maxVersion = version.Version
		}
	}

	// Return all versions for the entity
	response := &models.VersionHistoryResponse{
		TranslationID:  0, // 0 indicates multiple translations
		CurrentVersion: maxVersion,
		Versions:       versions, // Return all versions for all fields/languages
		TotalVersions:  len(versions),
	}

	return response, nil
}

// GetVersionDiff compares two translation versions
func (s *Service) GetVersionDiff(ctx context.Context, versionID1, versionID2 int) (*models.VersionDiff, error) {
	s.logger.Info().
		Int("version1", versionID1).
		Int("version2", versionID2).
		Msg("Getting version diff")

	diff, err := s.translationRepo.GetVersionDiff(ctx, versionID1, versionID2)
	if err != nil {
		return nil, fmt.Errorf("failed to get version diff: %w", err)
	}

	return diff, nil
}

// RollbackVersion rolls back a translation to a previous version
func (s *Service) RollbackVersion(ctx context.Context, req *models.RollbackRequest, userID int) error {
	s.logger.Info().
		Int("translation_id", req.TranslationID).
		Int("version_id", req.VersionID).
		Int("user_id", userID).
		Str("comment", req.Comment).
		Msg("Rolling back version")

	// Perform rollback
	err := s.translationRepo.RollbackToVersion(ctx, req.TranslationID, req.VersionID, userID)
	if err != nil {
		return fmt.Errorf("failed to rollback version: %w", err)
	}

	// Log additional audit entry with comment if provided
	if req.Comment != "" {
		auditLog := &models.TranslationAuditLog{
			UserID:     &userID,
			Action:     "rollback_comment",
			EntityType: strPtr("translation"),
			EntityID:   &req.TranslationID,
			NewValue:   &req.Comment,
		}

		if err := s.auditRepo.LogAction(ctx, auditLog); err != nil {
			s.logger.Error().Err(err).Msg("Failed to log rollback comment")
		}
	}

	return nil
}

// GetAuditLogs retrieves audit logs with filters
func (s *Service) GetAuditLogs(ctx context.Context, filters map[string]interface{}) ([]models.TranslationAuditLog, error) {
	s.logger.Info().Msg("Getting audit logs")

	// For now, use the basic method from auditRepo
	// In a full implementation, you'd pass filters to the repository method
	limit, ok := filters["limit"].(int)
	if !ok {
		limit = 100
	}

	logs, err := s.auditRepo.GetRecentLogs(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, nil
}

// GetAuditStatistics retrieves audit statistics
func (s *Service) GetAuditStatistics(ctx context.Context) (*models.AuditStatistics, error) {
	s.logger.Info().Msg("Getting audit statistics")

	// Get recent logs to calculate statistics
	logs, err := s.auditRepo.GetRecentLogs(ctx, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for statistics: %w", err)
	}

	stats := &models.AuditStatistics{
		TotalActions:  len(logs),
		ActionsByType: make(map[string]int),
		ActionsByUser: make(map[int]int),
		RecentActions: []models.TranslationAuditLog{},
	}

	// Calculate statistics
	for _, log := range logs {
		stats.ActionsByType[log.Action]++
		if log.UserID != nil {
			stats.ActionsByUser[*log.UserID]++
		}
	}

	// Get recent actions (last 10)
	if len(logs) > 10 {
		stats.RecentActions = logs[:10]
	} else {
		stats.RecentActions = logs
	}

	return stats, nil
}
