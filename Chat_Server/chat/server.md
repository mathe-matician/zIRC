```
A → B: PASS hunter2 TS 6 :A                 # A authenticates to B
A → B: CAPAB :QS EX ...                     # A advertises capabilities
A → B: SERVER A 1 :Zach's Server            # A introduces itself

🟦 B adds A to its internal server graph:
     graph.add(serverName="A", hopcount=1, description="Zach's Server", viaConnection=A)

B → A: CAPAB :QS EX ...                     # B replies with its capabilities

B → A: SERVER B 1 :Server B Desc            # B introduces itself (as part of burst)
🟦 A adds B to its internal server graph:
     graph.add(serverName="B", hopcount=1, description="Server B Desc", viaConnection=B)

B → A: SERVER C 2 :Server C Desc            # B introduces its downstream server
🟦 A adds C to its server graph:
     graph.add(serverName="C", hopcount=2, viaConnection=B)

B → A: UID userC 2 ... :Charlie             # B introduces a user on C
🟦 A adds Charlie to user table, associated with C

🟩 A now knows the connection to B is accepted

🔁 A → X: SERVER B 1 :Server B Desc          # A propagates B to its downstream server X
🔁 A → X: SERVER C 2 :Server C Desc          # A propagates C to X
🔁 A → X: UID userC 2 ... :Charlie           # A propagates Charlie to X

📤 A → B: SERVER X 2 :Server X Desc          # A bursts its own topology to B
🟦 B adds X to its server graph:
     graph.add(serverName="X", hopcount=2, viaConnection=A)

📤 A → B: UID userX 2 ... :Xander            # A introduces a user from X
🟦 B adds Xander to user table, associated with X
```