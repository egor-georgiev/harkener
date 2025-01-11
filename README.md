# Harkener: The Internet Sonification Project
**tl;dr: [try it out](https://nouuid4u.com/harkener/)**

## What It Does:
- Monitors TCP SYN packets that may indicate port scans on an internet-exposed server
- Gathers data on destination ports
- Transforms the collected data into sound

## TO-DO (what I plan to do myself):
- [ ] client:
    - [ ] rework sound synthesis and refactor
    - [ ] consider visualization
    - [ ] add a clear "about" section
- [ ] server:
    - [ ] simplify the websocket part (drop the hub design and register connections directly with capture)
    - [ ] add metrics:
        - [ ] collect data about packet distribution (and serve it to the client?)
        - [ ] count active clients/session time
        - [ ] count goroutines/channels/memory usage/learn about profiling
    - [ ] improve logging

