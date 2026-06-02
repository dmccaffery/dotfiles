// Mount asciinema-player on any <div data-cast="…"> container.
//
// Re-runs on every Material instant-navigation (via the document$ observable)
// so the player still initialises when arriving on a page client-side, and
// guards against double-mounting the same element.
// The recording's terminal font is Iosevka Nerd Font (see extras.css). asciinema-player
// measures glyph metrics when it mounts, so the web font must be loaded first — otherwise the
// grid is sized against the fallback and the Nerd Font glyphs misalign once they swap in.
async function ensureDemoFontLoaded() {
  if (!document.fonts || !document.fonts.load) return Promise.resolve();
  var family = '"Iosevka NF"';
  try {
    return await Promise.all([
      document.fonts.load("1em " + family),
      document.fonts.load("bold 1em " + family),
      document.fonts.load("italic 1em " + family),
      document.fonts.load("italic bold 1em " + family),
    ]);
  } catch {}
}

function initAsciinemaPlayers() {
  if (typeof AsciinemaPlayer === "undefined") {
    // Player bundle not ready yet (CDN still loading) — retry shortly.
    window.setTimeout(initAsciinemaPlayers, 100);
    return;
  }
  ensureDemoFontLoaded().then(mountAsciinemaPlayers);
}

function mountAsciinemaPlayers() {
  document.querySelectorAll("div[data-cast]").forEach(function (el) {
    if (el.dataset.mounted) return;
    el.dataset.mounted = "1";
    AsciinemaPlayer.create(el.dataset.cast, el, {
      // Render with the terminal's Nerd Font (loaded via @font-face in extras.css). The
      // player measures its own glyph metrics, so the font must be named here — a CSS
      // font-family override alone leaves it on the default font, tofu-ing the glyphs.
      terminalFontFamily: '"Iosevka NF", monospace',
      // Show the t=2:00 frame (the lazygit-in-tmux view) as the poster.
      poster: "npt:2:16",
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
