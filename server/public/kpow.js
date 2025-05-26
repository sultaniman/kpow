(function () {
  const mQuery = matchMedia("(prefers-color-scheme: dark)");

  // Track color theme change
  mQuery.onchange = ({ matches }) => {
    if (matches) {
      localStorage.setItem("theme", "dark");
      document.documentElement.setAttribute("data-theme", "dark");
    } else {
      localStorage.setItem("theme", "light");
      document.documentElement.setAttribute("data-theme", "light");
    }
  };

  const getColorScheme = () => {
    // If media query doesn't match then light is default
    let colorScheme = "light";
    if ("theme" in localStorage) {
      return localStorage.getItem("theme");
    } else {
      if (mQuery.matches) {
        colorScheme = "dark";
      }

      localStorage.setItem("theme", colorScheme);
    }

    return colorScheme;
  };

  // Detect color theme and set the relevant value
  window.onload = () => {
    document.documentElement.setAttribute("data-theme", getColorScheme());
  };
}());
