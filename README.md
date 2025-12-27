#BackGroundCC

It's is a Linux user-space background data transfer service designed to move bulk, delay-tolerant data while preserving low latency for foreground network traffic.
It is a user space re derivation of LEDBAT++.
The goal is not to maximize throughput, but to ensure that background transfers remain effectively invisible to interactive applications such as web browsing, SSH, video calls, gaming, or remote desktops. 
BackgroundCC behaves as a lower-than-best-effort transport, opportunistically using spare bandwidth when the network is idle and yielding smoothly and conservatively whenever other traffic competes for capacity.

The congestion control logic in BackgroundCC is delay-based and inspired by LEDBAT and LEDBAT++, as described in RFC 6817 (https://datatracker.ietf.org/doc/html/rfc6817) and the LEDBAT++ experimental research draft (https://datatracker.ietf.org/doc/draft-irtf-iccrg-ledbat-plus-plus/).
The sender continuously measures round-trip time (RTT) and maintains a dynamic baseline to estimate queueing delay. Rate adaptation is driven by a target-based control loop that increases the sending rate when queueing delay is low and backs off early as delay approaches or exceeds the target. Packet loss is treated as a secondary safety signal rather than the primary driver of congestion control.
BackgroundCC incorporates key ideas from LEDBAT++, including reduced gain, multiplicative decrease on sustained congestion, periodic slowdowns to refresh baseline delay estimates, and constrained startup behavior. Packet transmission is strictly time-paced to avoid micro-bursts that could temporarily inflate queues or distort delay measurements. 
The system is designed to run entirely in user space and uses UDP with application-level mechanisms for reliability, sequencing, and retransmission.

In addition to following the core principles of LEDBAT and LEDBAT++, BackgroundCC includes practical engineering choices to improve robustness on real networks. 
RTT samples are filtered to handle jitter from Wi-Fi, scheduling noise, and ACK compression, and congestion decisions are based on persistent delay trends rather than individual samples. 
When delay measurements become unreliable, the controller falls back to a conservative safe mode that prioritizes protecting foreground traffic. These behaviors are implementation-level extensions intended to make the system stable and predictable under heterogeneous and noisy network conditions.
