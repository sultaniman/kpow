(function () {
  const mQuery = matchMedia("(prefers-color-scheme: dark)");
  const themeKey = "theme";
  const themeAttr = "data-theme";
  const dark = "dark";
  const light = "light";

  // Track color theme change
  mQuery.onchange = ({ matches }) => {
    if (matches) {
      setColorTheme(dark)
    } else {
      setColorTheme(light)
    }
  };

  const getColorScheme = () => {
    // If media query doesn't match then light is default
    if (themeKey in localStorage) {
      return localStorage.getItem(themeKey);
    }

    if (mQuery.matches) {
      return dark;
    }

    return light;
  };

  const setColorTheme = (colorMode) => {
    if (typeof colorMode === "undefined") {
      colorMode = getColorScheme();
    }

    localStorage.setItem(themeKey, colorMode);
    document.documentElement.setAttribute(themeAttr, colorMode);
    document.querySelectorAll(".mode-icon").forEach((el) => {
      if (el.dataset.colorTheme !== colorMode) {
        el.removeAttribute("data-checked");
      } else {
        el.setAttribute("data-checked", "true");
      }
    })
  }

  // Detect color theme and set the relevant value
  window.onload = () => {
    setColorTheme();

    document.querySelectorAll(".mode-icon").forEach((el) => {
      el.addEventListener("click", (evt) => {
        evt.preventDefault();
        const targetEl = evt.target.closest("span");
        setColorTheme(targetEl.dataset.colorTheme);
      })
    })

    const showKeyBtn = document.getElementById("copyKey");
    if (showKeyBtn) {
      showKeyBtn.addEventListener("click", function (evt) {
        evt.preventDefault();
        console.log({ evt });
      })
    }
  };
}());
