document.addEventListener("DOMContentLoaded", function () {
  var tabs = document.querySelectorAll(".install-tab");
  var panels = document.querySelectorAll(".install-panel");

  tabs.forEach(function (tab) {
    tab.addEventListener("click", function () {
      var target = tab.getAttribute("data-tab");

      tabs.forEach(function (t) {
        t.classList.remove("active");
        t.setAttribute("aria-selected", "false");
      });
      panels.forEach(function (p) {
        p.classList.remove("active");
      });

      tab.classList.add("active");
      tab.setAttribute("aria-selected", "true");
      document.querySelector('[data-panel="' + target + '"]').classList.add("active");
    });
  });
});
