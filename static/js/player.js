async function main() {
  let res = await fetch("/api/v1/settings")
  let defaults = await res.json()

  const player = new Plyr('#player', {
    keyboard: { focused: true, global: true },
    speed: {
      selected: defaults.speed
    }
  });

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
}
main()