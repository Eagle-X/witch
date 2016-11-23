# witch
RESTful process supervisor

# Install
```
go get github.com/eleme/witch
```

# Usage
```
Usage of witch:
  -c string
       	Config file (default "witch.yaml")
```

# Config file
```
# Listen address, default: :5671.
listen: :5671
# Specify the process control system, available controls buildin, supervisor and systemd.
# Default: buildin
control: buildin
# Only if control is supervisor or systemd, service MUST be given.
service:
# Only if control is buidin, command MUST be given.
command: sleep 3600
# The pid file of the process to be supervised, MUST change different one.
# Only if control is buildin, pid_file MUST be given.
pid_file: /var/run/witch/witch.pid
# Connection authentication username and password,
# the format is {username: password, ...}. default: {noadmin: noADMIN}.
auth: {noadmin: noADMIN}
```

# Exmaple
start witch
```
witch -c witch.ymal
```
control
```
curl -u noadmin:noADMIN -XPUT -d '{"name":"is_alive"}' http://127.0.0.1:5671/api/app/actions
curl -u noadmin:noADMIN -XPUT -d '{"name":"start"}' http://127.0.0.1:5671/api/app/actions
curl -u noadmin:noADMIN -XPUT -d '{"name":"stop"}' http://127.0.0.1:5671/api/app/actions
curl -u noadmin:noADMIN -XPUT -d '{"name":"restart"}' http://127.0.0.1:5671/api/app/actions
```


