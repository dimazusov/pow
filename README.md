# TCP-server with protection from DDOS attacks, based on Proof of Work

## 1. Description
Test task for Server Engineer
Design and implement “Word of Wisdom” tcp server.
• TCP server should be protected from DDOS attacks with the Prof of Work
(https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should
be used.
• The choice of the POW algorithm should be explained.
• After Prof Of Work verification, server should send one of the quotes from “word of
wisdom” book or any other collection of the quotes.
• Docker file should be provided both for the server and for the client that solves the
POW challenge.

## 2. Getting started
### 2.1 Requirements
+ [Docker](https://docs.docker.com/engine/install/) installed (to run docker-compose)

### 2.2 Start pow server and client by docker-compose:
```
make start
```

### 2.3 Start server:
```
make start-server
```

### 2.3 Stop server:
```
make stop-server
```

### 2.4 Start client:
```
make start-client
```

### 2.5 Generate protobuf schemes:
```
make gen-proto
```

### 2.6 Run tests:
```
make test
```

### 2.7 Run linter:
```
make lint
```

### 2.8 Example
[!example](https://ibb.co/P4VP9T2)

## 3. Protocol definition
This solution uses TCP-based protocol.
The server exchanges messages with clients.
Message consists of: 
- message length, first 4 bytes, this is integer number
- encoded protobuff structure, you can see this messages /api/api.proto

### 3.1 Types of requests
Solution supports 4 types of requests, switching by header:
+ 1 - InitRequest - from client to server - request new challenge from server
+ 2 - InitResponse - from server to client - message with challenge for client
+ 3 - ChallengeRequest - from client to server - message with solved challenge, for verification
+ 4 - ChallengeResource - from server to client - message with random quote

### 3.2 Examples of protocol message
Here i provide examples for all types of requests:
```
message InitRequest {
  int32 ProtocolVersion = 1;
}

message InitResponse {  
  bytes Challenge = 1;    
  int32 Difficulty = 2;
}

message ChallengeRequest { 
  uint64 Nonce = 1;
}

message ChallengeResponse {
  string Quote = 1;
} 
```

## 4. Proof of Work
Idea of Proof of Work for DDOS protection is that client, which wants to get some resource from server, 
should firstly solve some challenge from server. 
This challenge should require more computational work on client side and verification of challenge's solution - much less on the server side.

### 4.1 Selection of an algorithm
There is some different algorithms of Proof Work. 
I compared next three algorithms as more understandable and having most extensive documentation:

- [Hashcash](https://en.wikipedia.org/wiki/Hashcash)
- [Guided tour puzzle protocol](https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol)
- [Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree)

After evaluating the options, I opted for Hashcash. The other algorithms present the following drawbacks:
- With Merkle tree, the server has to perform excessive work to verify the client's solution. For a tree with 4 leaves and a depth of 3, the server would need to carry out 3 hash calculations.
- The guided tour puzzle requires the client to frequently ask the server for the next part of the guide, complicating the protocol's logic.

Hashcash, instead has next advantages:
- Easy to implement
- Abundant documentation and descriptive articles
- Straightforward server-side validation
- Adjustable client difficulty through modification of the required number of leading zeros

Of course Hashcash also has disadvantages like:
- The computation time is influenced by the client's machine power. For instance, very weak clients might not be able to solve the challenge, while extremely powerful computers could carry out DDOS attacks. However, the challenge's difficulty can be dynamically adjusted by the server by altering the required number of leading zeros.
- Pre-computing challenges in advance of a DDOS attack is another concern. Some clients might analyze the protocol and precompute numerous challenges to use them all at once. This issue can be mitigated by additional validation of Hashcash parameters on the server. For example, when creating a challenge, the server could store the rand value in a Redis cache and check its existence during the verification step (as implemented in this solution).

But all of those disadvantages could be solved in real production environment. 

## 5. Structure of the project
Project structure implements [Go-layout](https://github.com/golang-standards/project-layout) pattern.
Existing directories:
+ api/api.proto protobuf entities description
+ cmd/client/main.go main file for client
+ cmd/server/main.go main file for server
+ internal/config/server.go - config server
+ internal/config/client.go - config client
+ internal/client - all logic of client
+ internal/server - all logic of server
+ internal/connection - wrapper for connection
+ internal/hashcash - logic of chosen PoW algorithm (Hashcash)
+ internal/storage - storage for values, for saving a challenge, and quotes

## 6. How to improve
1. Adaptive Difficulty Adjustment:
- Implement a dynamic difficulty adjustment mechanism that responds more rapidly to changes in network conditions. This helps balance the load and prevents both overburdening weaker clients and under-challenging stronger ones.
2. Enhanced Security Measures:
- Introduce additional validation steps, such as tracking and limiting the number of challenge requests per client to prevent abuse.
- Use rate limiting to reduce the risk of pre-computed challenges being used for DDOS attacks.
3. Cache Optimization:
- Optimize the storage and retrieval of nonce values (e.g., rand values) in the server cache. Use efficient data structures and algorithms to minimize latency and ensure quick verification.
4. Client-Side Rate Limiting:
- Implement client-side rate limiting to reduce the number of challenge requests a single client can make within a given time frame. This can help prevent abuse and reduce the risk of DDOS attacks.
5. Better Resource Utilization:
- Develop algorithms to monitor and balance server load, distributing validation tasks more efficiently across multiple servers if necessary.
6. Enhanced Logging and Monitoring:
- Implement comprehensive logging and monitoring systems to detect unusual patterns or potential attacks early, allowing for quick responses to potential threats.
7. Client Authentication:
- Require clients to authenticate themselves before issuing challenges. This can involve using additional security layers like tokens or certificates.