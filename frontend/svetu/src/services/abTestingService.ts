import { v4 as uuidv4 } from 'uuid';
import configManager from '@/config';

interface Experiment {
  id: string;
  name: string;
  description: string;
  variants: Variant[];
  targetAudience?: TargetAudience;
  trafficAllocation: number; // 0-100 percentage
  startDate: Date;
  endDate?: Date;
  status: 'draft' | 'running' | 'paused' | 'completed';
  metrics: string[];
  winnerVariant?: string;
  statisticalSignificance?: number;
}

interface Variant {
  id: string;
  name: string;
  description: string;
  weight: number; // Distribution weight
  config: Record<string, any>;
  metrics?: VariantMetrics;
}

interface VariantMetrics {
  impressions: number;
  conversions: number;
  conversionRate: number;
  revenue?: number;
  averageOrderValue?: number;
  bounceRate?: number;
  engagementTime?: number;
  customMetrics?: Record<string, number>;
}

interface TargetAudience {
  segments?: string[];
  countries?: string[];
  devices?: ('mobile' | 'tablet' | 'desktop')[];
  browsers?: string[];
  languages?: string[];
  newUsers?: boolean;
  customRules?: Array<{
    field: string;
    operator: 'equals' | 'contains' | 'greater' | 'less';
    value: any;
  }>;
}

interface UserContext {
  userId?: string;
  sessionId: string;
  device: string;
  browser: string;
  country: string;
  language: string;
  isNewUser: boolean;
  customProperties?: Record<string, any>;
}

interface ExperimentAssignment {
  experimentId: string;
  variantId: string;
  assignedAt: Date;
  sticky: boolean;
}

interface ConversionEvent {
  experimentId: string;
  variantId: string;
  eventName: string;
  value?: number;
  metadata?: Record<string, any>;
  timestamp: Date;
}

/**
 * A/B Testing Service
 * Manages experiments, variant assignments, and conversion tracking
 */
class ABTestingService {
  private static instance: ABTestingService;
  private experiments: Map<string, Experiment> = new Map();
  private userAssignments: Map<string, ExperimentAssignment[]> = new Map();
  private conversionEvents: ConversionEvent[] = [];
  private featureFlags: Map<string, boolean | string | number> = new Map();
  private analyticsQueue: any[] = [];
  private config = {
    cookieName: 'ab_experiments',
    cookieExpiry: 90, // days
    minSampleSize: 100,
    confidenceLevel: 0.95,
    apiEndpoint: '/api/v1/experiments',
    enableRemoteConfig: true,
    enableAnalytics: true,
  };

  private constructor() {
    this.initialize();
  }

  static getInstance(): ABTestingService {
    if (!this.instance) {
      this.instance = new ABTestingService();
    }
    return this.instance;
  }

  /**
   * Initialize the service
   */
  private async initialize() {
    // Load experiments from remote config
    if (this.config.enableRemoteConfig) {
      await this.loadRemoteExperiments();
    }

    // Load user assignments from cookie
    this.loadUserAssignments();

    // Start analytics batch processor
    if (this.config.enableAnalytics) {
      this.startAnalyticsBatcher();
    }

    // Listen for visibility changes to pause/resume
    document.addEventListener('visibilitychange', () => {
      if (document.hidden) {
        this.flushAnalytics();
      }
    });
  }

  /**
   * Load experiments from remote config
   */
  private async loadRemoteExperiments() {
    try {
      const apiUrl = configManager.get('api.url');
      const response = await fetch(`${apiUrl}${this.config.apiEndpoint}/active`);
      const experiments = await response.json();

      experiments.forEach((exp: Experiment) => {
        this.experiments.set(exp.id, exp);
      });
    } catch (error) {
      console.error('Failed to load experiments:', error);
      // Fall back to local experiments
      this.loadLocalExperiments();
    }
  }

  /**
   * Load local experiment configurations
   */
  private loadLocalExperiments() {
    // Hardcoded experiments for fallback
    const localExperiments: Experiment[] = [
      {
        id: 'unified-attributes-ui',
        name: 'Unified Attributes UI Test',
        description: 'Testing new unified attributes interface',
        variants: [
          {
            id: 'control',
            name: 'Control',
            description: 'Current attributes UI',
            weight: 50,
            config: { useUnifiedUI: false },
          },
          {
            id: 'variant-a',
            name: 'Unified UI',
            description: 'New unified attributes UI',
            weight: 50,
            config: { useUnifiedUI: true },
          },
        ],
        trafficAllocation: 100,
        startDate: new Date('2025-09-01'),
        status: 'running',
        metrics: ['conversion_rate', 'engagement_time', 'bounce_rate'],
      },
    ];

    localExperiments.forEach((exp) => {
      this.experiments.set(exp.id, exp);
    });
  }

