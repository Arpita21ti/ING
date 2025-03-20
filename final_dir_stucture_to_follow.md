# The main dir structure for scale-able design

```text
TNP_RGPV_SERVER/
├── cmd/                                 # Command-line entry points
│   ├── api/                             # REST API server entry point
│   │   └── main.go                      # Server initialization, config loading, router setup
│   ├── worker/                          # Background job processor
│   │   └── main.go                      # Worker initialization, queue consumption
│   ├── migrations/                      # Database migration tool
│   │   └── main.go                      # DB migration execution utility
│   └── seed/                             # Seed command entry point
│      └── main.go                       # CLI for running seeds
├── seeds/                              # Database seeding scripts
│   ├── development/                    # Development environment seeds
│   │   ├── 001_roles.go                # Basic role definitions
│   │   ├── 002_admin_users.go          # Admin user accounts
│   │   └── 003_sample_data.go          # Sample data for development
│   ├── testing/                        # Test environment seeds
│   │   ├── test_users.go               # Test user accounts
│   │   └── test_scenarios.go           # Specific test scenarios
│   ├── production/                     # Production environment seeds
│   │   ├── initial_roles.go            # Initial roles for production
│   │   └── default_settings.go         # Default system settings
│   └── seed_runner.go                  # Script to run seeds in sequence
├── internal/                            # Private application code
|   ├── security/                             # Security framework
|   │   └── scanner/                          # Security scanning
|   │       └── integration_scanner.go        # Scans for security issues in integrations
|   ├── analytics/                            # Analytics framework
|   │   ├── usage/                            # Usage analytics
|   │   │   └── integration_usage.go          # Tracks integration usage patterns
|   │   └── cost/                             # Cost analysis
|   │       └── service_cost_analyzer.go      # Analyzes costs of external services
|   ├── tenant/                               # Multi-tenant support. If different departments or faculties need different integration configurations.
|   │   └── integration_config.go             # Tenant-specific integration settings
|   ├── compliance/                           # Compliance management
|   │   ├── gdpr/                             # GDPR compliance handling
|   │   │   └── data_processor.go             # Handles data subject rights
|   │   ├── data_sharing/                     # Data sharing controls
|   │   │   └── consent_manager.go            # Manages user consent for data sharing
|   │   |── audit/                            # Audit logging for compliance
│   |   └── integration_auditor.go        # Audits interactions with external services
│   ├── domain/                          # Core domain models & business logic
│   │   ├── student/                     # Student domain
│   │   │   ├── model.go                 # Student entity definitions, value objects
│   │   │   ├── repository.go            # Student repository interface
│   │   │   └── service.go               # Student business logic, validation rules
│   │   ├── quiz/                        # Quiz management domain
│   │   │   ├── model.go                 # Quiz entities and relationships
│   │   │   ├── question/                # Question type hierarchy
│   │   │   │   ├── base.go              # Question interface and common functionality
│   │   │   │   ├── mcq.go               # Multiple choice question implementation
│   │   │   │   ├── fill_blank.go        # Fill-in-the-blank question implementation
│   │   │   │   └── true_false.go        # True/False question implementation
│   │   │   ├── repository.go            # Quiz repository interface
│   │   │   └── service.go               # Quiz business logic, scoring rules
│   │   ├── event/                       # Event management domain
│   │   │   ├── model.go                 # Event entities, scheduling entities
│   │   │   ├── repository.go            # Event repository interface
│   │   │   └── service.go               # Event business logic, validation rules
│   │   ├── coordinator/                 # Coordinator domain
│   │   │   ├── model.go                 # Coordinator entities and relationships
│   │   │   ├── repository.go            # Coordinator repository interface
│   │   │   └── service.go               # Coordinator-specific business logic
│   │   ├── role/                        # Role-based permissions system
│   │   │   ├── model.go                 # Role definitions, permission entities
│   │   │   ├── repository.go            # Role repository interface
│   │   │   └── service.go               # Permission checking, role assignment logic
|   |   ├── integration/                   # Integration interfaces
|   |   │   ├── notification_provider.go   # Abstract notification service
|   |   │   ├── file_storage_provider.go   # Abstract file storage
|   |   │   ├── ai_provider.go             # Abstract AI service
|   |   │   └── design_provider.go         # Abstract design service
|   │   ├── ai/                                # AI features domain
|   │   │   ├── model.go                       # AI model entities
|   │   │   ├── repository.go                  # AI model repository
|   │   │   └── service.go                     # AI services (using Hugging Face)
|   │   ├── design/                            # Design features domain
|   │   │   ├── model.go                       # Design entities
|   │   │   ├── repository.go                  # Design repository
|   │   │   └── service.go                     # Design services (using Canva)
│   │   └── [other domains]              # Additional business domains
│   ├── api/                             # API layer
│   │   ├── rest/                        # REST API
│   │   │   ├── handler/                 # Request handlers
│   │   │   │   ├── student/             # Student-related endpoint handlers
│   │   │   │   │   ├── auth_handler.go  # Student authentication endpoints
│   │   │   │   │   └── profile_handler.go # Student profile management endpoints
│   │   │   │   ├── quiz/                # Quiz-related endpoint handlers
│   │   │   │   │   ├── quiz_handler.go  # Quiz CRUD operations
│   │   │   │   │   └── attempt_handler.go # Quiz attempt submission/retrieval
│   │   │   │   └── [other domains]      # Additional domain endpoint handlers
│   │   │   ├── middleware/              # HTTP middleware
│   │   │   │   ├── auth.go              # Authentication/authorization middleware
│   │   │   │   ├── logging.go           # Request logging middleware
│   │   │   │   └── rate_limit.go        # Rate limiting middleware
│   │   │   ├── dto/                     # Data transfer objects
│   │   │   │   ├── request/             # Request models for API endpoints
│   │   │   │   └── response/            # Response models for API endpoints
│   │   │   |── router/                  # Router configuration
│   │   │   |   ├── router.go            # Main router setup
│   │   │   |   ├── student_routes.go    # Student endpoint registration
│   │   │   |   └── quiz_routes.go       # Quiz endpoint registration
|   │   │   ├── ai/                        # AI feature endpoints
|   │   │   └── design/                    # Design feature endpoints
│   │   |── graphql/                     # Future GraphQL API
│   │   |   ├── resolver/                # GraphQL resolvers for each domain
│   │   |   ├── schema/                  # GraphQL schema definitions
│   │   |  └── directive/               # Custom GraphQL directives
|   │   └── webhook/                           # Webhook handlers
|   │       ├── google_webhook.go              # Google webhook handlers
|   │       ├── onesignal_webhook.go           # OneSignal webhook handlers
|   │       └── canva_webhook.go               # Canva webhook handlers
│   ├── infrastructure/                  # External dependencies
|   │   ├── migration/                        # Feature migration framework
|   │   │   └── service_switcher.go           # Handles switching between service providers
|   │   ├── offline/                          # Offline operation support
|   │   │   └── offline_queue.go              # Queues operations for when online
|   │   ├── discovery/                        # Service discovery
|   │   │   └── service_registry.go           # Dynamic service endpoint discovery
|   │   ├── quota/                            # API quota management
|   │   │   ├── limit_tracker.go              # Tracks API usage against limits
|   │   │   └── throttling_service.go         # Implements request throttling
|   │   ├── monitoring/                       # Health monitoring
|   │   │   ├── integration_health.go         # Integration health checks
|   │   │   └── alerts/                       # Alert system for service issues
|   │   │       └── service_degradation.go    # Handles service degradation alerts
|   │   ├── fallback/
|   |   │   ├── notification_fallback.go       # Fallback notification delivery
|   |   │   └── storage_fallback.go            # Fallback storage options
|   │   ├── sdk/                               # SDK wrappers
|   │   │   ├── google_sdk_wrapper.go          # Google SDK wrapper
|   │   │   └── huggingface_sdk_wrapper.go     # Hugging Face SDK wrapper
│   │   ├── database/                    # Database adapters
│   │   │   ├── postgres/                # PostgreSQL adapter
│   │   │   │   ├── connection.go        # Connection pool management
│   │   │   │   ├── migrations/          # Postgres-specific migrations
│   │   │   │   └── repositories/        # Concrete repository implementations
│   │   │   │       ├── student_repo.go  # Student repository PostgreSQL implementation
│   │   │   │       └── quiz_repo.go     # Quiz repository PostgreSQL implementation
│   │   │   └── mongodb/                 # MongoDB adapter
│   │   │       ├── connection.go        # MongoDB connection management
│   │   │       └── repositories/        # MongoDB repository implementations
│   │   ├── auth/                        # Auth providers
│   │   │   └── jwt_provider.go          # JWT token generation/validation
|   │   │   ├── oauth_client.go                # OAuth client implementation
|   │   │   └── api_key_manager.go             # API key management
│   │   ├── email/                       # Email service
│   │   │   ├── smtp_provider.go         # SMTP email sender implementation
│   │   │   └── template_engine.go       # Email template rendering
│   │   ├── storage/                     # File storage
│   │   │   ├── local_storage.go         # Local filesystem storage
│   │   │   └── s3_storage.go            # S3 compatible storage
│   │   |── cache/                       # Caching layer
│   │   │    ├── redis_cache.go           # Redis cache implementation
│   │   │    └── local_cache.go           # In-memory cache implementation
|   |   └── integration/                       # Third-party service integrations
|   │   │   ├── versioning/                   # API version management
|   │   │   │   └── version_manager.go        # Handles API version transitions
|   |   |   ├── google/                        # Google services integration
|   |   |   │   ├── auth.go                    # Google OAuth implementation
|   |   |   │   ├── calendar.go                # Google Calendar integration
|   |   |   │   └── drive.go                   # Google Drive integration
|   |   |   ├── onesignal/                     # OneSignal push notifications
|   |   |   │   └── notification_client.go     # OneSignal API client
|   |   |   ├── huggingface/                   # Hugging Face AI models
|   |   |   │   ├── client.go                  # API client for Hugging Face
|   |   |   │   └── model_manager.go           # AI model management
|   |   |   |── canva/                         # Canva design integration
|   |   |   |   └── design_client.go           # Canva API client
|   |   │   └── factory/
|   |   │       ├── notification_factory.go        # Creates notification providers
|   |   │       ├── storage_factory.go             # Creates storage providers
|   |   │       └── ai_factory.go                  # Creates AI service providers
│   ├── worker/                          # Background jobs
│   │   ├── email/                       # Email sending workers
│   │   │   └── notification_worker.go   # Sends email notifications 
│   │   ├── notification/                # Push notification workers
│   │   │   └── push_worker.go           # Sends push notifications
│   │   └── analytics/                   # Analytics processing
│   │       └── report_worker.go         # Generates periodic reports
│   ├── config/                          # Configuration
│   │   ├── app_config.go                # Application configuration loader
│   │   ├── environment.go               # Environment variable handling
│   │   |── feature_flags.go             # Feature toggle support
|   │   ├── integration_config.go        # Third-party integration config
|   │   └── credentials.go               # Secure credential management
│   ├── common/                          # Cross-cutting concerns
│   │   ├── validator/                   # Validation logic
│   │   │   ├── password_validator.go    # Password strength validation
│   │   │   └── input_sanitizer.go       # Input sanitization functions
│   │   ├── errors/                      # Error types & handling
│   │   │   ├── domain_errors.go         # Domain-specific error types
│   │   │   └── http_errors.go           # HTTP error handling utilities
│   │   └── utils/                       # Shared utilities
│   │       ├── pagination.go            # Pagination helpers
│   │       └── security.go              # Security-related helpers
│   └── eventbus/                        # Event-driven architecture
│       ├── events/                      # Event definitions
│       │   ├── student_events.go        # Student-related events
│       │   └── quiz_events.go           # Quiz-related events
│       ├── publisher/                   # Event publishers
│       │   └── kafka_publisher.go       # Kafka implementation of event publishing
│       └── subscriber/                  # Event subscribers
│           └── kafka_subscriber.go      # Kafka implementation of event consumption
├── pkg/                                 # Public packages
│   ├── logger/                          # Logging utility
│   │   └── logger.go                    # Structured logging implementation
│   ├── pagination/                      # Reusable pagination
│   │   └── paginator.go                 # Generic pagination implementation
│   └── security/                        # Security utilities
│       └── password.go                  # Password hashing & verification
├── migrations/                          # Database migrations
│   ├── postgres/                        # PostgreSQL migrations
│   │   ├── 000001_create_students.up.sql   # Migration to create students table
│   │   └── 000001_create_students.down.sql # Rollback for students table
│   └── mongodb/                         # MongoDB migrations/seed scripts
├── scripts/                             # CI/CD, deployment scripts
│   ├── deploy/                          # Deployment scripts
│   │   ├── kubernetes/                  # K8s deployment manifests
│   │   └── docker-compose.yml           # Docker compose for local dev
│   ├── ci/                              # CI pipeline scripts
│   │   ├── build.sh                     # Build script for CI
│   │   └── test.sh                      # Test execution script
│   └── local/                           # Local development scripts
│       └── setup.sh                     # Local dev environment setup
├── docs/                                # Documentation
|   ├── integrations/                         # Integration documentation
|   │   ├── setup_guides/                     # Setup guides for each integration
|   │   └── api_references/                   # API reference for integration points
│   ├── api/                             # API documentation
│   │   └── openapi.yaml                 # OpenAPI/Swagger specification
│   ├── architecture/                    # Architecture docs
│   │   └── domain_model.md              # Domain model documentation
│   └── guides/                          # Guides for developers
│       └── getting_started.md           # Developer onboarding guide
├── test/                                # Test helpers and fixtures
│   ├── integration/                     # Integration test helpers
│   │   |── db_test_helper.go            # Database test utilities
|   │   └── external/                         # External service tests
|   │       ├── google_test.go                # Google API integration tests
|   │       ├── onesignal_test.go             # OneSignal integration tests
|   │       └── mock_services/                # Mock implementations for testing
|   │           └── mock_google_api.go        # Mock Google API responses
│   ├── mocks/                           # Mock implementations
│   │   └── repositories/                # Mock repository implementations
│   └── fixtures/                        # Test data fixtures
│       └── students.json                # Sample student data for tests
├── .env.example                         # Example environment variables
├── .gitignore                           # Git ignore file
├── go.mod                               # Go module definition
└── go.sum                               # Go module checksums
```
