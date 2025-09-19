# AGENTS.md

This document guides AI agents and contributors working in this repository. It defines project structure, conventions, workflows, and safety rules. Its scope covers the entire repository. If a more deeply nested AGENTS.md exists, it overrides conflicting guidance for its subtree.

## Repository Structure & Assets
- `bin/`: Build artifacts emitted by the Makefile targets (created on demand).
- `cmd/app`: Main HTTP entrypoint and wiring; builds to `bin/app`.
- `cmd/kit`: Helper CLIs (`jwt`, `http_route`, `rdbms_seeder`, `permission_seeder`).
- `cmd/template`: Minimal bootstrap executable that wires registrars; useful for scaffolding.
- `internal`: Domain logic and supporting modules with colocated tests (`*_test.go`).
  - `http/`: Gin HTTP stack — `controller/`, `middleware/`, `response/`, `route/`, `server.go`.
  - `scheduler/`: Background jobs (cron-like tasks) and example jobs under `job/`.
  - `registrar/`: Service and dependency registration for bootstrapping.
  - `pkg/`: Project-specific packages (database, models, permission logic, user domain helpers, etc.).
  - `migration/`: Database migrations split into `rdbms/` and `clickhouse/` trees.
  - `test/`: Fixtures and helpers for tests (admin, permission, migration, HTTP utilities).
- `docs`: Swagger outputs (`swagger.json|yaml`, `docs.go`).
- `deployment/`: Docker assets (e.g., `deployment/docker/`).
- `storage`: Runtime artifacts such as `storage/log/access.log` and `storage/log/app.log`.
- Configuration & env templates live at the root: `.env*`, `.air.toml`, `config*.yml`, `.golangci.yml`.

## Key Technologies
- Language: Go (Gin web framework).
- Data: GORM (MySQL/PostgreSQL/SQLite), ClickHouse, Redis.
- Config: Viper (file-based + environment variables).
- AuthN/Z: JWT with Casbin-based authorization policies.
- Docs: Swagger annotations generate specs in `docs/`.
- Scheduling: Cron-style background tasks.

## Build, Run, and Tooling
- `make`: Compiles the main service and helper binaries with the `jsoniter` build tag enabled.
- `make all`: Builds `bin/app` alongside helper CLIs.
- `make run`: Runs the service locally with the race detector enabled.
- `go run cmd/app/*`: Spins up the API without the Makefile; use `make air` for hot reloading (requires the `air` tool).
- CLI builds: `make jwt`, `make http_route`, `make rdbms_seeder`, `make permission_seeder` build individual helpers into `bin/`.
- Docs: `make swagger` regenerates Swagger artifacts after updating annotations.
- Quality: `make linter` (golangci-lint with `errcheck` and `gosec`) and `golangci-lint run ./...` should pass before commits.
- Tests & benchmarks: `make test [args=./...]`, `make benchmark`, or `go test ./...` (ensure SQLite and `.env.testing` are available).

## Coding Style & Naming Conventions
- Run `gofmt`, `goimports`, and `go vet ./...` prior to review. Indentation must use tabs (gofmt default, width 4 spaces).
- Keep package names lowercase with no underscores; prefer descriptive names (`internal/http/controller/user`).
- Use natural exported identifiers (e.g., `UserLoginRequest`). For camelCase names containing "id", spell it as `Id` (e.g., `userId`).
- Avoid import aliases unless required for conflict resolution or clarity; rely on the default import name whenever possible.
- Reuse established helper patterns (e.g., response builders in `internal/http/response`) rather than handcrafting JSON payloads.
- Respect the `jsoniter` build tag where the Makefile enables it.
- Follow directory-local conventions and keep reusable code in `internal/pkg/`; leave HTTP wiring and other application layers under `internal/`.

### Example: indentation must use tabs
```golang
// use tab is Good example
package main

func main() {
	fmt.Println("example")
}

// use space is Bad example
package main

func main() {
  fmt.Println("example")
}
```

