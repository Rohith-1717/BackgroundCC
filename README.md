# BackGroundCC

BackgroundCC is a Linux user space background data transfer service designed to move bulk, delay tolerant data without interfering with interactive network traffic.
It is a user space re derivation of LEDBAT++. The primary objective is simple: background transfers should remain invisible to the user experience. Applications such as web browsing, SSH, gaming, VoIP, video conferencing, and remote desktops must not experience noticeable latency increases because of background data movement.
BackgroundCC uses spare bandwidth when the network is idle and yields smoothly and conservatively as soon as competing traffic appears.

## Design 
The system is built around a latency first philosophy.
Instead of aggressively chasing throughput, BackgroundCC protects foreground traffic. It assumes that bulk transfers are less time sensitive than interactive flows and therefore must adapt quickly and gracefully under congestion.

The design is guided by the following principles:

- Latency protection: Interactive applications must remain responsive even during large background transfers.
- Delay based congestion control: Queueing delay is treated as the primary congestion signal rather than packet loss.
- User space simplicity: The entire system runs in user space using UDP. No kernel modifications are required.
- Stability under real networks: The controller is engineered to handle jitter, wireless noise, scheduling delays, and heterogeneous paths.
- Research grounded but practical: The implementation follows LEDBAT++ principles while incorporating practical extensions necessary for real deployments.

## Delay Based Congestion Model
The congestion controller is built around a target queueing delay model.

Let:

RTT(t) = measured round trip time at time t
BaseDelay = minimum RTT observed over a sliding window
QueueDelay(t) = RTT(t) − BaseDelay

BaseDelay approximates propagation delay without queueing. QueueDelay represents congestion induced buffering.
A configurable target delay T is defined. Typical values are in the range of 25 to 100 ms depending on configuration.

### The control objective is:
QueueDelay(t) is approx. T
The error term is:

e(t) = T − QueueDelay(t)

- If e(t) > 0
The queue is under target -> increase rate.
- If e(t) < 0
The queue exceeds target -> decrease rate.

## Window and rate update rule

BackgroundCC maintains a congestion window cwnd measured in packets.

### A simplified additive update rule is:
  
cwnd = cwnd + G*(e(t)/T)
where G is a gain parameter controlling aggressiveness.

### This can also be expressed as:
  
delta_cwnd proportional to (T − QueueDelay)/T

When delay is far below target, the increase is larger. As delay approaches T, the increase becomes smaller.

- Under sustained congestion:
  cwnd = cwnd * beta
where 0 < beta < 1.
This multiplicative decrease ensures fast yielding when congestion persists.


## How to Run
 
BackgroundCC is a two-sided tool: one machine must run in `receive` mode before the other starts `send`ing, since the receiver is what performs reassembly and generates the ACKs the sender's congestion control depends on.
 
### Build
 
```bash
git clone https://github.com/Rohith-1717/BackgroundCC.git
cd BackgroundCC/backgroundcc
go build -o bgcc .
```
 
### Run a transfer
 
**On the receiving machine**, start listening first:
 
```bash
./bgcc receive <local_ip>:<port> <output_file>
```
 
Example:
```bash
./bgcc receive 0.0.0.0:9000 received_file.dat
```
- `0.0.0.0:9000` — local address/port to bind and listen on (any interface, port 9000)
- `received_file.dat` — path where the reassembled file will be written
**On the sending machine**, once the receiver is up:
 
```bash
./bgcc send <local_ip>:<port> <remote_ip>:<port> <file_to_send>
```
 
Example:
```bash
./bgcc send 0.0.0.0:9001 192.168.1.50:9000 myfile.dat
```
- `0.0.0.0:9001` — sender's own local bind address/port
- `192.168.1.50:9000` — the receiver's address (must match what it bound to)
- `myfile.dat` — the file to transfer

## Reference
LEDBAT: https://datatracker.ietf.org/doc/html/rfc6817

LEDBAT++: https://datatracker.ietf.org/doc/draft-irtf-iccrg-ledbat-plus-plus/
