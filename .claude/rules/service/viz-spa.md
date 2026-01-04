---
description: viz-spa specific guardrails for GraphQL schema alignment and component quality
paths: services/cursor-viz-spa/**
---

# viz-spa Rules

Service-specific constraints for P6 (viz-spa React dashboard).

---

## NEVER

- **Manually define GraphQL types** in P6 (use codegen from P5 schema when available)
- **Hardcode API queries** (define all queries in `.graphql` files)
- **Use custom CSS** (Tailwind CSS only)
- **Create class components** (functional components + hooks only)
- **Import non-Tailwind fonts** without justification
- **Skip accessibility checks** (WCAG 2.1 AA minimum)
- **Use prop drilling** excessively (3+ levels → context/hook)
- **Render without error boundaries** for large sections

---

## ALWAYS

- **Verify GraphQL queries match P5 schema** before implementation
- **Use GraphQL codegen** (when available) to generate types
- **Test components with renderWithProviders** (includes Apollo + Router)
- **Include ARIA labels** for interactive components
- **Use semantic HTML**: `<button>` not `<div>`
- **Support keyboard navigation**: Tab, Enter, Escape
- **Provide loading and error states**: Every GraphQL query
- **Target 80%+ test coverage** for components
- **Follow React hooks best practices**: Dependency arrays, no conditional hooks

---

## GraphQL Query Alignment

### Schema Verification
**Before implementing a query**:
1. Check P5 GraphQL schema (`.graphql` file or introspection)
2. Verify field names match exactly
3. Verify argument types match
4. Verify return types match

### Query Definition
- **Location**: `src/graphql/queries/` or similar
- **Format**: `.graphql` files (not string literals)
- **Naming**: Descriptive names (DashboardQuery, DevelopersListQuery)
- **Example**:
```graphql
query GetDashboard($startDate: String!, $endDate: String!) {
  dashboard(startDate: $startDate, endDate: $endDate) {
    kpis { ... }
    topPerformers { ... }
  }
}
```

### Codegen Integration
- **When available**: Use GraphQL codegen to generate types
- **Never manually define**: `interface Developer {}`
- **Always use generated**: `import { Developer } from '@generated/types'`
- **Update on schema change**: Regenerate types immediately

---

## React Components

### Structure
- **Functional components**: No class components
- **Props interface**: Define clear input types
- **Return JSX**: One component per file

### Hooks
- **useState**: For local component state
- **useEffect**: For side effects with dependency array
- **Custom hooks**: For shared logic
- **NO conditional hooks**: Never call hooks conditionally

### Testing
- **Unit tests**: Each component in isolation
- **Integration tests**: With Apollo provider + Router
- **Test file**: `Component.test.tsx` next to `Component.tsx`
- **Use renderWithProviders**: From test utilities

---

## Styling (Tailwind Only)

### Utility Classes
- **Stack utilities**: `flex`, `grid`, `space-x-4`
- **Responsive**: `md:grid-cols-3`, `sm:px-2`
- **Dark mode**: `dark:bg-gray-900`
- **Hover/focus**: `hover:bg-blue-600`, `focus:ring-2`

### No Custom CSS
- **Never create**: `.css` or `<style>` blocks
- **Never use**: `styled-components`, `CSS Modules`
- **Override with**: Tailwind `@apply` only if critical (rare)

### Responsive Design
- **Mobile-first**: Design for mobile, add breakpoints
- **Breakpoints**: `sm` (640px), `md` (768px), `lg` (1024px), `xl`
- **Test all breakpoints**: Use browser dev tools

---

## Accessibility (WCAG 2.1 AA)

### Semantic HTML
- **Use**: `<button>`, `<input>`, `<label>`, `<nav>`, `<main>`
- **Avoid**: `<div onClick>`, generic containers

### ARIA Labels
- **Interactive elements**: `aria-label="Close modal"`
- **Complex elements**: `aria-describedby`, `aria-owns`
- **Live regions**: `aria-live="polite"` for updates
- **Current page**: `aria-current="page"` in nav

### Keyboard Navigation
- **Tab order**: Logical progression (use tabindex sparingly)
- **Focus visible**: Always visible focus outline
- **Escape key**: Close modals/popovers
- **Test with**: Keyboard only, no mouse

### Color & Contrast
- **Text contrast**: 4.5:1 for normal text, 3:1 for large
- **Color alone**: Don't use color only to convey information
- **Test with**: WAVE accessibility auditor

### Form Accessibility
- **Labels**: `<label htmlFor="id">` for inputs
- **Error messages**: Associated with inputs
- **Validation**: Clear error text, not just red
- **Focus management**: After form submission

---

## Performance

### Code Splitting
- **Routes**: Lazy load page components
- **Heavy components**: Split large components
- **Bundle size**: Monitor with webpack-bundle-analyzer

### Apollo Client
- **Fetch policy**: `cache-and-network` for freshness
- **Refetch**: Use `refetch()` for manual updates
- **Cache**: Let Apollo handle automatic caching
- **Error handling**: Show errors to user

---

## Testing Requirements

### Unit Tests
- Component render: Happy path + error states
- Props validation: Different prop combinations
- User interactions: Click, type, submit
- Coverage: 80%+ for components

### Integration Tests
- With Apollo provider: Real queries
- With React Router: Navigation
- Data flow: Query → Component → Display

### Visual Testing
- Screenshot baseline with Playwright
- Responsive layouts at breakpoints
- Dark mode appearance
- Focus states

---

## State Management

### When to Use What
- **Component state**: useState (single component)
- **Global state**: Apollo cache or Context (rarely)
- **URL state**: useSearchParams (filters, pagination)
- **Router state**: useLocation (current page)

### Apollo Cache
- **Trust cache**: Update automatically on mutations
- **Manual updates**: `cache.modify()` when needed
- **Eviction**: Clear cache on logout

---

## Documentation

### Component Comments
- **Props documentation**: JSDoc comments
- **Usage examples**: In Storybook or comments
- **Accessibility**: Special requirements

### Commit Messages
```
feat(viz-spa): add {component/feature}

Implements {description}.

## Changes
- Added Component: {name}
- Schema: Verified against P5 {query}

## Testing
- Tests: X new cases
- Coverage: Y%
- A11y: WCAG 2.1 AA verified
- Responsive: Tested at sm/md/lg
```

---

## See Also

- React/Vite patterns in `.claude/skills/react-vite-patterns/`
- GraphQL patterns in `.claude/skills/typescript-graphql-patterns/`
- API contract in `.claude/skills/api-contract/SKILL.md`
- Global coding standards in `03-coding-standards.md`
- SPEC.md: `services/cursor-viz-spa/SPEC.md`
