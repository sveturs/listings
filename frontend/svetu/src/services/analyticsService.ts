interface AnalyticsEvent {
  name: string;
  category: string;
  action?: string;
  label?: string;
  value?: number;
  metadata?: Record<string, any>;
  timestamp: Date;
  sessionId: string;
  userId?: string;
}

interface PageView {
  path: string;
  title: string;
  referrer: string;
  duration?: number;
  exitRate?: number;
  bounceRate?: number;
  timestamp: Date;
}

interface UserProfile {
  userId: string;
  traits: Record<string, any>;
  segments: string[];
  createdAt: Date;
  lastSeen: Date;
}

interface ConversionGoal {
  id: string;
  name: string;
  type: 'event' | 'pageview' | 'duration' | 'custom';
  conditions: GoalCondition[];
  value?: number;
  isCompleted: boolean;
}

interface GoalCondition {
  field: string;
  operator: 'equals' | 'contains' | 'greater' | 'less' | 'regex';
  value: any;
}

interface FunnelStep {
  name: string;
  event: string;
  completed: boolean;
  dropoffRate?: number;
  averageTime?: number;
}

interface Funnel {
  id: string;
  name: string;
  steps: FunnelStep[];
  conversionRate: number;
  totalTime: number;
}

interface HeatmapData {
  element: string;
  x: number;
  y: number;
  clicks: number;
  hovers: number;
  scrollDepth: number;
}

interface SessionRecording {
  sessionId: string;
  userId?: string;
  startTime: Date;
  endTime?: Date;
  events: RecordedEvent[];
  metadata: Record<string, any>;
}

interface RecordedEvent {
  type: 'click' | 'scroll' | 'input' | 'navigation' | 'error';
  target?: string;
  data: any;
  timestamp: number;
}

/**
 * Comprehensive Analytics Service
 */
class AnalyticsService {
  private static instance: AnalyticsService;
  private events: AnalyticsEvent[] = [];
  private pageViews: PageView[] = [];
  private goals: Map<string, ConversionGoal> = new Map();
  private funnels: Map<string, Funnel> = new Map();
  private heatmapData: HeatmapData[] = [];
  private sessionRecording: SessionRecording | null = null;
  private userProfile: UserProfile | null = null;
  private sessionId: string;
  private pageStartTime: number = Date.now();
  private isRecording: boolean = false;
  private eventQueue: any[] = [];
  private config = {
    apiEndpoint: '/api/v1/analytics',
    batchSize: 20,
    batchInterval: 5000,
    enableHeatmap: true,
    enableRecording: false,
    enableAutoTrack: true,
    cookiePrefix: 'analytics_',
    sessionTimeout: 30 * 60 * 1000, // 30 minutes
    providers: {
      googleAnalytics: null as any,
      mixpanel: null as any,
      amplitude: null as any,
      segment: null as any,
    },
  };

  private constructor() {
    this.sessionId = this.getOrCreateSessionId();
    this.initialize();
  }

  static getInstance(): AnalyticsService {
    if (!this.instance) {
      this.instance = new AnalyticsService();
    }
    return this.instance;
  }

  /**
   * Initialize analytics
   */
  private initialize() {
    // Load user profile
    this.loadUserProfile();

    // Setup auto-tracking
    if (this.config.enableAutoTrack) {
      this.setupAutoTracking();
    }

    // Setup heatmap tracking
    if (this.config.enableHeatmap) {
      this.setupHeatmapTracking();
    }

    // Start batch processor
    this.startBatchProcessor();

    // Track initial page view
    this.trackPageView();

    // Setup page visibility tracking
    this.setupVisibilityTracking();

    // Initialize third-party providers
    this.initializeProviders();
  }

