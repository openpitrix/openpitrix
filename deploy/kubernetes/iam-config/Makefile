
default: secret-key.txt

secret-key.txt:
	cat /dev/urandom | env LC_CTYPE=C tr -dc 'a-zA-Z0-9' | fold -w 99 | head -n 1 > secret-key.txt