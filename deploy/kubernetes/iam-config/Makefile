
default: secret-key.txt client-id.txt client-secret.txt admin-password.txt

secret-key.txt:
	cat /dev/urandom | env LC_ALL=C tr -dc 'a-zA-Z0-9' | fold -w 99 | head -n 1 | tr -d '\n' > secret-key.txt


client-id.txt:
	cat /dev/urandom | env LC_ALL=C tr -dc 'a-zA-Z0-9' | fold -w 24 | head -n 1 | tr -d '\n' > client-id.txt


client-secret.txt:
	cat /dev/urandom | env LC_ALL=C tr -dc 'a-zA-Z0-9' | fold -w 24 | head -n 1 | tr -d '\n' > client-secret.txt

admin-password.txt:
	cat /dev/urandom | env LC_ALL=C tr -dc 'a-zA-Z0-9' | fold -w 12 | head -n 1 | tr -d '\n' > admin-password.txt