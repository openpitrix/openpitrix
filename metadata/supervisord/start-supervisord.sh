#!/bin/sh

SERVICE=$1

case $SERVICE in
	"frontgate")
		echo "files = /etc/supervisor.d/frontgate.ini /etc/supervisor.d/metad.ini" >> /etc/supervisord.conf
		exec supervisord
		;;
	"drone")
		echo "files = /etc/supervisor.d/drone.ini" >> /etc/supervisord.conf
		exec supervisord
		;;
	*)
		echo "Usage: start-supervisord.sh <frontgate|drone>"
		;;
esac
