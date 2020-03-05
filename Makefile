import-passwords:
	echo "Downloading CrackStation's Password Cracking List"
	curl -L -o ./import-passwords/crackstation-human-only.txt.gz https://crackstation.net/files/crackstation-human-only.txt.gz

setup-passwords:
	cd ./import-passwords && go build && ./import-passwords

maildev:
	docker run --rm -it -p 2525:1025 -p 8080:1080 djfarrelly/maildev --incoming-user=test --incoming-pass=test

run:
	docker run -d -p 2525:1025 -p 8080:1080 --name zendesk-maildev djfarrelly/maildev --incoming-user=test --incoming-pass=test
	-rm ./zendesk-product_security_challenge
	rice embed-go
	go build
	-./zendesk-product_security_challenge
	docker rm -f zendesk-maildev

pkg:
	echo "Building macOS amd64 zip..."
	-rm ./zendesk-product_security_challenge ./build.zip ./build.macos.amd64.zip ./build.linux.amd64.zip ./build.win.amd64.zip ./build.companions.zip
	rice embed-go
	GOOS=darwin GOARCH=amd64 go build -o zendesk-product_security_challenge
	zip -r ./build.macos.amd64.zip -X ./zendesk-product_security_challenge ./.env ./common-passwords.db
	echo "Building linux amd64 zip..."
	-rm zendesk-product_security_challenge
	GOOS=linux GOARCH=amd64 go build -o zendesk-product_security_challenge
	zip -r ./build.linux.amd64.zip -X ./zendesk-product_security_challenge ./.env ./common-passwords.db
	echo "Building windows amd64 zip..."
	-rm zendesk-product_security_challenge
	GOOS=windows GOARCH=amd64 go build -o zendesk-product_security_challenge.exe
	zip -r ./build.win.amd64.zip -X ./zendesk-product_security_challenge.exe ./.env ./common-passwords.db
	echo "Building companion zip for source building..."
	-rm zendesk-product_security_challenge.exe
	zip -r ./build.companions.zip -X ./.env ./common-passwords.db