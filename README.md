# Harkener
Harkener transforms the 'noise' of the internet — unseen connection attempts from port scanners, web crawlers, and other probes — into sound.
You can try it out via the **[link](https://nouuid4u.com/harkener/)**.

## What it Does Technically:
- listens for TCP SYN packets
- gathers data on the source IP address and the destination port number
- transforms the collected data into sound, using the destination port number to modulate the frequency

## TODO:
- [x] drop TLS support
- [ ] switch to SSE
  - handle re-connects on the client side
- [ ] add metrics
- [ ] serve the source IP as well as destination port
- [ ] sonification improvements:
  - [ ] dev-version with bells and whistles for tuning:
    - min, max freq
    - envelope
    - gain
    - volume
- [ ] hilbert curve visualization
- [ ] drop cobra for pflag
