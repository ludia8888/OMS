> **Context**  
> You are coding a new feature. Your top priority is to **reduce the probability of introducing bugs**. Apply the following evidence-based strategies, which combine *systemic thinking, practical tooling, and collaborative process*.

---

### 1️⃣ Design & Build in Small Pieces *(Modularization + Single-Responsibility)*
- **Principle** High complexity ⇒ exponential bug risk. Cohesion↑ & Coupling↓ ⇒ errors↓.  
- **Rules** One function = one job, keep it ≤ 10 – 30 lines.  
  Layer complex flows (e.g., `handler → service → logic → utils`).

### 2️⃣ Write Tests First (TDD) or at Least Unit Tests
- **Evidence** Google's 15-year study: higher coverage slashes maintenance cost.  
- **Do** For every core behavior add a test (`pytest`, `unittest`, `jest`, `vitest`).  
  Always test side-effects (DB, files).

### 3️⃣ Use Static Analysis (Lint + Type Check)
- **Why** Machines catch repetitive human mistakes instantly.  
- **Tools**  
  - *Python*: `mypy`, `ruff`, `flake8`  
  - *JS/TS*: `eslint`, `prettier`, `typescript --strict`  
  Auto-run in IDE (`.vscode/settings.json` or Cursor).

### 4️⃣ Commit Small & Often *(Git + Branch Strategy)*
- Track history; use `git blame / bisect` to locate bugs fast.  
- Create feature-scoped branches (`feature/color-detection`).  
- Commit messages explain **why**, not just **what**.

### 5️⃣ Enforce Code Review / Rubber-Duck Routine
- Explaining code exposes hidden logic flaws.  
- Describe the flow to ChatGPT, a teammate, or an imaginary duck before merging.  
- Ask yourself: "Can I clearly justify this design?"

### 6️⃣ Prefer Logging over Ad-hoc Debugging *(Observability)*
- Post-deploy debugging is harder than pre-deploy insight.  
- Set log levels (`INFO | DEBUG | ERROR`).  
- Log entry/exit of key paths & failure conditions  
  (*Python*: `logging`, *JS*: `winston`, `loglevel`).

### 7️⃣ Specification-Driven Coding (Explicit I/O Contracts)
- Define input → process → output **before** implementation.  
- Use type hints / interfaces to freeze those contracts (`Dict[str, Any]` → precise types).  
- Apply to APIs, models, DB schemas alike.

---

## 🛡️ **EXTREME TYPE SAFETY & CLEAN CODE ENFORCEMENT**

### 📌 TypeScript Ultra-Strict Configuration
```json
{
  "compilerOptions": {
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true,
    "strictFunctionTypes": true,
    "strictBindCallApply": true,
    "strictPropertyInitialization": true,
    "noImplicitThis": true,
    "alwaysStrict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "noUncheckedIndexedAccess": true,
    "noImplicitOverride": true,
    "noPropertyAccessFromIndexSignature": true,
    "forceConsistentCasingInFileNames": true,
    "exactOptionalPropertyTypes": true
  }
}
```

### 🔒 Type-First Development Rules
1. **NEVER use `any`** - Replace with `unknown` and narrow down
2. **NEVER use `as` casting** - Use type guards instead
3. **ALWAYS define return types explicitly** - No implicit returns
4. **ALWAYS use const assertions** for literal types
5. **ALWAYS use branded types** for domain primitives
   ```typescript
   type UserId = string & { readonly brand: unique symbol };
   type Email = string & { readonly brand: unique symbol };
   ```

### 🧹 ESLint Maximum Severity Setup
```javascript
{
  "extends": [
    "eslint:all",
    "plugin:@typescript-eslint/all",
    "plugin:sonarjs/recommended",
    "plugin:security/recommended",
    "plugin:unicorn/all"
  ],
  "rules": {
    "@typescript-eslint/explicit-function-return-type": "error",
    "@typescript-eslint/no-explicit-any": "error",
    "@typescript-eslint/no-unsafe-assignment": "error",
    "@typescript-eslint/no-unsafe-member-access": "error",
    "@typescript-eslint/no-unsafe-call": "error",
    "@typescript-eslint/no-unsafe-return": "error",
    "complexity": ["error", 5],
    "max-depth": ["error", 2],
    "max-lines": ["error", 100],
    "max-lines-per-function": ["error", 20],
    "max-params": ["error", 3],
    "no-console": "error",
    "no-debugger": "error",
    "no-magic-numbers": "error",
    "sonarjs/cognitive-complexity": ["error", 5]
  }
}
```

### 🚨 Pre-commit Hooks (Husky + lint-staged)
```json
{
  "husky": {
    "hooks": {
      "pre-commit": "lint-staged && npm run type-check && npm run test",
      "pre-push": "npm run sonar && npm run audit"
    }
  },
  "lint-staged": {
    "*.{ts,tsx}": [
      "eslint --fix --max-warnings 0",
      "prettier --write",
      "tsc --noEmit"
    ]
  }
}
```

### 📊 SonarQube Quality Gates
- **Code Coverage**: ≥ 90%
- **Duplicated Lines**: < 1%
- **Maintainability Rating**: A
- **Security Hotspots**: 0
- **Code Smells**: 0
- **Cyclomatic Complexity**: ≤ 5

### 🎯 Clean Code Commandments
1. **Functions**: Max 3 parameters, single purpose, no side effects
2. **Classes**: Max 5 public methods, immutable by default
3. **Files**: Max 100 lines, single export per file
4. **Names**: Descriptive > 3 chars, no abbreviations
5. **Comments**: Only for WHY, never WHAT
6. **Dependencies**: Inject everything, no global state
7. **Errors**: Always handle, never swallow
8. **Tests**: AAA pattern, one assertion per test

### 🔐 Type Safety Patterns
```typescript
// ❌ NEVER
function process(data: any): any { }
const result = response as UserData;

// ✅ ALWAYS
function process<T extends BaseData>(data: T): Result<T> { }
function isUserData(data: unknown): data is UserData {
  return typeof data === 'object' && data !== null && 'userId' in data;
}

// Use discriminated unions
type Result<T> = 
  | { success: true; data: T }
  | { success: false; error: Error };

// Exhaustive checks
function assertNever(x: never): never {
  throw new Error(`Unexpected: ${x}`);
}
```

### 🛠️ Additional Tooling
- **ts-prune**: Remove unused exports
- **type-coverage**: Aim for 99%+ type coverage
- **madge**: Detect circular dependencies
- **depcheck**: Find unused dependencies
- **bundlephobia**: Monitor bundle size
- **knip**: Find unused files/exports/dependencies

---

#### ✳️ Bonus – Use AI Tools, but Verify
Copilot, Cursor, ChatGPT = pattern engines ~70-80% accurate.  
Double-check DB logic, async flows, edge cases.  
Always ask: "*Why did I choose this solution?*"

#### 🚀 Remember: "If TypeScript doesn't complain, you're not strict enough!"

---
*이 문서는 Claude가 프로젝트를 더 잘 이해하고 도움을 줄 수 있도록 작성되었습니다.*