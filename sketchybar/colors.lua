return {
  black = 0xff181819,
  white = 0xffe2e2e3,
  orange = 0xfff39660,
  magenta = 0xffb39df3,
  grey = 0xff7f8490,
  transparent = 0x00000000,

  bar = {
    bg = 0x00000000,
    border = 0x00000000,
  },
  popup = {
    bg = 0xc02c2e34,
    border = 0xff7f8490,
  },
  bg1 = 0xff363944,
  bg2 = 0xff414550,

  base = 0xff24273a,
  mantle = 0xff1e2030,
  crust = 0xff181926,
  text = 0xffcad3f5,
  subtext0 = 0xffb8c0e0,
  subtext1 = 0xffa5adcb,
  surface0 = 0xff363a4f,
  surface1 = 0xff494d64,
  surface2 = 0xff5b6078,
  overlay0 = 0xff6e738d,
  overlay1 = 0xff8087a2,
  overlay2 = 0xff939ab7,
  blue = 0xff8aadf4,
  lavender = 0xffb7bdf8,
  sapphire = 0xff7dc4e4,
  sky = 0xff91d7e3,
  teal = 0xff8bd5ca,
  green = 0xffa6da95,
  yellow = 0xffeed49f,
  peach = 0xfff5a97f,
  maroon = 0xffee99a0,
  red = 0xffed8796,
  mauve = 0xffc6a0f6,
  pink = 0xfff5bde6,
  flamingo = 0xfff0c6c6,
  rosewater = 0xfff4dbd6,

  rainbow = {
    0xffff007c,
    0xffc53b53,
    0xffff757f,
    0xff41a6b5,
    0xff4fd6be,
    0xffc3e88d,
    0xffffc777,
    0xff9d7cd8,
    0xffff9e64,
    0xffbb9af7,
    0xff7dcfff,
    0xff7aa2f7,
  },

  with_alpha = function(color, alpha)
    if alpha > 1.0 or alpha < 0.0 then
      return color
    end
    return (color & 0x00ffffff) | (math.floor(alpha * 255.0) << 24)
  end,
}
