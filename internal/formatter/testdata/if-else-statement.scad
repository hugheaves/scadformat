if (part == "adapter") {
  adapter();
} else {
  side_weight();
  translate([0,20,0]) mirror([1,0,0]) side_weight();
}