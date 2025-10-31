# 🤝 Contributing Guide

## Code Standards

### TypeScript
```typescript
// ✅ Good
export async function getEmployeesAPI(
  filter?: EmployeeFilter
): Promise<ApiResponse<PaginatedResponse<Employee>>> {
  // Implementation
}

// ❌ Bad
export function getEmployees(f: any) {
  // Implementation
}
```

### React Components
```typescript
// ✅ Functional component with hooks
export const EmployeeCard: React.FC<EmployeeCardProps> = ({ employee }) => {
  const [data, setData] = useState(null);
  return <div>{employee.name}</div>;
};

// ✅ Use proper typing
interface EmployeeCardProps {
  employee: Employee;
  onEdit?: (employee: Employee) => void;
}
```

### File Structure
```
- Component file: PascalCase (EmployeeForm.tsx)
- Hook file: camelCase with use prefix (useEmployeeForm.ts)
- API file: camelCase with Api suffix (employees.api.ts)
- Type file: index.ts in types folder
- Style file: same name as component (.module.css)
```

## Commit Message Format

```
<type>: <subject>

<body>

<footer>
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Code style (formatting)
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Build/tooling

### Examples
```
✅ feat: add employee face recognition
✅ fix: resolve login token expiry issue
✅ docs: update API guide
✅ refactor: simplify attendance logic
❌ update stuff
❌ fix bug
```

## Pull Request Process

1. **Create feature branch**
   ```bash
   git checkout -b feat/your-feature-name
   ```

2. **Implement changes**
   - Follow code standards
   - Add proper TypeScript types
   - Add comments for complex logic

3. **Test locally**
   ```bash
   npm run dev
   npm run build
   npm run lint
   npm run type-check
   ```

4. **Push and create PR**
   ```bash
   git push origin feat/your-feature-name
   ```

5. **PR Description**
   - What changed
   - Why changed
   - How to test
   - Screenshots if UI changes

## Code Review Checklist

- [ ] Code follows style guide
- [ ] TypeScript types correct
- [ ] No console.log or debug code
- [ ] Error handling implemented
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No breaking changes

## Testing

### Unit Tests
```typescript
import { render, screen } from '@testing-library/react';
import { EmployeeForm } from './EmployeeForm';

describe('EmployeeForm', () => {
  it('should render form fields', () => {
    render(<EmployeeForm />);
    expect(screen.getByLabelText('Name')).toBeInTheDocument();
  });
});
```

Run tests:
```bash
npm run test
npm run test -- --watch
```

### Integration Tests
- Test with real API
- Verify data flow
- Check error handling

## Performance Guidelines

- Use React.memo for expensive components
- Implement proper loading states
- Optimize images and assets
- Use pagination for large lists
- Debounce search/filter inputs

## Accessibility (a11y)

- Use semantic HTML
- Add alt text to images
- Keyboard navigation support
- ARIA labels where needed
- Color contrast ratios

```typescript
// ✅ Good
<button aria-label="Delete employee" onClick={handleDelete}>
  <TrashIcon />
</button>

// ❌ Bad
<div onClick={handleDelete}>
  <TrashIcon />
</div>
```

## Documentation

- Update README.md for major changes
- Add JSDoc comments to functions
- Document complex logic
- Keep API_GUIDE.md updated

```typescript
/**
 * Get list of employees with filters
 * @param filter - Filter options (page, limit, search, etc.)
 * @returns Paginated list of employees
 * @example
 * const result = await getEmployeesAPI({ page: 1, limit: 10 });
 */
export async function getEmployeesAPI(filter?: EmployeeFilter) {
  // ...
}
```

## Common Issues & Solutions

### ESLint Errors
```bash
npm run lint -- --fix
```

### TypeScript Errors
```bash
npm run type-check
# Fix type issues in your editor
```

### Build Failures
```bash
npm run clean
npm install
npm run build
```

## Questions?

- Check [README.md](README.md)
- Review [QUICK_START.md](QUICK_START.md)
- Check [SETUP_GUIDE.md](SETUP_GUIDE.md)
- Create GitHub issue
