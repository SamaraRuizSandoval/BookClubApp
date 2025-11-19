import '@testing-library/jest-dom'; // or '@testing-library/jest-dom/vitest'
import { beforeEach, afterEach, vi } from 'vitest';

// --- Control timers so Ionic's timeouts don't fire after teardown ---
beforeEach(() => {
  // Use fake timers so we control any setTimeout / setInterval that Ionic sets up
  vi.useFakeTimers();
});

afterEach(() => {
  // Run anything pending (like the ion-app timeout) *while window still exists*
  vi.runOnlyPendingTimers();

  // Back to real timers for next test/file
  vi.useRealTimers();

  // Optional: clean up DOM between tests
  document.body.innerHTML = '';
  document.head.innerHTML = '';
});

// --- Your existing JSDOM polyfills/mocks ---

// JSDOM polyfills/mocks
if (!window.matchMedia) {
  // Minimal mock that Ionic's SplitPane expects (addListener/removeListener)
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
