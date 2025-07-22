# TODO
# default probably shouldn't be dev build
dev:
	@echo "Starting dev build"
	./Chat_Server/scripts/gen-self-signed-certs.sh
	docker compose up -d --build --force-recreate

zoo:
	@echo "Starting zirc_zoo server"
	docker compose --profile multi up zirc_zoo -d --build --force-recreate

animal:
	@echo "Starting zirc_zoo server"
	docker compose --profile multi up zirc_animalhouse -d --build --force-recreate

devdown:
	@echo "Tearing dev build down"
	docker compose down

testclient:
	@echo "Starting test client"
	node ./Chat_Server/TestClient/test_client.js