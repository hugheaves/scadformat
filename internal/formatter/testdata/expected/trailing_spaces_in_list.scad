include <BOSL2/std.scad>
include <BOSL2/screws.scad>

$fn = 15;

module qwe(skrew_len) {
  spec = [
    ["system", "ISO"],
    ["type", "screw_info"],
    ["pitch", 2.3],
    ["head", "flat"],
    ["head_size", 5.8],
    ["head_size_sharp", 5.8],
    ["head_angle", 74],
    ["diameter", 1.7],
    ["length", skrew_len]
  ];

  screw_hole(spec, anchor = TOP, counterbore = 0.5);
}

qwe(10);
