# FossaBot Points [<img align="right" src="https://cdn.tcole.me/logo.png" height="35">](hhttps://github.com/TimothyCole/FossaPoints)

[![Discord](https://img.shields.io/discord/313591755180081153.svg?label=Personal%20Discord&colorB=308bcd&maxAge=3600)](https://discordapp.com/invite/YFtfGwq)
[![Follow on Twitter](https://img.shields.io/twitter/follow/modesttim.svg?style=popout&label=Follow%20on%20Twitter)](https://twitter.com/intent/follow?screen_name=modesttim)

---

### Information
This is a collection of HTTP Handlers written in Go that I host on [Google Cloud Functions Go Alpha](https://medium.com/google-cloud/google-cloud-functions-for-go-57e4af9b10da) that I use to add point functionality to [FossaBot](https://fossabot.com/).

In order to automatically give points for an interval I call the GiveAll functions from Google's [Cloud Scheduler](https://cloud.google.com/scheduler/).
