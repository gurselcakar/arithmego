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

  var copyIcon = '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>';
  var checkIcon = '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>';

  document.querySelectorAll(".copy-btn").forEach(function (btn) {
    btn.addEventListener("click", function () {
      var text = btn.getAttribute("data-copy");
      navigator.clipboard.writeText(text).then(function () {
        btn.innerHTML = checkIcon;
        setTimeout(function () {
          btn.innerHTML = copyIcon;
        }, 2000);
      });
    });
  });

  // Autoplay videos when they enter the viewport (mobile Safari fix)
  var videos = document.querySelectorAll("video");
  if ("IntersectionObserver" in window) {
    var observer = new IntersectionObserver(function (entries) {
      entries.forEach(function (entry) {
        if (entry.isIntersecting) {
          entry.target.play();
        }
      });
    }, { threshold: 0.25 });
    videos.forEach(function (video) {
      observer.observe(video);
    });
  } else {
    videos.forEach(function (video) {
      video.play();
    });
  }
});