  /**
   * Initialize third-party analytics providers
   */
  private initializeProviders() {
    // Google Analytics 4
    if (typeof window !== 'undefined' && (window as any).gtag) {
      this.config.providers.googleAnalytics = (window as any).gtag;
    }

    // Mixpanel
    if (typeof window !== 'undefined' && (window as any).mixpanel) {
      this.config.providers.mixpanel = (window as any).mixpanel;
    }

    // Amplitude
    if (typeof window !== 'undefined' && (window as any).amplitude) {
      this.config.providers.amplitude = (window as any).amplitude;
    }

    // Segment
    if (typeof window !== 'undefined' && (window as any).analytics) {
      this.config.providers.segment = (window as any).analytics;
    }
  }

  /**
   * Track custom event
   */
  track(
    eventName: string,
    properties?: Record<string, any>,
    options?: {
      category?: string;
      action?: string;
      label?: string;
      value?: number;
    }
  ) {
    const event: AnalyticsEvent = {
      name: eventName,
      category: options?.category || 'custom',
      action: options?.action,
      label: options?.label,
      value: options?.value,
      metadata: properties,
      timestamp: new Date(),
      sessionId: this.sessionId,
      userId: this.userProfile?.userId,
    };

    this.events.push(event);
    this.queueEvent(event);

    // Check goal completion
    this.checkGoalCompletion(event);

    // Update funnel progress
    this.updateFunnelProgress(event);

    // Send to third-party providers
    this.sendToProviders('track', eventName, properties);
  }

  /**
   * Track page view
   */
  trackPageView(path?: string, title?: string) {
    const pageView: PageView = {
      path: path || window.location.pathname,
      title: title || document.title,
      referrer: document.referrer,
      timestamp: new Date(),
    };

    this.pageViews.push(pageView);

    // Calculate page duration
    const now = Date.now();
    const duration = now - this.pageStartTime;
    if (this.pageViews.length > 1) {
      this.pageViews[this.pageViews.length - 2].duration = duration;
    }
    this.pageStartTime = now;

    // Queue event
    this.queueEvent({
      type: 'pageview',
      ...pageView,
    });

    // Send to providers
    this.sendToProviders('page', pageView.path, { title: pageView.title });
  }

  /**
   * Identify user
   */
  identify(userId: string, traits?: Record<string, any>) {
    this.userProfile = {
      userId,
      traits: traits || {},
      segments: this.calculateSegments(traits),
      createdAt: new Date(),
      lastSeen: new Date(),
    };

    // Save to storage
    this.saveUserProfile();

    // Queue identification event
    this.queueEvent({
      type: 'identify',
      userId,
      traits,
    });

    // Send to providers
    this.sendToProviders('identify', userId, traits);
  }

  /**
   * Track conversion goal
   */
  trackGoal(goalId: string, value?: number) {
    const goal = this.goals.get(goalId);
    if (!goal) {
      console.warn(`Goal ${goalId} not found`);
      return;
    }

    goal.isCompleted = true;
    goal.value = value;

    this.track('goal_completed', {
      goalId,
      goalName: goal.name,
      value,
    });
  }

  /**
   * Define conversion goal
   */
  defineGoal(goal: ConversionGoal) {
    this.goals.set(goal.id, goal);
  }

  /**
   * Check if event completes any goals
   */
  private checkGoalCompletion(event: AnalyticsEvent) {
    this.goals.forEach((goal) => {
      if (goal.isCompleted) return;

      if (
        goal.type === 'event' &&
        this.matchesConditions(event, goal.conditions)
      ) {
        goal.isCompleted = true;
        this.track('goal_completed', {
          goalId: goal.id,
          goalName: goal.name,
          triggerEvent: event.name,
        });
      }
    });
  }

  /**
   * Define funnel
   */
  defineFunnel(funnel: Funnel) {
    this.funnels.set(funnel.id, funnel);
  }

