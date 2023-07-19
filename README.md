### ⚠️ **Notice:** Librarian has been archived. https://bcow.xyz/posts/archiving-librarian/

#### LibreOdysee
An alternative frontend for LBRY/Odysee. Inspired by [Invidious](https://github.com/iv-org/invidious) and [Libreddit](https://github.com/spikecodes/libreddit).

## Table of Contents
- [Features](#features)
- [Comparison](#comparison)
  - [Speed](#speed)
  - [Privacy](#privacy)
    - [Odysee](#odysee)
    - [LibreOdysee](#libreodysee)
- [Instances](#instances)
  - [Clearnet](#clearnet)
  - [Tor](#tor)
- [Automatically redirect links](#automatically-redirect-links)
  - [LibRedirect](#libredirect)
  - [GreaseMonkey script](#greasemonkey-script)
- [Install](#install)

## Features

- JavaScript not required*
- Lightweight
- No ads
- No tracking
- No crypto garbage

⚠️ * JavaScript is required to play livestreams except on Apple devices.

⚠️ **Notice: All content is from Odysee. LibreOdysee does not allow uploads or comments. Any issues with content should be reported to Odysee.**

## Demo

[🎞️ Video](https://watch.whateveritworks.org/@RetroMusic:d/1987-Rick-Astley-Never-Gonna-Give-You-Up-1920x1080:f)
[📺 Channel](https://watch.whateveritworks.org/@DistroTube:2)
[📰 Article](https://watch.whateveritworks.org/@Odysee:8/spooky:b)

## Comparison
Comparing LibreOdysee to Odysee.

### Speed
Tested using [Google PageSpeed Insights](https://pagespeed.web.dev/).

|             | [LibreOdysee](https://pagespeed.web.dev/analysis/https-watch-whateveritworks-org/98d1lc08yr) | [Odysee](https://pagespeed.web.dev/report?url=https%3A%2F%2Fodysee.com%2F) |
| ----------- | --------- | ------ |
| Performance | 99 | 27 |
| Request count | 17 | 470 |
| Resource Size | 702 KiB | 2,457 KiB |
| Time to Interactive | 0.9s | 18.4s |

### Privacy

#### Odysee
<a href="https://tosdr.org/en/service/2391">
  <img alt="Odysee Privacy Grade" src="https://shields.tosdr.org/en_2391.svg">
</a>

Odysee has admitted to using browser fingerprinting for ads and loads multiple ads, trackers, and an annoying cookie banner.

> We and our partners process data to provide:
Use precise geolocation data. Actively scan device characteristics for identification. Store and/or access information on a device. Personalised ads and content, ad and content measurement, audience insights and product development.

They also use your data for these purposes and you cannot opt-out of it.
- Ensure security, prevent fraud, and debug
- Technically deliver ads or content
- Match and combine offline data sources
- Link different devices
- Receive and use automatically-sent device characteristics for identification

See what trackers and cookies they use: https://themarkup.org/blacklight.?url=odysee.com

#### LibreOdysee
LibreOdysee itself does not collect any data but instance operators may collect data. You can view a "privacy nutrition label" by clicking on the "Privacy" link at the bottom.

## Instances

Open an issue to have your instance listed here!

⭐ Instances with a star don't collect data, don't use Cloudflare, support livestreams, and proxy videos.

### Clearnet

| URL                                                             | Country      | Provider         | Privacy               | Livestreams | Proxy | Notes |
| :-------------------------------------------------------------- | :----------- | :--------------- | :-------------------- | :---------- | :---- | :---- |
| [odysee.076.ne.jp](https://odysee.076.ne.jp)                    | 🇯🇵 JP        | GMOグローバルサイン | ✅ Data not collected | ✅️ | ❌️ | [Edited theme](https://git.076.ne.jp/TechnicalSuwako/Librarian-mod) |
| [librarian.pussthecat.org](https://librarian.pussthecat.org/)   | 🇩🇪 DE        | Hetzner          | ⚠️ Data collected     | ✅️ | ✅️ |  |
| [lbry.projectsegfau.lt](https://lbry.projectsegfau.lt/)         | 🇱🇺 LU, 🇺🇸 US, 🇮🇳 IN | See below | ✅ Data not collected | ✅️ | -  |  |
| [lbry.us.projectsegfau.lt](https://lbry.us.projectsegfau.lt/)   | 🇺🇸 US | DigitalOcean            | ✅ Data not collected | ✅️ | ❌️ |  |
| ⭐ [lbry.eu.projectsegfau.lt](https://lbry.eu.projectsegfau.lt/)| 🇱🇺 LU | FranTech Solutions      | ✅ Data not collected | ✅️ | ✅️ |  |
| [lbry.in.projectsegfau.lt](https://lbry.in.projectsegfau.lt/)   | 🇮🇳 IN | Airtel                  | ✅ Data not collected | ✅️ | ❌️ |  |
| ⭐ [librarian.esmailelbob.xyz](https://librarian.esmailelbob.xyz/) | 🇨🇦 CA     | OVH              | ✅ Data not collected | ✅️ | ✅️ |
| [lbry.vern.cc](https://lbry.vern.cc/)                           | 🇺🇸 US        | OVHCloud         | ✅ Data not collected | ❌️ | ❌️ |[Edited theme](https://git.vern.cc/root/modifications/src/branch/master/librarian) |
| [lbry.slipfox.xyz](https://lbry.slipfox.xyz)                    | 🇺🇸 US        | Hetzner          | ✅ Data not collected | ❌️ | ❌️ |  |
| [lbry.mywire.org](https://lbry.mywire.org)                      | 🇷🇺 RU        | justhost.ru      | ✅ Data not collected | ❌️ | ❌️ |  |
| [lbry.ooguy.com](https://lbry.ooguy.com)                        | 🇸🇰 SK        | STARK INDUSTRIES | ✅ Data not collected | ❌️ | ❌️ |  |
| [lbn.frail.duckdns.org](https://lbn.frail.duckdns.org/)         | 🇧🇷 BR        | WSNET TELECOM    | ✅ Data not collected | ✅️ | ❌️ |  |
| [lbry.ramondia.net](https://lbry.ramondia.net/)		          | 🇩🇪 DE	       | Hetzner	      | ✅ Data not collected | ✅ | ✅ |  |

### Tor

| URL | Privacy               | Live streams | Notes |
| :-- | :-------------------- | :----------- | :---- |
| ⭐ [librarian.esmail5pdn24shtvieloeedh7ehz3nrwcdivnfhfcedl7gf4kwddhkqd.onion](http://librarian.esmail5pdn24shtvieloeedh7ehz3nrwcdivnfhfcedl7gf4kwddhkqd.onion/) | ✅ Data not collected | ✅️ | Onion of librarian.esmailelbob.xyz |
| [lbry.vernccvbvyi5qhfzyqengccj7lkove6bjot2xhh5kajhwvidqafczrad.onion](http://lbry.vernccvbvyi5qhfzyqengccj7lkove6bjot2xhh5kajhwvidqafczrad.onion/) | ✅ Data not collected | ❌️ | Onion of lbry.vern.cc. [Edited theme](https://git.vern.cc/root/modifications/src/branch/master/librarian) |
| [5znbzx2xcymhddzekfjib3isgqq4ilcyxa2bsq6vqmnvbtgu4f776lqd.onion](http://5znbzx2xcymhddzekfjib3isgqq4ilcyxa2bsq6vqmnvbtgu4f776lqd.onion/) | ✅ Data not collected | ❌️ | Onion of lbry.slipfox.xyz |
| [bxewpsswttslepw27w2hhxhlizwm7l7y54x3jw5cfrb64hb6lgc557ad.onion](http://bxewpsswttslepw27w2hhxhlizwm7l7y54x3jw5cfrb64hb6lgc557ad.onion/) | ✅ Data not collected | ❌️ | Onion of lbry.ooguy.com |
| ⭐ [lbry.pjsfkvpxlinjamtawaksbnnaqs2fc2mtvmozrzckxh7f3kis6yea25ad.onion](http://lbry.pjsfkvpxlinjamtawaksbnnaqs2fc2mtvmozrzckxh7f3kis6yea25ad.onion/) | ✅ Data not collected | ✅️ | Onion of lbry.eu.projectsegfau.lt |

### I2P

| URL | Privacy               | Live streams | Notes |
| :-- | :-------------------- | :----------- | :---- |
| ⭐ [pjsf7uucpqf2crcmfo3nvwdmjhirxxjfyuvibdfp5x3af2ghqnaa.b32.i2p](http://pjsf7uucpqf2crcmfo3nvwdmjhirxxjfyuvibdfp5x3af2ghqnaa.b32.i2p/) | ✅ Data not collected | ✅️ | lbry.eu.projectsegfau.lt on I2P |

### Automatically redirect links

#### LibRedirect
Use [LibRedirect](https://github.com/libredirect/libredirect) to automatically redirect Odysee links to Librarian! This needs to be enabled in settings.
- [Firefox](https://addons.mozilla.org/firefox/addon/libredirect/)
- [Chromium-based browsers (Brave, Google Chrome)](https://github.com/libredirect/libredirect#install-in-chromium-brave-and-chrome)
- [Edge](https://microsoftedge.microsoft.com/addons/detail/libredirect/aodffkeankebfonljgbcfbbaljopcpdb)

#### GreaseMonkey script
There is a script to redirect Odysee links to Librarian.
https://github.com/WhateverItWorks/LibreOdysee-Redirects

### Install
```
git clone https://github.com/WhateverItWorks/LibreOdysee.git
cd LibreOdysee
mkdir data
cp config.example.yml data/config.yml
nano data/config.yml
docker-compose up -d --build
```
LibreOdysee should now be running at http://localhost:8080.
