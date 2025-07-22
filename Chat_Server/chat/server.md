```
ZOO → ZIRC: PASS hunter2 TS 6 :ZOO                 # ZOO authenticates to ZIRC
ZOO → ZIRC: CAPAB :QS EX ...                     # ZOO advertises capabilities
ZOO → ZIRC: SERVER ZOO 1 :Zach's Server            # ZOO introduces itself

🟦 ZIRC adds ZOO to its internal server graph

ZIRC → ZOO: CAPAB :QS EX ...                     # ZIRC replies with its capabilities

ZIRC: Sets its conn state to BURST_SEND
ZIRC → ZOO: SERVER ZIRC 1 :Server ZIRC Desc            # ZIRC introduces itself (as part of burst)
🟦 ZOO adds ZIRC to its internal server graph:
🟦 ZOO updates its Conn state for ZIRC saying that it is BURST_RECV

ZIRC → ZOO: SERVER C 2 :Server C Desc            # ZIRC introduces its downstream server
🟦 ZOO adds C to its server graph:

ZIRC → ZOO: UID userC 2 ... :Charlie             # ZIRC introduces a user on C
🟦 ZOO adds Charlie to user table, associated with C

ZIRC → ZOO: PING ZOO
Singles to ZOO that bursting is complete

🟩 ZOO now knows the connection to ZIRC is accepted

🔁 ZOO → X: SERVER ZIRC 1 :Server ZIRC Desc          # ZOO propagates ZIRC to its downstream server X
🔁 ZOO → X: SERVER C 2 :Server C Desc          # ZOO propagates C to X
🔁 ZOO → X: UID userC 2 ... :Charlie           # ZOO propagates Charlie to X

📤 ZOO → ZIRC: SERVER X 2 :Server X Desc          # ZOO bursts its own topology to ZIRC
🟦 ZIRC adds X to its server graph:
     graph.add(serverName="X", hopcount=2, viaConnection=ZOO)

📤 ZOO → ZIRC: UID userX 2 ... :Xander            # ZOO introduces a user from X
🟦 ZIRC adds Xander to user table, associated with X
```