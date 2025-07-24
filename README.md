# Scout
Sentinal is an Open Source Peer-to-Peer Observability tool made in golang. 
It has two binaries, one that works as a hub (receives system metrics from other connections) and a publisher, 
which runs in the background of any machine to stream data to the hub.

scout works via websocket connections, which can be customized to a wanted refresh rate of the streamed data.

Here are some milestones to develop:
web dashboard,
more metrics, 
install as cli tool, (ex: scout pub --dest xyz --refresh 500),
docker support,
logs,

