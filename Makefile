watch_cli_api:
	wgo run ./cmd/api/ & wgo run ./cmd/cli/ wait

watch_api:
	wgo run ./cmd/api/

watch_cli:
	wgo run ./cmd/cli/

