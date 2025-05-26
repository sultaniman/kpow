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
    } else {
      if (mQuery.matches) {
        return dark;
      }
    }

    return light;
  };

  const setColorTheme = (colorMode) => {
    localStorage.setItem(themeKey, colorMode);
    document.documentElement.setAttribute(themeAttr, colorMode);
    document.querySelectorAll(".mode-icon").forEach((el) => {
      const colorTheme = el.dataset.colorTheme;
      if (colorTheme !== colorMode) {
        el.removeAttribute("data-checked");
      } else {
        el.setAttribute("data-checked", "true");
      }
    })
  }

  // Detect color theme and set the relevant value
  window.onload = () => {
    setColorTheme(getColorScheme());

    document.querySelectorAll(".mode-icon").forEach((el) => {
      el.addEventListener("click", (evt) => {
        evt.preventDefault();
        const targetEl = evt.target.closest("span");
        setColorTheme(targetEl.dataset.colorTheme);
      })
    })
  };
}());
