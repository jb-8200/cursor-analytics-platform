# Design Document: Performance Testing with Lighthouse

**Feature ID**: P6-F06
**Epic**: P6 - cursor-viz-spa (Testing Enhancement)
**Created**: January 4, 2026
**Status**: PROPOSED

---

## Overview

Implement Lighthouse CI for automated performance testing and regression detection in cursor-viz-spa.

---

## Architecture

### Performance Testing Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│  CI Pipeline (on PR)                                                   │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  1. Build P6                                                           │
│     npm run build                                                      │
│                                                                         │
│  2. Start Preview Server                                               │
│     npm run preview (serves dist/)                                     │
│                                                                         │
│  3. Run Lighthouse CI                                                  │
│     lhci autorun                                                       │
│     ┌───────────────────────────────────────────────────────────────┐   │
│     │ Metrics:                                                      │   │
│     │ • First Contentful Paint: 1.2s ✅                             │   │
│     │ • Time to Interactive: 3.5s ✅                                │   │
│     │ • Speed Index: 2.8s ✅                                        │   │
│     │ • Cumulative Layout Shift: 0.05 ✅                            │   │
│     └───────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  4. Compare to Budget                                                  │
│     ┌───────────────────────────────────────────────────────────────┐   │
│     │ Budget Check:                                                 │   │
│     │ • FCP < 2000ms: PASS (1200ms)                                 │   │
│     │ • TTI < 5000ms: PASS (3500ms)                                 │   │
│     │ • Performance Score > 80: PASS (87)                           │   │
│     └───────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  5. Upload Report                                                      │
│     HTML report → GitHub Artifacts                                     │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Configuration

### Lighthouse CI Config

```javascript
// services/cursor-viz-spa/lighthouserc.js
module.exports = {
  ci: {
    collect: {
      url: ['http://localhost:4173/dashboard'],
      startServerCommand: 'npm run preview',
      startServerReadyPattern: 'Local:',
      numberOfRuns: 3,
      settings: {
        preset: 'desktop',
        throttling: {
          cpuSlowdownMultiplier: 1,
        },
      },
    },
    assert: {
      preset: 'lighthouse:recommended',
      assertions: {
        // Core Web Vitals
        'first-contentful-paint': ['warn', { maxNumericValue: 2000 }],
        'interactive': ['error', { maxNumericValue: 5000 }],
        'speed-index': ['warn', { maxNumericValue: 3500 }],
        'cumulative-layout-shift': ['warn', { maxNumericValue: 0.1 }],
        // Performance score
        'categories:performance': ['warn', { minScore: 0.8 }],
        // Accessibility (bonus)
        'categories:accessibility': ['warn', { minScore: 0.9 }],
      },
    },
    upload: {
      target: 'filesystem',
      outputDir: './lighthouse-results',
    },
  },
};
```

### Performance Budget

```json
// services/cursor-viz-spa/lighthouse-budget.json
[
  {
    "path": "/dashboard",
    "resourceSizes": [
      { "resourceType": "script", "budget": 300 },
      { "resourceType": "stylesheet", "budget": 100 },
      { "resourceType": "image", "budget": 200 },
      { "resourceType": "total", "budget": 800 }
    ],
    "timings": [
      { "metric": "first-contentful-paint", "budget": 2000 },
      { "metric": "interactive", "budget": 5000 },
      { "metric": "speed-index", "budget": 3500 }
    ]
  }
]
```

---

## NPM Scripts

```json
{
  "scripts": {
    "lighthouse": "lhci autorun",
    "lighthouse:collect": "lhci collect --config=lighthouserc.js",
    "lighthouse:assert": "lhci assert --config=lighthouserc.js",
    "preview": "vite preview --port 4173"
  }
}
```

---

## CI Integration

### GitHub Actions

```yaml
# .github/workflows/lighthouse.yml
name: Lighthouse CI

on:
  pull_request:
    paths:
      - 'services/cursor-viz-spa/**'

jobs:
  lighthouse:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'

      - name: Install dependencies
        run: |
          cd services/cursor-viz-spa
          npm ci

      - name: Build
        run: |
          cd services/cursor-viz-spa
          npm run build

      - name: Run Lighthouse CI
        run: |
          cd services/cursor-viz-spa
          npm install -g @lhci/cli
          lhci autorun

      - name: Upload Lighthouse results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: lighthouse-results
          path: services/cursor-viz-spa/lighthouse-results/
          retention-days: 7
```

---

## Metrics Explained

| Metric | Target | Description |
|--------|--------|-------------|
| First Contentful Paint (FCP) | < 2s | Time until first content visible |
| Time to Interactive (TTI) | < 5s | Time until fully interactive |
| Speed Index | < 3.5s | How quickly content is visually displayed |
| Cumulative Layout Shift (CLS) | < 0.1 | Visual stability |
| Performance Score | > 80 | Overall Lighthouse score |

---

## Success Metrics

| Metric | Target | Why |
|--------|--------|-----|
| Dashboard FCP | < 2s | Users see content quickly |
| Dashboard TTI | < 5s | Dashboard is usable quickly |
| Performance Score | > 80 | Good overall performance |
| Regression Detection | Automated | Catch performance issues in PR |

---

## References

- [Lighthouse CI Documentation](https://github.com/GoogleChrome/lighthouse-ci)
- [Core Web Vitals](https://web.dev/vitals/)
- `docs/e2e-testing-strategy.md` (Phase 4)
