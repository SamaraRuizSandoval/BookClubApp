import '@testing-library/jest-dom';

// JSDOM polyfills/mocks
if (!window.matchMedia) {
  // Minimal mock that Ionic's SplitPane expects (addListener/removeListener)
  // Vitest + JSDOM don’t implement matchMedia by default.
  // This keeps IonSplitPane from throwing during tests.
  // You can extend this if you need more behavior.
   
  window.matchMedia = ((query: string): any => {
    return {
      matches: false,
      media: query,
      onchange: null,
      addListener: () => {}, // deprecated but used by Ionic
      removeListener: () => {}, // deprecated but used by Ionic
      addEventListener: () => {},
      removeEventListener: () => {},
      dispatchEvent: () => false,
    };
  }) as unknown as typeof window.matchMedia;
}

// Some Ionic components also touch ResizeObserver in certain layouts.
// Safe no-op mock to avoid surprises later.
if (!('ResizeObserver' in window)) {
   
  (window as any).ResizeObserver = class {
    observe() {}
    unobserve() {}
    disconnect() {}
  };
}