  /**
   * Update funnel progress
   */
  private updateFunnelProgress(event: AnalyticsEvent) {
    this.funnels.forEach((funnel) => {
      const currentStepIndex = funnel.steps.findIndex(
        (step) => !step.completed && step.event === event.name
      );

      if (currentStepIndex !== -1) {
        funnel.steps[currentStepIndex].completed = true;

        // Calculate metrics
        const completedSteps = funnel.steps.filter((s) => s.completed).length;
        funnel.conversionRate = (completedSteps / funnel.steps.length) * 100;

        // Track funnel progress
        this.track('funnel_step_completed', {
          funnelId: funnel.id,
          funnelName: funnel.name,
          stepIndex: currentStepIndex,
          stepName: funnel.steps[currentStepIndex].name,
          conversionRate: funnel.conversionRate,
        });
      }
    });
  }

  /**
   * Setup auto-tracking
   */
  private setupAutoTracking() {
    // Click tracking
    document.addEventListener('click', (e) => {
      const target = e.target as HTMLElement;
      const trackingData = this.getElementTrackingData(target);

      if (trackingData) {
        this.track('click', trackingData, {
          category: 'interaction',
          action: 'click',
        });
      }
    });

    // Form tracking
    document.addEventListener('submit', (e) => {
      const form = e.target as HTMLFormElement;
      this.track(
        'form_submit',
        {
          formId: form.id,
          formName: form.name,
          formAction: form.action,
        },
        {
          category: 'interaction',
          action: 'form_submit',
        }
      );
    });

    // Error tracking
    window.addEventListener('error', (e) => {
      this.track(
        'error',
        {
          message: e.message,
          filename: e.filename,
          lineno: e.lineno,
          colno: e.colno,
          stack: e.error?.stack,
        },
        {
          category: 'error',
          action: 'javascript_error',
        }
      );
    });
  }

  /**
   * Setup heatmap tracking
   */
  private setupHeatmapTracking() {
    // Track clicks for heatmap
    document.addEventListener('click', (e) => {
      const target = e.target as HTMLElement;
      const selector = this.getElementSelector(target);

      const data: HeatmapData = {
        element: selector,
        x: e.pageX,
        y: e.pageY,
        clicks: 1,
        hovers: 0,
        scrollDepth: this.getScrollDepth(),
      };

      this.heatmapData.push(data);

      // Aggregate data
      if (this.heatmapData.length >= 100) {
        this.sendHeatmapData();
      }
    });

    // Track scroll depth
    let maxScrollDepth = 0;
    window.addEventListener('scroll', () => {
      const depth = this.getScrollDepth();
      if (depth > maxScrollDepth) {
        maxScrollDepth = depth;
        this.track(
          'scroll_depth',
          { depth },
          {
            category: 'engagement',
            value: depth,
          }
        );
      }
    });
  }

  /**
   * Start session recording
   */
  startRecording() {
    if (this.isRecording) return;

    this.isRecording = true;
    this.sessionRecording = {
      sessionId: this.sessionId,
      userId: this.userProfile?.userId,
      startTime: new Date(),
      events: [],
      metadata: {
        userAgent: navigator.userAgent,
        screenResolution: `${screen.width}x${screen.height}`,
        viewport: `${window.innerWidth}x${window.innerHeight}`,
      },
    };

    // Record DOM mutations - disabled for type issues
    // TODO: Fix type issue with 'dom_change'
    /*
    const observer = new MutationObserver((mutations) => {
      this.recordEvent({
        type: 'dom_change',
        data: mutations.map((m) => ({
          type: m.type,
          target: this.getElementSelector(m.target as HTMLElement),
        })),
      });
    });

    observer.observe(document.body, {
      childList: true,
      subtree: true,
      attributes: true,
    });
    */

    // Record interactions
    this.recordInteractions();
  }

  /**
   * Stop session recording
   */
  stopRecording() {
    if (!this.isRecording) return;

    this.isRecording = false;
    if (this.sessionRecording) {
      this.sessionRecording.endTime = new Date();
      this.sendRecording(this.sessionRecording);
      this.sessionRecording = null;
    }
  }