  /**
   * Load user assignments from cookie
   */
  private loadUserAssignments() {
    const cookie = this.getCookie(this.config.cookieName);
    if (cookie) {
      try {
        const assignments = JSON.parse(decodeURIComponent(cookie));
        Object.entries(assignments).forEach(([userId, userAssignments]) => {
          this.userAssignments.set(
            userId,
            userAssignments as ExperimentAssignment[]
          );
        });
      } catch (error) {
        console.error('Failed to parse assignments cookie:', error);
      }
    }
  }

  /**
   * Save user assignments to cookie
   */
  private saveUserAssignments() {
    const assignments: Record<string, ExperimentAssignment[]> = {};
    this.userAssignments.forEach((value, key) => {
      assignments[key] = value;
    });

    this.setCookie(
      this.config.cookieName,
      encodeURIComponent(JSON.stringify(assignments)),
      this.config.cookieExpiry
    );
  }

  /**
   * Get variant for a user in an experiment
   */
  getVariant(experimentId: string, userContext: UserContext): Variant | null {
    const experiment = this.experiments.get(experimentId);
    if (!experiment || experiment.status !== 'running') {
      return null;
    }

    // Check if experiment targets this user
    if (!this.isUserEligible(experiment, userContext)) {
      return null;
    }

    // Check traffic allocation
    if (!this.isInTrafficAllocation(experiment, userContext)) {
      return null;
    }

    // Get or create assignment
    const userId = userContext.userId || userContext.sessionId;
    const assignment = this.getOrCreateAssignment(experiment, userId);

    // Find and return the variant
    return (
      experiment.variants.find((v) => v.id === assignment.variantId) || null
    );
  }

  /**
   * Check if user is eligible for experiment
   */
  private isUserEligible(
    experiment: Experiment,
    userContext: UserContext
  ): boolean {
    const audience = experiment.targetAudience;
    if (!audience) return true;

    // Check device
    if (
      audience.devices &&
      !audience.devices.includes(userContext.device as any)
    ) {
      return false;
    }

    // Check country
    if (
      audience.countries &&
      !audience.countries.includes(userContext.country)
    ) {
      return false;
    }

    // Check language
    if (
      audience.languages &&
      !audience.languages.includes(userContext.language)
    ) {
      return false;
    }

    // Check new user status
    if (
      audience.newUsers !== undefined &&
      audience.newUsers !== userContext.isNewUser
    ) {
      return false;
    }

    // Check custom rules
    if (audience.customRules) {
      for (const rule of audience.customRules) {
        const value = userContext.customProperties?.[rule.field];
        if (!this.evaluateRule(rule, value)) {
          return false;
        }
      }
    }

    return true;
  }

  /**
   * Evaluate custom targeting rule
   */
  private evaluateRule(rule: any, value: any): boolean {
    switch (rule.operator) {
      case 'equals':
        return value === rule.value;
      case 'contains':
        return String(value).includes(String(rule.value));
      case 'greater':
        return Number(value) > Number(rule.value);
      case 'less':
        return Number(value) < Number(rule.value);
      default:
        return false;
    }
  }

  /**
   * Check if user is in traffic allocation
   */
  private isInTrafficAllocation(
    experiment: Experiment,
    userContext: UserContext
  ): boolean {
    const userId = userContext.userId || userContext.sessionId;
    const hash = this.hashString(userId + experiment.id);
    const bucket = (hash % 100) + 1;
    return bucket <= experiment.trafficAllocation;
  }

  /**
   * Get or create assignment for user
   */
  private getOrCreateAssignment(
    experiment: Experiment,
    userId: string
  ): ExperimentAssignment {
    // Check existing assignments
    const userAssignments = this.userAssignments.get(userId) || [];
    const existing = userAssignments.find(
      (a) => a.experimentId === experiment.id
    );

    if (existing) {
      return existing;
    }

    // Create new assignment
    const variantId = this.selectVariant(experiment, userId);
    const assignment: ExperimentAssignment = {
      experimentId: experiment.id,
      variantId,
      assignedAt: new Date(),
      sticky: true,
    };

    // Save assignment
    userAssignments.push(assignment);
    this.userAssignments.set(userId, userAssignments);
    this.saveUserAssignments();

    // Track assignment event
    this.trackEvent('experiment_assigned', {
      experimentId: experiment.id,
      variantId,
      userId,
    });

    return assignment;
  }

