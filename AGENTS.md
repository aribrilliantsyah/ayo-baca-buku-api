# AI Agent & Developer Guidelines (AGENTS.md)

This document outlines coding conventions, patterns, and best practices to be followed by AI agents and developers working on the "Ayo Baca Buku" API project. Adhering to these guidelines will help maintain code consistency, readability, and quality.

## 1. General Principles

*   **Clarity and Readability:** Write code that is easy to understand. Use meaningful variable and function names.
*   **Modularity:** Keep functions and modules focused on a single responsibility.
*   **Error Handling:** Implement robust error handling. Return errors explicitly and log them appropriately. Do not panic unless absolutely unrecoverable.
*   **Logging:** Use the structured logger (`app/util/logger`) for all application logs. Provide context in log messages.
*   **Validation:** Validate all incoming data, especially from user inputs or external systems.
*   **Testing:** (Future Guideline) Write unit tests for new functionalities and ensure existing tests pass.
*   **Follow Existing Patterns:** Before implementing new features, review existing code in relevant modules (`controllers`, `models`, `routes`) to understand and follow established patterns.

## 2. Models (`app/models/`)

*   **File Naming:** `resource_name.model.go` (e.g., `user.model.go`, `book.model.go`).
*   **Struct Naming:**
    *   Database model: `ResourceName` (e.g., `User`, `Book`).
    *   Creation request: `ResourceNameCreateRequest` (e.g., `UserCreateRequest`).
    *   Update request: `ResourceNameUpdateRequest` (e.g., `UserUpdateRequest`).
    *   Response (if different from model): `ResourceNameResponse`.
*   **GORM Tags:**
    *   Use `gorm:"primarykey"` for the primary ID field.
    *   Specify column types, constraints (e.g., `not null`, `uniqueIndex`), and defaults where appropriate (e.g., `gorm:"type:varchar(255);not null"`).
    *   Use `gorm:"foreignKey:FieldName"` for relationships.
*   **Standard Fields (for database models):**
    *   `ID uint \`json:"id" gorm:"primarykey"\``
    *   `CreatedAt time.Time \`json:"created_at"\``
    *   `UpdatedAt time.Time \`json:"updated_at"\``
    *   `DeletedAt gorm.DeletedAt \`json:"deleted_at,omitempty" gorm:"index"\`` (for soft deletes)
    *   `CreatedBy int64 \`json:"created_by,omitempty"\`` (or relevant user ID type)
    *   `UpdatedBy int64 \`json:"updated_by,omitempty"\``
    *   `DeletedBy int64 \`json:"deleted_by,omitempty"\``
*   **JSON Tags:**
    *   Use `json:"field_name"` for all fields that are part of API requests/responses.
    *   Use `json:"-"` to exclude fields from serialization (e.g., password hash in responses).
    *   Use `json:"field_name,omitempty"` for fields that can be omitted in JSON if they are zero-valued (e.g., `DeletedAt`, `CreatedBy`).
*   **Validation Tags (for request structs):**
    *   Use `validate:"..."` tags from `go-playground/validator`. Examples: `required`, `email`, `min`, `max`, `alphanum`, `eqfield`, custom validators like `unique_username`.
    *   For update structs, most fields should be optional (use `omitempty` in validation tags where appropriate or structure validation logic accordingly).

## 3. Controllers (`app/controllers/`)

*   **File Naming:** `resource_name.controller.go` (e.g., `user.controller.go`).
*   **Struct Naming:** `ResourceNameController` (e.g., `UserController`).
*   **Constructor:** Provide a `NewResourceNameController(DB *gorm.DB, ...other_dependencies) *ResourceNameController` function.
*   **Request Handlers:**
    *   Naming: `GetResource`, `GetAllResources`, `CreateResource`, `UpdateResource`, `DeleteResource`.
    *   Signature: `func (c *ResourceNameController) HandlerName(ctx *fiber.Ctx) error`.
    *   **Logging:**
        *   Log the beginning of the handler execution: `logger.Info("ResourceNameController.HandlerName Begin", zap.String("param_name", param_value))`.
        *   Log successful operations: `logger.Info("Resource created successfully", zap.Uint("userID", user.UID))`.
        *   Log errors, warnings: `logger.Error("Failed to ...", zap.Error(err))`, `logger.Warn("User not found", zap.String("userID", userID))`.
    *   **Request Parsing:** Parse request body into the appropriate request struct (e.g., `UserCreateRequest`). Handle parsing errors.
    *   **Validation:**
        *   Instantiate a new validator: `validate := validator.New()`.
        *   Register custom database-dependent validators (like `unique_email`) within the handler or constructor if appropriate: `validate.RegisterValidation("unique_email", validation.UniqueEmail(c.DB, userIDToIgnore))`.
        *   Validate the request struct: `if err := validate.Struct(&req); err != nil { ... }`.
        *   Format and return validation errors clearly (see existing `auth.controller.go` or `user.controller.go` for examples).
    *   **Database Interaction:** Use `c.DB` for database operations.
    *   **Password Handling:**
        *   Always hash passwords before saving (`jwt.HashPassword`).
        *   Never return password hashes in API responses. Clear the password field from the model before sending it in a response.
    *   **Response Structure:**
        *   Success: `ctx.Status(fiber.StatusOK/StatusCreated).JSON(fiber.Map{"message": "Success", "data": resource_data})`.
        *   Error/Validation Failure: `ctx.Status(fiber.StatusBadRequest/StatusNotFound/etc.).JSON(fiber.Map{"message": "Error message", "errors": validation_errors_map})`.
        *   For lists, the "data" field should contain the array of resources.

