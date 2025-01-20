# blue orange green

Radio or Sattelite (Iridium)

Data:
- gps position
- accelerometer
- depth
- sea surface temp
- "COG" - course over ground
- "SOG" - speed over ground
- range
- last uplink
- next expected uplink

Display underlying ocean currents. 

Unsupported: Sonic Modems, Ropeless Buoy/Traps

"Something funny is happening"
- pulled underwater
- started wandering


Plotter link
- avoid already set
- trail on map "just for that set", must be manually reset?
- Looks like it's just writing  to a folder?

Avoid already set gear.


## Weather 

Fun idea:
- Boats are slow
- Nowcast for area area where you currently are
- 1hr forecast 10 miles out, 2hr 20, etc. 


## Ideas on what to build

Backend is some way to reliably assemble a shared log.
Assumptions are "at least once" delivery, where buoys could presumably store tons of old messages and have them uploaded once they're charging or something?
Frontend basically parses and displays this log.

### Buoysim 

- Each buoy is a little goroutine
- Each routine is just a simple state machine
- After a timeout they execute some transition function, then upload some protocol buffer

### Backend

- At-least-once + idempotence 
- Assemble a log, maybe just in postgres? 
- query it and serve up cool tracks and such.

### Frontend

- Display a map
- Put markers on the map?
- Svelte sounds fun? 


