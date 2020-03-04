install:
	yarn global add maildev
	echo "Downloading CrackStation's Password Cracking List"
	curl -L -o ./import-passwords/crackstation-human-only.txt.gz https://crackstation.net/files/crackstation-human-only.txt.gz

setup_pwdb:
	cd ./import-passwords && go build && ./import-passwords

run:
	maildev -s 2525 -w 8080 &
	-rm ./zendesk-product_security_challenge
	go build
	-./zendesk-product_security_challenge
	kill $$(ps aux | grep maildev | grep -v grep | awk '{print $$2}')