## Error Responses
- Define canonical error messages and codes exclusively in `internal/http/response/error_message.go`.
- Message constants describe user-facing text. Update `MessageToCode` with a unique `<status>-<sequence>` string (e.g., `400-001`).
- Organize the map by HTTP status headers (`// 400`, `// 401`, etc.) so related errors stay grouped.
- Ensure handlers return responses whose status aligns with the declared message and add corresponding tests when introducing new errors.

## Development Workflow
When adding a new API endpoint:
1. Extend or create handler files under `internal/http/controller/<area>` (for example, `internal/http/controller/user/user_register.go`). Follow the existing `<feature>_<action>.go` naming and add matching `*_test.go` coverage.
2. Add or adjust business logic in `internal/pkg/<domain>` (e.g., `internal/pkg/user`, `internal/pkg/permission`). Keep helpers reusable and include unit tests where practical.
3. Wire routes in the corresponding router under `internal/http/route/<area>/api_route.go`, ensuring it implements `route.Router` and guards handlers with middleware as needed.
4. If you introduce a brand-new router, register it in `internal/http/route/api_route.go` by appending it to the `routers` slice so it participates in `AttachRoutes`.
5. For admin/protected capabilities, synchronise seeds between runtime and tests: update `cmd/kit/permission_seeder/permission_seeder.go` and mirror the same permissions/roles in `internal/test/permission_service.go` (and adjust `internal/test/admin_service.go` when admin fixtures change).

General guidance:
- Keep changes minimal and localized; avoid unrelated refactors.
- Favor composition over duplication; reuse helpers under `internal/pkg/` and existing shared utilities.
- Update Swagger comments whenever API shapes change, then run `make swagger`.
- Wrap work in a transaction only when multiple insert/update/delete statements need to succeed together.

## Testing Guidelines
- Test framework: standard `testing`; `testify` is available for assertions.
- Location: place `*_test.go` files alongside the code they cover; prefer table-driven tests or testify suites.
- Bootstrapping: leverage `internal/test` utilities for environment loading, seeded users, migrations, and CSRF helpers instead of reimplementing setup logic.
- Migrations: run required migrations first (e.g., `make sqlite-migration args=up`). Tests expect `.env.testing` to provide DSNs and secrets.
- Coverage: exercise both success and failure paths. Document skipped integration tests in PR descriptions. Add controller/service tests whenever adding new endpoints.

## Migrations
- Use Make targets with environment loaded from `.env`: `make mysql-migration args="up"`, `make pgsql-migration args="up"`, `make sqlite-migration args="up"`, or `make clickhouse-migration args="up"`.
- Ensure relevant `DB_*` variables are set for the target database.
- Keep migrations idempotent and reversible; document non-trivial data movements.

## Security & Configuration
- Never commit secrets. Place credentials in `storage/...` as described in the README.
- Provide `.env.dev` and `.env.testing` locally; avoid committing real values.
- Validate configuration via `make run` and integration tests before deployment.

## Commit & PR Guidelines
- Follow Conventional Commits (`feat(scope): ...`, `fix(scope): ...`, `chore: ...`, `docs: ...`, `refactor: ...`, etc.).
- Keep commits focused; include related migrations or Swagger updates in the same commit when relevant.
- Before opening a PR, run `make linter` and `make test`, summarize behavioral changes, link relevant issues, and attach API diffs or screenshots when endpoints or docs change.

## Helper CLIs (under `cmd/kit`)
- `jwt`: Issue and inspect JWTs for local/dev.
- `http_route`: Generate or lint HTTP route scaffolding.
- `rdbms_seeder`: Seed relational databases with base data.
- `permission_seeder`: Seed Casbin roles, permissions, and grouping policies.

## Agent Notes & Precedence
- Scope: This file applies to the entire repository.
- Precedence: More deeply nested AGENTS.md files override this one for their subtree.
- Behavior: Be precise, safe, and helpful. Do not fix unrelated bugs. Ask for clarification when requirements are ambiguous.
- Formatting: Follow existing code style. Avoid intrusive refactors not required by the task.

