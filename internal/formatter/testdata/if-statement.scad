if (part == "adapter") {
  adapter();
} else if (part == "back_weight") {
  back_weight();
} else if (part == "retaining_clip") {
  retaining_clip();
} else if (part == "side_weight") {
  side_weight();
  translate([0,20,0]) mirror([1,0,0]) side_weight();
}