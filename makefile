testclient:
	@echo "Starting test client"
	node ./Chat_Server/TestClient/test_client.js

dev:
	@echo "Starting dev build"
	./Chat_Server/scripts/gen-self-signed-certs.sh
	docker compose up -d --build --force-recreate

dev_down:
	@echo "Tearing dev build down"
	docker compose down