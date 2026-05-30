// Mount asciinema-player on any <div data-cast="…"> container.
//
// Re-runs on every Material instant-navigation (via the document$ observable)
// so the player still initialises when arriving on a page client-side, and
// guards against double-mounting the same element.
function initAsciinemaPlayers() {
  if (typeof AsciinemaPlayer === "undefined") {
    // Player bundle not ready yet (CDN still loading) — retry shortly.
    window.setTimeout(initAsciinemaPlayers, 100);
    return;
  }
  document.querySelectorAll("div[data-cast]").forEach(function (el) {
    if (el.dataset.mounted) return;
    el.dataset.mounted = "1";
    AsciinemaPlayer.create(el.dataset.cast, el, {
      // Show the t=2:00 frame (the lazygit-in-tmux view) as the poster.
      poster: "npt:2:00",
      // Collapse long idle pauses in the recording to 2s.
      idleTimeLimit: 2,
      fit: "width",
    });
  });
}

if (typeof document$ !== "undefined" && document$.subscribe) {
  document$.subscribe(initAsciinemaPlayers);
} else {
  document.addEventListener("DOMContentLoaded", initAsciinemaPlayers);
}
