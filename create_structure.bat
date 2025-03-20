@echo off

REM cmd directory
mkdir cmd\api cmd\worker cmd\migrations cmd\seed
type nul > cmd\api\main.go
type nul > cmd\worker\main.go
type nul > cmd\migrations\main.go
type nul > cmd\seed\main.go

REM seeds directory
mkdir seeds\development seeds\testing seeds\production
type nul > seeds\development\001_roles.go
type nul > seeds\development\002_admin_users.go
type nul > seeds\development\003_sample_data.go
type nul > seeds\testing\test_users.go
type nul > seeds\testing\test_scenarios.go
type nul > seeds\production\initial_roles.go
type nul > seeds\production\default_settings.go
type nul > seeds\seed_runner.go

REM internal directory - security, analytics, tenant, compliance
mkdir internal\security\scanner internal\analytics\usage internal\analytics\cost internal\tenant internal\compliance\gdpr internal\compliance\data_sharing internal\compliance\audit
type nul > internal\security\scanner\integration_scanner.go
type nul > internal\analytics\usage\integration_usage.go
type nul > internal\analytics\cost\service_cost_analyzer.go
type nul > internal\tenant\integration_config.go
type nul > internal\compliance\gdpr\data_processor.go
type nul > internal\compliance\data_sharing\consent_manager.go
type nul > internal\compliance\integration_auditor.go

REM internal directory - domain
mkdir internal\domain\student internal\domain\quiz\question internal\domain\event internal\domain\coordinator internal\domain\role internal\domain\integration internal\domain\ai internal\domain\design
type nul > internal\domain\student\model.go
type nul > internal\domain\student\repository.go
type nul > internal\domain\student\service.go
type nul > internal\domain\quiz\model.go
type nul > internal\domain\quiz\repository.go
type nul > internal\domain\quiz\service.go
type nul > internal\domain\quiz\question\base.go
type nul > internal\domain\quiz\question\mcq.go
type nul > internal\domain\quiz\question\fill_blank.go
type nul > internal\domain\quiz\question\true_false.go
type nul > internal\domain\event\model.go
type nul > internal\domain\event\repository.go
type nul > internal\domain\event\service.go
type nul > internal\domain\coordinator\model.go
type nul > internal\domain\coordinator\repository.go
type nul > internal\domain\coordinator\service.go
type nul > internal\domain\role\model.go
type nul > internal\domain\role\repository.go
type nul > internal\domain\role\service.go
type nul > internal\domain\integration\notification_provider.go
type nul > internal\domain\integration\file_storage_provider.go
type nul > internal\domain\integration\ai_provider.go
type nul > internal\domain\integration\design_provider.go
type nul > internal\domain\ai\model.go
type nul > internal\domain\ai\repository.go
type nul > internal\domain\ai\service.go
type nul > internal\domain\design\model.go
type nul > internal\domain\design\repository.go
type nul > internal\domain\design\service.go

REM internal directory - api
mkdir internal\api\rest\handler\student internal\api\rest\handler\quiz internal\api\rest\middleware internal\api\rest\dto\request internal\api\rest\dto\response internal\api\rest\router internal\api\rest\ai internal\api\rest\design internal\api\graphql\resolver internal\api\graphql\schema internal\api\graphql\directive internal\api\webhook
type nul > internal\api\rest\handler\student\auth_handler.go
type nul > internal\api\rest\handler\student\profile_handler.go
type nul > internal\api\rest\handler\quiz\quiz_handler.go
type nul > internal\api\rest\handler\quiz\attempt_handler.go
type nul > internal\api\rest\middleware\auth.go
type nul > internal\api\rest\middleware\logging.go
type nul > internal\api\rest\middleware\rate_limit.go
type nul > internal\api\rest\router\router.go
type nul > internal\api\rest\router\student_routes.go
type nul > internal\api\rest\router\quiz_routes.go
type nul > internal\api\webhook\google_webhook.go
type nul > internal\api\webhook\onesignal_webhook.go
type nul > internal\api\webhook\canva_webhook.go

