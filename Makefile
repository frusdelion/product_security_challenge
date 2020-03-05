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

