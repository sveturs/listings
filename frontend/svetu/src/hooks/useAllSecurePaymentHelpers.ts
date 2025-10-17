// Helper functions for navigation - can be mocked in tests

export const navigateToUrl = (url: string) => {
  window.location.href = url;
};

export const getLocationOrigin = () => {
  return typeof window !== 'undefined' ? window.location.origin : '';
};