REM internal directory - infrastructure
mkdir internal\infrastructure\migration internal\infrastructure\offline internal\infrastructure\discovery internal\infrastructure\quota internal\infrastructure\monitoring\alerts internal\infrastructure\fallback internal\infrastructure\sdk internal\infrastructure\database\postgres\migrations internal\infrastructure\database\postgres\repositories internal\infrastructure\database\mongodb\repositories internal\infrastructure\auth internal\infrastructure\email internal\infrastructure\storage internal\infrastructure\cache internal\infrastructure\integration\versioning internal\infrastructure\integration\google internal\infrastructure\integration\onesignal internal\infrastructure\integration\huggingface internal\infrastructure\integration\canva internal\infrastructure\integration\factory
type nul > internal\infrastructure\migration\service_switcher.go
type nul > internal\infrastructure\offline\offline_queue.go
type nul > internal\infrastructure\discovery\service_registry.go
type nul > internal\infrastructure\quota\limit_tracker.go
type nul > internal\infrastructure\quota\throttling_service.go
type nul > internal\infrastructure\monitoring\integration_health.go
type nul > internal\infrastructure\monitoring\alerts\service_degradation.go
type nul > internal\infrastructure\fallback\notification_fallback.go
type nul > internal\infrastructure\fallback\storage_fallback.go
type nul > internal\infrastructure\sdk\google_sdk_wrapper.go
type nul > internal\infrastructure\sdk\huggingface_sdk_wrapper.go
type nul > internal\infrastructure\database\postgres\connection.go
type nul > internal\infrastructure\database\postgres\repositories\student_repo.go
type nul > internal\infrastructure\database\postgres\repositories\quiz_repo.go
type nul > internal\infrastructure\database\mongodb\connection.go
type nul > internal\infrastructure\auth\jwt_provider.go
type nul > internal\infrastructure\auth\oauth_client.go
type nul > internal\infrastructure\auth\api_key_manager.go
type nul > internal\infrastructure\email\smtp_provider.go
type nul > internal\infrastructure\email\template_engine.go
type nul > internal\infrastructure\storage\local_storage.go
type nul > internal\infrastructure\storage\s3_storage.go
type nul > internal\infrastructure\cache\redis_cache.go
type nul > internal\infrastructure\cache\local_cache.go
type nul > internal\infrastructure\integration\versioning\version_manager.go
type nul > internal\infrastructure\integration\google\auth.go
type nul > internal\infrastructure\integration\google\calendar.go
type nul > internal\infrastructure\integration\google\drive.go
type nul > internal\infrastructure\integration\onesignal\notification_client.go
type nul > internal\infrastructure\integration\huggingface\client.go
type nul > internal\infrastructure\integration\huggingface\model_manager.go
type nul > internal\infrastructure\integration\canva\design_client.go
type nul > internal\infrastructure\integration\factory\notification_factory.go
type nul > internal\infrastructure\integration\factory\storage_factory.go
type nul > internal\infrastructure\integration\factory\ai_factory.go

REM internal directory - worker
mkdir internal\worker\email internal\worker\notification internal\worker\analytics
type nul > internal\worker\email\notification_worker.go
type nul > internal\worker\notification\push_worker.go
type nul > internal\worker\analytics\report_worker.go

REM internal directory - config
mkdir internal\config
type nul > internal\config\app_config.go
type nul > internal\config\environment.go
type nul > internal\config\feature_flags.go
type nul > internal\config\integration_config.go
type nul > internal\config\credentials.go

REM internal directory - common
mkdir internal\common\validator internal\common\errors internal\common\utils
type nul > internal\common\validator\password_validator.go
type nul > internal\common\validator\input_sanitizer.go
type nul > internal\common\errors\domain_errors.go
type nul > internal\common\errors\http_errors.go
type nul > internal\common\utils\pagination.go
type nul > internal\common\utils\security.go

REM internal directory - eventbus
mkdir internal\eventbus\events internal\eventbus\publisher internal\eventbus\subscriber
type nul > internal\eventbus\events\student_events.go
type nul > internal\eventbus\events\quiz_events.go
type nul > internal\eventbus\publisher\kafka_publisher.go
type nul > internal\eventbus\subscriber\kafka_subscriber.go

REM pkg directory
mkdir pkg\logger pkg\pagination pkg\security
type nul > pkg\logger\logger.go
type nul > pkg\pagination\paginator.go
type nul > pkg\security\password.go

REM migrations directory
mkdir migrations\postgres migrations\mongodb
type nul > migrations\postgres\000001_create_students.up.sql
type nul > migrations\postgres\000001_create_students.down.sql

REM scripts directory
mkdir scripts\deploy\kubernetes scripts\ci scripts\local
type nul > scripts\deploy\docker-compose.yml
type nul > scripts\ci\build.sh
type nul > scripts\ci\test.sh
type nul > scripts\local\setup.sh

REM docs directory
mkdir docs\integrations\setup_guides docs\integrations\api_references docs\api docs\architecture docs\guides
type nul > docs\api\openapi.yaml
type nul > docs\architecture\domain_model.md
type nul > docs\guides\getting_started.md

REM test directory
mkdir test\integration\external\mock_services test\mocks\repositories test\fixtures
type nul > test\integration\db_test_helper.go
type nul > test\integration\external\google_test.go
type nul > test\integration\external\onesignal_test.go
type nul > test\integration\external\mock_services\mock_google_api.go
type nul > test\fixtures\students.json

REM Add root files
type nul > .env.example
type nul > .gitignore

echo Directory structure created successfully!