  /**
   * Record interactions for session replay
   */
  private recordInteractions() {
    if (!this.isRecording) return;

    // Record clicks
    document.addEventListener('click', (e) => {
      this.recordEvent({
        type: 'click',
        target: this.getElementSelector(e.target as HTMLElement),
        data: { x: e.pageX, y: e.pageY },
      });
    });

    // Record scrolls
    let scrollTimeout: any;
    window.addEventListener('scroll', () => {
      clearTimeout(scrollTimeout);
      scrollTimeout = setTimeout(() => {
        this.recordEvent({
          type: 'scroll',
          data: { x: window.scrollX, y: window.scrollY },
        });
      }, 100);
    });

    // Record input changes
    document.addEventListener('input', (e) => {
      const target = e.target as HTMLInputElement;
      this.recordEvent({
        type: 'input',
        target: this.getElementSelector(target),
        data: {
          value: target.type === 'password' ? '***' : target.value,
        },
      });
    });
  }

  /**
   * Record event for session replay
   */
  private recordEvent(event: Partial<RecordedEvent>) {
    if (!this.sessionRecording) return;

    this.sessionRecording.events.push({
      type: event.type as any,
      target: event.target,
      data: event.data,
      timestamp: Date.now() - this.sessionRecording.startTime.getTime(),
    });
  }

  /**
   * Utility methods
   */
  private getElementTrackingData(
    element: HTMLElement
  ): Record<string, any> | null {
    const data: Record<string, any> = {
      tagName: element.tagName.toLowerCase(),
      text: element.textContent?.slice(0, 100),
      classList: Array.from(element.classList).join(' '),
      id: element.id,
    };

    // Check for data attributes
    Array.from(element.attributes).forEach((attr) => {
      if (attr.name.startsWith('data-track-')) {
        data[attr.name.replace('data-track-', '')] = attr.value;
      }
    });

    return Object.keys(data).length > 0 ? data : null;
  }

  private getElementSelector(element: HTMLElement): string {
    const path: string[] = [];

    while (element && element !== document.body) {
      let selector = element.tagName.toLowerCase();

      if (element.id) {
        selector += `#${element.id}`;
        path.unshift(selector);
        break;
      } else if (element.className) {
        selector += `.${element.className.split(' ').join('.')}`;
      }

      path.unshift(selector);
      element = element.parentElement!;
    }

    return path.join(' > ');
  }

  private getScrollDepth(): number {
    const scrolled = window.scrollY + window.innerHeight;
    const height = document.documentElement.scrollHeight;
    return Math.round((scrolled / height) * 100);
  }

  private matchesConditions(event: any, conditions: GoalCondition[]): boolean {
    return conditions.every((condition) => {
      const value = event[condition.field] || event.metadata?.[condition.field];

      switch (condition.operator) {
        case 'equals':
          return value === condition.value;
        case 'contains':
          return String(value).includes(String(condition.value));
        case 'greater':
          return Number(value) > Number(condition.value);
        case 'less':
          return Number(value) < Number(condition.value);
        case 'regex':
          return new RegExp(condition.value).test(String(value));
        default:
          return false;
      }
    });
  }

  private calculateSegments(traits?: Record<string, any>): string[] {
    const segments: string[] = [];

    if (!traits) return segments;

    // Example segmentation logic
    if (traits.plan === 'premium') segments.push('premium_users');
    if (traits.totalPurchases > 10) segments.push('power_users');
    if (
      traits.lastPurchase &&
      Date.now() - traits.lastPurchase < 30 * 24 * 60 * 60 * 1000
    ) {
      segments.push('active_buyers');
    }

    return segments;
  }

  /**
   * Batch processing
   */
  private queueEvent(event: any) {
    this.eventQueue.push({
      ...event,
      timestamp: new Date().toISOString(),
      sessionId: this.sessionId,
      userId: this.userProfile?.userId,
    });

    if (this.eventQueue.length >= this.config.batchSize) {
      this.flushEvents();
    }
  }

  private startBatchProcessor() {
    setInterval(() => {
      if (this.eventQueue.length > 0) {
        this.flushEvents();
      }
    }, this.config.batchInterval);
  }