## Conventional Commits 1.0.0

## 摘要

Conventional Commits 規範是在 commit message 之上的一種輕量級約定。它提供了一組簡單的規則來建立明確的提交歷史；這使得在其之上編寫自動化工具變得更加容易。這個約定與 SemVer 相吻合，透過在 commit message 中描述功能、修復和破壞性變更。

commit message 的結構應該如下：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

提交包含以下結構化元素，以向您的函式庫的使用者傳達意圖：

1.  **fix:** `fix` 類型的提交修補了您程式碼庫中的一個錯誤 (這對應於語意化版本中的 `PATCH`)。
2.  **feat:** `feat` 類型的提交為程式碼庫引入了一個新功能 (這對應於語意化版本中的 `MINOR`)。
3.  **BREAKING CHANGE:** 一個包含 `BREAKING CHANGE:` 註腳，或在類型/範圍後面附加 `!` 的提交，引入了一個破壞性的 API 變更 (對應於語意化版本中的 `MAJOR`)。一個 BREAKING CHANGE 可以是任何類型提交的一部分。
4.  除了 `fix:` 和 `feat:` 之外，也允許使用其他類型，例如 `@commitlint/config-conventional` (基於 Angular 約定) 推薦 `build:`, `chore:`, `ci:`, `docs:`, `style:`, `refactor:`, `perf:`, `test:` 等。
5.  除了 `BREAKING CHANGE: <description>` 之外，也可以提供其他註腳，並遵循類似 git trailer 格式的約定。

額外的類型並非 Conventional Commits 規範所強制要求的，並且在語意化版本中沒有隱含的影響 (除非它們包含 BREAKING CHANGE)。
可以為提交的類型提供一個範圍，以提供額外的上下文資訊，並包含在括號內，例如 `feat(parser): add ability to parse arrays`。

## 範例

### 包含描述和重大變更註腳的提交訊息

```
feat: allow provided config object to extend other configs

BREAKING CHANGE: `extends` key in config file is now used for extending other config files
```

### 包含 `!` 以提醒有重大變更的提交訊息

```
feat!: send an email to the customer when a product is shipped
```

### 包含範圍和 `!` 以提醒有重大變更的提交訊息

```
feat(api)!: send an email to the customer when a product is shipped
```

### 同時包含 `!` 和 BREAKING CHANGE 註腳的提交訊息

```
chore!: drop support for Node 6

BREAKING CHANGE: use JavaScript features not available in Node 6.
```

### 沒有內文的提交訊息

```
docs: correct spelling of CHANGELOG
```

### 包含範圍的提交訊息

```
feat(lang): add Polish language
```

### 包含多段落內文和多個註腳的提交訊息

```
fix: prevent racing of requests

Introduce a request id and a reference to latest request.
Dismiss incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.

Reviewed-by: Z
Refs: #123
```

## 規格

本文件中的關鍵詞「MUST」、「MUST NOT」、「REQUIRED」、「SHALL」、「SHALL NOT」、「SHOULD」、「SHOULD NOT」、「RECOMMENDED」、「MAY」和「OPTIONAL」應根據 RFC 2119 中的描述進行解釋。

