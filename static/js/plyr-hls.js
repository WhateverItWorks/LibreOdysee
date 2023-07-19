document.addEventListener('DOMContentLoaded', async () => {
  let res = await fetch("/api/v1/settings")
  let defaults = await res.json()

  const defaultOptions = {
    keyboard: { focused: true, global: true },
    speed: {
      selected: defaults.speed
    }
  };
  
  const video = document.querySelector('video');
  const source = video.getElementsByTagName("source")[0].src;

  // Preview thumbnails
  if (!document.getElementById("isLive")) {
    const vttUrl = source.replace("master.m3u8", "stream_sprite.vtt")
    defaultOptions.previewThumbnails = {
      enabled: true,
      src: vttUrl
    }
  }

  if (!Hls.isSupported()) {
    video.src = source;
    const player = new Plyr(video, defaultOptions);

    // Keyboard shortcuts
    document.addEventListener('keydown', (event) => {
      if (event.target.id === "searchBar") return;
      event.preventDefault()
      switch (event.key) {
        case 'j':
          player.rewind(15);
          break;
        case ' ':
          player.togglePlay();
          break;
        case 'l':
          player.forward(15);
          break;
      }
    });
  } else {
    // For more Hls.js options, see https://github.com/dailymotion/hls.js
    const hls = new Hls();
    hls.loadSource(source);

    // From the m3u8 playlist, hls parses the manifest and returns
    // all available video qualities. This is important, in this approach,
    // we will have one source on the Plyr player.
    hls.on(Hls.Events.MANIFEST_PARSED, function (event, data) {

      // Transform available levels into an array of integers (height values).
      const availableQualities = hls.levels.map((l) => l.height)
      availableQualities.unshift(0) //prepend 0 to quality array

      // Add new qualities to option
      defaultOptions.quality = {
        default: defaults.quality, //Default - AUTO
        options: availableQualities,
        forced: true,
        onChange: (e) => updateQuality(e),
      }
      // Add Auto Label 
      defaultOptions.i18n = {
        qualityLabel: {
          0: 'Auto',
        },
      }

      hls.on(Hls.Events.LEVEL_SWITCHED, function (event, data) {
        var span = document.querySelector(".plyr__menu__container [data-plyr='quality'][value='0'] span")
        if (hls.autoLevelEnabled) {
          span.innerHTML = `Auto (${hls.levels[data.level].height}p)`
        } else {
          span.innerHTML = `Auto`
        }
      })

      // Initialize new Plyr player with quality options
      const player = new Plyr(video, defaultOptions);

      // Keyboard shortcuts
      document.addEventListener('keydown', (event) => {
        if (event.target.id === "searchBar") return;
        event.preventDefault()
        switch (event.key) {
          case 'j':
            player.rewind(15);
            break;
          case ' ':
            player.togglePlay();
            break;
          case 'l':
            player.forward(15);
            break;
        }
      });

      if (location.hash) {
        player.on('loadeddata', () => { player.currentTime = location.hash.replace("#", "") * 1 })
      }

      const urlParams = new URLSearchParams(location.search);
      if (urlParams.get("t")) {
        player.on('loadeddata', () => { player.currentTime = urlParams.get("t") * 1 })
      }

      window.addEventListener('hashchange', () => {
        player.currentTime = location.hash.replace("#", "") * 1
      })
    });

    hls.attachMedia(video);
    window.hls = hls;
  }

  function updateQuality(newQuality) {
    if (newQuality === 0) {
      window.hls.currentLevel = -1; //Enable AUTO quality if option.value = 0
    } else {
      window.hls.levels.forEach((level, levelIndex) => {
        if (level.height === newQuality) {
          console.log("Found quality match with " + newQuality);
          window.hls.currentLevel = levelIndex;
        }
      });
    }
  }
});