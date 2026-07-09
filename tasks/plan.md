# Implementation Plan: Migrate from Swagger UI to Scalar

## Overview
Replace the current Swagger UI (swaggo/gin-swagger) with Scalar API reference UI. The existing swaggo annotations on controllers will be kept for OpenAPI spec generation, but the UI layer will be switched to Scalar.

## Architecture Decisions
- **Keep swaggo annotations** - The `@Summary`, `@Param`, `@Success` etc. annotations on controllers are useful for generating the OpenAPI spec. We keep these.
- **Keep swaggo CLI for spec generation** - Continue using `swag init` to generate `docs/swagger.json`.
- **Replace swagger UI with Scalar** - Remove `gin-swagger` and `swaggo/files` dependencies. Serve the spec via Scalar's CDN-based UI.
- **Endpoint: `/docs`** - Change from `/swagger/*any` to `/docs` for Scalar UI.

## What changes
1. **Remove Swagger UI deps** from `routes.go` (gin-swagger, swaggo/files imports)
2. **Add Scalar HTML endpoint** in `routes.go` that serves Scalar's CDN UI pointing to the swagger.json
3. **Serve swagger.json statically** so Scalar can load it
4. **Remove unused deps** from go.mod (swaggo/files, gin-swagger) after replacing
5. **Keep docs/ folder** and swaggo annotations - they're still needed for spec generation

## Task List

### Phase 1: Foundation
- [ ] Task 1: Create spec file and plan (this task)

### Phase 2: Migration
- [ ] Task 2: Update routes.go - remove swagger UI imports, add Scalar HTML endpoint and spec endpoint
- [ ] Task 3: Clean up go.mod - remove unused swaggo dependencies
- [ ] Task 4: Verify build succeeds

### Phase 3: Verification
- [ ] Checkpoint: Build passes, no compilation errors