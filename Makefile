db/seed:
	docker compose --profile seed run --rm seed

test/auth:
	API_BASE_URL=http://localhost:8080 npx vitest run test/tests/auth.test.ts

test/authRequired:
	API_BASE_URL=http://localhost:8080 npx vitest run test/tests/authRequired.test.ts

test/user:
	API_BASE_URL=http://localhost:8080 npx vitest run test/tests/user.test.ts

test/projects:
	API_BASE_URL=http://localhost:8080 npx vitest run test/tests/projects.test.ts

test/tasks:
	API_BASE_URL=http://localhost:8080 npx vitest run test/tests/tasks.test.ts
