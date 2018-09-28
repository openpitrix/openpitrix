
default: secret-key.txt client-id.txt client-secret.txt

secret-key.txt:
	cat /dev/urandom | env LC_CTYPE=C tr -dc 'a-zA-Z0-9' | fold -w 99 | head -n 1 | tr -d '\n' > secret-key.txt


client-id.txt:
	cat /dev/urandom | env LC_CTYPE=C tr -dc 'a-zA-Z0-9' | fold -w 24 | head -n 1 | tr -d '\n' > client-id.txt


client-secret.txt:
	cat /dev/urandom | env LC_CTYPE=C tr -dc 'a-zA-Z0-9' | fold -w 24 | head -n 1 | tr -d '\n' > client-secret.txt