slowgrog
========

Spent a good amount of time trying to find out why some code was causing performance issues with Redis. This was not Redis fault per say, however it would have been nice to have a tool that could hint to what usage was wrong. Go rocks and this seemed like the kind of tool that I could hack up real quick to surface some information about the goings on in Redis. 
Currently the core features are:
 
 * Configurable sampling of commands from the Redis MONITOR command
 * Gets the most recent 10 SLOWLOG 
 * Configurable sampling of the most recent X number of commands from monitor
 * INFO in JSON format

Command Arg:
------------
Tried to keep the arguments for the redis-cli the same for auth, port and host. The other arguments are more for configuring the sampling of the data. 

```
-h="127.0.0.1": redis host
-p=6379: redis port
-a="": Redis password

-cmdlimit=100: number of recent commands the MONITOR will store
-frequency=10000: Number of miliseconds to delay between samples INFO, SLOWLOG
-monsamplen=1000: Length of miliseconds that the monitor is sampled (0 will be coninuous however this is very costly to performance)
-slogsize=10: slowlog size
```


Vision:
-------
Given some more time I would like to add a number of features that I think would be extremely helpful findout out information about what is going on with Redis.

 * REST API that can just expose specific aspects of the JSON payload
 * Lua hooks to allow for extending the tool
 * UI would be nice to display this information in a more human friendly way.
