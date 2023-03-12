.PHONY: create-keypair

PWD = $(shell pwd)
AUTHTPATH = $(PWD)/auth-service
MIGRATIONPATH = $(AUTHTPATH)/migration

create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(AUTHTPATH)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(AUTHTPATH)/rsa_private_$(ENV).pem -pubout -out $(AUTHTPATH)/rsa_public_$(ENV).pem

migrate-up:
	@echo "Up migration"
	cd $(MIGRATIONPATH)
	goose postgres "user=pguser password=pgpwd dbname=testdb sslmode=disable" up

migrate-down:
	@echo "Down migration"
	cd $(MIGRATIONPATH)
	goose postgres "user=pguser password=pgpwd dbname=testdb sslmode=disable" down