1.  提交 **必須 (MUST)** 以一個類型作為前綴，該類型由一個名詞組成，如 `feat`、`fix` 等，後面跟著 **可選的 (OPTIONAL)** 範圍、**可選的 (OPTIONAL)** `!`，以及 **必需的 (REQUIRED)** 冒號和空格。
2.  當提交為您的應用程式或函式庫新增功能時，**必須 (MUST)** 使用 `feat` 類型。
3.  當提交代表對您的應用程式的錯誤修復時，**必須 (MUST)** 使用 `fix` 類型。
4.  在類型之後 **可以 (MAY)** 提供一個範圍。範圍 **必須 (MUST)** 由一個描述程式碼庫某個區塊的名詞組成，並用括號包圍，例如 `fix(parser):`
5.  在類型/範圍前綴後的冒號和空格之後，**必須 (MUST)** 立即跟著一個描述。描述是程式碼變更的簡短摘要，例如 *fix: array parsing issue when multiple spaces were contained in string.*
6.  在簡短描述之後，**可以 (MAY)** 提供一個更長的提交內文，提供有關程式碼變更的額外上下文資訊。內文 **必須 (MUST)** 在描述後空一行開始。
7.  提交內文是自由格式的，**可以 (MAY)** 由任意數量的以換行符分隔的段落組成。
8.  在內文之後，**可以 (MAY)** 提供一個或多個註腳，並在內文後空一行。每個註腳 **必須 (MUST)** 由一個單詞 token 組成，後面跟著 `:<space>` 或 `<space>#` 分隔符，然後是一個字串值 (這受到 git trailer 約定的啟發)。
9.  註腳的 token **必須 (MUST)** 使用 `-` 代替空白字元，例如 `Acked-by` (這有助於將註腳部分與多段落的內文區分開)。`BREAKING CHANGE` 是一個例外，它也 **可以 (MAY)** 作為 token 使用。
10. 註腳的值 **可以 (MAY)** 包含空格和換行符，當解析到下一個有效的註腳 token/分隔符對時，解析 **必須 (MUST)** 終止。
11. 重大變更 **必須 (MUST)** 在提交的類型/範圍前綴中標示，或作為註腳中的一個條目。
12. 如果作為註腳包含，重大變更 **必須 (MUST)** 由大寫文字 `BREAKING CHANGE` 組成，後面跟著一個冒號、空格和描述，例如 *BREAKING CHANGE: environment variables now take precedence over config files.*
13. 如果包含在類型/範圍前綴中，重大變更 **必須 (MUST)** 由一個緊接在 `:` 前面的 `!` 來表示。如果使用了 `!`，`BREAKING CHANGE:` **可以 (MAY)** 從註腳部分省略，提交描述 **應 (SHALL)** 用於描述重大變更。
14. 除了 `feat` 和 `fix` 之外，**可以 (MAY)** 在您的提交訊息中使用其他類型，例如 `docs: update ref docs.`
15. 構成 Conventional Commits 的資訊單元，實作者 **不得 (MUST NOT)** 將其視為區分大小寫，但 `BREAKING CHANGE` **必須 (MUST)** 為大寫。
16. 當在註腳中作為 token 使用時，`BREAKING-CHANGE` **必須 (MUST)** 與 `BREAKING CHANGE` 同義。

## 為什麼要使用 Conventional Commits？

*   自動產生 CHANGELOG。
*   自動決定語意化版本升級 (基於提交的類型)。
*   向團隊成員、公眾和其他利害關係人傳達變更的性質。
*   觸發建構和發布流程。
*   透過讓他們探索更有結構的提交歷史，使人們更容易為您的專案做出貢獻。

## 常見問題

**在初始開發階段，我應該如何處理提交訊息？**

我們建議您就像已經發布了產品一樣進行。通常會有人，即使是您的軟體開發同事，正在使用您的軟體。他們會想知道修復了什麼、破壞了什麼等等。

**提交標題中的類型是大寫還是小寫？**

任何大小寫都可以使用，但最好保持一致。

**如果提交符合多種提交類型，我該怎麼辦？**

盡可能地回去做多次提交。Conventional Commits 的部分好處是它能夠驅使我們做出更有組織的提交和 PR。

**這不會阻礙快速開發和快速迭代嗎？**

它阻礙的是以無組織的方式快速行動。它可以幫助您在多個專案和不同貢獻者之間長期快速地行動。

**Conventional Commits 會不會導致開發人員限制他們所做的提交類型，因為他們會按照所提供的類型來思考？**

Conventional Commits 鼓勵我們多做某些類型的提交，例如修復。除此之外，Conventional Commits 的靈活性允許您的團隊提出自己的類型並隨著時間的推移更改這些類型。

**這與 SemVer 有何關係？**

