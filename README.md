<div align="center">
  <h1>Terminal Typeracer</h1>
  <i> Typeracer in the terminal over ssh!</i>
</div>

---

Made in <24 hours for a hackathon. I won the "Spirit of the Hackathon" award :tada:


Try it out at
```bash
ssh 45.55.159.44 -p 23234
```
Use ctrl-c to quit


[Demo](https://github.com/user-attachments/assets/ee0b2d70-2e5a-4f47-bfba-9d2c25dce2ee)



### Inspiration
[typeracer](https://play.typeracer.com/) and [terminal.shop](https://www.terminal.shop/)

### What it does
You can play typeracer in the terminal over ssh


### How I built it
[Bubbletea](https://github.com/charmbracelet/bubbletea) and lots of love

### Challenges I ran into
I Severely underestimated how much time I had to do this. Luckily I managed to get an MVP without too much trouble.

### Accomplishments that I'm proud of
- Wrote my own networking protocol using just raw TCP.
- Working ssh server that I can go onto whenever I want

### What I learned
i got some more experience with go and I learned how to host things outside of my local network

### What's next for Terminal Typeracer
Clean up the UI, make it easier to join and leave lobbies


### Running locally
Clone the repo and cd into the project directory. Then, open two terminals and run
```
# Terminal 1
go run ./cmd/server

# Terminal 2
go run ./cmd/ssh
```

Alternatively, you can do the following in one terminal session

```
go run ./cmd/server & go run ./cmd/ssh
```

After starting the server and ssh handler, you can do `ssh localhost -p 23234` to connect.
