# catl

Kibana is an excellent tool for visualising and performing analysis on your log data but I've found that the Discover pane isn't that effective when you are trying to pull contextual information about a particular event in a log (what happened before and after a particular error?) or in a number of other situations where commandline utilities like `grep` have traditionally been useful.

`catl` is my attempt to address this problem. It's a simple program that executes a search query against an elasticsearch cluster and returns the message field in timestamp order on stdout. This allows you to use utilities like `grep`, `awk` and `sed` to parse your logs for information as if it was a file on your local disk.


## Build
```
$ go build
```

## Usage
```
$ catl -h
usage: catl [<flags>] <query>

Writes a logstash elasticsearch query to stdout as if it were a logfile

Flags:
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
  -i, --index="logstash-*"       Index pattern.
  -u, --url="http://localhost:9200"
                                 Logstash server URL.
  -m, --message-field="message"  Field to be returned
  -s, --sort-field="@timestamp"  Field to sort the results by

Args:
  <query>  Elasticseach query string.
```
You can set any of the flags with environment variables in the form `CATL_$FLAG`:
```shell
export CATL_INDEX=logstash-2018.07.02
export CATL_MESSAGE_FIELD=logmessage
```

## Examples
Returns ERROR level logs along with the preceding message:
```
$ catl 'meta.service: "foobar" AND meta.environment: "live" AND type: "java_logback"' | grep 'ERROR' -B 1
2018-07-03 00:00:00,522 [http-nio-8080-exec-7] DEBUG c.b.b.f.c.SillyClass - Doing something...
2018-07-03 00:00:00,522 [http-nio-8080-exec-4] ERROR  c.b.b.f.SillyErrorClass - Something foolish happened
```
Save the output to a file:
```
$ catl 'meta.service: "foobar" AND meta.environment: "live" AND type: "java_logback" AND loglevel: "ERROR"' > errors.log
```