`fix` 類型的提交應轉換為 `PATCH` 版本。`feat` 類型的提交應轉換為 `MINOR` 版本。無論類型如何，在提交中帶有 `BREAKING CHANGE` 的提交都應轉換為 `MAJOR` 版本。

**我應該如何對我的 Conventional Commits 規範擴充進行版本控制，例如 `@jameswomack/conventional-commit-spec`？**

我們建議使用 SemVer 來發布您對此規範的擴充 (並鼓勵您進行這些擴充！)。

**如果我不小心使用了錯誤的提交類型怎麼辦？**

*   **當您使用了規範中的類型但不是正確的類型時，例如 `fix` 而不是 `feat`**：在合併或發布錯誤之前，我們建議使用 `git rebase -i` 來編輯提交歷史。發布後，清理工作將根據您使用的工具和流程而有所不同。
*   **當您使用了非規範的類型時，例如 `feet` 而不是 `feat`**：在最壞的情況下，如果提交不符合 Conventional Commits 規範，也不是世界末日。這僅意味著該提交將被基於該規範的工具所忽略。

**我的所有貢獻者都需要使用 Conventional Commits 規範嗎？**

不！如果您在 Git 上使用基於 squash 的工作流程，主要維護者可以在合併時清理提交訊息——這不會給臨時提交者增加任何工作量。一個常見的工作流程是讓您的 git 系統自動 squash 來自 pull request 的提交，並為主要維護者提供一個表單來輸入合併的正確 git 提交訊息。

**Conventional Commits 如何處理還原提交？**

還原程式碼可能很複雜：您是在還原多個提交嗎？如果您還原一個功能，下一個版本應該是補丁嗎？

Conventional Commits 沒有明確定義還原行為。相反，我們將其留給工具作者使用類型和註腳的靈活性來開發他們處理還原的邏輯。

一個建議是使用 `revert` 類型，以及一個引用正在被還原的提交 SHA 的註腳：

```
revert: let us never again speak of the noodle incident

Refs: 676104e, a215868
```

## 授權

Creative Commons - CC BY 3.0

## Commit Message 撰寫指導方針

撰寫 commit message 的首選方法是**結合「對話歷史」與「程式碼變更」**。這能讓開發助理最完整地理解變更的「意圖」和「實作」。如果沒有對話歷史，則退回使用備用方法。

### 首選方法：結合對話與 `git diff`

這是最推薦、也最簡單的作法。開發助理會自動分析整個互動過程以及最終的程式碼變更，產生最精準的 commit message。

**操作流程：**
1.  完成一項任務後，將變更加入暫存區 (`git add .`)。
2.  直接要求開發助理產生 commit message。

**範例 Prompt:**
> 「好了，我們完成了。請幫我產生一個 commit message。」

### 備用方法：僅使用 `git diff`

如果沒有相關的對話歷史（例如，你是離線完成開發，現在才要提交），開發助理也能單獨分析 `git diff` 的內容來產生 commit message。

**操作流程:**
1.  將變更加入暫存區 (`git add .`)。
2.  執行 `git diff --staged`。
3.  將完整的 `diff` 輸出結果提供給開發助理。

**範例 Prompt:**
> 「請幫我為以下的 `git diff` 產生 commit message：」
> ```diff
> diff --git a/src/utils/math.js b/src/utils/math.js
> index 6e9b2f7..8b4e6ad 100644
> --- a/src/utils/math.js
> +++ b/src/utils/math.js
> @@ -1,5 +1,9 @@
>  function add(a, b) {
>    return a + b;
>  }
> +
> +function subtract(a, b) {
> +  return a - b;
> +}
>  
> -module.exports = { add };
> +module.exports = { add, subtract };
> ```

### 語言與文法 (Language and Grammar)

- **使用英文 (Use English):** Commit messages 應全部使用英文撰寫，以利於國際協作和工具處理。
- **文法正確性 (Grammatical Correctness):** 在提交之前，請檢查您的 commit message，確保文法正確、拼寫無誤且標點符號使用得當。一個清晰、專業的 message 有助於他人理解您的變更。
