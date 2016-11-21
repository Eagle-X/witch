# witch
RESTful process supervisor

# Install
```
go get github.com/eleme/witch
```

# Usage
```
Usage: witch [config file] [cmd path] [cmd arguments]
```

# Config file
```
# Listen address, default: :5671.
listen: :5671
# The pid file of the process to be supervised, MUST change different one.
pid_file: /var/run/witch/witch.pid
# Connection authentication username and password,
# the format is {username: password, ...}. default: {noadmin: noADMIN}.
auth: {noadmin: noADMIN}
```

# Exmaple
start witch
```
witch witch.ymal sleep 10000
```
control
```
curl -u noadmin:noADMIN -XPUT -d '{"name":"is_alive"}' http://127.0.0.1:5671/api/app/actions
curl -u noadmin:noADMIN -XPUT -d '{"name":"start"}' http://127.0.0.1:5671/api/app/actions
curl -u noadmin:noADMIN -XPUT -d '{"name":"stop"}' http://127.0.0.1:5671/api/app/actions
curl -u noadmin:noADMIN -XPUT -d '{"name":"restart"}' http://127.0.0.1:5671/api/app/actions
```