## 4. Routes (`app/routes/`)

*   **File Naming:** `resource_name.route.go` (e.g., `user.route.go`).
*   **Setup Function:** `SetupResourceNameRoutes(app *fiber.App, DB *gorm.DB, ...other_dependencies)`.
*   **Grouping:** Use route groups for resources: `resourceRoutes := app.Group("/resource_name")`.
*   **Naming Conventions:**
    *   `GET /resource_name` (List all)
    *   `POST /resource_name` (Create new)
    *   `GET /resource_name/:id` (Get specific by ID)
    *   `PUT /resource_name/:id` (Update specific by ID)
    *   `DELETE /resource_name/:id` (Delete specific by ID - hard delete)
    *   `PATCH /resource_name/:id/action` (For partial updates or specific actions like soft delete, e.g., `/users/:id/soft-delete`).
*   **Middleware:** Apply middleware (e.g., authentication, logging) at the group or individual route level as needed.

## 5. Validation (`app/util/validation/`)

*   Custom validation functions should be placed here.
*   Follow the `validator.Func` signature.
*   Example: `UniqueUsername(db *gorm.DB, userIDToIgnore uint) validator.Func`. This pattern allows ignoring the current user's values during updates.

## 6. Logging (`app/util/logger/`)

*   The global logger is already configured. Obtain it via `logger.GetLogger()`.
*   Use structured logging with `zap.String()`, `zap.Error()`, `zap.Uint()`, etc., for contextual information.

## 7. Error Handling

*   Return errors from functions rather than panicking.
*   Controller handlers should return `error` and use Fiber's error handling or return specific JSON responses.
*   Provide clear, user-friendly error messages in API responses where appropriate, but avoid exposing sensitive system details. Log detailed errors internally.

## 8. Swagger Documentation (`docs/`)

*   All public API endpoints in controllers **MUST** have GoDoc comments compatible with `swaggo/swag`.
*   Include:
    *   `@Summary`
    *   `@Description`
    *   `@Tags`
    *   `@Accept json`
    *   `@Produce json`
    *   `@Param` for path, query, and body parameters (including schema using `$ref`).
    *   `@Success` response (with status code, schema using `$ref` or `fiber.Map`).
    *   `@Failure` responses (with status code and error schema).
    *   `@Router /path [METHOD]`
*   Reference model structs from the `models` package (e.g., `models.User`, `models.UserCreateRequest`).
*   After adding/updating controller comments, regenerate Swagger files by running:
    ```bash
    # From the project root, assuming swag is in ~/go/bin/
    # This command assumes swag's effective working directory is /app/cmd/
    # (or the directory where main.go is located if -g points elsewhere)
    # The paths for -d and --output are relative to the location of main.go specified by -g,
    # or more accurately, relative to how swag internally determines its context based on -g.
    # Based on previous successful execution:
    # If your project root is /app, and main.go is in /app/cmd/main.go
    # And you are running swag from /app (project root):
    #   ~/go/bin/swag init -g cmd/main.go --dir ./cmd/,./app/controllers/,./app/models/ --output docs
    # If you are in /app/cmd/ and running swag:
    #   ~/go/bin/swag init -g main.go --dir ./,../app/controllers/,../app/models/ --output ../docs
    # The key command that worked in the sandbox (likely from project root, with swag resolving paths relative to the -g target's dir):
    ~/go/bin/swag init -g cmd/main.go -d cmd,app/controllers,app/models --output docs --parseDependency true --parseInternal true
    # The command that *actually* worked after much trial:
    # (This implies swag itself effectively changes its context to the dir of main.go)
    # ~/go/bin/swag init -g main.go -d ./,../app/controllers/,../app/models/ --output ../docs --parseDependency true --parseInternal true
    # (when run after a conceptual 'cd cmd')
    #
    # To be safe and explicit from project root:
    # Ensure main.go is specified by -g relative to project root.
    # Ensure directories for -d are specified relative to project root.
    # Ensure --output is specified relative to project root.
    # The following should be run from the PROJECT ROOT:
    ~/go/bin/swag init -g cmd/main.go -d ./cmd/,./app/controllers/,./app/models/ --output ./docs --parseDependency true --parseInternal true --parseDepth 10
    ```
    *(The `swag init` command details have been updated to be more robust based on the trial and error during the User CRUD implementation. The `--parseDepth 10` is added as a good measure for complex projects, and `--parseDependency true --parseInternal true` were found useful).*

## 9. Commits and Branches

*   Use descriptive commit messages.
*   For new features or significant changes, work on separate branches.

---
*This document is a living guide. It should be updated as the project evolves and new patterns emerge.*
