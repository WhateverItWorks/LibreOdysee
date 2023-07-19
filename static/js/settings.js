const cookieSettings = [
  { name: "theme", type: "select", user: null, default: null },
  { name: "showRelated", type: "switch", user: null, default: null },
  { name: "nsfw", type: "switch", user: null, default: null },
  { name: "autoplay", type: "switch", user: null, default: null },
  { name: "commentWarning", type: "switch", user: null, default: null },
]
const sbCategories = ["sponsor", "selfpromo", "interaction", "intro", "outro", "preview", "filler"]

async function main() {
  // Apply the default settings
  let res = await fetch("/api/v1/settings")
  let defaults = await res.json()

  cookieSettings.forEach(setting => {
    setting.default = defaults[setting.name]
    applySetting(setting.name, setting.type, defaults[setting.name])
  })

  let sbDef = defaults.sponsorblock
  let defaultCategories = `${sbDef.sponsor ? 'sponsor,' : ''}${sbDef.selfpromo ? 'selfpromo,' : ''}${sbDef.interaction ? 'interaction,' : ''}${sbDef.intro ? 'intro,' : ''}${sbDef.outro ? 'outro,' : ''}${sbDef.preview ? 'preview,' : ''}${sbDef.filler ? 'filler' : ''}`;
  let re = /,$/g;
  defaultCategories = defaultCategories.replace(re, "")
  sbCategories.forEach(category => {
    let elem = document.getElementById(category)
    
    if (defaultCategories.includes(category)) {
      elem.checked = true
    }
  })

  // Apply user settings
  cookieSettings.forEach(setting => {
    let cookie = getCookie(setting.name)
    if (cookie) {
      setting.user = cookie
      applySetting(setting.name, setting.type, cookie)
    }
  })

  sbCategories.forEach(category => {
    let elem = document.getElementById(category)
    let userCategories = localStorage.getItem("sb_categories")
    
    if (userCategories && userCategories.includes(category)) {
      elem.checked = true
    }
  })

  // Listen for user changes
  for (let i = 0; i < cookieSettings.length; i++) {
    const setting = cookieSettings[i];
    let elem = document.getElementById(setting.name)
    
    elem.addEventListener("change", () => {
      let data = elem.checked
      if (setting.type == "select") {
        data = elem.value
      }
      setting.user = data
      setCookie(setting.name, data)
    })
  }
  for (let i = 0; i < sbCategories.length; i++) {
    const category = sbCategories[i];
    let elem = document.getElementById(category)
    
    elem.addEventListener("change", () => updateSBSetting(category))
  }

  // Allow user changes
  cookieSettings.forEach(setting => {
    let elem = document.getElementById(setting.name)
    elem.removeAttribute("disabled")
  })
  sbCategories.forEach(category => {
    let elem = document.getElementById(category)
    elem.removeAttribute("disabled")
  })
}
main()

function applySetting(name, type, data) {
  let elem = document.getElementById(name)
  switch (type) {
    case "switch":
      elem.checked = data === 'true' ? true : false
      break;
    case "select":
      elem.value = data
      break;
  }
}

function getCookie(name) {
  let cookies = document.cookie.split("; ");
  let cookie = cookies.filter(cookie => cookie.includes(name))[0]
  return cookie ? cookie.split("=")[1] : "";
}

function setCookie(name, data) {
  document.cookie = `${name}=${data}; path=/; SameSite=Strict; max-age=2147483647`
}

// Reload on theme change
document.getElementById("theme").addEventListener("change", () => {
  location.reload()
})

function updateSBSetting(category) {
  let categories = localStorage.getItem("sb_categories") || "";
  if (categories.includes(category)) {
    let re = new RegExp(`,?${category}`)
    localStorage.setItem("sb_categories", categories.replace(re, ""));
  } else if (categories.length == 0) {
    localStorage.setItem("sb_categories", categories + category);
  } else {
    localStorage.setItem("sb_categories", categories + "," + category);
  }

  let newCategories = localStorage.getItem("sb_categories")
  if (newCategories.startsWith(",")) {
    localStorage.setItem("sb_categories", newCategories.substring(1, 999));
  }
}