  /**
   * Select variant based on weights
   */
  private selectVariant(experiment: Experiment, userId: string): string {
    const hash = this.hashString(userId + experiment.id + 'variant');
    const totalWeight = experiment.variants.reduce(
      (sum, v) => sum + v.weight,
      0
    );
    const random = (hash % totalWeight) + 1;

    let cumulative = 0;
    for (const variant of experiment.variants) {
      cumulative += variant.weight;
      if (random <= cumulative) {
        return variant.id;
      }
    }

    return experiment.variants[0].id;
  }

  /**
   * Track conversion event
   */
  trackConversion(
    experimentId: string,
    eventName: string,
    value?: number,
    metadata?: Record<string, any>
  ) {
    const userId = this.getCurrentUserId();
    const userAssignments = this.userAssignments.get(userId) || [];
    const assignment = userAssignments.find(
      (a) => a.experimentId === experimentId
    );

    if (!assignment) {
      return; // User not in experiment
    }

    const event: ConversionEvent = {
      experimentId,
      variantId: assignment.variantId,
      eventName,
      value,
      metadata,
      timestamp: new Date(),
    };

    this.conversionEvents.push(event);

    // Queue for analytics
    this.queueAnalytics({
      type: 'conversion',
      ...event,
    });

    // Update variant metrics
    this.updateVariantMetrics(
      experimentId,
      assignment.variantId,
      eventName,
      value
    );
  }

  /**
   * Update variant metrics
   */
  private updateVariantMetrics(
    experimentId: string,
    variantId: string,
    eventName: string,
    value?: number
  ) {
    const experiment = this.experiments.get(experimentId);
    if (!experiment) return;

    const variant = experiment.variants.find((v) => v.id === variantId);
    if (!variant) return;

    if (!variant.metrics) {
      variant.metrics = {
        impressions: 0,
        conversions: 0,
        conversionRate: 0,
      };
    }

    // Update metrics
    if (eventName === 'impression') {
      variant.metrics.impressions++;
    } else if (eventName === 'conversion') {
      variant.metrics.conversions++;
      if (value) {
        variant.metrics.revenue = (variant.metrics.revenue || 0) + value;
      }
    }

    // Calculate conversion rate
    if (variant.metrics.impressions > 0) {
      variant.metrics.conversionRate =
        (variant.metrics.conversions / variant.metrics.impressions) * 100;
    }

    // Check for statistical significance
    this.checkStatisticalSignificance(experiment);
  }

  /**
   * Check statistical significance
   */
  private checkStatisticalSignificance(experiment: Experiment) {
    if (experiment.variants.length !== 2) return; // Only for A/B tests

    const control = experiment.variants[0].metrics;
    const variant = experiment.variants[1].metrics;

    if (!control || !variant) return;

    // Minimum sample size check
    if (
      control.impressions < this.config.minSampleSize ||
      variant.impressions < this.config.minSampleSize
    ) {
      return;
    }

    // Calculate Z-score
    const p1 = control.conversionRate / 100;
    const p2 = variant.conversionRate / 100;
    const n1 = control.impressions;
    const n2 = variant.impressions;

    const pooledP = (control.conversions + variant.conversions) / (n1 + n2);
    const se = Math.sqrt(pooledP * (1 - pooledP) * (1 / n1 + 1 / n2));
    const zScore = Math.abs(p1 - p2) / se;

    // Check significance (95% confidence = 1.96 Z-score)
    const isSignificant = zScore >= 1.96;

    if (isSignificant) {
      experiment.statisticalSignificance = this.zScoreToConfidence(zScore);
      experiment.winnerVariant =
        p2 > p1 ? experiment.variants[1].id : experiment.variants[0].id;

      // Notify about winner
      this.notifyExperimentComplete(experiment);
    }
  }

  /**
   * Convert Z-score to confidence percentage
   */
  private zScoreToConfidence(zScore: number): number {
    // Simplified conversion
    if (zScore >= 3.29) return 99.9;
    if (zScore >= 2.58) return 99;
    if (zScore >= 1.96) return 95;
    if (zScore >= 1.64) return 90;
    return Math.min(zScore * 30, 89);
  }

  /**
   * Notify about experiment completion
   */
  private notifyExperimentComplete(experiment: Experiment) {
    console.log(
      `Experiment ${experiment.name} has a winner: ${experiment.winnerVariant}`
    );

    // Send to analytics
    this.trackEvent('experiment_complete', {
      experimentId: experiment.id,
      winnerVariant: experiment.winnerVariant,
      significance: experiment.statisticalSignificance,
    });

    // Update status
    experiment.status = 'completed';

    // Save to remote
    this.saveExperimentResults(experiment);
  }