  private async flushEvents() {
    if (this.eventQueue.length === 0) return;

    const events = [...this.eventQueue];
    this.eventQueue = [];

    try {
      await fetch(`${this.config.apiEndpoint}/events`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ events }),
      });
    } catch (error) {
      console.error('Failed to send analytics events:', error);
      // Re-queue events
      this.eventQueue.unshift(...events);
    }
  }

  private async sendHeatmapData() {
    const data = [...this.heatmapData];
    this.heatmapData = [];

    try {
      await fetch(`${this.config.apiEndpoint}/heatmap`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          data,
          page: window.location.pathname,
        }),
      });
    } catch (error) {
      console.error('Failed to send heatmap data:', error);
    }
  }

  private async sendRecording(recording: SessionRecording) {
    try {
      await fetch(`${this.config.apiEndpoint}/recording`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(recording),
      });
    } catch (error) {
      console.error('Failed to send session recording:', error);
    }
  }

  /**
   * Third-party provider integration
   */
  private sendToProviders(method: string, ...args: any[]) {
    // Google Analytics
    if (this.config.providers.googleAnalytics) {
      try {
        this.config.providers.googleAnalytics('event', ...args);
      } catch (error) {
        console.debug('GA tracking error:', error);
      }
    }

    // Mixpanel
    if (this.config.providers.mixpanel) {
      try {
        this.config.providers.mixpanel[method](...args);
      } catch (error) {
        console.debug('Mixpanel tracking error:', error);
      }
    }

    // Amplitude
    if (this.config.providers.amplitude) {
      try {
        if (method === 'track') {
          this.config.providers.amplitude.logEvent(...args);
        } else if (method === 'identify') {
          this.config.providers.amplitude.setUserId(args[0]);
        }
      } catch (error) {
        console.debug('Amplitude tracking error:', error);
      }
    }

    // Segment
    if (this.config.providers.segment) {
      try {
        this.config.providers.segment[method](...args);
      } catch (error) {
        console.debug('Segment tracking error:', error);
      }
    }
  }

  /**
   * Session management
   */
  private getOrCreateSessionId(): string {
    let sessionId = sessionStorage.getItem('analytics_session_id');

    if (!sessionId) {
      sessionId = `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      sessionStorage.setItem('analytics_session_id', sessionId);
    }

    return sessionId;
  }

  /**
   * User profile management
   */
  private loadUserProfile() {
    const stored = localStorage.getItem('analytics_user_profile');
    if (stored) {
      try {
        this.userProfile = JSON.parse(stored);
      } catch (error) {
        console.error('Failed to load user profile:', error);
      }
    }
  }

  private saveUserProfile() {
    if (this.userProfile) {
      localStorage.setItem(
        'analytics_user_profile',
        JSON.stringify(this.userProfile)
      );
    }
  }

  /**
   * Visibility tracking
   */
  private setupVisibilityTracking() {
    document.addEventListener('visibilitychange', () => {
      if (document.hidden) {
        this.track('page_hidden', {
          timeOnPage: Date.now() - this.pageStartTime,
        });
        this.flushEvents();
      } else {
        this.track('page_visible');
      }
    });
  }

  /**
   * Public API
   */
  getEvents(): AnalyticsEvent[] {
    return [...this.events];
  }

  getPageViews(): PageView[] {
    return [...this.pageViews];
  }

  getGoals(): ConversionGoal[] {
    return Array.from(this.goals.values());
  }

  getFunnels(): Funnel[] {
    return Array.from(this.funnels.values());
  }

  getUserProfile(): UserProfile | null {
    return this.userProfile;
  }

  getSessionId(): string {
    return this.sessionId;
  }

  reset() {
    this.events = [];
    this.pageViews = [];
    this.goals.clear();
    this.funnels.clear();
    this.heatmapData = [];
    this.sessionRecording = null;
    this.flushEvents();
  }
}

// Export singleton instance
export const analytics = AnalyticsService.getInstance();

// Export types
export type {
  AnalyticsEvent,
  PageView,
  UserProfile,
  ConversionGoal,
  Funnel,
  HeatmapData,
  SessionRecording,
};
