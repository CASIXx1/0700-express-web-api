db/seed:
	docker compose --profile seed run --rm seed

test/auth:
	API_BASE_URL=http://localhost:8080 npx vitest run test/tests/auth.test.ts