  /**
   * Save experiment results
   */
  private async saveExperimentResults(experiment: Experiment) {
    if (!this.config.enableRemoteConfig) return;

    try {
      const apiUrl = configManager.get('api.url');
      await fetch(`${apiUrl}${this.config.apiEndpoint}/${experiment.id}/complete`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          winnerVariant: experiment.winnerVariant,
          significance: experiment.statisticalSignificance,
          variants: experiment.variants.map((v) => ({
            id: v.id,
            metrics: v.metrics,
          })),
        }),
      });
    } catch (error) {
      console.error('Failed to save experiment results:', error);
    }
  }

  /**
   * Feature flags management
   */
  isFeatureEnabled(flagName: string, defaultValue: boolean = false): boolean {
    const value = this.featureFlags.get(flagName);
    if (value === undefined) {
      return defaultValue;
    }
    return Boolean(value);
  }

  getFeatureValue(flagName: string, defaultValue: any = null): any {
    return this.featureFlags.get(flagName) ?? defaultValue;
  }

  setFeatureFlag(flagName: string, value: boolean | string | number) {
    this.featureFlags.set(flagName, value);
  }

  /**
   * Analytics integration
   */
  private queueAnalytics(event: any) {
    if (!this.config.enableAnalytics) return;

    this.analyticsQueue.push({
      ...event,
      timestamp: new Date().toISOString(),
      sessionId: this.getSessionId(),
    });

    // Batch send every 10 events or 5 seconds
    if (this.analyticsQueue.length >= 10) {
      this.flushAnalytics();
    }
  }

  private startAnalyticsBatcher() {
    setInterval(() => {
      if (this.analyticsQueue.length > 0) {
        this.flushAnalytics();
      }
    }, 5000);
  }

  private async flushAnalytics() {
    if (this.analyticsQueue.length === 0) return;

    const events = [...this.analyticsQueue];
    this.analyticsQueue = [];

    try {
      const apiUrl = configManager.get('api.url');
      await fetch(`${apiUrl}/api/v1/analytics/events`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ events }),
      });
    } catch (error) {
      console.error('Failed to send analytics:', error);
      // Re-queue events
      this.analyticsQueue.unshift(...events);
    }
  }

  private trackEvent(eventName: string, data: any) {
    this.queueAnalytics({
      type: 'event',
      name: eventName,
      data,
    });
  }

  /**
   * Utility methods
   */
  private hashString(str: string): number {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      const char = str.charCodeAt(i);
      hash = (hash << 5) - hash + char;
      hash = hash & hash;
    }
    return Math.abs(hash);
  }

  private getCurrentUserId(): string {
    // Try to get from auth context
    const authUser = this.getAuthenticatedUser();
    if (authUser?.id) {
      return authUser.id;
    }

    // Fall back to session ID
    return this.getSessionId();
  }

  private getSessionId(): string {
    let sessionId = sessionStorage.getItem('ab_session_id');
    if (!sessionId) {
      sessionId = uuidv4();
      sessionStorage.setItem('ab_session_id', sessionId);
    }
    return sessionId;
  }

  private getAuthenticatedUser(): any {
    // This would integrate with your auth system
    return null;
  }

  private getCookie(name: string): string | null {
    const match = document.cookie.match(
      new RegExp('(^| )' + name + '=([^;]+)')
    );
    return match ? match[2] : null;
  }

  private setCookie(name: string, value: string, days: number) {
    const expires = new Date();
    expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);
    document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/`;
  }

  /**
   * Public API
   */
  getAllExperiments(): Experiment[] {
    return Array.from(this.experiments.values());
  }

  getExperimentResults(experimentId: string): any {
    const experiment = this.experiments.get(experimentId);
    if (!experiment) return null;

    return {
      experiment: experiment.name,
      status: experiment.status,
      variants: experiment.variants.map((v) => ({
        name: v.name,
        metrics: v.metrics,
      })),
      winner: experiment.winnerVariant,
      significance: experiment.statisticalSignificance,
    };
  }

  getUserExperiments(): ExperimentAssignment[] {
    const userId = this.getCurrentUserId();
    return this.userAssignments.get(userId) || [];
  }

  resetUserAssignments() {
    const userId = this.getCurrentUserId();
    this.userAssignments.delete(userId);
    this.saveUserAssignments();
  }
}

// Export singleton instance
export const abTesting = ABTestingService.getInstance();

// Export types
export type { Experiment, Variant, UserContext, ConversionEvent };
