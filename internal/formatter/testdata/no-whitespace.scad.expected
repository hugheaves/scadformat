include <gripper_common.scad>
$fn = 64;
servo_width = 12.0;
servo_length = 23;
servo_tab_offset = 4.1;
servo_tab_length = 5;
servo_top_height = 3.4;
screw_hole_offset = 2.5;
screw_head_height = 3;
screw_head_diameter = 3.8;
bottom_mounting_tab_thickness = 6;
top_mounting_tab_thickness = 2;
mounting_tab_width = 6;
front_spindle_gap = 1;
screw_hole_diameter = 2.5;
standoff_diameter = 6;
top_thickness = servo_top_height - spindle_plate_thickness;
core_length = gear_diameter + gear_tooth_diameter + 5;
core_length_center = core_length / 2;
core_width = gear_tooth_diameter;
core_width_center = core_width / 2;
core_height = servo_tab_offset + top_thickness;
core_height_center = core_height / 2;
center_height = spindle_plate_thickness * 2 + spindle_height;
overall_height = center_height + core_height * 2;
standoff_height = center_height + core_height - screw_head_height - 2;
total_spindle_height = center_height + spindle_insertion_depth;
echo("overall_height=", overall_height);
echo("standoff_height=", standoff_height);
echo("total_spindle_height=", total_spindle_height);
back_left_spindle_length_offset = servo_tab_length + servo_width / 2 + gear_diameter;
back_right_spindle_length_offset = servo_tab_length + servo_width / 2;
front_spindle_plate_diameter = spindle_diameter + 4;
front_spindle_body_width_center = core_width + front_spindle_plate_diameter / 2;
front_spindle_body_length = front_spindle_plate_diameter * 2 + front_spindle_gap;
front_spindle_body_width = front_spindle_plate_diameter + front_spindle_gap * 2;
front_spindle_width_offset = front_spindle_body_width_center + front_spindle_gap * 2;
echo("core_height =", core_height);
echo("core_width =", core_width);
echo("core_length =", core_length);
echo("spindle width distance = ", core_width_center - front_spindle_width_offset);
module top_screw_hole() {
  translate([0, 0, core_height - screw_head_height])
    cylinder(d = screw_head_diameter, h = 100);
  cylinder(d = screw_hole_diameter, h = 100, $fn = 16);
}
module bottom_screw_hole() {
  cylinder(d = screw_hole_diameter, h = 100, $fn = 16);
  cylinder(d = screw_head_diameter, h = screw_head_height, $fn = 16);
}
module mounting_tab_screw_hole() {
  cylinder(d = screw_hole_diameter, h = bottom_mounting_tab_thickness + 0.0001, $fn = 16);
}
module top_mounting_tabs() {
  translate([0, -mounting_tab_offset - bottom_mounting_tab_thickness, 0]) {
    difference() {
      union() {
        translate([-mounting_tab_width, 0, 0])
          cube([core_length + mounting_tab_width * 2, bottom_mounting_tab_thickness, core_height]);
        cube([core_length, bottom_mounting_tab_thickness + mounting_tab_offset, core_height]);
      }
      translate([-mounting_tab_width, 0, 0])
        cube([mounting_tab_width, bottom_mounting_tab_thickness - top_mounting_tab_thickness, core_height]);
      translate([core_length, 0, 0])
        cube([mounting_tab_width, bottom_mounting_tab_thickness - top_mounting_tab_thickness, core_height]);
      mirror([0, 1, -1]) {
        translate([0, mounting_tab_width / 2, 0]) {
          translate([-mounting_tab_width / 2, 0, 0]) {
            mounting_tab_screw_hole();
            translate([0, overall_height - mounting_tab_width, 0])
              mounting_tab_screw_hole();
          }
          translate([core_length + mounting_tab_width / 2, 0, 0]) {
            mounting_tab_screw_hole();
            translate([0, overall_height - mounting_tab_width, 0])
              mounting_tab_screw_hole();
          }
        }
      }
    }
  }
}
module servo_cutout() {
  translate([-servo_width / 2, 0, 0]) {
    translate([-screw_hole_offset, 0, 0]) {
      top_screw_hole();
    }
    translate([servo_length + screw_hole_offset, 0, 0]) {
      top_screw_hole();
    }
  }
  translate([-servo_width / 2, -servo_width / 2, 0]) {
    cube([servo_length, servo_width, servo_tab_offset]);
  }
  hull() {
    translate([0, 0, 0]) {
      cylinder(d = servo_width, h = 100);
      translate([7, 0, 0])
        cylinder(d = 6, h = 100, $fn = 16);
    }
  }
}
module bottom_mounting_tabs() {
  difference() {
    union() {
      translate([-mounting_tab_width, 0, 0]) {
        cube([mounting_tab_width, bottom_mounting_tab_thickness, overall_height - core_height]);
        cube([mounting_tab_width, bottom_mounting_tab_thickness - top_mounting_tab_thickness, overall_height]);
      }
      translate([core_length, 0, 0]) {
        cube([mounting_tab_width, bottom_mounting_tab_thickness, overall_height - core_height]);
        cube([mounting_tab_width, bottom_mounting_tab_thickness - top_mounting_tab_thickness, overall_height]);
      }
    }
    mirror([0, 1, -1]) {
      translate([0, mounting_tab_width / 2, 0]) {
        translate([-mounting_tab_width / 2, 0, 0]) {
          mounting_tab_screw_hole();
          translate([0, overall_height - mounting_tab_width, 0])
            mounting_tab_screw_hole();
        }
        translate([core_length + mounting_tab_width / 2, 0, 0]) {
          mounting_tab_screw_hole();
          translate([0, overall_height - mounting_tab_width, 0])
            mounting_tab_screw_hole();
        }
      }
    }
  }
}
module bottom_spindle(plate_diameter) {
  cylinder(d = spindle_diameter, h = total_spindle_height);
  cylinder(d = plate_diameter, h = spindle_plate_thickness);
}
module core() {
  cube([core_length, core_width, core_height]);
  translate([0, -standoff_diameter / 2 - 1, 0])
    cube([core_length, core_width, core_height]);
  translate([(core_length - front_spindle_body_length) / 2, core_width, 0])
    cube([front_spindle_body_length, front_spindle_body_width, core_height]);
}
module bottom() {
  difference() {
    union() {
      core();
      bottom_mounting_tabs();
      translate([core_length_center + gear_diameter / 2, core_width_center, core_height - 0.0001]) {
        bottom_spindle(gear_tooth_diameter);
      }
      translate([core_length_center - gear_diameter / 2, core_width_center, core_height - 0.0001]) {
        bottom_spindle(gear_tooth_diameter);
      }
      translate([0, front_spindle_width_offset, core_height - 0.0001]) {
        translate([core_length_center + front_spindle_plate_diameter / 2 + front_spindle_gap / 2, 0, 0]) {
          bottom_spindle(front_spindle_plate_diameter);
        }
        translate([core_length_center - front_spindle_plate_diameter / 2 - front_spindle_gap / 2, 0, 0]) {
          bottom_spindle(front_spindle_plate_diameter);
        }
      }
      translate([core_length_center, core_width + 0.5, core_height - 0.0001])
        cylinder(d = standoff_diameter, h = standoff_height);
      standoff_offset = sin(45) * (gear_tooth_diameter + standoff_diameter) / 2;
      translate([core_length_center - gear_diameter / 2 - standoff_offset, core_width_center - standoff_offset, core_height - 0.0001])
        cylinder(d = standoff_diameter, h = standoff_height);
    }
    translate([core_length_center, core_width_center - 12.5, 0])
      bottom_screw_hole();
    translate([core_length_center, core_width_center + 12.5, 0])
      bottom_screw_hole();
    translate([core_length_center - gear_diameter / 2, core_width_center, 0])
      servo_cutout();
  }
}
module top() {
  difference() {
    union() {
      body();
      top_mounting_tab();
      translate([0, 0, -spindle_plate_thickness]) {
        translate([back_left_spindle_length_offset, core_width_center, 0]) {
          cylinder(d = core_width, h = spindle_plate_thickness);
        }
        translate([back_right_spindle_length_offset, core_width_center, 0]) {
          cylinder(d = core_width, h = spindle_plate_thickness);
        }
        translate([0, front_spindle_width_offset, 0]) {
          translate([front_left_spindle_length_offset, 0, 0]) {
            cylinder(d = front_spindle_plate_diameter, h = spindle_plate_thickness);
          }
          translate([front_right_spindle_length_offset, 0, 0]) {
            cylinder(d = front_spindle_plate_diameter, h = spindle_plate_thickness);
          }
        }
      }
    }
    spindle_insert_height = spindle_insertion_depth + spindle_plate_thickness;
    translate([0, 0, -spindle_plate_thickness]) {
      translate([back_left_spindle_length_offset, core_width_center, 0]) {
        cylinder(d = spindle_diameter + horizontal_gap, h = spindle_insert_height);
      }
      translate([back_right_spindle_length_offset, core_width_center, 0]) {
        cylinder(d = servo_gear_spindle_diameter + horizontal_gap, h = spindle_insert_height);
      }
      translate([0, front_spindle_width_offset, 0]) {
        translate([front_left_spindle_length_offset, 0, 0]) {
          cylinder(d = spindle_diameter + horizontal_gap, h = spindle_insert_height);
        }
        translate([front_right_spindle_length_offset, 0, 0]) {
          cylinder(d = spindle_diameter + horizontal_gap, h = spindle_insert_height);
        }
      }
    }
    translate([core_length_center, core_width_center + 12.5, 0]) {
      cylinder(d = standoff_diameter + horizontal_gap, h = core_height - screw_head_height - 2);
      top_screw_hole();
    }
    translate([core_length_center, core_width_center - 12.5, 0]) {
      cylinder(d = standoff_diameter + horizontal_gap, h = core_height - screw_head_height - 2);
      top_screw_hole();
    }
    translate([core_length_center, core_width_center, 0]) {
      cylinder(d = standoff_diameter + horizontal_gap, h = core_height - screw_head_height - 2);
      top_screw_hole();
    }
  }
}
top();
rotate([180, 0, 0])
  bottom();
